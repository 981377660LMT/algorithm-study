// // RangeToRangeGraph (区间图)
// // !原图的连通分量/最短路在新图上仍然等价
// // 线段树优化建图

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// )

// const INF WeightType = 1e18

// func main() {
// 	CF786B()
// 	// yuki1868()
// }

// // https://www.luogu.com.cn/problem/CF786B
// func CF786B() {
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, q, start int32
// 	fmt.Fscan(in, &n, &q, &start)
// 	start--
// 	G := NewRangeToRangeGraph(n)
// 	for i := int32(0); i < q; i++ {
// 		var op int32
// 		fmt.Fscan(in, &op)
// 		if op == 1 {
// 			var from, to int32
// 			var weight WeightType
// 			fmt.Fscan(in, &from, &to, &weight)
// 			from--
// 			to--
// 			G.Add(from, to, weight)
// 		} else if op == 2 {
// 			var from, l, r int32
// 			var weight WeightType
// 			fmt.Fscan(in, &from, &l, &r, &weight)
// 			from--
// 			l--
// 			G.AddToRange(from, l, r, weight)
// 		} else if op == 3 {
// 			var to, l, r int32
// 			var weight WeightType
// 			fmt.Fscan(in, &to, &l, &r, &weight)
// 			to--
// 			l--
// 			G.AddFromRange(l, r, to, weight)
// 		}
// 	}

// 	newGraph, size := G.Build()
// 	res := DijkstraInt32(size, newGraph, start)
// 	for i := int32(0); i < n; i++ {
// 		fmt.Fprint(out, res[i], " ")
// 	}
// }

// func yuki1868() {
// 	// https://yukicoder.me/problems/no/1868
// 	// !给定一张有向图,每个点i可以向右达到i+1,i+2,...,targets[i]。求从0到n-1的最短路。
// 	// 解法1：每个点i连接targets[i],边权为1,所有i到i-1连边,边权为0。然后跑最短路。(前后缀优化建图)
// 	// 解法2：RangeToRangeGraph。每个点i连接i+1,i+2,...,targets[i]。然后跑最短路。
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n int32
// 	fmt.Fscan(in, &n)
// 	targets := make([]int32, n-1) // !从i可以到 i+1, i+2, ..., targets[i]
// 	for i := range targets {
// 		fmt.Fscan(in, &targets[i])
// 		targets[i]-- // [0,n-1]内
// 	}

// 	R := NewRangeToRangeGraph(n)
// 	for i := int32(0); i < n-1; i++ {
// 		R.AddToRange(i, i+1, targets[i]+1, 1) // 左闭右开
// 	}
// 	adjList, newN := R.Build()

// 	dist, queue := make([]WeightType, newN), NewDeque(newN)
// 	for i := range dist {
// 		dist[i] = INF
// 	}
// 	dist[0] = 0
// 	queue.Append(0)
// 	for queue.Size() > 0 {
// 		cur := queue.PopLeft()
// 		nexts := adjList[cur]
// 		for i := 0; i < len(nexts); i++ {
// 			e := &nexts[i]
// 			next, weight := e.to, e.weight
// 			cand := dist[cur] + weight
// 			if cand < dist[next] {
// 				dist[next] = cand
// 				if weight == 0 {
// 					queue.AppendLeft(next)
// 				} else {
// 					queue.Append(next)
// 				}
// 			}
// 		}
// 	}

// 	fmt.Fprintln(out, dist[n-1])
// }

// func jump(nums []int) int {
// 	// 45. 跳跃游戏 II
// 	// https://leetcode.cn/problems/jump-game-ii/
// 	n := len(nums)
// 	G := NewRangeToRangeGraph(int32(n))
// 	for i := 0; i < n; i++ {
// 		G.AddToRange(int32(i), int32(i+1), min32(int32(i+nums[i]+1), int32(n)), 1)
// 	}
// 	adjList, _ := G.Build()
// 	bfs := func(start int32, adjList [][]neighbor) []WeightType {
// 		n := len(adjList)
// 		dist := make([]WeightType, n)
// 		for i := 0; i < n; i++ {
// 			dist[i] = INF
// 		}
// 		dist[start] = 0
// 		queue := []int32{start}
// 		for len(queue) > 0 {
// 			cur := queue[0]
// 			queue = queue[1:]
// 			nexts := adjList[cur]
// 			for i := 0; i < len(nexts); i++ {
// 				e := &nexts[i]
// 				next, weight := e.to, e.weight
// 				cand := dist[cur] + weight
// 				if cand < dist[next] {
// 					dist[next] = cand
// 					queue = append(queue, next)
// 				}
// 			}
// 		}

// 		return dist
// 	}
// 	dist := bfs(0, adjList)
// 	return int(dist[n-1])
// }

// type WeightType = int
// type edge struct {
// 	from, to int32
// 	weight   WeightType
// }
// type neighbor struct {
// 	to     int32
// 	weight WeightType
// }

// type RangeToRangeGraph struct {
// 	n     int32
// 	nNode int32
// 	edges []edge
// }

// func NewRangeToRangeGraph(n int32) *RangeToRangeGraph {
// 	g := &RangeToRangeGraph{
// 		n:     n,
// 		nNode: n * 3,
// 	}
// 	for i := int32(2); i < n+n; i++ {
// 		g.edges = append(g.edges, edge{from: g.toUpperIdx(i / 2), to: g.toUpperIdx(i), weight: 0})
// 	}
// 	for i := int32(2); i < n+n; i++ {
// 		g.edges = append(g.edges, edge{from: g.toLowerIdx(i), to: g.toLowerIdx(i / 2), weight: 0})
// 	}
// 	return g
// }

// // 添加有向边 from -> to, 权重为 weight.
// func (g *RangeToRangeGraph) Add(from, to int32, weight WeightType) {
// 	g.edges = append(g.edges, edge{from: from, to: to, weight: weight})
// }

// // 从区间 [fromStart, fromEnd) 中的每个点到 to 都添加一条有向边，权重为 weight.
// func (g *RangeToRangeGraph) AddFromRange(fromStart, fromEnd, to int32, weight WeightType) {
// 	l, r := int32(fromStart)+g.n, int32(fromEnd)+g.n
// 	for l < r {
// 		if l&1 == 1 {
// 			g.Add(g.toLowerIdx(l), to, weight)
// 			l++
// 		}
// 		if r&1 == 1 {
// 			r--
// 			g.Add(g.toLowerIdx(r), to, weight)
// 		}
// 		l >>= 1
// 		r >>= 1
// 	}
// }

// // 从 from 到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
// func (g *RangeToRangeGraph) AddToRange(from, toStart, toEnd int32, weight WeightType) {
// 	l, r := int32(toStart)+g.n, int32(toEnd)+g.n
// 	for l < r {
// 		if l&1 == 1 {
// 			g.Add(from, g.toUpperIdx(l), weight)
// 			l++
// 		}
// 		if r&1 == 1 {
// 			r--
// 			g.Add(from, g.toUpperIdx(r), weight)
// 		}
// 		l >>= 1
// 		r >>= 1
// 	}
// }

// // 从区间 [fromStart, fromEnd) 中的每个点到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
// func (g *RangeToRangeGraph) AddRangeToRange(fromStart, fromEnd, toStart, toEnd int32, weight WeightType) {
// 	newNode := g.nNode
// 	g.nNode++
// 	g.AddFromRange(fromStart, fromEnd, newNode, weight)
// 	g.AddToRange(newNode, toStart, toEnd, 0)
// }

// // 返回`新图的有向邻接表和新图的节点数`.
// func (g *RangeToRangeGraph) Build() (graph [][]neighbor, vertex int32) {
// 	graph = make([][]neighbor, g.nNode)
// 	for i := 0; i < len(g.edges); i++ {
// 		e := &g.edges[i]
// 		u, v, w := e.from, e.to, e.weight
// 		graph[u] = append(graph[u], neighbor{v, w})
// 	}
// 	return graph, g.nNode
// }

// func (g *RangeToRangeGraph) toUpperIdx(i int32) int32 {
// 	if i >= g.n {
// 		return i - g.n
// 	}
// 	return g.n + i
// }

// func (g *RangeToRangeGraph) toLowerIdx(i int32) int32 {
// 	if i >= g.n {
// 		return i - g.n
// 	}
// 	return g.n + g.n + i
// }

// type D = int32
// type Deque struct{ l, r []D }

// func NewDeque(cap int32) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

// func (q Deque) Empty() bool {
// 	return len(q.l) == 0 && len(q.r) == 0
// }

// func (q Deque) Size() int {
// 	return len(q.l) + len(q.r)
// }

// func (q *Deque) AppendLeft(v D) {
// 	q.l = append(q.l, v)
// }

// func (q *Deque) Append(v D) {
// 	q.r = append(q.r, v)
// }

// func (q *Deque) PopLeft() (v D) {
// 	if len(q.l) > 0 {
// 		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
// 	} else {
// 		v, q.r = q.r[0], q.r[1:]
// 	}
// 	return
// }

// func (q *Deque) Pop() (v D) {
// 	if len(q.r) > 0 {
// 		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
// 	} else {
// 		v, q.l = q.l[0], q.l[1:]
// 	}
// 	return
// }

// func (q Deque) Front() D {
// 	if len(q.l) > 0 {
// 		return q.l[len(q.l)-1]
// 	}
// 	return q.r[0]
// }

// func (q Deque) Back() D {
// 	if len(q.r) > 0 {
// 		return q.r[len(q.r)-1]
// 	}
// 	return q.l[0]
// }

// // 0 <= i < q.Size()
// func (q Deque) At(i int) D {
// 	if i < len(q.l) {
// 		return q.l[len(q.l)-1-i]
// 	}
// 	return q.r[i-len(q.l)]
// }

// // 如果不存在则返回 -1.
// func DijkstraInt32(n int32, graph [][]neighbor, start int32) []WeightType {
// 	pq := NewHeap(func(a, b H) bool { return a.dist < b.dist }, []H{{0, start}})
// 	dist := make([]WeightType, n)
// 	for i := range dist {
// 		dist[i] = INF
// 	}
// 	dist[start] = 0
// 	for pq.Len() > 0 {
// 		cur := pq.Pop()
// 		curDist, curNode := cur.dist, cur.node
// 		if curDist > dist[curNode] {
// 			continue
// 		}
// 		nexts := graph[curNode]
// 		for i := 0; i < len(nexts); i++ {
// 			e := &nexts[i]
// 			next, weight := e.to, e.weight
// 			if tmp := curDist + weight; tmp < dist[next] {
// 				dist[next] = tmp
// 				pq.Push(H{tmp, next})
// 			}

// 		}
// 	}
// 	for i := range dist {
// 		if dist[i] == INF {
// 			dist[i] = -1
// 		}
// 	}
// 	return dist
// }

// type H = struct {
// 	dist WeightType
// 	node int32
// }

// func NewHeap(less func(a, b H) bool, nums []H) *Heap {
// 	nums = append(nums[:0:0], nums...)
// 	heap := &Heap{less: less, data: nums}
// 	heap.heapify()
// 	return heap
// }

// type Heap struct {
// 	data []H
// 	less func(a, b H) bool
// }

// func (h *Heap) Push(value H) {
// 	h.data = append(h.data, value)
// 	h.pushUp(h.Len() - 1)
// }

// func (h *Heap) Pop() (value H) {
// 	if h.Len() == 0 {
// 		panic("heap is empty")
// 	}
// 	value = h.data[0]
// 	h.data[0] = h.data[h.Len()-1]
// 	h.data = h.data[:h.Len()-1]
// 	h.pushDown(0)
// 	return
// }

// func (h *Heap) Top() (value H) {
// 	value = h.data[0]
// 	return
// }

// func (h *Heap) Len() int { return len(h.data) }

// func (h *Heap) heapify() {
// 	n := h.Len()
// 	for i := (n >> 1) - 1; i > -1; i-- {
// 		h.pushDown(i)
// 	}
// }

// func (h *Heap) pushUp(root int) {
// 	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
// 		h.data[root], h.data[parent] = h.data[parent], h.data[root]
// 		root = parent
// 	}
// }

// func (h *Heap) pushDown(root int) {
// 	n := h.Len()
// 	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
// 		right := left + 1
// 		minIndex := root
// 		if h.less(h.data[left], h.data[minIndex]) {
// 			minIndex = left
// 		}
// 		if right < n && h.less(h.data[right], h.data[minIndex]) {
// 			minIndex = right
// 		}
// 		if minIndex == root {
// 			return
// 		}
// 		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
// 		root = minIndex
// 	}
// }
// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func min32(a, b int32) int32 {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

//	func max32(a, b int32) int32 {
//		if a > b {
//			return a
//		}
//		return b
//	}
package main

func main() {

}
