package main

import (
	"math/bits"
	"sort"
)

const INF32 int32 = 1e9 + 10

// https://leetcode.cn/problems/find-the-number-of-ways-to-place-people-ii/
// 如果围栏的 内部 或者 边缘 上有任何其他人，Alice 都会难过。
func numberOfPairs(points [][]int) int {
	wm := NewWaveletMatrix2DDynamicAbelGroup(
		int32(len(points)),
		func(i int32) (x, y XY, w bool) {
			return int32(points[i][0]), int32(points[i][1]), true
		},
		false,
		false,
	)

	res := int32(0)
	for i := int32(0); i < int32(len(points)); i++ {
		x1, y1 := int32(points[i][0]), int32(points[i][1])
		for j := int32(0); j < int32(len(points)); j++ {
			if i == j {
				continue
			}
			x2, y2 := int32(points[j][0]), int32(points[j][1])
			count := wm.Count(x1, x2+1, y2, y1+1)
			if count == 2 {
				res++
			}
		}
	}

	return int(res)
}

type XY = int32

// y维度升序排列，x维度按照BinaryTrie构建.
type WaveletMatrix2DDynamicBit01 struct {
	smallX, smallY bool
	xToIdx, yToIdx *ToIdx
	n, lg          int32
	mid            []int32
	bv             []*bitVector
	newIdx         []int32
	A              []int32
	dat            []*bitArray0132
}

func NewWaveletMatrix2DDynamicAbelGroup(
	n int32, f func(i int32) (x, y XY, w bool),
	smallX, smallY bool,
) *WaveletMatrix2DDynamicBit01 {
	res := &WaveletMatrix2DDynamicBit01{smallX: smallX, smallY: smallY, xToIdx: NewToIdx(), yToIdx: NewToIdx()}
	res._build(n, f)
	return res
}

// pid: 初始化时的点的索引.
func (wm *WaveletMatrix2DDynamicBit01) Add(pid int32) {
	pid = wm.newIdx[pid]
	a := wm.A[pid]
	for d := wm.lg - 1; d >= 0; d-- {
		if (a>>d)&1 == 1 {
			pid = wm.mid[d] + wm.bv[d].Rank(pid, true)
		} else {
			pid = wm.bv[d].Rank(pid, false)
		}
		wm.dat[d].Add(pid)
	}
}

func (wm *WaveletMatrix2DDynamicBit01) Remove(pid int32) {
	pid = wm.newIdx[pid]
	a := wm.A[pid]
	for d := wm.lg - 1; d >= 0; d-- {
		if (a>>d)&1 == 1 {
			pid = wm.mid[d] + wm.bv[d].Rank(pid, true)
		} else {
			pid = wm.bv[d].Rank(pid, false)
		}
		wm.dat[d].Remove(pid)
	}
}

func (wm *WaveletMatrix2DDynamicBit01) Count(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	if x1 >= x2 || y1 >= y2 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._countInner(y1, y2, x2) - wm._countInner(y1, y2, x1)
}

func (wm *WaveletMatrix2DDynamicBit01) Query(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	if x1 >= x2 || y1 >= y2 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._sumInner(y1, y2, x2) - wm._sumInner(y1, y2, x1)
}

func (wm *WaveletMatrix2DDynamicBit01) QueryPrefix(x, y XY) int32 {
	if wm.n == 0 {
		return 0
	}
	x = wm.xToIdx.Get(x, wm.smallX)
	y = wm.yToIdx.Get(y, wm.smallY)
	return wm._sumInner(0, y, x)
}

func (wm *WaveletMatrix2DDynamicBit01) _build(n int32, f func(i int32) (x, y XY, w bool)) {
	wm.n = n
	if n == 0 {
		wm.lg = 0
		return
	}
	tmp, Y, S := make([]XY, n), make([]XY, n), make([]bool, n)
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
	wm.dat = make([]*bitArray0132, wm.lg)
	wm.A = make([]int32, n)
	for i := int32(0); i < n; i++ {
		wm.A[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}

	A0, A1 := make([]XY, n), make([]XY, n)
	S0, S1 := make([]bool, n), make([]bool, n)
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
		wm.dat[d] = newBitArray0132From(n, func(i int32) bool { return S[i] })
	}
	for i := int32(0); i < n; i++ {
		wm.A[i] = wm.xToIdx.Get(tmp[i], wm.smallX)
	}
}

func (wm *WaveletMatrix2DDynamicBit01) _countInner(L, R, x int32) int32 {
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

func (wm *WaveletMatrix2DDynamicBit01) _sumInner(L, R, x int32) int32 {
	if x == 0 {
		return 0
	}
	sum := int32(0)
	for d := wm.lg - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
		if (x>>d)&1 == 1 {
			sum += wm.dat[d].QueryRange(l0, r0)
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

// 01树状数组.
type bitArray0132 struct {
	n    int32
	size int32 // data、bit的长度
	data []uint64
	bit  *bitArray32
}

func newBitArray0132(n int32) *bitArray0132 {
	size := int32(n>>6 + 1)
	data := make([]uint64, size)
	bit := newBitArray32(size)
	return &bitArray0132{n: n, size: size, data: data, bit: bit}
}

func newBitArray0132From(n int32, f func(index int32) bool) *bitArray0132 {
	size := n>>6 + 1
	data := make([]uint64, size)
	for i := int32(0); i < n; i++ {
		if f(i) {
			data[i>>6] |= 1 << (i & 63)
		}
	}
	bit := newBitArray32From(size, func(i int32) int32 { return int32(bits.OnesCount64(data[i])) })
	return &bitArray0132{n: n, size: size, data: data, bit: bit}
}
func (bit01 *bitArray0132) QueryPrefix(end int32) int32 {
	i, j := end>>6, end&63
	res := bit01.bit.QueryPrefix(i)
	res += int32(bits.OnesCount64(bit01.data[i] & ((1 << j) - 1)))
	return res
}

func (bit01 *bitArray0132) QueryRange(start, end int32) int32 {
	if start >= end {
		return 0
	}
	if start == 0 {
		return bit01.QueryPrefix(end)
	}
	res := int32(0)
	res -= int32(bits.OnesCount64(bit01.data[start>>6] & ((1 << (start & 63)) - 1)))
	res += int32(bits.OnesCount64(bit01.data[end>>6] & ((1 << (end & 63)) - 1)))
	res += bit01.bit.QueryRange(start>>6, end>>6)
	return res
}

func (bit01 *bitArray0132) Add(index int32) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 1 {
		return false
	}
	bit01.data[i] |= 1 << j
	bit01.bit.Add(i, 1)
	return true
}

func (bit01 *bitArray0132) Remove(index int32) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 0 {
		return false
	}
	bit01.data[i] ^= 1 << j
	bit01.bit.Add(i, -1)
	return true
}

func (bit01 *bitArray0132) Has(index int32) bool {
	i, j := index>>6, index&63
	return (bit01.data[i]>>j)&1 == 1
}

// !Point Add Range Sum, 0-based.
type bitArray32 struct {
	n    int32
	data []int32
}

func newBitArray32(n int32) *bitArray32 {
	res := &bitArray32{n: n, data: make([]int32, n)}
	return res
}

func newBitArray32From(n int32, f func(i int32) int32) *bitArray32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &bitArray32{n: n, data: data}
}

func (b *bitArray32) Add(index int32, v int32) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *bitArray32) QueryPrefix(end int32) int32 {
	res := int32(0)
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *bitArray32) QueryRange(start, end int32) int32 {
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := int32(0), int32(0)
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
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
