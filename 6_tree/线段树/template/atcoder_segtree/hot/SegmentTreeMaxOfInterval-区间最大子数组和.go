// MaxInterval/SegmentTreeMaxOfInterval/RangeMaxSubSum
// 最大子段和/最大子数组和
// 最小子段和/最小子数组和
// 小白逛公园 : 动态子段和/区间最大子段和

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	小白逛公园()
}

// https://www.luogu.com.cn/problem/P4513
func 小白逛公园() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		var value int
		fmt.Fscan(in, &value)
		leaves[i] = FromElement(value)
	}

	seg := NewSegmentTreeInterval(leaves)

	for i := 0; i < m; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var start, end int
			fmt.Fscan(in, &start, &end)
			if start > end {
				start, end = end, start
			}
			start--
			res := seg.Query(start, end)
			fmt.Fprintln(out, res.max)
		} else {
			var pos, newValue int
			fmt.Fscan(in, &pos, &newValue)
			pos--
			seg.Set(pos, FromElement(newValue))
		}
	}
}

// 2382. 删除操作后的最大子段和
// https://leetcode.cn/problems/maximum-segment-sum-after-removals/description/
const INF int = 2e15

func maximumSegmentSum(nums []int, removeQueries []int) []int64 {
	leave := make([]E, len(nums))
	for i, v := range nums {
		leave[i] = FromElement(v)
	}
	seg := NewSegmentTreeInterval(leave)
	res := make([]int64, len(removeQueries))
	for i, v := range removeQueries {
		seg.Set(v, FromElement(-INF))
		res[i] = int64(seg.QueryAll().max)
		if res[i] < 0 {
			res[i] = 0
		}
	}
	return res
}

// 区间前缀和/区间后缀和.
type Interval struct {
	sum  int   // 区间和
	len  int32 // 区间长度
	max  int   // !连续子数组和的最大值(非空区间)
	lmax int   // !前缀和的最大值
	rmax int   // !后缀和的最大值
	min  int   // !连续子数组和的最小值(非空区间)
	lmin int   // !前缀和的最小值
	rmin int   // !后缀和的最小值
}

// 建立长度为1, 值为value的区间.
func FromElement(value int) Interval {
	return Interval{
		sum: value, len: 1,
		max: value, lmax: value, rmax: value,
		min: value, lmin: value, rmin: value,
	}
}

// 建立长度为length, 值为value的区间.
func FromElementWithLength(value int, length int) Interval {
	if length == 0 {
		return Interval{}
	}
	sum := value * length
	if sum > INF {
		sum = INF
	}
	if sum < -INF {
		sum = -INF
	}
	tmp1 := value
	if value > 0 {
		tmp1 *= length
	}
	tmp2 := value
	if value < 0 {
		tmp2 *= length
	}
	return Interval{
		sum: sum, len: int32(length),
		max: tmp1, lmax: tmp1, rmax: tmp1,
		min: tmp2, lmin: tmp2, rmin: tmp2,
	}
}

// 区间合并.
func Merge(this, other Interval) Interval {
	if this.len == 0 {
		return other
	}
	if other.len == 0 {
		return this
	}

	candSum := this.sum + other.sum
	if candSum > INF {
		candSum = INF
	}
	if candSum < -INF {
		candSum = -INF
	}
	return Interval{
		sum:  candSum,
		len:  this.len + other.len,
		max:  max(max(this.max, other.max), this.rmax+other.lmax),
		lmax: max(this.lmax, this.sum+other.lmax),
		rmax: max(other.rmax, other.sum+this.rmax),
		min:  min(min(this.min, other.min), this.rmin+other.lmin),
		lmin: min(this.lmin, this.sum+other.lmin),
		rmin: min(other.rmin, other.sum+this.rmin),
	}
}

type E = Interval

func (*SegmentTreeInterval) e() E        { return Interval{} }
func (*SegmentTreeInterval) op(a, b E) E { return Merge(a, b) }

type SegmentTreeInterval struct {
	n, size int
	seg     []E
}

func NewSegmentTreeInterval(leaves []E) *SegmentTreeInterval {
	res := &SegmentTreeInterval{}
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
func (st *SegmentTreeInterval) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeInterval) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeInterval) Update(index int, value E) {
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
func (st *SegmentTreeInterval) Query(start, end int) E {
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
func (st *SegmentTreeInterval) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *SegmentTreeInterval) MaxRight(left int, predicate func(E) bool) int {
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
func (st *SegmentTreeInterval) MinLeft(right int, predicate func(E) bool) int {
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

func (tree *SegmentTreeInterval) String() string {
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

// (松)离散化.
//
//	rank: 给定一个数，返回它的排名`(0-count)`.
//	count: 离散化(去重)后的元素个数.
func sortedSet(nums []int) (rank func(int) int, count int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	allNums := make([]int, 0, count)
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	rank = func(x int) int { return sort.SearchInts(allNums, x) }
	return
}
