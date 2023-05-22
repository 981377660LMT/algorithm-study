package kruskal

import (
	"cmnx/src/unionfind/unionfindarray"
	"fmt"
	"sort"
	"strings"
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

type AdjListEdge struct{ to, weight int }
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

// !给定无向图的边，求出一个最小生成树(如果不存在,则求出的是森林中的多个最小生成树)
//  这个树也叫Kruskal重构树.
func Kruskal(n int, edges [][]int) (treeEdges [][3]int, ok bool) {
	type edge struct {
		u, v   int
		weight int
	}

	sortedEdges := make([]edge, len(edges))
	for i := range edges {
		e := edges[i]
		sortedEdges[i] = edge{u: e[0], v: e[1], weight: e[2]}
	}
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].weight < sortedEdges[j].weight
	})

	uf := NewUnionFindArray(n)
	count := 0
	for i := range sortedEdges {
		edge := &sortedEdges[i]
		root1, root2 := uf.Find(edge.u), uf.Find(edge.v)
		if root1 != root2 {
			uf.Union(edge.u, edge.v)
			treeEdges = append(treeEdges, [3]int{edge.u, edge.v, edge.weight})
			count++
			if count == n-1 {
				ok = true
				return
			}
		}
	}
	return
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		Rank:   rank,
		size:   n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int
	// 每个连通分量的大小
	Rank []int

	size   int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
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

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.size; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
