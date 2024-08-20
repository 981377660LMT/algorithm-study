// 在线版猫树.
// https://codeforces.com/blog/entry/79108

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/staticrmq
	// 5e5 - 750ms
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	e := func() int32 { return INF32 }
	op := func(a, b int32) int32 { return min32(a, b) }
	st := NewDisjointSparseTableFast(int32(n), func(i int32) int32 { return nums[i] }, e, op)
	for i := 0; i < q; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)
		fmt.Fprintln(out, st.Query(l, r))
	}
}

const INF32 int32 = 1 << 30

// Static RMQ, O(n)预处理, O(1)查询.
type DisjointSparseTableFast[E any] struct {
	n        int32
	leaves   []E
	pre, suf []E
	st       *DisjointSparseTable[E]
	e        func() E
	op       func(E, E) E
}

func NewDisjointSparseTableFast[E comparable](n int32, f func(int32) E, e func() E, op func(E, E) E) *DisjointSparseTableFast[E] {
	res := &DisjointSparseTableFast[E]{}
	bNum := n >> 4
	leaves := make([]E, 0, n)
	for i := int32(0); i < n; i++ {
		leaves = append(leaves, f(i))
	}
	pre, suf := append(leaves[:0:0], leaves...), append(leaves[:0:0], leaves...)
	for i := int32(1); i < n; i++ {
		if i&15 != 0 {
			pre[i] = op(pre[i-1], leaves[i])
		}
	}
	for i := n - 1; i > 0; i-- {
		if i&15 != 0 {
			suf[i-1] = op(leaves[i-1], suf[i])
		}
	}
	st := NewDisjointSparse(bNum, func(i int32) E { return suf[i<<4] }, e, op)
	res.n = n
	res.leaves = leaves
	res.pre, res.suf = pre, suf
	res.st = st
	res.e = e
	res.op = op
	return res
}

func NewDisjointSparseTableFastFrom[E comparable](leaves []E, e func() E, op func(E, E) E) *DisjointSparseTableFast[E] {
	return NewDisjointSparseTableFast(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

func (st *DisjointSparseTableFast[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	end--
	a, b := start>>4, end>>4
	if a < b {
		x := st.st.Query(a+1, b)
		x = st.op(st.suf[start], x)
		x = st.op(x, st.pre[end])
		return x
	}
	x := st.leaves[start]
	for i := start + 1; i <= end; i++ {
		x = st.op(x, st.leaves[i])
	}
	return x
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTableFast[E]) MaxRight(left int32, check func(e E) bool) int32 {
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
func (ds *DisjointSparseTableFast[E]) MinLeft(right int32, check func(e E) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
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

type DisjointSparseTable[E any] struct {
	n    int32
	e    func() E
	op   func(E, E) E
	data [][]E
}

// DisjointSparseTable 支持幺半群的区间静态查询.
//
//	eg: 区间乘积取模/区间仿射变换...
func NewDisjointSparse[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *DisjointSparseTable[E] {
	res := &DisjointSparseTable[E]{}
	log := int32(1)
	for (1 << log) < n {
		log++
	}
	data := make([][]E, log)
	data[0] = make([]E, 0, n)
	for i := int32(0); i < n; i++ {
		data[0] = append(data[0], f(i))
	}
	for i := int32(1); i < log; i++ {
		data[i] = append(data[i], data[0]...)
		tmp := data[i]
		b := int32(1 << i)
		for m := b; m <= n; m += 2 * b {
			l, r := m-b, min32(m+b, n)
			for j := m - 1; j >= l+1; j-- {
				tmp[j-1] = op(tmp[j-1], tmp[j])
			}
			for j := m; j < r-1; j++ {
				tmp[j+1] = op(tmp[j], tmp[j+1])
			}
		}
	}
	res.n = n
	res.e = e
	res.op = op
	res.data = data
	return res
}

func (ds *DisjointSparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return ds.e()
	}
	end--
	if start == end {
		return ds.data[0][start]
	}
	lca := bits.Len32(uint32(start^end)) - 1
	return ds.op(ds.data[lca][start], ds.data[lca][end])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTable[E]) MaxRight(left int32, check func(e E) bool) int32 {
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
func (ds *DisjointSparseTable[E]) MinLeft(right int32, check func(e E) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
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

func topbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.Len32(uint32(x)) - 1)
}

func lowbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.TrailingZeros32(uint32(x)))
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
