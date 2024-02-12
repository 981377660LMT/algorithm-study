package main

import "fmt"

func main() {
	n := 10
	edges := [][]int{{0, 1, 4}, {1, 2, 8}, {2, 3, 3}, {3, 4, 3}, {4, 5, 1}, {4, 6, 1}, {3, 7, 4}, {2, 8, 1}, {1, 9, 1}}
	points := []int{5, 6, 7, 8, 9}
	adjList := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
		adjList[v] = append(adjList[v], [2]int{u, w})
	}
	fmt.Println(MinDistToOther(adjList, points))
}

const INF int = 1e18

// !给定一个无负权的带权图和一个点集，对点集内的每个点V，求V到点集内其他点距离的最小值,以及到V最近的点是谁.
// 换一种描述方法，给定一张图，有一些黑点和白点，对每个黑点，求出它到其他黑点的最近距离.
// 按照points中点的顺序返回答案.
//
// !以黑点集合为源做一次多源次短路dij，然后每个黑点的次短路就是答案.
// 注意次短路的出发点不为能自己.
// 类似 abc245G - Foreign Friends-简洁写法.go
func MinDistToOther(adjList [][][2]int, points []int) (dist []int, nearest []int) {
	n := len(adjList)
	dist = make([]int, n)
	source1, source2 := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		source1[i], source2[i] = -1, -1
	}

	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, nil)
	for _, v := range points {
		pq.Push(H{dist: 0, node: v, source: v})
	}

	for pq.Len() > 0 {
		item := pq.Pop()
		curDist, cur, curSource := item.dist, item.node, item.source
		if curSource == source1[cur] || curSource == source2[cur] {
			continue
		}
		if source1[cur] == -1 {
			source1[cur] = curSource
		} else if source2[cur] == -1 {
			source2[cur] = curSource
		} else {
			continue
		}

		if curSource != cur { // 出发点不为自己时，更新距离
			dist[cur] = min(dist[cur], curDist)
		}
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			nextDist := curDist + weight
			pq.Push(H{nextDist, next, curSource})
		}
	}

	nearest = source2
	for i, v := range points {
		dist[i] = dist[v]
		nearest[i] = nearest[v]
	}
	dist = dist[:len(points)]
	nearest = nearest[:len(points)]
	return
}

type H = struct{ dist, node, source int }

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
