package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1054
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	uf := NewWeightedUnionFind(n)
	for i := 0; i < q; i++ {
		var op, a, b int
		fmt.Fscan(in, &op, &a, &b)
		if op == 1 {
			a, b = a-1, b-1
			uf.Union(a, b)
		} else if op == 2 {
			a -= 1
			uf.AddGroup(a, b)
		} else {
			a -= 1
			fmt.Fprintln(out, uf.Get(a))
		}
	}

}

// 维护分量和的并查集.
type WeightedUnionFind struct {
	Part   int
	parent []int
	value  []int
	delta  []int
	total  []int
}

// NewUnionFindWeighted
func NewWeightedUnionFind(n int) *WeightedUnionFind {
	uf := &WeightedUnionFind{
		Part:   n,
		parent: make([]int, n),
		value:  make([]int, n),
		delta:  make([]int, n),
		total:  make([]int, n),
	}
	for i := range uf.parent {
		uf.parent[i] = -1
	}
	return uf
}

// u的值加上delta.
func (uf *WeightedUnionFind) Add(u, delta int) {
	uf.value[u] += delta
	uf.total[uf.Find(u)] += delta
}

// u所在集合的值加上delta.
func (uf *WeightedUnionFind) AddGroup(u, delta int) {
	root := uf.Find(u)
	uf.delta[root] += delta
	uf.total[root] -= uf.parent[root] * delta
}

// u的值.
func (uf *WeightedUnionFind) Get(u int) int {
	_, delta := uf._find(u)
	return uf.value[u] + delta
}

// u所在集合的值.
func (uf *WeightedUnionFind) GetGroup(u int) int {
	return uf.total[uf.Find(u)]
}

func (uf *WeightedUnionFind) Union(u, v int) bool {
	u = uf.Find(u)
	v = uf.Find(v)
	if u == v {
		return false
	}
	if uf.parent[u] > uf.parent[v] {
		u ^= v
		v ^= u
		u ^= v
	}
	uf.parent[u] += uf.parent[v]
	uf.parent[v] = u
	uf.delta[v] -= uf.delta[u]
	uf.total[u] += uf.total[v]
	uf.Part--
	return true
}

func (uf *WeightedUnionFind) UnionWithCallback(u, v int, f func(u, v int)) bool {
	u = uf.Find(u)
	v = uf.Find(v)
	if u == v {
		return false
	}
	if uf.parent[u] > uf.parent[v] {
		u ^= v
		v ^= u
		u ^= v
	}
	uf.parent[u] += uf.parent[v]
	uf.parent[v] = u
	uf.delta[v] -= uf.delta[u]
	uf.total[u] += uf.total[v]
	uf.Part--
	if f != nil {
		f(u, v)
	}
	return true
}

// u的根.
func (uf *WeightedUnionFind) Find(u int) int {
	root, _ := uf._find(u)
	return root
}

// u和v是否连通.
func (uf *WeightedUnionFind) IsConnected(u, v int) bool {
	return uf.Find(u) == uf.Find(v)
}

// u所在集合的大小.
func (uf *WeightedUnionFind) GetSize(u int) int {
	return -uf.parent[uf.Find(u)]
}

// 所有分量.
func (uf *WeightedUnionFind) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := range uf.parent {
		root := uf.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (uf *WeightedUnionFind) _find(u int) (root, delta int) {
	if uf.parent[u] < 0 {
		return u, uf.delta[u]
	}
	root, delta = uf._find(uf.parent[u])
	delta += uf.delta[u]
	uf.parent[u] = root
	uf.delta[u] = delta - uf.delta[root]
	return
}
