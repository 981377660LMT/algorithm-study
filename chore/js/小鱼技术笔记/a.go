package main

import (
	"fmt"
	"sync"
)

func doing(i int) {
	fmt.Println(i)
}

func main() {
	// defer的作用域与变量是不一样的，它是函数作用域，不是块作用域
	for i := 0; i != 10; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("mc")

	waitgroup := sync.WaitGroup{}

	for i := 0; i != 5; i++ {
		waitgroup.Add(1)
		go func(i int) {
			defer waitgroup.Done()
			doing(i)
		}(i)
	}

	waitgroup.Wait()
}
