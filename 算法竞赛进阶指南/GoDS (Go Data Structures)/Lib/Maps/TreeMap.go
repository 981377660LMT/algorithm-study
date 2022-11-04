package maps

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
)

func e() {
	treeMap := treemap.NewWithIntComparator()
	treeMap.Put(10, "a")
	treeMap.Put(20, "b")
	treeMap.Put(0, "c")
	fmt.Println(treeMap.Get(1)) // a true

	// floor ceiling
	fmt.Println(treeMap.Floor(0))    // 0 c
	fmt.Println(treeMap.Floor(15))   // 10 a
	fmt.Println(treeMap.Floor(20))   // 20 b
	fmt.Println(treeMap.Ceiling(15)) // 20 b
	fmt.Println(treeMap.Ceiling(20)) // 20 b
	fmt.Println(treeMap.Ceiling(21)) // <nil> <nil>

	// min max
	fmt.Println(treeMap.Min())
	fmt.Println(treeMap.Max())
}
