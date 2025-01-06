package main

import "github.com/qujing226/cases/demo/mullinkcheck/check"

func main() {
	err := check.ProcessLinks()
	if err != nil {
		panic(err)
	}
}
