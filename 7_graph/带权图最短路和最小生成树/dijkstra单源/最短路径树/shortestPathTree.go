// https://www.luogu.com.cn/blog/LawrenceSivan/cf545e-paths-and-trees-zui-duan-lu-jing-shu-post
// 求出最短路径树
// !给定一张带正权的无向图和一个源点，求边权和最小的最短路径树。
// 输出该树权值和和和每条边的编号

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}
	var start int
	fmt.Fscan(in, &start)
	start--
	weightSum, edgeId := CF545E(n, edges, start)
	fmt.Fprintln(out, weightSum)
	for _, id := range edgeId {
		fmt.Fprint(out, id+1, " ")
	}
}

// 给定一张带正权的无向图和一个源点，求边权和最小的最短路径树。
// 输出该树权值和和和每条边的编号.
func CF545E(n int, edges [][3]int, start int) (weightSum int, edgeId []int) {
	adjList := make([][][3]int, n)
	for i, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [3]int{v, w, i})
		adjList[v] = append(adjList[v], [3]int{u, w, i})
	}
	_, _, preEdge := DijkstraShortestPathTree(n, adjList, start)
	for i := 0; i < n; i++ {
		if i == start {
			continue
		}
		pre := preEdge[i]
		weightSum += edges[pre][2]
		edgeId = append(edgeId, pre)
	}
	return
}

// dijkstra求出路径上每个点的前驱点和前驱边.
// 其中邻接表的每个元素是一个三元组，分别是邻接点，边权，边的编号.
func DijkstraShortestPathTree(n int, adjList [][][3]int, start int) (dist []int, preVertex []int, preEdge []int) {
	dist, preVertex, preEdge = make([]int, n), make([]int, n), make([]int, n)
	for i := range dist {
		dist[i] = INF
		preVertex[i] = -1
		preEdge[i] = -1
	}
	dist[start] = 0
	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, []H{{0, start}})
	for pq.Len() > 0 {
		cur := pq.Pop()
		curDist, curVertex := cur.dist, cur.vertex
		if dist[curVertex] < curDist {
			continue
		}
		for _, e := range adjList[curVertex] {
			next, weight, edgeId := e[0], e[1], e[2]
			if tmp := dist[curVertex] + weight; tmp < dist[next] {
				dist[next] = tmp
				preEdge[next] = edgeId
				preVertex[next] = curVertex
				pq.Push(H{tmp, next})
			} else if tmp == dist[next] { // 在最短路相等的情况下，扩展到同一个节点，后出堆的点连的边权值一定更小
				preEdge[next] = edgeId
				preVertex[next] = curVertex
			}
		}
	}
	return
}

type H = struct{ dist, vertex int }

func NewHeap(less func(a, b H) bool, data []H) *Heap {
	data = append(data[:0:0], data...)
	heap := &Heap{less: less, data: data}
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
	if h.Len() == 0 {
		panic("heap is empty")
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
