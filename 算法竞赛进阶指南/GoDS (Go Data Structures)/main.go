package main

import "fmt"

func main() {
	i := 0
	defer func(i int) {
		fmt.Println(i) // 0
	}(i)

	i++
}
