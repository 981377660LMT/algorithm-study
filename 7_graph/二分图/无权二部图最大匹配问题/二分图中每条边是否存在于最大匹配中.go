//!   - 判断每条边是否`一定`存在于最大匹配中
//      对残量网络进行强连通分量分解，对每条边 (u,v),按以下顺序判断 (一定在最大匹配中/一定不在最大匹配中/不确定):
//      1. u,v在同一个强连通分量中 => 不确定
//      2. 边(u,v)存在于最大匹配中 => 一定在最大匹配中
//      3. 否则u和v属于不同的强连通分量 => 一定不在最大匹配中
//      删边的情况参考 二分图中每条边是否存在于最大匹配中

// 二分图上的dinic是O(nsqrt(n))的

package main

import "fmt"

func main() {
	// 2
	// 1 2
	n := 4
	edges := [][2]int{{0, 1}, {2, 3}, {2, 1}}
	fmt.Println(solve(n, edges))
}

// 判断每条边是否`存在于最大匹配中,保证图是二分图
//  0:不确定,1:一定在最大匹配中,2:一定不在最大匹配中
func solve(n int, edges [][2]int) []int {
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

	ids := make([]int, n) // 原图中的点在二分图中的编号 左边:0-L-1 右边:L-n-1
	c1, c2 := 0, 0
	for i := 0; i < n; i++ {
		if colors[i] == 0 {
			ids[i] = c1
			c1++
		} else {
			ids[i] = c2 + L
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

	match := bf.MaxMatching()
	matchedEdge := make(map[int]struct{})
	for _, e := range match {
		u, v := e[0], e[1]+L
		matchedEdge[u*n+v] = struct{}{}
	}
	g := bf.BuildRisidualGraph() // 残量网络
	scc := NewStronglyConnectedComponents(len(g))
	for i := 0; i < len(g); i++ {
		for _, v := range g[i] {
			scc.AddEdge(i, v)
		}
	}
	scc.Build()

	states := make([]int, len(edges)) // 0:未知 1:一定在最大匹配中 -1:一定不在最大匹配中
	for i, e := range edges {
		u, v := e[0], e[1]
		if colors[u] == 1 {
			u, v = v, u
		}
		if scc.CompId[ids[u]] == scc.CompId[ids[v]] {
			break
		}
		if _, ok := matchedEdge[ids[u]*n+ids[v]]; ok {
			states[i] = 1
		} else {
			states[i] = -1
		}
	}

	return states
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

// SCC
type StronglyConnectedComponents struct {
	G      [][]int // 原图
	CompId []int   // 每个顶点所属的强连通分量的编号
	rg     [][]int
	order  []int
	used   []bool
	eid    int
}

func NewStronglyConnectedComponents(n int) *StronglyConnectedComponents {
	return &StronglyConnectedComponents{G: make([][]int, n)}
}

func (scc *StronglyConnectedComponents) AddEdge(from, to int) {
	scc.G[from] = append(scc.G[from], to)
}

func (scc *StronglyConnectedComponents) Build() {
	scc.rg = make([][]int, len(scc.G))
	for i := range scc.G {
		for _, next := range scc.G[i] {
			scc.rg[next] = append(scc.rg[next], i)
		}
	}

	scc.CompId = make([]int, len(scc.G))
	for i := range scc.CompId {
		scc.CompId[i] = -1
	}
	scc.used = make([]bool, len(scc.G))
	for i := range scc.G {
		scc.dfs(i)
	}
	for i, j := 0, len(scc.order)-1; i < j; i, j = i+1, j-1 {
		scc.order[i], scc.order[j] = scc.order[j], scc.order[i]
	}

	ptr := 0
	for _, v := range scc.order {
		if scc.CompId[v] == -1 {
			scc.rdfs(v, ptr)
			ptr++
		}
	}

}

// 获取顶点k所属的强连通分量的编号
func (scc *StronglyConnectedComponents) Get(k int) int {
	return scc.CompId[k]
}

func (scc *StronglyConnectedComponents) dfs(idx int) {
	tmp := scc.used[idx]
	scc.used[idx] = true
	if tmp {
		return
	}
	for _, next := range scc.G[idx] {
		scc.dfs(next)
	}
	scc.order = append(scc.order, idx)
}

func (scc *StronglyConnectedComponents) rdfs(idx int, cnt int) {
	if scc.CompId[idx] != -1 {
		return
	}
	scc.CompId[idx] = cnt
	for _, next := range scc.rg[idx] {
		scc.rdfs(next, cnt)
	}
}
