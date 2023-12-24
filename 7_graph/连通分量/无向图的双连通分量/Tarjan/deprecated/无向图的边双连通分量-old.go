// 二辺連結成分分解とも. 二重辺連結成分とは,
// !1本の辺を取り除いても連結である部分グラフである.
// つまり, 橋を含まない部分グラフなので, 橋を列挙することで二重辺連結成分を列挙できる.
// !二重辺連結成分で縮約後の頂点と橋からなるグラフは森になっている.

// build(): 二重辺連結成分分解する.
//  tree には縮約後の頂点からなる森が格納される.
//  comp には各頂点が属する二重辺連結成分の頂点番号が格納される.
//  group には各二重辺連結成分について, それに属する頂点が格納される.

// !注意这里把孤立的点也当作一个点双
// 桥的两个端点所在的连通分量不同

// 问题特点:边`只能经过一次`(不能走桥)、a到b的路径是否`存在且唯一`

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	缩点成树()
}
func yosupu() {
	// https://judge.yosupo.jp/submission/125538
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	g := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], Edge{u, v, 1, i})
		g[v] = append(g[v], Edge{v, u, 1, i})
	}
	tec := NewTwoEdgeConnectedComponents(g)
	tec.Build()
	fmt.Fprintln(out, len(tec.Group))
	for _, p := range tec.Group {
		fmt.Fprint(out, len(p))
		for _, v := range p {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

// 新树的边为原图的桥
func 缩点成树() {
	n := 5
	edges := [][2]int{{0, 1}, {1, 2}, {2, 0}, {1, 3}, {1, 4}}
	graph := make([][]Edge, n)
	for i := 0; i < len(edges); i++ {
		u, v := edges[i][0], edges[i][1]
		graph[u] = append(graph[u], Edge{u, v, 1, i})
		graph[v] = append(graph[v], Edge{v, u, 1, i})
	}
	EBCC := NewTwoEdgeConnectedComponents(graph)
	EBCC.Build()

	tree := make([][]int, len(EBCC.Group))
	for _, e := range edges {
		u, v := e[0], e[1]
		id1, id2 := EBCC.CompId[u], EBCC.CompId[v]
		if id1 != id2 { // 桥
			tree[id1] = append(tree[id1], id2)
			tree[id2] = append(tree[id2], id1)
		}
	}
	fmt.Println(tree)
}

type Edge = struct{ from, to, cost, index int }
type TwoEdgeConnectedComponents struct {
	Tree    [][]Edge // 缩点后各个顶点形成的树
	CompId  []int    // 每个点所属的边双连通分量的编号
	Group   [][]int  // 每个边双连通分量里的点
	g       [][]Edge
	lowLink *LowLink
	k       int
}

func NewTwoEdgeConnectedComponents(g [][]Edge) *TwoEdgeConnectedComponents {
	return &TwoEdgeConnectedComponents{
		g:       g,
		lowLink: NewLowLink(g),
	}
}

func (tec *TwoEdgeConnectedComponents) Build() {
	tec.lowLink.Build()
	tec.CompId = make([]int, len(tec.g))
	for i := 0; i < len(tec.g); i++ {
		tec.CompId[i] = -1
	}
	for i := 0; i < len(tec.g); i++ {
		if tec.CompId[i] == -1 {
			tec.dfs(i, -1)
		}
	}
	tec.Group = make([][]int, tec.k)
	for i := 0; i < len(tec.g); i++ {
		tec.Group[tec.CompId[i]] = append(tec.Group[tec.CompId[i]], i)
	}
	tec.Tree = make([][]Edge, tec.k)
	for _, e := range tec.lowLink.Bridge {
		tec.Tree[tec.CompId[e.from]] = append(tec.Tree[tec.CompId[e.from]], Edge{tec.CompId[e.from], tec.CompId[e.to], e.cost, e.index})
		tec.Tree[tec.CompId[e.to]] = append(tec.Tree[tec.CompId[e.to]], Edge{tec.CompId[e.to], tec.CompId[e.from], e.cost, e.index})
	}
}

// 每个点所属的边双连通分量的编号.
func (tec *TwoEdgeConnectedComponents) Get(k int) int { return tec.CompId[k] }

func (tec *TwoEdgeConnectedComponents) dfs(idx, par int) {
	if par >= 0 && tec.lowLink.ord[par] >= tec.lowLink.low[idx] {
		tec.CompId[idx] = tec.CompId[par]
	} else {
		tec.CompId[idx] = tec.k
		tec.k++
	}
	for _, e := range tec.g[idx] {
		if tec.CompId[e.to] == -1 {
			tec.dfs(e.to, idx)
		}
	}
}

type LowLink struct {
	Articulation []int  // 関節点
	Bridge       []Edge // 橋
	g            [][]Edge
	ord, low     []int
	used         []bool
}

func NewLowLink(g [][]Edge) *LowLink {
	return &LowLink{g: g}
}

func (ll *LowLink) Build() {
	ll.used = make([]bool, len(ll.g))
	ll.ord = make([]int, len(ll.g))
	ll.low = make([]int, len(ll.g))
	k := 0
	for i := 0; i < len(ll.g); i++ {
		if !ll.used[i] {
			k = ll.dfs(i, k, -1)
		}
	}
}

func (ll *LowLink) dfs(idx, k, par int) int {
	ll.used[idx] = true
	ll.ord[idx] = k
	k++
	ll.low[idx] = ll.ord[idx]
	isArticulation := false
	beet := false
	cnt := 0
	for _, e := range ll.g[idx] {
		if e.to == par {
			tmp := beet
			beet = true
			if !tmp {
				continue
			}
		}
		if !ll.used[e.to] {
			cnt++
			k = ll.dfs(e.to, k, idx)
			ll.low[idx] = min(ll.low[idx], ll.low[e.to])
			if par >= 0 && ll.low[e.to] >= ll.ord[idx] {
				isArticulation = true
			}
			if ll.ord[idx] < ll.low[e.to] {
				ll.Bridge = append(ll.Bridge, e)
			}
		} else {
			ll.low[idx] = min(ll.low[idx], ll.ord[e.to])
		}
	}

	if par == -1 && cnt > 1 {
		isArticulation = true
	}
	if isArticulation {
		ll.Articulation = append(ll.Articulation, idx)
	}
	return k
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
