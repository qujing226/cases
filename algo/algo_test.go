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

// 逆波兰函数
func Test_evalRPN(t *testing.T) {
	num := evalRPN([]string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"})
	t.Log(num)
}

func Test_kmp(t *testing.T) {
	str := "mississippi"
	pat := "issip"
	t.Log(KMP(str, pat))
}

func TestMove(t *testing.T) {
	fmt.Println(Move([]int{0, 1, 0, 0, 3, 12}))
}

func Move(nums []int) []int {
	length := len(nums)
	fast, slow := 0, 0
	for ; fast < length; fast++ {
		if nums[fast] == 0 {
			continue
		} else {
			nums[slow], nums[fast] = nums[fast], nums[slow]
			slow++
		}
	}
	return nums
}

func main() {
	loop := 0
	fmt.Scanln(&loop)
	for ; loop > 0; loop-- {
		flag := false
		n, m, k := 0, 0, 0
		fmt.Scan(&n, &m, &k)
		min_step := max(n, m)
		if k < min_step {
			fmt.Println(-1)
			continue
		}
		remain := k - min_step
		// 我的思路是先到达目标点，然后如果是剩余奇数步，就走一个三角形，先有一个斜边再走两个直角边到达原点
		if remain%2 == 1 {
			flag = true
			remain -= 3
		}
		maxDiagonal := min(n, m) + remain
		if flag {
			fmt.Println(maxDiagonal + 1)
		} else {
			fmt.Println(maxDiagonal)
		}

	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func TestQ(t *testing.T) {
	ch := make(chan int, 7)
	for i := 0; i < 8; i++ {
		ch <- i
	}
}
