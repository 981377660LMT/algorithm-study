package main

import (
	"fmt"
	"time"
)

func main() {
	mh := NewHeap[int](func(a, b int) bool { return a < b })
	time1 := time.Now()
	for i := 0; i < 1000000; i++ {
		mh.Push(i)
	}
	for i := 0; i < 1000000; i++ {
		mh.Pop()
	}
	fmt.Println(time.Since(time1)) // 0.03
}

type IHeap[T any] interface {
	Push(value T)
	Pop() T
	Replace(v T) T
	PushPop(v T) T

	Top() T
	Len() int
}

var _ IHeap[any] = (*Heap[any])(nil)

func NewHeap[H any](less func(a, b H) bool, nums ...H) *Heap[H] {
	res := &Heap[H]{less: less, data: append(nums[:0:0], nums...)}
	if len(nums) > 1 {
		res.heapify()
	}
	return res
}

func NewHeapWithCapacity[H any](less func(a, b H) bool, capacity int32, nums ...H) *Heap[H] {
	if n := int32(len(nums)); capacity < n {
		capacity = n
	}
	res := &Heap[H]{less: less, data: make([]H, 0, capacity)}
	res.data = append(res.data, nums...)
	if len(nums) > 1 {
		res.heapify()
	}
	return res
}

type Heap[H any] struct {
	less func(a, b H) bool
	data []H
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.up(h.Len() - 1)
}

func (h *Heap[H]) Pop() H {
	n := h.Len() - 1
	h.data[0], h.data[n] = h.data[n], h.data[0]
	h.down(0, n)
	res := h.data[n]
	h.data = h.data[:n]
	return res
}

func (h *Heap[H]) Top() H {
	return h.data[0]
}

// replace 弹出并返回堆顶，同时将 v 入堆.
// 需保证 h 非空.
func (h *Heap[H]) Replace(v H) H {
	top := h.Top()
	h.data[0] = v
	h.fix(0)
	return top
}

// pushPop 先将 v 入堆，然后弹出并返回堆顶.
func (h *Heap[H]) PushPop(v H) H {
	data, less := h.data, h.less
	if len(data) > 0 && less(data[0], v) {
		v, data[0] = data[0], v
		h.fix(0)
	}
	return v
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.down(i, n)
	}
}

func (h *Heap[H]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(h.data[j], h.data[i]) {
			break
		}
		h.data[i], h.data[j] = h.data[j], h.data[i]
		j = i
	}
}

func (h *Heap[H]) down(i0, n int) bool {
	i := i0
	for {
		j1 := (i << 1) | 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(h.data[j2], h.data[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(h.data[j], h.data[i]) {
			break
		}
		h.data[i], h.data[j] = h.data[j], h.data[i]
		i = j
	}
	return i > i0
}

func (h *Heap[H]) fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}
