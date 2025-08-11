package Pool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPool(t *testing.T) {
	// 初始化一个实例
	bufferpool := &sync.Pool{
		New: func() interface{} {
			println("create new instance")
			return make([]byte, 1024)
		},
	}

	buffer := bufferpool.Get()
	bufferpool.Put(buffer)

}

// 这是一个对象池和非对象池初始化的一个对比
func TestCreateNum(t *testing.T) {
	var numCalcsCreated int32

	bufferpool := getBufferFromPool()
	var createbuffer = func() any {
		atomic.AddInt32(&numCalcsCreated, 1)
		buffer := make([]byte, 1024)
		return &buffer
	}

	bufferpool = &sync.Pool{
		New: createbuffer,
	}

	numWorkers := 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			//buffer := bufferpool.Get()
			buffer := createbuffer()
			_ = buffer.(*[]byte)
			defer bufferpool.Put(buffer)
		}()
	}
	wg.Wait()
	fmt.Printf("%d buffer objects were created.\n", numCalcsCreated)
}

func getBufferFromPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			println("create new instance")
			return make([]byte, 1024)
		},
	}
}
