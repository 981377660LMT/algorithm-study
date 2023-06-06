package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

func main() {
	set := NewWAryTrie(1e6)
	set.Insert(0)
	set.Insert(1)
	set.Insert(2)
	fmt.Println(set.Has(1))
	fmt.Println(set)
	fmt.Println(set.Max(), set.Prev(1025*32+1), set.Next(0))
	set.Discard(2)
	fmt.Println(set.Min(), set.Max())

}

// W叉Trie树.
type WAryTrie struct {
	n  int
	a1 []uint32
	a2 []uint32
	a3 []uint32
	a4 uint32
}

// 建立一个元素范围为[0,n)的W叉Trie树.
//  !n<2^20.
func NewWAryTrie(n int) *WAryTrie {
	return &WAryTrie{
		n:  n,
		a1: make([]uint32, (n>>5)+1),
		a2: make([]uint32, (n>>10)+1),
		a3: make([]uint32, (n>>15)+1),
	}
}

func (wat *WAryTrie) Has(x int) bool {
	return (wat.a1[x>>5]>>(x&31))&1 == 1
}

func (wat *WAryTrie) Insert(x int) {
	wat.a1[x>>5] |= 1 << (x & 31)
	wat.a2[x>>10] |= 1 << ((x >> 5) & 31)
	wat.a3[x>>15] |= 1 << ((x >> 10) & 31)
	wat.a4 |= 1 << (x >> 15)
}

// 返回是否成功删除(元素是否存在).
func (wat *WAryTrie) Discard(x int) (ok bool) {
	bit0 := uint32(1) << (x & 31)
	if wat.a1[x>>5]&bit0 == 0 {
		return
	}
	ok = true
	wat.a1[x>>5] -= bit0
	if wat.a1[x>>5] > 0 {
		return
	}
	bit1 := uint32(1) << ((x >> 5) & 31)
	wat.a2[x>>10] -= bit1
	if wat.a2[x>>10] > 0 {
		return
	}
	bit2 := uint32(1) << ((x >> 10) & 31)
	wat.a3[x>>15] -= bit2
	if wat.a3[x>>15] > 0 {
		return
	}
	wat.a4 -= uint32(1) << (x >> 15)
	return
}

// 返回集合中的最小值.如果不存在, 返回-1.
func (wat *WAryTrie) Min() int {
	if wat.a4 == 0 {
		return -1
	}
	x := wat._minBit(wat.a4)
	x = (x << 5) + wat._minBit(wat.a3[x])
	x = (x << 5) + wat._minBit(wat.a2[x])
	return (x << 5) + wat._minBit(wat.a1[x])
}

// 返回集合中的最大值.如果不存在, 返回n.
func (wat *WAryTrie) Max() int {
	if wat.a4 == 0 {
		return wat.n
	}
	x := wat._maxBit(wat.a4)
	x = (x << 5) + wat._maxBit(wat.a3[x])
	x = (x << 5) + wat._maxBit(wat.a2[x])
	return (x << 5) + wat._maxBit(wat.a1[x])
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (wat *WAryTrie) Prev(x int) int {
	if x < 0 {
		return -1
	}
	if x >= wat.n {
		x = wat.n - 1
	}
	if wat.Has(x) {
		return x
	}

	if tmp := wat._prevBit(wat.a1[x>>5], x); tmp != 0 {
		// 低 5 位设置为零
		return (x & 0xFFFFFFE0) + wat._maxBit(tmp)
	}
	x >>= 5
	if tmp := wat._prevBit(wat.a2[x>>5], x); tmp != 0 {
		x = (x & 0xFFFFFFE0) + wat._maxBit(tmp)
		return (x << 5) + wat._maxBit(wat.a1[x])
	}
	x >>= 5
	if tmp := wat._prevBit(wat.a3[x>>5], x); tmp != 0 {
		x = (x & 0xFFFFFFE0) + wat._maxBit(tmp)
		x = (x << 5) + wat._maxBit(wat.a2[x])
		return (x << 5) + wat._maxBit(wat.a1[x])
	}
	x >>= 5
	if tmp := wat._prevBit(wat.a4, x); tmp != 0 {
		x = wat._maxBit(tmp)
		x = (x << 5) + wat._maxBit(wat.a3[x])
		x = (x << 5) + wat._maxBit(wat.a2[x])
		return (x << 5) + wat._maxBit(wat.a1[x])
	}
	return -1
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (wat *WAryTrie) Next(x int) int {
	if x < 0 {
		x = 0
	}
	if x >= wat.n {
		return wat.n
	}
	if wat.Has(x) {
		return x
	}

	if a := wat.a1[x>>5]; wat._nextBit(a, x) > 1 {
		return x + 1 + wat._minBit(wat._nextBit(a, x+1))
	}
	x >>= 5
	if a := wat.a2[x>>5]; wat._nextBit(a, x) > 1 {
		x += 1 + wat._minBit(wat._nextBit(a, x+1))
		return (x << 5) + wat._minBit(wat.a1[x])
	}
	x >>= 5
	if a := wat.a3[x>>5]; wat._nextBit(a, x) > 1 {
		x += 1 + wat._minBit(wat._nextBit(a, x+1))
		x = (x << 5) + wat._minBit(wat.a2[x])
		return (x << 5) + wat._minBit(wat.a1[x])
	}
	x >>= 5
	if wat._nextBit(wat.a4, x) > 1 {
		x += 1 + wat._minBit(wat._nextBit(wat.a4, x+1))
		x = (x << 5) + wat._minBit(wat.a3[x])
		x = (x << 5) + wat._minBit(wat.a2[x])
		return (x << 5) + wat._minBit(wat.a1[x])
	}
	return wat.n
}

// 遍历[start,end)区间内的元素.
func (wat *WAryTrie) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = wat.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (wat *WAryTrie) String() string {
	res := []string{}
	wat.Enumerate(0, wat.n, func(i int) {
		res = append(res, strconv.Itoa(i))
	})
	return fmt.Sprintf("WAryTrie{%v}", strings.Join(res, ", "))
}

func (wat *WAryTrie) _maxBit(x uint32) int {
	return 31 - bits.LeadingZeros32(x)
}

func (wat *WAryTrie) _minBit(x uint32) int {
	return bits.TrailingZeros32(x)
}

func (wat *WAryTrie) _prevBit(x uint32, y int) uint32 {
	return x & (1<<(y&31) - 1)
}

func (wat *WAryTrie) _nextBit(x uint32, y int) uint32 {
	return x >> (y & 31)
}
