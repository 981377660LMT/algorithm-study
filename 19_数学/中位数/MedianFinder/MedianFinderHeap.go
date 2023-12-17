package main

import "fmt"

func main() {
	M := NewMedianFinderHeap()
	M.Add(1)
	M.Add(3)
	M.Add(5)
	M.Add(2)
	fmt.Println(M.Query())
}

// 对顶堆维护中位数.
type MedianFinderHeap struct {
	left, right *Heap
}

func NewMedianFinderHeap() *MedianFinderHeap {
	return &MedianFinderHeap{
		left:  NewHeap(func(a, b H) bool { return a > b }, nil),
		right: NewHeap(func(a, b H) bool { return a < b }, nil),
	}
}

func (m *MedianFinderHeap) Add(num int) {
	if m.left.Len() == m.right.Len() {
		m.right.Push(num)
		m.left.Push(m.right.Pop())
	} else {
		m.left.Push(num)
		m.right.Push(m.left.Pop())
	}
}

// 查询结果会向下取整.
func (m *MedianFinderHeap) Query() int {
	if m.Len() == 0 {
		panic("MedianFinderHeap is empty")
	}

	if m.left.Len() == m.right.Len() {
		return (m.left.Peek() + m.right.Peek()) >> 1
	} else {
		return m.left.Peek()
	}
}

func (m *MedianFinderHeap) Len() int {
	return m.left.Len() + m.right.Len()
}

type H = int

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
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
