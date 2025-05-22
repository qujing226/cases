package godis_distributed_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	lock := NewDistributedLock(rdb, NewSpinConfig(
		5, 500*time.Microsecond, 1.5, 10*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if ok, err := lock.AcquireWithSpinAndRenew(ctx); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatal("lock failed")
		return
	}
	defer lock.Release()

}
