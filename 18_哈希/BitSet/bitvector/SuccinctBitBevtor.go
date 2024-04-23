// !deprecated

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

func main() {
	demo()
	test()
	testTime()
}

func demo() {
	bv := NewSuccinctBitVector(100)
	bv.Set(1, true)
	bv.Set(2, true)
	bv.Build()
	fmt.Println(bv.Count(3, false))
	fmt.Println(bv.Kth(0, true))
	fmt.Println(bv.Kth(1, true))
}

const NOT_FOUND = ^uint64(0)

type SuccinctBitVector struct {
	size  uint64
	ones  uint64
	large []uint64 // 大块
	small []uint16 // 小块
	bv    []uint16 // bitVector
}

func NewSuccinctBitVector(n uint64) *SuccinctBitVector {
	res := &SuccinctBitVector{size: n}
	res.bv = make([]uint16, (n+15)>>4+1)
	res.large = make([]uint64, n>>9+1)
	res.small = make([]uint16, n>>4+1)
	return res
}

func (sbv *SuccinctBitVector) Set(pos uint64, bit bool) {
	blockPos := pos >> 4
	offset := pos & 15
	if bit {
		sbv.bv[blockPos] |= 1 << offset
	} else {
		sbv.bv[blockPos] &= ^(1 << offset)
	}
}

func (sbv *SuccinctBitVector) Get(pos uint64) bool {
	blockPos := pos >> 4
	offset := pos & 15
	return (sbv.bv[blockPos]>>offset)&1 == 1
}

func (sbv *SuccinctBitVector) Build() {
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

func (sbv *SuccinctBitVector) Count(end uint64, bit bool) uint64 {
	if bit {
		return sbv.large[end>>9] + uint64(sbv.small[end>>4]) + uint64(bits.OnesCount16(sbv.bv[end>>4]&(1<<(end&15)-1)))
	}
	return end - sbv.Count(end, true)
}

// !kth 从0开始.
func (sbv *SuccinctBitVector) Kth(k uint64, bit bool) uint64 {
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

func (sbv *SuccinctBitVector) _selectInBlock(x uint64, rank uint64) uint64 {
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
	for i := 0; i < 100; i++ {
		n := uint64(rand.Intn(2e3) + 10)
		bits := make([]bool, n)
		for i := uint64(0); i < n; i++ {
			bits[i] = rand.Intn(2) == 1
		}

		get := func(index uint64) bool {
			return bits[index]
		}

		count := func(end uint64, bit bool) uint64 {
			res := uint64(0)
			for i := uint64(0); i < end; i++ {
				if get(i) == bit {
					res++
				}
			}
			return res
		}

		kth := func(k uint64, bit bool) uint64 {
			cnt := uint64(0)
			for i := uint64(0); i < n; i++ {
				if get(i) == bit {
					if cnt == k {
						return i
					}
					cnt++
				}
			}
			return NOT_FOUND
		}
		_ = kth

		sbv := NewSuccinctBitVector(n)
		for i := uint64(0); i < n; i++ {
			sbv.Set(i, bits[i])
		}
		sbv.Build()

		for i := uint64(0); i < n; i++ {
			if sbv.Get(i) != bits[i] {
				panic("Get Error")
			}
			if sbv.Count(i, true) != count(i, true) {
				fmt.Println(sbv.Count(i, true), count(i, true))
				panic("Count1 Error")
			}
			if sbv.Count(i, false) != count(i, false) {
				fmt.Println(sbv.Count(i, false), count(i, false))
				panic("Count0 Error")
			}
			if sbv.Kth(i, true) != kth(i, true) {
				fmt.Println(sbv.Kth(i, true), kth(i, true))
				panic("Kth1 Error")
			}
			if sbv.Kth(i, false) != kth(i, false) {
				fmt.Println(sbv.Kth(i, false), kth(i, false))
				panic("Kth0 Error")
			}
		}
	}

	fmt.Println("OK")
}

func testTime() {
	n := uint64(1e7)
	bv := NewSuccinctBitVector(n)
	for i := uint64(0); i < n; i++ {
		if i%4 != 0 {
			bv.Set(i, true)
		}
	}

	time1 := time.Now()
	bv.Build()
	fmt.Println(time.Since(time1))

	for i := uint64(0); i < n; i++ {
		bv.Get(i)
		bv.Count(i, true)
		bv.Count(i, false)
		bv.Kth(i, true)
		bv.Kth(i, false)
	}

	fmt.Println(time.Since(time1))
}
