package main

import (
	"fmt"
	"strings"
)

func main() {
	ds := NewDifferSystem(3)
	ds.GreaterThanOrEqualTo(0, 1, 10)
	ds.LessThanOrEqualTo(1, 2, 2)
	fmt.Println(ds.HasSolution(3))
	fmt.Println(ds)
	fmt.Println(ds.PossibleSolutionOf(2))
	fmt.Println(ds.FindMinDifferenceBetween(0, 1))
	// fmt.Println(ds.FindMaxDifferenceBetween(0, 1))
}

const INF int = 2e18

// 差分约束. TODO: 习题练习
type DifferSystem struct {
	nodes       []*DifferNode
	queue       []*DifferNode
	n           int32
	allPos      bool
	hasSolution bool
}

func NewDifferSystem(n int32) *DifferSystem {
	res := &DifferSystem{
		nodes:  make([]*DifferNode, n),
		queue:  make([]*DifferNode, 0, n),
		n:      n,
		allPos: true,
	}
	for i := int32(0); i < n; i++ {
		node := NewNode()
		node.id = i
		res.nodes[i] = node
	}
	return res
}

// d[i] - d[j] <= d
func (ds *DifferSystem) LessThanOrEqualTo(i, j int32, d int) {
	ds.nodes[j].adj = append(ds.nodes[j].adj, NewNeighbor(ds.nodes[i], d))
	ds.allPos = ds.allPos && d >= 0
}

// d[i] - d[j] >= d
func (ds *DifferSystem) GreaterThanOrEqualTo(i, j int32, d int) {
	ds.LessThanOrEqualTo(j, i, -d)
}

// d[i] - d[j] = d
func (ds *DifferSystem) EqualTo(i, j int32, d int) {
	ds.GreaterThanOrEqualTo(i, j, d)
	ds.LessThanOrEqualTo(i, j, d)
}

// d[i] - d[j] < d
func (ds *DifferSystem) LessThan(i, j int32, d int) {
	ds.LessThanOrEqualTo(i, j, d-1)
}

// d[i] - d[j] > d
func (ds *DifferSystem) GreaterThan(i, j int32, d int) {
	ds.GreaterThanOrEqualTo(i, j, d+1)
}

func (ds *DifferSystem) HasSolution(i int32) bool {
	ds._prepare(0)
	for i := int32(0); i < ds.n; i++ {
		ds.nodes[i].inQueue = true
		ds.queue = append(ds.queue, ds.nodes[i])
	}
	ds.hasSolution = ds._spfa()
	return ds.hasSolution
}

// Find max(ai - aj), if INF is returned, it means no constraint between ai and aj.
func (ds *DifferSystem) FindMaxDifferenceBetween(i, j int32) int {
	ds.RunSince(j)
	return ds.nodes[i].dist
}

// Find min(ai - aj), if INF is returned, it means no constraint between ai and aj.
func (ds *DifferSystem) FindMinDifferenceBetween(i, j int32) int {
	r := ds.FindMaxDifferenceBetween(j, i)
	if r == INF {
		return INF
	}
	return -r
}

// After invoking this method, the value of i is max(ai - aj).
func (ds *DifferSystem) RunSince(j int32) bool {
	ds._prepare(INF)
	ds.queue = ds.queue[:0]
	ds.queue = append(ds.queue, ds.nodes[j])
	ds.nodes[j].dist = 0
	ds.nodes[j].inQueue = true
	ds.hasSolution = ds._spfa()
	return ds.hasSolution
}

func (ds *DifferSystem) PossibleSolutionOf(i int32) int {
	return ds.nodes[i].dist
}

func (ds *DifferSystem) Clear(n int32) {
	ds.n = n
	ds.allPos = true
	for i := int32(0); i < n; i++ {
		ds.nodes[i].adj = ds.nodes[i].adj[:0]
	}
}

func (ds *DifferSystem) _spfa() bool {
	if ds.allPos {
		ds._dijkstra()
		return true
	}
	for len(ds.queue) > 0 {
		head := ds.queue[0]
		ds.queue = ds.queue[1:]
		head.inQueue = false
		if head.times >= ds.n {
			return false
		}
		for _, edge := range head.adj {
			node := edge.next
			if node.dist <= edge.weight+head.dist {
				continue
			}
			node.dist = edge.weight + head.dist
			if node.inQueue {
				continue
			}
			node.times++
			node.inQueue = true
			ds.queue = append(ds.queue, node)
		}
	}
	return true
}

func (ds *DifferSystem) _prepare(initDist int) {
	ds.queue = ds.queue[:0]
	for i := int32(0); i < ds.n; i++ {
		node := ds.nodes[i]
		node.dist = initDist
		node.times = 0
		node.inQueue = false
	}
}

func (ds *DifferSystem) _dijkstra() {
	pq := NewErasableHeapGeneric[*DifferNode](
		func(a, b *DifferNode) bool {
			if a.dist == b.dist {
				return a.id < b.id
			}
			return a.dist < b.dist
		},
		ds.queue...,
	)
	for pq.Len() > 0 {
		head := pq.Pop()
		for _, e := range head.adj {
			if e.next.dist <= head.dist+e.weight {
				continue
			}
			pq.Erase(e.next)
			e.next.dist = head.dist + e.weight
			pq.Push(e.next)
		}
	}
}

func (ds *DifferSystem) String() string {
	sb := strings.Builder{}
	for i := int32(0); i < ds.n; i++ {
		for _, edge := range ds.nodes[i].adj {
			sb.WriteString(edge.String())
			sb.WriteByte('\n')
		}
	}
	sb.WriteString("-------------\n")
	if !ds.hasSolution {
		sb.WriteString("impossible")
	} else {
		for i := int32(0); i < ds.n; i++ {
			sb.WriteString(fmt.Sprintf("a%d=%d\n", i, ds.nodes[i].dist))
		}
	}
	return sb.String()
}

type Neighbor struct {
	next   *DifferNode
	weight int
}

func NewNeighbor(next *DifferNode, weight int) *Neighbor {
	return &Neighbor{next: next, weight: weight}
}

func (n *Neighbor) String() string {
	return fmt.Sprintf("next=%d, weight=%d", n.next.id, n.weight)
}

type DifferNode struct {
	adj     []*Neighbor
	dist    int
	inQueue bool
	times   int32
	id      int32
}

func NewNode() *DifferNode {
	return &DifferNode{}
}

type ErasableHeapGeneric[H comparable] struct {
	data   *HeapGeneric[H]
	erased *HeapGeneric[H]
}

func NewErasableHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *ErasableHeapGeneric[H] {
	return &ErasableHeapGeneric[H]{NewHeapGeneric(less, nums...), NewHeapGeneric(less)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeapGeneric[H]) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
}

func (h *ErasableHeapGeneric[H]) Push(value H) {
	h.data.Push(value)
	h.normalize()
}

func (h *ErasableHeapGeneric[H]) Pop() (value H) {
	value = h.data.Pop()
	h.normalize()
	return
}

func (h *ErasableHeapGeneric[H]) Peek() (value H) {
	value = h.data.Top()
	return
}

func (h *ErasableHeapGeneric[H]) Len() int {
	return h.data.Len()
}

func (h *ErasableHeapGeneric[H]) Clear() {
	h.data.Clear()
	h.erased.Clear()
}

func (h *ErasableHeapGeneric[H]) normalize() {
	for h.data.Len() > 0 && h.erased.Len() > 0 && h.data.Top() == h.erased.Top() {
		h.data.Pop()
		h.erased.Pop()
	}
}

type HeapGeneric[H comparable] struct {
	data []H
	less func(a, b H) bool
}

func NewHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *HeapGeneric[H] {
	nums = append(nums[:0:0], nums...)
	heap := &HeapGeneric[H]{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

func (h *HeapGeneric[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *HeapGeneric[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *HeapGeneric[H]) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *HeapGeneric[H]) Len() int { return len(h.data) }

func (h *HeapGeneric[H]) Clear() {
	h.data = h.data[:0]
}

func (h *HeapGeneric[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *HeapGeneric[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *HeapGeneric[H]) pushDown(root int) {
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
