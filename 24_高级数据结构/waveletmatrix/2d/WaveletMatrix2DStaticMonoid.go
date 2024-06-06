package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/rectangle_sum
// 1250ms
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	X, Y, W := make([]int32, n), make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &X[i], &Y[i], &W[i])
	}
	wm := NewWaveletMatrix2DStaticMonoid(
		n, func(i int32) (x, y XY, w Monoid) { return X[i], Y[i], Monoid(W[i]) },
		func(m int32, f func(i int32) Monoid) IRMQ[Monoid] { return NewDisjointSparseTable(m, f, e, op) },
		false, false,
	)
	for i := int32(0); i < q; i++ {
		var a, c, b, d int32
		fmt.Fscan(in, &a, &c, &b, &d)
		fmt.Fprintln(out, wm.Query(a, b, c, d))
	}
}

func demo() {
	points := [][]int32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	wm := NewWaveletMatrix2DStaticMonoid(
		int32(len(points)), func(i int32) (x, y XY, w Monoid) { return points[i][0], points[i][1], Monoid(points[i][2]) },
		func(m int32, f func(i int32) Monoid) IRMQ[Monoid] { return NewDisjointSparseTable(m, f, e, op) },
		false, false,
	)
	fmt.Println(wm.Count(1, 5, 2, 8)) // 2
	fmt.Println(wm.Query(1, 5, 2, 8)) // 9
}

const ST_LOG int32 = 4 // 数据较大5，数据较小4

type XY = int32
type Weight = int32
type Monoid = int

func e() Monoid             { return 0 }
func op(a, b Monoid) Monoid { return a + b }

// y维度升序排列，x维度按照BinaryTrie构建.
type WaveletMatrix2DStaticMonoid struct {
	smallX, smallY bool
	xToIdx, yToIdx *_toIdx
	n, lg          int32
	mid            []int32
	bv             []*_bitVector
	newIdx         []int32
	a              []int32
	dat            []IRMQ[Monoid]
}

func NewWaveletMatrix2DStaticMonoid(
	n int32, f func(i int32) (x, y XY, w Monoid),
	createRMQ func(m int32, f func(i int32) Monoid) IRMQ[Monoid],
	smallX, smallY bool,
) *WaveletMatrix2DStaticMonoid {
	res := &WaveletMatrix2DStaticMonoid{smallX: smallX, smallY: smallY, xToIdx: _newToIdx(), yToIdx: _newToIdx()}
	res._build(n, f, createRMQ)
	return res
}

func (wm *WaveletMatrix2DStaticMonoid) Count(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	if x1 >= x2 || y1 >= y2 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._prefixCount(y1, y2, x2) - wm._prefixCount(y1, y2, x1)
}

func (wm *WaveletMatrix2DStaticMonoid) Query(x1, x2, y1, y2 XY) Monoid {
	if wm.n == 0 {
		return e()
	}
	if x1 >= x2 || y1 >= y2 {
		return e()
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	res := e()
	wm._queryDfs(y1, y2, x1, x2, wm.lg-1, &res)
	return res
}

func (wm *WaveletMatrix2DStaticMonoid) _build(
	n int32, f func(i int32) (x, y XY, w Monoid),
	createRMQ func(m int32, f func(i int32) Monoid) IRMQ[Monoid],
) {
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
	for i := int32(0); i < n; i++ {
		wm.newIdx[order[i]] = i
	}

	if len(tmp) > 0 {
		tmpMax := wm.xToIdx.Get(maxs(tmp, 0)+1, wm.smallX)
		wm.lg = int32(bits.Len32(uint32(tmpMax)))
	}
	wm.mid = make([]int32, wm.lg)
	wm.bv = make([]*_bitVector, wm.lg)
	for i := range wm.bv {
		wm.bv[i] = _newBitVector(n)
	}
	wm.dat = make([]IRMQ[Monoid], wm.lg)
	wm.a = make([]int32, n)
	for i := int32(0); i < n; i++ {
		wm.a[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}

	A0, A1 := make([]XY, n), make([]XY, n)
	S0, S1 := make([]Monoid, n), make([]Monoid, n)
	for d := wm.lg - 1; d >= 0; d-- {
		p0, p1 := int32(0), int32(0)
		for i := int32(0); i < n; i++ {
			f := (wm.a[i]>>d)&1 == 1
			if !f {
				S0[p0], A0[p0] = S[i], wm.a[i]
				p0++
			} else {
				S1[p1], A1[p1] = S[i], wm.a[i]
				wm.bv[d].Set(i)
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		wm.a, A0 = A0, wm.a
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			wm.a[p0+i], S[p0+i] = A1[i], S1[i]
		}
		wm.dat[d] = _newStaticRangeProduct[Monoid](
			n, func(i int32) Monoid { return S[i] }, e, op, ST_LOG, createRMQ,
		)
	}
	for i := int32(0); i < n; i++ {
		wm.a[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}
}

func (wm *WaveletMatrix2DStaticMonoid) _prefixCount(L, R, x int32) int32 {
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

func (wm *WaveletMatrix2DStaticMonoid) _queryDfs(L, R, x1, x2, d int32, res *Monoid) {
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

type IRMQ[E any] interface {
	Query(start, end int32) E
}

type _StaticRangeProduct[E any] struct {
	n, log        int32
	arr, pre, suf []E // inclusive
	e             func() E
	op            func(a, b E) E
	rmq           IRMQ[E]
}

// log: 一般为4.
func _newStaticRangeProduct[E any](
	n int32, f func(i int32) E, e func() E, op func(a, b E) E, log int32,
	createRMQ func(n int32, f func(i int32) E) IRMQ[E],
) *_StaticRangeProduct[E] {
	bNum := n >> log
	arr := make([]E, n)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i)
	}
	pre := append(arr[:0:0], arr...)
	suf := append(arr[:0:0], arr...)
	mask := int32((1 << log) - 1)
	for i := int32(1); i < n; i++ {
		if i&mask != 0 {
			pre[i] = op(pre[i-1], arr[i])
		}
	}
	for i := n - 1; i >= 1; i-- {
		if i&mask != 0 {
			suf[i-1] = op(arr[i-1], suf[i])
		}
	}
	rmq := createRMQ(bNum, func(i int32) E { return suf[i<<log] })
	return &_StaticRangeProduct[E]{n: n, log: log, arr: arr, pre: pre, suf: suf, e: e, op: op, rmq: rmq}
}

func (s *_StaticRangeProduct[E]) Query(start, end int32) E {
	if start >= end {
		return s.e()
	}
	end--
	a, b := start>>s.log, end>>s.log
	if a < b {
		x := s.rmq.Query(a+1, b)
		x = s.op(s.suf[start], x)
		x = s.op(x, s.pre[end])
		return x
	}
	x := s.arr[start]
	for i := start + 1; i <= end; i++ {
		x = s.op(x, s.arr[i])
	}
	return x
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (s *_StaticRangeProduct[E]) MaxRight(left int32, check func(e E) bool) int32 {
	if left == s.n {
		return s.n
	}
	ok, ng := left, s.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(s.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (s *_StaticRangeProduct[E]) MinLeft(right int32, check func(e E) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(s.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// SparseTable 稀疏表: st[j][i] 表示区间 [i, i+2^j) 的贡献值.
type SparseTable[E any] struct {
	st [][]E
	e  func() E
	op func(E, E) E
	n  int32
}

func NewSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *SparseTable[E] {
	res := &SparseTable[E]{}
	b := bits.Len32(uint32(n))
	st := make([][]E, b)
	for i := range st {
		st[i] = make([]E, n)
	}
	for i := int32(0); i < n; i++ {
		st[0][i] = f(i)
	}
	for i := 1; i < b; i++ {
		for j := int32(1); j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	res.st = st
	res.e = e
	res.op = op
	res.n = n
	return res
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	b := bits.Len32(uint32(end-start)) - 1 // log2
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

type DisjointSparseTable[E any] struct {
	n    int32
	e    func() E
	op   func(E, E) E
	data [][]E
}

// DisjointSparseTable 支持幺半群的区间静态查询.
//
//	eg: 区间乘积取模/区间仿射变换...
func NewDisjointSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *DisjointSparseTable[E] {
	res := &DisjointSparseTable[E]{}
	log := int32(1)
	for (1 << log) < n {
		log++
	}
	data := make([][]E, log)
	data[0] = make([]E, 0, n)
	for i := int32(0); i < n; i++ {
		data[0] = append(data[0], f(i))
	}
	for i := int32(1); i < log; i++ {
		data[i] = append(data[i], data[0]...)
		tmp := data[i]
		b := int32(1 << i)
		for m := b; m <= n; m += 2 * b {
			l, r := m-b, min32(m+b, n)
			for j := m - 1; j >= l+1; j-- {
				tmp[j-1] = op(tmp[j-1], tmp[j])
			}
			for j := m; j < r-1; j++ {
				tmp[j+1] = op(tmp[j], tmp[j+1])
			}
		}
	}
	res.n = n
	res.e = e
	res.op = op
	res.data = data
	return res
}

func (ds *DisjointSparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return ds.e()
	}
	end--
	if start == end {
		return ds.data[0][start]
	}
	lca := bits.Len32(uint32(start^end)) - 1
	return ds.op(ds.data[lca][start], ds.data[lca][end])
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
