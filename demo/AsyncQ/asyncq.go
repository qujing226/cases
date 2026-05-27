package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// 替换为你的 EloqKV 地址
const redisAddr = "127.0.0.1:26279"

func main() {
	// 1. 连接测试
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// 2. 核心功能测试：延迟任务 (依赖 ZSET 和 Lua)
	// 如果 EloqKV 的 Lua 脚本处理有问题，这里会报错
	task := asynq.NewTask("email:send", []byte("hello"))
	info, err := client.Enqueue(task, asynq.ProcessIn(2*time.Second))
	if err != nil {
		log.Fatalf("❌ 入队失败 (可能是 Lua/ZSET 问题): %v", err)
	}
	fmt.Printf("✅ 任务入队成功 ID: %s (将在2秒后处理)\n", info.ID)

	// 3. 启动 Worker 消费 (依赖 LIST 和 BRPOP)
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 1},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:send", func(ctx context.Context, t *asynq.Task) error {
		fmt.Printf("✅ 成功消费任务: %s\n", string(t.Payload()))
		return nil
	})

	fmt.Println("等待 Worker 拉取任务...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("❌ Worker 运行失败: %v", err)
	}

	// 在 main 中：
	task = asynq.NewTask("task:fail", nil)
	// 设置：最多重试 3 次，每次失败后立刻重试（这里为了测试方便）
	_, err = client.Enqueue(task, asynq.MaxRetry(3))
}

// ... client初始化同前 ...

// 模拟一个必然失败的任务
func HandleFailTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println("🔥 模拟任务处理失败，抛出错误...")
	return fmt.Errorf("intentional error")
}

// 启动 Worker，观察日志：
// 你应该看到：
// 1. 任务失败
// 2. Asynq 将其放入 Retry 队列 (EloqKV 必须支持 ZADD/ZPOPMIN 等操作的原子性)
// 3. 过一会儿再次被捞起重试
