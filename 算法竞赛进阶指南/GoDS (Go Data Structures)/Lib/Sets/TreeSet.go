package sets

import (
	"github.com/emirpasic/gods/sets/treeset"
)

func ts() {
	set := treeset.NewWithIntComparator()
	set.Add(1, 2, 3)
	it1 := set.Iterator()
	it2 := set.Iterator()
	it1.Next()
	it2.Next()
	it2.Next()

}
