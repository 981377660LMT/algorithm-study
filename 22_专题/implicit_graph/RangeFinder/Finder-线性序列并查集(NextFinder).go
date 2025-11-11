// https://maspypy.github.io/library/ds/decremental_fastset.hpp

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
	uf := NewNextFinder(10)
	fmt.Println(uf)
	uf.Erase(2)
	uf.Erase(3)
	uf.Erase(4)
	fmt.Println(uf)
	fmt.Println(uf.Next(0), uf.Next(2))
}

// LinearSequenceUnionFind 线性序列并查集(NextFinder).
type NextFinder struct {
	n     int
	right []int
	data  []uint64
}

func NewNextFinder(n int) *NextFinder {
	len := (n >> 6) + 1
	f := &NextFinder{
		n:     n,
		right: make([]int, len),
		data:  make([]uint64, len),
	}
	MASK := uint64(1<<64 - 1)
	for i := range f.right {
		f.right[i] = i
		f.data[i] = MASK
	}
	return f
}

// Next 下一个
//
//	如果不存在，返回n.
func (f *NextFinder) Next(x int) int {
	if x < 0 {
		x = 0
	}
	n := f.n
	if x >= n {
		return n
	}
	div := x >> 6
	mod := x & 63
	mask := f.data[div] >> mod
	if mask != 0 {
		return ((div << 6) | mod) + bits.TrailingZeros64(mask)
	}
	div = f.findNext(div + 1)
	return (div << 6) + bits.TrailingZeros64(f.data[div])
}

// Erase 删除
func (f *NextFinder) Erase(x int) {
	div := x >> 6
	mod := x & 63
	if (f.data[div]>>mod)&1 != 0 { // flip
		f.data[div] ^= 1 << mod
	}
	if f.data[div] == 0 {
		f.right[div] = div + 1 // union to right
	}
}

func (f *NextFinder) Has(x int) bool {
	if x < 0 || x >= f.n {
		return false
	}
	return (f.data[x>>6]>>(x&63))&1 != 0
}

func (f *NextFinder) String() string {
	sb := []string{}
	for i := 0; i < f.n; i++ {
		if f.Has(i) {
			sb = append(sb, fmt.Sprintf("%d", i))
		}
	}
	return "Finder(" + strings.Join(sb, ",") + ")"
}

func (f *NextFinder) findNext(x int) int {
	if f.right[x] == x {
		return x
	}
	f.right[x] = f.findNext(f.right[x])
	return f.right[x]
}
