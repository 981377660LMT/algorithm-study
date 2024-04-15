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
// 943 ms
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
	wm := NewWaveletMatrix2DStaticAbelGroup(n, func(i int32) (x, y XY, w AbelGroup) {
		return X[i], Y[i], AbelGroup(W[i])
	}, false, false)
	for i := int32(0); i < q; i++ {
		var a, c, b, d int32
		fmt.Fscan(in, &a, &c, &b, &d)
		fmt.Fprintln(out, wm.Query(a, b, c, d))
	}
}

func demo() {
	points := [][]int32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	wm := NewWaveletMatrix2DStaticAbelGroup(int32(len(points)), func(i int32) (x, y XY, w AbelGroup) {
		return XY(points[i][0]), XY(points[i][1]), AbelGroup(points[i][2])
	}, false, false)
	fmt.Println(wm.Count(1, 5, 2, 8)) // 2
	fmt.Println(wm.Query(1, 5, 2, 8)) // 9
}

type XY = int32
type Weight = int32
type AbelGroup = int

func e() AbelGroup                { return 0 }
func op(a, b AbelGroup) AbelGroup { return a + b }
func inv(a AbelGroup) AbelGroup   { return -a }

// y维度升序排列，x维度按照BinaryTrie构建.
type WaveletMatrix2DStaticAbelGroup struct {
	smallX, smallY bool
	xToIdx, yToIdx *ToIdx
	n, lg          int32
	mid            []int32
	bv             []*bitVector
	dat            [][]AbelGroup
}

func NewWaveletMatrix2DStaticAbelGroup(
	n int32, f func(i int32) (x, y XY, w AbelGroup),
	smallX, smallY bool,
) *WaveletMatrix2DStaticAbelGroup {
	res := &WaveletMatrix2DStaticAbelGroup{smallX: smallX, smallY: smallY, xToIdx: NewToIdx(), yToIdx: NewToIdx()}
	res._build(n, f)
	return res
}

func (wm *WaveletMatrix2DStaticAbelGroup) Count(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._prefixCount(y1, y2, x2) - wm._prefixCount(y1, y2, x1)
}

func (wm *WaveletMatrix2DStaticAbelGroup) Query(x1, x2, y1, y2 XY) AbelGroup {
	if wm.n == 0 {
		return e()
	}
	if x1 > x2 || y1 > y2 {
		return e()
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	add := wm._prefixSum(y1, y2, x2)
	sub := wm._prefixSum(y1, y2, x1)
	return op(add, inv(sub))
}

func (wm *WaveletMatrix2DStaticAbelGroup) _build(n int32, f func(i int32) (x, y XY, w AbelGroup)) {
	wm.n = n
	if n == 0 {
		wm.lg = 0
		return
	}
	A, Y, S := make([]XY, n), make([]XY, n), make([]AbelGroup, n)
	for i := int32(0); i < n; i++ {
		A[i], Y[i], S[i] = f(i)
	}
	order := argSort(Y)
	A = reArrange(A, order)
	Y = reArrange(Y, order)
	S = reArrange(S, order)
	wm.xToIdx.Build(A, wm.smallX)
	wm.yToIdx.Build(Y, wm.smallY)

	tmp := wm.xToIdx.Get(maxs(A, 0)+1, wm.smallX)
	wm.lg = int32(bits.Len32(uint32(tmp)))
	wm.mid = make([]int32, wm.lg)
	wm.bv = make([]*bitVector, wm.lg)
	for i := range wm.bv {
		wm.bv[i] = newBitVector(n)
	}
	wm.dat = make([][]AbelGroup, wm.lg+1)
	for i := range wm.dat {
		nums := make([]AbelGroup, n+1)
		for j := range nums {
			nums[j] = e()
		}
		wm.dat[i] = nums
	}
	for i := int32(0); i < n; i++ {
		A[i] = wm.xToIdx.Get(A[i], wm.smallX)
	}

	A0, A1 := make([]XY, n), make([]XY, n)
	S0, S1 := make([]AbelGroup, n), make([]AbelGroup, n)
	for d := wm.lg - 1; d >= -1; d-- {
		p0, p1 := int32(0), int32(0)
		tmp := wm.dat[d+1]
		for i := int32(0); i < n; i++ {
			tmp[i+1] = op(tmp[i], S[i])
		}
		if d == -1 {
			break
		}
		for i := int32(0); i < n; i++ {
			f := (A[i]>>d)&1 == 1
			if !f {
				S0[p0], A0[p0] = S[i], A[i]
				p0++
			} else {
				S1[p1], A1[p1] = S[i], A[i]
				wm.bv[d].Set(i)
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		A, A0 = A0, A
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			A[p0+i], S[p0+i] = A1[i], S1[i]
		}
	}
}

func (wm *WaveletMatrix2DStaticAbelGroup) _prefixCount(L, R, x int32) int32 {
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

func (wm *WaveletMatrix2DStaticAbelGroup) _prefixSum(L, R, x int32) AbelGroup {
	add, sub := e(), e()
	for d := wm.lg - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
		if (x>>d)&1 == 1 {
			add = op(add, wm.dat[d][r0])
			sub = op(sub, wm.dat[d][l0])
			L += wm.mid[d] - l0
			R += wm.mid[d] - r0
		} else {
			L, R = l0, r0
		}
	}
	return op(add, inv(sub))
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
