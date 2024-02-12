// https://ei1333.github.io/library/graph/tree/rmq-lowest-common-ancestor.hpp
// オイラーツアーとスパーステーブルによって最小共通祖先を求める.
// O(1) 求 LCA，O(nlogn) 预处理

package main

import (
	"math/bits"
)

type LCARMQ struct {
	Depth     []int
	order, in []int
	g         [][]int
	st        *_St
}

func NewLCA(tree [][]int, roots []int) *LCARMQ {
	res := &LCARMQ{g: tree}
	n := len(res.g)
	res.order = make([]int, 0, n*2-1)
	res.Depth = make([]int, 0, n*2-1)
	res.in = make([]int, n)
	for _, root := range roots {
		res.dfs(root, -1, 0)
	}
	vs := make([]int, 2*n-1)
	for i := range vs {
		vs[i] = i
	}
	res.st = NewSparseTable(vs, func(i, j int) int {
		if res.Depth[i] < res.Depth[j] {
			return i
		}
		return j
	})
	return res
}

func (lca *LCARMQ) LCA(u, v int) int {
	if u == v {
		return u
	}
	if lca.in[u] > lca.in[v] {
		u, v = v, u
	}
	return lca.order[lca.st.Query(lca.in[u], lca.in[v]+1)]
}

func (lca *LCARMQ) Dist(u, v int) int {
	return lca.Depth[lca.in[u]] + lca.Depth[lca.in[v]] - 2*lca.Depth[lca.in[lca.LCA(u, v)]]
}

func (lca *LCARMQ) dfs(cur, pre, dep int) {
	lca.in[cur] = len(lca.order)
	lca.order = append(lca.order, cur)
	lca.Depth = append(lca.Depth, dep)
	for _, next := range lca.g[cur] {
		if next != pre {
			lca.dfs(next, cur, dep+1)
			lca.order = append(lca.order, cur)
			lca.Depth = append(lca.Depth, dep)
		}
	}
}

type S = int

type _St struct {
	st     [][]S
	lookup []int
	op     func(S, S) S
}

func NewSparseTable(nums []S, op func(S, S) S) *_St {
	res := &_St{}
	n := len(nums)
	b := bits.Len(uint(n))
	st := make([][]S, b)
	for i := range st {
		st[i] = make([]S, n)
	}
	for i := range nums {
		st[0][i] = nums[i]
	}
	for i := 1; i < b; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	lookup := make([]int, n+1)
	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}
	res.st = st
	res.lookup = lookup
	res.op = op
	return res
}

func (st *_St) Query(start, end int) S {
	b := st.lookup[end-start]
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
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
