package main

import "fmt"

func main() {
	i := 0
	defer func(i int) {
		fmt.Println(i) // 0
	}(i)

	i++

	fmt.Println(-12 % 29)

	type pair struct{ x, y int }
	p1 := pair{1, 2}
	p2 := pair{1, 2}
	fmt.Println(p1 == p2)
}
