// https://www.luogu.com.cn/problem/P2865
// P2865 [USACO06NOV] Roadblocks G
// !1.如果 dist1[next] > dist1[cur] + d(cur,next)，则更新 dist1[next]
// !2.如果 dist1[next] < dist1[cur] + d(cur,next) < dist2[next]，则更新 dist2[next]
//    注意不能取等，否则dist2[cur]和dist1[cur]可能相等

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	_, dist2 := MinDist(n, edges, 0)
	fmt.Fprintln(out, dist2[n-1])
}

const INF int = 1e18

// 给定一个无向带权图,求start到其他点的最短路和严格次短路.
// 不存在则为INF.
func MinDist(n int, edges [][3]int, start int) (dist1, dist2 []int) {
	adjList := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
		adjList[v] = append(adjList[v], [2]int{u, w})
	}
	dist1, dist2 = make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		dist1[i], dist2[i] = INF, INF
	}
	dist1[start] = 0
	pq := NewHeap(func(a, b H) bool { return a[0] < b[0] }, nil)
	pq.Push(H{0, start})
	for pq.Len() > 0 {
		item := pq.Pop()
		curDist, cur := item[0], item[1]
		if curDist > dist2[cur] { // !注意是dist2
			continue
		}
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := curDist + weight
			if cand < dist1[next] {
				dist1[next], cand = cand, dist1[next]
				pq.Push(H{dist1[next], next})
			}
			if dist1[next] < cand && cand < dist2[next] { // dist1[next] < cand ，严格次短路
				dist2[next] = cand
				pq.Push(H{dist2[next], next})
			}
		}
	}

	return
}

type H = [2]int

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
