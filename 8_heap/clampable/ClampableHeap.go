// RangeChminRangeSum
// RangeChmaxRangeSum
// ChminHeap/ChmaxHeap/ClampHeap
// 设计一个高效的数据结构，支持三种操作：
// 1. 添加一个数；2.容器内所有数与x取min；3.求容器内元素之和
//
// 一般配合后缀数组的height数组使用.

package main

import "fmt"

func main() {
	C := NewClampableHeap(false)
	C.Add(1)
	C.Add(2)
	C.Add(3)
	fmt.Println(C.Sum()) // 6
	C.Clamp(2)
	fmt.Println(C.Sum()) // 5
	C.Clamp(1)
	C.Add(2)
	fmt.Println(C.Sum()) // 5

	C = NewClampableHeap(true)
	C.Add(1)
	C.Add(2)
	C.Add(3)
	fmt.Println(C.Sum()) // 6
	C.Clamp(2)
	fmt.Println(C.Sum()) // 7
	C.Clamp(3)
	C.Add(2)
	fmt.Println(C.Sum()) // 11
}

type ClampableHeap struct {
	clampMin bool
	total    int
	count    int
	heap     *Heap
}

// clampMin：
//  为true时，调用Clamp(x)后，容器内所有数最小值被截断(小于x的数变成x)；
//  为false时，调用Clamp(x)后，容器内所有数最大值被截断(大于x的数变成x).
//  如果需要同时支持两种操作，可以使用双端堆.
func NewClampableHeap(clampMin bool) *ClampableHeap {
	var less func(a, b H) bool
	if clampMin {
		less = func(a, b H) bool { return a.value <= b.value }
	} else {
		less = func(a, b H) bool { return a.value >= b.value }
	}
	return &ClampableHeap{clampMin: clampMin, heap: NewHeap(less)}
}

func (h *ClampableHeap) Add(x int) {
	h.heap.Push(H{value: x, count: 1})
	h.total += x
	h.count++
}

func (h *ClampableHeap) Clamp(x int) {
	newCount := 0
	if h.clampMin {
		for h.heap.Len() > 0 {
			top := h.heap.Top()
			if top.value > x {
				break
			}
			h.heap.Pop()
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	} else {
		for h.heap.Len() > 0 {
			top := h.heap.Top()
			if top.value < x {
				break
			}
			h.heap.Pop()
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	}
	h.total += x * newCount
	h.heap.Push(H{value: x, count: int32(newCount)})
}

func (h *ClampableHeap) Sum() int {
	return h.total
}

func (h *ClampableHeap) Len() int {
	return h.count
}

func (h *ClampableHeap) Clear() {
	h.heap.Clear()
	h.total = 0
	h.count = 0
}

type H = struct {
	value int
	count int32
}

func NewHeap(less func(a, b H) bool, nums ...H) *Heap {
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

func (h *Heap) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) Clear() { h.data = h.data[:0] }

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
