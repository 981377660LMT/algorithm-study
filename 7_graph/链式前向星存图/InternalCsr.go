package main

import "fmt"

func main() {
	graph := NewInternalCsr(4, 4)
	graph.AddDirectedEdge(0, 1, 1)
	graph.AddDirectedEdge(0, 2, 1)
	graph.AddDirectedEdge(1, 3, 1)
	graph.AddDirectedEdge(2, 3, 1)
	dist := DijkstraInternalCsr(4, graph, 0)
	fmt.Println(1)
	fmt.Println(dist)
}

type InternalCsr struct {
	to     []int32
	next   []int32
	weight []int
	head   []int32
	eid    int32
}

func NewInternalCsr(n int, m int) *InternalCsr {
	res := &InternalCsr{
		to:     make([]int32, m),
		next:   make([]int32, m),
		weight: make([]int, m),
		head:   make([]int32, n),
		eid:    0,
	}
	for i := range res.head {
		res.head[i] = -1
	}
	return res
}

func (csr *InternalCsr) AddDirectedEdge(from, to int, weight int) {
	csr.to[csr.eid] = int32(to)
	csr.next[csr.eid] = csr.head[from]
	csr.weight[csr.eid] = weight
	csr.head[from] = csr.eid
	csr.eid++
}

func (csr *InternalCsr) EnumerateNeighbors(cur int, f func(next int, weight int)) {
	for i := csr.head[cur]; i != -1; i = csr.next[i] {
		f(int(csr.to[i]), csr.weight[i])
	}
}

func (csr *InternalCsr) Clear() {
	csr.eid = 0
	for i := range csr.head {
		csr.head[i] = -1
	}
}

func (csr *InternalCsr) Copy() *InternalCsr {
	clone := &InternalCsr{}
	clone.eid = csr.eid
	clone.to = make([]int32, len(csr.to))
	copy(clone.to, csr.to)
	clone.next = make([]int32, len(csr.next))
	copy(clone.next, csr.next)
	clone.weight = make([]int, len(csr.weight))
	copy(clone.weight, csr.weight)
	clone.head = make([]int32, len(csr.head))
	copy(clone.head, csr.head)
	return clone
}

const INF int = 1e18

type Neighbor struct{ to, weight int }

func DijkstraInternalCsr(n int, adjList *InternalCsr, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := nhp(func(a, b H) int {
		return a.dist - b.dist
	}, []H{{start, 0}})

	for pq.Len() > 0 {
		curNode := pq.Pop()
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		adjList.EnumerateNeighbors(cur, func(next, weight int) {
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		})
	}

	return
}

type H = struct{ node, dist int }

// Should return a number:
//
//	negative , if a < b
//	zero		 , if a == b
//	positive , if a > b
type Comparator = func(a, b H) int

func nhp(comparator Comparator, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{comparator: comparator, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data       []H
	comparator Comparator
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		return
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
