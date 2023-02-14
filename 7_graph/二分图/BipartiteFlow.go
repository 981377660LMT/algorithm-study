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
// build_residual_graph() 残余グラフを構築する.
//  左侧顶点为[0,n)，右侧顶点为[n,n+m), 起点为n+m, 终点为n+m+1.
//  对残量网络进行强连通分量分解，对每条边 (u,v):
//  1. 如果u,v在同一个强连通分量中:不一定在最大匹配中
//  2. 如果边(u,v)在最大匹配中:在最大匹配中一定用到
//  3. 否则,u和v属于不同的强连通分量:一定不在最大匹配中

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
