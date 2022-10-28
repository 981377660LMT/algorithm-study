package main

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
)

func main() {
	mp := treebidimap.NewWith(utils.IntComparator, utils.StringComparator)
	mp.Put(1, "a")
	mp.Put(2, "b")
	mp.Put(3, "c")

	fmt.Println(mp.Get(1)) // a true
}
