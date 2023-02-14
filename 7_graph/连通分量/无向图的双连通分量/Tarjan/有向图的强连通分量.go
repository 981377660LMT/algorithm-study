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
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	scc := NewStronglyConnectedComponents(n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		scc.AddEdge(a, b, 1)
	}
	scc.Build()
	fmt.Fprintln(out, len(scc.Group))
	for i := 0; i < len(scc.Group); i++ {
		fmt.Fprint(out, len(scc.Group[i]))
		for j := 0; j < len(scc.Group[i]); j++ {
			fmt.Fprint(out, " ", scc.Group[i][j])
		}
		fmt.Fprintln(out)
	}
}

type WeightedEdge struct{ from, to, cost, index int }
type StronglyConnectedComponents struct {
	G      [][]WeightedEdge // 原图
	Dag    [][]WeightedEdge // 强连通分量缩点后的顶点和边组成的DAG
	CompId []int            //每个顶点所属的强连通分量的编号
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
	for i := range scc.G {
		for _, e := range scc.G[i] {
			x, y := scc.CompId[e.from], scc.CompId[e.to]
			if x == y {
				continue
			}
			dag[x] = append(dag[x], WeightedEdge{x, y, e.cost, e.index})
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
