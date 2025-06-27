package DistributedLock

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math"
	"time"
)

const (
	lockKey       = "my_distributed_lock"
	lockValve     = "locked"
	lockDuration  = 10 * time.Second
	renewInterval = 3 * time.Second
)

// 如果需要实现可重入，可以使用计数器，每次获取锁时，计数器加1，释放锁时，计数器减1，计数器为0时，删除锁。

type SpinConfig struct {
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 基础重试间隔
	BackoffFactor float64       // 退避因子（指数退避）
	TimeOut       time.Duration // 总超时时间
}

func NewSpinConfig(maxRetries int, retryInterval time.Duration, backoffFactor float64, timeOut time.Duration) *SpinConfig {
	return &SpinConfig{
		MaxRetries:    maxRetries,
		RetryInterval: retryInterval,
		BackoffFactor: backoffFactor,
		TimeOut:       timeOut,
	}
}

type DistributedLock struct {
	client       redis.Cmdable
	ctx          context.Context
	cancelRenew  context.CancelFunc
	lockAcquired bool

	spinConfig     SpinConfig    // 自旋
	lastRetryDelay time.Duration // 当前退避时间
}

func NewDistributedLock(client redis.Cmdable, spinCfg *SpinConfig) *DistributedLock {
	return &DistributedLock{
		client:     client,
		ctx:        context.Background(),
		spinConfig: *spinCfg,
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

// AcquireWithSpin 自旋锁
func (dl *DistributedLock) AcquireWithSpin(ctx context.Context) (bool, error) {
	// 创建超时控制
	timeoutCtx, cancel := context.WithTimeout(ctx, dl.spinConfig.TimeOut)
	defer cancel()

	attempt := 0
	for {
		select {
		case <-timeoutCtx.Done():
			return false, timeoutCtx.Err()
		default:
			acquired, err := dl.tryAcquire()
			if err != nil {
				return false, err
			}
			if acquired {
				return true, nil
			}
			// 计算下次重试的时间间隔
			delay := dl.calculateDelay(attempt)
			timer := time.NewTimer(delay)
			select {
			case <-timeoutCtx.Done():
				return false, timeoutCtx.Err()
			case <-timer.C:
				attempt++
				if dl.spinConfig.MaxRetries > 0 && attempt >= dl.spinConfig.MaxRetries {
					return false, fmt.Errorf("max retries exceeded")
				}
			}
			timer.Stop()
		}
	}
}

func (dl *DistributedLock) tryAcquire() (bool, error) {
	result, err := dl.client.SetNX(dl.ctx, lockKey, lockValve, lockDuration).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return result, err
}

// calculateDelay 智能退避算法
func (dl *DistributedLock) calculateDelay(attempt int) time.Duration {
	if dl.spinConfig.BackoffFactor <= 1 {
		return dl.spinConfig.RetryInterval
	}
	maxDelay := dl.spinConfig.TimeOut / 2
	delay := time.Duration(float64(dl.spinConfig.RetryInterval) * math.Pow(dl.spinConfig.BackoffFactor, float64(attempt)))
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}

// AcquireWithSpinAndRenew 看门狗 + 自旋
func (dl *DistributedLock) AcquireWithSpinAndRenew(ctx context.Context) (bool, error) {
	acquired, err := dl.AcquireWithSpin(ctx)
	if !acquired || err != nil {
		return false, err
	}
	// 启动看们狗
	renewCtx, cancel := context.WithCancel(context.Background())
	dl.cancelRenew = cancel
	go dl.watchdog(renewCtx)

	return true, nil
}
