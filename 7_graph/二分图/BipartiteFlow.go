// https://ei1333.github.io/library/graph/flow/bipartite-flow.hpp
// 二分图网络流

// 概要
// 二部グラフに対するフロー.
// 最大マッチングは Hopcroft-Karp に基づく実装.
// 最大流を求める Dinic と同じアルゴリズムだが,
// Hopcroft-Karp はこれを二部グラフ用に書き換えたもので定数倍が軽い.
// 残余グラフをBFSして各頂点までの最短距離を計算し, 最短距離のパスをDFSで見つけてフローを流す.
// 計算量が良いため基本的にはこれを使うと良い.

// DAG的最小路径覆盖定义：在一个有向图中，找出最少的路径，使得这些路径经过了所有的点。

// 使い方
// BipartiteFlow(n, m) 左側の頂点数 n, 右側の頂点数 m で初期化する.
// add_edge(u, v) 頂点 u, v 間に辺を張る.
// max_matching():= 最大マッチングを返す.
// erase_edge(a, b):= 頂点 u, v 間にある辺を削除する.
// min_vertex_cover():= 最小頂点被覆を返す. (最小点覆盖)
// max_independent_set():= 最大安定集合を返す. (最大独立集合)
// min_edge_cover():= 最小辺被覆を返す. (最小辺覆盖)
// lex_max_matching():= 辞書順最小の最大マッチングを返す. (字典序最小的最大匹配)
// lex_min_vertex_cover(ord):= 辞書順(優先度順)最小頂点被覆を返す.
//  (按照ord里的点的优先顺序排序的最小点覆盖)
//  ord里表示左侧顶点[0,n)和右侧顶点[n,n+m)的优先顺序
// BuildRisidualGraph() 残余グラフを構築する.
//  左侧顶点为[0,n)，右侧顶点为[n,n+m), 起点为n+m, 终点为n+m+1.
// !残量网络的性质
//!   - 判断每条边是否`一定`存在于最大匹配中
//      对残量网络进行强连通分量分解，对每条边 (u,v),按以下顺序判断 (一定在最大匹配中/一定不在最大匹配中/不确定):
//      1. u,v在同一个强连通分量中 => 不确定
//      2. 边(u,v)存在于最大匹配中 => 一定在最大匹配中
//      3. 否则u和v属于不同的强连通分量 => 一定不在最大匹配中
//      删边的情况参考 二分图中每条边是否存在于最大匹配中
//!   - 判断`每个点`是否`一定`存在于最大匹配中 => 二分图博弈问题:
//      朴素的想法:去除这个顶点之后,如果最大流发生变化,则这个顶点一定在最大匹配中
//!     进一步的方法:在残量网络中,对二分图左边的点,从虚拟源点出发bfs,路径上的左侧的点`一定不在`最大匹配中(去除这些点之后,最大流不变)
//      反之则一定在最大匹配中
//!     对二分图右边的点,在残量图的反图上，从虚拟汇点出发bfs,路径上的右侧的点`一定不在`最大匹配中
//      删点的情况参考 G - Vertex Deletion-二分图最大匹配中一定包含的点的个数

// 加边是O(1)的 删边是O(V)的
// 其余操作都是O(Esqrt(V))的

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func Amidakuji() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0334
	// 字典序最小匹配

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	n--
	bf := NewBipartiteFlow(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var x int
			fmt.Fscan(in, &x)
			if x != 0 {
				bf.AddEdge(j, i)
			}
		}
	}

	v := bf.LexMaxMatching()
	if len(v) < n {
		fmt.Fprintln(out, "no")
	} else {
		fmt.Fprintln(out, "yes")
		for _, e := range v {
			fmt.Fprintln(out, e[1]+1)
		}
	}
}

func ArrangementOfPieces() {
	// 删边操作+判断完美匹配(每行每列恰好一个棋子)
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3198
	// n*n的棋盘上有m个给定的棋子 (n<=5000)
	// 现在有q个操作 (xi,yi)
	// 如果原来位置可以放棋子，就变成不能放棋子，否则变成可以放棋子
	// 每次操作后,能否做到在棋盘上放n个棋子,使得每行每列都恰好有一个棋子 (完美匹配)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	flow := NewBipartiteFlow(n, n)
	mas := map[struct{ x, y int }]struct{}{}
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x, y = x-1, y-1
		flow.AddEdge(x, y)
		mas[struct{ x, y int }{x, y}] = struct{}{}
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x, y = x-1, y-1
		if _, ok := mas[struct{ x, y int }{x, y}]; ok {
			flow.EraseEdge(x, y)
			delete(mas, struct{ x, y int }{x, y})
		} else {
			flow.AddEdge(x, y)
			mas[struct{ x, y int }{x, y}] = struct{}{}
		}
		if len(flow.MaxMatching()) == n {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

func MatchingOnBipartiteGraph() {
	// https://judge.yosupo.jp/problem/bipartitematching
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var L, R, M int
	fmt.Fscan(in, &L, &R, &M)
	bf := NewBipartiteFlow(L, R)
	for i := 0; i < M; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		bf.AddEdge(u, v)
	}
	res := bf.MaxMatching()
	fmt.Fprintln(out, len(res))
	for _, e := range res {
		fmt.Fprintln(out, e[0], e[1])
	}
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
//  !0<=u<n,0<=v<m.
func (bf *BipartiteFlow) AddEdge(u, v int) {
	bf.g[u] = append(bf.g[u], v)
	bf.rg[v] = append(bf.rg[v], u)
}

// /* http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3198 */
func (bt *BipartiteFlow) EraseEdge(u, v int) {
	if bt.MatchL[u] == v {
		bt.MatchL[u] = -1
		bt.MatchR[v] = -1
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
//  返回(左侧点,右侧点)的匹配对.
//  !0<=左侧点<n,0<=右侧点<m.
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
	for i := 0; i < bt.N; i++ {
		if bt.MatchL[i] == -1 || !bt.alive[i] {
			continue
		}
		bt.MatchR[bt.MatchL[i]] = -1
		bt.MatchL[i] = -1
		bt.timeStamp++
		bt.findAugmentPath(i)
		bt.alive[i] = false
		es = append(es, [2]int{i, bt.MatchL[i]})
	}
	return es
}

// 最小点覆盖.
func (bt *BipartiteFlow) MinVertexCover() []int {
	visited := bt.findResidualPath()
	res := []int{}
	for i := 0; i < (bt.N + bt.M); i++ {
		if visited[i] != (i < bt.N) {
			res = append(res, i)
		}
	}
	return res
}

// 字典序(ord优先度顺序)最小点覆盖.
//  /* https://atcoder.jp/contests/utpc2013/tasks/utpc2013_11 */
func (bt *BipartiteFlow) LexMinVertexCover(ord []int) []int {
	if len(ord) != bt.N+bt.M {
		panic("len(ord) != bt.n+bt.m")
	}
	res := bt.BuildRisidualGraph()
	rRes := make([][]int, bt.N+bt.M+2)
	for i := 0; i < bt.N+bt.M+2; i++ {
		for _, to := range res[i] {
			rRes[to] = append(rRes[to], i)
		}
	}

	que := []int{}
	visited := make([]int8, bt.N+bt.M+2)
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

	expandRight(bt.N + bt.M)
	expandLeft(bt.N + bt.M + 1)
	ret := []int{}
	for _, v := range ord {
		if v < bt.N {
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
	for i := 0; i < (bt.N + bt.M); i++ {
		if visited[i] != (i >= bt.N) {
			res = append(res, i)
		}
	}
	return res
}

// 最小边覆盖.
func (bt *BipartiteFlow) MinEdgeCover() [][2]int {
	es := bt.MaxMatching()
	for i := 0; i < bt.N; i++ {
		if bt.MatchL[i] >= 0 {
			continue
		}
		if len(bt.g[i]) == 0 {
			return [][2]int{}
		}
		es = append(es, [2]int{i, bt.g[i][0]})
	}
	for i := 0; i < bt.M; i++ {
		if bt.MatchR[i] >= 0 {
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

func (bf *BipartiteFlow) findAugmentPath(a int) bool {
	bf.used[a] = bf.timeStamp
	for _, b := range bf.g[a] {
		c := bf.MatchR[b]
		if c < 0 || (bf.alive[c] && bf.used[c] != bf.timeStamp && bf.findAugmentPath(c)) {
			bf.MatchR[b] = a
			bf.MatchL[a] = b
			return true
		}
	}
	return false
}
