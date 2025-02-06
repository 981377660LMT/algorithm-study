// DAG最小路径覆盖与最小不相交路径覆盖
// https://ei1333.github.io/library/graph/flow/bipartite-flow.hpp
// https://www.luogu.com.cn/problem/P2764
// https://zhuanlan.zhihu.com/p/125759333
// https://www.cnblogs.com/justPassBy/p/5369930.html

// DAG最小路径覆盖(DAG の最小パス被覆)
// DAG最小路径覆盖可以归结为二分图最大匹配问题

// 覆盖DAG所有点的`路径的集合`叫做DAG的路径覆盖,注意路径不相交
// 路径覆盖的路径数最少的集合叫做DAG的最小路径覆盖

// 做法:
// !原图每个点拆成出点和入点,如果有一条有向边u->v,那么就连一条边u.int->v.out
// 得到一个二分图.
// 跑一遍最大流，便能得到最大合并路径数，再用点数去减即得最小路径覆盖数。
// 从in点到out'点的每一条流，都代表着一次合并。而从源点只给每个点输送1单位流量，又保证了每个点只被经过一次。
// !最小路径覆盖=原图的结点数-新图的最大匹配数。
// n<=150 m<=6000

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	P2764()
}

// https://www.luogu.com.cn/problem/P2764
func P2764() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, m)
	for i := 0; i < m; i++ {
		edges[i] = make([]int, 2)
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	count, paths := MinimumPathCovering1(n, edges)
	for _, path := range paths {
		for _, v := range path {
			fmt.Fprint(out, v+1, " ")
		}
		fmt.Fprintln(out)
	}
	fmt.Fprintln(out, count)
}

const INF int = 1e18

// !DAG最小不相交路径覆盖(最小路径点覆盖)
func MinimumPathCovering1(n int, directedEdges [][]int) (count int, paths [][]int) {
	newEdges := make([][]int, 0, len(directedEdges))
	bf := NewBipartiteFlow(n, n)
	for _, edge := range directedEdges {
		u, v := edge[0], edge[1]
		newEdges = append(newEdges, []int{u, v + n}) // 拆点 , A'in => B'out,从A点到B'点的每一条流，都代表着一次合并
		bf.AddEdge(u, v)
	}

	maxMathing := bf.MaxMatching()
	count = n - len(maxMathing)

	uf := NewUnionFindArray(n)
	ml := bf.MatchL
	for i := 0; i < n; i++ {
		if ml[i] == -1 {
			continue
		}
		uf.Union(i, ml[i])
	}

	groups := make(map[int][]int)
	for i := 0; i < n; i++ {
		groups[uf.Find(i)] = append(groups[uf.Find(i)], i)
	}
	for _, group := range groups {
		paths = append(paths, group)
	}
	return
}

// !DAG最小可相交路径覆盖.
// TODO
// https://atcoder.jp/contests/abc374/tasks/abc374_g
func MinimumPathCovering2(n int, directedEdges [][]int) (count int, paths [][]int) {
	adjMatrix := make([][]bool, n)
}

type BipartiteFlow struct {
	N, M           int
	MatchL, MatchR []int
	timeStamp      int
	g, rg          [][]int
	dist           []int
	used           []int
	alive          []bool
	matched        bool
}

// 指定左侧点数n，右侧点数m，初始化二分图最大流.
func NewBipartiteFlow(n, m int) *BipartiteFlow {
	g, rg := make([][]int, n), make([][]int, m)
	matchL, matchR := make([]int, n), make([]int, m)
	used, alive := make([]int, n), make([]bool, n)
	for i := 0; i < n; i++ {
		matchL[i] = -1
		alive[i] = true
	}
	for i := 0; i < m; i++ {
		matchR[i] = -1
	}

	return &BipartiteFlow{
		N:      n,
		M:      m,
		g:      g,
		rg:     rg,
		MatchL: matchL,
		MatchR: matchR,
		used:   used,
		alive:  alive,
	}
}

// 增加一条边u-v.u属于左侧点集，v属于右侧点集.
//
//	!0<=u<n,0<=v<m.
func (bf *BipartiteFlow) AddEdge(u, v int) {
	bf.g[u] = append(bf.g[u], v)
	bf.rg[v] = append(bf.rg[v], u)
}

// 求最大匹配.
//
//	返回(左侧点,右侧点)的匹配对.
//	!0<=左侧点<n,0<=右侧点<m.
func (bf *BipartiteFlow) MaxMatching() [][2]int {
	bf.matched = true
	for {
		bf.buildAugmentPath()
		bf.timeStamp++
		flow := 0
		for i := 0; i < bf.N; i++ {
			if bf.MatchL[i] == -1 {
				tmp := bf.findMinDistAugmentPath(i)
				if tmp {
					flow++
				}
			}
		}

		if flow == 0 {
			break
		}
	}

	res := [][2]int{}
	for i := 0; i < bf.N; i++ {
		if bf.MatchL[i] >= 0 {
			res = append(res, [2]int{i, bf.MatchL[i]})
		}
	}
	return res
}

// 构建残量图.
//
//	left: [0,n), right: [n,n+m), S: n+m, T: n+m+1
func (bf *BipartiteFlow) BuildRisidualGraph() [][]int {
	if !bf.matched {
		bf.MaxMatching()
	}

	S := bf.N + bf.M
	T := bf.N + bf.M + 1
	ris := make([][]int, bf.N+bf.M+2)
	for i := 0; i < bf.N; i++ {
		if bf.MatchL[i] == -1 {
			ris[S] = append(ris[S], i)
		} else {
			ris[i] = append(ris[i], S)
		}
	}

	for i := 0; i < bf.M; i++ {
		if bf.MatchR[i] == -1 {
			ris[i+bf.N] = append(ris[i+bf.N], T)
		} else {
			ris[T] = append(ris[T], i+bf.N)
		}
	}

	for i := 0; i < bf.N; i++ {
		for _, j := range bf.g[i] {
			if bf.MatchL[i] == j {
				ris[j+bf.N] = append(ris[j+bf.N], i)
			} else {
				ris[i] = append(ris[i], j+bf.N)
			}
		}
	}

	return ris
}

func (bf *BipartiteFlow) findResidualPath() []bool {
	res := bf.BuildRisidualGraph()
	que := []int{}
	visited := make([]bool, bf.N+bf.M+2)
	que = append(que, bf.N+bf.M)
	visited[bf.N+bf.M] = true
	for len(que) > 0 {
		idx := que[0]
		que = que[1:]
		for _, to := range res[idx] {
			if visited[to] {
				continue
			}
			visited[to] = true
			que = append(que, to)
		}
	}
	return visited
}

func (bf *BipartiteFlow) buildAugmentPath() {
	que := []int{}
	bf.dist = make([]int, len(bf.g))
	for i := 0; i < len(bf.g); i++ {
		bf.dist[i] = -1
	}
	for i := 0; i < bf.N; i++ {
		if bf.MatchL[i] == -1 {
			que = append(que, i)
			bf.dist[i] = 0
		}
	}
	for len(que) > 0 {
		a := que[0]
		que = que[1:]
		for _, b := range bf.g[a] {
			c := bf.MatchR[b]
			if c >= 0 && bf.dist[c] == -1 {
				bf.dist[c] = bf.dist[a] + 1
				que = append(que, c)
			}
		}
	}
}

func (bf *BipartiteFlow) findMinDistAugmentPath(a int) bool {
	bf.used[a] = bf.timeStamp
	for _, b := range bf.g[a] {
		c := bf.MatchR[b]
		if c < 0 || (bf.used[c] != bf.timeStamp && bf.dist[c] == bf.dist[a]+1 && bf.findMinDistAugmentPath(c)) {
			bf.MatchR[b] = a
			bf.MatchL[a] = b
			return true
		}
	}
	return false
}

type UnionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

type Edge = [2]int

func Dijkstra(n int, adjList [][]Edge, start int) (dist, preV []int) {
	type pqItem struct{ node, dist int }
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	preV = make([]int, n)
	for i := range preV {
		preV[i] = -1
	}

	pq := nhp(func(a, b H) int {
		return a.(pqItem).dist - b.(pqItem).dist
	}, nil)
	pq.Push(pqItem{start, 0})

	for pq.Len() > 0 {
		curNode := pq.Pop().(pqItem)
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge[0], edge[1]
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				preV[next] = cur
				pq.Push(pqItem{next, cand})
			}
		}
	}

	return
}

type H = interface{}

// Should return a number:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
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

func StronglyConnectedComponent(graph [][]int) (compCount int, belong []int) {
	n32 := int32(len(graph))
	compId := int32(0)
	comp := make([]int32, n32)
	low := make([]int32, n32)
	ord := make([]int32, n32)
	for i := range ord {
		ord[i] = -1
	}
	path := []int32{}
	now := int32(0)

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		ord[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if ord[to] == -1 {
				dfs(int32(to))
				if low[v] > low[to] {
					low[v] = low[to]
				}
			} else if low[v] > ord[to] {
				low[v] = ord[to]
			}
		}
		if low[v] == ord[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				ord[u] = n32
				comp[u] = compId
				if u == v {
					break
				}
			}
			compId++
		}
	}

	for v := int32(0); v < n32; v++ {
		if ord[v] == -1 {
			dfs(v)
		}
	}

	compCount = int(compId)
	belong = make([]int, n32)
	for v := int32(0); v < n32; v++ {
		belong[v] = compCount - 1 - int(comp[v])
	}
	return
}

func SccDag(graph [][]int, compCount int, belong []int) (dag [][]int) {
	unique := func(nums []int32) []int32 {
		set := make(map[int32]struct{})
		for _, v := range nums {
			set[v] = struct{}{}
		}
		res := make([]int32, 0, len(set))
		for k := range set {
			res = append(res, k)
		}
		return res
	}

	edges := make([][]int32, compCount)
	for cur, nexts := range graph {
		curComp := belong[cur]
		for _, next := range nexts {
			nextComp := belong[next]
			if curComp != nextComp {
				edges[curComp] = append(edges[curComp], int32(nextComp))
			}
		}
	}

	dag = make([][]int, compCount)
	for cur := 0; cur < compCount; cur++ {
		edges[cur] = unique(edges[cur])
		for _, next := range edges[cur] {
			dag[cur] = append(dag[cur], int(next))
		}
	}

	return
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
