// UniqueProductQuery
// https://maspypy.github.io/library/ds/offline_query/uniqueproductquery.hpp
// !离线计算区间[start,end)内所有元素去重后的贡献值(op)
// 求区间出现次数偶数次数的异或和: https://codeforces.com/contest/703/submission/144458323

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	// yosupo()
	CF703D()
}

func yosupo() {
	// https://judge.yosupo.jp/problem/static_range_count_distinct
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	U := NewUniqueProductQuery(nums, func() E { return 0 }, func(a, b E) E { return a + b })
	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		U.AddQuery(start, end)
	}
	rangeUniqueCount := U.Run(func(i int) E { return 1 })

	for _, v := range rangeUniqueCount {
		fmt.Fprintln(out, v)
	}
}

// Mishka and Interesting sum
// https://www.luogu.com.cn/problem/CF703D
// 求区间出现次数偶数次数的元素的异或和
//
// 等价于 区间异或^去重后的异或
// !区间中出现偶数次的数的异或和就是0，此时再异或上区间内不同的数。此时这些出现偶数次的数的异或再次出现。而那些出现单次的数就消失了.
func CF703D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	preXor := make([]int, n+1)
	for i := 0; i < n; i++ {
		preXor[i+1] = preXor[i] ^ nums[i]
	}
	U := NewUniqueProductQuery(nums, func() E { return 0 }, func(a, b E) E { return a ^ b })

	var q int
	fmt.Fscan(in, &q)
	rangeXor := make([]int, q)
	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		start--
		U.AddQuery(start, end)
		rangeXor[i] = preXor[end] ^ preXor[start]
	}
	rangeUniqueXor := U.Run(func(i int) E { return nums[i] })
	res := make([]int, q)
	for i := 0; i < q; i++ {
		res[i] = rangeXor[i] ^ rangeUniqueXor[i]
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type E = int

type UniqueProductQuery struct {
	n     int
	nums  []int
	query [][2]int32
	e     func() E
	op    func(a, b E) E
}

func NewUniqueProductQuery(nums []int, e func() E, op func(a, b E) E) *UniqueProductQuery {
	return &UniqueProductQuery{
		n:    len(nums),
		nums: nums,
		e:    e,
		op:   op,
	}
}

// [start, end)
//
//	0 <= start <= end <= n
func (upq *UniqueProductQuery) AddQuery(start, end int) {
	upq.query = append(upq.query, [2]int32{int32(start), int32(end)})
}

// f: 位于下标为i处的元素对答案的贡献.
func (upq *UniqueProductQuery) Run(f func(i int) E) []E {
	q := int32(len(upq.query))
	res := make([]E, q)
	ids := make([][]int32, upq.n+1)
	for qi := int32(0); qi < q; qi++ {
		ids[upq.query[qi][1]] = append(ids[upq.query[qi][1]], qi)
	}
	for _, q := range ids[0] {
		res[q] = upq.e()
	}
	seg := NewSegmentTree(upq.n, func(i int) E { return upq.e() }, upq.e, upq.op)
	last := make(map[int]int, upq.n)
	for i, x := range upq.nums {
		if p, ok := last[x]; ok {
			seg.Set(p, upq.e())
		}
		last[x] = i
		seg.Set(i, f(i))
		for _, qi := range ids[i+1] {
			start, end := upq.query[qi][0], upq.query[qi][1]
			res[qi] = seg.Query(int(start), int(end))
		}
	}

	return res
}

type SegmentTree struct {
	n, size int
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func NewSegmentTree(n int, f func(int) E, e func() E, op func(a, b E) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	res.e = e
	res.op = op
	return res
}

func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
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
