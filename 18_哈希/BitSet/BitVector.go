// WaveletMatrix 里用到的BitVector
// reference: https://tiramister.net/blog/posts/bitvector/
// 简洁数据结构

package main

import "math/bits"

type BitVector struct {
	n     int
	block []int
	sum   []int // block ごとに立っている 1 の数の累積和
}

func NewBitVector(n int) *BitVector {
	blockCount := (n + 63) >> 6
	return &BitVector{
		n:     n,
		block: make([]int, blockCount),
		sum:   make([]int, blockCount),
	}
}

func (f *BitVector) Set(i int) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *BitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

func (f *BitVector) Get(i int) int {
	return (f.block[i>>6] >> uint(i&63)) & 1
}

// 统计 [0,end) 中 value 的个数
func (f *BitVector) Count(value, end int) int {
	if value == 1 {
		return f.count1(end)
	}
	return end - f.count1(end)
}

// 统计 [0,end) 中第 k(0-indexed) 个 value 的位置
func (f *BitVector) Index(value, k int) int {
	if k < 0 || f.Count(value, f.n) <= k {
		return -1
	}

	left, right := 0, f.n
	for right-left > 1 {
		mid := (left + right) >> 1
		if f.Count(value, mid) >= k+1 {
			right = mid
		} else {
			left = mid
		}
	}
	return right - 1
}

func (f *BitVector) IndexWithStart(value, k, start int) int {
	return f.Index(value, k+f.Count(value, start))
}

func (f *BitVector) count1(k int) int {
	mask := (1 << uint(k&63)) - 1
	return f.sum[k>>6] + bits.OnesCount(uint(f.block[k>>6]&mask))
}
