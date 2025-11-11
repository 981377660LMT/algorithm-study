package main

import (
	"slices"
)

// https://leetcode.cn/problems/sum-of-subarray-minimums/description/
func sumSubarrayMins(arr []int) int {
	defaultVal := 0
	changesHistory := RangeMinChange(arr, func(a, b int) bool { return a < b }, defaultVal)

	const mod int = 1e9 + 7
	res := 0
	cur := 0 // 以当前下标为右端点的所有子数组最小值之和
	for _, changes := range changesHistory {
		for _, c := range changes {
			cnt := c.lEnd - c.lStart
			if c.oldMin == defaultVal {
				c.oldMin = 0
			}
			delta := (c.newMin - c.oldMin) * cnt
			cur = (cur + delta) % mod
		}
		if cur < 0 {
			cur += mod
		}
		res = (res + cur) % mod
	}
	return res
}

// https://leetcode.cn/problems/minimum-difficulty-of-a-job-schedule/
// 你需要制定一份 d 天的工作计划表。工作之间存在依赖，要想执行第 i 项工作，你必须完成全部 j 项工作（ 0 <= j < i）。
// 你每天 至少 需要完成一项任务。工作计划的总难度是这 d 天每一天的难度之和，而一天的工作难度是当天应该完成工作的最大难度。
// 给你一个整数数组 jobDifficulty 和一个整数 d，分别代表工作难度和需要计划的天数。第 i 项工作的难度是 jobDifficulty[i]。
// 返回整个工作计划的 最小难度 。如果无法制定工作计划，则返回 -1 。
// 1 <= jobDifficulty.length <= 300
// 0 <= jobDifficulty[i] <= 1000
// 1 <= d <= 10
// O(d·n·log n)
func minDifficulty(jobDifficulty []int, d int) int {
	n := len(jobDifficulty)
	if n < d {
		return -1
	}

	dp := make([]int, n+1) // dp[i] 表示在用 day 天覆盖前 i 个任务的最小难度
	dp[0] = 0
	mx := 0
	for i := 1; i <= n; i++ {
		if jobDifficulty[i-1] > mx {
			mx = jobDifficulty[i-1]
		}
		dp[i] = mx
	}
	if d == 1 {
		return dp[n]
	}

	const defaultVal = -1
	changesHistory := RangeMinChange(jobDifficulty, func(a, b int) bool { return a > b }, defaultVal)

	for day := 2; day <= d; day++ {
		seg := NewSegmentTreeRangeAddRangeMin(n, func(i int) int { return dp[i] })

		ndp := make([]int, n+1)
		for i := 0; i <= n; i++ {
			ndp[i] = INF
		}

		for i := 1; i <= n; i++ {
			for _, c := range changesHistory[i-1] {
				old := c.oldMin
				if old == defaultVal {
					old = 0
				}
				delta := c.newMin - old
				seg.UpdateRange(c.lStart, c.lEnd, delta)
			}
			// dp[day][i] = min_{j∈[day-1, i-1]} (dpPrev[j] + max(job[j..i-1]))
			ndp[i] = seg.Query(day-1, i)
		}

		dp = ndp
	}

	return dp[n]
}

type changeInfo[T any] struct {
	lStart, lEnd   int
	oldMin, newMin T
}

// 维护区间最小值的变化历史。
// 返回：res[i]，表示右端点r=i+1时，所有受影响区间[l,r)的最小值变化记录：(l, r, oldMin, newMin)
// 每次右端点推进，所有被当前元素“刷新”最小值的区间都会被记录下来，适用于区间DP、单调栈优化等场景。
func RangeMinChange[T any](arr []T, less func(a, b T) bool, defaultVal T) [][]changeInfo[T] {
	type stackInfo[T any] struct {
		lStart, lEnd int
		newMin       T
	}
	n := len(arr)
	res := make([][]changeInfo[T], n)
	var stack []stackInfo[T]
	for i, v := range arr {
		res[i] = append(res[i], changeInfo[T]{lStart: i, lEnd: i + 1, oldMin: defaultVal, newMin: v})
		ptr := i
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if !less(v, top.newMin) {
				break
			}
			res[i] = append(res[i], changeInfo[T]{lStart: top.lStart, lEnd: top.lEnd, oldMin: top.newMin, newMin: v})
			ptr = top.lStart
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, stackInfo[T]{lStart: ptr, lEnd: i + 1, newMin: v})
		slices.Reverse(res[i])
	}
	return res
}

// Utils

const INF int = 1e18

type summax = struct {
	sum, max int
}

type summin = struct {
	sum, min int
}

type SegmentTreeRangeAddRangeMax struct {
	n    int
	lazy int
	seg  *SegmentTreeGeneric[summax]
}

func NewSegmentTreeRangeAddRangeMax(n int, f func(int) int) *SegmentTreeRangeAddRangeMax {
	res := &SegmentTreeRangeAddRangeMax{}
	res.build(n, f)
	return res
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMax) Query(l, r int) int {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return -INF
	}
	res := st.seg.Query(l, r).max
	for l += st.seg.size; l > 0; l >>= 1 {
		if l&1 == 1 {
			l--
			res += st.seg.seg[l].sum
		}
	}
	return res + st.lazy
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMax) UpdateRange(l, r int, x int) {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return
	}
	if x == 0 {
		return
	}
	st.UpdateSuffix(l, x)
	st.UpdateSuffix(r, -x)
}

// [0, r)
func (st *SegmentTreeRangeAddRangeMax) UpdatePrefix(r int, x int) {
	if r > st.n {
		r = st.n
	}
	if r <= 0 {
		return
	}
	if x == 0 {
		return
	}
	st.lazy += x
	st.UpdateSuffix(r, -x)
}

// [l, n)
func (st *SegmentTreeRangeAddRangeMax) UpdateSuffix(l int, x int) {
	if l < 0 {
		l = 0
	}
	if l >= st.n {
		return
	}
	if x == 0 {
		return
	}
	t := st.seg.Get(l).sum + x
	st.seg.Set(l, summax{t, t})
}

func (st *SegmentTreeRangeAddRangeMax) UpdateAll(x int) {
	st.lazy += x
}

func (st *SegmentTreeRangeAddRangeMax) Set(i int, x int) {
	if i < 0 || i >= st.n {
		return
	}
	cur := st.Query(i, i+1)
	st.UpdateRange(i, i+1, x-cur)
}

func (st *SegmentTreeRangeAddRangeMax) Update(i int, x int) {
	if i < 0 || i >= st.n {
		return
	}
	cur := st.Query(i, i+1)
	if cur < x {
		st.UpdateRange(i, i+1, x-cur)
	}
}

func (st *SegmentTreeRangeAddRangeMax) build(n int, f func(int) int) {
	st.lazy = 0
	st.n = n
	pre := 0
	st.seg = NewSegmentTreeGeneric(
		n,
		func(i int) summax {
			t := f(i) - pre
			pre += t
			return summax{t, t}
		},
		func() summax { return summax{max: -2 * INF} },
		func(a, b summax) summax {
			a.max = max(a.max, a.sum+b.max)
			a.sum += b.sum
			return a
		},
	)
}

type SegmentTreeRangeAddRangeMin struct {
	n    int
	lazy int
	seg  *SegmentTreeGeneric[summin]
}

func NewSegmentTreeRangeAddRangeMin(n int, f func(int) int) *SegmentTreeRangeAddRangeMin {
	res := &SegmentTreeRangeAddRangeMin{}
	res.build(n, f)
	return res
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMin) Query(l, r int) int {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return INF
	}
	res := st.seg.Query(l, r).min
	for l += st.seg.size; l > 0; l >>= 1 {
		if l&1 == 1 {
			l--
			res += st.seg.seg[l].sum
		}
	}
	return res + st.lazy
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMin) UpdateRange(l, r int, x int) {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return
	}
	if x == 0 {
		return
	}
	st.UpdateSuffix(l, x)
	st.UpdateSuffix(r, -x)
}

// [0, r)
func (st *SegmentTreeRangeAddRangeMin) UpdatePrefix(r int, x int) {
	if r > st.n {
		r = st.n
	}
	if r <= 0 {
		return
	}
	if x == 0 {
		return
	}
	st.lazy += x
	st.UpdateSuffix(r, -x)
}

// [l, n)
func (st *SegmentTreeRangeAddRangeMin) UpdateSuffix(l int, x int) {
	if l < 0 {
		l = 0
	}
	if l >= st.n {
		return
	}
	if x == 0 {
		return
	}
	t := st.seg.Get(l).sum + x
	st.seg.Set(l, summin{t, t})
}

func (st *SegmentTreeRangeAddRangeMin) UpdateAll(x int) {
	if x == 0 {
		return
	}
	st.lazy += x
}

func (st *SegmentTreeRangeAddRangeMin) Set(i int, x int) {
	if i < 0 || i >= st.n {
		return
	}
	cur := st.Query(i, i+1)
	st.UpdateRange(i, i+1, x-cur)
}

func (st *SegmentTreeRangeAddRangeMin) Update(i int, x int) {
	if i < 0 || i >= st.n {
		return
	}
	cur := st.Query(i, i+1)
	if cur > x {
		st.UpdateRange(i, i+1, x-cur)
	}
}
func (st *SegmentTreeRangeAddRangeMin) build(n int, f func(int) int) {
	st.lazy = 0
	st.n = n
	pre := 0
	st.seg = NewSegmentTreeGeneric(
		n,
		func(i int) summin {
			t := f(i) - pre
			pre += t
			return summin{t, t}
		},
		func() summin { return summin{min: 2 * INF} },
		func(a, b summin) summin {
			a.min = min(a.min, a.sum+b.min)
			a.sum += b.sum
			return a
		},
	)
}

// SegmentTreeGeneric
type SegmentTreeGeneric[E any] struct {
	n, size int
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func NewSegmentTreeGeneric[E any](
	n int, f func(int) E,
	e func() E, op func(a, b E) E,
) *SegmentTreeGeneric[E] {
	res := &SegmentTreeGeneric[E]{e: e, op: op}
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
func (st *SegmentTreeGeneric[E]) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeGeneric[E]) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeGeneric[E]) Update(index int, value E) {
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
func (st *SegmentTreeGeneric[E]) Query(start, end int) E {
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
func (st *SegmentTreeGeneric[E]) QueryAll() E { return st.seg[1] }
func (st *SegmentTreeGeneric[E]) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
