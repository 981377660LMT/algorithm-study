// WaveletMatrix 里用到的BitVector
// reference: https://tiramister.net/blog/posts/bitvector/
// 简洁数据结构

// 静态的bitset,添加元素需要build构建

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	bv := NewBitVector(65)
	for i := 0; i < 65; i++ {
		bv.Set(i)
	}
	bv.Build()
	fmt.Println(bv.Count(0, 65))
}

type BitVector struct {
	n     int
	block []int
	sum   []int
}

func NewBitVector(n int) *BitVector {
	blockCount := n>>6 + 1
	return &BitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
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

// 第 k(0-indexed) 个 value 的位置
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
