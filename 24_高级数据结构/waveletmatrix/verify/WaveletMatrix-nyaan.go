// Usage:
// !Attention: nums[i] >= 0

// Count(start, end, value) – [start, end) 中值为 value 的数的个数.
// CountRange(start, end, lower, upper) – [start, end) 中值在 [lower, upper) 内的数的个数.

// Kth(start, end, k) – [start, end) 中第 k(0-indexed) 小的数.
// KthMax(start, end, k) – [start, end) 中第 k(0-indexed) 大的数.

// Lower(start, end, value) – [start, end) 中值小于 value 的最大值.不存在的话返回 -1.
// Higher(start, end, value) – [start, end) 中值大于 value 的最小值.不存在的话返回 -1.
// Floor(start, end, value) – [start, end) 中值不超过 value 的最大值.不存在的话返回 -1.
// Ceiling(start, end, value) – [start, end) 中值不小于 value 的最小值.不存在的话返回 -1.

package main

import (
	"fmt"
	"math/bits"
	"sort"
)

const INF int = 1e18

func main() {
	nums := []int{3, 5, 4, 6, 2, 1, 7, 8, 9, 0}
	wm := NewWaveletMatrix(nums)
	checkKth := func(start, end, k int) {
		cur := append([]int{}, nums[start:end]...)
		sort.Ints(cur)
		if cur[k] != wm.Kth(start, end, k) {
			panic("KthSmallest")
		}
	}
	checkFreq := func(start, end, lower, upper int) {
		cur := append([]int{}, nums[start:end]...)
		res := 0
		for _, v := range cur {
			if v >= lower && v < upper {
				res++
			}
		}
		if res != wm.CountRange(start, end, lower, upper) {
			panic("RangeFreq")
		}
	}
	checkPrev := func(start, end, upper int) {
		cur := append([]int{}, nums[start:end]...)
		res := -1
		for _, v := range cur {
			if v < upper {
				res = max(res, v)
			}
		}
		if res != wm.Lower(start, end, upper) {
			panic("PrevValue")
		}
	}

	for i := 0; i < len(nums); i++ {
		for j := i + 1; j <= len(nums); j++ {
			for k := 0; k < j-i; k++ {
				checkKth(i, j, k)
				checkFreq(i, j, k, k+100)
				checkPrev(i, j, k)
			}
		}
	}

	fmt.Println(wm.Lower(0, 10, 1))
	fmt.Println(wm.Ceiling(0, 10, 1))

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type WaveletMatrix struct {
	n, log int
	data   []int
	bv     []*BV
}

// 给定非负整数数组 nums 构建一个 WaveletMatrix.
func NewWaveletMatrix(data []int) *WaveletMatrix {
	n := uint32(len(data))
	cur := make([]int, n)
	dataCopy := make([]int, n)
	max_ := 0
	for i, v := range data {
		if v > max_ {
			max_ = v
		}
		cur[i] = v
		dataCopy[i] = v
	}

	maxLog := bits.Len(uint(max_)) + 1
	bv := make([]*BV, maxLog)
	for i := 0; i < maxLog; i++ {
		bv[i] = NewBV(n)
	}

	nxt := make([]int, n)
	for h := maxLog - 1; h >= 0; h-- {
		for i := uint32(0); i < n; i++ {
			if (cur[i]>>h)&1 == 1 {
				bv[h].Set(i)
			}
		}

		bv[h].Build()
		it := [2]uint32{0, bv[h].zeros}
		for i := uint32(0); i < n; i++ {
			pos := bv[h].Get(i)
			nxt[it[pos]] = cur[i]
			it[pos]++
		}
		cur, nxt = nxt, cur
	}

	return &WaveletMatrix{
		n:    int(n),
		log:  maxLog,
		data: dataCopy,
		bv:   bv,
	}
}

func (wm *WaveletMatrix) Kth(start, end, k int) int {
	res := 0
	ul, ur, uk := uint32(start), uint32(end), uint32(k)
	for h := wm.log - 1; h >= 0; h-- {
		lo, r0 := wm.bv[h].Rank0(ul), wm.bv[h].Rank0(ur)
		if uk < r0-lo {
			ul, ur = lo, r0
		} else {
			uk -= r0 - lo
			res |= 1 << h
			ul += wm.bv[h].zeros - lo
			ur += wm.bv[h].zeros - r0
		}
	}
	return res
}

func (wm *WaveletMatrix) KthMax(start, end, k int) int {
	return wm.Kth(start, end, end-start-k-1)
}

func (wm *WaveletMatrix) Count(start, end, x int) int {
	if x >= 1<<wm.log {
		return end - start
	}
	ul, ur := uint32(start), uint32(end)
	res := uint32(0)
	for h := wm.log - 1; h >= 0; h-- {
		f := (x>>h)&1 == 1
		l0, r0 := wm.bv[h].Rank0(ul), wm.bv[h].Rank0(ur)
		if f {
			res += r0 - l0
			ul += wm.bv[h].zeros - l0
			ur += wm.bv[h].zeros - r0
		} else {
			ul = l0
			ur = r0
		}
	}
	return int(res)
}

func (wm *WaveletMatrix) CountRange(start, end, floor, higher int) int {
	return wm.Count(start, end, higher) - wm.Count(start, end, floor)
}

func (wm *WaveletMatrix) Lower(start, end, x int) int {
	count := wm.Count(start, end, x)
	if count == 0 {
		return -1
	}
	return wm.Kth(start, end, count-1)
}

func (wm *WaveletMatrix) Ceiling(start, end, x int) int {
	count := wm.Count(start, end, x)
	if count == end-start {
		return -1
	}
	return wm.Kth(start, end, count)
}

type BV struct {
	n, zeros uint32
	block    []uint64
	count    []uint32
}

func NewBV(n uint32) *BV {
	bv := &BV{}
	bv.init(n)
	return bv
}

func (bv *BV) Set(i uint32) {
	bv.block[i/64] |= 1 << (i % 64)
}

func (bv *BV) Build() {
	for i := 1; i < (len(bv.block)); i++ {
		bv.count[i] = bv.count[i-1] + uint32(bits.OnesCount64(bv.block[i-1]))
	}
	bv.zeros = bv.Rank0(bv.n)
}

func (bv *BV) Get(i uint32) uint32 {
	return uint32((bv.block[i/64] >> (i % 64)) & 1)
}

func (bv *BV) init(n uint32) {
	bv.n = uint32(n)
	bv.zeros = uint32(n)
	bv.block = make([]uint64, n/64+1)
	bv.count = make([]uint32, len(bv.block))
}

// [0,i]中0的个数
func (bv *BV) Rank0(i uint32) uint32 {
	return i - bv.Rank1(i)
}

// [0,i]中1的个数
func (bv *BV) Rank1(i uint32) uint32 {
	return bv.count[i/64] + uint32(bits.OnesCount64((bv.block[i/64] & ((1 << (i % 64)) - 1))))
}
