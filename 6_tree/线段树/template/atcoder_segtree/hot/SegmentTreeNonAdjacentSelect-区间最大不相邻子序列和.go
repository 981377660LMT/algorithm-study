// 打家劫舍带修改版
// SegmentTreeNonAdjacentSelect-区间最大不相邻子序列和(不相邻选数)
// 线性序列最大独立集/最大子段和/线段树维护最大独立集
// https://leetcode.cn/problems/maximum-sum-of-subsequence-with-non-adjacent-elements/solutions/2790603/fen-zhi-si-xiang-xian-duan-shu-pythonjav-xnhz/

package main

import (
	"bufio"
	"fmt"
	"os"
)

// P3097 [USACO13DEC] Optimal Milking G
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	seg := NewSegmentTree(int32(len(nums)), func(i int32) E { return FromElement(nums[i]) })
	res := 0
	for i := 0; i < q; i++ {
		var pos, v int
		fmt.Fscan(in, &pos, &v)
		seg.Set(int32(pos-1), FromElement(v))
		cur := seg.QueryAll()
		res += cur.f11
	}
	fmt.Fprintln(out, res)
}

// 100306. 不包含相邻元素的子序列的最大和
// https://leetcode.cn/problems/maximum-sum-of-subsequence-with-non-adjacent-elements/description/
func maximumSumSubsequence(nums []int, queries [][]int) int {
	const MOD int = 1e9 + 7

	seg := NewSegmentTree(int32(len(nums)), func(i int32) E { return FromElement(nums[i]) })
	res := 0
	for _, query := range queries {
		pos, v := query[0], query[1]
		seg.Set(int32(pos), FromElement(v))
		cur := seg.QueryAll()
		res += cur.f11 // !f11 没有任何限制，就是打家劫舍的答案
	}
	return res % MOD
}

const INF int = 1e18

// f00: 左右两个元素一定不选
// f01: 左边元素一定不选，右边元素可选可不选
// f10: 左边元素可选可不选，右边元素一定不选
// f11: 左右两个元素可选可不选，没有任何限制，就是打家劫舍的答案
type E = struct{ f00, f01, f10, f11 int }

func FromElement(x int) E {
	return E{f11: max(0, x)}
}
func (*SegmentTreeNonAdjacentSelect) e() E {
	return E{}
}
func (*SegmentTreeNonAdjacentSelect) op(a, b E) E {
	res := E{}
	res.f00 = max(a.f00+b.f10, a.f01+b.f00)
	res.f01 = max(a.f00+b.f11, a.f01+b.f01)
	res.f10 = max(a.f10+b.f10, a.f11+b.f00)
	res.f11 = max(a.f10+b.f11, a.f11+b.f01)
	return res
}
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

type SegmentTreeNonAdjacentSelect struct {
	n, size int32
	seg     []E
}

func NewSegmentTree(n int32, f func(int32) E) *SegmentTreeNonAdjacentSelect {
	res := &SegmentTreeNonAdjacentSelect{}
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
func (st *SegmentTreeNonAdjacentSelect) Get(index int32) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeNonAdjacentSelect) Set(index int32, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeNonAdjacentSelect) Update(index int32, value E) {
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
func (st *SegmentTreeNonAdjacentSelect) Query(start, end int32) E {
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
func (st *SegmentTreeNonAdjacentSelect) QueryAll() E { return st.seg[1] }
func (st *SegmentTreeNonAdjacentSelect) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
