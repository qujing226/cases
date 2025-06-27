package test

import (
	"fmt"
	"math"
	"testing"
)

func minOperations(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}

	minDiff := math.MaxInt32
	minOps := -1

	for i := 0; i < n; i++ {
		currMin := nums[i]
		currMax := nums[i]
		for j := i; j < n; j++ {
			if nums[j] < currMin {
				currMin = nums[j]
			}
			if nums[j] > currMax {
				currMax = nums[j]
			}
			currDiff := currMax - currMin
			ops := i + (n - 1 - j)

			if currDiff < minDiff {
				minDiff = currDiff
				minOps = ops
			} else if currDiff == minDiff && ops < minOps {
				minOps = ops
			}
		}
	}

	return minOps
}

func Test_findAnagrams(t *testing.T) {
	fmt.Println(minOperations([]int{1, 2}))
}
