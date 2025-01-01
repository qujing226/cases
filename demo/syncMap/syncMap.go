package syncMap

import (
	"fmt"
	"peninsula12/cases/demo/bloomFilter"
	"sync"
)

// 1. 使用 读写锁 实现
// 2. 加入布隆过滤器      需指定预估map容量和布隆过滤器正确率
// 3. 使用 channel 实现

/*
	Case 1：map读写锁
*/

// Map Case 1：map读写锁
type Map[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewMap 初始化map
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

// Load 获取数据
func (m *Map[K, V]) Load(key K) (val V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

// Store 存储数据
func (m *Map[K, V]) Store(key K, val V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = val
}

// Delete 删除数据
func (m *Map[K, V]) Delete(key K) {
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
type BFMap[K comparable, V any] struct {
	mu sync.RWMutex
	bf *bloomFilter.BloomFilter
	m  map[K]V
}

// NewBFMap 初始化布隆过滤器Map
func NewBFMap[k comparable, v any](size int, rate float64) *BFMap[k, v] {
	return &BFMap[k, v]{
		bf: bloomFilter.NewBloomFilter(size, rate),
		m:  make(map[k]v),
	}
}

// Load 加载数据
func (m *BFMap[K, V]) Load(key K) (val V, flag bool) {
	if !m.bf.Test(fmt.Sprintf("%v", key)) {
		return val, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[key], true
}

// Store 存储数据
func (m *BFMap[K, V]) Store(key K, val V) {
	m.mu.Lock()
	m.m[key] = val
	m.mu.Unlock()
	m.bf.Add(fmt.Sprintf("%v", key))
}

// Delete 删除数据
func (m *BFMap[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.m, key)
	m.mu.Unlock()
}

/*
	Case 3: channel实现
	1. 创建一个channel，用于接收操作
	2. 启动一个goroutine，用于接收操作，并根据操作类型进行操作
	3. 使用channel进行读写操作
	result[V] 类型用来Load方法中标识key是否存在
*/

type operation string

const (
	load  operation = "load"
	store operation = "store"
	del   operation = "del"
)

type operate[K comparable, V any] struct {
	op  operation
	key K
	val V
	ret chan result[V]
}

type result[V any] struct {
	val V
	ok  bool
}

type chMap[K comparable, V any] struct {
	ch   chan operate[K, V]
	m    map[K]V
	done chan struct{}
	// 单例模式 ——>  close
	once sync.Once
}

func NewChMap[K comparable, V any]() *chMap[K, V] {
	m := &chMap[K, V]{
		ch:   make(chan operate[K, V]),
		m:    make(map[K]V),
		done: make(chan struct{}),
	}
	go func() {
		for {
			var operate operate[K, V]
			select {
			case operate = <-m.ch:
				switch operate.op {
				case load:
					val, ok := m.m[operate.key]
					operate.ret <- result[V]{val, ok}
				case store:
					m.m[operate.key] = operate.val
				case del:
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

func (m *chMap[K, V]) Load(key K) (val V, flag bool) {
	operate := operate[K, V]{
		op:  load,
		key: key,
		ret: make(chan result[V]),
	}
	m.ch <- operate      // 将数据包装在Operate中进行传送
	res := <-operate.ret // 等待回写
	return res.val, res.ok
}
func (m *chMap[K, V]) Store(key K, val V) {
	operate := operate[K, V]{
		op:  store,
		key: key,
		val: val,
	}
	m.ch <- operate
}

func (m *chMap[K, V]) Delete(key K) {
	operate := operate[K, V]{
		op:  del,
		key: key,
	}
	m.ch <- operate
}
func (m *chMap[K, V]) Close() {
	m.once.Do(func() {
		close(m.done)
	})
}
