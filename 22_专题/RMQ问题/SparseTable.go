// st表, 查询区间最大值以及对应的下标(多个最大值时取最小的下标).

package main

import (
	"fmt"
	"math/bits"
)

const INF int = 1e18

func main() {
	leaves := []S{{1, 0}, {2, 1}, {3, 2}, {4, 3}, {5, 4}, {6, 5}, {7, 6}, {8, 7}, {9, 8}, {10, 9}}
	e := func() S { return S{max: -INF, index: -1} }
	op := func(s1, s2 S) S {
		if s1.max > s2.max {
			return s1
		}
		if s1.max < s2.max {
			return s2
		}
		return S{max: s1.max, index: min(s1.index, s2.index)}
	}

	st := NewSparseTable(leaves, e, op)
	fmt.Println(st.Query(0, 9))
	fmt.Println(st.Query(0, 8))
	fmt.Println(st.MaxRight(0, func(s S) bool { return s.max < 5 }))
}

// RangeMaxWIthIndex

type S = struct{ max, index int }

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j) 的贡献值.
type SparseTable struct {
	st     [][]S
	lookup []int
	e      func() S
	op     func(S, S) S
	n      int
}

func NewSparseTable(leaves []S, e func() S, op func(S, S) S) *SparseTable {
	res := &SparseTable{}
	n := len(leaves)
	b := bits.Len(uint(n))
	st := make([][]S, b)
	for i := range st {
		st[i] = make([]S, n)
	}
	for i := range leaves {
		st[0][i] = leaves[i]
	}
	for i := 1; i < b; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	lookup := make([]int, n+1)
	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}
	res.st = st
	res.lookup = lookup
	res.e = e
	res.op = op
	res.n = n
	return res
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable) Query(start, end int) S {
	if start >= end {
		return st.e()
	}
	b := st.lookup[end-start]
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *SparseTable) MaxRight(left int, check func(e S) bool) int {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *SparseTable) MinLeft(right int, check func(e S) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
