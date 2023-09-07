// WaveletMatrix 里用到的BitVector
// reference: https://tiramister.net/blog/posts/bitvector/
// 简洁数据结构

// 静态的bitset,添加元素需要build构建
package main

import "math/bits"

func main() {

}

type BV struct {
	n, zeros uint32
	block    []uint64
	count    []uint32
}

func NewBV(n int) *BV {
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

func (bv *BV) Get(i uint32) bool {
	return (bv.block[i/64]>>(i%64))&1 == 1
}

func (bv *BV) init(n int) {
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
