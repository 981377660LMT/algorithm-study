package main

import (
	"fmt"
)

func main() {
	n := 5
	edges := [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}
	g := make([][]Edge, n)
	lca := NewLCA(n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		g[u] = append(g[u], Edge{v, 1})
		g[v] = append(g[v], Edge{u, 1})
		lca.AddEdge(u, v)
	}
	lca.Build(0)
	_, _, parents := CentroidDecompositionWithParents(g)

	fmt.Println(lca.LCA(0, 4)) // 原图上的lca
	fmt.Println(parents)       // 点分树上的父亲结点
}

type Edge = struct{ to, cost int }

// 树的重心分解, 返回点分树和点分树的根
//  g: 原图
//  centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//  root: 重心树的根
//	parents: 每个点的父亲,不存在则为-1
func CentroidDecompositionWithParents(g [][]Edge) (centTree [][]int, root int, parents []int) {
	n := len(g)
	subSize := make([]int, n)
	removed := make([]bool, n)
	centTree = make([][]int, n)
	parents = make([]int, n)
	for i := range parents {
		parents[i] = -1
	}

	var getSize func(cur, parent int) int
	var getCentroid func(cur, parent, mid int) int
	var build func(cur int) int
	getSize = func(cur, parent int) int {
		subSize[cur] = 1
		for _, e := range g[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			subSize[cur] += getSize(next, cur)
		}
		return subSize[cur]
	}
	getCentroid = func(cur, parent, mid int) int {
		for _, e := range g[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			if subSize[next] > mid {
				return getCentroid(next, cur, mid)
			}
		}
		return cur
	}
	build = func(cur int) int {
		centroid := getCentroid(cur, -1, getSize(cur, -1)/2)
		removed[centroid] = true
		for _, e := range g[centroid] {
			next := e.to
			if !removed[next] {
				u, v := centroid, build(next)
				centTree[u] = append(centTree[u], v)
				parents[v] = u
			}
		}
		removed[centroid] = false
		return centroid
	}

	root = build(0)
	return
}

type LCAHLD struct {
	Depth, Parent      []int
	Tree               [][]int
	dfn, top, heavySon []int
	dfnId              int
}

func NewLCA(n int) *LCAHLD {
	tree := make([][]int, n)
	dfn := make([]int, n)      // vertex => dfn
	top := make([]int, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	return &LCAHLD{
		Tree:     tree,
		dfn:      dfn,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
	}
}

// 添加无向边 u-v.
func (hld *LCAHLD) AddEdge(u, v int) {
	hld.Tree[u] = append(hld.Tree[u], v)
	hld.Tree[v] = append(hld.Tree[v], u)
}

func (hld *LCAHLD) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

func (hld *LCAHLD) LCA(u, v int) int {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.Parent[hld.top[v]]
	}
}

func (hld *LCAHLD) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

func (hld *LCAHLD) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, next := range hld.Tree[cur] {
		if next != pre {
			nextSize := hld.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAHLD) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld.markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld.markTop(next, next)
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
