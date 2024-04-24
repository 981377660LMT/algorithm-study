package main

import "fmt"

func main() {
	D := NewTopoLogicalSortOfDAG(6)
	D.AddDirectedEdge(0, 1)
	D.AddDirectedEdge(1, 2)
	D.AddDirectedEdge(1, 3)
	D.AddDirectedEdge(2, 4)
	D.AddDirectedEdge(3, 4)
	D.AddDirectedEdge(4, 5)
	D.Sort()
	fmt.Println(D.TopoOrder)
	fmt.Println(D.LongestDistance())
}

type TopoLogicalSortOfDAG struct {
	AdjList   [][]int32
	TopoOrder []int32 // 字典序最小拓扑序列
	n         int32
}

func NewTopoLogicalSortOfDAG(n int32) *TopoLogicalSortOfDAG {
	return &TopoLogicalSortOfDAG{
		AdjList: make([][]int32, n),
		n:       n,
	}
}

func NewTopoLogicalSortOfDAGFromAdjList(adjList [][]int32) *TopoLogicalSortOfDAG {
	return &TopoLogicalSortOfDAG{
		AdjList: adjList,
		n:       int32(len(adjList)),
	}
}

func (t *TopoLogicalSortOfDAG) AddDirectedEdge(u, v int32) {
	t.AdjList[u] = append(t.AdjList[u], v)
}

// O(E + V log V)
func (t *TopoLogicalSortOfDAG) Sort() bool {
	indegree := make([]int32, t.n)
	for u := int32(0); u < t.n; u++ {
		for _, v := range t.AdjList[u] {
			indegree[v]++
		}
	}

	pq := NewHeap(func(a, b int32) bool { return a < b }, nil)
	used := make([]bool, t.n)
	for u := int32(0); u < t.n; u++ {
		if indegree[u] == 0 {
			pq.Push(u)
			used[u] = true
		}
	}

	t.TopoOrder = make([]int32, 0, t.n)
	for pq.Len() > 0 {
		u := pq.Pop()
		t.TopoOrder = append(t.TopoOrder, u)
		for _, v := range t.AdjList[u] {
			indegree[v]--
			if indegree[v] == 0 && !used[v] {
				used[v] = true
				pq.Push(v)
			}
		}
	}

	return len(t.TopoOrder) == int(t.n)
}

// DAG最长路径.
// 如果拓扑序列唯一，则返回值为n-1，原图为含有哈密顿路径的图.
func (t *TopoLogicalSortOfDAG) LongestDistance() int32 {
	res := int32(0)
	dist := make([]int32, t.n)
	for _, u := range t.TopoOrder {
		for _, v := range t.AdjList[u] {
			dist[v] = max32(dist[v], dist[u]+1)
			res = max32(res, dist[v])
		}
	}
	return res
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
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
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
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
