// 二维线段树
// 基于 fractional cascading 实现
// https://dpair2005.github.io/articles/fc/
// https://maspypy.github.io/library/ds/segtree/segtree_2d.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF32 int32 = 1e9 + 10

func main() {
	abc360_g()
}

// func yosupoPointAddRectangleSum() {
// 	// https://judge.yosupo.jp/problem/point_add_rectangle_sum
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, q int32
// 	fmt.Fscan(in, &n, &q)
// 	xs, ys, ws := make([]int32, n), make([]int32, n), make([]int, n)
// 	for i := int32(0); i < n; i++ {
// 		fmt.Fscan(in, &xs[i], &ys[i], &ws[i])
// 	}
// 	queries := make([][4]int32, q)
// 	for i := int32(0); i < q; i++ {
// 		var t int32
// 		fmt.Fscan(in, &t)
// 		if t == 0 {
// 			var x, y, w int32
// 			fmt.Fscan(in, &x, &y, &w)
// 			xs = append(xs, x)
// 			ys = append(ys, y)
// 			ws = append(ws, 0)
// 			queries[i] = [4]int32{-1, x, y, w}
// 		} else {
// 			var a, b, c, d int32
// 			fmt.Fscan(in, &a, &b, &c, &d)
// 			queries[i] = [4]int32{a, c, b, d}
// 			xs = append(xs, 0)
// 			ys = append(ys, 0)
// 			ws = append(ws, 0)
// 		}
// 	}

// 	tree := NewSegmentTree2DSparse32WithWeights(xs, ys, ws, true)
// 	for i := int32(0); i < q; i++ {
// 		a, b, c, d := queries[i][0], queries[i][1], queries[i][2], queries[i][3]
// 		if a == -1 {
// 			tree.Update(n+i, int(d))
// 		} else {
// 			fmt.Fprintln(out, tree.Query(a, b, c, d))
// 		}
// 	}
// }

// G - Suitable Edit for LIS (修改一个数的LIS)
// https://atcoder.jp/contests/abc360/tasks/abc360_g
// 给定一个数组，进行一次操作，使得最长上升子序列的长度最大。
// 操作为，将任意一个数改为任意一个数。
//
// 前后缀分解
func abc360_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int32
	fmt.Fscan(in, &N)
	A := make([]int32, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &A[i])
	}

	dp1, dp2 := make([]int32, N), make([]int32, N)

	seg1 := NewSegmentTree2DSparse32Fast(N, func(i int32) (int32, int32, int32) { return i, A[i], -INF32 }, true)
	seg2 := NewSegmentTree2DSparse32Fast(N, func(i int32) (int32, int32, int32) { return i, A[i], -INF32 }, true)

	res := int32(0)
	for i := int32(0); i < N; i++ {
		x1, x2 := int32(1), int32(1)
		if i >= 1 {
			x2 = 2
		}
		if tmp := seg1.Query(0, i, 0, A[i]) + 1; tmp > x1 {
			x1 = tmp
		}
		if tmp := seg2.Query(0, i, 0, A[i]) + 1; tmp > x2 {
			x2 = tmp
		}
		if i >= 1 {
			if tmp := seg1.Query(0, i-1, 0, A[i]-1) + 2; tmp > x2 {
				x2 = tmp
			}
		}
		dp1[i], dp2[i] = x1, x2
		if x2 > res {
			res = x2
		}
		if i < N-1 {
			if tmp := x1 + 1; tmp > res {
				res = tmp
			}
		}

		seg1.Set(i, dp1[i])
		seg2.Set(i, dp2[i])
	}

	fmt.Fprintln(out, res)
}

// 2907.最大递增三元组的和
// https://leetcode.cn/problems/maximum-profitable-triplets-with-increasing-prices-ii/
func maxProfit(prices []int, profits []int) int {
	n := int32(len(prices))
	xs := make([]int32, n)
	for i := int32(0); i < n; i++ {
		xs[i] = int32(i)
	}
	prices32 := make([]int32, n)
	for i := int32(0); i < n; i++ {
		prices32[i] = int32(prices[i])
	}
	profits32 := make([]int32, n)
	for i := int32(0); i < n; i++ {
		profits32[i] = int32(profits[i])
	}

	tree := NewSegmentTree2DSparse32FastWithWeightsFrom(xs, prices32, profits32, false)
	res := int32(-1)
	for i := int32(0); i < n; i++ {
		x, y := i, prices32[i]
		max1 := tree.Query(0, x, 0, y)
		if max1 == 0 {
			continue
		}
		max2 := tree.Query(x+1, n, y+1, INF32)
		if max2 == 0 {
			continue
		}
		res = max32(res, max1+max2+profits32[i])
	}

	return int(res)
}

// 需要满足交换律.
type E = int32

func e() E        { return -INF32 }
func op(a, b E) E { return max32(a, b) } // TODO

type SegmentTree2DSparse32Fast struct {
	n          int32
	keyX       []int32
	keyY       []int32
	minX       int32
	allY       []int32
	pos        []int32
	indptr     []int32
	size       int32
	data       []E
	discretize bool
	unit       E
	toLeft     []int32
}

func NewSegmentTree2DSparse32Fast(n int32, f func(i int32) (int32, int32, E), discretize bool) *SegmentTree2DSparse32Fast {
	xs := make([]int32, n)
	ys := make([]int32, n)
	ws := make([]E, n)
	for i := int32(0); i < n; i++ {
		xs[i], ys[i], ws[i] = f(i)
	}
	return NewSegmentTree2DSparse32FastWithWeightsFrom(xs, ys, ws, discretize)
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(自动所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewSegmentTree2DSparse32FastFrom(xs, ys []int32, discretize bool) *SegmentTree2DSparse32Fast {
	res := &SegmentTree2DSparse32Fast{discretize: discretize, unit: e()}
	ws := make([]E, len(xs))
	for i := range ws {
		ws[i] = res.unit
	}
	res._build(xs, ys, ws)
	return res
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(自动所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewSegmentTree2DSparse32FastWithWeightsFrom(xs, ys []int32, ws []E, discretize bool) *SegmentTree2DSparse32Fast {
	res := &SegmentTree2DSparse32Fast{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

func (t *SegmentTree2DSparse32Fast) Update(rawIndex int32, value E) {
	i := int32(1)
	p := t.pos[rawIndex]
	indPtr, toLeft := t.indptr, t.toLeft
	for {
		t._update(i, p-indPtr[i], value)
		if i >= t.size {
			break
		}
		lc := toLeft[p] - toLeft[indPtr[i]]
		rc := p - indPtr[i] - lc
		if toLeft[p+1] > toLeft[p] {
			p = indPtr[i<<1] + lc
			i <<= 1
		} else {
			p = indPtr[i<<1|1] + rc
			i = i<<1 | 1
		}
	}
}

func (t *SegmentTree2DSparse32Fast) Set(rawIndex int32, value E) {
	i := int32(1)
	p := t.pos[rawIndex]
	indPtr, toLeft := t.indptr, t.toLeft
	for {
		t._set(i, p-indPtr[i], value)
		if i >= t.size {
			break
		}
		lc := toLeft[p] - toLeft[indPtr[i]]
		rc := p - indPtr[i] - lc
		if toLeft[p+1] > toLeft[p] {
			p = indPtr[i<<1] + lc
			i <<= 1
		} else {
			p = indPtr[i<<1|1] + rc
			i = i<<1 | 1
		}
	}
}

// [lx,rx) * [ly,ry)
func (t *SegmentTree2DSparse32Fast) Query(lx, rx, ly, ry int32) E {
	L := t._xtoi(lx)
	R := t._xtoi(rx)
	res := t.unit
	indPtr, toLeft := t.indptr, t.toLeft
	var dfs func(i, l, r, a, b int32)
	dfs = func(i, l, r, a, b int32) {
		if a == b || R <= l || r <= L {
			return
		}
		if L <= l && r <= R {
			res = op(res, t._query(i, a, b))
			return
		}
		la := toLeft[indPtr[i]+a] - toLeft[indPtr[i]]
		ra := a - la
		lb := toLeft[indPtr[i]+b] - toLeft[indPtr[i]]
		rb := b - lb
		m := (l + r) >> 1
		dfs(i<<1, l, m, la, lb)
		dfs(i<<1|1, m, r, ra, rb)
	}
	dfs(1, 0, t.size, bisectLeft(t.allY, ly, 0, int32(len(t.allY)-1)), bisectLeft(t.allY, ry, 0, int32(len(t.allY)-1)))
	return res
}

// nlogn.
func (seg *SegmentTree2DSparse32Fast) Count(lx, rx, ly, ry int32) int32 {
	L := seg._xtoi(lx)
	R := seg._xtoi(rx)
	res := int32(0)
	indPtr, toLeft := seg.indptr, seg.toLeft
	var dfs func(i, l, r, a, b int32)
	dfs = func(i, l, r, a, b int32) {
		if a == b || R <= l || r <= L {
			return
		}
		if L <= l && r <= R {
			res += b - a
			return
		}
		la := toLeft[indPtr[i]+a] - toLeft[indPtr[i]]
		ra := a - la
		lb := toLeft[indPtr[i]+b] - toLeft[indPtr[i]]
		rb := b - lb
		m := (l + r) >> 1
		dfs(i<<1, l, m, la, lb)
		dfs(i<<1|1, m, r, ra, rb)
	}
	dfs(1, 0, seg.size, bisectLeft(seg.allY, ly, 0, int32(len(seg.allY)-1)), bisectLeft(seg.allY, ry, 0, int32(len(seg.allY)-1)))
	return res
}

func (t *SegmentTree2DSparse32Fast) _update(i int32, y int32, val E) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	offset := lid << 1
	y += n
	for y > 0 {
		t.data[offset+y] = op(t.data[offset+y], val)
		y >>= 1
	}
}

func (seg *SegmentTree2DSparse32Fast) _set(i, y int32, val E) {
	lid := seg.indptr[i]
	n := seg.indptr[i+1] - seg.indptr[i]
	off := lid << 1
	y += n
	seg.data[off+y] = val
	for y > 1 {
		y >>= 1
		seg.data[off+y] = op(seg.data[off+y<<1], seg.data[off+y<<1|1])
	}
}

func (t *SegmentTree2DSparse32Fast) _query(i int32, ly, ry int32) E {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	offset := lid << 1
	left, right := n+ly, n+ry
	val := t.unit
	for left < right {
		if left&1 == 1 {
			val = op(val, t.data[offset+left])
			left++
		}
		if right&1 == 1 {
			right--
			val = op(t.data[offset+right], val)
		}
		left >>= 1
		right >>= 1
	}
	return val
}

func (seg *SegmentTree2DSparse32Fast) _build(X, Y []int32, wt []E) {
	if len(X) != len(Y) || len(X) != len(wt) {
		panic("Lengths of X, Y, and wt must be equal.")
	}

	if seg.discretize {
		seg.keyX = unique(X)
		seg.n = int32(len(seg.keyX))
	} else {
		if len(X) > 0 {
			min_, max_ := int32(0), int32(0)
			for _, x := range X {
				if x < min_ {
					min_ = x
				}
				if x > max_ {
					max_ = x
				}
			}
			seg.minX = min_
			seg.n = max_ - min_ + 1
		}
	}

	log := int32(0)
	for 1<<log < seg.n {
		log++
	}
	size := int32(1 << log)
	seg.size = size

	orderX := make([]int32, len(X))
	for i := range orderX {
		orderX[i] = seg._xtoi(X[i])
	}
	seg.indptr = make([]int32, 2*size+1)
	for _, i := range orderX {
		i += size
		for i > 0 {
			seg.indptr[i+1]++
			i >>= 1
		}
	}
	for i := int32(1); i <= 2*size; i++ {
		seg.indptr[i] += seg.indptr[i-1]
	}
	seg.data = make([]E, 2*seg.indptr[2*size])
	for i := range seg.data {
		seg.data[i] = seg.unit
	}

	seg.toLeft = make([]int32, seg.indptr[size]+1)
	ptr := append([]int32(nil), seg.indptr...)
	order := argSort(Y)
	seg.pos = make([]int32, len(X))
	for i, v := range order {
		seg.pos[v] = int32(i)
	}
	for _, rawIdx := range order {
		i := orderX[rawIdx] + size
		j := int32(-1)
		for i > 0 {
			p := ptr[i]
			ptr[i]++
			seg.data[seg.indptr[i+1]+p] = wt[rawIdx]
			if j != -1 && j&1 == 0 {
				seg.toLeft[p+1] = 1
			}
			j = i
			i >>= 1
		}
	}
	for i := int32(1); i < int32(len(seg.toLeft)); i++ {
		seg.toLeft[i] += seg.toLeft[i-1]
	}

	for i := int32(0); i < 2*size; i++ {
		off := 2 * seg.indptr[i]
		n := seg.indptr[i+1] - seg.indptr[i]
		for j := n - 1; j >= 1; j-- {
			seg.data[off+j] = op(seg.data[off+j<<1], seg.data[off+j<<1|1])
		}
	}

	allY := append([]int32(nil), Y...)
	sort.Slice(allY, func(i, j int) bool { return allY[i] < allY[j] })
	seg.allY = allY
}

func (seg *SegmentTree2DSparse32Fast) _xtoi(x int32) int32 {
	if seg.discretize {
		return bisectLeft(seg.keyX, x, 0, int32(len(seg.keyX)-1))
	}
	tmp := x - seg.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > seg.n {
		tmp = seg.n
	}
	return tmp
}

func bisectLeft(nums []int32, x int32, left, right int32) int32 {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func unique(nums []int32) (sorted []int32) {
	set := make(map[int32]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int32, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func argSort(nums []int32) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}
