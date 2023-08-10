package main

// https://leetcode.cn/problems/design-graph-with-shortest-path-calculator/
type Graph struct {
	g  [][][2]int
	rg [][][2]int
}

func Constructor(n int, edges [][]int) Graph {
	g, rg := make([][][2]int, n), make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		g[u] = append(g[u], [2]int{v, w})
		rg[v] = append(rg[v], [2]int{u, w})
	}
	return Graph{g, rg}
}

func (this *Graph) AddEdge(edge []int) {
	u, v, w := edge[0], edge[1], edge[2]
	this.g[u] = append(this.g[u], [2]int{v, w})
	this.rg[v] = append(this.rg[v], [2]int{u, w})
}

func (this *Graph) ShortestPath(node1 int, node2 int) int {
	return FastDijkstra(len(this.g), this.g, this.rg, node1, node2)
}

//
//
// 双向dijkstra 求start到end的最短路.
func FastDijkstra(n int, g, rg [][][2]int, start, end int) int {
	if start == end {
		return 0
	}
	dist := make([]int, n)
	drev := make([]int, n)
	for i := range dist {
		dist[i] = -1
		drev[i] = -1
	}
	dist[start] = 0
	drev[end] = 0

	pq1 := NewHeap(func(a, b H) bool { return a[0] < b[0] })
	pq2 := NewHeap(func(a, b H) bool { return a[0] < b[0] })
	pq1.Push(H{0, start})
	pq2.Push(H{0, end})

	res := -1
	for pq1.Len() > 0 && pq2.Len() > 0 {
		d1 := pq1.Top()[0]
		d2 := pq2.Top()[0]
		if res >= 0 && d1+d2 >= res {
			break
		}
		if d1 <= d2 {
			u := pq1.Pop()[1]
			if dist[u] > d1 {
				continue
			}
			for _, e := range g[u] {
				v, w := e[0], e[1]
				cand := dist[u] + w
				if dist[v] >= 0 && dist[v] <= cand {
					continue
				}
				dist[v] = cand
				pq1.Push(H{dist[v], v})
				if drev[v] >= 0 {
					nu := dist[v] + drev[v]
					if res < 0 || res > nu {
						res = nu
					}
				}
			}
		} else {
			u := pq2.Pop()[1]
			if drev[u] > d2 {
				continue
			}
			for _, e := range rg[u] {
				v, w := e[0], e[1]
				cand := drev[u] + w
				if drev[v] >= 0 && drev[v] <= cand {
					continue
				}
				drev[v] = cand
				pq2.Push(H{drev[v], v})
				if dist[v] >= 0 {
					nu := dist[v] + drev[v]
					if res < 0 || res > nu {
						res = nu
					}
				}
			}
		}
	}

	return res
}

type H = [2]int

func NewHeap(less func(a, b H) bool, nums ...H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Top() (value H) {
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
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
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
