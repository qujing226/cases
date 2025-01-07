package main

import (
	"reflect"
	"testing"
)

// TestAddRange Range模块
func TestAddRange(t *testing.T) {
	tests := []struct {
		rm       RangeModule
		start    int
		end      int
		expected [][]int
	}{
		{RangeModule{[][]int{{1, 3}}}, 2, 5, [][]int{{1, 5}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 2, 5, [][]int{{1, 5}, {6, 9}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 10, 12, [][]int{{1, 3}, {6, 9}, {10, 12}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 4, 5, [][]int{{1, 3}, {4, 5}, {6, 9}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 5, 6, [][]int{{1, 3}, {5, 9}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 0, 1, [][]int{{0, 3}, {6, 9}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 0, 0, [][]int{{0, 0}, {1, 3}, {6, 9}}},
		{RangeModule{[][]int{{1, 3}, {6, 9}}}, 10, 10, [][]int{{1, 3}, {6, 9}, {10, 10}}},
	}

	for _, test := range tests {
		test.rm.AddRange(test.start, test.end)
		if !reflect.DeepEqual(test.rm.intervals, test.expected) {
			t.Errorf("AddRange(%d, %d) = %v, want %v", test.start, test.end, test.rm.intervals, test.expected)
		}
	}
}
