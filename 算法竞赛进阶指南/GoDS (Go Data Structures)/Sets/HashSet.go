package main

import (
	"fmt"

	"github.com/emirpasic/gods/sets/hashset"
)

func main() {
	set := hashset.New()
	set.Add("a")
	set.Add("b")
	fmt.Println(set.Contains("a")) // true
	set.Remove("a")
	println(set.Contains("a")) // false
	newSet := set.Union(hashset.New("b", "c"))
	for _, v := range newSet.Values() {
		fmt.Println(v)
	}
}
