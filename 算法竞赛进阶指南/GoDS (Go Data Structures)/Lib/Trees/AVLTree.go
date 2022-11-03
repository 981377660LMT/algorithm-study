package main

import (
	"fmt"

	"github.com/emirpasic/gods/trees/avltree"
)

func main() {
	tree := avltree.NewWithIntComparator()
	tree.Put(1, "a")
	tree.Put(2, "b")
	tree.Put(3, "c")
	tree.Put(4, "d")
	fmt.Println(tree.Get(1)) // a true
}
