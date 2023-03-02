// https://atcoder.jp/contests/abc223/tasks/abc223_g
// G - Vertex Deletion
//!   - 判断二分图中每个点是否`一定`存在于最大匹配中 => 二分图博弈问题:
//    朴素的想法:去除这个顶点之后,如果最大流发生变化,则这个顶点一定在最大匹配中
//!    进一步的方法:在残量网络中,对二分图左边的点,从虚拟源点出发bfs,路径上的点`一定不在`最大匹配中(去除这些点之后,最大流不变)
//     反之则一定在最大匹配中
//!    对二分图右边的点，在残量图的反图上，从虚拟汇点出发bfs，路径上的二分图右侧的点一定不在最大匹配中。

// 二分图上的dinic是O(nsqrt(n))的

// ps: 二分图博弈(另一种不带颜色的是公平博弈(grundy数)):
// https://zhuanlan.zhihu.com/p/555764217
// https://zhuanlan.zhihu.com/p/359334008
// 两个人轮流移动无向图上的一个点，已经游览过的点不能再次访问，询问是先手还是后手胜利问题。
// !如果起点是二分图最大匹配中一定包含的点，那么先手必胜；反正就是后手必胜。
// 一个经典的二分图博弈模型是在国际象棋`棋盘`上，双方轮流移动一个士兵，不能走已经走过的格子，问谁先无路可走。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://atcoder.jp/contests/abc223/tasks/abc223_g
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([][]int, 0, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v})
	}

	important := solve(n, edges)
	canRemove := 0
	for _, v := range important {
		if !v {
			canRemove++
		}
	}
	fmt.Fprintln(out, canRemove)
}

func P4055游戏(grid []string) [][2]int {
	// https://www.luogu.com.cn/problem/P4055
	// 返回先手能够赢得游戏的起点(不一定在最大匹配中的空地)
	ROW, COL := len(grid), len(grid[0])
	edges := [][]int{}
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			if grid[i][j] == '#' {
				continue
			}
			if i+1 < ROW && grid[i+1][j] == '.' {
				edges = append(edges, []int{i*COL + j, (i+1)*COL + j})
			}
			if j+1 < COL && grid[i][j+1] == '.' {
				edges = append(edges, []int{i*COL + j, i*COL + j + 1})
			}
		}
	}

	states := solve(ROW*COL, edges)
	res := [][2]int{}

	for i, ok := range states {
		if !ok && grid[i/COL][i%COL] == '.' { // !不一定在最大匹配中的空地
			res = append(res, [2]int{i / COL, i % COL})
		}
	}
	return res
}

// 判断二分图中每个点是否`一定`存在于最大匹配中(删去这个点之后,最大流减小)
//  保证图是二分图
func solve(n int, edges [][]int) []bool {
	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	L, R := 0, 0
	colors := make([]int, n) // 0-left, 1-right
	for i := 0; i < n; i++ {
		colors[i] = -1
	}
	var dfs func(u, c int)
	dfs = func(u, c int) {
		colors[u] = c
		if c == 0 {
			L++
		} else {
			R++
		}
		for _, v := range adjList[u] {
			if colors[v] == -1 {
				dfs(v, c^1)
			}
		}
	}
	for i := 0; i < n; i++ {
		if colors[i] == -1 {
			dfs(i, 0)
		}
	}

	ids := make([]int, n)  // 原图中的点在二分图中的编号 左边:0-L-1 右边:L-n-1
	rids := make([]int, n) // 二分图中的点在原图中的编号 0-n-1 => 0-n-1
	c1, c2 := 0, 0
	for i := 0; i < n; i++ {
		if colors[i] == 0 {
			ids[i] = c1
			rids[c1] = i
			c1++
		} else {
			ids[i] = c2 + L
			rids[c2+L] = i
			c2++
		}
	}

	bf := NewBipartiteFlow(L, R)
	for _, e := range edges {
		u, v := e[0], e[1]
		if colors[u] == 1 {
			u, v = v, u
		}
		bf.AddEdge(ids[u], ids[v]-L)
	}

	g := bf.BuildRisidualGraph() // 残量网络
	rg := make([][]int, len(g))  // 残量网络的反图
	for i := 0; i < len(g); i++ {
		for _, v := range g[i] {
			rg[v] = append(rg[v], i)
		}
	}

	important := make([]bool, n)
	for i := range important {
		important[i] = true
	}
	bfs := func(graph [][]int, start, begin, end int) {
		vis := make([]bool, len(graph))
		q := []int{start}
		vis[start] = true
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			if begin <= cur && cur < end {
				important[rids[cur]] = false // 从虚拟源点出发能到达的左侧点/从虚拟汇点出发能到达的右侧点
			}
			for _, v := range graph[cur] {
				if !vis[v] {
					vis[v] = true
					q = append(q, v)
				}
			}
		}

	}

	bfs(g, n, 0, L)
	bfs(rg, n+1, L, n)
	return important
}

type BipartiteFlow struct {
	N, M, timeStamp int
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
		N:      n,
		M:      m,
		g:      g,
		rg:     rg,
		matchL: matchL,
		matchR: matchR,
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
	for i := 0; i < bf.N; i++ {
		if bf.matchL[i] >= 0 {
			res = append(res, [2]int{i, bf.matchL[i]})
		}
	}
	return res
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
		if bf.matchL[i] == -1 {
			ris[S] = append(ris[S], i)
		} else {
			ris[i] = append(ris[i], S)
		}
	}

	for i := 0; i < bf.M; i++ {
		if bf.matchR[i] == -1 {
			ris[i+bf.N] = append(ris[i+bf.N], T)
		} else {
			ris[T] = append(ris[T], i+bf.N)
		}
	}

	for i := 0; i < bf.N; i++ {
		for _, j := range bf.g[i] {
			if bf.matchL[i] == j {
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
