// Dominator Tree of a Directed Graph
// O(ElogV)

// python版本:
// https://tjkendev.github.io/procon-library/python/graph/dominator-lengauer-tarjan.html

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 求出有向图的支配树, 返回每个节点的`直接支配顶点idom`.
//
//	n: 顶点数
//	edges: 有向边
//	root: 支配树的根节点
//
// 考虑支配树中从root到达结点i的路径, 该路径上的结点都是i的支配者，要到达i就必须经过这些结点(有向图的割点).
// 其中, 从root到达结点i的路径上的最后一个结点就是i的直接支配者idom[i]，也是距离i`最近的支配者`.
// !idom[i]=-1：从根节点无法到达结点i
// !idom[i]=root：从根节点到达结点i的路径上不存在i的支配者(割点)
// !idom[i]为其他值：i的最近的支配者(割点)为idom[i]
func BuildDominatorTree(n int32, edges [][]int32, root int32) []int32 {
	graph := make([][]int32, n)
	rGraph := make([][]int32, n)
	for _, e := range edges { // addEdge
		u, v := e[0], e[1]
		graph[u] = append(graph[u], v)
		rGraph[v] = append(rGraph[v], u)
	}

	// build
	par := make([]int32, n)
	idom, semi := make([]int32, n), make([]int32, n) // idom: immediate dominator
	for i := int32(0); i < n; i++ {
		idom[i], semi[i] = -1, -1
	}
	ord := make([]int32, 0, n)

	var dfs func(v int32)
	dfs = func(v int32) {

		semi[v] = int32(len(ord))
		ord = append(ord, v)
		for _, u := range graph[v] {
			if semi[u] != -1 {
				continue
			}
			par[u] = v
			dfs(u)
		}
	}
	dfs(root)

	bkt := make([][]int32, n)
	uf := newUnionFind(semi)
	us := make([]int32, n)

	for i := len(ord) - 1; i >= 0; i-- {
		v := ord[i]
		for _, u := range rGraph[v] {
			if semi[u] < 0 {
				continue
			}
			u := uf.eval(u)
			if semi[u] < semi[v] {
				semi[v] = semi[u]
			}
		}
		bkt[ord[semi[v]]] = append(bkt[ord[semi[v]]], v)
		for _, u := range bkt[par[v]] {
			us[u] = uf.eval(u)
		}
		bkt[par[v]] = bkt[par[v]][:0] // clear
		uf.link(par[v], v)
	}

	for i := 1; i < len(ord); i++ {
		v := ord[i]
		u := us[v]
		if semi[v] == semi[u] {
			idom[v] = semi[v]
		} else {
			idom[v] = idom[u]
		}
	}

	for i := 1; i < len(ord); i++ {
		v := ord[i]
		idom[v] = ord[idom[v]]
	}

	idom[root] = root

	return idom
}

type unionFind struct {
	semi   []int32
	ps, ms []int32
}

func newUnionFind(semi []int32) *unionFind {
	n := int32(len(semi))
	ps, ms := make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		ps[i], ms[i] = i, i
	}
	return &unionFind{semi, ps, ms}
}

func (uf *unionFind) find(v int32) int32 {
	if uf.ps[v] == v {
		return v
	}
	r := uf.find(uf.ps[v])
	if uf.semi[uf.ms[v]] > uf.semi[uf.ms[uf.ps[v]]] {
		uf.ms[v] = uf.ms[uf.ps[v]]
	}
	uf.ps[v] = r
	return r
}

func (uf *unionFind) eval(v int32) int32 {
	uf.find(v)
	return uf.ms[v]
}

func (uf *unionFind) link(p, v int32) {
	uf.ps[v] = p
}

// 捕获赤牛大作战
// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0294
// 给定一个带有n个结点的`连通的有向图`,基地位于0号点
// 给定q个询问
// 每个询问给定一个赤牛的位置qi,你需要选择一个点作为埋伏点来捕获赤牛
// 选择的点的要求为:
// 1. 不能在基地埋伏(太危险)
// 2. 埋伏点必须是从基地到达赤牛位置的所有路径中一定会经过的点
// 3. 在满足1,2的条件下,埋伏点距离赤牛的距离越近越好;如果不存在这样的点,就直接埋伏在赤牛的位置
// 对每个询问,输出一个整数,表示你选择的埋伏点的编号
// n<=1e5, m<=3e5, q<=1e5
func catchAkabeko(n int32, edges [][]int32, queries []int32) []int32 {
	idom := BuildDominatorTree(n, edges, 0) // 每个点最近的支配者
	res := make([]int32, 0, len(queries))
	for _, q := range queries {
		if idom[q] == 0 || idom[q] == -1 { // 不存在割点或者无法到达,直接埋伏在赤牛的位置
			res = append(res, q)
		} else {
			res = append(res, idom[q]) // 埋伏在最近的支配者
		}
	}
	return res
}

// https://leetcode.cn/problems/disconnect-path-in-a-binary-matrix-by-at-most-one-flip/
func isPossibleToCutPath(grid [][]int) bool {
	ROW, COL := int32(len(grid)), int32(len(grid[0]))
	edges := [][]int32{}
	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			if grid[i][j] == 0 {
				continue
			}
			if i+1 < ROW && grid[i+1][j] == 1 {
				edges = append(edges, []int32{i*COL + j, (i+1)*COL + j})
			}
			if j+1 < COL && grid[i][j+1] == 1 {
				edges = append(edges, []int32{i*COL + j, i*COL + j + 1})
			}
		}
	}

	idom := BuildDominatorTree(ROW*COL, edges, 0)
	return idom[ROW*COL-1] != 0
}

// 求每个点能支配的个数(包括自己)
// https://www.luogu.com.cn/problem/P5180
func DominatorTreeSize(n int32, edges [][]int32, root int32) []int32 {
	parents := BuildDominatorTree(n, edges, root)
	adjList := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		p := parents[i]
		if p != -1 && p != i {
			adjList[p] = append(adjList[p], i)
		}
	}

	subSize := make([]int32, n)
	var dfs func(cur int32)
	dfs = func(cur int32) {
		subSize[cur] = 1
		for _, next := range adjList[cur] {
			dfs(next)
			subSize[cur] += subSize[next]
		}
	}
	dfs(root)

	return subSize
}

// Team Rocket Rises Again (最短路图+支配树)
// https://www.luogu.com.cn/problem/CF757F
// 给定 n 个点， m 条边的带权无向图，一个起点 s 。
// 在删掉一个点的情况下，使尽可能多的点到起点的最短距离被改变，求最多的点数。
// n<=1e5, m<=3e5
//
// !1.转为最短路图
//
//	最短路图的简单性质：
//	- 从 s 到任意一点 u 的任意一条最短路的每一条边一定存在在新图中。除此之外的其他边一定不存在。
//	- 最短路图是一个 DAG。
//
// !2.DAG 上的支配树(DAG 上任意一点 u 删除后，起点 s 无法到达的点的个数)
// !删掉点 u 后到 s 的最短距离改变的点的数量，等价于在最短路图上删掉点 u 后 s 无法到达的数量。
func CF757F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, start int32
	fmt.Fscan(in, &n, &m, &start)
	start--

	edges := make([][]int32, 0, m) // [u, v, w]
	for i := int32(0); i < m; i++ {
		var u, v, w int32
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int32{u, v, w})
	}

	adjList := make([][]Neighbor, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], Neighbor{v, w})
		adjList[v] = append(adjList[v], Neighbor{u, w})
	}

	dist := Dijkstra(n, adjList, start)

	newGrapEdges := make([][]int32, 0, m) // 最短路径图的边
	for i := int32(0); i < n; i++ {
		if dist[i] == INF {
			continue
		}
		for _, edge := range adjList[i] {
			to, weight := edge.to, edge.weight
			if dist[to] == dist[i]+int(weight) {
				newGrapEdges = append(newGrapEdges, []int32{i, to})
			}
		}
	}

	parents := BuildDominatorTree(n, newGrapEdges, start)
	dominatorTree := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		p := parents[i]
		if p != -1 && p != i {
			dominatorTree[p] = append(dominatorTree[p], i)
		}
	}

	subSize := make([]int32, n)
	var dfs func(cur int32)
	dfs = func(cur int32) {
		subSize[cur] = 1
		for _, next := range dominatorTree[cur] {
			dfs(next)
			subSize[cur] += subSize[next]
		}
	}
	dfs(start)

	res := int32(0)
	for _, next := range dominatorTree[start] {
		if subSize[next] > res {
			res = subSize[next]
		}
	}
	fmt.Fprintln(out, res)
}

const INF int = 1e18

type Neighbor struct{ to, weight int32 }

func Dijkstra(n int32, adjList [][]Neighbor, start int32) (dist []int) {
	dist = make([]int, n)
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
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge.to, edge.weight
			if cand := curDist + int(weight); cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		}
	}

	return
}

type H = struct {
	node int32
	dist int
}

// Should return a number:
//
//	negative , if a < b
//	zero		 , if a == b
//	positive , if a > b
type Comparator = func(a, b H) int

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

func main() {
	CF757F()
}
