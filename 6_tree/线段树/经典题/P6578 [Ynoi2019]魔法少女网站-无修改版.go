// https://www.luogu.com.cn/problem/P6578
// P6578 [Ynoi2019]魔法少女网站
//

package main

import (
	"sort"
)

// 给定一个数组nums和一些查询(start,end,x)，对每个查询回答区间[start,end)内有多少个子数组最大值不超过x.
// nums.length<=1e5,nums[i]<=1e5,查询个数<=1e5
// 离线查询+线段树.
// 将所有查询按照x从小到大排序,然后从小到大依次处理，即维护一个01数组.
// 长为len的极长连续段的贡献为len*(len+1)/2.
// 用线段树维护.
func 魔法少女网站无修改版本(nums []int, queries [][3]int) []int {
	type queryWithId struct{ start, end, x, id int }
	type pair struct{ value, index int }

	n, q := len(nums), len(queries)
	sortedQueries := make([]queryWithId, q)
	for qi, query := range queries {
		sortedQueries[qi] = queryWithId{start: query[0], end: query[1], x: query[2], id: qi}
	}
	sort.Slice(sortedQueries, func(i, j int) bool { return sortedQueries[i].x < sortedQueries[j].x })
	sortedPairs := make([]pair, n)
	for i, num := range nums {
		sortedPairs[i] = pair{value: num, index: i}
	}
	sort.Slice(sortedPairs, func(i, j int) bool { return sortedPairs[i].value < sortedPairs[j].value })
	res := make([]int, q)

	leaves := make([]E, n)
	for i := range leaves {
		leaves[i] = FromElement(0)
	}
	tree := NewSegmentTreeLongestOne(leaves)

	ptr := 0
	for _, query := range sortedQueries {
		qs, qe, qx, qid := query.start, query.end, query.x, query.id
		for ptr < n && sortedPairs[ptr].value <= qx {
			tree.Set(sortedPairs[ptr].index, FromElement(1))
			ptr++
		}
		res[qid] = tree.Query(qs, qe).pairCount
	}
	return res
}

type V = int

type E = struct {
	size                       int
	preOne, sufOne, longestOne int // 前缀1的个数, 后缀1的个数, 最长1的个数
	leftV, rightV              V   // 左端点值, 右端点值

	pairCount int // 区间内所有极长连续1段的贡献和 sum(len_i*(len_i+1)/2)
}

func FromElement(v V) E {
	if v == 1 {
		return E{size: 1, preOne: 1, sufOne: 1, longestOne: 1, leftV: 1, rightV: 1, pairCount: 1}
	}
	return E{size: 1, leftV: v, rightV: v}
}

func (*SegmentTreeLongestOne) e() E { return E{} }
func (*SegmentTreeLongestOne) op(a, b E) E {
	res := E{leftV: a.leftV, rightV: b.rightV, size: a.size + b.size}
	if a.rightV == b.leftV {
		res.preOne = a.preOne
		if a.preOne == a.size {
			res.preOne += b.preOne
		}
		res.sufOne = b.sufOne
		if b.sufOne == b.size {
			res.sufOne += a.sufOne
		}
		res.longestOne = max(max(a.longestOne, b.longestOne), a.sufOne+b.preOne)
	} else {
		res.preOne = a.preOne
		res.sufOne = b.sufOne
		res.longestOne = max(a.longestOne, b.longestOne)
	}
	return res
}

// 维护区间最长1的个数, 区间前缀1的个数，区间后缀1的个数.
type SegmentTreeLongestOne struct {
	n, size int
	seg     []E
}

func NewSegmentTreeLongestOne(leaves []E) *SegmentTreeLongestOne {
	res := &SegmentTreeLongestOne{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTreeLongestOne) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeLongestOne) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeLongestOne) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTreeLongestOne) Query(start, end int) E {
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
func (st *SegmentTreeLongestOne) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeLongestOne) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := st.op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeLongestOne) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := st.op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
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
