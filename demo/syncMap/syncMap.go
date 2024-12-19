package syncMap

import (
	"fmt"
	"sync"
	"websocket/demo/bloomFilter"
)

// 1. 使用 读写锁 实现
// 2. 加入布隆过滤器      需指定预估map容量和布隆过滤器正确率
// 3. 使用 channel 实现

/*
	Case 1：map读写锁
*/

// Map Case 1：map读写锁
type Map[k comparable, v any] struct {
	mu sync.RWMutex
	m  map[k]v
}

func NewMap[k comparable, v any]() *Map[k, v] {
	return &Map[k, v]{
		m: make(map[k]v),
	}
}

func (m *Map[k, v]) Load(key k) (v, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.m[key]
	if ok {
		return v, true
	}
	return v, false
}

func (m *Map[k, v]) Store(key k, val v) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = val
}

func (m *Map[k, v]) Delete(key k) {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()
}

/*
	Case 2: 加入布隆过滤器
	初始化过程中的size用来标识map的预估数据量
	只有在读多而且很多数据无法查到的情况下效率才高
	优化点：
		通过算法使得读操作的goroutine无需获取锁即可判断数据是否不存在
		查不存在是准确的，查存在是不准确的
*/

// BFMap 布隆过滤器
type BFMap[k comparable, v any] struct {
	mu sync.RWMutex
	bf *bloomFilter.BloomFilter
	m  map[k]v
}

func NewBFMap[k comparable, v any](size int, rate float64) *BFMap[k, v] {
	return &BFMap[k, v]{
		bf: bloomFilter.NewBloomFilter(size, rate),
		m:  make(map[k]v),
	}
}

func (m *BFMap[k, v]) Load(key k) (val v, flag bool) {
	if !m.bf.Test(fmt.Sprintf("%v", key)) {
		return val, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[key], true
}

func (m *BFMap[k, v]) Store(key k, val v) {
	m.mu.Lock()
	m.m[key] = val
	m.mu.Unlock()
	m.bf.Add(fmt.Sprintf("%v", key))
}

func (m *BFMap[k, v]) Delete(key k) {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()
}

/*
	Case 3: channel实现
	1. 创建一个channel，用于接收操作
	2. 启动一个goroutine，用于接收操作，并根据操作类型进行操作
	3. 使用channel进行读写操作
	result[v] 类型用来Load方法中标识key是否存在
*/

// Operation 标识操作
type Operation string

const (
	Load   Operation = "load"
	Store  Operation = "store"
	Delete Operation = "delete"
)

type Operate[k comparable, v any] struct {
	Operation
	key    k
	val    v
	Return chan result[v]
}

type result[v any] struct {
	val v
	ok  bool
}

type ChMap[k comparable, v any] struct {
	ch   chan Operate[k, v]
	m    map[k]v
	done chan struct{}
}

func NewChMap[k comparable, v any]() *ChMap[k, v] {
	m := &ChMap[k, v]{
		ch:   make(chan Operate[k, v]),
		m:    make(map[k]v),
		done: make(chan struct{}),
	}
	go func() {
		for {
			var operate Operate[k, v]
			select {
			case operate = <-m.ch:
				switch operate.Operation {
				case Load:
					val, ok := m.m[operate.key]
					operate.Return <- result[v]{val, ok}
				case Store:
					m.m[operate.key] = operate.val
				case Delete:
					delete(m.m, operate.key)
				}
			case <-m.done:
				close(m.ch)
				return
			}
		}

	}()
	return m
}

func (m *ChMap[k, v]) Load(key k) (val v, flag bool) {
	operate := Operate[k, v]{
		Operation: Load,
		key:       key,
		Return:    make(chan result[v]),
	}
	m.ch <- operate         // 将数据包装在Operate中进行传送
	res := <-operate.Return // 等待回写
	return res.val, res.ok
}
func (m *ChMap[k, v]) Store(key k, val v) {
	operate := Operate[k, v]{
		Operation: Store,
		key:       key,
		val:       val,
	}
	m.ch <- operate
}

func (m *ChMap[k, v]) Delete(key k) {
	operate := Operate[k, v]{
		Operation: Delete,
		key:       key,
	}
	m.ch <- operate
}
func (m *ChMap[k, v]) Close() {
	close(m.done)
}
