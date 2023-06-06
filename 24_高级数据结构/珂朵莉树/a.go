// 珂朵莉树(ODT)/Intervals

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

const INF int = 1e18

func demo() {
	odt := NewODT(10, -INF)
	odt.Set(0, 3, 1)
	odt.Set(3, 5, 2)
	fmt.Println(odt.Len, odt.Count, odt)
}

func UnionOfInterval() {
	// https://atcoder.jp/contests/abc256/tasks/abc256_d
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	odt := NewODT(2e5+10, -INF)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		odt.Set(l, r, 1)
	}
	odt.EnumerateAll(func(l, r int, x int) {
		if x == 1 {
			fmt.Fprintln(out, l, r)
		}
	})
}

func main() {
	UnionOfInterval()
}

type Value = int

type ODT struct {
	Len        int // 区间数
	Count      int // 区间元素个数之和
	llim, rlim int
	noneValue  Value
	data       []Value
	ss         *WAryTrie
}

// 指定区间长度 n 和哨兵 noneValue 建立一个 ODT.
//  区间为[0,n).
func NewODT(n int, noneValue Value) *ODT {
	res := &ODT{}
	dat := make([]Value, n)
	for i := 0; i < n; i++ {
		dat[i] = noneValue
	}
	ss := NewWAryTrie(n)
	ss.Insert(0)

	res.rlim = n
	res.noneValue = noneValue
	res.data = dat
	res.ss = ss
	return res
}

// 返回包含 x 的区间的信息.
func (odt *ODT) Get(x int, erase bool) (start, end int, value Value) {
	start, end = odt.ss.Prev(x), odt.ss.Next(x+1)
	value = odt.data[start]
	if erase && value != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
		odt.data[start] = odt.noneValue
		odt.mergeAt(start)
		odt.mergeAt(end)
	}
	return
}

func (odt *ODT) Set(start, end int, value Value) {
	odt.EnumerateRange(start, end, func(l, r int, x Value) {}, true)
	odt.ss.Insert(start)
	odt.data[start] = value
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}
	odt.mergeAt(start)
	odt.mergeAt(end)
}

func (odt *ODT) EnumerateAll(f func(start, end int, value Value)) {
	odt.EnumerateRange(0, odt.rlim, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *ODT) EnumerateRange(start, end int, f func(start, end int, value Value), erase bool) {
	if !(odt.llim <= start && start <= end && end <= odt.rlim) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue
	if !erase {
		l := odt.ss.Prev(start)
		for l < end {
			r := odt.ss.Next(l + 1)
			f(max(l, start), min(r, end), odt.data[l])
			l = r
		}
		return
	}

	// 分割区间
	p := odt.ss.Prev(start)
	if p < start {
		odt.ss.Insert(start)
		odt.data[start] = odt.data[p]
		if odt.data[start] != NONE {
			odt.Len++
		}
	}
	p = odt.ss.Next(end)
	if end < p {
		odt.data[end] = odt.data[odt.ss.Prev(end)]
		odt.ss.Insert(end)
		if odt.data[end] != NONE {
			odt.Len++
		}
	}
	p = start
	for p < end {
		q := odt.ss.Next(p + 1)
		x := odt.data[p]
		f(p, q, x)
		if odt.data[p] != NONE {
			odt.Len--
			odt.Count -= q - p
		}
		odt.ss.Discard(p)
		p = q
	}
	odt.ss.Insert(start)
	odt.data[start] = NONE
}

func (odt *ODT) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int, value Value) {
		var v interface{} = value
		if value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
}

func (odt *ODT) mergeAt(p int) {
	if p <= 0 || odt.rlim <= p {
		return
	}
	q := odt.ss.Prev(p - 1)
	if odt.data[p] == odt.data[q] {
		if odt.data[p] != odt.noneValue {
			odt.Len--
		}
		odt.ss.Discard(p)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
