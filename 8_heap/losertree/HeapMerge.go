// 参考 python heapq.merge
// 多路归并合并k个有序数据结构.

package main

import "fmt"

// 示例使用
func main() {
	it1 := NewSliceIterator([]int{1, 4, 7, 10})
	it2 := NewSliceIterator([]int{2, 5, 8, 11})
	it3 := NewSliceIterator([]int{3, 6, 9, 12})

	merged := HeapMerge([]Iterator[int]{it1, it2, it3}, func(a, b int) bool {
		return a < b
	})

	for merged.Next() {
		fmt.Println(merged.Value())
	}
}

type Iterator[E any] interface {
	Next() bool // Advances and returns true if there is a value at this new position.
	Value() E
}

type heapElement[E any] struct {
	value E
	it    Iterator[E]
}

type Merged[E any] struct {
	heap *Heap[heapElement[E]]
	cur  E
}

func HeapMerge[E any](iters []Iterator[E], less func(E, E) bool) Iterator[E] {
	heapLess := func(a, b heapElement[E]) bool {
		return less(a.value, b.value)
	}

	h := NewHeap(heapLess, nil)
	for _, it := range iters {
		if it.Next() {
			h.Push(heapElement[E]{value: it.Value(), it: it})
		}
	}
	return &Merged[E]{heap: h}
}

func (m *Merged[E]) Next() bool {
	if m.heap.Len() == 0 {
		return false
	}
	elem := m.heap.Pop()
	m.cur = elem.value
	if elem.it.Next() {
		m.heap.Push(heapElement[E]{value: elem.it.Value(), it: elem.it})
	}
	return true
}

func (m *Merged[E]) Value() E {
	return m.cur
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
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
	for i := (n >> 1) - 1; i >= 0; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; root > 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
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

// !示例迭代器实现
func NewSliceIterator[E any](slice []E) Iterator[E] {
	return &SliceIterator[E]{slice: slice, index: -1}
}

type SliceIterator[E any] struct {
	slice []E
	index int
}

func (it *SliceIterator[E]) Next() bool {
	it.index++
	return it.index < len(it.slice)
}

func (it *SliceIterator[E]) Value() E {
	return it.slice[it.index]
}
