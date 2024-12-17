package service

import (
	"context"
	"webook/webook/internal/repository"
	"webook/webook/internal/service/sms"
)

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeServicer interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) error
}

type CodeService struct {
	repo   repository.CodeRepositoryer
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepositoryer, smsSvc sms.Service) CodeServicer {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := sms.GenerateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = svc.smsSvc.Send(ctx, biz, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) error {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
