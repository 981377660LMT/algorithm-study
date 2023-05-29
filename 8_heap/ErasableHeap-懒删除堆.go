// 懒删除堆/可删除堆

package main

import "fmt"

func main() {
	ehp := NewErasableHeap(func(a, b H) bool { return a < b }, []H{1, 2, 3, 4, 5})
	ehp.Erase(1)
	ehp.Erase(4)
	fmt.Println(ehp.Pop())
	fmt.Println(ehp.Len())
	fmt.Println(ehp.Pop())
	fmt.Println(ehp.Pop())
}

type H = int

type ErasableHeap struct {
	base   *Heap
	erased *Heap
}

func NewErasableHeap(less func(a, b H) bool, nums []H) *ErasableHeap {
	return &ErasableHeap{NewHeap(less, nums), NewHeap(less, nil)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeap) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
}

func (h *ErasableHeap) Push(value H) {
	h.base.Push(value)
	h.normalize()
}

func (h *ErasableHeap) Pop() (value H) {
	value = h.base.Pop()
	h.normalize()
	return
}

func (h *ErasableHeap) Peek() (value H) {
	value = h.base.Top()
	return
}

func (h *ErasableHeap) Len() int {
	return h.base.Len()
}

func (h *ErasableHeap) normalize() {
	for h.base.Len() > 0 && h.erased.Len() > 0 && h.base.Top() == h.erased.Top() {
		h.base.Pop()
		h.erased.Pop()
	}
}

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
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

func (h *Heap) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
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
