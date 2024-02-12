package main

import "fmt"

func main() {
	// dist [2 2 8 1000000000000000000 10] adjList [[[1 4]] [[0 4] [2 8] [9 1]] [[1 8] [3 3] [8 1]] [[2 3] [4 3] [7 4]] [[3 3] [5 1] [6 1]] [[4 1]] [[4 1]] [[3 4]] [[2 1]] [[1 1]]] points [5 6 7 8 9]
	n := 10
	edges := [][]int{{0, 1, 4}, {1, 2, 8}, {2, 3, 3}, {3, 4, 3}, {4, 5, 1}, {4, 6, 1}, {3, 7, 4}, {2, 8, 1}, {1, 9, 1}}
	points := []int{5, 6, 7, 8, 9}
	adjList := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
		adjList[v] = append(adjList[v], [2]int{u, w})
	}
	fmt.Println(MinDistToOther(adjList, points)) // [1 1 1]
	// dist, nearest := MinDistToOtherWithNearest(adjList, points)
	// fmt.Println(dist, nearest) // [1 1 5] [1 0 1]
}

const INF int = 1e18

// !给定一个无负权的带权图和一个点集，对点集内的每个点V，求V到点集内其他点距离的最小值.
// 换一种描述方法，给定一张图，有一些黑点和白点，对每个黑点，求出它到其他黑点的最近距离.
// 按照points中点的顺序返回答案.
//
// !以黑点集合为源做一次多源次短路dij，然后每个黑点的次短路就是答案.
// !这个次短路指的是“前两个不同节点到达这里的最短路”.
func MinDistToOther(adjList [][][2]int, points []int) (dist []int) {
	n := len(adjList)
	dist1, dist2 := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		dist1[i], dist2[i] = INF, INF
	}
	type H = struct {
		dist         int
		node, source int32
	}
	pq := NewHeap(func(a, b interface{}) bool { return a.(H).dist < b.(H).dist }, nil)
	for _, v := range points {
		dist1[v] = 0
		pq.Push(H{dist: 0, node: int32(v), source: int32(v)})
	}

	for pq.Len() > 0 {
		item := pq.Pop().(H)
		curDist, cur, curSource := item.dist, item.node, item.source
		// if curDist > dist2[cur] {
		// 	continue
		// }
		for _, e := range adjList[cur] {
			next, weight := int32(e[0]), e[1]
			nextDist := curDist + weight
			if nextDist < dist1[next] {
				dist1[next], nextDist = nextDist, dist1[next]
				pq.Push(H{dist1[next], next, curSource})
			}
			if nextDist < dist2[next] && next != curSource { // 非严格次短路
				dist2[next] = nextDist
				pq.Push(H{dist2[next], next, curSource})
			}
		}
	}

	fmt.Println(dist1, dist2)
	dist = make([]int, len(points))
	for i, v := range points {
		dist[i] = dist2[v]
	}
	return
}

// !给定一个无负权的带权图和一个点集，对点集内的每个点V，求V到点集内其他点距离的最小值，以及到V最近的点是谁.
func MinDistToOtherWithNearest(adjList [][][2]int, points []int) (dist []int, nearest []int) {
	n := len(adjList)
	dist1, dist2 := make([]int, n), make([]int, n)
	source1, source2 := make([]int32, n), make([]int32, n)
	for i := 0; i < n; i++ {
		dist1[i], dist2[i] = INF, INF
		source1[i], source2[i] = -1, -1
	}
	type H = struct {
		dist         int
		node, source int32
	}
	pq := NewHeap(func(a, b interface{}) bool { return a.(H).dist < b.(H).dist }, nil)
	for _, v := range points {
		v32 := int32(v)
		dist1[v32] = 0
		source1[v32] = v32
		pq.Push(H{dist: 0, node: v32, source: v32})
	}

	for pq.Len() > 0 {
		item := pq.Pop().(H)
		curDist, cur, curSource := item.dist, item.node, item.source
		if dist1[cur] != curDist && dist2[cur] != curDist {
			continue
		}
		for _, e := range adjList[cur] {
			next, weight := int32(e[0]), e[1]
			nextDist := curDist + weight
			if nextDist < dist1[next] {
				if source1[next] != next {
					dist2[next] = dist1[next]
					source2[next] = source1[next]
				}
				dist1[next] = nextDist
				source1[next] = curSource
				pq.Push(H{dist1[next], next, curSource})
			}
			if nextDist < dist2[next] && next != curSource { // 非严格次短路
				dist2[next] = nextDist
				source2[next] = curSource
				pq.Push(H{dist2[next], next, curSource})
			}
		}
	}

	dist = make([]int, len(points))
	nearest = make([]int, len(points))
	for i, v := range points {
		dist[i] = dist2[v]
		nearest[i] = int(source2[v])
	}
	return
}

type H = interface{}

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
