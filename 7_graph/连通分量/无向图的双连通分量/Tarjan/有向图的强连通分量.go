// https://ei1333.github.io/library/graph/connected-components/strongly-connected-components.hpp
// 与えられた有向グラフを強連結成分分解する.
// グラフの任意の 2頂点間に有効路が存在するとき, 有向グラフが強連結であるとよぶ.
// 強連結成分は, 極大で強連結な部分グラフである.
// 適当な頂点からDFSをして, 帰りがけ順に頂点を列挙することを, 未訪問の頂点がある間繰り返す.
// 次に辺をすべて逆向きにしたグラフについて, 列挙した頂点の逆順にDFS する.
// 1回の DFS で到達できた頂点が1つの強連結成分となる.
// 強連結成分を縮約後の頂点とそれらを結ぶ辺からなるグラフはDAGになっている.

// build(): 強連結成分分解する.
//  dag には縮約後の頂点と辺からなるDAGが格納される.
//  comp には各頂点が属する強連結成分の頂点番号が格納される.(相等则说明两点可以互相到达)
//  group には各強連結成分について, それに属する頂点が格納される.(group组之间按照拓扑序排序)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1813
	// 不等关系:有向边; 全部相等:强连通(环)
	// 给定一个DAG 求将DAG变为一个环(强连通分量)的最少需要添加的边数
	// !答案为 `max(入度为0的点的个数, 出度为0的点的个数)`

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	scc := NewStronglyConnectedComponents(n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		scc.AddEdge(u-1, v-1, 1)
	}
	scc.Build()

	if len(scc.Group) == 1 { // 缩成一个点了,说明是强连通的
		fmt.Fprintln(out, 0)
		return
	}

	g := len(scc.Group)
	indeg, outDeg := make([]int, g), make([]int, g)
	for i := 0; i < g; i++ {
		for _, next := range scc.Dag[i] {
			indeg[next]++
			outDeg[i]++
		}
	}

	in0, out0 := 0, 0
	for i := 0; i < g; i++ {
		if indeg[i] == 0 {
			in0++
		}
		if outDeg[i] == 0 {
			out0++
		}
	}

	fmt.Fprintln(out, max(in0, out0))
}

func yuki1293() {
	// https://yukicoder.me/problems/no/1293
	// No.1293 2種類の道路-SCC
	// 无向图中有两种路径,各有road1,road2条
	// 求有多少个二元组(a,b),满足从a到b经过 '若干条第一种路径+若干条第二种路径'

	// !每个点i拆成点2*i和点2*i+1,2*i->2*i+1
	// !第一种路径: 2*i<->2*j
	// !第二种路径: 2*i+1<->2*j+1
	// 然后对每个顶点求出有多少个可以到达自己

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, road1, road2 int
	fmt.Fscan(in, &n, &road1, &road2)
	scc := NewStronglyConnectedComponents(2 * n)
	for i := 0; i < road1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		scc.AddEdge(2*a, 2*b, 1)
		scc.AddEdge(2*b, 2*a, 1)
	}
	for i := 0; i < road2; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		scc.AddEdge(2*a+1, 2*b+1, 1)
		scc.AddEdge(2*b+1, 2*a+1, 1)
	}
	for i := 0; i < n; i++ {
		scc.AddEdge(2*i, 2*i+1, 1)
	}
	scc.Build()

	v := len(scc.Group)
	dp := make([]int, v)
	for i := 0; i < n; i++ {
		dp[scc.CompId[2*i]]++
	}
	for i := 0; i < v; i++ {
		for _, to := range scc.Dag[i] {
			dp[to] += dp[i]
		}
	}

	res := 0
	for i := 0; i < n; i++ {
		res += dp[scc.CompId[2*i+1]] - 1 // !减去自己到自己的路径1
	}
	fmt.Fprintln(out, res)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type WeightedEdge struct{ from, to, cost, index int }
type StronglyConnectedComponents struct {
	G      [][]WeightedEdge // 原图
	Dag    [][]int          // 强连通分量缩点后的DAG(有向图邻接表)
	CompId []int            // 每个顶点所属的强连通分量的编号
	Group  [][]int          // 每个强连通分量所包含的顶点
	rg     [][]WeightedEdge
	order  []int
	used   []bool
	eid    int
}

func NewStronglyConnectedComponents(n int) *StronglyConnectedComponents {
	return &StronglyConnectedComponents{G: make([][]WeightedEdge, n)}
}

func (scc *StronglyConnectedComponents) AddEdge(from, to, cost int) {
	scc.G[from] = append(scc.G[from], WeightedEdge{from, to, cost, scc.eid})
	scc.eid++
}

func (scc *StronglyConnectedComponents) Build() {
	scc.rg = make([][]WeightedEdge, len(scc.G))
	for i := range scc.G {
		for _, e := range scc.G[i] {
			scc.rg[e.to] = append(scc.rg[e.to], WeightedEdge{e.to, e.from, e.cost, e.index})
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

	dag := make([][]int, ptr)
	visited := make(map[int]struct{}) // 边去重
	for i := range scc.G {
		for _, e := range scc.G[i] {
			x, y := scc.CompId[e.from], scc.CompId[e.to]
			if x == y {
				continue // 原来的边 x->y 的顶点在同一个强连通分量内,可以汇合同一个 SCC 的权值
			}
			hash := x*len(scc.G) + y
			if _, ok := visited[hash]; !ok {
				dag[x] = append(dag[x], y)
				visited[hash] = struct{}{}
			}
		}
	}
	scc.Dag = dag

	scc.Group = make([][]int, ptr)
	for i := range scc.G {
		scc.Group[scc.CompId[i]] = append(scc.Group[scc.CompId[i]], i)
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
	for _, e := range scc.G[idx] {
		scc.dfs(e.to)
	}
	scc.order = append(scc.order, idx)
}

func (scc *StronglyConnectedComponents) rdfs(idx int, cnt int) {
	if scc.CompId[idx] != -1 {
		return
	}
	scc.CompId[idx] = cnt
	for _, e := range scc.rg[idx] {
		scc.rdfs(e.to, cnt)
	}
}
