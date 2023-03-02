// https://ei1333.github.io/library/structure/heap/persistent-leftist-heap.hpp
// clone 函数拷贝一份结点就可以持久化

package main

import (
	"fmt"
	"strings"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	heapq := NewLeftistHeap(func(v1, v2 int) bool { return v1 < v2 }, true)
	pq1 := heapq.Alloc(0, -1)
	for i := 0; i < len(nums); i++ {
		pq1 = heapq.Push(pq1, nums[i], i)
	}

	fmt.Println(pq1.String())
}

type V = int

type leftistHeap struct {
	less         func(v1, v2 V) bool
	isPersistent bool
}

type heap struct {
	Value       V
	Id          int
	height      int // 维持平衡
	left, right *heap
}

// less: 小根堆返回 v1 < v2, 大根堆返回 v1 > v2
// isPersistent: 是否持久化, 如果是, 则每次合并操作都会返回一个新的堆, 否则会直接修改原堆.
func NewLeftistHeap(less func(v1, v2 V) bool, isPersisitent bool) *leftistHeap {
	return &leftistHeap{less: less, isPersistent: isPersisitent}
}

func (lh *leftistHeap) Alloc(value V, id int) *heap {
	res := &heap{Value: value, Id: id, height: 1}
	return res
}

func (lh *leftistHeap) Build(nums []V) []*heap {
	res := make([]*heap, len(nums))
	for i, num := range nums {
		res[i] = lh.Alloc(num, i)
	}
	return res
}

func (lh *leftistHeap) Push(heap *heap, value V, id int) *heap {
	return lh.Meld(heap, lh.Alloc(value, id))
}

func (lh *leftistHeap) Pop(heap *heap) *heap {
	return lh.Meld(heap.left, heap.right)
}

func (lh *leftistHeap) Top(heap *heap) V {
	return heap.Value
}

// 合并两个堆,返回合并后的堆.
func (lh *leftistHeap) Meld(heap1, heap2 *heap) *heap {
	if heap1 == nil {
		return heap2
	}
	if heap2 == nil {
		return heap1
	}
	if lh.less(heap2.Value, heap1.Value) {
		heap1, heap2 = heap2, heap1
	}
	heap1 = lh.clone(heap1)
	heap1.right = lh.Meld(heap1.right, heap2)
	if heap1.left == nil || heap1.left.height < heap1.right.height {
		heap1.left, heap1.right = heap1.right, heap1.left
	}
	heap1.height = 1
	if heap1.right != nil {
		heap1.height += heap1.right.height
	}
	return heap1
}

func (h *heap) String() string {
	var sb []string
	var dfs func(h *heap)
	dfs = func(h *heap) {
		if h == nil {
			return
		}
		sb = append(sb, fmt.Sprintf("%d", h.Value))
		dfs(h.left)
		dfs(h.right)
	}
	dfs(h)
	return strings.Join(sb, " ")
}

// 持久化,拷贝一份结点.
func (lh *leftistHeap) clone(h *heap) *heap {
	if h == nil || !lh.isPersistent {
		return h
	}
	res := &heap{height: h.height, Value: h.Value, Id: h.Id, left: h.left, right: h.right}
	return res
}
