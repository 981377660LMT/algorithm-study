// 检查边长度限制的路径是否存在
// https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths/solution/zai-xian-zuo-fa-shu-shang-bei-zeng-lca-b-lzjq/

// 1.求出森林中的多个最小生成树是最优的
// 2.求出多个最小生成树后，再求出两个点路径上的最大边权
package main

import (
	"math/bits"
	"sort"
)

const INF int = 1e18

func distanceLimitedPathsExist(n int, edgeList [][]int, queries [][]int) []bool {
	forest, _ := Kruskal(n, edgeList)
	uf := NewUnionFindArray(n)
	for _, e := range edgeList {
		uf.Union(e[0], e[1])
	}
	lca := NewLCA(n, forest)
	res := make([]bool, len(queries))
	for i, q := range queries {
		if uf.IsConnected(q[0], q[1]) {
			res[i] = lca.QueryMaxWeight(q[0], q[1]) < q[2]
		}
	}
	return res
}

type AdjListEdge struct{ to, weight int }

type LCA struct {
	n       int
	bitLen  int
	tree    [][]AdjListEdge
	depth   []int
	visited []bool

	// dp[i][j] 表示节点j向上跳2^i步的父节点,dpWeight[i][j] 表示节点j向上跳2^i步经过的最大边权
	dp, dpWeight [][]int
}

func NewLCA(n int, tree [][]AdjListEdge) *LCA {
	lca := &LCA{
		n:       n,
		bitLen:  bits.Len(uint(n)),
		tree:    tree,
		depth:   make([]int, n),
		visited: make([]bool, n),
	}

	lca.dp, lca.dpWeight = makeDp(lca)
	for i := 0; i < n; i++ {
		if !lca.visited[i] {
			lca.dfsAndInitDp(i, -1, 0)
		}
	}
	lca.fillDp()
	return lca
}

// 查询树节点两点路径上最大边权
func (lca *LCA) QueryMaxWeight(root1, root2 int) int {
	res := -INF
	if lca.depth[root1] < lca.depth[root2] {
		root1, root2 = root2, root1
	}

	toDepth := lca.depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.depth[root1]-toDepth)&(1<<i) > 0 {
			res = max(res, lca.dpWeight[i][root1])
			root1 = lca.dp[i][root1]
		}
	}

	if root1 == root2 {
		return res
	}

	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			res = max(res, max(lca.dpWeight[i][root1], lca.dpWeight[i][root2]))
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}

	res = max(res, max(lca.dpWeight[0][root1], lca.dpWeight[0][root2]))
	return res
}

func (lca *LCA) dfsAndInitDp(cur, pre, dep int) {
	lca.visited[cur] = true
	lca.depth[cur] = dep
	lca.dp[0][cur] = pre
	for _, e := range lca.tree[cur] {
		if next := e.to; next != pre {
			lca.dpWeight[0][next] = e.weight
			lca.dfsAndInitDp(next, cur, dep+1)
		}
	}
}

func makeDp(lca *LCA) (dp, dpWeight [][]int) {
	dp, dpWeight = make([][]int, lca.bitLen), make([][]int, lca.bitLen)
	for i := 0; i < lca.bitLen; i++ {
		dp[i], dpWeight[i] = make([]int, lca.n), make([]int, lca.n)
		for j := 0; j < lca.n; j++ {
			dp[i][j] = -1
			dpWeight[i][j] = -INF
		}
	}
	return
}

func (lca *LCA) fillDp() {
	for i := 0; i < lca.bitLen-1; i++ {
		for j := 0; j < lca.n; j++ {
			if lca.dp[i][j] == -1 {
				lca.dp[i+1][j] = -1
			} else {
				lca.dp[i+1][j] = lca.dp[i][lca.dp[i][j]]
				lca.dpWeight[i+1][j] = max(lca.dpWeight[i][j], lca.dpWeight[i][lca.dp[i][j]])
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

type KruskalEdge struct {
	u, v   int
	weight int
}

// !给定无向图的边，求出一个最小生成树(如果不存在,则求出的是森林中的多个最小生成树)
func Kruskal(n int, edges [][]int) (tree [][]AdjListEdge, ok bool) {
	sortedEdges := make([]KruskalEdge, len(edges))
	for i := range edges {
		e := edges[i]
		sortedEdges[i] = KruskalEdge{u: e[0], v: e[1], weight: e[2]}
	}
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].weight < sortedEdges[j].weight
	})

	tree = make([][]AdjListEdge, n)
	uf := NewUnionFindArray(n)
	count := 0
	for i := range sortedEdges {
		edge := &sortedEdges[i]
		root1, root2 := uf.Find(edge.u), uf.Find(edge.v)
		if root1 != root2 {
			uf.Union(edge.u, edge.v)
			tree[edge.u] = append(tree[edge.u], AdjListEdge{to: edge.v, weight: edge.weight})
			tree[edge.v] = append(tree[edge.v], AdjListEdge{to: edge.u, weight: edge.weight})
			count++
			if count == n-1 {
				return tree, true
			}
		}
	}

	return tree, false
}

func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		size:   n,
		Rank:   rank,
		parent: parent,
	}
}

type UnionFindArray struct {
	size   int
	Part   int
	Rank   []int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.Rank[root1] > ufa.Rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.Rank[root2] += ufa.Rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}
