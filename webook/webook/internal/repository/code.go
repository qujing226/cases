package repository

import (
	"context"
	"webook/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepositoryer interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) error
}

type CodeRepository struct {
	cache cache.CodeCacher
}

func NewCodeRepository(client cache.CodeCacher) CodeRepositoryer {
	return &CodeRepository{
		cache: client,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Verify(ctx, biz, phone, code)
}
