package main

import (
	"fmt"
	"github.com/qujing226/cases/algo/person"
)

type Node struct {
	val int
}

func main() {
	p := &person.Person{"Bob", 25}
	// 编译时错误: age字段是私有的，无法直接访问
	fmt.Println(p.age)
}
