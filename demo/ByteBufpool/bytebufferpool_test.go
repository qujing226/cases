package ByteBufpool

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/valyala/bytebufferpool"
	"io"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	bb := bytebufferpool.Get()

	//bb.WriteString("first line\n")
	//bb.Write([]byte("second line\n"))
	//bb.B = append(bb.B, "third line\n"...)
	err := binary.Write(bb, binary.BigEndian, uint32(42))
	if err != nil {
		fmt.Println(err)
		return
	}
	var num uint32
	bytess := bytes.NewBuffer(bb.B)
	err = binary.Read(bytess, binary.BigEndian, &num)
	fmt.Printf("bytebuffer contents=%q\n", bb.B)
	fmt.Println(num)

	// It is safe to release byte buffer now, since it is
	// no longer used.
	bytebufferpool.Put(bb)
}

// ... existing code ...
// ... existing code ...

func TestBufferPerformanceComparison(t *testing.T) {
	const (
		iterations = 10000
		strSize    = 10 * 1024 // 10KB
	)

	// 准备测试数据
	shortBytes := []byte("example data")
	number := uint32(0x12345678)

	// 定义测试执行器
	type bufferTester struct {
		name         string
		createBuffer func() interface{}
		writeFunc    func(interface{})
	}

	testers := []bufferTester{
		// binary.Write with bytes.Buffer
		{
			name: "binary.Write/bytes.Buffer",
			createBuffer: func() interface{} {
				return new(bytes.Buffer)
			},
			writeFunc: func(buf interface{}) {
				_ = binary.Write(buf.(io.Writer), binary.BigEndian, number)
			},
		},
		{
			name: "binary.Write/bytes.Buffer",
			createBuffer: func() interface{} {
				return bytebufferpool.Get()
			},
			writeFunc: func(buf interface{}) {
				_ = binary.Write(buf.(io.Writer), binary.BigEndian, number)
			},
		},
		// bytes.Buffer direct write
		{
			name: "bytes.Buffer.Write",
			createBuffer: func() interface{} {
				return new(bytes.Buffer)
			},
			writeFunc: func(buf interface{}) {
				_, _ = buf.(*bytes.Buffer).Write(shortBytes)
			},
		},
		// bytebufferpool write
		{
			name: "bytebufferpool.Write",
			createBuffer: func() interface{} {
				return bytebufferpool.Get()
			},
			writeFunc: func(buf interface{}) {
				_, _ = buf.(*bytebufferpool.ByteBuffer).Write(shortBytes)
			},
		},
		// binary.Write with bytebufferpool
		{
			name: "binary.Write/bytebufferpool",
			createBuffer: func() interface{} {
				return bytebufferpool.Get()
			},
			writeFunc: func(buf interface{}) {
				_ = binary.Write(buf.(io.Writer), binary.BigEndian, number)
			},
		},
	}

	// 执行性能测试
	results := make([]struct {
		name       string
		totalTime  time.Duration
		avgTime    time.Duration
		allocCount int64
	}, len(testers))

	for i, tester := range testers {
		var totalAlloc int64
		start := time.Now()

		for j := 0; j < iterations; j++ {
			// 每次测试创建新缓冲区
			buf := tester.createBuffer()

			// 执行写入操作
			tester.writeFunc(buf)

			// 释放资源并统计内存分配
			if bb, ok := buf.(*bytebufferpool.ByteBuffer); ok {
				bytebufferpool.Put(bb)
			}
			// 记录内存分配情况
			// 实际测试应使用pprof进行精确内存统计
		}

		elapsed := time.Since(start)

		results[i] = struct {
			name       string
			totalTime  time.Duration
			avgTime    time.Duration
			allocCount int64
		}{
			name:       tester.name,
			totalTime:  elapsed,
			avgTime:    elapsed / time.Duration(iterations),
			allocCount: totalAlloc,
		}
	}

	// 输出对比结果
	t.Log("\n缓冲区性能对比分析:")
	t.Logf("%-30s | %-15s | %-15s | %-15s", "测试项", "总耗时", "单次耗时", "内存分配")
	t.Logf("%-30s-+-%-15s-+-%-15s-+-%-15s",
		strings.Repeat("-", 30), strings.Repeat("-", 15),
		strings.Repeat("-", 15), strings.Repeat("-", 15))

	for _, res := range results {
		t.Logf("%-30s | %-15s | %-15s | %d",
			res.name, res.totalTime, res.avgTime, res.allocCount)
	}
}

// 辅助函数：生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

// 辅助函数：生成随机字节切片
func generateRandomBytes(length int) []byte {
	return []byte(generateRandomString(length))
}
