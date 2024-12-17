package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码次数过多")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数过多")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCacher interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) error
	Key(biz, phone string) string
}

type CodeCacheRedis struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCacher {
	return &CodeCacheRedis{
		client: client,
	}
}

func (c *CodeCacheRedis) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
	case -1:
		return errors.New("internal error")
	case -2:
		return ErrCodeSendTooMany
	}
	return nil
}
func (c *CodeCacheRedis) Verify(ctx context.Context, biz, phone, inputCode string) error {
	code, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, inputCode).Int()
	if err != nil {
		fmt.Println(err)
		return err
	}
	switch code {
	case 0:
	case -1:
		return ErrCodeVerifyTooManyTimes
	case -2:
		return errors.New("internal error")
	}
	return nil
}

func (c *CodeCacheRedis) Key(biz, phone string) string {
	return "phone_code:" + biz + ":" + phone
}
