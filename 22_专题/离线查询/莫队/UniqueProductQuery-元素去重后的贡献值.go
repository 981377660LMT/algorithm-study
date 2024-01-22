// UniqueProductQuery
// https://maspypy.github.io/library/ds/offline_query/uniqueproductquery.hpp
// !离线计算区间[start,end)内所有元素去重后的贡献值(op)

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
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

	U := NewUniqueProductQuery(nums)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		U.AddQuery(l, r)
	}
	res := U.Run(func(i int) E { return 1 })
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
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

type UniqueProductQuery struct {
	n     int
	nums  []int
	query [][2]int
}

func NewUniqueProductQuery(nums []int) *UniqueProductQuery {
	return &UniqueProductQuery{
		n:    len(nums),
		nums: nums,
	}
}

// [start, end)
//
//	0 <= start <= end <= n
func (upq *UniqueProductQuery) AddQuery(start, end int) {
	upq.query = append(upq.query, [2]int{start, end})
}

// f: 位于下标为i处的元素对答案的贡献值
func (upq *UniqueProductQuery) Run(f func(i int) E) []E {
	q := len(upq.query)
	res := make([]E, q)
	ids := make([][]int, upq.n+1)
	for qi := 0; qi < q; qi++ {
		ids[upq.query[qi][1]] = append(ids[upq.query[qi][1]], qi)
	}
	for _, q := range ids[0] {
		res[q] = e()
	}
	seg := NewSegmentTree(make([]E, upq.n))
	pos := make(map[int]int, upq.n)
	for i := 0; i < upq.n; i++ {
		x := upq.nums[i]
		if p, ok := pos[x]; ok {
			seg.Set(p, e())
		}
		pos[x] = i
		seg.Set(i, f(i))
		for _, qi := range ids[i+1] {
			start, end := upq.query[qi][0], upq.query[qi][1]
			res[qi] = seg.Query(start, end)
		}
	}

	return res
}

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
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
		return e()
	}
	leftRes, rightRes := e(), e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}
