package sarama_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	//cfg.Producer.Partitioner = sarama.NewHashPartitioner
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
	//		sarama.OffsetOldest,"")
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

func TestGenerateBitcoinAddress(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "test",
			input:  "0fe57f1532728ea9f1891b0bce90ba3f9c3c64f0cda0e439e9c2fa56553014b9",
			output: "15gQHFAHyvSAMZZ7tZLfbNyngtJ6fQe3Hm",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := GenerateBitcoinAddress(tc.input)
			assert.Equal(t, tc.output, res)
		})
	}
}

func GenerateBitcoinAddress(privateKeyHex string) string {
	// 1. 解码私钥
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	// 2. 生成压缩公钥
	pubKey := privateKey.PubKey()
	x := pubKey.X()
	y := pubKey.Y()

	var prefix byte
	if y.Bit(0) == 0 { // 根据 y 的奇偶性选择前缀
		prefix = 0x02
	} else {
		prefix = 0x03
	}

	// 将 x 填充到 32 字节
	xBytes := x.Bytes()
	paddedX := make([]byte, 32)
	copy(paddedX[32-len(xBytes):], xBytes) // ✅ 高位补零
	compressedPublicKey := append([]byte{prefix}, paddedX...)

	// 3. 计算 SHA-256 + RIPEMD-160
	sha256Hash := sha256.Sum256(compressedPublicKey)
	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash[:])
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// 4. 添加版本字节 0x00
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// 5. 计算校验和
	firstSHA := sha256.Sum256(versionedPayload)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	// 6. Base58 编码
	finalPayload := append(versionedPayload, checksum...)
	bitcoinAddress := base58.Encode(finalPayload)
	return bitcoinAddress
}

// GenerateBitcoinAddress generates a compressed Bitcoin P2PKH address from a private key.
func GenerateBitcoinAddressV1(privateKeyHex string) string {
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	// 2. 生成压缩公钥
	compressedPublicKey := privateKey.PubKey().SerializeCompressed()

	// 3. 生成 P2PKH 地址（Legacy 格式）
	addressPubKeyHash, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(compressedPublicKey), &chaincfg.MainNetParams)
	return addressPubKeyHash.EncodeAddress()
}
