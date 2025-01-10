// api:
//  1. Set(i int32)
//  2. Get(i int32) int32
//  3. Build()
//  4. Count0(end int32) int32
//  5. Count1(end int32) int32
//  6. Count(end int32, v int32) int32
//  7. Kth0(k int32) int32
//  8. Kth1(k int32) int32
//  9. Kth(k int32, v int32) int32

package main

import (
	"fmt"
	"math/bits"
	"time"
)

func main() {
	testTime()
}

func demo() {
	bv := NewBitVector(10)
	bv.Set(3)
	bv.Set(8)
	bv.Build()
	println(bv.Get(3))
	fmt.Println(bv.Count0(4))
	fmt.Println(bv.Count1(4))
	fmt.Println(bv.Kth0(0))
	fmt.Println(bv.Kth1(0))
	fmt.Println(bv.Kth1(1))
	fmt.Println(bv.Kth1(2))
	fmt.Println(bv.GetAll())
}

type BitVector struct {
	n      int32
	size   int32
	bit    []uint64
	preSum []int32
}

func NewBitVector(n int32) *BitVector {
	size := (n + 63) >> 6
	bit := make([]uint64, size+1)
	preSum := make([]int32, size+1)
	return &BitVector{n: n, size: size, bit: bit, preSum: preSum}
}

func (bv *BitVector) Set(i int32) {
	bv.bit[i>>6] |= 1 << (i & 63)
}

func (bv *BitVector) Build() {
	for i := int32(0); i < bv.size; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bit[i]))
	}
}

func (bv *BitVector) Get(i int32) int32 {
	return int32(bv.bit[i>>6] >> (i & 63) & 1)
}

func (bv *BitVector) Count0(end int32) int32 {
	return end - (bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1))))
}

func (bv *BitVector) Count1(end int32) int32 {
	return bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1)))
}

func (bv *BitVector) Count(end int32, value int32) int32 {
	if value == 1 {
		return bv.Count1(end)
	}
	return end - bv.Count1(end)
}

func (bv *BitVector) Kth0(k int32) int32 {
	if k < 0 || bv.Count0(bv.n) <= k {
		return -1
	}
	l, r := int32(0), bv.size+1
	for r-l > 1 {
		m := (l + r) >> 1
		if m<<6-bv.preSum[m] > k {
			r = m
		} else {
			l = m
		}
	}
	indx := l << 6
	k -= (l<<6 - bv.preSum[l]) - bv.Count0(indx)
	l, r = indx, indx+64
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.Count0(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}

// k>=0.
func (bv *BitVector) Kth1(k int32) int32 {
	if k < 0 || bv.Count1(bv.n) <= k {
		return -1
	}
	l, r := int32(0), bv.size+1
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.preSum[m] > k {
			r = m
		} else {
			l = m
		}
	}
	indx := l << 6
	k -= bv.preSum[l] - bv.Count1(indx)
	l, r = indx, indx+64
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.Count1(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}

func (bv *BitVector) Kth(k int32, v int32) int32 {
	if v == 1 {
		return bv.Kth1(k)
	}
	return bv.Kth0(k)
}

func (bv *BitVector) GetAll() []int32 {
	res := make([]int32, 0, bv.n)
	for i := int32(0); i < bv.n; i++ {
		res = append(res, bv.Get(i))
	}
	return res
}

func testTime() {
	n := int32(1e7)
	bv := NewBitVector(n)
	for i := int32(0); i < n; i++ {
		if i%4 != 0 {
			bv.Set(i)
		}
	}

	time1 := time.Now()
	bv.Build()
	fmt.Println(time.Since(time1))

	for i := int32(0); i < n; i++ {
		bv.Get(i)
		bv.Count(i, 1)
		bv.Count(i, 0)
		bv.Kth(i, 1)
		bv.Kth(i, 0)
	}

	fmt.Println(time.Since(time1))
}
