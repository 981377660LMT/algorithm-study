package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

// 给定一个n个点m条边的无向带权图.其中0~k-1为休息点(基地),k~n-1为普通点.
// q 次询问，每次给定(u, v, t)，如果连续的不在房子中的时间超过t，则无法到达.
// 求解每次询问是否能到达.
// n,m,q<=2e5
// u,v<k
//
// 解法:
// 1. 多源dijkstra算出点离他最近的基地 形成若干个组(每个点属于哪个基地)
// 那么从一个组走到另一个组后 最好先马上去休息 (最好先去离这个点最近的基地)
// 问题就转化为基地之间走了
// 2. 边排序；离线查询+并查集 判断两个组是否连通
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}
	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1], &queries[i][2])
		queries[i][0]--
		queries[i][1]--
	}

	adjList := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], [2]int{v, w})
		adjList[v] = append(adjList[v], [2]int{u, w})
	}
	starts := make([]int, k) // 起点为所有休息点
	for i := 0; i < k; i++ {
		starts[i] = i
	}
	dist, _, root := DijkstraMultiStart(n, adjList, starts)

	toMerge := [][3]int{}
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		if root[u] != root[v] {
			toMerge = append(toMerge, [3]int{dist[u] + dist[v] + w, root[u], root[v]})
		}
	}
	sort.Slice(toMerge, func(i, j int) bool { return toMerge[i][0] < toMerge[j][0] })

	order := make([]int, q) // 按照查询的时间排序
	for i := 0; i < q; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return queries[order[i]][2] < queries[order[j]][2] })
	uf := NewUnionFindArraySimple(k)

	ptr := 0
	for i := 0; i < q; i++ {
		qi := order[i]
		u, v, t := queries[qi][0], queries[qi][1], queries[qi][2]
		for ptr < len(toMerge) {
			weight, root1, root2 := toMerge[ptr][0], toMerge[ptr][1], toMerge[ptr][2]
			if weight <= t {
				ptr++
				uf.Union(root1, root2)
			} else {
				break
			}
		}

		if uf.Find(u) == uf.Find(v) {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}

}

type UnionFindArraySimple struct {
	Part int
	n    int
	data []int32
}

func NewUnionFindArraySimple(n int) *UnionFindArraySimple {
	data := make([]int32, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple) Union(key1 int, key2 int) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *UnionFindArraySimple) Find(key int) int {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = int32(u.Find(int(u.data[key])))
	return int(u.data[key])
}

func (u *UnionFindArraySimple) GetSize(key int) int {
	return int(-u.data[u.Find(key)])
}

// 多源最短路, 返回(距离, 前驱, 根节点).
// 用于求出离每个点最近的起点.
func DijkstraMultiStart(n int, graph [][][2]int, roots []int) (dist []int, pre []int, belongRoot []int) {
	dist = make([]int, n)
	pre = make([]int, n)
	belongRoot = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
		belongRoot[i] = -1
	}
	pq := NewSiftHeap(n, func(i, j int32) bool { return dist[i] < dist[j] })
	for _, v := range roots {
		dist[v] = 0
		belongRoot[v] = v
		pq.Push(v)
	}
	for pq.Size() > 0 {
		cur := pq.Pop()
		for _, e := range graph[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				belongRoot[next] = belongRoot[cur]
				pre[next] = cur
				pq.Push(next)
			}
		}
	}
	return
}

type SiftHeap struct {
	heap []int32
	pos  []int32
	less func(i, j int32) bool
	ptr  int32
}

func NewSiftHeap(n int, less func(i, j int32) bool) *SiftHeap {
	pos := make([]int32, n)
	for i := 0; i < n; i++ {
		pos[i] = -1
	}
	return &SiftHeap{
		heap: make([]int32, n),
		pos:  pos,
		less: less,
	}
}

func (h *SiftHeap) Push(i int) {
	if h.pos[i] == -1 {
		h.pos[i] = h.ptr
		h.heap[h.ptr] = int32(i)
		h.ptr++
	}
	h._siftUp(int32(i))
}

// 如果不存在,则返回-1.
func (h *SiftHeap) Pop() int {
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
	return int(res)
}

// 如果不存在,则返回-1.
func (h *SiftHeap) Peek() int {
	if h.ptr == 0 {
		return -1
	}
	return int(h.heap[0])
}

func (h *SiftHeap) Size() int {
	return int(h.ptr)
}

func (h *SiftHeap) _siftUp(i int32) {
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

func (h *SiftHeap) _siftDown(i int32) {
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
