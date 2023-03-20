// O(nlogn)

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

	R := NewRerootingSubTree(n)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		edges[i] = [2]int{a, b}
		R.AddEdge(a, b, Edge{from: a, to: b, cost: 1})
	}

	// !E: (每个点处的直径,每个点到最远点的距离)
	R.ReRooting(
		func(root int) E { return E{0, 0} },
		func(dp1, dp2 E) E {
			return E{max(max(dp1.diam, dp2.diam), dp1.dist+dp2.dist), max(dp1.dist, dp2.dist)}
		},
		func(dp E, edge Edge) E {
			return E{dp.diam, dp.dist + 1}
		},
	)

	res := n
	for _, e := range edges {
		u, v := e[0], e[1]
		dia1, dia2 := R.SubTree(u, v).diam, R.SubTree(v, u).diam
		res = min(res, max(max(dia1, dia2), (dia1+1)/2+(dia2+1)/2+1))
	}
	fmt.Fprintln(out, res)
}

type E = struct{ diam, dist int }
type Edge = struct{ from, to, cost int }

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
	rr.G[u] = append(rr.G[u], Node{to: v, data: e})
	rr.G[v] = append(rr.G[v], Node{to: u, data: e})
}

// 当边的方向会影响结果时, 需要给正向和反向分别添加不同的边
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

func (rr *ReRootingSubTree) search(vs []Node, vid int) int {
	return sort.Search(len(vs), func(i int) bool { return vs[i].to >= vid })
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
