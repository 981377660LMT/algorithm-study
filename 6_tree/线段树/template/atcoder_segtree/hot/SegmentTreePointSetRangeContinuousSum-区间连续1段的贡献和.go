// P9474 [yLOI2022] 长安幻世绘/长安浮世绘
// https://www.luogu.com.cn/problem/P9474
// 在元素互不相同的数组 a 中选出一个长度为 m 的元素互不相邻的子列，使得子列的极差最小。
// 不相邻选数的题目，可以用线段树维护
//
// 对nums排序，当灯亮度最小值为bl时，选出m个不相邻的灯需要的亮度最大的灯的亮度最小值为br.
// l增加时r单调不减，因此可以考虑双指针.
//
// !维护一个数据结构，支持单点修改01，查询所有极长1段的贡献之和(每个1段的贡献是ceil(len/2))
// SegTreeG-区间连续1段的贡献和
// !一个01数组，维护区间极长1段的贡献，这里的贡献取决于区间1段的长度

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// https://leetcode.cn/problems/count-stable-subarrays/description/
func countStableSubarrays(nums []int, queries [][]int) []int64 {
	n := len(nums)
	res := make([]int64, len(queries))

	m := n - 1
	nums01 := make([]bool, m)
	for i := 0; i < m; i++ {
		nums01[i] = nums[i] <= nums[i+1]
	}

	type E struct {
		preOnes int32
		sufOnes int32
		size    int32
		sum     int
	}
	eval := func(size int32) int {
		v := int(size)
		return v * (v + 1) / 2
	}
	of := func(b bool) E {
		if b {
			return E{preOnes: 1, sufOnes: 1, size: 1, sum: eval(1)}
		}
		return E{size: 1, sum: eval(0)}
	}

	e := func() E { return E{sum: eval(0)} }
	op := func(a, b E) E {
		res := E{preOnes: a.preOnes, sufOnes: b.sufOnes, size: a.size + b.size}
		if a.preOnes == a.size {
			res.preOnes += b.preOnes
		}
		if b.sufOnes == b.size {
			res.sufOnes += a.sufOnes
		}
		res.sum = a.sum + b.sum - eval(a.sufOnes) - eval(b.preOnes) + eval(a.sufOnes+b.preOnes)
		return res
	}

	seg := NewSegmentTreeG(e, op, m, func(i int) E { return of(nums01[i]) })

	set := func(i int, v int) {
		if nums[i] == v {
			return
		}
		nums[i] = v
		if i > 0 {
			b := nums[i-1] <= nums[i]
			nums01[i-1] = b
			seg.Set(i-1, of(b))
		}
		if i < n-1 {
			b := nums[i] <= nums[i+1]
			nums01[i] = b
			seg.Set(i, of(b))
		}
	}
	_ = set

	for i, q := range queries {
		left, right := q[0], q[1]
		cur := (right - left + 1) + seg.Query(left, right).sum
		res[i] = int64(cur)
	}
	return res
}

func luoguP9474() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	const INF int = 1e18

	type E = struct {
		preOnes int32 // 区间前缀1的个数
		sufOnes int32 // 区间后缀1的个数
		size    int32 // 区间长度
		sum     int   // 区间内连续1段的贡献和
	}
	eval := func(size int32) int { return int((size + 1) / 2) } // !计算区间连续1段的贡献
	of := func(b bool) E {
		if b {
			return E{preOnes: 1, sufOnes: 1, size: 1, sum: eval(1)}
		} else {
			return E{size: 1, sum: eval(0)}
		}
	}

	e := func() E { return E{sum: eval(0)} }
	op := func(a, b E) E {
		res := E{preOnes: a.preOnes, sufOnes: b.sufOnes, size: a.size + b.size}
		if a.preOnes == a.size {
			res.preOnes += b.preOnes
		}
		if b.sufOnes == b.size {
			res.sufOnes += a.sufOnes
		}
		res.sum = a.sum + b.sum - eval(a.sufOnes) - eval(b.preOnes) + eval(a.sufOnes+b.preOnes)
		return res
	}

	seg := NewSegmentTreeG(e, op, n, func(i int) E { return of(false) })

	type pair struct{ v, id int }
	numsWithId := make([]pair, n)
	for i := range nums {
		numsWithId[i] = pair{v: nums[i], id: i}
	}
	sort.Slice(numsWithId, func(i, j int) bool { return numsWithId[i].v < numsWithId[j].v })
	res := INF
	left := 0
	for right := 0; right < n; right++ {
		seg.Set(numsWithId[right].id, of(true))
		for left <= right && seg.QueryAll().sum >= m {
			res = min(res, numsWithId[right].v-numsWithId[left].v)
			seg.Set(numsWithId[left].id, of(false))
			left++
		}
	}

	fmt.Fprintln(out, res)
}

// SegmentTreePointSetRangeContinuousOnesSum

type SegmentTreeG[E any] struct {
	n, size int
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func NewSegmentTreeG[E any](
	e func() E, op func(a, b E) E,
	n int, f func(int) E,
) *SegmentTreeG[E] {
	res := &SegmentTreeG[E]{e: e, op: op}
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
func (st *SegmentTreeG[E]) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeG[E]) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeG[E]) Update(index int, value E) {
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
func (st *SegmentTreeG[E]) Query(start, end int) E {
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
func (st *SegmentTreeG[E]) QueryAll() E { return st.seg[1] }
func (st *SegmentTreeG[E]) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
