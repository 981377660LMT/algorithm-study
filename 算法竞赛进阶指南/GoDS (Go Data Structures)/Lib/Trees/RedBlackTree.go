package trees

import (
	"fmt"

	"github.com/emirpasic/gods/trees/redblacktree"
)

func a() {
	tree := redblacktree.NewWithIntComparator()
	tree.Put(1, "a")
	tree.Put(2, "b")

	fmt.Println(tree.Get(1)) // a true
}
