// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/heap.go

package main

import (
	"container/heap"
	"fmt"
)

func main() {
	mh := NewHeap[int](0, func(a, b int) bool { return a < b })
	mh.Push(1)
	mh.Push(3)
	mh.Push(2)
	fmt.Println(mh.Pop()) // 1
	fmt.Println(mh.Pop()) // 2
	fmt.Println(mh.Pop()) // 3
}

type Heap[T any] struct{ wrapper *wrapper[T] }

func NewHeap[T any](initCapacity int32, less func(a, b T) bool) *Heap[T] {
	if initCapacity < 0 {
		initCapacity = 0
	}
	wrapper := &wrapper[T]{make([]T, 0, initCapacity), less}
	return &Heap[T]{wrapper}
}

func (h *Heap[T]) Push(v T)   { heap.Push(h.wrapper, v) }
func (h *Heap[T]) Pop() T     { return heap.Pop(h.wrapper).(T) }
func (h *Heap[T]) Top() T     { return h.wrapper.data[0] }
func (h *Heap[T]) Len() int32 { return int32(len(h.wrapper.data)) }
func (h *Heap[T]) Init()      { heap.Init(h.wrapper) }

// replace 弹出并返回堆顶，同时将 v 入堆
// 需保证 h 非空
func (h *Heap[T]) Replace(v T) T {
	top := h.Top()
	h.wrapper.data[0] = v
	heap.Fix(h.wrapper, 0)
	return top
}

// pushPop 先将 v 入堆，然后弹出并返回堆顶
func (h *Heap[T]) PushPop(v T) T {
	data, less := h.wrapper.data, h.wrapper.less
	if len(data) > 0 && less(data[0], v) {
		v, data[0] = data[0], v
		heap.Fix(h.wrapper, 0)
	}
	return v
}

var _ heap.Interface = &wrapper[any]{}

type wrapper[T any] struct {
	data []T
	less func(a, b T) bool
}

func (w *wrapper[T]) Len() int           { return len(w.data) }
func (w *wrapper[T]) Less(i, j int) bool { return w.less(w.data[i], w.data[j]) }
func (w *wrapper[T]) Swap(i, j int)      { w.data[i], w.data[j] = w.data[j], w.data[i] }
func (w *wrapper[T]) Push(v any)         { w.data = append(w.data, v.(T)) }
func (w *wrapper[T]) Pop() any {
	res := w.data[len(w.data)-1]
	w.data = w.data[:len(w.data)-1]
	return res
}
