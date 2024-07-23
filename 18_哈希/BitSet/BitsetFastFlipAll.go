// https://leetcode.cn/problems/design-bitset/description/

package main

import (
	"bytes"
	"math/bits"
)

// https://leetcode.cn/problems/design-bitset/description/
type Bitset struct {
	n  int
	bs *BitsetFastFlipAll
}

func Constructor(size int) Bitset {
	return Bitset{n: size, bs: NewBitsetFastFlipAll(size, false)}
}

func (this *Bitset) Fix(idx int) {
	this.bs.Add(idx)
}

func (this *Bitset) Unfix(idx int) {
	this.bs.Discard(idx)
}

func (this *Bitset) Flip() {
	this.bs.FlipAll()
}

func (this *Bitset) All() bool {
	return this.bs.OnesCount() == this.n
}

func (this *Bitset) One() bool {
	return this.bs.OnesCount() > 0
}

func (this *Bitset) Count() int {
	return this.bs.OnesCount()
}

func (this *Bitset) ToString() string {
	res := bytes.Repeat([]byte{'0'}, this.n)
	this.bs.ForEach(func(p int) { res[p] = '1' })
	return string(res)
}

// ----------------

type BitsetFastFlipAll struct {
	flip      bool
	n         int
	onesCount int
	data      []uint64
}

func NewBitsetFastFlipAll(n int, filled bool) *BitsetFastFlipAll {
	data := make([]uint64, n>>6+1)
	onesCount := 0
	if filled {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= (len(data) << 6) - n
		}
		onesCount = n
	}
	return &BitsetFastFlipAll{n: n, data: data, onesCount: onesCount}
}

func (b *BitsetFastFlipAll) FlipAll() {
	b.flip = !b.flip
	b.onesCount = b.n - b.onesCount
}

func (b *BitsetFastFlipAll) Add(i int) bool {
	if b.data[i>>6]>>(i&63)&1 == 1 != b.flip {
		return false
	}
	b.data[i>>6] ^= 1 << (i & 63)
	b.onesCount++
	return true
}

func (b *BitsetFastFlipAll) Discard(i int) bool {
	if b.data[i>>6]>>(i&63)&1 == 1 == b.flip {
		return false
	}
	b.data[i>>6] ^= 1 << (i & 63)
	b.onesCount--
	return true
}

func (b *BitsetFastFlipAll) Flip(i int) {
	if b.data[i>>6]>>(i&63)&1 == 1 == b.flip {
		b.data[i>>6] ^= 1 << (i & 63)
		b.onesCount++
	} else {
		b.data[i>>6] ^= 1 << (i & 63)
		b.onesCount--
	}
}

func (b *BitsetFastFlipAll) Has(i int) bool {
	return b.data[i>>6]>>(i&63)&1 == 1 != b.flip
}

func (b *BitsetFastFlipAll) OnesCount() int {
	return b.onesCount
}

func (bs *BitsetFastFlipAll) ForEach(f func(p int)) {
	if bs.flip {
		bs.flipData()
		bs.flip = false
	}
	for i, v := range bs.data {
		for ; v != 0; v &= v - 1 {
			f(i<<6 | bits.TrailingZeros64(v))
		}
	}
}

func (bs *BitsetFastFlipAll) flipData() {
	start, end := 0, bs.n
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] ^= maskL ^ maskR
		return
	}
	bs.data[i] ^= maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = ^bs.data[i]
	}
	bs.data[i] ^= ^maskR
}
