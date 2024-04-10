package main

// https://leetcode.cn/problems/network-delay-time/submissions/
func networkDelayTime(times [][]int, n int, k int) int {
	graph := make([][]Neighbour, n)
	for _, e := range times {
		u, v, w := e[0]-1, e[1]-1, e[2]
		graph[u] = append(graph[u], Neighbour{int32(v), int32(w)})
	}
	dist := DijkstraSiftHeap1(int32(n), graph, int32(k-1))
	res := 0
	for _, d := range dist {
		if d == INF {
			return -1
		}
		if d > res {
			res = d
		}
	}
	return res
}

const INF int = 1e18

// 采用SiftHeap加速的dijkstra算法.求出起点到各点的最短距离.
type Neighbour struct {
	next   int32
	weight int32
}

func DijkstraSiftHeap1(n int32, graph [][]Neighbour, start int32) []int {
	dist := make([]int, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF
	}
	pq := NewSiftHeap32(n, func(i, j int32) bool { return dist[i] < dist[j] })
	dist[start] = 0
	pq.Push(start)
	for pq.Size() > 0 {
		cur := pq.Pop()
		for _, e := range graph[cur] {
			next, weight := e.next, e.weight
			cand := dist[cur] + int(weight)
			if cand < dist[next] {
				dist[next] = cand
				pq.Push(next)
			}
		}
	}
	return dist
}

// 采用SiftHeap加速的dijkstra算法.求出一条路径.
//  如果不存在,则返回(INF, nil).
func DijkstraSiftHeap232(n int32, graph [][]Neighbour, start, end int32) (res int, path []int32) {
	dist := make([]int, n)
	pre := make([]int32, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF
		pre[i] = -1
	}
	pq := NewSiftHeap32(n, func(i, j int32) bool { return dist[i] < dist[j] })
	dist[start] = 0
	pq.Push(start)
	for pq.Size() > 0 {
		cur := pq.Pop()
		for _, e := range graph[cur] {
			next, weight := e.next, e.weight
			cand := dist[cur] + int(weight)
			if cand < dist[next] {
				dist[next] = cand
				pq.Push(next)
				pre[next] = cur
			}
		}
	}

	if dist[end] == INF {
		return INF, nil
	}

	res = dist[end]
	cur := end
	for pre[cur] != -1 {
		path = append(path, cur)
		cur = pre[cur]
	}
	path = append(path, start)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return
}

// 多源最短路, 返回(距离, 前驱, 根节点).
// 用于求出离每个点最近的起点.
func DijkstraMultiStart32(n int32, graph [][]Neighbour, starts []int32) (dist []int, pre []int32, roots []int32) {
	dist = make([]int, n)
	pre = make([]int32, n)
	roots = make([]int32, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF
		pre[i] = -1
		roots[i] = -1
	}
	pq := NewSiftHeap32(n, func(i, j int32) bool { return dist[i] < dist[j] })
	for _, v := range starts {
		dist[v] = 0
		roots[v] = v
		pq.Push(v)
	}
	for pq.Size() > 0 {
		cur := pq.Pop()
		for _, e := range graph[cur] {
			next, weight := e.next, e.weight
			cand := dist[cur] + int(weight)
			if cand < dist[next] {
				dist[next] = cand
				roots[next] = roots[cur]
				pre[next] = cur
				pq.Push(next)
			}
		}
	}
	return
}

type SiftHeap32 struct {
	heap []int32
	pos  []int32
	less func(i, j int32) bool
	ptr  int32
}

func NewSiftHeap32(n int32, less func(i, j int32) bool) *SiftHeap32 {
	pos := make([]int32, n)
	for i := int32(0); i < n; i++ {
		pos[i] = -1
	}
	return &SiftHeap32{
		heap: make([]int32, n),
		pos:  pos,
		less: less,
	}
}

func (h *SiftHeap32) Push(i int32) {
	if h.pos[i] == -1 {
		h.pos[i] = h.ptr
		h.heap[h.ptr] = i
		h.ptr++
	}
	h._siftUp(i)
}

// 如果不存在,则返回-1.
func (h *SiftHeap32) Pop() int32 {
	if h.ptr == 0 {
		return -1
	}
	res := h.heap[0]
	h.pos[res] = -1
	h.ptr--
	ptr := h.ptr
	if ptr > 0 {
		tmp := h.heap[ptr]
		h.pos[tmp] = 0
		h.heap[0] = tmp
		h._siftDown(tmp)
	}
	return res
}

// 如果不存在,则返回-1.
func (h *SiftHeap32) Peek() int32 {
	if h.ptr == 0 {
		return -1
	}
	return h.heap[0]
}

func (h *SiftHeap32) Size() int32 {
	return h.ptr
}

func (h *SiftHeap32) _siftUp(i int32) {
	curPos := h.pos[i]
	p := int32(0)
	for curPos != 0 {
		p = h.heap[(curPos-1)>>1]
		if !h.less(i, p) {
			break
		}
		h.pos[p] = curPos
		h.heap[curPos] = p
		curPos = (curPos - 1) >> 1
	}
	h.pos[i] = curPos
	h.heap[curPos] = i
}

func (h *SiftHeap32) _siftDown(i int32) {
	curPos := h.pos[i]
	c := int32(0)
	for {
		c = (curPos << 1) | 1
		if c >= h.ptr {
			break
		}
		if c+1 < h.ptr && h.less(h.heap[c+1], h.heap[c]) {
			c++
		}
		if !h.less(h.heap[c], i) {
			break
		}
		tmp := h.heap[c]
		h.heap[curPos] = tmp
		h.pos[tmp] = curPos
		curPos = c
	}
	h.pos[i] = curPos
	h.heap[curPos] = i
}
