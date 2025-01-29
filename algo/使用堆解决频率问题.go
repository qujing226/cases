package main

import "container/heap"

func topKFrequent(nums []int, k int) []int {
	maps := make(map[int]int)
	for _, val := range nums {
		maps[val]++
	}
	h := &IHeap{}
	heap.Init(h)
	// 所有元素入堆，堆长为k
	for key, val := range maps {
		heap.Push(h, [2]int{key, val})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	res := make([]int, k)
	// 按顺序返回堆中的元素
	for i := 0; i < k; i++ {
		res[k-i-1] = heap.Pop(h).([2]int)[0]
	}
	return res
}

type IHeap [][2]int

func (h *IHeap) Len() int {
	return len(*h)
}

func (h *IHeap) Less(i, j int) bool {
	return (*h)[i][1] < (*h)[j][1]
}

func (h *IHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}
func (h *IHeap) Push(x any) {
	*h = append(*h, x.([2]int))
}

func (h *IHeap) Pop() any {
	item := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return item
}
