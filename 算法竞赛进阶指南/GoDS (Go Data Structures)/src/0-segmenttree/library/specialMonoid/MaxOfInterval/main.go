package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {

}

// https://leetcode.cn/problems/longest-substring-of-one-repeating-character/
func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {

}

// https://yukicoder.me/problems/no/2281
func K101Flip() {}

// https://yukicoder.me/problems/no/2333
// 1 pos val: 将第pos个数修改为val
// 2 l r: 查询[l, r)区间内的相同值组成的最大连续子段和.
// !解法:离线预处理修改操作,将要修改的位置预先进行分割.
func SlimeStructure() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	slimes := make([][2]int, n) // (value,count)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &slimes[i][0], &slimes[i][1])
	}
	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var op, a, b int
		fmt.Fscan(in, &op, &a, &b)
		a--
		queries[i] = [3]int{op, a, b}
	}

	// TODO

}

// 区间前缀和/区间后缀和.
type Interval struct {
	sum  int // 区间和
	len  int // 区间长度
	max  int // !连续子数组和的最大值(非空区间)
	lmax int // !前缀和的最大值
	rmax int // !后缀和的最大值
	min  int // !连续子数组和的最小值(非空区间)
	lmin int // !前缀和的最小值
	rmin int // !后缀和的最小值
}

// 建立长度为1, 值为value的区间.
func NewInterval(value int) *Interval {
	return &Interval{
		sum: value, len: 1,
		max: value, lmax: value, rmax: value,
		min: value, lmin: value, rmin: value,
	}
}

// 建立长度为length, 值为value的区间.
func NewIntervalWithLength(value int, length int) *Interval {
	if length == 0 {
		return &Interval{}
	}
	sum := value * length
	tmp1 := value
	if value > 0 {
		tmp1 *= length
	}
	tmp2 := value
	if value < 0 {
		tmp2 *= length
	}
	return &Interval{
		sum: sum, len: length,
		max: tmp1, lmax: tmp1, rmax: tmp1,
		min: tmp2, lmin: tmp2, rmin: tmp2,
	}
}

func (this *Interval) Equals(other *Interval) bool {
	return this.sum == other.sum && this.len == other.len &&
		this.max == other.max && this.lmax == other.lmax && this.rmax == other.rmax &&
		this.min == other.min && this.lmin == other.lmin && this.rmin == other.rmin
}

func (this *Interval) IsEmpty() bool {
	return this.sum == 0 && this.len == 0 &&
		this.max == 0 && this.lmax == 0 && this.rmax == 0 &&
		this.min == 0 && this.lmin == 0 && this.rmin == 0
}

// 区间合并.
func (this *Interval) Merge(other *Interval) *Interval {
	if this.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return this
	}
	return &Interval{
		sum: this.sum + other.sum, len: this.len + other.len,
		max:  max(max(this.max, other.max), this.rmax+other.lmax),
		lmax: max(this.lmax, this.sum+other.lmax),
		rmax: max(other.rmax, other.sum+this.rmax),
		min:  min(min(this.min, other.min), this.rmin+other.lmin),
		lmin: min(this.lmin, this.sum+other.lmin),
		rmin: min(other.rmin, other.sum+this.rmin),
	}
}

type E = *Interval

func (*SegmentTreeInterval) e() E        { return &Interval{} }
func (*SegmentTreeInterval) op(a, b E) E { return a.Merge(b) }

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
//  rank: 给定一个数，返回它的排名`(0-count)`.
//  count: 离散化(去重)后的元素个数.
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
