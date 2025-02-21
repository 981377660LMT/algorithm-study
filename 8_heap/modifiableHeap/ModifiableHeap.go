// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/heap.go

package main

import (
	"container/heap"
	"fmt"
)

func main() {
	mh := NewModifiableHeap[int](0, func(a, b int) bool { return a < b })

	{
		token := mh.Push(1)
		mh.Push(2)
		mh.Push(3)
		mh.Modify(token, 4)
		fmt.Println(mh.Pop().value) // 2
		fmt.Println(mh.Pop().value) // 3
		fmt.Println(mh.Pop().value) // 4
	}

	{
		mh.Push(3)
		token := mh.Push(2)
		mh.Remove(token)
		fmt.Println(mh.Top().value) // 3
	}

}

type ModifiableHeap[T any] struct{ wrapper *pq[T] }

// 支持修改、删除指定元素的堆
// 用法：调用 push 会返回一个 *viPair 指针，记作 p
// 将 p 存于他处（如 slice 或 map），可直接在外部修改 p.v 后调用 fix(p.hi)，从而做到修改堆中指定元素
// heap.Remove(p.hi) 通过 swap_remove 实现，复杂度 O(logn).
func NewModifiableHeap[T any](initCapacity int32, less func(a, b T) bool) *ModifiableHeap[T] {
	if initCapacity < 0 {
		initCapacity = 0
	}
	wrapper := &pq[T]{make([]*viPair[T], 0, initCapacity), less}
	return &ModifiableHeap[T]{wrapper}
}

func (h *ModifiableHeap[T]) Push(v T) *viPair[T] {
	p := &viPair[T]{v, int32(len(h.wrapper.data))}
	heap.Push(h.wrapper, p)
	return p
}
func (h *ModifiableHeap[T]) Pop() *viPair[T]     { return heap.Pop(h.wrapper).(*viPair[T]) }
func (h *ModifiableHeap[T]) Top() *viPair[T]     { return h.wrapper.data[0] }
func (h *ModifiableHeap[T]) Len() int32          { return int32(len(h.wrapper.data)) }
func (h *ModifiableHeap[T]) Init()               { heap.Init(h.wrapper) }
func (h *ModifiableHeap[T]) Remove(p *viPair[T]) { heap.Remove(h.wrapper, int(p.heapIndex)) }
func (h *ModifiableHeap[T]) Modify(p *viPair[T], v T) {
	p.value = v
	heap.Fix(h.wrapper, int(p.heapIndex))
}
func (h *ModifiableHeap[T]) Clear() {
	h.wrapper.Clear()
}

var _ heap.Interface = &pq[any]{}

type viPair[T any] struct {
	value     T
	heapIndex int32
}

type pq[T any] struct {
	data []*viPair[T]
	less func(a, b T) bool
}

func (w *pq[T]) Len() int           { return len(w.data) }
func (w *pq[T]) Less(i, j int) bool { return w.less(w.data[i].value, w.data[j].value) }
func (w *pq[T]) Swap(i, j int) {
	data := w.data
	data[i], data[j] = data[j], data[i]
	data[i].heapIndex = int32(i)
	data[j].heapIndex = int32(j)
}
func (w *pq[T]) Push(v any) { w.data = append(w.data, v.(*viPair[T])) }
func (w *pq[T]) Pop() any {
	res := w.data[len(w.data)-1]
	w.data = w.data[:len(w.data)-1]
	return res
}
func (w *pq[T]) Clear() {
	w.data = w.data[:0]
}
