package binaryheap

import (
	"cmnx/src/heap"
	"fmt"
)

// Assert Heap implementation
var _ heap.Heap = (*BinaryHeap)(nil)

func main() {
	heap := NewBinaryHeap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	}, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	heap.Push(0)
	heap.Push(-1)
	heap.Push(-2)
	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
	fmt.Println(heap.Pop())
}

func NewBinaryHeap(comparator Comparator, nums []interface{}) *BinaryHeap {
	numsCopy := append([]interface{}{}, nums...)
	heap := &BinaryHeap{comparator: comparator, data: numsCopy}
	heap.heapify()
	return heap
}

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b interface{}) int

type BinaryHeap struct {
	data       []interface{}
	comparator Comparator
}

func (h *BinaryHeap) Push(value interface{}) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *BinaryHeap) Pop() (value interface{}) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *BinaryHeap) Peek() (value interface{}) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	return
}

func (h *BinaryHeap) Len() int { return len(h.data) }

func (h *BinaryHeap) heapify() {
	for i := (h.Len() >> 1) - 1; i >= 0; i-- {
		h.pushDown(i)
	}
}

func (h *BinaryHeap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *BinaryHeap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
