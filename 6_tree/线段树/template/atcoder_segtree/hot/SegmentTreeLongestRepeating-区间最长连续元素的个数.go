// SegmentTreeLongestRepeating-区间最长连续元素的个数

package main

import (
	"fmt"
	"strings"
)

// https://leetcode.cn/problems/longest-substring-of-one-repeating-character/
func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	n := len(s)
	leaves := make([]E, len(s))
	for i := 0; i < n; i++ {
		leaves[i] = FromElement(s[i])
	}
	seg := NewSegmentTreeLongestRepeating(leaves)
	res := make([]int, len(queryIndices))
	for i := 0; i < len(queryIndices); i++ {
		pos := queryIndices[i]
		char := queryCharacters[i]
		seg.Set(pos, FromElement(char))
		res[i] = seg.QueryAll().max
	}
	return res
}

type V = byte

type E = struct {
	size                int
	preMax, sufMax, max int // 前缀最大值(连续长度), 后缀最大值, 区间最大值
	leftV, rightV       V   // 左端点值, 右端点值
}

func FromElement(v V) E {
	return E{
		size:   1,
		preMax: 1,
		sufMax: 1,
		max:    1,
		leftV:  v,
		rightV: v,
	}
}

func (*SegmentTreeLongestRepeating) e() E { return E{} }
func (*SegmentTreeLongestRepeating) op(a, b E) E {
	res := E{leftV: a.leftV, rightV: b.rightV, size: a.size + b.size}
	if a.rightV == b.leftV {
		res.preMax = a.preMax
		if a.preMax == a.size {
			res.preMax += b.preMax
		}
		res.sufMax = b.sufMax
		if b.sufMax == b.size {
			res.sufMax += a.sufMax
		}
		res.max = max(max(a.max, b.max), a.sufMax+b.preMax)
	} else {
		res.preMax = a.preMax
		res.sufMax = b.sufMax
		res.max = max(a.max, b.max)
	}
	return res
}

type SegmentTreeLongestRepeating struct {
	n, size int
	seg     []E
}

func NewSegmentTreeLongestRepeating(leaves []E) *SegmentTreeLongestRepeating {
	res := &SegmentTreeLongestRepeating{}
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
func (st *SegmentTreeLongestRepeating) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeLongestRepeating) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeLongestRepeating) Update(index int, value E) {
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
func (st *SegmentTreeLongestRepeating) Query(start, end int) E {
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
func (st *SegmentTreeLongestRepeating) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeLongestRepeating) MaxRight(left int, predicate func(E) bool) int {
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
func (st *SegmentTreeLongestRepeating) MinLeft(right int, predicate func(E) bool) int {
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

func (tree *SegmentTreeLongestRepeating) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := 0; i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
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
