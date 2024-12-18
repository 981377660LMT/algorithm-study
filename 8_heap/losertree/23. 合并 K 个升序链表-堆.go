// https://leetcode.cn/problems/merge-k-sorted-lists/description/

package main

type ListNode struct {
	Val  int
	Next *ListNode
}

type NodeIterator struct {
	cur   *ListNode
	value int
}

func NewNodeIterator(head *ListNode) *NodeIterator {
	return &NodeIterator{cur: head}
}

func (ls *NodeIterator) Next() bool {
	if ls.cur != nil {
		ls.value = ls.cur.Val
		ls.cur = ls.cur.Next
		return true
	}
	return false
}

func (ls *NodeIterator) Value() int {
	return ls.value
}

func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	var iters []Iterator[int]
	for _, list := range lists {
		if list != nil {
			iters = append(iters, NewNodeIterator(list))
		}
	}
	if len(iters) == 0 {
		return nil
	}

	less := func(a, b int) bool { return a < b }
	iter := HeapMerge(iters, less)

	dummy := &ListNode{}
	cur := dummy
	for iter.Next() {
		cur.Next = &ListNode{Val: iter.Value()}
		cur = cur.Next
	}
	return dummy.Next
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
	heap    *Heap[heapElement[E]]
	current E
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
	m.current = elem.value
	if elem.it.Next() {
		m.heap.Push(heapElement[E]{value: elem.it.Value(), it: elem.it})
	}
	return true
}

func (m *Merged[E]) Value() E {
	return m.current
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
