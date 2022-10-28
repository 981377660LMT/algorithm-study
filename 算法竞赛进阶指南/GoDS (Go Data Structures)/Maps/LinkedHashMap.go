package main

import (
	"fmt"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func main() {
	linkedHashMap := linkedhashmap.New()
	linkedHashMap.Put("a", 1)
	linkedHashMap.Put("b", 2)
	linkedHashMap.Put(1, "a")
	linkedHashMap.Put(2, "b")
	linkedHashMap.Put(3, "c")
	linkedHashMap.Each(func(key interface{}, value interface{}) {
		fmt.Println(key, value)
	})
}
