// https://github.com/MitI-7/WaveletMatrix/blob/master/WaveletMatrix/WaveletMatrix.hpp
// WaveletMatrixStaticOmni/StaticWaveletMatrixOmni
// api:
//  1. PrefixCount(end uint64, x uint64) uint64
//  2. RangeCount(start, end uint64, x uint64) uint64
//  3. RangeFreq(start, end uint64, x, y int) uint64
//  4. Kth(k uint64, x int) uint64
//  5. KthSmallest(start, end uint64, k uint64) int
//     KthSmallestIndex(start, end uint64, k uint64) int
//  6. KthLargest(start, end uint64, k uint64) int
//     KthLargestIndex(start, end uint64, k uint64) int
//  7. CountAll(start, end uint64, x int) uint64 -> 统计[start, end)中==x, <x, >x的个数
//  8. CountLessThan(start, end uint64, x int) uint64
//  9. CountMoreThan(start, end uint64, x int) uint64
//  10. Floor(start, end uint64, x int) (int, bool)
//  11. Lower(start, end uint64, x int) (int, bool)
//  12. Ceil(start, end uint64, x int) (int, bool)
//  13. Higher(start, end uint64, x int) (int, bool)

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"time"
)

func main() {
	// demo()
	// test()
	testTime()
}
func demo() {

	nums := []uint64{3, 4, 0, 2, 1}
	wm := NewWaveletMatrixStaticFast(uint64(len(nums)), func(i uint64) uint64 { return nums[i] })
	_ = wm

}

type WaveletMatrixStaticFast struct {
	bvs           []*succinctBitVector
	beginOne      []uint64
	beginAlphabet map[uint64]uint64
	size          uint64
	maxElement    uint64
	bitSize       uint64
}

func NewWaveletMatrixStaticFast(n uint64, f func(uint64) uint64) *WaveletMatrixStaticFast {
	if n <= 0 {
		panic("n must be positive")
	}
	data := make([]uint64, n)
	for i := uint64(0); i < n; i++ {
		data[i] = f(i)
	}
	maxElement := uint64(0)
	for i := uint64(0); i < n; i++ {
		if data[i] > maxElement {
			maxElement = data[i]
		}
	}
	maxElement++
	bitSize := uint64(bits.Len64(maxElement))
	if bitSize == 0 {
		bitSize = 1
	}
	bvs := make([]*succinctBitVector, bitSize)
	for i := uint64(0); i < bitSize; i++ {
		bvs[i] = newSuccinctBitVector(n)
	}
	beginOne := make([]uint64, bitSize)
	beginAlphabet := make(map[uint64]uint64, n)

	zero, one := make([]uint64, n), make([]uint64, n)
	for i := int8(bitSize - 1); i >= 0; i-- {
		bv := bvs[i]
		p, q := uint64(0), uint64(0)
		for j := uint64(0); j < n; j++ {
			c := data[j]
			bit := (c >> i) & 1
			if bit == 0 {
				zero[p] = c
				p++
			} else {
				one[q] = c
				q++
				bv.Set(j, true)
			}
		}
		beginOne[i] = p
		bv.Build()
		data, zero = zero, data
		copy(data[p:], one[:q])
	}

	for i := int32(n - 1); i >= 0; i-- {
		beginAlphabet[data[i]] = uint64(i)
	}

	return &WaveletMatrixStaticFast{
		bvs:           bvs,
		beginOne:      beginOne,
		beginAlphabet: beginAlphabet,
		size:          n,
		maxElement:    maxElement,
		bitSize:       bitSize,
	}
}

func (wm *WaveletMatrixStaticFast) Get(pos uint64) uint64 {
	if pos >= wm.size {
		return NOT_FOUND
	}
	c := uint64(0)
	for i := uint64(0); i < wm.bitSize; i++ {
		bit := wm.bvs[i].Get(pos)
		c <<= 1
		if bit {
			c |= 1
		}
		pos = wm.bvs[i].Count(pos, bit)
		if bit {
			pos += wm.beginOne[i]
		}
	}
	return c
}

func (wm *WaveletMatrixStaticFast) PrefixCount(end uint64, v uint64) uint64 {
	if end <= 0 {
		return 0
	}
	if end > wm.size {
		end = wm.size
	}
	if v >= wm.maxElement {
		return 0
	}
	beginPos, ok := wm.beginAlphabet[v]
	if !ok {
		return 0
	}
	for i := int8(wm.bitSize - 1); i >= 0; i-- {
		bit := v >> i & 1
		end = wm.bvs[i].Count(end, bit == 1)
		if bit == 1 {
			end += wm.beginOne[i]
		}
	}
	return end - beginPos
}

func (wm *WaveletMatrixStaticFast) RangeCount(start, end uint64, v uint64) uint64 {
	return wm.PrefixCount(end, v) - wm.PrefixCount(start, v)
}

func (wm *WaveletMatrixStaticFast) RangeFreq(start, end uint64, floor, higher uint64) uint64 {
	if end > wm.size || start >= end || floor >= higher || floor >= wm.maxElement {
		return 0
	}
	return wm.CountLessThan(start, end, higher) - wm.CountLessThan(start, end, floor)
}

// k: 0-indexed.
func (wm *WaveletMatrixStaticFast) Kth(k uint64, v uint64) uint64 {
	if v >= wm.maxElement {
		return NOT_FOUND
	}
	var s uint64
	if tmp, ok := wm.beginAlphabet[v]; !ok {
		return NOT_FOUND
	} else {
		s = tmp + k
	}
	for i := uint64(0); i < wm.bitSize; i++ {
		bit := (v >> i) & 1
		if bit == 1 {
			s = wm.bvs[i].Kth(s-wm.bvs[i].Count(wm.size, false), true)
		} else {
			s = wm.bvs[i].Kth(s, false)
		}
		if s == NOT_FOUND {
			return NOT_FOUND
		}
	}
	return s
}

func (wm *WaveletMatrixStaticFast) KthSmallestIndex(start, end uint64, k uint64) uint64 {
	if end > wm.size || start >= end || k >= end-start {
		return NOT_FOUND
	}
	val := uint64(0)
	for i := int8(wm.bitSize - 1); i >= 0; i-- {
		numOfZeroBegin := wm.bvs[i].Count(start, false)
		numOfZeroEnd := wm.bvs[i].Count(end, false)
		numOfZero := numOfZeroEnd - numOfZeroBegin
		bit := uint64(0)
		if k >= numOfZero {
			bit = 1
		}
		if bit == 1 {
			k -= numOfZero
			start = wm.beginOne[i] + start - numOfZeroBegin
			end = wm.beginOne[i] + end - numOfZeroEnd
		} else {
			start = numOfZeroBegin
			end = numOfZeroBegin + numOfZero
		}
		val = (val << 1) | bit
	}

	left := uint64(0)
	for i := int8(wm.bitSize - 1); i >= 0; i-- {
		bit := (val >> i) & 1
		left = wm.bvs[i].Count(left, bit == 1)
		if bit == 1 {
			left += wm.beginOne[i]
		}
	}
	rank := start + k - left
	return wm.Kth(rank, val)
}

func (wm *WaveletMatrixStaticFast) KthSmallest(start, end uint64, k uint64) uint64 {
	if end > wm.size || start >= end || k >= end-start {
		return NOT_FOUND
	}
	res := uint64(0)
	for bit := int8(wm.bitSize - 1); bit >= 0; bit-- {
		l0, r0 := wm.bvs[bit].Count(start, false), wm.bvs[bit].Count(end, false)
		if c := r0 - l0; c <= k {
			res |= 1 << bit
			k -= c
			start += wm.beginOne[bit] - l0
			end += wm.beginOne[bit] - r0
		} else {
			start, end = l0, r0
		}
	}
	return res
}

func (wm *WaveletMatrixStaticFast) KthLargest(start, end uint64, k uint64) uint64 {
	return wm.KthSmallest(start, end, end-start-k-1)
}

func (wm *WaveletMatrixStaticFast) KthLargestIndex(start, end uint64, k uint64) uint64 {
	return wm.KthSmallestIndex(start, end, end-start-k-1)
}

// 区间[start, end)中等于v的个数、小于v的个数、大于v的个数.
func (wm *WaveletMatrixStaticFast) CountAll(start, end uint64, v uint64) (uint64, uint64, uint64) {
	if start < 0 {
		start = 0
	}
	if end > wm.size {
		end = wm.size
	}
	if start >= end {
		return 0, 0, 0
	}
	num := end - start
	if v >= wm.maxElement {
		return 0, num, 0
	}
	rankLessThan, rankMoreThan := uint64(0), uint64(0)
	for i := int8(wm.bitSize - 1); i >= 0 && start < end; i-- {
		bit := v >> i & 1
		rank0Begin := wm.bvs[i].Count(start, false)
		rank0End := wm.bvs[i].Count(end, false)
		rank1Begin := start - rank0Begin
		rank1End := end - rank0End
		if bit == 1 {
			rankLessThan += rank0End - rank0Begin
			start = wm.beginOne[i] + rank1Begin
			end = wm.beginOne[i] + rank1End
		} else {
			rankMoreThan += rank1End - rank1Begin
			start = rank0Begin
			end = rank0End
		}
	}
	rank := num - rankLessThan - rankMoreThan
	return rank, rankLessThan, rankMoreThan
}

func (wm *WaveletMatrixStaticFast) CountLessThan(start, end uint64, v uint64) uint64 {
	_, less, _ := wm.CountAll(start, end, v)
	return less
}

func (wm *WaveletMatrixStaticFast) CountMoreThan(start, end uint64, v uint64) uint64 {
	_, _, more := wm.CountAll(start, end, v)
	return more
}

func (wm *WaveletMatrixStaticFast) Floor(start, end uint64, x uint64) uint64 {
	same, less, _ := wm.CountAll(start, end, x)
	if same > 0 {
		return x
	}
	if less == 0 {
		return NOT_FOUND
	}
	return wm.KthSmallest(start, end, less-1)
}
func (wm *WaveletMatrixStaticFast) Lower(start, end uint64, x uint64) uint64 {
	_, less, _ := wm.CountAll(start, end, x)
	if less == 0 {
		return NOT_FOUND
	}
	return wm.KthSmallest(start, end, less-1)
}

func (wm *WaveletMatrixStaticFast) Ceil(start, end uint64, x uint64) uint64 {
	same, less, _ := wm.CountAll(start, end, x)
	if same > 0 {
		return x
	}
	if less == end-start {
		return NOT_FOUND
	}
	return wm.KthSmallest(start, end, less)
}

func (wm *WaveletMatrixStaticFast) Higher(start, end uint64, x uint64) uint64 {
	if start >= end {
		return NOT_FOUND
	}
	_, less, _ := wm.CountAll(start, end, x+1)
	if less == end-start {
		return NOT_FOUND
	}
	return wm.KthSmallest(start, end, less)
}

const NOT_FOUND = ^uint64(0)

type succinctBitVector struct {
	size  uint64
	ones  uint64
	large []uint64 // 大块
	small []uint16 // 小块
	bv    []uint16 // bitVector
}

func newSuccinctBitVector(n uint64) *succinctBitVector {
	res := &succinctBitVector{size: n}
	res.bv = make([]uint16, (n+15)>>4+1)
	res.large = make([]uint64, n>>9+1)
	res.small = make([]uint16, n>>4+1)
	return res
}

func (sbv *succinctBitVector) Set(pos uint64, bit bool) {
	blockPos := pos >> 4
	offset := pos & 15
	if bit {
		sbv.bv[blockPos] |= 1 << offset
	} else {
		sbv.bv[blockPos] &= ^(1 << offset)
	}
}

func (sbv *succinctBitVector) Get(pos uint64) bool {
	blockPos := pos >> 4
	offset := pos & 15
	return (sbv.bv[blockPos]>>offset)&1 == 1
}

func (sbv *succinctBitVector) Build() {
	num := uint64(0)
	for i := uint64(0); i <= sbv.size; i++ {
		if i&511 == 0 {
			sbv.large[i>>9] = num
		}
		if i&15 == 0 {
			sbv.small[i>>4] = uint16(num - sbv.large[i>>9])
		}
		if i != sbv.size && i&15 == 0 {
			num += uint64(bits.OnesCount16(sbv.bv[i>>4]))
		}
	}
	sbv.ones = num
}

func (sbv *succinctBitVector) Count(end uint64, bit bool) uint64 {
	if end <= 0 {
		return 0
	}
	if end >= sbv.size {
		end = sbv.size
	}
	if bit {
		return sbv.large[end>>9] + uint64(sbv.small[end>>4]) + uint64(bits.OnesCount16(sbv.bv[end>>4]&(1<<(end&15)-1)))
	}
	return end - sbv.Count(end, true)
}

// !kth 从0开始.
func (sbv *succinctBitVector) Kth(k uint64, bit bool) uint64 {
	k++
	if !bit && k > sbv.size-sbv.ones {
		return NOT_FOUND
	}
	if bit && k > sbv.ones {
		return NOT_FOUND
	}

	// 大块内搜索
	largeIndex := uint64(0)
	{
		left, right := uint64(0), uint64(len(sbv.large))
		for right-left > 1 {
			mid := (left + right) >> 1
			var r uint64
			if bit {
				r = sbv.large[mid]
			} else {
				r = mid<<9 - sbv.large[mid]
			}
			if r < k {
				left = mid
				largeIndex = mid
			} else {
				right = mid
			}
		}
	}

	// 小块内搜索
	smallIndex := (largeIndex << 9) >> 4
	{
		left, right := (largeIndex<<9)>>4, min64(((largeIndex+1)<<9)>>4, uint64(len(sbv.small)))
		for right-left > 1 {
			mid := (left + right) >> 1
			r := sbv.large[largeIndex] + uint64(sbv.small[mid])
			if !bit {
				r = mid<<4 - r
			}
			if r < k {
				left = mid
				smallIndex = mid
			} else {
				right = mid
			}
		}
	}

	// bitVector内搜索
	rankPos := uint64(0)
	{
		beginBlockIndex := (smallIndex << 4) >> 4
		totalBit := sbv.large[largeIndex] + uint64(sbv.small[smallIndex])
		if !bit {
			totalBit = smallIndex<<4 - totalBit
		}
		for i := uint64(0); ; i++ {
			b := uint64(bits.OnesCount16(sbv.bv[beginBlockIndex+i]))
			if !bit {
				b = 16 - b
			}
			if totalBit+b >= k {
				block := uint64(sbv.bv[beginBlockIndex+i])
				if !bit {
					block = ^block
				}
				rankPos = (beginBlockIndex+i)<<4 + sbv._selectInBlock(block, k-totalBit)
				break
			}
			totalBit += b
		}
	}

	return rankPos
}

func (sbv *succinctBitVector) _selectInBlock(x uint64, rank uint64) uint64 {
	x1 := x - ((x & 0xAAAAAAAAAAAAAAAA) >> 1)
	x2 := (x1 & 0x3333333333333333) + ((x1 >> 2) & 0x3333333333333333)
	x3 := (x2 + (x2 >> 4)) & 0x0F0F0F0F0F0F0F0F
	pos := uint64(0)
	for ; ; pos += 8 {
		rankNext := (x3 >> pos) & 0xFF
		if rank <= rankNext {
			break
		}
		rank -= rankNext
	}
	v2 := (x2 >> pos) & 0xF
	if rank > v2 {
		rank -= v2
		pos += 4
	}
	v1 := (x1 >> pos) & 0x3
	if rank > v1 {
		rank -= v1
		pos += 2
	}
	v0 := (x >> pos) & 0x1
	if v0 < rank {
		rank -= v0
		pos += 1
	}
	return pos
}

func min64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func test() {
	for i := 0; i < 10; i++ {
		nums := make([]uint64, 5)
		for j := 0; j < len(nums); j++ {
			nums[j] = uint64(rand.Intn(5))
		}
		wm := NewWaveletMatrixStaticFast(uint64(len(nums)), func(u uint64) uint64 {
			return nums[u]
		})

		prefixBf := func(end uint64, x uint64) uint64 {
			res := uint64(0)
			for i := uint64(0); i < end; i++ {
				if nums[i] == x {
					res++
				}
			}
			return res
		}

		rangeBf := func(start, end uint64, x uint64) uint64 {
			res := uint64(0)
			for i := start; i < end; i++ {
				if nums[i] == x {
					res++
				}
			}
			return res
		}

		rangeFreqBf := func(start, end uint64, x, y uint64) uint64 {
			res := uint64(0)
			for i := start; i < end; i++ {
				if nums[i] >= x && nums[i] < y {
					res++
				}
			}
			return res
		}

		kthBf := func(k uint64, x uint64) uint64 {
			cnt := uint64(0)
			for i := uint64(0); i < uint64(len(nums)); i++ {
				if nums[i] == x {
					if cnt == k {
						return i
					}
					cnt++
				}
			}
			return NOT_FOUND
		}
		_ = kthBf

		kthSmallestBf := func(start, end, k uint64) uint64 {
			arr := make([]uint64, 0, end-start)
			for i := start; i < end; i++ {
				arr = append(arr, nums[i])
			}
			sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
			if int(k) >= len(arr) {
				return NOT_FOUND
			}
			return arr[k]
		}
		_ = kthSmallestBf

		kthSmallestIndexBf := func(start, end, k uint64) uint64 {
			valueWithIndex := make([][2]uint64, 0, end-start)
			for i := start; i < end; i++ {
				valueWithIndex = append(valueWithIndex, [2]uint64{nums[i], i})
			}
			sort.Slice(valueWithIndex, func(i, j int) bool {
				if valueWithIndex[i][0] != valueWithIndex[j][0] {
					return valueWithIndex[i][0] < valueWithIndex[j][0]
				}
				return valueWithIndex[i][1] < valueWithIndex[j][1]
			})
			if int(k) >= len(valueWithIndex) {
				return NOT_FOUND
			}
			return valueWithIndex[k][1]
		}
		_ = kthSmallestIndexBf

		kthLargestBf := func(start, end, k uint64) uint64 {
			arr := make([]uint64, 0, end-start)
			for i := start; i < end; i++ {
				arr = append(arr, nums[i])
			}
			sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
			if int(k) >= len(arr) {
				return NOT_FOUND
			}
			return arr[len(arr)-1-int(k)]
		}
		_ = kthLargestBf

		countAllBf := func(start, end, v uint64) (uint64, uint64, uint64) {
			rank, rankLessThan, rankMoreThan := uint64(0), uint64(0), uint64(0)
			for i := start; i < end; i++ {
				if nums[i] == v {
					rank++
				}
				if nums[i] < v {
					rankLessThan++
				}
				if nums[i] > v {
					rankMoreThan++
				}
			}
			return rank, rankLessThan, rankMoreThan
		}

		floorBf := func(start, end uint64, x uint64) uint64 {
			res := NOT_FOUND
			for i := start; i < end; i++ {
				if nums[i] <= x {
					if res == NOT_FOUND || nums[i] > res {
						res = nums[i]
					}
				}
			}
			return res
		}

		lowerBf := func(start, end uint64, x uint64) uint64 {
			res := NOT_FOUND
			for i := start; i < end; i++ {
				if nums[i] < x {
					if res == NOT_FOUND || nums[i] > res {
						res = nums[i]
					}
				}
			}
			return res
		}

		ceilBf := func(start, end uint64, x uint64) uint64 {
			res := NOT_FOUND
			for i := start; i < end; i++ {
				if nums[i] >= x {
					if res == NOT_FOUND || nums[i] < res {
						res = nums[i]
					}
				}
			}
			return res
		}

		higherBf := func(start, end uint64, x uint64) uint64 {
			res := NOT_FOUND
			for i := start; i < end; i++ {
				if nums[i] > x {
					if res == NOT_FOUND || nums[i] < res {
						res = nums[i]
					}
				}
			}
			return res
		}
		_ = higherBf

		for j := 0; j < 1000; j++ {
			start, end := uint64(rand.Intn(len(nums)+1)), uint64(rand.Intn(len(nums)+1))
			if start > end {
				start, end = end, start
			}
			x := uint64(rand.Intn(len(nums) + 1))
			if prefixBf(uint64(end), x) != wm.PrefixCount(uint64(end), x) {
				fmt.Println(prefixBf(uint64(end), x), wm.PrefixCount(uint64(end), x), end, x, nums)
				panic("prefixBf")
			}
			if rangeBf(uint64(start), uint64(end), x) != wm.RangeCount(uint64(start), uint64(end), x) {
				panic("rangeBf")
			}

			y := uint64(rand.Intn(len(nums) + 1))
			if res1, res2 := rangeFreqBf(uint64(start), uint64(end), x, y), wm.RangeFreq(uint64(start), uint64(end), x, y); res1 != res2 {
				fmt.Println(res1, res2, start, end, x, y)
				panic("rangeFreqBf")
			}

			k := uint64(rand.Intn(len(nums) + 1))
			if res1, res2 := kthBf(uint64(k), x), wm.Kth(uint64(k), x); res1 != res2 {
				fmt.Println(res1, res2, k, x, nums)
				panic("kthBf")
			}

			if res1, res2 := kthSmallestBf(uint64(start), uint64(end), uint64(k)), wm.KthSmallest(uint64(start), uint64(end), uint64(k)); res1 != res2 {
				fmt.Println(res1, res2, start, end, k)
				panic("kthSmallestBf")
			}

			if res1, res2 := kthSmallestIndexBf(uint64(start), uint64(end), uint64(k)), wm.KthSmallestIndex(uint64(start), uint64(end), uint64(k)); res1 != res2 {
				fmt.Println(res1, res2, start, end, k)
				panic("kthSmallestIndexBf")
			}

			if res1, res2 := kthLargestBf(uint64(start), uint64(end), uint64(k)), wm.KthLargest(uint64(start), uint64(end), uint64(k)); res1 != res2 {
				fmt.Println(res1, res2, start, end, k)
				panic("kthLargestBf")
			}

			c1, c2, c3 := countAllBf(uint64(start), uint64(end), x)
			c4, c5, c6 := wm.CountAll(uint64(start), uint64(end), x)
			if c1 != c4 || c2 != c5 || c3 != c6 {
				fmt.Println(c1, c2, c3, c4, c5, c6, start, end, x)
				panic("countAllBf")
			}

			funcs1 := []func(uint64, uint64, uint64) uint64{floorBf, lowerBf, ceilBf, higherBf}
			funcs2 := []func(uint64, uint64, uint64) uint64{wm.Floor, wm.Lower, wm.Ceil, wm.Higher}
			for i := 0; i < len(funcs1); i++ {
				res1 := funcs1[i](uint64(start), uint64(end), x)
				res2 := funcs2[i](uint64(start), uint64(end), x)
				if res1 != res2 {
					fmt.Println(res1, res2, start, end, x, nums)
					fmt.Println(runtime.FuncForPC(reflect.ValueOf(funcs2[i]).Pointer()).Name())
					panic("funcs")
				}
			}
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n := uint64(2e5)
	nums := make([]uint64, n)
	for i := uint64(0); i < n; i++ {
		nums[i] = uint64(rand.Intn(int(n)))
	}

	time1 := time.Now()
	wm := NewWaveletMatrixStaticFast(uint64(len(nums)), func(u uint64) uint64 { return nums[u] })
	fmt.Println("build time:", time.Since(time1))

	for i := uint64(0); i < n; i++ {
		// wm.PrefixCount(i, nums[i])
		// wm.RangeCount(0, i, nums[i])
		// wm.RangeFreq(0, i, nums[i], nums[i]+1)
		// wm.Kth(i, nums[i])
		// wm.KthSmallest(0, i, i)
		// wm.KthLargest(0, i, i)
		wm.Floor(0, i, nums[i])
		wm.Lower(0, i, nums[i])
		wm.Ceil(0, i, nums[i])
		wm.Higher(0, i, nums[i])
	}

	fmt.Println(time.Since(time1)) // 575.7913ms
}
