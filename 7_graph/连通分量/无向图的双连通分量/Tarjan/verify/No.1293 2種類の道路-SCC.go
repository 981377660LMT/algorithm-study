// https://yukicoder.me/problems/no/1293
// 无向图中有两种路径,各有road1,road2条
// 求有多少个二元组(a,b),满足从a到b经过 '若干条第一种路径+若干条第二种路径'

// !每个点i拆成点2*i和点2*i+1,2*i->2*i+1
// !第一种路径: 2*i<->2*j
// !第二种路径: 2*i+1<->2*j+1
// 然后对每个顶点求出有多少个可以到达自己

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
		for _, e := range scc.Dag[i] {
			dp[e.to] += dp[e.from]
		}
	}

	res := 0
	for i := 0; i < n; i++ {
		res += dp[scc.CompId[2*i+1]] - 1 // !减去自己到自己的路径1
	}
	fmt.Fprintln(out, res)
}

type WeightedEdge struct{ from, to, cost, index int }
type StronglyConnectedComponents struct {
	G      [][]WeightedEdge // 原图
	Dag    [][]WeightedEdge // 强连通分量缩点后的顶点和边组成的DAG
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

	dag := make([][]WeightedEdge, ptr)
	visited := make(map[int]struct{}) // !去重
	for i := range scc.G {
		for _, e := range scc.G[i] {
			x, y := scc.CompId[e.from], scc.CompId[e.to]
			if x == y {
				continue
			}
			hash := x*len(scc.G) + y
			if _, ok := visited[hash]; !ok {
				dag[x] = append(dag[x], WeightedEdge{x, y, e.cost, e.index})
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
