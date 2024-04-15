package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// yosupo1()
	yosupo2()
}

// RectangleSum
// https://judge.yosupo1.jp/problem/rectangle_sum
// 1085 ms
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
	wm := NewWaveletMatrix2DDynamicAbelGroup(n, func(i int32) (x, y XY, w AbelGroup) {
		return X[i], Y[i], AbelGroup(W[i])
	}, false, false)
	for i := int32(0); i < q; i++ {
		var a, c, b, d int32
		fmt.Fscan(in, &a, &c, &b, &d)
		fmt.Fprintln(out, wm.Query(a, b, c, d))
	}
}

// PointAddRectangleSum
// https://judge.yosupo.jp/problem/point_add_rectangle_sum
// 524 ms
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

	wm := NewWaveletMatrix2DDynamicAbelGroup(
		int32(len(X)), func(i int32) (x, y XY, w AbelGroup) { return X[i], Y[i], AbelGroup(W[i]) },
		false, false,
	)
	ptr := n
	for _, query := range queries {
		if query[0] == -1 {
			wm.Add(ptr, AbelGroup(query[3]))
			ptr++
		} else {
			fmt.Fprintln(out, wm.Query(query[0], query[1], query[2], query[3]))
		}
	}
}

func demo() {
	points := [][]int32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	wm := NewWaveletMatrix2DDynamicAbelGroup(int32(len(points)), func(i int32) (x, y XY, w AbelGroup) {
		return XY(points[i][0]), XY(points[i][1]), AbelGroup(points[i][2])
	}, false, false)
	fmt.Println(wm.Count(1, 5, 2, 8)) // 2
	fmt.Println(wm.Query(1, 5, 2, 8)) // 9
	fmt.Println(wm.QueryPrefix(2, 3)) // 3
	wm.Add(0, 1)
	fmt.Println(wm.QueryPrefix(2, 3)) // 4
}

type XY = int32
type Weight = int32
type AbelGroup = int

func e() AbelGroup                { return 0 }
func op(a, b AbelGroup) AbelGroup { return a + b }
func inv(a AbelGroup) AbelGroup   { return -a }

// y维度升序排列，x维度按照BinaryTrie构建.
type WaveletMatrix2DDynamicAbelGroup struct {
	smallX, smallY bool
	xToIdx, yToIdx *ToIdx
	n, lg          int32
	mid            []int32
	bv             []*bitVector
	newIdx         []int32
	A              []int32
	dat            []*bitGroup
}

func NewWaveletMatrix2DDynamicAbelGroup(
	n int32, f func(i int32) (x, y XY, w AbelGroup),
	smallX, smallY bool,
) *WaveletMatrix2DDynamicAbelGroup {
	res := &WaveletMatrix2DDynamicAbelGroup{smallX: smallX, smallY: smallY, xToIdx: NewToIdx(), yToIdx: NewToIdx()}
	res._build(n, f)
	return res
}

// pointIndex: 初始化时的点的索引.
// w: 权值.
func (wm *WaveletMatrix2DDynamicAbelGroup) Add(pointIndex int32, w AbelGroup) {
	pointIndex = wm.newIdx[pointIndex]
	a := wm.A[pointIndex]
	for d := wm.lg - 1; d >= 0; d-- {
		if (a>>d)&1 == 1 {
			pointIndex = wm.mid[d] + wm.bv[d].Rank(pointIndex, true)
		} else {
			pointIndex = wm.bv[d].Rank(pointIndex, false)
		}
		wm.dat[d].Update(pointIndex, w)
	}
}

func (wm *WaveletMatrix2DDynamicAbelGroup) Count(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._countInner(y1, y2, x2) - wm._countInner(y1, y2, x1)
}

func (wm *WaveletMatrix2DDynamicAbelGroup) Query(x1, x2, y1, y2 XY) AbelGroup {
	if wm.n == 0 {
		return e()
	}
	if x1 > x2 || y1 > y2 {
		return e()
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	add := wm._sumInner(y1, y2, x2)
	sub := wm._sumInner(y1, y2, x1)
	return op(add, inv(sub))
}

func (wm *WaveletMatrix2DDynamicAbelGroup) QueryPrefix(x, y XY) AbelGroup {
	if wm.n == 0 {
		return e()
	}
	x = wm.xToIdx.Get(x, wm.smallX)
	y = wm.yToIdx.Get(y, wm.smallY)
	return wm._sumInner(0, y, x)
}

func (wm *WaveletMatrix2DDynamicAbelGroup) _build(n int32, f func(i int32) (x, y XY, w AbelGroup)) {
	wm.n = n
	if n == 0 {
		wm.lg = 0
		return
	}
	tmp, Y, S := make([]XY, n), make([]XY, n), make([]AbelGroup, n)
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
	wm.bv = make([]*bitVector, wm.lg)
	for i := range wm.bv {
		wm.bv[i] = newBitVector(n)
	}
	wm.dat = make([]*bitGroup, wm.lg)
	wm.A = make([]int32, n)
	for i := int32(0); i < n; i++ {
		wm.A[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}

	A0, A1 := make([]XY, n), make([]XY, n)
	S0, S1 := make([]AbelGroup, n), make([]AbelGroup, n)
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
		wm.dat[d] = NewBITGroupFrom(n, func(i int32) AbelGroup { return S[i] })
	}
	for i := int32(0); i < n; i++ {
		wm.A[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}
}

func (wm *WaveletMatrix2DDynamicAbelGroup) _countInner(L, R, x int32) int32 {
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

func (wm *WaveletMatrix2DDynamicAbelGroup) _sumInner(L, R, x int32) AbelGroup {
	if x == 0 {
		return e()
	}
	sum := e()
	for d := wm.lg - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
		if (x>>d)&1 == 1 {
			sum = op(sum, wm.dat[d].QueryRange(l0, r0))
			L += wm.mid[d] - l0
			R += wm.mid[d] - r0
		} else {
			L, R = l0, r0
		}
	}
	return sum
}

type ToIdx struct {
	key    []XY
	mi, ma XY
	dat    []int32
}

func NewToIdx() *ToIdx { return &ToIdx{} }
func (ti *ToIdx) Build(X []XY, small bool) {
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

func (ti *ToIdx) Get(x XY, small bool) int32 {
	if small {
		return ti.dat[clamp(x-ti.mi, 0, ti.ma-ti.mi+1)]
	} else {
		return lb(ti.key, x)
	}
}

type bitVector struct {
	bits   []uint64
	preSum []int32
}

func newBitVector(n int32) *bitVector {
	return &bitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *bitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *bitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *bitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
}

type bitGroup struct {
	n     int32
	data  []AbelGroup
	total AbelGroup
}

func newBITGroup(n int32) *bitGroup {
	data := make([]AbelGroup, n)
	for i := range data {
		data[i] = e()
	}
	return &bitGroup{n: n, data: data, total: e()}
}

func NewBITGroupFrom(n int32, f func(index int32) AbelGroup) *bitGroup {
	total := e()
	data := make([]AbelGroup, n)
	for i := range data {
		data[i] = f(int32(i))
		total = op(total, data[i])
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &bitGroup{n: n, data: data, total: total}
}

func (fw *bitGroup) Update(i int32, x AbelGroup) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *bitGroup) QueryAll() AbelGroup { return fw.total }

// [0, end)
func (fw *bitGroup) QueryPrefix(end int32) AbelGroup {
	if end > fw.n {
		end = fw.n
	}
	res := e()
	for end > 0 {
		res = op(res, fw.data[end-1])
		end &= end - 1
	}
	return res
}

// [start, end)
func (fw *bitGroup) QueryRange(start, end int32) AbelGroup {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	if start > end {
		return e()
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, fw.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
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
