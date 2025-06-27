package BloomFilter

import (
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	bitmap  []bool
	k       int
	m       int
	n       int
	hashFns []func(data []byte) uint32
}

// NewBloomFilter New
func NewBloomFilter(expectedItems int, falsePositiveRate float64) *BloomFilter {
	m := optimalSize(expectedItems, falsePositiveRate)
	k := optimalHashFunctions(m, expectedItems)
	bf := &BloomFilter{
		bitmap:  make([]bool, m),
		k:       k,
		m:       m,
		n:       expectedItems,
		hashFns: make([]func(data []byte) uint32, k),
	}
	for i := 0; i < k; i++ {
		bf.hashFns[i] = generateHashFunction(i)
	}
	return bf
}

// Add 数据添加
func (bf *BloomFilter) Add(data string) {
	bytes := []byte(data)
	for _, hashFn := range bf.hashFns {
		index := hashFn(bytes) % uint32(bf.m)
		bf.bitmap[index] = true
	}
}

// Test 检查数据是否可能位于布隆过滤器
func (bf *BloomFilter) Test(data string) bool {
	bytes := []byte(data)
	for _, hashFn := range bf.hashFns {
		index := hashFn(bytes) % uint32(bf.m)
		if !bf.bitmap[index] {
			return false
		}
	}
	return true
}

// optimalSize 计算最优的布隆过滤器大小
func optimalSize(n int, p float64) int {
	return int(-float64(n) * math.Log(p) / (math.Ln2 * math.Ln2))
}

// optimalHashFunctions 计算最优的hashFunc个数
func optimalHashFunctions(m, n int) int {
	return int(float64(m) / float64(n) * math.Ln2)
}

// generateHashFunction 生成哈希函数
func generateHashFunction(seed int) func(data []byte) uint32 {
	return func(data []byte) uint32 {
		h := fnv.New32a()
		h.Write([]byte{byte(seed)})
		h.Write(data)
		return h.Sum32()
	}
}
