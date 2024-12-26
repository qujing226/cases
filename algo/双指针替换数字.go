package main

import "fmt"

func replaceNumber(strByte []byte) string {
	// 查看有多少字符
	numCount, oldSize := 0, len(strByte)
	for i := 0; i < len(strByte); i++ {
		if (strByte[i] <= '9') && (strByte[i] >= '0') {
			numCount++
		}
	}
	// 增加长度
	for i := 0; i < numCount; i++ {
		strByte = append(strByte, []byte("     ")...)
	}
	tmpBytes := []byte("number")
	// 双指针从后遍历
	leftP, rightP := oldSize-1, len(strByte)-1
	for leftP < rightP {
		rightShift := 1
		// 如果是数字则加入number
		if (strByte[leftP] <= '9') && (strByte[leftP] >= '0') {
			for i, tmpByte := range tmpBytes {
				strByte[rightP-len(tmpBytes)+i+1] = tmpByte
			}
			rightShift = len(tmpBytes)
		} else {
			strByte[rightP] = strByte[leftP]
		}
		// 更新指针
		rightP -= rightShift
		leftP -= 1
	}
	return string(strByte)
}

func replaceNumberMain() {
	var strByte []byte
	_, err := fmt.Scanln(&strByte)
	if err != nil {
		return
	}

	newString := replaceNumber(strByte)

	fmt.Println(newString)
}

func replace() {
	var s []byte
	_, err := fmt.Scanln(&s)
	if err != nil {
		panic(err)
	}
	num, l := 0, len(s)-1
	for _, v := range s {
		if v >= '0' && v <= '9' {
			num++
		}
	}
	for i := 0; i < num; i++ {
		s = append(s, "     "...)
	}
	tem := []byte("number")
	tl := 5
	left, right := l, len(s)-1
	for left >= 0 {
		if s[left] < '0' || s[left] > '9' {
			s[right] = s[left]
		} else {
			for i, val := range tem {
				s[right-tl+i] = val
			}
			right -= tl
		}
		left--
		right--
	}
	fmt.Println(string(s))
}
