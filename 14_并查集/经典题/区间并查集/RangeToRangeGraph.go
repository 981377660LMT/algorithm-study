// RangeToRangeGraph (区间图)
// !原图的连通分量/最短路在新图上仍然等价
// 线段树优化建图

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	CF786B()
	// yuki1868()
}

// https://www.luogu.com.cn/problem/CF786B
func CF786B() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, start int32
	fmt.Fscan(in, &n, &q, &start)
	start--
	G := NewRangeToRangeGraph(n, 0)
	newGraph := make([][]neighbor, G.Size())
	G.Init(func(from, to int32) { newGraph[from] = append(newGraph[from], neighbor{to, 0}) })
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var from, to int32
			var weight int32
			fmt.Fscan(in, &from, &to, &weight)
			from--
			to--
			G.Add(from, to, func(from, to int32) {
				newGraph[from] = append(newGraph[from], neighbor{to, weight})
			})
		} else if op == 2 {
			var from, l, r int32
			var weight int32
			fmt.Fscan(in, &from, &l, &r, &weight)
			from--
			l--
			G.AddToRange(from, l, r, func(from, to int32) {
				newGraph[from] = append(newGraph[from], neighbor{to, weight})
			})
		} else if op == 3 {
			var to, l, r int32
			var weight int32
			fmt.Fscan(in, &to, &l, &r, &weight)
			to--
			l--
			G.AddFromRange(l, r, to, func(from, to int32) {
				newGraph[from] = append(newGraph[from], neighbor{to, weight})
			})
		}
	}

	res := DijkstraSiftHeap1(int32(len(newGraph)), newGraph, start)
	for i := int32(0); i < n; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

func yuki1868() {
	// https://yukicoder.me/problems/no/1868
	// !给定一张有向图,每个点i可以向右达到i+1,i+2,...,targets[i]。求从0到n-1的最短路。
	// 解法1：每个点i连接targets[i],边权为1,所有i到i-1连边,边权为0。然后跑最短路。(前后缀优化建图)
	// 解法2：RangeToRangeGraph。每个点i连接i+1,i+2,...,targets[i]。然后跑最短路。
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	targets := make([]int32, n-1) // !从i可以到 i+1, i+2, ..., targets[i]
	for i := range targets {
		fmt.Fscan(in, &targets[i])
		targets[i]-- // [0,n-1]内
	}

	R := NewRangeToRangeGraph(n, 0)
	adjList := make([][]neighbor, R.Size())
	R.Init(func(from, to int32) {
		adjList[from] = append(adjList[from], neighbor{to, 0})
	})
	for i := int32(0); i < n-1; i++ {
		R.AddToRange(i, i+1, targets[i]+1, func(from, to int32) {
			adjList[from] = append(adjList[from], neighbor{to, 1})
		})
	}

	dist, queue := make([]int, int32(len(adjList))), NewDeque(int32(len(adjList)))
	for i := range dist {
		dist[i] = INF
	}
	dist[0] = 0
	queue.Append(0)
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		nexts := adjList[cur]
		for i := 0; i < len(nexts); i++ {
			e := &nexts[i]
			next, weight := e.next, e.weight
			cand := dist[cur] + int(weight)
			if cand < dist[next] {
				dist[next] = cand
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}

	fmt.Fprintln(out, dist[n-1])
}

func jump(nums []int) int {
	// 45. 跳跃游戏 II
	// https://leetcode.cn/problems/jump-game-ii/
	n := int32(len(nums))
	G := NewRangeToRangeGraph(int32(n), 0)
	adjList := make([][]neighbor, G.Size())
	G.Init(func(from, to int32) { adjList[from] = append(adjList[from], neighbor{to, 0}) })
	for i := int32(0); i < n; i++ {
		G.AddToRange(i, i+1, min32(i+1+int32(nums[i]), n), func(from, to int32) {
			adjList[from] = append(adjList[from], neighbor{to, 1})
		})
	}

	bfs := func(start int32, adjList [][]neighbor) []int32 {
		n := len(adjList)
		dist := make([]int32, n)
		for i := 0; i < n; i++ {
			dist[i] = 1e9
		}
		dist[start] = 0
		queue := []int32{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			nexts := adjList[cur]
			for i := 0; i < len(nexts); i++ {
				e := &nexts[i]
				next, weight := e.next, e.weight
				cand := dist[cur] + int32(weight)
				if cand < dist[next] {
					dist[next] = cand
					queue = append(queue, next)
				}
			}
		}
		return dist
	}
	dist := bfs(0, adjList)
	return int(dist[n-1])
}

type RangeToRangeGraph struct {
	n        int32
	maxSize  int32
	allocPtr int32
}

// 新建一个区间图，n 为原图的节点数，rangeToRangeOpCount 为区间到区间的最大操作次数.
// 最后得到的新图的节点数为 n*3 + rangeToRangeOpCount，前n个节点为原图的节点。
func NewRangeToRangeGraph(n int32, rangeToRangeOpCount int32) *RangeToRangeGraph {
	g := &RangeToRangeGraph{
		n:        n,
		maxSize:  n*3 + rangeToRangeOpCount,
		allocPtr: n * 3,
	}
	return g
}

// 新图的结点数.
func (g *RangeToRangeGraph) Size() int32 { return g.maxSize }

func (g *RangeToRangeGraph) Init(f func(from, to int32)) {
	n := g.n
	for i := int32(2); i < n+n; i++ {
		f(g.toUpperIdx(i>>1), g.toUpperIdx(i))
		f(g.toLowerIdx(i), g.toLowerIdx(i>>1))
	}
}

// 添加有向边 from -> to.
func (g *RangeToRangeGraph) Add(from, to int32, f func(from, to int32)) {
	f(from, to)
}

// 从区间 [fromStart, fromEnd) 中的每个点到 to 都添加一条有向边.
func (g *RangeToRangeGraph) AddFromRange(fromStart, fromEnd, to int32, f func(from, to int32)) {
	l, r := fromStart+g.n, fromEnd+g.n
	for l < r {
		if l&1 == 1 {
			f(g.toLowerIdx(l), to)
			l++
		}
		if r&1 == 1 {
			r--
			f(g.toLowerIdx(r), to)
		}
		l >>= 1
		r >>= 1
	}
}

// 从 from 到区间 [toStart, toEnd) 中的每个点都添加一条有向边.
func (g *RangeToRangeGraph) AddToRange(from, toStart, toEnd int32, f func(from, to int32)) {
	l, r := toStart+g.n, toEnd+g.n
	for l < r {
		if l&1 == 1 {
			f(from, g.toUpperIdx(l))
			l++
		}
		if r&1 == 1 {
			r--
			f(from, g.toUpperIdx(r))
		}
		l >>= 1
		r >>= 1
	}
}

// 从区间 [fromStart, fromEnd) 中的每个点到区间 [toStart, toEnd) 中的每个点都添加一条有向边.
func (g *RangeToRangeGraph) AddRangeToRange(fromStart, fromEnd, toStart, toEnd int32, f func(from, to int32)) {
	newNode := g.allocPtr
	g.allocPtr++
	g.AddFromRange(fromStart, fromEnd, newNode, f)
	g.AddToRange(newNode, toStart, toEnd, f)
}

func (g *RangeToRangeGraph) toUpperIdx(i int32) int32 {
	if i >= g.n {
		return i - g.n
	}
	return g.n + i
}

func (g *RangeToRangeGraph) toLowerIdx(i int32) int32 {
	if i >= g.n {
		return i - g.n
	}
	return g.n + g.n + i
}

type D = int32
type Deque struct{ l, r []D }

func NewDeque(cap int32) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

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

// 采用SiftHeap加速的dijkstra算法.求出起点到各点的最短距离.
type neighbor struct {
	next   int32
	weight int32
}

// 不存在则返回-1.
func DijkstraSiftHeap1(n int32, graph [][]neighbor, start int32) []int {
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
	for i := int32(0); i < n; i++ {
		if dist[i] == INF {
			dist[i] = -1
		}
	}
	return dist
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
