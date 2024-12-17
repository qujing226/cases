package sms

import (
	"context"
	"fmt"
	"math/rand/v2"
)

type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
	Verify(ctx context.Context, biz string, code string, number string) (bool, error)
}

type Client struct {
	tpl string
}

func GenerateCode() string {
	number := rand.N[int](1000000)
	return fmt.Sprintf("%06d", number)
}
