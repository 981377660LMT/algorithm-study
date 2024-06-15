// 3177. 求出最长好子序列 II (值域线段树优化dp)
// https://leetcode.cn/problems/find-the-maximum-length-of-a-good-subsequence-ii/solutions/2805181/chu-li-xiang-lin-yuan-su-bu-tong-de-dpwe-gobq/
// !求数组nums的最长子序列，满足相邻不等元素的对数<=k，
// !dp[i][v]表示i个相邻不同，结尾为v时的最长子序列长度
// 结尾相同时(v==u)：dp[i][v] = max(dp[i][u])+1
// 结尾不同时(v!=u)：dp[i][v] = max(dp[i-1][u])+1

package main

import "sort"

func maximumLength(nums []int, k int) int {
	nums32, size := SortedSet32(nums)
	dp := make([]*SegmentTree, k+1)
	for i := range dp {
		dp[i] = NewSegmentTree(size, func(int32) int32 { return 0 })
	}
	for _, v := range nums32 {
		for j := k; j >= 0; j-- {
			dp[j].Update(v, dp[j].Get(v)+1) // 结尾相同时
			if j > 0 {
				leftMax := dp[j-1].Query(0, v) // 结尾不同时
				rightMax := dp[j-1].Query(v+1, size)
				dp[j].Update(v, max32(dp[j].Get(v), max32(leftMax, rightMax)+1))
			}
		}
	}
	return int(dp[k].QueryAll())
}

func SortedSet32(nums []int) (newNums []int32, size int32) {
	tmp := append(nums[:0:0], nums...)
	sort.Ints(tmp)
	for fast := 0; fast < len(tmp); fast++ {
		if tmp[fast] != tmp[size] {
			size++
			tmp[size] = tmp[fast]
		}
	}
	size++
	tmp = tmp[:size]
	newNums = make([]int32, len(nums))
	for i, v := range nums {
		newNums[i] = int32(sort.SearchInts(tmp, v))
	}
	return
}

const INF32 int32 = 1 << 30

// PointSetRangeMin

type E = int32

func (*SegmentTree) e() E        { return 0 }
func (*SegmentTree) op(a, b E) E { return max32(a, b) }
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

type SegmentTree struct {
	n, size int32
	seg     []E
}

func NewSegmentTree(n int32, f func(int32) E) *SegmentTree {
	res := &SegmentTree{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
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
func (st *SegmentTree) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree) Update(index int32, value E) {
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
func (st *SegmentTree) Query(start, end int32) E {
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
