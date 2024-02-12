// 次短路计数
// https://www.acwing.com/solution/content/12246/

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

	var T int
	fmt.Fscan(in, &T)
	for t := 0; t < T; t++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		edges := make([][3]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
			edges[i][0]--
			edges[i][1]--
		}
		var start, end int
		fmt.Fscan(in, &start, &end)
		start--
		end--
		dist, count := Solve(n, edges, start)
		count1, count2 := count[end][0], count[end][1]
		if dist[end][0]+1 == dist[end][1] {
			count1 += count2
		}
		fmt.Fprintln(out, count1)
	}
}

const INF int = 1e18

// 求最短路、次短路的距离与路径数.
func Solve(
	n int, directedEdges [][3]int, start int,
) (dist, count [][2]int) {
	adjList := make([][][2]int, n)
	for _, e := range directedEdges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
	}
	dist, count = make([][2]int, n), make([][2]int, n)
	for i := 0; i < n; i++ {
		dist[i][0], dist[i][1] = INF, INF
	}
	dist[start][0], count[start][0] = 0, 1
	pq := NewHeap(func(a, b H) bool { return a[0] < b[0] }, nil)
	pq.Push(H{0, start, 0}) // (dist, node, type)
	for pq.Len() > 0 {
		item := pq.Pop()
		curDist, cur, curType := item[0], item[1], item[2]
		if dist[cur][curType] < curDist {
			continue
		}
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := curDist + weight
			if cand < dist[next][0] {
				dist[next][1], count[next][1] = dist[next][0], count[next][0]
				pq.Push(H{dist[next][1], next, 1})
				dist[next][0], count[next][0] = cand, count[cur][curType]
				pq.Push(H{dist[next][0], next, 0})
			} else if dist[next][0] == cand {
				count[next][0] += count[cur][curType]
			} else if cand < dist[next][1] {
				dist[next][1], count[next][1] = cand, count[cur][curType]
				pq.Push(H{dist[next][1], next, 1})
			} else if cand == dist[next][1] {
				count[next][1] += count[cur][curType]
			}
		}
	}
	return
}

type H = [3]int

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
