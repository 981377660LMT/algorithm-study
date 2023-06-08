// Bipartite Graph Edge Coloring(二部グラフの辺彩色)
// https://ei1333.github.io/library/graph/others/bipartite-graph-edge-coloring.hpp
// 二分图边着色,使得每个顶点的所有相邻边的颜色不同
// 二分图的边的彩色数为最大度数D,一般图的边彩色为D或者D+1

// add_edge(a, b):
//  aからbに辺を張る. aは二部グラフの左側, bは右側の頂点を指す.
// build():
//  二部グラフの辺彩色を返す.
//  !同じ色に塗るべき辺の番号が同じ配列に格納される.
//  !辺の番号は add_edge() を呼び出した順に 0-indexed で付与される.
// O(E*Sqrt(V)*log(max(deg)))

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// https://judge.yosupo.jp/problem/bipartite_edge_coloring
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var L, R, m int
	fmt.Fscan(in, &L, &R, &m)
	ecbg := NewBipartiteGraphEdgeColoring()
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		ecbg.AddEdge(a, b)
	}
	res := ecbg.Build()
	fmt.Fprintln(out, len(res))
	color := make([]int, m)
	for i, v := range res {
		for _, e := range v {
			color[e] = i
		}
	}
	for _, c := range color {
		fmt.Fprintln(out, c)
	}
}

func NewBipartiteGraphEdgeColoring() *BipartiteGraphEdgeColoring {
	return &BipartiteGraphEdgeColoring{}
}

// aからbに辺を張る. aは二部グラフの左側, bは右側の頂点を指す.
func (bgec *BipartiteGraphEdgeColoring) AddEdge(left, right int) {
	bgec._A = append(bgec._A, left)
	bgec._B = append(bgec._B, right)
	bgec._L = max(bgec._L, left+1)
	bgec._R = max(bgec._R, right+1)
}

// 同じ色に塗るべき辺の番号が同じ配列に格納される.
//  辺の番号は add_edge() を呼び出した順に 0-indexed で付与される.
func (bgec *BipartiteGraphEdgeColoring) Build() [][]int {
	bgec.g = bgec.buildKRegularGraph()
	ord := make([]int, len(bgec.g._A))
	for i := range ord {
		ord[i] = i
	}
	bgec.rec(ord, bgec.g.k)
	res := [][]int{}
	for i := 0; i < len(bgec.ans); i++ {
		res = append(res, []int{})
		for _, j := range bgec.ans[i] {
			if j < len(bgec._A) {
				res[len(res)-1] = append(res[len(res)-1], j)
			}
		}
	}
	return res
}

func conctract(deg []int, k int) *_unionFindArray {
	que := nhp(func(a, b H) int {
		return a[0] - b[0] // TODO
	}, nil)
	for i, d := range deg {
		que.Push(H{d, i})
	}
	uf := newUnionFindArray(len(deg))
	for que.Len() > 1 {
		p := que.Pop()
		q := que.Pop()
		if p[0]+q[0] > k {
			continue
		}
		p[0] += q[0]
		uf.Union(p[1], q[1])
		que.Push(p)
	}
	return uf
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type regularGraph struct {
	k, n   int
	_A, _B []int
}

type BipartiteGraphEdgeColoring struct {
	_L, _R int
	_A, _B []int
	ans    [][]int
	g      *regularGraph
}

func (bgec *BipartiteGraphEdgeColoring) buildKRegularGraph() *regularGraph {
	deg := [][]int{make([]int, bgec._L), make([]int, bgec._R)}
	for _, a := range bgec._A {
		deg[0][a]++
	}
	for _, b := range bgec._B {
		deg[1][b]++
	}

	k := max(maxs(deg[0]...), maxs(deg[1]...))

	// step1
	uf := []*_unionFindArray{conctract(deg[0], k), conctract(deg[1], k)}
	ptr := []int{0, 0}
	id := [][]int{make([]int, bgec._L), make([]int, bgec._R)}
	for i := 0; i < bgec._L; i++ {
		if uf[0].Find(i) == i {
			id[0][i] = ptr[0]
			ptr[0]++
		}
	}
	for i := 0; i < bgec._R; i++ {
		if uf[1].Find(i) == i {
			id[1][i] = ptr[1]
			ptr[1]++
		}
	}

	// step2
	N := max(ptr[0], ptr[1])
	deg[0] = make([]int, N)
	deg[1] = make([]int, N)

	// step3
	C, D := make([]int, 0, N*k), make([]int, 0, N*k)
	for i := 0; i < len(bgec._A); i++ {
		u := id[0][uf[0].Find(bgec._A[i])]
		v := id[1][uf[1].Find(bgec._B[i])]
		deg[0][u]++
		deg[1][v]++
		C = append(C, u)
		D = append(D, v)
	}

	j := 0
	for i := 0; i < N; i++ {
		for deg[0][i] < k {
			for deg[1][j] == k {
				j++
			}
			deg[0][i]++
			deg[1][j]++
			C = append(C, i)
			D = append(D, j)
		}
	}

	return &regularGraph{k: k, n: N, _A: C, _B: D}
}

func (bgec *BipartiteGraphEdgeColoring) rec(ord []int, k int) {
	if k == 0 {
		return
	}

	if k == 1 {
		bgec.ans = append(bgec.ans, ord)
		return
	}

	if k&1 == 0 {
		et := NewEulerianTrail(bgec.g.n+bgec.g.n, false)
		for _, p := range ord {
			et.AddEdge(bgec.g._A[p], bgec.g._B[p]+bgec.g.n)
		}
		paths := et.EnumerateEulerianTrail()
		path := []int{}
		for _, ps := range paths {
			for _, e := range ps {
				path = append(path, ord[e])
			}
		}

		beet := [][]int{make([]int, 0, len(path)/2), make([]int, 0, len(path)/2)}
		for i := 0; i < len(path); i++ {
			beet[i&1] = append(beet[i&1], path[i])
		}
		bgec.rec(beet[0], k/2)
		bgec.rec(beet[1], k/2)
	} else {
		flow := NewBipartiteFlow(bgec.g.n, bgec.g.n)
		for _, p := range ord {
			flow.AddEdge(bgec.g._A[p], bgec.g._B[p])
		}
		flow.MaxMatching()
		beet := []int{}
		bgec.ans = append(bgec.ans, []int{})
		for _, p := range ord {
			if flow.matchL[bgec.g._A[p]] == bgec.g._B[p] {
				flow.matchL[bgec.g._A[p]] = -1
				bgec.ans[len(bgec.ans)-1] = append(bgec.ans[len(bgec.ans)-1], p)
			} else {
				beet = append(beet, p)
			}
		}
		bgec.rec(beet, k-1)
	}
}

type EulerianTrail struct {
	g          [][][2]int
	es         [][2]int
	m          int
	usedVertex []bool
	usedEdge   []bool
	deg        []int
	directed   bool
}

func NewEulerianTrail(n int, directed bool) *EulerianTrail {
	res := &EulerianTrail{
		g:          make([][][2]int, n),
		usedVertex: make([]bool, n),
		deg:        make([]int, n),
		directed:   directed,
	}
	return res
}

func (e *EulerianTrail) AddEdge(a, b int) {
	e.es = append(e.es, [2]int{a, b})
	e.g[a] = append(e.g[a], [2]int{b, e.m})
	if e.directed {
		e.deg[a]++
		e.deg[b]--
	} else {
		e.g[b] = append(e.g[b], [2]int{a, e.m})
		e.deg[a]++
		e.deg[b]++
	}
	e.m++
}

// 枚举所有连通块的`欧拉回路`,返回边的编号.
//  如果连通块内不存在欧拉回路,返回空.
func (e *EulerianTrail) EnumerateEulerianTrail() [][]int {
	if e.directed {
		for _, d := range e.deg {
			if d != 0 {
				return [][]int{}
			}
		}
	} else {
		for _, d := range e.deg {
			if d&1 == 1 {
				return [][]int{}
			}
		}
	}

	e.usedEdge = make([]bool, e.m)
	res := [][]int{}
	for i := 0; i < len(e.g); i++ {
		if !e.usedVertex[i] && len(e.g[i]) > 0 {
			res = append(res, e.work(i))
		}
	}

	return res
}

// 枚举所有连通块的`欧拉路径`(半欧拉回路),返回边的编号.
//  如果连通块内不存在欧拉路径,返回空.
func (e *EulerianTrail) EnumerateSemiEulerianTrail() [][]int {
	uf := newUnionFindArray(len(e.g))
	for _, es := range e.es {
		uf.Union(es[0], es[1])
	}
	group := make([][]int, len(e.g))
	for i := 0; i < len(e.g); i++ {
		group[uf.Find(i)] = append(group[uf.Find(i)], i)
	}

	res := [][]int{}
	e.usedEdge = make([]bool, e.m)
	for _, vs := range group {
		if len(vs) == 0 {
			continue
		}

		latte, malta := -1, -1
		if e.directed {
			for _, p := range vs {
				if abs(e.deg[p]) > 1 {
					return [][]int{}
				} else if e.deg[p] == 1 {
					if latte >= 0 {
						return [][]int{}
					}
					latte = p
				}
			}
		} else {
			for _, p := range vs {
				if e.deg[p]&1 == 1 {
					if latte == -1 {
						latte = p
					} else if malta == -1 {
						malta = p
					} else {
						return [][]int{}
					}
				}
			}
		}

		var cur []int
		if latte == -1 {
			cur = e.work(vs[0])
		} else {
			cur = e.work(latte)
		}

		if len(cur) > 0 {
			res = append(res, cur)
		}
	}

	return res
}

func (e *EulerianTrail) GetEdge(index int) (int, int) {
	return e.es[index][0], e.es[index][1]
}

func (e *EulerianTrail) work(s int) []int {
	st := [][2]int{}
	ord := []int{}
	st = append(st, [2]int{s, -1})
	for len(st) > 0 {
		index := st[len(st)-1][0]
		e.usedVertex[index] = true
		if len(e.g[index]) == 0 {
			ord = append(ord, st[len(st)-1][1])
			st = st[:len(st)-1]
		} else {
			e_ := e.g[index][len(e.g[index])-1]
			e.g[index] = e.g[index][:len(e.g[index])-1]
			if e.usedEdge[e_[1]] {
				continue
			}
			e.usedEdge[e_[1]] = true
			st = append(st, [2]int{e_[0], e_[1]})
		}
	}

	ord = ord[:len(ord)-1]
	for i, j := 0, len(ord)-1; i < j; i, j = i+1, j-1 {
		ord[i], ord[j] = ord[j], ord[i]
	}
	return ord
}

type BipartiteFlow struct {
	n, m, timeStamp int
	g, rg           [][]int
	matchL, matchR  []int
	dist            []int
	used            []int
	alive           []bool
	matched         bool
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
		n:      n,
		m:      m,
		g:      g,
		rg:     rg,
		matchL: matchL,
		matchR: matchR,
		used:   used,
		alive:  alive,
	}
}

// 增加一条边u-v.u属于左侧点集，v属于右侧点集.
func (bf *BipartiteFlow) AddEdge(u, v int) {
	bf.g[u] = append(bf.g[u], v)
	bf.rg[v] = append(bf.rg[v], u)
}

// /* http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3198 */
func (bt *BipartiteFlow) EraseEdge(u, v int) {
	if bt.matchL[u] == v {
		bt.matchL[u] = -1
		bt.matchR[v] = -1
	}
	// remove v in bt.g[u]
	// remove u in bt.rg[v]
	for i := 0; i < len(bt.g[u]); i++ {
		if bt.g[u][i] == v {
			bt.g[u] = append(bt.g[u][:i], bt.g[u][i+1:]...)
			break
		}
	}
	for i := 0; i < len(bt.rg[v]); i++ {
		if bt.rg[v][i] == u {
			bt.rg[v] = append(bt.rg[v][:i], bt.rg[v][i+1:]...)
			break
		}
	}
}

// 求最大匹配.
func (bf *BipartiteFlow) MaxMatching() [][2]int {
	bf.matched = true
	for {
		bf.buildAugmentPath()
		bf.timeStamp++
		flow := 0
		for i := 0; i < bf.n; i++ {
			if bf.matchL[i] == -1 {
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
	for i := 0; i < bf.n; i++ {
		if bf.matchL[i] >= 0 {
			res = append(res, [2]int{i, bf.matchL[i]})
		}
	}
	return res
}

// 字典序最小的最大匹配.
// /* http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0334 */
func (bt *BipartiteFlow) LexMaxMatching() [][2]int {
	if !bt.matched {
		bt.MaxMatching()
	}
	for _, vs := range bt.g {
		sort.Ints(vs) // 字典序最小
	}
	es := [][2]int{}
	for i := 0; i < bt.n; i++ {
		if bt.matchL[i] == -1 || !bt.alive[i] {
			continue
		}
		bt.matchR[bt.matchL[i]] = -1
		bt.matchL[i] = -1
		bt.timeStamp++
		bt.findAugmentPath(i)
		bt.alive[i] = false
		es = append(es, [2]int{i, bt.matchL[i]})
	}
	return es
}

// 最小点覆盖.
func (bt *BipartiteFlow) MinVertexCover() []int {
	visited := bt.findResidualPath()
	res := []int{}
	for i := 0; i < (bt.n + bt.m); i++ {
		if visited[i] != (i < bt.n) {
			res = append(res, i)
		}
	}
	return res
}

// 字典序(ord优先度顺序)最小点覆盖.
//  /* https://atcoder.jp/contests/utpc2013/tasks/utpc2013_11 */
func (bt *BipartiteFlow) LexMinVertexCover(ord []int) []int {
	if len(ord) != bt.n+bt.m {
		panic("len(ord) != bt.n+bt.m")
	}
	res := bt.BuildRisidualGraph()
	rRes := make([][]int, bt.n+bt.m+2)
	for i := 0; i < bt.n+bt.m+2; i++ {
		for _, to := range res[i] {
			rRes[to] = append(rRes[to], i)
		}
	}

	que := []int{}
	visited := make([]int8, bt.n+bt.m+2)
	for i := range visited {
		visited[i] = -1
	}

	expandLeft := func(t int) {
		if visited[t] != -1 {
			return
		}
		que = append(que, t)
		visited[t] = 1
		for len(que) > 0 {
			v := que[0]
			que = que[1:]
			for _, to := range rRes[v] {
				if visited[to] != -1 {
					continue
				}
				que = append(que, to)
				visited[to] = 1
			}
		}
	}
	expandRight := func(t int) {
		if visited[t] != -1 {
			return
		}
		que = append(que, t)
		visited[t] = 0
		for len(que) > 0 {
			v := que[0]
			que = que[1:]
			for _, to := range res[v] {
				if visited[to] != -1 {
					continue
				}
				que = append(que, to)
				visited[to] = 0
			}
		}
	}

	expandRight(bt.n + bt.m)
	expandLeft(bt.n + bt.m + 1)
	ret := []int{}
	for _, v := range ord {
		if v < bt.n {
			expandLeft(v)
			if visited[v]&1 != 0 { // visited[v] != 0
				ret = append(ret, v)
			}
		} else {
			expandRight(v)
			if (^visited[v] & 1) != 0 { // visited[v] == 0
				ret = append(ret, v)
			}
		}
	}
	return ret
}

// 最大独立集.
func (bt *BipartiteFlow) MaxIndependentSet() []int {
	visited := bt.findResidualPath()
	res := []int{}
	for i := 0; i < (bt.n + bt.m); i++ {
		if visited[i] != (i >= bt.n) {
			res = append(res, i)
		}
	}
	return res
}

// 最小边覆盖.
func (bt *BipartiteFlow) MinEdgeCover() [][2]int {
	es := bt.MaxMatching()
	for i := 0; i < bt.n; i++ {
		if bt.matchL[i] >= 0 {
			continue
		}
		if len(bt.g[i]) == 0 {
			return [][2]int{}
		}
		es = append(es, [2]int{i, bt.g[i][0]})
	}
	for i := 0; i < bt.m; i++ {
		if bt.matchR[i] >= 0 {
			continue
		}
		if len(bt.rg[i]) == 0 {
			return [][2]int{}
		}
		es = append(es, [2]int{bt.rg[i][0], i})
	}
	return es
}

// 构建残量图.
//  left: [0,n), right: [n,n+m), S: n+m, T: n+m+1
func (bf *BipartiteFlow) BuildRisidualGraph() [][]int {
	if !bf.matched {
		bf.MaxMatching()
	}

	S := bf.n + bf.m
	T := bf.n + bf.m + 1
	ris := make([][]int, bf.n+bf.m+2)
	for i := 0; i < bf.n; i++ {
		if bf.matchL[i] == -1 {
			ris[S] = append(ris[S], i)
		} else {
			ris[i] = append(ris[i], S)
		}
	}

	for i := 0; i < bf.m; i++ {
		if bf.matchR[i] == -1 {
			ris[i+bf.n] = append(ris[i+bf.n], T)
		} else {
			ris[T] = append(ris[T], i+bf.n)
		}
	}

	for i := 0; i < bf.n; i++ {
		for _, j := range bf.g[i] {
			if bf.matchL[i] == j {
				ris[j+bf.n] = append(ris[j+bf.n], i)
			} else {
				ris[i] = append(ris[i], j+bf.n)
			}
		}
	}

	return ris
}

func (bf *BipartiteFlow) findResidualPath() []bool {
	res := bf.BuildRisidualGraph()
	que := []int{}
	visited := make([]bool, bf.n+bf.m+2)
	que = append(que, bf.n+bf.m)
	visited[bf.n+bf.m] = true
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
	for i := 0; i < bf.n; i++ {
		if bf.matchL[i] == -1 {
			que = append(que, i)
			bf.dist[i] = 0
		}
	}
	for len(que) > 0 {
		a := que[0]
		que = que[1:]
		for _, b := range bf.g[a] {
			c := bf.matchR[b]
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
		c := bf.matchR[b]
		if c < 0 || (bf.used[c] != bf.timeStamp && bf.dist[c] == bf.dist[a]+1 && bf.findMinDistAugmentPath(c)) {
			bf.matchR[b] = a
			bf.matchL[a] = b
			return true
		}
	}
	return false
}

func (bf *BipartiteFlow) findAugmentPath(a int) bool {
	bf.used[a] = bf.timeStamp
	for _, b := range bf.g[a] {
		c := bf.matchR[b]
		if c < 0 || (bf.alive[c] && bf.used[c] != bf.timeStamp && bf.findAugmentPath(c)) {
			bf.matchR[b] = a
			bf.matchL[a] = b
			return true
		}
	}
	return false
}

func newUnionFindArray(n int) *_unionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_unionFindArray{
		Part:   n,
		size:   n,
		rank:   rank,
		parent: parent,
	}
}

type _unionFindArray struct {
	size   int
	Part   int
	rank   []int
	parent []int
}

func (ufa *_unionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_unionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_unionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_unionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

type H = []int

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
