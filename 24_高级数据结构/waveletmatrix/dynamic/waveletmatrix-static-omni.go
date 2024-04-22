// https://github.com/MitI-7/WaveletMatrix/blob/master/WaveletMatrix/WaveletMatrix.hpp
package main

import "math/bits"

func main() {

}

type WaveletMatrixStaticOmni struct {
}
type bitVector struct {
	n      int32
	size   int32
	bit    []uint64
	preSum []int32
}

func newBitVector(n int32) *bitVector {
	size := (n + 63) >> 6
	bit := make([]uint64, size+1)
	preSum := make([]int32, size+1)
	return &bitVector{n: n, size: size, bit: bit, preSum: preSum}
}

func (bv *bitVector) Set(i int32) {
	bv.bit[i>>6] |= 1 << (i & 63)
}

func (bv *bitVector) Build() {
	for i := int32(0); i < bv.size; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bit[i]))
	}
}

func (bv *bitVector) Get(i int32) bool {
	return bv.bit[i>>6]>>(i&63)&1 == 1
}

func (bv *bitVector) Count0(end int32) int32 {
	return end - (bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1))))
}

func (bv *bitVector) Count1(end int32) int32 {
	return bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1)))
}

func (bv *bitVector) Count(end int32, value int32) int32 {
	if value == 1 {
		return bv.Count1(end)
	}
	return end - bv.Count1(end)
}

func (bv *bitVector) Kth0(k int32) int32 {
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
func (bv *bitVector) Kth1(k int32) int32 {
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

func (bv *bitVector) Kth(k int32, v int32) int32 {
	if v == 1 {
		return bv.Kth1(k)
	}
	return bv.Kth0(k)
}
