package godis_distributed_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	lockKey       = "my_distributed_lock"
	lockValve     = "locked"
	lockDuration  = 10 * time.Second
	renewInterval = 3 * time.Second
)

type DistributedLock struct {
	client       redis.Cmdable
	ctx          context.Context
	cancelRenew  context.CancelFunc
	lockAcquired bool
}

func NewDistributedLock(client redis.Cmdable) *DistributedLock {
	return &DistributedLock{
		client: client,
		ctx:    context.Background(),
	}
}

// Acquire 获取锁 （带看门狗）
func (dl *DistributedLock) Acquire() (bool, error) {
	if dl.lockAcquired {
		return false, nil
	}
	ok, err := dl.client.SetNX(dl.ctx, lockKey, lockValve, lockDuration).Result()
	if err != nil || !ok {
		return false, err
	}

	var renewCtx context.Context
	renewCtx, dl.cancelRenew = context.WithCancel(context.Background())

	go dl.watchdog(renewCtx)

	dl.lockAcquired = true
	return true, nil
}

func (dl *DistributedLock) watchdog(ctx context.Context) {
	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// Release 时执行cancelRenew 会导致 renewCtx 被取消(Done)
			return
		case <-ticker.C:
			script := redis.NewScript(`
				if redis.call("GET",KEYS[1]) == ARGV[1] then 
					return redis.call("PEXPIRE",KEYS[1],ARGV[2])
				else
					return 0
				end
			`)
			result, err := script.Run(dl.ctx, dl.client, []string{lockKey}, lockValve, int(lockDuration/time.Millisecond)).Int()
			if err != nil || result == 0 {
				dl.Release()
				return
			}
		}
	}
}

func (dl *DistributedLock) Release() error {
	if !dl.lockAcquired {
		return nil
	}
	// 停止看门狗
	if dl.cancelRenew != nil {
		dl.cancelRenew()
	}
	// lua脚本保证原子性删除
	script := redis.NewScript(`
		if redis.call("GET",KEYS[1]) == ARGV[1] then 
			return redis.call("DEL",KEYS[1])
		else 
			return 0
		end
	`)
	_, err := script.Run(dl.ctx, dl.client, []string{lockKey}, lockValve).Int()
	dl.lockAcquired = false
	return err
}
