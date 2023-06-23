// https://yukicoder.me/problems/no/1435
// 三角形子数组的个数
// 给定一个数组,求有多少个长度>=2的子数组满足
// !子数组升序排序后,前两个数之和大于等于最后一个数
// n<=2e5 1<=nums[i]<=1e9

// !子数组最大/最小值的单调性(窗口扩张时,最大值不减,最小值不增)
// !左端点固定时,可以二分右端点的位置
// !线段树维护区间最小、第二小、最大值

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	fmt.Fprintln(out, countTriangleSubArray(nums))
}

func countTriangleSubArray(nums []int) int {
	n := len(nums)
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		leaves[i] = E{nums[i], nums[i], INF}
	}

	seg := NewSegmentTree(leaves)
	res := 0
	for left := 0; left < n; left++ {
		right := seg.MaxRight(left, func(e E) bool { return e.max1 <= e.min1+e.min2 })
		res += right - left - 1
	}
	return res
}

const INF int = 1e18

type E = struct{ max1, min1, min2 int } // max/min1/min2

func (*SegmentTree) e() E { return E{-INF, -INF, -INF} }
func (*SegmentTree) op(a, b E) E {
	aMax1, aMin1, aMin2 := a.max1, a.min1, a.min2
	bMax1, bMin1, bMin2 := b.max1, b.min1, b.min2
	if aMax1 == -INF {
		return b
	}
	if bMax1 == -INF {
		return a
	}
	if aMin1 < bMin1 {
		return E{max(aMax1, bMax1), aMin1, min(aMin2, bMin1)}
	}
	return E{max(aMax1, bMax1), bMin1, min(aMin1, bMin2)}
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
