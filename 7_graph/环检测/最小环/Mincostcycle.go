// https://yukicoder.me/problems/no/1320
// n,m<=2000,wi<=1e9

// !枚举每条边删除：`O(E*(V+E)*logV)` 找最小环,不存在返回INF

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

	var directed, n, m int
	fmt.Fscan(in, &directed, &n, &m)
	edges := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		edges = append(edges, Edge{u, v, w})
	}
	res := MincostCycle(n, edges, directed == 1)
	if res == INF {
		res = -1
	}
	fmt.Fprintln(out, res)
}

const INF int = 1e18

type Edge struct{ from, to, weight int }

// 返回最小环的权值和. 不存在时返回INF.
//  !把每条边断开,然后求从断开的边的to到from的最短路.
func MincostCycle(n int, edges []Edge, directed bool) int {
	adjList := make([][][2]int, n)
	maxWeight := 0
	for _, e := range edges {
		u, v, w := e.from, e.to, e.weight
		if w > maxWeight {
			maxWeight = w
		}
		adjList[u] = append(adjList[u], [2]int{v, w})
		if !directed {
			adjList[v] = append(adjList[v], [2]int{u, w})
		}
	}
	res := INF
	for _, e := range edges {
		from_, to, weight := e.from, e.to, e.weight
		var dist int
		if maxWeight <= 1 {
			dist = bfs01Point(n, adjList, to, from_)
		} else {
			dist = dijkstraPoint(n, adjList, to, from_)
		}
		cand := weight + dist
		if cand < res {
			res = cand
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bfs01Point(n int, adjList [][][2]int, start, target int) int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}

	queue := Deque{}
	queue.Append(start)
	dist[start] = 0
	for !queue.Empty() {
		cur := queue.PopLeft()
		if cur == target {
			return dist[cur]
		}
		for _, e := range adjList[cur] {
			next, cost := e[0], e[1]
			if (next == start && cur == target) || (next == target && cur == start) {
				continue
			}
			if cand := dist[cur] + cost; dist[next] > cand {
				dist[next] = cand
				if cost == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}

	return INF
}

type D = int
type Deque struct{ l, r []D }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func dijkstraPoint(n int, adjList [][][2]int, start, target int) int {
	dist := make([]int, n)
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
		if cur == target {
			return curDist
		}
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge[0], edge[1]
			if (next == start && cur == target) || (next == target && cur == start) {
				continue
			}
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		}
	}

	return INF
}

type H = struct{ node, dist int }

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b H) int

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
