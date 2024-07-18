// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/heap.go
// !deprecated, slow

package main

import (
	"container/heap"
	"fmt"
	"time"
)

func main() {
	mh := NewHeapHeapStandard[int](0, func(a, b int) bool { return a < b })
	mh.Push(1)
	mh.Push(3)
	mh.Push(2)
	fmt.Println(mh.Pop()) // 1
	fmt.Println(mh.Pop()) // 2
	fmt.Println(mh.Pop()) // 3

	time1 := time.Now()
	for i := 0; i < 1000000; i++ {
		mh.Push(i)
	}
	for i := 0; i < 1000000; i++ {
		mh.Pop()
	}
	fmt.Println(time.Since(time1)) // 159.028834ms
}

type HeapStandard[T any] struct{ wrapper *wrapper[T] }

func NewHeapHeapStandard[T any](initCapacity int32, less func(a, b T) bool) *HeapStandard[T] {
	if initCapacity < 0 {
		initCapacity = 0
	}
	wrapper := &wrapper[T]{make([]T, 0, initCapacity), less}
	return &HeapStandard[T]{wrapper}
}

func (h *HeapStandard[T]) Push(v T)   { heap.Push(h.wrapper, v) }
func (h *HeapStandard[T]) Pop() T     { return heap.Pop(h.wrapper).(T) }
func (h *HeapStandard[T]) Top() T     { return h.wrapper.data[0] }
func (h *HeapStandard[T]) Len() int32 { return int32(len(h.wrapper.data)) }
func (h *HeapStandard[T]) Init()      { heap.Init(h.wrapper) }

// replace 弹出并返回堆顶，同时将 v 入堆
// 需保证 h 非空
func (h *HeapStandard[T]) Replace(v T) T {
	top := h.Top()
	h.wrapper.data[0] = v
	heap.Fix(h.wrapper, 0)
	return top
}

// pushPop 先将 v 入堆，然后弹出并返回堆顶
func (h *HeapStandard[T]) PushPop(v T) T {
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
