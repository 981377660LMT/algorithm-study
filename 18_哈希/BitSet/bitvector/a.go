package main

import "fmt"

func main() {
	fmt.Println(1<<2 | 3 ^ 1)
	fmt.Println((1 << 2) | (3 ^ 1))
	fmt.Println(1<<3 - 1)
}
