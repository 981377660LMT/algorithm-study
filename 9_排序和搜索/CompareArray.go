package main

import (
	"fmt"
	"slices"
)

func main() {
	fmt.Println(slices.Compare([]int{1, 2, 3}, []int{1, 2, 3}))
}
