package main

type MinHeap struct {
	data []int
	size int
}

func NewMinHeap(size int) *MinHeap {
	d := make([]int, 0)
	return &MinHeap{data: d, size: size}
}

func (mh *MinHeap) Size() int {
	return mh.size
}

func (mh *MinHeap) siftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if mh.data[index] < mh.data[parent] {
			mh.data[index], mh.data[parent] = mh.data[parent], mh.data[index]
			index = parent
		} else {
			break
		}
	}
}

func (mh *MinHeap) siftDown(index int) {
	n := len(mh.data)
	for {
		left, right := index*2+1, index*2+2
		smallest := left

		if left < n && mh.data[left] < mh.data[smallest] {
			smallest = left
		}
		if right < n && mh.data[right] < mh.data[smallest] {
			smallest = right
		}
		if smallest != index {
			mh.data[index], mh.data[smallest] = mh.data[smallest], mh.data[index]
			index = smallest
		} else {
			break
		}
	}
}

func (mh *MinHeap) Push(value int) {
	if len(mh.data) < mh.size {
		mh.data = append(mh.data, value)
		mh.siftUp(mh.size)
	} else if value > mh.data[0] {
		mh.data[0] = value
		mh.siftDown(0)
	}
}
