// https://ei1333.github.io/library/graph/mst/boruvka.hpp
// Boruvka(最小全域木) 在线最小生成树
// 不预先给出图，
// 而是给定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。

package main

import (
	"fmt"
	"strings"
)

// Brouvka
// 陽にグラフを作らず、何らかのデータ構造で未訪問の行き先を探す想定。
//  !find_unused(u)：unused なうちで、u と最小コストで結べる点を探す。
func OnlineMST(
	n int,
	setUsed func(u int), setUnused func(u int), findUnused func(u int) (v int, cost int),
) (res [][3]int, ok bool) {
	uf := NewUnionFindArray(n)
	for {
		upd := false
		comp := make([][]int, n)
		cand := make([][3]int, n)
		for v := 0; v < n; v++ {
			cand[v] = [3]int{-1, -1, -1}
		}
		for v := 0; v < n; v++ {
			comp[uf.Find(v)] = append(comp[uf.Find(v)], v)
		}
		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			for i := range comp[v] {
				setUsed(comp[v][i])
			}
			for i := range comp[v] {
				x := comp[v][i]
				y, cost := findUnused(x)
				if y == -1 {
					continue
				}
				a, c := cand[v][0], cand[v][2]
				if a == -1 || cost < c {
					cand[v] = [3]int{x, y, cost}
				}
				for i := range comp[v] {
					setUnused(comp[v][i])
				}
			}
		}
		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			a, b, c := cand[v][0], cand[v][1], cand[v][2]
			if a == -1 {
				continue
			}
			upd = true
			if uf.Union(a, b) {
				res = append(res, [3]int{a, b, c})
			}
		}

		if !upd {
			break
		}
	}

	ok = true
	return
}

func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
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

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
