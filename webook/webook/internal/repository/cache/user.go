package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/webook/internal/domain"
)

var (
	ErrKeyNotExist = redis.Nil
)

type UserCacher interface {
	Set(ctx context.Context, u domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
}

type UserCache struct {
	// 传单机redis可以 传cluster的redis也可以
	client     redis.Cmdable
	expiration time.Duration
}

// NewUserCache
// A用到B B一定是接口  B一定是A的字段   A绝对不初始化B而是依赖注入
func NewUserCache(cmd redis.Cmdable) UserCacher {
	return &UserCache{
		client:     cmd,
		expiration: time.Minute * 15,
	}
}

func (uc *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	val, err := uc.client.Get(ctx, uc.Key(id)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	if err := json.Unmarshal(val, &u); err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (uc *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return uc.client.Set(ctx, uc.Key(u.Id), val, uc.expiration).Err()
}

func (uc *UserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
