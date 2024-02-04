// Auxiliary Tree 虚树
// https://oi-wiki.org/graph/virtual-tree/
// https://cmwqf.github.io/2020/04/17/%E6%B5%85%E8%B0%88%E8%99%9A%E6%A0%91/
// https://tjkendev.github.io/procon-library/python/graph/auxiliary_tree.html
// 指定された頂点たちの最小共通祖先関係を保って木を圧縮してできる補助的な木
// !有的时候，题目给你一棵树，然后有q组询问，每组询问为与树上有关的k个点的相关问题，
// 而这个问题如果只有一个询问的话，可以用树形dp轻松解决
// 那么我们可以考虑每次只在这k个点及相关的点构成的树上进行dp
// 往往需要虚树上进行树形 DP
// !把一些有用的点给拿出来，然后通过最少的lca把这些节点给穿到一起，在新的树上做树形dp。
// 点集个数为k时,最多2*k-1个顶点

package main

import (
	"sort"
)

func main() {

}

type AuxiliaryTree struct {
	G             [][]int // 原图邻接表(无向边)
	vg            [][]int // 虚树邻接表(有向边)
	s             []int
	fs, ls, depth []int

	lg []int
	st [][]int
}

// 给定顶点个数n和无向边集(u,v)构建.
//
//	O(nlogn)
func NewAuxiliaryTree(n int, edges [][]int) *AuxiliaryTree {
	g := make([][]int, n)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	res := &AuxiliaryTree{
		G:     g,
		vg:    make([][]int, n),
		s:     []int{},
		fs:    make([]int, n),
		ls:    make([]int, n),
		depth: make([]int, n),
	}

	res.dfs(0, -1, 0)
	res.buildSt()
	return res
}

// 指定点集,返回虚树的(有向图邻接表,虚树的根).
//
//	如果虚树不存在`(len(points<=1))`,返回空邻接表和-1.
//	O(klogk) 构建虚树.
func (t *AuxiliaryTree) Query(points []int) ([][]int, int) {
	k := len(points)
	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool {
		return t.fs[points[i]] < t.fs[points[j]]
	})

	for i := 0; i < k-1; i++ {
		x, y := t.fs[points[i]], t.fs[points[i+1]]
		l := t.lg[y-x+1]
		w := t.st[l][x]
		if t.depth[t.st[l][y-(1<<l)+1]] < t.depth[t.st[l][x]] {
			w = t.st[l][y-(1<<l)+1]
		}
		points = append(points, w)
	}

	sort.Slice(points, func(i, j int) bool {
		return t.fs[points[i]] < t.fs[points[j]]
	})

	stk := []int{}
	pre := -1
	root := -1
	for _, v := range points {
		if pre == v {
			continue
		}
		for len(stk) > 0 && t.ls[stk[len(stk)-1]] < t.fs[v] {
			stk = stk[:len(stk)-1]
		}
		if len(stk) > 0 {
			parent := stk[len(stk)-1]
			t.vg[parent] = append(t.vg[parent], v)
			if root == -1 {
				root = parent
			}
		}

		t.vg[v] = t.vg[v][:0]
		stk = append(stk, v)
		pre = v
	}

	return t.vg, root
}

func (t *AuxiliaryTree) dfs(v, p, d int) {
	t.depth[v] = d
	t.fs[v] = len(t.s)
	t.s = append(t.s, v)
	for _, w := range t.G[v] {
		if w == p {
			continue
		}
		t.dfs(w, v, d+1)
		t.s = append(t.s, v)
	}
	t.ls[v] = len(t.s)
	t.s = append(t.s, v)
}

func (t *AuxiliaryTree) buildSt() {
	l := len(t.s)
	lg := make([]int, l+1)
	for i := 2; i <= l; i++ {
		lg[i] = lg[i>>1] + 1
	}
	st := make([][]int, lg[l]+1)
	for i := range st {
		st[i] = make([]int, l-(1<<i)+1)
		for j := range st[i] {
			st[i][j] = l
		}
	}

	copy(st[0], t.s)
	b := 1
	for i := 0; i < lg[l]; i++ {
		st0, st1 := st[i], st[i+1]
		for j := 0; j < l-(b<<1)+1; j++ {
			st1[j] = st0[j]
			if t.depth[st0[j+b]] < t.depth[st0[j]] {
				st1[j] = st0[j+b]
			}
		}
		b <<= 1
	}

	t.lg = lg
	t.st = st
}
