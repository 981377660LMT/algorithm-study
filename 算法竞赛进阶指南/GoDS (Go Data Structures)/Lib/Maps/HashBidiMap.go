// A map based on two hashmaps. Keys are unordered.

package main

import (
	"fmt"

	"github.com/emirpasic/gods/maps/hashbidimap"
)

func main() {
	bidiMap := hashbidimap.New()
	bidiMap.Put("a", 1)
	bidiMap.Put("b", 2)
	bidiMap.Put(1, "a")
	fmt.Println(bidiMap.Get("a"))  // 1 true
	fmt.Println(bidiMap.GetKey(1)) // a true
	fmt.Println(bidiMap.Keys()...)
	fmt.Println(bidiMap.Values()...)
}
