// https://ei1333.github.io/library/graph/others/block-cut-tree.hpp
// https://ei1333.hateblo.jp/entry/2020/03/25/010057
// 二重連結成分分解とも. 二重頂点連結成分とは,
// 1個の頂点を取り除いても連結である部分グラフである.

// 関節点は, その頂点とそれを端点とする辺を削除したときの部分グラフが非連結になるような頂点を指す.
//  したがって, 関節点を列挙した後に頑張ると列挙できる.
// build():
//  二重頂点連結成分分解する. bc には各二重頂点連結成分に属する辺が格納される.
//  注意不考虑节点数为 1 的图

//  !如果原图中某个连通分量只有一个点，则需要具体情况具体分析，
//  !这里不将孤立点当成点双.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/biconnected_components
	// !这道题把孤立点也当作点双
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

	tec := NewBiConnectedComponents(g)
	tec.Build()
	visited := make([]bool, n)
	for _, p := range tec.BCC {
		for _, v := range p {
			visited[v.from] = true
			visited[v.to] = true
		}
	}

	isolate := 0
	for _, v := range visited {
		if !v {
			isolate++
		}
	}

	fmt.Fprintln(out, len(tec.BCC)+isolate)
	for _, p := range tec.BCC {
		points := map[int]struct{}{}
		for _, v := range p {
			points[v.from] = struct{}{}
			points[v.to] = struct{}{}
			visited[v.from] = true
			visited[v.to] = true
		}
		fmt.Fprint(out, len(points))
		for k := range points {
			fmt.Fprint(out, " ", k)
		}
		fmt.Fprintln(out)
	}

	for i, v := range visited {
		if !v {
			fmt.Fprintln(out, 1, i)
		}
	}
}

type Edge = struct{ from, to, cost, index int }
type BiConnectedComponents struct {
	BCC     [][]Edge // 每个边双连通分量中的边
	g       [][]Edge
	lowLink *LowLink
	used    []bool
	tmp     []Edge
}

func NewBiConnectedComponents(g [][]Edge) *BiConnectedComponents {
	return &BiConnectedComponents{
		g:       g,
		lowLink: NewLowLink(g),
	}
}

func (bcc *BiConnectedComponents) Build() {
	bcc.lowLink.Build()
	bcc.used = make([]bool, len(bcc.g))
	for i := 0; i < len(bcc.used); i++ {
		if !bcc.used[i] {
			bcc.dfs(i, -1)
		}
	}
}

func (bcc *BiConnectedComponents) dfs(idx, par int) {
	bcc.used[idx] = true
	beet := false
	for _, next := range bcc.g[idx] {
		if next.to == par {
			b := beet
			beet = true
			if !b {
				continue
			}
		}

		if !bcc.used[next.to] || bcc.lowLink.ord[next.to] < bcc.lowLink.ord[idx] {
			bcc.tmp = append(bcc.tmp, next)
		}

		if !bcc.used[next.to] {
			bcc.dfs(next.to, idx)
			if bcc.lowLink.low[next.to] >= bcc.lowLink.ord[idx] {
				bcc.BCC = append(bcc.BCC, []Edge{})
				for {
					e := bcc.tmp[len(bcc.tmp)-1]
					bcc.BCC[len(bcc.BCC)-1] = append(bcc.BCC[len(bcc.BCC)-1], e)
					bcc.tmp = bcc.tmp[:len(bcc.tmp)-1]
					if e.index == next.index {
						break
					}
				}
			}
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
