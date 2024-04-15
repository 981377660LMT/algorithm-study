package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	yosupo1()
	// yosupo2()
}

// RectangleSum
// https://judge.yosupo1.jp/problem/rectangle_sum
// 1219 ms
func yosupo1() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	X, Y, W := make([]int32, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &X[i], &Y[i], &W[i])
	}
	wm := NewWaveletMatrix2DDynamicMonoid(n, func(i int32) (x, y XY, w Monoid) {
		return X[i], Y[i], Monoid(W[i])
	}, false, false)
	for i := int32(0); i < q; i++ {
		var a, c, b, d int32
		fmt.Fscan(in, &a, &c, &b, &d)
		fmt.Fprintln(out, wm.Query(a, b, c, d))
	}
}

// PointAddRectangleSum
// https://judge.yosupo.jp/problem/point_add_rectangle_sum
// 716 ms
func yosupo2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	X, Y, W := make([]int32, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &X[i], &Y[i], &W[i])
	}
	type query = [4]int32
	queries := make([]query, q)
	for i := int32(0); i < q; i++ {
		var t int32
		fmt.Fscan(in, &t)
		if t == 0 {
			var x, y, w int32
			fmt.Fscan(in, &x, &y, &w)
			X = append(X, x)
			Y = append(Y, y)
			W = append(W, 0)
			queries[i] = query{-1, x, y, w}
		} else {
			var a, b, c, d int32
			fmt.Fscan(in, &a, &b, &c, &d)
			queries[i] = query{a, c, b, d}
		}
	}

	wm := NewWaveletMatrix2DDynamicMonoid(
		int32(len(X)), func(i int32) (x, y XY, w Monoid) { return X[i], Y[i], Monoid(W[i]) },
		false, false,
	)
	ptr := n
	for _, query := range queries {
		if query[0] == -1 {
			wm.Upadte(ptr, Monoid(query[3]))
			ptr++
		} else {
			fmt.Fprintln(out, wm.Query(query[0], query[1], query[2], query[3]))
		}
	}
}

func demo() {
	points := [][]int32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	wm := NewWaveletMatrix2DDynamicMonoid(
		int32(len(points)), func(i int32) (x, y XY, w Monoid) {
			return XY(points[i][0]), XY(points[i][1]), Monoid(points[i][2])
		}, false, false)
	fmt.Println(wm.Count(1, 5, 2, 8)) // 2
	fmt.Println(wm.Query(1, 5, 2, 8)) // 9
	wm.Set(0, 1)
	fmt.Println(wm.Query(1, 5, 2, 8)) // 7
	wm.Upadte(0, 100)
	fmt.Println(wm.Query(1, 5, 2, 8)) // 107
}

type XY = int32
type Weight = int32
type Monoid = int

func e() Monoid             { return 0 }
func op(a, b Monoid) Monoid { return a + b }

// y维度升序排列，x维度按照BinaryTrie构建.
type WaveletMatrix2DDynamicMonoid struct {
	smallX, smallY bool
	xToIdx, yToIdx *_toIdx
	n, lg          int32
	mid            []int32
	bv             []*_bitVector
	newIdx         []int32
	A              []int32
	dat            []*_segmentTree
}

func NewWaveletMatrix2DDynamicMonoid(
	n int32, f func(i int32) (x, y XY, w Monoid),
	smallX, smallY bool,
) *WaveletMatrix2DDynamicMonoid {
	res := &WaveletMatrix2DDynamicMonoid{smallX: smallX, smallY: smallY, xToIdx: _newToIdx(), yToIdx: _newToIdx()}
	res._build(n, f)
	return res
}

func (wm *WaveletMatrix2DDynamicMonoid) Count(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._prefixCount(y1, y2, x2) - wm._prefixCount(y1, y2, x1)
}

func (wm *WaveletMatrix2DDynamicMonoid) Query(x1, x2, y1, y2 XY) Monoid {
	if wm.n == 0 {
		return e()
	}
	if x1 > x2 || y1 > y2 {
		return e()
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	res := e()
	wm._queryDfs(y1, y2, x1, x2, wm.lg-1, &res)
	return res
}

func (wm *WaveletMatrix2DDynamicMonoid) Set(pointId int32, x Monoid) {
	pointId = wm.newIdx[pointId]
	a := wm.A[pointId]
	for d := wm.lg - 1; d >= 0; d-- {
		if (a>>d)&1 == 1 {
			pointId = wm.mid[d] + wm.bv[d].Rank(pointId, true)
		} else {
			pointId = wm.bv[d].Rank(pointId, false)
		}
		wm.dat[d].Set(pointId, x)
	}
}

func (wm *WaveletMatrix2DDynamicMonoid) Upadte(pointId int32, x Monoid) {
	pointId = wm.newIdx[pointId]
	a := wm.A[pointId]
	for d := wm.lg - 1; d >= 0; d-- {
		if (a>>d)&1 == 1 {
			pointId = wm.mid[d] + wm.bv[d].Rank(pointId, true)
		} else {
			pointId = wm.bv[d].Rank(pointId, false)
		}
		wm.dat[d].Update(pointId, x)
	}
}

func (wm *WaveletMatrix2DDynamicMonoid) _build(n int32, f func(i int32) (x, y XY, w Monoid)) {
	wm.n = n
	if n == 0 {
		wm.lg = 0
		return
	}
	tmp, Y, S := make([]XY, n), make([]XY, n), make([]Monoid, n)
	for i := int32(0); i < n; i++ {
		tmp[i], Y[i], S[i] = f(i)
	}
	order := argSort(Y)
	tmp = reArrange(tmp, order)
	Y = reArrange(Y, order)
	S = reArrange(S, order)
	wm.xToIdx.Build(tmp, wm.smallX)
	wm.yToIdx.Build(Y, wm.smallY)
	wm.newIdx = make([]int32, n)
	for i := range wm.newIdx {
		wm.newIdx[order[i]] = int32(i)
	}

	tmpMax := wm.xToIdx.Get(maxs(tmp, 0)+1, wm.smallX)
	wm.lg = int32(bits.Len32(uint32(tmpMax)))
	wm.mid = make([]int32, wm.lg)
	wm.bv = make([]*_bitVector, wm.lg)
	for i := range wm.bv {
		wm.bv[i] = _newBitVector(n)
	}
	wm.dat = make([]*_segmentTree, wm.lg)
	wm.A = make([]int32, n)
	for i := int32(0); i < n; i++ {
		wm.A[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}

	A0, A1 := make([]XY, n), make([]XY, n)
	S0, S1 := make([]Monoid, n), make([]Monoid, n)
	for d := wm.lg - 1; d >= 0; d-- {
		p0, p1 := int32(0), int32(0)
		for i := int32(0); i < n; i++ {
			f := (wm.A[i]>>d)&1 == 1
			if !f {
				S0[p0], A0[p0] = S[i], wm.A[i]
				p0++
			} else {
				S1[p1], A1[p1] = S[i], wm.A[i]
				wm.bv[d].Set(i)
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		wm.A, A0 = A0, wm.A
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			wm.A[p0+i] = A1[i]
			S[p0+i] = S1[i]
		}
		wm.dat[d] = _newSegmentTree(n, func(i int32) Monoid { return S[i] })
	}
	for i := int32(0); i < n; i++ {
		wm.A[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}
}

func (wm *WaveletMatrix2DDynamicMonoid) _prefixCount(L, R, x int32) int32 {
	cnt := int32(0)
	for d := wm.lg - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
		if (x>>d)&1 == 1 {
			cnt += r0 - l0
			L += wm.mid[d] - l0
			R += wm.mid[d] - r0
		} else {
			L, R = l0, r0
		}
	}
	return cnt
}

func (wm *WaveletMatrix2DDynamicMonoid) _queryDfs(L, R, x1, x2, d int32, res *Monoid) {
	if x1 < 0 {
		x1 = 0
	}
	if tmp := int32(1 << (d + 1)); x2 > tmp {
		x2 = tmp
	}
	if x1 >= x2 {
		return
	}
	if x1 == 0 && x2 == (1<<(d+1)) {
		*res = op(*res, wm.dat[d+1].Query(L, R))
		return
	}
	l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
	wm._queryDfs(l0, r0, x1, x2, d-1, res)
	wm._queryDfs(L+wm.mid[d]-l0, R+wm.mid[d]-r0, x1-(1<<d), x2-(1<<d), d-1, res)
}

type _toIdx struct {
	key    []XY
	mi, ma XY
	dat    []int32
}

func _newToIdx() *_toIdx { return &_toIdx{} }
func (ti *_toIdx) Build(X []XY, small bool) {
	if small {
		ti.mi, ti.ma = mins(X, 0), maxs(X, 0)
		ti.dat = make([]int32, ti.ma-ti.mi+2)
		for _, x := range X {
			ti.dat[x-ti.mi+1]++
		}
		for i := 1; i < len(ti.dat); i++ {
			ti.dat[i] += ti.dat[i-1]
		}
	} else {
		ti.key = make([]XY, len(X))
		copy(ti.key, X)
		sort.Slice(ti.key, func(i, j int) bool { return ti.key[i] < ti.key[j] })
	}
}

func (ti *_toIdx) Get(x XY, small bool) int32 {
	if small {
		return ti.dat[clamp(x-ti.mi, 0, ti.ma-ti.mi+1)]
	} else {
		return lb(ti.key, x)
	}
}

type _bitVector struct {
	bits   []uint64
	preSum []int32
}

func _newBitVector(n int32) *_bitVector {
	return &_bitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *_bitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *_bitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *_bitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
}

type _segmentTree struct {
	n, size int32
	seg     []Monoid
}

func _newSegmentTree(n int32, f func(int32) Monoid) *_segmentTree {
	res := &_segmentTree{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]Monoid, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *_segmentTree) Get(index int32) Monoid {
	if index < 0 || index >= st.n {
		return e()
	}
	return st.seg[index+st.size]
}
func (st *_segmentTree) Set(index int32, value Monoid) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *_segmentTree) Update(index int32, value Monoid) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *_segmentTree) Query(start, end int32) Monoid {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return e()
	}
	leftRes, rightRes := e(), e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}
func (st *_segmentTree) QueryAll() Monoid { return st.seg[1] }
func (st *_segmentTree) GetAll() []Monoid {
	res := make([]Monoid, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

func mins(nums []XY, defaultValue XY) XY {
	if len(nums) == 0 {
		return defaultValue
	}
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < res {
			res = nums[i]
		}
	}
	return res
}

func maxs(nums []XY, defaultValue XY) XY {
	if len(nums) == 0 {
		return defaultValue
	}
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > res {
			res = nums[i]
		}
	}
	return res
}

func clamp(x, mi, ma XY) XY {
	if x < mi {
		return mi
	}
	if x > ma {
		return ma
	}
	return x
}

func lb(nums []XY, x XY) int32 {
	left, right := int32(0), int32(len(nums)-1)
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

func argSort(nums []XY) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reArrange[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}
