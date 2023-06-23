// 2213. 由单个字符重复的最长子字符串
// https://leetcode.cn/problems/longest-substring-of-one-repeating-character/
// ![l,r]区间的最大连续长度就是
// !左区间的最大连续长度,右区间最大连续长度,以及左右两区间结合在一起中间的最大连续长度.
package main

func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	n := len(s)
	leaves := make([]E, len(s))
	for i := 0; i < n; i++ {
		leaves[i] = E{1, 1, 1, 1, s[i], s[i]}
	}
	seg := NewSegmentTree(leaves)
	res := make([]int, len(queryIndices))
	for i := 0; i < len(queryIndices); i++ {
		pos := queryIndices[i]
		char := queryCharacters[i]
		seg.Set(pos, E{1, 1, 1, 1, char, char})
		res[i] = seg.QueryAll().max
	}
	return res
}

const INF int = 1e18

type E = struct {
	size                int
	preMax, sufMax, max int  // 前缀最大值,后缀最大值,区间最大值
	lc, rc              byte // 区间左端点字符,右端点字符
}

func (*SegmentTree) e() E { return E{} }
func (*SegmentTree) op(a, b E) E {
	res := E{lc: a.lc, rc: b.rc, size: a.size + b.size}
	if a.rc == b.lc {
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

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
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

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTree) MaxRight(left int, predicate func(E) bool) int {
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
				if predicate(st.op(res, st.seg[left])) {
					res = st.op(res, st.seg[left])
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
func (st *SegmentTree) MinLeft(right int, predicate func(E) bool) int {
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
				if predicate(st.op(st.seg[right], res)) {
					res = st.op(st.seg[right], res)
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
