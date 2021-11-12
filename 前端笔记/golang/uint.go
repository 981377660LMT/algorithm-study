package main

import "fmt"

func main() {
	const a uint = 1
	const b uint = 2
	// a - b (constant -1 of type uint) overflows uint
	fmt.Println(a - b)
}
