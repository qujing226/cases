package main

import (
	"fmt"
	"strconv"
)

func evalRPN(tokens []string) int {
	lt := len(tokens)
	stack := make([]int, 0, lt)
	for _, val := range tokens {
		flag := IsOperator(val)
		if flag != 0 {
			numA, numB := stack[len(stack)-2], stack[len(stack)-1]
			numC := Calculate(flag, numA, numB)
			fmt.Println(numA, val, numB, "=", numC)
			stack = stack[:len(stack)-2]
			stack = append(stack, numC)
		} else {
			num, _ := strconv.Atoi(val)
			stack = append(stack, num)
		}
	}
	return stack[0]
}

func IsOperator(b string) int {
	if b == "+" {
		return 1
	} else if b == "-" {
		return 2
	} else if b == "*" {
		return 3
	} else if b == "/" {
		return 4
	} else {
		return 0
	}
}

func Calculate(flag int, numA, numB int) int {
	switch flag {
	case 1:
		return numA + numB
	case 2:
		return numA - numB
	case 3:
		return numA * numB
	case 4:
		return numA / numB
	default:
		return 0
	}
}
