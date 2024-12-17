package ioc

import "github.com/redis/go-redis/v9"

func InitRDB() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	return rdb
}
