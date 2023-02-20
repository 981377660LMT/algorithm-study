// https://ei1333.github.io/library/structure/heap/persistent-leftist-heap.hpp
// clone 函数拷贝一份结点就可以持久化

// TODO 有问题
package main

import (
	"fmt"
	"strings"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	heapq := NewLeftistHeap(func(v1, v2 int) bool { return v1 < v2 }, true)
	pq1 := heapq.Alloc(0)
	for i := 0; i < len(nums); i++ {
		pq1 = heapq.Push(pq1, nums[i])
	}

	fmt.Println(pq1.String())
}

type V = int

type LeftistHeap struct {
	less         func(v1, v2 V) bool
	isPersistent bool
	pid          int
}

type Heap struct {
	Value       V
	Id          int
	height      int // 维持平衡
	left, right *Heap
}

func (h *Heap) String() string {
	var sb []string
	var dfs func(h *Heap)
	dfs = func(h *Heap) {
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

// less: 小根堆返回 v1 < v2, 大根堆返回 v1 > v2
// isPersistent: 是否持久化, 如果是, 则每次合并操作都会返回一个新的堆, 否则会直接修改原堆.
func NewLeftistHeap(less func(v1, v2 V) bool, isPersisitent bool) *LeftistHeap {
	return &LeftistHeap{less: less, isPersistent: isPersisitent}
}

func (lh *LeftistHeap) Alloc(value V) *Heap {
	res := &Heap{Value: value, Id: lh.pid, height: 1}
	lh.pid++
	return res
}

func (lh *LeftistHeap) Build(nums []V) []*Heap {
	res := make([]*Heap, len(nums))
	for i, num := range nums {
		res[i] = lh.Alloc(num)
	}
	return res
}

func (lh *LeftistHeap) Push(heap *Heap, value V) *Heap {
	return lh.Meld(heap, lh.Alloc(value))
}

func (lh *LeftistHeap) Pop(heap *Heap) *Heap {
	return lh.Meld(heap.left, heap.right)
}

func (lh *LeftistHeap) Top(heap *Heap) V {
	return heap.Value
}

// 合并两个堆,返回合并后的堆.
func (lh *LeftistHeap) Meld(heap1, heap2 *Heap) *Heap {
	if heap1 == nil {
		return heap2
	}
	if heap2 == nil {
		return heap1
	}
	heap1 = lh.clone(heap1)
	heap2 = lh.clone(heap2)
	if lh.less(heap2.Value, heap1.Value) {
		heap1, heap2 = heap2, heap1
	}

	if heap1.right != nil {
		heap1.right = lh.Meld(heap1.right, heap2)
	} else {
		heap1.right = heap2
	}
	if heap1.left == nil || heap1.left.height < heap1.right.height {
		heap1.left, heap1.right = heap1.right, heap1.left
	}
	heap1.height = 1
	if heap1.right != nil {
		heap1.height += heap1.right.height
	}
	return heap1
}

// 持久化,拷贝一份结点.
func (lh *LeftistHeap) clone(heap *Heap) *Heap {
	if heap == nil || !lh.isPersistent {
		return heap
	}
	res := &Heap{height: heap.height, Value: heap.Value, Id: lh.pid, left: heap.left, right: heap.right}
	lh.pid++
	return res
}

func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(value int) int {
	for ufa.parent[value] != value {
		ufa.parent[value] = ufa.parent[ufa.parent[value]]
		value = ufa.parent[value]
	}
	return value
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(value int) int {
	return ufa.rank[ufa.Find(value)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
