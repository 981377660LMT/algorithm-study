// 静态摩尔投票.查询区间绝对众数.

package main

import "sort"

// https://leetcode.cn/problems/online-majority-element-in-subarray/description/
type MajorityChecker struct {
	M *MajorityVoting
}

func Constructor(arr []int) MajorityChecker {
	return MajorityChecker{NewMajorityVoting(len(arr), func(i int) int { return arr[i] })}
}

func (this *MajorityChecker) Query(left int, right int, threshold int) int {
	res, freq := this.M.Query(left, right+1)
	if freq < threshold {
		return -1
	}
	return res
}

type MajorityVoting struct {
	seg *SegmentTree
	pos map[int][]int
}

func NewMajorityVoting(n int, f func(int) int) *MajorityVoting {
	seg := NewSegmentTree(n, func(i int) E { return of(f(i)) })
	pos := make(map[int][]int)
	for i := 0; i < n; i++ {
		v := f(i)
		pos[v] = append(pos[v], i)
	}
	return &MajorityVoting{seg: seg, pos: pos}
}

func (mv *MajorityVoting) Query(start, end int) (majority int, freq int) {
	majority = mv.seg.Query(start, end).value
	indexes := mv.pos[majority]
	freq = sort.SearchInts(indexes, end) - sort.SearchInts(indexes, start)
	return
}

const INF int = 1e18

type E = struct {
	value    int
	freqDiff int
}

func of(value int) E {
	return E{value: value, freqDiff: 1}
}

func (*SegmentTree) e() E { return E{} }
func (*SegmentTree) op(a, b E) E {
	if a.value == b.value {
		a.freqDiff += b.freqDiff
		return a
	}
	k := min(a.freqDiff, b.freqDiff)
	a.freqDiff -= k
	b.freqDiff -= k
	if a.freqDiff >= b.freqDiff {
		return a
	}
	return b
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

type SegmentTree struct {
	n, size int
	seg     []E
}

func NewSegmentTree(n int, f func(int) E) *SegmentTree {
	res := &SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
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
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int, value E) {
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
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
func (st *SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
