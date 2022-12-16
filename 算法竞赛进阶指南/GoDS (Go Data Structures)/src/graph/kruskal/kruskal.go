package kruskal

import (
	"cmnx/src/unionfind/unionfindarray"
	"fmt"
	"sort"
)

func demo() {
	n := 5
	edges := [][3]int{
		{0, 1, 1},
		{1, 2, 2},
		{2, 3, 1},
		{3, 4, 2},
		{0, 4, 3},
		{0, 2, 4},
	}

	sortedEdges := make([]KruskalEdge, len(edges))
	for i := range edges {
		edge := &edges[i]
		sortedEdges[i] = KruskalEdge{edge[0], edge[1], edge[2]}
	}

	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].weight < sortedEdges[j].weight
	})

	fmt.Println(Kruskal1(n, sortedEdges))
}

type KruskalEdge struct {
	u, v   int
	weight int
	// !某些题目需要
	// id int
}

// Kruskal1 算法 O(mlogm)
//  求出最小生成树的边权之和,如果不能构成生成树，返回 -1
func Kruskal1(n int, sortedEdges []KruskalEdge) int {
	uf := unionfindarray.NewUnionFindArray(n)
	count, res := 0, 0
	for i := range sortedEdges {
		edge := &sortedEdges[i]
		root1, root2 := uf.Find(edge.u), uf.Find(edge.v)
		if root1 != root2 {
			uf.Union(edge.u, edge.v)
			res += edge.weight
			count++
		}
	}

	if count != n-1 {
		return -1
	}

	return res
}

type AdjListEdge struct{ to, weight int }

// !给定无向图的边，求出一个最小生成树(如果不存在,则求出的是森林中的多个最小生成树)
func Kruskal2(n int, edges [][]int) (tree [][]AdjListEdge, ok bool) {
	sortedEdges := make([]KruskalEdge, len(edges))
	for i := range edges {
		e := edges[i]
		sortedEdges[i] = KruskalEdge{u: e[0], v: e[1], weight: e[2]}
	}
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].weight < sortedEdges[j].weight
	})

	tree = make([][]AdjListEdge, n)
	uf := unionfindarray.NewUnionFindArray(n)
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
