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
//  不预先给出图，而是指定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。
//  findUnused(u)：返回 unused 中与 u 权值最小的点 v 和边权 cost
//                如果不存在，返回 (-1,*)
func OnlineMST(
	n int,
	setUsed func(u int), setUnused func(u int), findUnused func(u int) (v int, cost int),
) (res [][3]int, ok bool) {
	uf := NewUnionFindArray(n)
	for {
		updated := false
		groups := make([][]int, n)
		cand := make([][3]int, n) // [u, v, cost]
		for v := 0; v < n; v++ {
			cand[v] = [3]int{-1, -1, -1}
		}

		for v := 0; v < n; v++ {
			groups[uf.Find(v)] = append(groups[uf.Find(v)], v)
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			for _, x := range groups[v] {
				setUsed(x)
			}
			for _, x := range groups[v] {
				y, cost := findUnused(x)
				if y == -1 {
					continue
				}
				a, c := cand[v][0], cand[v][2]
				if a == -1 || cost < c {
					cand[v] = [3]int{x, y, cost}
				}
			}
			for _, x := range groups[v] {
				setUnused(x)
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
			updated = true
			if uf.Union(a, b) {
				res = append(res, [3]int{a, b, c})
			}
		}

		if !updated {
			break
		}
	}

	if len(res) != n-1 {
		return nil, false
	}
	return res, true
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
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
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
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
