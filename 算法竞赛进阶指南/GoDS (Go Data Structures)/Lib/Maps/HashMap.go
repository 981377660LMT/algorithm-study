package maps

// A map based on hash tables. Keys are unordered.

import (
	"fmt"

	"github.com/emirpasic/gods/maps/hashmap"
)

func b() {
	mp := hashmap.New()
	mp.Put("a", 1)
	mp.Put("b", 2)
	mp.Put(1, "a")
	fmt.Println(mp.Get("a")) // 1 true
	fmt.Println(mp.Get("b")) // 2 true
	fmt.Println(mp.Get(1))   // a true
	for _, key := range mp.Keys() {
		fmt.Println(key)
		fmt.Println(mp.Get(key))
	}

	for _, value := range mp.Values() {
		fmt.Println(value)
	}
}
