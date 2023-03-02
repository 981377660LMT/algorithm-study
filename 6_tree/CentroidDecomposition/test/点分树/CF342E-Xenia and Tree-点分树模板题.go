package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &ops[i][0], &ops[i][1])
		ops[i][1]--
	}

	res := XeniaAndTree(n, edges, ops)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/CF342E
// 给定一棵 n 个节点的树，初始时 1 号节点为红色，其余为蓝色。
// 要求支持如下操作：
// 1.将一个节点变为红色。
// 2.询问节点 u 到最近红色节点的距离。
// !n,q<=1e5
//
// https://blog.csdn.net/qq_35755187/article/details/113445005
// 首先建出点分树，考虑对于每一个分治中心维护离它最近的红点。
// !根据性质，点分树上(u,v)两点的lca一定在原树中(u,v)的路径上
// !这就意味着原树中的路径都能以他们两个端点在点分树中的lca为终点，分解为两条路径。
// 用一个dist数组记录点分树中每个结点到它点分树子树中最近红色结点的距离。
// !查询某个点到最近红色结点的距离时，只要查询它在点分树中的祖先的dist+它和点分树中祖先在原树中的距离即可
func XeniaAndTree(n int, edges [][2]int, operations [][2]int) []int {
	const INF int = 1e18
	g := make([][]Edge, n)
	lca := NewLCA(n)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], Edge{e[1], 1})
		g[e[1]] = append(g[e[1]], Edge{e[0], 1})
		lca.AddEdge(e[0], e[1])
	}
	lca.Build(0)

	_, _, parents := CentroidDecompositionWithParents(g)

	res := []int{}
	dist := make([]int, n) // dist[v]表示点分树中点v到其子树中最近的红色结点的(原图上的)距离
	for i := range dist {
		dist[i] = INF
	}
	modify := func(u int) {
		cur := u
		for cur != -1 {
			dist[cur] = min(dist[cur], lca.Dist(u, cur))
			cur = parents[cur]
		}
	}
	query := func(u int) int {
		res, cur := INF, u
		for cur != -1 {
			res = min(res, dist[cur]+lca.Dist(u, cur))
			cur = parents[cur]
		}
		return res
	}

	modify(0) // !开始时1号节点为红色
	for _, op := range operations {
		if op[0] == 1 {
			modify(op[1])
		} else {
			res = append(res, query(op[1]))
		}
	}
	return res
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
