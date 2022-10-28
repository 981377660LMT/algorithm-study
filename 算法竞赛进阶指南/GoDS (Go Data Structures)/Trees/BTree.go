package main

import (
	"fmt"

	"github.com/emirpasic/gods/trees/btree"
)

func main() {
	tree := btree.NewWithIntComparator(4) // 3 is the minimum degree (defines the range for number of children nodes)
	tree.Put(1, "a")
	tree.Put(2, "b")
	tree.Put(3, "c")
	tree.Put(4, "d")
	fmt.Println(tree.Get(1)) // a true
}
