package main

/*
你有一个模块，需要记录一些整数的区间，这些区间会随着时间变化进行增、删、查询。
简单点说，就是你要管理一个整数的区间，可以插入一个新的区间，也可以删除一个区间，还可以查询某个整数是否在当前所有区间内。

接下来说一下我对这个问题的理解和解决方案。首先，我们可以把这个问题当做一个区间合并问题来思考。
比如，两个区间 [1, 3] 和 [2, 5]，我们可以把它们合并成 [1, 5]。
而当我们删除一个区间时，也是从这个区间中删除一个部分或者一个完整区间。

问题的关键是如何高效地执行这些操作。直接暴力实现肯定不行，性能可能会很差，尤其是删除和查询操作。我们需要找到一种更高效的数据结构。
*/

/*
解决方案：
区间合并：
	我们可以使用一个有序列表来存储所有的区间，列表中的区间按照区间的左端点排序。
	这样，当我们插入一个新的区间时，我们只需要遍历这个列表，找到合适的位置进行合并或者插入。

删除区间：
	删除一个区间就是找到对应的区间，然后根据具体情况进行拆分或者删除。
	比如，你要删除一个区间 [4, 7]，可能会有以下几种情况：
		- 区间完全包含要删除的区间。
		- 区间和要删除的区间有交集。
		- 区间和要删除的区间完全不重叠。
查询整数：
	查询某个整数是否在当前的区间内，我们只需要遍历区间列表，判断这个整数是否在某个区间内即可。
*/

type RangeModule struct {
	intervals [][]int
}

func Constructor() RangeModule {
	return RangeModule{
		intervals: [][]int{},
	}
}

// AddRange adds a range [start, end) to the module.
func (rm *RangeModule) AddRange(start int, end int) {
	var newInterval [][]int
	inserted := false
	for _, interval := range rm.intervals {
		if start > interval[1] {
			newInterval = append(newInterval, interval)
		} else if end < interval[0] {
			if !inserted {
				newInterval = append(newInterval, []int{start, end})
				inserted = true
			}
			newInterval = append(newInterval, interval)
		} else {
			start = min(start, interval[0])
			end = max(end, interval[1])
			newInterval = append(newInterval, []int{start, end})
		}
	}
	if !inserted {
		newInterval = append(newInterval, []int{start, end})
	}
	rm.intervals = newInterval
}

// RemoveRange removes a range [start, end) from the module.
func (rm *RangeModule) RemoveRange(start int, end int) {
	newIntervals := [][]int{}
	for _, interval := range rm.intervals {
		if interval[1] < start || interval[0] > end {
			newIntervals = append(newIntervals, interval)
		} else {
			if interval[0] < start {
				newIntervals = append(newIntervals, []int{interval[0], start})
			}
			if interval[1] > end {
				newIntervals = append(newIntervals, []int{end, interval[1]})
			}
		}
	}
	rm.intervals = newIntervals
}

func (rm *RangeModule) Query(num int) bool {
	for _, interval := range rm.intervals {
		if num >= interval[0] && num <= interval[1] {
			return true
		}
	}
	return false
}
