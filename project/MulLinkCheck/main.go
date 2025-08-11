package main

import (
	"github.com/qujing226/cases/project/MulLinkCheck/check"
)

func main() {
	err := check.ProcessLinks()
	if err != nil {
		panic(err)
	}
}
