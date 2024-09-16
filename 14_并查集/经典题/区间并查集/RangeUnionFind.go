// 区间并查集 RangeUnionFind/UnionFindRange

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	abc364_f()
}

// F - Range Connect MST (abc364 F) 区间并查集
// https://atcoder.jp/contests/abc364/tasks/abc364_f
// 给定n+q个点以及q次操作.
// 第i次操作，连接n+i与区间[Li,Ri]内的所有点，边权为Wi.
// 问最后是否连通，联通请求出最小生成树权值.
//
// 按照边权排序，区间并查集合并即可.
func abc364_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int
	fmt.Fscan(in, &N, &Q)
	operations := make([][3]int, Q)
	for i := 0; i < Q; i++ {
		fmt.Fscan(in, &operations[i][0], &operations[i][1], &operations[i][2])
		operations[i][0]--
		operations[i][1]--
	}
	order := make([]int, Q)
	for i := 0; i < Q; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		i, j = order[i], order[j]
		return operations[i][2] < operations[j][2]
	})

	res := 0
	uf := NewUnionFindRange(N + Q)
	for _, qi := range order {
		vNode, left, right, weight := N+qi, operations[qi][0], operations[qi][1], operations[qi][2]
		uf.Union(vNode, left, func(_, _ int) { res += weight })
		uf.UnionRange(left, right, func(_, _ int) { res += weight })
	}
	if uf.Part == 1 {
		fmt.Fprintln(out, res)
	} else {
		fmt.Fprintln(out, -1)
	}
}

// https://leetcode.cn/problems/amount-of-new-area-painted-each-day/
func amountPainted(paint [][]int) []int {
	uf := NewUnionFindRange(1e5 + 10)
	res := make([]int, len(paint))
	for i, v := range paint {
		left, right := v[0], v[1]
		res[i] = uf.UnionRange(left, right, nil)
	}
	return res
}

type UnionFindRange struct {
	Part   int
	n      int
	parent []int
	rank   []int
}

func NewUnionFindRange(n int) *UnionFindRange {
	uf := &UnionFindRange{
		Part:   n,
		n:      n,
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 1
	}
	return uf
}

func (uf *UnionFindRange) Find(x int) int {
	for x != uf.parent[x] {
		uf.parent[x] = uf.parent[uf.parent[x]]
		x = uf.parent[x]
	}
	return x
}

// Union 后, 大的编号的组会指向小的编号的组.
func (uf *UnionFindRange) Union(x, y int, beforeUnion func(big, small int)) bool {
	if x < y {
		x, y = y, x
	}
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return false
	}
	if beforeUnion != nil {
		beforeUnion(rootY, rootX)
	}
	uf.parent[rootX] = rootY
	uf.rank[rootY] += uf.rank[rootX]
	uf.Part--
	return true
}

// UnionRange 合并`闭区间`[left, right] 的所有元素, 返回合并次数.
func (uf *UnionFindRange) UnionRange(left, right int, beforeUnion func(big, small int)) int {
	if left >= right {
		return 0
	}
	leftRoot := uf.Find(left)
	rightRoot := uf.Find(right)
	unionCount := 0
	for rightRoot != leftRoot {
		unionCount++
		uf.Union(rightRoot, rightRoot-1, beforeUnion)
		rightRoot = uf.Find(rightRoot - 1)
	}
	return unionCount
}

func (uf *UnionFindRange) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFindRange) GetSize(x int) int {
	return uf.rank[uf.Find(x)]
}

func (uf *UnionFindRange) GetGroups() map[int][]int {
	group := make(map[int][]int)
	for i := 0; i < uf.n; i++ {
		group[uf.Find(i)] = append(group[uf.Find(i)], i)
	}
	return group
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
