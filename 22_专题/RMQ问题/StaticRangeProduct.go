// 分块，每一块长度为2^log，每个块维护前后缀，st表维护所有块.
// 预处理O(nlog(n)/2^log)，查询O(1)或O(2^log).

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int32 = 1e9 + 10

// https://judge.yosupo.jp/problem/staticrmq
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	type E = struct {
		min   int32
		index int32
	}
	e := func() E { return E{min: INF, index: -1} }
	op := func(s1, s2 E) E {
		if s1.min < s2.min {
			return s1
		}
		if s1.min > s2.min {
			return s2
		}
		return E{min: s1.min, index: min32(s1.index, s2.index)}
	}
	rmq := NewStaticRangeProduct(
		n, func(i int32) E { return E{min: nums[i], index: int32(i)} }, e, op, 4,
		func(n int32, f func(i int32) E) IRMQ[E] { return NewDisjointSparseTable(n, f, e, op) },
	)

	for i := int32(0); i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		fmt.Fprintln(out, rmq.Query(start, end).min)
	}
}

type IRMQ[E any] interface {
	Query(start, end int32) E
}

type StaticRangeProduct[E any] struct {
	n, log        int32
	arr, pre, suf []E // inclusive
	e             func() E
	op            func(a, b E) E
	rmq           IRMQ[E]
}

// log: 一般为4.
func NewStaticRangeProduct[E any](
	n int32, f func(i int32) E, e func() E, op func(a, b E) E, log int32,
	createRMQ func(n int32, f func(i int32) E) IRMQ[E],
) *StaticRangeProduct[E] {
	bNum := n >> log
	arr := make([]E, n)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i)
	}
	pre := append(arr[:0:0], arr...)
	suf := append(arr[:0:0], arr...)
	mask := int32((1 << log) - 1)
	for i := int32(1); i < n; i++ {
		if i&mask != 0 {
			pre[i] = op(pre[i-1], arr[i])
		}
	}
	for i := n - 1; i >= 1; i-- {
		if i&mask != 0 {
			suf[i-1] = op(arr[i-1], suf[i])
		}
	}
	rmq := createRMQ(bNum, func(i int32) E { return suf[i<<log] })
	return &StaticRangeProduct[E]{n: n, log: log, arr: arr, pre: pre, suf: suf, e: e, op: op, rmq: rmq}
}

func (s *StaticRangeProduct[E]) Query(start, end int32) E {
	if start >= end {
		return s.e()
	}
	end--
	a, b := start>>s.log, end>>s.log
	if a < b {
		x := s.rmq.Query(a+1, b)
		x = s.op(s.suf[start], x)
		x = s.op(x, s.pre[end])
		return x
	}
	x := s.arr[start]
	for i := start + 1; i <= end; i++ {
		x = s.op(x, s.arr[i])
	}
	return x
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (s *StaticRangeProduct[E]) MaxRight(left int32, check func(e E) bool) int32 {
	if left == s.n {
		return s.n
	}
	ok, ng := left, s.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(s.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (s *StaticRangeProduct[E]) MinLeft(right int32, check func(e E) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(s.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j) 的贡献值.
type SparseTable[E any] struct {
	st [][]E
	e  func() E
	op func(E, E) E
	n  int32
}

func NewSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *SparseTable[E] {
	res := &SparseTable[E]{}
	b := bits.Len32(uint32(n))
	st := make([][]E, b)
	for i := range st {
		st[i] = make([]E, n)
	}
	for i := int32(0); i < n; i++ {
		st[0][i] = f(i)
	}
	for i := 1; i < b; i++ {
		for j := int32(1); j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	res.st = st
	res.e = e
	res.op = op
	res.n = n
	return res
}

func NewSparseTableFrom[E any](leaves []E, e func() E, op func(E, E) E) *SparseTable[E] {
	return NewSparseTable(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	b := bits.Len32(uint32(end-start)) - 1 // log2
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
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
func NewDisjointSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *DisjointSparseTable[E] {
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
