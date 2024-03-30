package main

import "fmt"

func main() {
	type node struct {
		left, right *node
		size        int32
	}

	buffer := NewBufferWithCleaner[*node](
		func() *node { return &node{} },
		func(node *node) {
			node.left = nil
			node.right = nil
			node.size = 0
		},
	)

	root := buffer.Alloc()
	fmt.Println(root)
}

// A buffer that recycles objects.
type Buffer[T any] struct {
	recycles    []T
	supplier    func() T // Supplier is a function that returns a new object.
	cleaner     func(T)  // Cleaner is a function that cleans up / initializes an object.
	allocTime   int32
	releaseTime int32
}

func NewBuffer[T any](supplier func() T) *Buffer[T] {
	return NewBufferWithCleaner(supplier, func(T) {})
}

func NewBufferWithCleaner[T any](supplier func() T, cleaner func(T)) *Buffer[T] {
	return NewBufferWithCleanerAndCapacity(supplier, cleaner, 0)
}

func NewBufferWithCleanerAndCapacity[T any](supplier func() T, cleaner func(T), capacity int32) *Buffer[T] {
	return &Buffer[T]{
		recycles: make([]T, 0, capacity),
		supplier: supplier,
		cleaner:  cleaner,
	}
}

func (b *Buffer[T]) Alloc() T {
	b.allocTime++
	if len(b.recycles) == 0 {
		res := b.supplier()
		b.cleaner(res)
		return res
	} else {
		res := b.recycles[0]
		b.recycles = b.recycles[1:]
		return res
	}
}

func (b *Buffer[T]) Release(e T) {
	b.releaseTime++
	b.cleaner(e)
	b.recycles = append(b.recycles, e)
}

func (b *Buffer[T]) Check() {
	if b.allocTime != b.releaseTime {
		panic(fmt.Sprintf("Buffer alloc %d but release %d", b.allocTime, b.releaseTime))
	}
}
