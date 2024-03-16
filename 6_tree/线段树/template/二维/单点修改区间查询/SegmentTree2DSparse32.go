// 二维线段树区间修改单点查询
// Query(lx, rx, ly, ry) 查询区间[lx, rx) * [ly, ry)的值.
// Update(x, y, val) 将点(x, y)的值加上(op) val.
// !每次修改/查询 O(logn*logn)

// 1e5:1200ms; 2e5:2200ms
// https://kopricky.github.io/code/SegmentTrees/rangetree_pointupdate.html
// RangeTree

package main

import (
	"sort"
)

const INF32 int32 = 1e9 + 10

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

	tree := NewSegmentTree2DSparse32WithWeights(xs, prices32, profits32, false)
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

func e() E        { return 0 }
func op(a, b E) E { return max32(a, b) }

type SegmentTree2DSparse32 struct {
	n          int32
	keyX       []int32
	keyY       []int32
	minX       int32
	indptr     []int32
	data       []E
	discretize bool
	unit       E
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(自动所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewSegmentTree2DSparse32(xs, ys []int32, discretize bool) *SegmentTree2DSparse32 {
	res := &SegmentTree2DSparse32{discretize: discretize, unit: e()}
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
func NewSegmentTree2DSparse32WithWeights(xs, ys []int32, ws []E, discretize bool) *SegmentTree2DSparse32 {
	res := &SegmentTree2DSparse32{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

// 点 (x,y) 的值加上 val.
func (t *SegmentTree2DSparse32) Update(x, y int32, val E) {
	i := t._xtoi(x)
	i += t.n
	for i > 0 {
		t._add(i, y, val)
		i >>= 1
	}
}

// [lx,rx) * [ly,ry)
func (t *SegmentTree2DSparse32) Query(lx, rx, ly, ry int32) E {
	L := t._xtoi(lx) + t.n
	R := t._xtoi(rx) + t.n
	val := t.unit
	for L < R {
		if L&1 == 1 {
			val = op(val, t._prodI(L, ly, ry))
			L++
		}
		if R&1 == 1 {
			R--
			val = op(t._prodI(R, ly, ry), val)
		}
		L >>= 1
		R >>= 1
	}
	return val
}

func (t *SegmentTree2DSparse32) _add(i int32, y int32, val E) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	j := bisectLeft(t.keyY, y, lid, lid+n-1) - lid
	offset := 2 * lid
	j += n
	for j > 0 {
		t.data[offset+j] = op(t.data[offset+j], val)
		j >>= 1
	}
}

func (t *SegmentTree2DSparse32) _prodI(i int32, ly, ry int32) E {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	left := bisectLeft(t.keyY, ly, lid, lid+n-1) - lid + n
	right := bisectLeft(t.keyY, ry, lid, lid+n-1) - lid + n
	offset := 2 * lid
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

func (seg *SegmentTree2DSparse32) _build(X, Y []int32, wt []E) {
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
		seg.keyX = make([]int32, seg.n)
		for i := int32(0); i < seg.n; i++ {
			seg.keyX[i] = seg.minX + i
		}
	}

	N := seg.n
	keyYRaw := make([][]int32, N+N)
	datRaw := make([][]E, N+N)
	indices := argSort(Y)

	for _, i := range indices {
		ix, y := seg._xtoi(X[i]), Y[i]
		ix += N
		for ix > 0 {
			KY := keyYRaw[ix]
			if len(KY) == 0 || KY[len(KY)-1] < y {
				keyYRaw[ix] = append(keyYRaw[ix], y)
				datRaw[ix] = append(datRaw[ix], wt[i])
			} else {
				datRaw[ix][len(datRaw[ix])-1] = op(datRaw[ix][len(datRaw[ix])-1], wt[i])
			}
			ix >>= 1
		}
	}

	seg.indptr = make([]int32, N+N+1)
	for i := int32(0); i < N+N; i++ {
		seg.indptr[i+1] = seg.indptr[i] + int32(len(keyYRaw[i]))
	}
	fullN := seg.indptr[N+N]
	seg.keyY = make([]int32, fullN)
	seg.data = make([]E, 2*fullN)

	for i := int32(0); i < N+N; i++ {
		off := 2 * seg.indptr[i]
		n := seg.indptr[i+1] - seg.indptr[i]
		for j := int32(0); j < n; j++ {
			seg.keyY[seg.indptr[i]+j] = keyYRaw[i][j]
			seg.data[off+n+j] = datRaw[i][j]
		}
		for j := n - 1; j >= 1; j-- {
			seg.data[off+j] = op(seg.data[off+2*j], seg.data[off+2*j+1])
		}
	}
}

func (seg *SegmentTree2DSparse32) _xtoi(x int32) int32 {
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
