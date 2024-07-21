// https://ei1333.github.io/library/other/offline-rmq.hpp
// 离线RMQ (区间最小值查询)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/staticrmq
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	queries := make([][2]int32, q)
	for i := int32(0); i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
	}
	minIndexes := OfflineRMQ(queries, func(i, j int32) bool {
		return nums[i] < nums[j]
	})

	for i := int32(0); i < q; i++ {
		fmt.Fprintln(out, nums[minIndexes[i]])
	}
}

// 每个查询是左闭右开区间[left, right).返回每个查询的最小值的下标.
//
//	0<=left<right<=n.
func OfflineRMQ(queries [][2]int32, less func(i, j int32) bool) []int32 {
	n := int32(0)
	for _, query := range queries {
		n = max32(n, query[1])
	}

	uf := make([]int32, n)
	for i := int32(0); i < n; i++ {
		uf[i] = -1
	}
	var find func(key int32) int32
	find = func(key int32) int32 {
		if uf[key] < 0 {
			return key
		}
		uf[key] = find(uf[key])
		return uf[key]
	}
	union := func(key1, key2 int32) {
		root1, root2 := find(key1), find(key2)
		if root1 == root2 {
			return
		}
		if uf[root1] > uf[root2] {
			root1, root2 = root2, root1
		}
		uf[root1] += uf[root2]
		uf[root2] = root1
	}

	st, mark, res := make([]int32, n), make([]int32, n), make([]int32, len(queries))
	top := -1
	for _, query := range queries {
		mark[query[1]-1]++
	}
	q := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		q[i] = make([]int32, 0, mark[i])
	}
	for i := int32(0); i < int32(len(queries)); i++ {
		right := queries[i][1] - 1
		q[right] = append(q[right], i)
	}
	for i := int32(0); i < n; i++ {
		for top >= 0 && !less(st[top], i) {
			union(st[top], i)
			top--
		}
		st[top+1] = i
		top++
		mark[find(i)] = i
		for _, j := range q[i] {
			res[j] = mark[find(queries[j][0])]
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
