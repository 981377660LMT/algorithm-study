// https://ei1333.github.io/library/graph/tree/offline-lca.hpp
// O(n*α(n))，LCA离线，一般用于树上莫队

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/lca
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	parents := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &parents[i])
	}

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		tree[parents[i]] = append(tree[parents[i]], i+1)
		tree[i+1] = append(tree[i+1], parents[i])
	}

	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		queries[i] = [2]int{u, v}
	}

	res := OfflineLCA(tree, queries, 0)
	for i := 0; i < q; i++ {
		fmt.Fprintln(out, res[i])
	}
}

// LCA离线.
func OfflineLCA(tree [][]int, queries [][2]int, root int) []int {
	n := len(tree)
	ufa := NewUnionFindArray(n)
	st, mark, ptr, res := make([]int, n), make([]int, n), make([]int, n), make([]int, len(queries))
	for i := 0; i < len(queries); i++ {
		res[i] = -1
	}
	top := 0
	st[top] = root
	for _, q := range queries {
		mark[q[0]]++
		mark[q[1]]++
	}
	q := make([][][2]int, n)
	for i := 0; i < n; i++ {
		q[i] = make([][2]int, 0, mark[i])
		mark[i] = -1
		ptr[i] = len(tree[i])
	}
	for i := range queries {
		u, v := queries[i][0], queries[i][1]
		q[u] = append(q[u], [2]int{v, i})
		q[v] = append(q[v], [2]int{u, i})
	}
	run := func(u int) bool {
		for ptr[u] != 0 {
			v := tree[u][ptr[u]-1]
			ptr[u]--
			if mark[v] == -1 {
				top++
				st[top] = v
				return true
			}
		}
		return false
	}

	for top != -1 {
		u := st[top]
		if mark[u] == -1 {
			mark[u] = u
		} else {
			ufa.Union(u, tree[u][ptr[u]])
			mark[ufa.Find(u)] = u
		}

		if !run(u) {
			for _, v := range q[u] {
				if mark[v[0]] != -1 && res[v[1]] == -1 {
					res[v[1]] = mark[ufa.Find(v[0])]
				}
			}
			top--
		}
	}

	return res
}

type UnionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}
