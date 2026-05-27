package main

/*
六、和堆解法为什么复杂度一样

这是一个面试中很容易被追问的问题。

分治：

每个元素参与 log k 次 merge
→ O(N log k)

小根堆：

每个元素：

push: O(log k)
pop : O(log k)

总：

N 次操作
→ O(N log k)
*/
func mergeArrays(arrays [][]int) []int {
	if len(arrays) <= 1 {
		return arrays[0]
	}
	var merge func(left, right int) []int
	merge = func(left, right int) []int {
		if left >= right {
			return arrays[left]
		}
		mid := left + (right-left)/2
		l := merge(left, mid)
		r := merge(mid+1, right)
		return mergeTwoArrays(l, r)
	}
	return merge(0, len(arrays)-1)
}

// 归并 O(m + n)
func mergeTwoArrays(a, b []int) []int {
	m, n := len(a), len(b)
	result := make([]int, m+n)
	i, j := 0, 0
	for i < m && j < n {
		if a[i] < b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	for i < m {
		result = append(result, a[i])
		i++
	}

	for j < n {
		result = append(result, b[j])
		j++
	}
	return result
}
