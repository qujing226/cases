package MQ_Sarama_test

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	//cfg.Producer.Partitioner = MQ-Sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("Hello, this is a message"),
		// 会在生产者与消费者之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		// 只作用于发送过程
		Metadata: "this is metadata",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)
	msgCh := producer.Input()
	msgCh <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Key:   sarama.StringEncoder("key"),
		Value: sarama.StringEncoder("Hello, this is a message"),
	}
	errCh := producer.Errors()
	succCh := producer.Successes()

	select {
	case err := <-errCh:
		t.Log("send error:", err.Err)
	case <-succCh:
		t.Log("send success:")
	}
}

type JSONEncoder struct {
	Data any
}

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addrs,
		"test_group_consumer", cfg)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*10, func() {
		cancel()
	})
	err = consumer.Consume(ctx,
		[]string{"test_topic"}, &testConsumerGroupHandler{})
	// 消费结束
	t.Log(err)
}

type testConsumerGroupHandler struct {
}

func (t *testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("Setup")
	// 更改偏移量，建议走离线渠道
	// topic => 偏移量  消费的位置
	//partitions := session.Claims()["test_topic"]
	//for _,part := range partitions {
	//	session.ResetOffset("test_topic",part,
	//		MQ-Sarama.OffsetOldest,"")
	//}
	return nil
}

func (t *testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

// 异步消费
func (t *testConsumerGroupHandler) ConsumeClaim(
	// session 代表会话（建立链接-彻底断掉）
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		done := false

		var eg errgroup.Group
		var lastMsg *sarama.ConsumerMessage
		for i := 0; i < batchSize && !done; i++ {
			select {
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					// 消费者被关闭了
					return nil
				}
				eg.Go(func() error {
					// 在此处消费
					time.Sleep(time.Second)
					log.Println("msg:", string(msg.Value))
					return nil
				})
				if i == batchSize-1 {
					lastMsg = msg
					break
				}
			case <-ctx.Done():
				done = true
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			// 记录日志

			continue
		}
		if lastMsg != nil {
			log.Println("Mark -> lastMsg:", string(lastMsg.Value))
			session.MarkMessage(lastMsg, "pass")
		}
	}
	return nil
}

// 同步消费
func (t *testConsumerGroupHandler) ConsumeClaimV1(
	// session 代表会话（建立链接-彻底断掉）
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var bizMsg MyBizMsg
		err := json.Unmarshal(msg.Value, &bizMsg)
		if err != nil {
			// 重试 + 记录日志

		}
		println(string(msg.Value))
		session.MarkMessage(msg, "pass")
	}
	// msgs 被人关了，需要退出消费逻辑
	return nil
}

type MyBizMsg struct {
	TraceId string `json:"trace_id"`
	Data    string `json:"data"`
}
