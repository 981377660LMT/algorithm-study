// https://nyaannyaan.github.io/library/shortest-path/dijkstra-skew-heap.hpp
// !node代替的数组实现的堆，在js/py这类语言中会慢很多，在golang中会稍快一些

package main

// 使用SkewHeap优化Dijkstra算法常数.
//   如果不存在,最短路长度为-1.
func DijkstraSkewHeap(n int, adjList [][][2]int, start int) []int {
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	q := NewSkewHeap(n)
	dist[start] = 0
	q.v[start].key = 0
	q.rt = start
	for !q.Empty() {
		dc := q.v[q.rt].key
		cur := q.rt
		q.Pop()
		for _, e := range adjList[cur] {
			next, cost := e[0], e[1]
			if dist[next] == -1 {
				dist[next] = dc + cost
				q.v[next].key = dc + cost
				q.rt = q.Meld(q.rt, next)
			} else if dc+cost < dist[next] {
				dist[next] = dc + cost
				q.Update(next, dc+cost)
			}
		}
	}
	return dist
}

// 使用SkewHeap优化Dijkstra算法常数.
//  求出从start到end的最短路.如果不存在,最短路长度为-1.
func DijkstraSkewHeap2(n int, adjList [][][2]int, start, end int) (int, []int) {
	dist := make([]int, n)
	pre := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
		pre[i] = -1
	}
	q := NewSkewHeap(n)
	dist[start] = 0
	q.v[start].key = 0
	q.rt = start

	for !q.Empty() {
		dc := q.v[q.rt].key
		cur := q.rt
		q.Pop()
		for _, e := range adjList[cur] {
			next, cost := e[0], e[1]
			if dist[next] == -1 {
				dist[next] = dc + cost
				q.v[next].key = dc + cost
				q.rt = q.Meld(q.rt, next)
				pre[next] = cur
			} else if dc+cost < dist[next] {
				dist[next] = dc + cost
				q.Update(next, dc+cost)
				pre[next] = cur
			}
		}
	}

	if dist[end] == -1 {
		return -1, nil
	}

	path := []int{}
	cur := end
	for pre[cur] != -1 {
		path = append(path, cur)
		cur = pre[cur]
	}
	path = append(path, start)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return dist[end], path
}

type P = int

type SkewHeap struct {
	v  []*SNode
	rt int
}

func NewSkewHeap(n int) *SkewHeap {
	nodes := make([]*SNode, n)
	for i := 0; i < n; i++ {
		nodes[i] = NewSNode()
	}
	return &SkewHeap{v: nodes, rt: -1}
}

func (h *SkewHeap) Meld(x, y int) int {
	if x == -1 || y == -1 {
		if x == -1 {
			return y
		}
		return x
	}
	if h.v[x].key > h.v[y].key {
		x ^= y
		y ^= x
		x ^= y
	}
	h.v[x].r = h.Meld(h.v[x].r, y)
	h.v[h.v[x].r].p = x
	h.v[x].l, h.v[x].r = h.v[x].r, h.v[x].l
	return x
}

func (h *SkewHeap) Pop() {
	h.rt = h.Meld(h.v[h.rt].l, h.v[h.rt].r)
}

func (h *SkewHeap) Update(x int, k P) {
	n := h.v[x]
	n.key = k
	if x == h.rt {
		return
	}
	p := h.v[n.p]
	if p.key <= k {
		return
	}
	if p.l == x {
		p.l = -1
	} else {
		p.r = -1
	}
	n.p = -1
	h.rt = h.Meld(h.rt, x)
}

func (h *SkewHeap) Empty() bool {
	return h.rt == -1
}

type SNode struct {
	key     P
	p, l, r int
}

func NewSNode() *SNode {
	return &SNode{p: -1, l: -1, r: -1}
}
