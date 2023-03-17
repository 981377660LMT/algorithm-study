package main

import (
	"sort"
)

const INF int = 1e18

// https://leetcode.cn/problems/count-number-of-possible-root-nodes/
func rootCount(edges [][]int, guesses [][]int, k int) int {
	n := len(edges) + 1
	R := NewRerootingSubTree(n)
	for _, e := range edges {
		R.AddEdge2(e[0], e[1], Edge{from: e[0], to: e[1]}, Edge{from: e[1], to: e[0]})
	}
	set := make(map[[2]int]struct{})
	for _, g := range guesses {
		set[[2]int{g[0], g[1]}] = struct{}{}
	}

	dp := R.ReRooting(
		func(root int) E { return 0 },
		func(dp1, dp2 int) int {
			return dp1 + dp2
		},
		func(dp int, edge Edge) int {
			pair := [2]int{edge.from, edge.to}
			if _, ok := set[pair]; ok {
				return dp + 1
			}
			return dp
		},
	)

	res := 0
	for _, d := range dp {
		if d >= k {
			res++
		}
	}
	return res
}

type E = int
type Edge = struct{ from, to int }

type ReRootingSubTree struct {
	G           [][]Node
	ld, rd      [][]E
	lp, rp      []int
	e           func(root int) E
	op          func(dp1, dp2 E) E
	composition func(dp E, e Edge) E
}
type Node struct {
	to, rev int
	data    Edge
}

func NewRerootingSubTree(n int) *ReRootingSubTree {
	res := &ReRootingSubTree{
		G:  make([][]Node, n),
		ld: make([][]E, n),
		rd: make([][]E, n),
		lp: make([]int, n),
		rp: make([]int, n),
	}
	return res
}

func (rr *ReRootingSubTree) AddEdge(u, v int, e Edge) {
	rr.AddEdge2(u, v, e, e)
}

func (rr *ReRootingSubTree) AddEdge2(u, v int, e, revE Edge) {
	rr.G[u] = append(rr.G[u], Node{to: v, data: e})
	rr.G[v] = append(rr.G[v], Node{to: u, data: revE})
}

func (rr *ReRootingSubTree) ReRooting(
	e func(root int) E,
	op func(dp1, dp2 E) E,
	compositionEdge func(dp E, e Edge) E,
) []E {
	rr.e = e
	rr.op = op
	rr.composition = compositionEdge
	n := len(rr.G)
	for i := 0; i < n; i++ {
		sort.Slice(rr.G[i], func(j, k int) bool {
			return rr.G[i][j].to < rr.G[i][k].to
		})
		rr.ld[i] = make([]E, len(rr.G[i])+1)
		rr.rd[i] = make([]E, len(rr.G[i])+1)
		for j := range rr.ld[i] {
			rr.ld[i][j] = e(i)
			rr.rd[i][j] = e(i)
		}
		rr.lp[i] = 0
		rr.rp[i] = len(rr.G[i]) - 1
	}

	for i := 0; i < n; i++ {
		for j := range rr.G[i] {
			rr.G[i][j].rev = rr.search(rr.G[rr.G[i][j].to], i)
		}
	}

	res := make([]E, n)
	for i := 0; i < n; i++ {
		res[i] = rr.dfs(i, -1)
	}
	return res
}

// !root 作为根节点时, 子树 v 的 dp 值
func (rr *ReRootingSubTree) SubTree(root, v int) E {
	k := rr.search(rr.G[root], v)
	return rr.composition(rr.dfs(v, rr.G[root][k].rev), rr.G[root][k].data)
}

func (rr *ReRootingSubTree) dfs(root, eid int) E {
	for rr.lp[root] != eid && rr.lp[root] < len(rr.G[root]) {
		e := rr.G[root][rr.lp[root]]
		rr.ld[root][rr.lp[root]+1] = rr.op(rr.ld[root][rr.lp[root]], rr.composition(rr.dfs(e.to, e.rev), e.data))
		rr.lp[root]++
	}
	for rr.rp[root] != eid && rr.rp[root] >= 0 {
		e := rr.G[root][rr.rp[root]]
		rr.rd[root][rr.rp[root]] = rr.op(rr.rd[root][rr.rp[root]+1], rr.composition(rr.dfs(e.to, e.rev), e.data))
		rr.rp[root]--
	}
	if eid < 0 {
		return rr.rd[root][0]
	}
	return rr.op(rr.ld[root][eid], rr.rd[root][eid+1])
}

func (rr *ReRootingSubTree) search(vs []Node, idx int) int {
	return sort.Search(len(vs), func(i int) bool { return vs[i].to >= idx })
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
