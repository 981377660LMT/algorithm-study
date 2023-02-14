// https://ei1333.github.io/library/graph/others/namori-graph.hpp
// Namori Graph
// n頂点 n辺からなる`連結`無向グラフは, サイクルが 1個だけあるグラフとなる。
// このグラフを, とある漫画家のアイコンにちなんで なもりグラフ と呼ばれることがるが,
// 学術的には Unicyclic Graph, Pseudoforest が正しい。
// !ここでは, このグラフを 1 つのサイクル と, サイクル内の頂点に付属する木に分解する。
// !またサイクルに含まれる頂点番号を, サイクルの頂点数を kとして [0,k)にふりなおし,
// !これを tree_id と呼ぶことにする。
// !また付属する木も同様に, 木の頂点数をlとして [0,l)にふりなおす。

// Build():
//  サイクルと木に分解する。頂点数と辺の本数が同じ無向連結グラフである必要がある。
// Forest():
//  分解した無向木が treeId の昇順に格納される。木の頂点番号は0から振り直している。
//  辺の from, to は振り直し後の頂点番号, cost,idx はもとのグラフの辺の値をコピーする。
// LoopEdges():
//  サイクルに含まれる辺が順に格納される。i番目の辺は tree_id iと i+1を結ぶ辺である。
//  辺の from, to, cost, idx はもとのグラフの辺の値をコピーする。
// GetId(k):
//  頂点 kについて, サイクルの tree_id, 振り直された木の頂点番号 id を返す。
// GetInvId(tree_id, k):
//  サイクルの tree_id に付属する頂点 kの`もとの頂点番号`を返す。
//  特に GetInvId(tree_id, 0) はサイクルに含まれていたもとの頂点番号を指す。

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	g := make([][]Edge, n)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		g[a] = append(g[a], Edge{a, b, 1, i})
		g[b] = append(g[b], Edge{b, a, 1, i})
	}

	ng := NewNamoriGraph(g)
	ng.Build()
	res := []int{}
	for _, e := range ng.LoopEdges {
		res = append(res, e.index+1)
	}
	sort.Ints(res)
	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

type Edge = struct{ from, to, cost, index int }
type Graph = [][]Edge
type NamoriGraph struct {
	Forest     []Graph
	LoopEdges  []Edge
	g          Graph
	iv         [][]int
	markId, id []int
}

func NewNamoriGraph(g Graph) *NamoriGraph {
	return &NamoriGraph{g: g}
}

func (ng *NamoriGraph) Build() {
	n := len(ng.g)
	deg := make([]int, n)
	used := make([]bool, n)
	que := []int{}
	for i := 0; i < n; i++ {
		deg[i] = len(ng.g[i])
		if deg[i] == 1 {
			que = append(que, i)
			used[i] = true
		}
	}

	for len(que) > 0 {
		idx := que[0]
		que = que[1:]
		for _, e := range ng.g[idx] {
			if used[e.to] {
				continue
			}
			deg[e.to]--
			if deg[e.to] == 1 {
				que = append(que, e.to)
				used[e.to] = true
			}
		}
	}

	mx := 0
	for _, edges := range ng.g {
		for _, e := range edges {
			mx = max(mx, e.index)
		}
	}

	edgeUsed := make([]bool, mx+1)
	loop := []int{}
	for i := 0; i < n; i++ {
		if used[i] {
			continue
		}
		for update := true; update; {
			update = false
			loop = append(loop, i)
			for _, e := range ng.g[i] {
				if used[e.to] || edgeUsed[e.index] {
					continue
				}
				edgeUsed[e.index] = true
				ng.LoopEdges = append(ng.LoopEdges, Edge{i, e.to, e.cost, e.index})
				i = e.to
				update = true
				break
			}
		}
		break
	}

	loop = loop[:len(loop)-1]
	ng.markId = make([]int, n)
	ng.id = make([]int, n)
	for i := 0; i < len(loop); i++ {
		pre := loop[(i+len(loop)-1)%len(loop)]
		nxt := loop[(i+1)%len(loop)]
		sz := 0
		ng.markId[loop[i]] = i
		ng.iv = append(ng.iv, []int{})
		ng.id[loop[i]] = sz
		sz++
		ng.iv[len(ng.iv)-1] = append(ng.iv[len(ng.iv)-1], loop[i])
		for _, e := range ng.g[loop[i]] {
			if e.to != pre && e.to != nxt {
				ng.markDfs(e.to, loop[i], i, &sz)
			}
		}
		tree := make(Graph, sz)
		for _, e := range ng.g[loop[i]] {
			if e.to != pre && e.to != nxt {
				tree[ng.id[loop[i]]] = append(tree[ng.id[loop[i]]], Edge{ng.id[loop[i]], ng.id[e.to], e.cost, e.index})
				tree[ng.id[e.to]] = append(tree[ng.id[e.to]], Edge{ng.id[e.to], ng.id[loop[i]], e.cost, e.index})
				ng.buildDfs(e.to, loop[i], tree)
			}
		}
		ng.Forest = append(ng.Forest, tree)
	}
}

// 頂点 kについて, サイクルの tree_id, 振り直された木の頂点番号 id を返す。
func (ng *NamoriGraph) GetId(k int) (treeId, id int) {
	return ng.markId[k], ng.id[k]
}

// サイクルの tree_id に付属する頂点 kの`もとの頂点番号`を返す。
// 特に GetInvId(tree_id, 0) はサイクルに含まれていたもとの頂点番号を指す。
func (ng *NamoriGraph) GetInvId(treeId, k int) int {
	return ng.iv[treeId][k]
}

// markDfs
func (ng *NamoriGraph) markDfs(idx, par, k int, l *int) {
	ng.markId[idx] = k
	ng.id[idx] = *l
	*l++
	ng.iv[len(ng.iv)-1] = append(ng.iv[len(ng.iv)-1], idx)
	for _, e := range ng.g[idx] {
		if e.to != par {
			ng.markDfs(e.to, idx, k, l)
		}
	}
}

// buildDfs
func (ng *NamoriGraph) buildDfs(idx, par int, tree Graph) {
	for _, e := range ng.g[idx] {
		if e.to != par {
			tree[ng.id[idx]] = append(tree[ng.id[idx]], Edge{ng.id[idx], ng.id[e.to], e.cost, e.index})
			tree[ng.id[e.to]] = append(tree[ng.id[e.to]], Edge{ng.id[e.to], ng.id[idx], e.cost, e.index})
			ng.buildDfs(e.to, idx, tree)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
