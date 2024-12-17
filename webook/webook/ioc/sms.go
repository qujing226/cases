package ioc

import (
	"webook/webook/internal/service/sms"
	"webook/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 使用内存实现
	return memory.NewService()
}
