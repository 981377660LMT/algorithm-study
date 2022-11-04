package cmnx

import (
	"fmt"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/emirpasic/gods/trees/binaryheap"
)

func foo() {
	fmt.Println("Hello, gods")

	// array/stack
	arrayList := arraylist.New()

	// queue/deque
	linkedList := doublylinkedlist.New()

	// map
	hashMap := hashmap.New()
	linkedHashMap := linkedhashmap.New()
	treeMap := treemap.NewWithIntComparator()

	// set
	hashSet := hashset.New()
	linkedHashSet := linkedhashset.New()
	treeSet := treemap.NewWithIntComparator()

	// heap
	heap := binaryheap.NewWithIntComparator()

	_ = []interface{}{
		arrayList,
		linkedList,
		hashMap, linkedHashMap, treeMap,
		hashSet, linkedHashSet, treeSet,
		heap,
	}

}
