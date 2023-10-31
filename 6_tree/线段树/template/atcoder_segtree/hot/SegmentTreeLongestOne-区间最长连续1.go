// SegmentTreeLongestOne-区间最长连续1

package main

import (
	"fmt"
	"strings"
)

func main() {
	nums := []int{1, 0, 1, 1, 0, 1, 1, 1}
	leaves := make([]E, len(nums))
	for i, v := range nums {
		leaves[i] = FromElement(v)
	}
	tree := NewSegmentTreeLongestOne(leaves)

	fmt.Println(tree.QueryAll().pairCount)
	fmt.Println(tree.Query(0, 6).pairCount)

}

type V = int

type E = struct {
	size                       int
	preOne, sufOne, longestOne int // 前缀1的个数, 后缀1的个数, 最长1的个数
	leftV, rightV              V   // 左端点值, 右端点值

	// !TODO
	pairCount int // 区间内所有极长连续1段的贡献和 sum(len_i*(len_i+1)/2)
}

func FromElement(v V) E {
	if v == 1 {
		return E{
			size: 1, preOne: 1, sufOne: 1, longestOne: 1, leftV: 1, rightV: 1,
			// !TODO
			pairCount: 1,
		}
	}

	return E{
		size: 1, leftV: v, rightV: v,
		// !TODO
		pairCount: 0,
	}
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

		// !TODO
		n1, n2 := a.sufOne, b.preOne
		n3 := n1 + n2
		res.pairCount = a.pairCount + b.pairCount + n3*(n3+1)/2 - n1*(n1+1)/2 - n2*(n2+1)/2
	} else {
		res.preOne = a.preOne
		res.sufOne = b.sufOne
		res.longestOne = max(a.longestOne, b.longestOne)

		// !TODO
		res.pairCount = a.pairCount + b.pairCount
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

func (tree *SegmentTreeLongestOne) String() string {
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
