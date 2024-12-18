package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	nums, _ := reader.ReadString('\n')
	nums = strings.TrimSpace(nums)
	num, _ := strconv.Atoi(nums)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	left, right := 0, num-1
	b := []byte(s)
	c := []byte(s)
	for left < len(s) {
		if right == len(s)-1 {
			right = 0
		}
		c[right] = b[left]
		left++
		right++
	}
	fmt.Println(string(c))

}
