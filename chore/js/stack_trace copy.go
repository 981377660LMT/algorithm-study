package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

func bar() {
	panic(errors.New("failed"))
}

func foo() {
	bar()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in main:", r)
			fmt.Println(string(debug.Stack()))
		}
	}()

	foo()
}
