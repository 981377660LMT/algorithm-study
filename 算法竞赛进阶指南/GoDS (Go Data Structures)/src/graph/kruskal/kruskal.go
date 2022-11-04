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

	sortedEdges := make([]Edge, len(edges))
	for i := range edges {
		edge := &edges[i]
		sortedEdges[i] = Edge{edge[0], edge[1], edge[2], i}
	}

	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].weight < sortedEdges[j].weight
	})

	fmt.Println(Kruskal(n, sortedEdges))
}

type Edge struct {
	u, v   int
	weight int
	// 某些题目需要
	id int
}

// Kruskal 算法 O(mlogm)
//  如果不能构成生成树，返回 -1
func Kruskal(n int, sortedEdges []Edge) int {
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
