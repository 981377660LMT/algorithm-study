package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
	uf := NewPrevFinder(10)
	fmt.Println(uf)
	uf.Erase(2)
	uf.Erase(3)
	uf.Erase(4)
	fmt.Println(uf)
	fmt.Println(uf.Prev(0), uf.Prev(2), uf.Prev(100), uf.Prev(4))
}

// LinearSequenceUnionFind 线性序列并查集(PrevFinder).
type PrevFinder struct {
	n     int
	right []int
	data  []uint64
}

func NewPrevFinder(n int) *PrevFinder {
	len := (n >> 6) + 1
	f := &PrevFinder{
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

// 找到x左侧第一个未被访问过的位置(包含x).
//
//	如果不存在，返回-1.
func (f *PrevFinder) Prev(x int) int {
	if x < 0 {
		return -1
	}
	n := f.n
	if x >= n {
		x = n - 1
	}
	x = n - 1 - x
	div := x >> 6
	mod := x & 63
	mask := f.data[div] >> mod
	if mask != 0 {
		res := ((div << 6) | mod) + bits.TrailingZeros64(mask)
		if res < n {
			return n - 1 - res
		}
		return -1
	}
	div = f.findNext(div + 1)
	res := (div << 6) + bits.TrailingZeros64(f.data[div])
	if res < n {
		return n - 1 - res
	}
	return -1
}

// Erase 删除
func (f *PrevFinder) Erase(x int) {
	x = f.n - 1 - x
	div := x >> 6
	mod := x & 63
	if (f.data[div]>>mod)&1 != 0 { // flip
		f.data[div] ^= 1 << mod
	}
	if f.data[div] == 0 {
		f.right[div] = div + 1 // union to right
	}
}

func (f *PrevFinder) Has(x int) bool {
	if x < 0 || x >= f.n {
		return false
	}
	x = f.n - 1 - x
	return (f.data[x>>6]>>(x&63))&1 != 0
}

func (f *PrevFinder) String() string {
	sb := []string{}
	for i := 0; i < f.n; i++ {
		if f.Has(i) {
			sb = append(sb, fmt.Sprintf("%d", i))
		}
	}
	return "Finder(" + strings.Join(sb, ",") + ")"
}

func (f *PrevFinder) findNext(x int) int {
	for right := f.right[x]; right != x; {
		f.right[x] = f.right[right]
		x = f.right[x]
	}
	return x
}
