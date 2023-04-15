// https://ei1333.github.io/library/structure/heap/persistent-leftist-heap.hpp
// 可持久化可并堆/可持久化左偏树/可持久化堆
// clone 函数拷贝一份结点就可以持久化

// Build(nums []E) []*SkewHeapNode
// Alloc(key E, index int) *SkewHeapNode

// Meld(x, y *SkewHeapNode) *SkewHeapNode

// Push(t *SkewHeapNode, key E, index int) *SkewHeapNode
// Pop(t *SkewHeapNode) *SkewHeapNode
// Top(t *SkewHeapNode) E

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

func demo() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	heapq := NewLeftistHeap(func(v1, v2 int) bool { return v1 < v2 }, true)
	pq1 := heapq.Alloc(0, -1)
	for i := 0; i < len(nums); i++ {
		pq1 = heapq.Push(pq1, nums[i], i)
	}
	for i := 0; i < len(nums); i++ {
		pq1 = heapq.Pop(pq1)
		fmt.Println(pq1.Value)
	}
	fmt.Println(pq1.String())
}

func main() {
	// check with brute Force
	n := rand.Intn(100) + 5
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Intn(100)
	}

	heapq := NewLeftistHeap(func(v1, v2 int) bool { return v1 < v2 }, true)
	heaps1 := heapq.Build(nums)
	heaps2 := make([][]int, n)
	for i := 0; i < n; i++ {
		heaps2[i] = []int{nums[i]}
	}

	// top
	for i := 0; i < n; i++ {
		if heapq.Top(heaps1[i]) != heaps2[i][0] {
			panic("error")
		}
	}

	// random meld
	for i := 0; i < 10; i++ {
		var i1, i2 int
		for i1 != i2 {
			i1 = rand.Intn(n)
			i2 = rand.Intn(n)
		}
		heaps1[i1] = heapq.Meld(heaps1[i1], heaps1[i2])
		heaps2[i1] = append(heaps2[i1], heaps2[i2]...)
		sort.Ints(heaps2[i1])
		if heapq.Top(heaps1[i1]) != heaps2[i1][0] {
			panic("error")
		}
	}

	// meld then pop
	for i := 0; i < n; i++ {
		// push
		for j := 0; j < 10; j++ {
			v := rand.Intn(100)
			heaps1[i] = heapq.Push(heaps1[i], v, i)
			heaps2[i] = append(heaps2[i], v)
			sort.Ints(heaps2[i])
			if heapq.Top(heaps1[i]) != heaps2[i][0] {
				panic("error")
			}
		}
		for len(heaps2[i]) > 0 {
			if heapq.Top(heaps1[i]) != heaps2[i][0] {
				panic("error")
			}
			heaps1[i] = heapq.Pop(heaps1[i])
			heaps2[i] = heaps2[i][1:]
		}
	}
}

type V = int

type leftistHeap struct {
	less         func(v1, v2 V) bool
	isPersistent bool
}

type heap struct {
	Value       V
	Id          int
	rank        int // 维持平衡
	left, right *heap
}

// less: 小根堆返回 v1 < v2, 大根堆返回 v1 > v2
// isPersistent: 是否持久化, 如果是, 则每次合并操作都会返回一个新的堆, 否则会直接修改原堆.
func NewLeftistHeap(less func(v1, v2 V) bool, isPersisitent bool) *leftistHeap {
	return &leftistHeap{less: less, isPersistent: isPersisitent}
}

func (lh *leftistHeap) Alloc(value V, id int) *heap {
	res := &heap{Value: value, Id: id, rank: 1}
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

// 合并两个堆,返回合并后的堆.(是否会修改原堆取决于 isPersistent)
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
	if heap1.left == nil || heap1.left.rank < heap1.right.rank {
		heap1.left, heap1.right = heap1.right, heap1.left
	}
	heap1.rank = 1
	if heap1.right != nil {
		heap1.rank += heap1.right.rank
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
	res := &heap{rank: h.rank, Value: h.Value, Id: h.Id, left: h.left, right: h.right}
	return res
}
