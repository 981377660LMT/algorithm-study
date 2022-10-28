package main

import (
	"fmt"

	"github.com/emirpasic/gods/sets/treeset"
)

func main() {
	set := treeset.NewWithIntComparator()
	set.Add(1)
	set.Add(2, 2, 3, 4, -1)
	set.Remove(1, 10)
	set.Each(func(index int, value interface{}) {
		fmt.Println(index, value)
	})
}
