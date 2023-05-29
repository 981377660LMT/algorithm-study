// https://ei1333.github.io/library/graph/tree/rmq-lowest-common-ancestor.hpp
// オイラーツアーとスパーステーブルによって最小共通祖先を求める.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var q int
	fmt.Fscan(in, &q)
	lca := NewLCA(n)
	for i := 1; i < n; i++ {
		var parent int
		fmt.Fscan(in, &parent)
		lca.AddEdge(i, parent)
	}
	lca.Build(0)

	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		fmt.Fprintln(out, lca.LCA(u, v))
	}
}

type LCARMQ struct {
	Depth     []int
	order, in []int
	g         [][]int
	st        *_St
}

func NewLCA(n int) *LCARMQ { return &LCARMQ{g: make([][]int, n)} }

func (lca *LCARMQ) AddEdge(u, v int) {
	lca.g[u] = append(lca.g[u], v)
	lca.g[v] = append(lca.g[v], u)
}

func (lca *LCARMQ) Build(root int) {
	n := len(lca.g)
	lca.order = make([]int, 0, n*2-1)
	lca.Depth = make([]int, 0, n*2-1)
	lca.in = make([]int, n)
	lca.dfs(root, -1, 0)
	vs := make([]int, 2*n-1)
	for i := range vs {
		vs[i] = i
	}
	lca.st = NewSparseTable(vs, func(i, j int) int {
		if lca.Depth[i] < lca.Depth[j] {
			return i
		}
		return j
	})
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
