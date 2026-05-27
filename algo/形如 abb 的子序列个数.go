package main

import "fmt"

/*
寻找一个字符串中所有形如abb的个数，a和b一定不一样，b和b一定一样

抽象问题：
看到第K个b，他可以和前面的某个第J个b组成pair，在J之前的任何一个非b的第I个字符都可以组成abb结构

pair[c]
ans += pair[c]
*/

func FindAbb() {
	input := "abbwddnmwnzsiwmnwdddapdxxjwm"
	var count [26]int
	var pair [26]int
	res := 0
	for idx, char := range input {
		res += pair[char-'a']
		pair[char-'a'] += idx - count[char-'a']
		count[char-'a']++
	}
	fmt.Println(res)
}
