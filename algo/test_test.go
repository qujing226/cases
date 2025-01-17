package main

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
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

func Test(t *testing.T) {
	input := "a1b2c3"
	reader := bufio.NewReader(strings.NewReader(input))
	var s string
	_, err := fmt.Fscanln(reader, &s)
	if err != nil {
		return
	}
	var target []byte = []byte("number")
	s = strings.TrimSpace(s)
	sb := []byte(s)
	num := 0
	length := len(sb)
	for _, val := range sb {
		if val > '0' && val < '9' {
			num++
		}
	}
	for i := 0; i < num; i++ {
		sb = append(sb, "     "...)
	}
	left, right := length-1, len(sb)-1
	for left != right {
		if sb[left] < '0' || sb[left] > '9' {
			sb[right] = sb[left]
		} else {
			for i, val := range target {
				sb[right-5+i] = val
			}
			right -= 5
		}
		left--
		right--
	}
	fmt.Println(string(sb))
}
