package maps

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
)

func e() {
	mp := treemap.NewWithIntComparator()
	mp.Put(10, "a")
	mp.Put(20, "b")
	mp.Put(0, "c")
	mp.Iterator()
	fmt.Println(mp.Get(1)) // a true

	// floor ceiling
	fmt.Println(mp.Floor(0))    // 0 c
	fmt.Println(mp.Floor(15))   // 10 a
	fmt.Println(mp.Floor(20))   // 20 b
	fmt.Println(mp.Ceiling(15)) // 20 b
	fmt.Println(mp.Ceiling(20)) // 20 b
	fmt.Println(mp.Ceiling(21)) // <nil> <nil>

	// min max
	fmt.Println(mp.Min())
	fmt.Println(mp.Max())
}
