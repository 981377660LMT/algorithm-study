
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

const INF int = 2e9

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

	rmq := NewSparseTableSqrt(nums, func() S { return INF }, min)
	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		fmt.Fprintln(out, rmq.Query(start, end))
	}
}

type S = int

type SparseTableSqrt struct {
	arr    []S
	e      func() S
	op     func(S, S) S
	belong func(index int) int
	st     *_SparseTable
	pre    []S
	suf    []S
}

func NewSparseTableSqrt(arr []S, e func() S, op func(S, S) S) *SparseTableSqrt {
	res := &SparseTableSqrt{}
	n := len(arr)
	blockSize := int(math.Sqrt(float64(n))) + 1
	belong := func(index int) int { return index / blockSize }
	blockStart := func(index int) int { return index * blockSize }
	blockEnd := func(index int) int { return min((index+1)*blockSize, n) }
	blockCount := 1 + (n / blockSize)

	blockRes := make([]S, blockCount)
	for i := range blockRes {
		blockRes[i] = e()
	}
	for i := 0; i < n; i++ {
		bid := belong(i)
		blockRes[bid] = op(blockRes[bid], arr[i])
	}
	st := _NewSparseTable(blockRes, e, op)

	pre := make([]S, n)
	for bid := 0; bid < blockCount; bid++ {
		res := e()
		for i := blockStart(bid); i < blockEnd(bid); i++ {
			res = op(res, arr[i])
			pre[i] = res
		}
	}

	suf := make([]S, n)
	for bid := 0; bid < blockCount; bid++ {
		res := e()
		for i := blockEnd(bid) - 1; i >= blockStart(bid); i-- {
			res = op(arr[i], res)
			suf[i] = res
		}
	}

	res.arr = arr
	res.e = e
	res.op = op
	res.belong = belong
	res.st = st
	res.pre = pre
	res.suf = suf
	return res
}

func (st *SparseTableSqrt) Query(start, end int) S {
	if start < 0 {
		start = 0
	}
	if end > len(st.arr) {
		end = len(st.arr)
	}
	if start >= end {
		return st.e()
	}

	bid1 := st.belong(start)
	bid2 := st.belong(end - 1)
	if bid1 == bid2 {
		res := st.e()
		for i := start; i < end; i++ {
			res = st.op(res, st.arr[i])
		}
		return res
	}

	res := st.suf[start]
	res = st.op(res, st.st.Query(bid1+1, bid2))
	res = st.op(res, st.pre[end-1])
	return res
}

type _SparseTable struct {
	st     [][]S
	lookup []int
	e      func() S
	op     func(S, S) S
	n      int
}

func _NewSparseTable(leaves []S, e func() S, op func(S, S) S) *_SparseTable {
	res := &_SparseTable{}
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
func (st *_SparseTable) Query(start, end int) S {
	if start >= end {
		return st.e()
	}
	b := st.lookup[end-start]
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
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
