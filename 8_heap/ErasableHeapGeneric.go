// 懒删除堆/可删除堆

package main

import "fmt"

func main() {
	ehp := NewErasableHeapGeneric(func(a, b int) bool { return a > b }, []int{1, 2, 3, 4, 5}...)
	ehp.Erase(1)
	ehp.Erase(4)
	fmt.Println(ehp.Pop())
	fmt.Println(ehp.Len())
	fmt.Println(ehp.Pop())
	fmt.Println(ehp.Pop())
}

type ErasableHeapGeneric[H comparable] struct {
	data   *HeapGeneric[H]
	erased *HeapGeneric[H]
	size   int
}

func NewErasableHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *ErasableHeapGeneric[H] {
	return &ErasableHeapGeneric[H]{NewHeapGeneric(less, nums...), NewHeapGeneric(less), len(nums)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeapGeneric[H]) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
	h.size--
}

func (h *ErasableHeapGeneric[H]) Push(value H) {
	h.data.Push(value)
	h.normalize()
}

func (h *ErasableHeapGeneric[H]) Pop() (value H) {
	value = h.data.Pop()
	h.normalize()
	h.size--
	return
}

func (h *ErasableHeapGeneric[H]) Peek() (value H) {
	value = h.data.Top()
	return
}

func (h *ErasableHeapGeneric[H]) Len() int {
	return h.size
}

func (h *ErasableHeapGeneric[H]) Clear() {
	h.data.Clear()
	h.erased.Clear()
	h.size = 0
}

func (h *ErasableHeapGeneric[H]) normalize() {
	for h.data.Len() > 0 && h.erased.Len() > 0 && h.data.Top() == h.erased.Top() {
		h.data.Pop()
		h.erased.Pop()
	}
}

type HeapGeneric[H comparable] struct {
	data []H
	less func(a, b H) bool
}

func NewHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *HeapGeneric[H] {
	nums = append(nums[:0:0], nums...)
	heap := &HeapGeneric[H]{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

func (h *HeapGeneric[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *HeapGeneric[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *HeapGeneric[H]) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *HeapGeneric[H]) Len() int { return len(h.data) }

func (h *HeapGeneric[H]) Clear() {
	h.data = h.data[:0]
}

func (h *HeapGeneric[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *HeapGeneric[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *HeapGeneric[H]) pushDown(root int) {
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
