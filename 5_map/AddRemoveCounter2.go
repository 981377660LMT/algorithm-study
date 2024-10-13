package main

import "fmt"

func main() {
	counter := make(AddRemoveCounter[int32], 10)

	onAdd := func(v int32, c int32) {
		fmt.Printf("Add: %d, %d\n", v, c)
	}

	onRemove := func(v int32, c int32) {
		fmt.Printf("Remove: %d, %d\n", v, c)
	}

	counter.Add(1, onAdd, onRemove)
	counter.Add(2, onAdd, onRemove)
	counter.Add(1, onAdd, onRemove)

	fmt.Println(len(counter))
}

type AddRemoveCounter[T comparable] map[T]int32

func NewAddRemoveCounter[T comparable](cap int32) AddRemoveCounter[T] {
	return make(AddRemoveCounter[T], cap)
}

func (c AddRemoveCounter[T]) Add(x T, onAdd, onRemove func(v T, c /* c > 0 */ int32)) {
	preFreq := c[x]
	c[x]++
	if preFreq != 0 {
		onRemove(x, preFreq)
	}
	onAdd(x, preFreq+1)
}

func (c AddRemoveCounter[T]) Remove(x T, onAdd, onRemove func(v T, c /* c > 0 */ int32)) bool {
	preFreq := c[x]
	if preFreq == 0 {
		return false
	}
	c[x]--
	onRemove(x, preFreq)
	if preFreq > 1 {
		onAdd(x, preFreq-1)
	} else {
		delete(c, x)
	}
	return true
}
