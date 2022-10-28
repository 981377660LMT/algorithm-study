package main

import (
	"fmt"

	"github.com/emirpasic/gods/sets/linkedhashset"
)

func main() {
	set := linkedhashset.New(1, 2, 3, 4, 5, 0)
	fmt.Println(set.Contains(1, 2)) // true
	set.Each(func(index int, value interface{}) {
		fmt.Println(index, value)
	})
}
