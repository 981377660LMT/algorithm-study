// https://maspypy.github.io/library/ds/offline_query/range_mex_query.hpp
// RangeMexQuery-离线查询区间mex

package main

import "fmt"

func main() {
	rmq := NewRangeMexQuery([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	rmq.AddQuery(1, 2)
	fmt.Println(rmq.Run(1))
}

const INF int = 1e18

type E = int

func e() E        { return INF }
func op(a, b E) E { return min(a, b) }
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

type RangeMexQuery struct {
	nums  []int
	query [][2]int
}

func NewRangeMexQuery(nums []int) *RangeMexQuery {
	return &RangeMexQuery{nums: nums}
}

// [start, end)
//  0 <= start <= end <= n
func (rmq *RangeMexQuery) AddQuery(start, end int) {
	rmq.query = append(rmq.query, [2]int{start, end})
}

// mexStart: mex的起始值(从0开始还是从1开始)
func (rmq *RangeMexQuery) Run(mexStart int) []int {
	n := len(rmq.nums)
	leaves := make([]E, n+2)
	for i := 0; i < n+2; i++ {
		leaves[i] = -1
	}
	seg := NewSegmentTree(leaves)

	q := len(rmq.query)
	res := make([]int, q)
	ids := make([][]int, n+1)
	for i := 0; i < q; i++ {
		end := rmq.query[i][1]
		ids[end] = append(ids[end], i)
	}

	for i := 0; i < n+1; i++ {
		for _, q := range ids[i] {
			start := rmq.query[q][0]
			mex := seg.MaxRight(mexStart, func(x int) bool { return x >= start })
			res[q] = mex
		}
		if i < n && rmq.nums[i] < n+2 {
			seg.Set(rmq.nums[i], i)
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

func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return e()
	}
	return st.seg[index+st.size]
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

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if predicate(op(res, st.seg[left])) {
					res = op(res, st.seg[left])
					left++
				}
			}
			return left - st.size
		}
		res = op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}
