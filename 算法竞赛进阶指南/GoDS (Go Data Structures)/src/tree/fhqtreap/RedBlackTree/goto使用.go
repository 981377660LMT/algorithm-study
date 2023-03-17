package main

import "fmt"

func foo() {
	fmt.Println(demo(2))
}

func demo(num int) int {
	if num&1 == 0 {
		goto even
	} else {
		goto odd
	}

even:
	println("even")
	return 0 // !如果这里不return,则会执行下面的odd
odd:
	println("odd")
	return 1
}
