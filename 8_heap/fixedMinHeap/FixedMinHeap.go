package main

import "fmt"

func main() {
	fHeap := NewFixedMinHeap(3, func(a, b int) bool { return a < b })
	fHeap.Push(1)
	fHeap.Push(2)
	fHeap.Push(3)
	fmt.Println(fHeap.PeekMax()) // 1
	fHeap.Push(1)
	fmt.Println(fHeap.PeekMax()) // 2
	fHeap.Push(4)
	fmt.Println(fHeap.PeekMax()) // 2
	fHeap.Push(0)
	fmt.Println(fHeap.PeekMax()) // 1
}

// 带最大容量限制的最小堆，用于维护topK.
type FixedMinHeap[H any] struct {
	pq      *Heap[H]
	maxSize int
	less    func(a, b H) bool
}

func NewFixedMinHeap[H any](maxSize int, less func(a, b H) bool) *FixedMinHeap[H] {
	if maxSize == 0 {
		panic("maxSize must be positive")
	}
	res := &FixedMinHeap[H]{maxSize: maxSize, less: less}
	reversedLess := func(a, b H) bool { return !less(a, b) }
	res.pq = NewHeap(reversedLess, nil)
	return res
}

func (f *FixedMinHeap[H]) Push(value H) bool {
	if f.pq.Len() < f.maxSize {
		f.pq.Push(value)
		return true
	}
	if f.less(value, f.pq.Top()) {
		f.pq.Pop()
		f.pq.Push(value)
		return true
	}
	return false
}

func (f *FixedMinHeap[H]) PopMax() H {
	return f.pq.Pop()
}

func (f *FixedMinHeap[H]) PeekMax() H {
	return f.pq.Top()
}

func (f *FixedMinHeap[H]) Len() int {
	return f.pq.Len()
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
