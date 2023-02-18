// https://ei1333.github.io/library/other/offline-rmq.hpp
// 离线RMQ (区间最小值查询)

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// https://judge.yosupo.jp/problem/staticrmq

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
	}
	minIndexes := OfflineRMQ(queries, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, nums[minIndexes[i]])
	}
}

// 每个查询是左闭右开区间[left, right).返回每个查询的最小值的下标.
//  0<=left<right<=n.
func OfflineRMQ(queries [][2]int, less func(i, j int) bool) []int {
	n := 0
	for _, query := range queries {
		n = max(n, query[1])
	}

	uf := NewUnionFindArray(n)
	st, mark, res := make([]int, n), make([]int, n), make([]int, len(queries))
	top := -1
	for _, query := range queries {
		mark[query[1]-1]++
	}
	q := make([][]int, n)
	for i := 0; i < n; i++ {
		q[i] = make([]int, 0, mark[i])
	}
	for i := 0; i < len(queries); i++ {
		q[queries[i][1]-1] = append(q[queries[i][1]-1], i)
	}
	for i := 0; i < n; i++ {
		for top >= 0 && !less(st[top], i) {
			uf.Union(st[top], i)
			top--
		}
		st[top+1] = i
		top++
		mark[uf.Find(i)] = i
		for _, j := range q[i] {
			res[j] = mark[uf.Find(queries[j][0])]
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
