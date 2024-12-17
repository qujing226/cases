package tencent

import (
	"webook/webook/internal/service/sms"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   client,
	}
}
