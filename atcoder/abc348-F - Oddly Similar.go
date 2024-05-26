// F - Oddly Similar
// https://atcoder.jp/contests/abc348/tasks/abc348_f
// 给定n个长度为m的数组.
// 求有多少对数组满足，对应位置的元素相等的个数为奇数.
// 1 <= n, m <= 2000
// 1 <= a[i][j] <= 999
// O(nmk/64)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m uint16
	fmt.Fscan(in, &n, &m)
	M := make([][]uint16, n)
	for i := uint16(0); i < n; i++ {
		M[i] = make([]uint16, m)
	}
	for i := uint16(0); i < n; i++ {
		for j := uint16(0); j < m; j++ {
			fmt.Fscan(in, &M[i][j])
		}
	}

	similars := make([]Bitset, n)
	for i := range similars {
		similars[i] = NewBitset(int32(n))
	}

	for col := uint16(0); col < m; col++ {
		groups := make([]Bitset, 1000)
		for i := range groups {
			groups[i] = NewBitset(int32(n))
		}
		for row := uint16(0); row < n; row++ {
			v := M[row][col]
			similars[row].IXor(groups[v])
			groups[v].Set(int32(row))
		}
	}

	res := int32(0)
	for i := uint16(0); i < n; i++ {
		res += similars[i].OnesCount()
	}
	fmt.Fprintln(out, res)

}

type Bitset []uint

func NewBitset(n int32) Bitset { return make(Bitset, n>>6+1) } // (n+64-1)>>6

func (b Bitset) Has(p int32) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b Bitset) Flip(p int32)     { b[p>>6] ^= 1 << (p & 63) }
func (b Bitset) Set(p int32)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b Bitset) Reset(p int32)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b Bitset) Copy() Bitset {
	res := make(Bitset, len(b))
	copy(res, b)
	return res
}

func (bs Bitset) Clear() {
	for i := range bs {
		bs[i] = 0
	}
}

func (bs Bitset) OnesCount() int32 {
	res := int32(0)
	for i := range bs {
		res += int32(bits.OnesCount64(uint64(bs[i])))
	}
	return res
}

func (bs Bitset) IXor(other Bitset) {
	for i := range bs {
		bs[i] ^= other[i]
	}
}
