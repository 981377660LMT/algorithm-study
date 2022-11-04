package trees

// 堆

import (
	"fmt"

	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/utils"
)

func c() {
	inverseComparator := func(a, b interface{}) int {
		return -utils.IntComparator(a, b)
	}

	// Create a new heap with a custom comparator
	heap := binaryheap.NewWith(inverseComparator)
	heap.Push(1, 3, 2)
	fmt.Println(heap.Values()) // [3 2 1]

}

// 自定义比较器
type Element struct {
	name     string
	priotity int
}

func comparator(e1, e2 interface{}) int {
	priority1 := e1.(Element).priotity
	priority2 := e2.(Element).priotity
	return -utils.IntComparator(priority1, priority2)
}

func main() {
	// Create a new heap with a custom comparator
	heap := binaryheap.NewWith(comparator)
	heap.Push(Element{"a", 1}, Element{"b", 2}, Element{"c", 3})
	fmt.Println(heap.Values()) // [3 2 1]
}
