// https://maspypy.github.io/library/graph/online_unionfind.hpp

// https://maspypy.github.io/library/test/yukicoder/1170_2.test.cpp
// 不预先给出图,而是通过函数来寻找下一个可以合并(到达)的点.

package main

import (
	"fmt"
	"strings"
)

// 在线并查集, 通过函数来寻找下一个可以合并(到达)的点.
//  setUsed(u) : 删除数据结构中的u, 用于标记u已经被使用了.
//  findUnused(u) : 返回u可以到达的下一个点, 如果没有返回-1.
func OnlineUnionFind(
	n int,
	setUsed func(u int), findUnused func(u int) (v int),
) (res *UnionFindArray, ok bool) {
	res = NewUnionFindArray(n)
	done := make([]bool, n)
	que := []int{}
	for v := 0; v < n; v++ {
		if done[v] {
			continue
		}

		que = append(que, v)
		done[v] = true
		setUsed(v)
		for len(que) > 0 {
			x := que[0]
			que = que[1:]
			setUsed(x)
			done[x] = true
			for {
				to := findUnused(x)
				if to == -1 {
					break
				}
				res.Union(v, to)
				que = append(que, to)
				done[to] = true
				setUsed(to)
			}
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
