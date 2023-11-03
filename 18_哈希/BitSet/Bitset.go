// bitset 位集
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

// API:
// NewBitset(n) 创建一个长度为 n 的位集

// Set(p) 置 1
// Reset(p) 置 0
// Has(p) 判断 p 位是否为 1
// Flip(p) 翻转 p 位

// Foreach(f) 遍历所有 1 的位置

// Fill(v) 将所有位置为 v

// Index0() 返回第一个 0 的下标，若不存在则返回一个不小于 n 的位置
// Index1() 返回第一个 1 的下标，若不存在则返回一个不小于 n 的位置（同 C++ 中的 _Find_first）
// Next1(p) 返回下标 >= p 的第一个 1 的下标，若不存在则返回一个不小于 n 的位置（类似 C++ 中的 _Find_next，这里是 >=）
// Prev1(p) 返回下标 <= p 的第一个 1 的下标，若不存在则返回一个不大于 -1 的位置（类似 C++ 中的 _Find_prev，这里是 <=）
// Next0(p) 返回下标 >= p 的第一个 0 的下标，若不存在则返回一个不小于 n 的位置（类似 C++ 中的 _Find_next，这里是 >=）
// Prev0(p) 返回下标 <= p 的第一个 0 的下标，若不存在则返回一个不大于 -1 的位置（类似 C++ 中的 _Find_prev，这里是 <=）

// Count1() 返回 1 的个数
// Count0() 返回 0 的个数

// And() 位与
// IAnd() 位与，修改自身
// Or() 位或
// IOr() 位或，修改自身

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
	b := NewBitset(100)
	b.Set(10)
	b.Set(33)
	b.Set(76)
	b.Set(80)
	b.Rsh(17)
	fmt.Println(b)
}

type Bitset []uint

func NewBitset(n int) Bitset { return make(Bitset, n>>6+1) } // (n+64-1)>>6

func (b Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

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

// 遍历所有 1 的位置
// 如果对范围有要求，可在 f 中 return p < n
func (b Bitset) ForEach(f func(p int) (shouldBreak bool)) {
	for i, v := range b {
		for ; v != 0; v &= v - 1 {
			j := i<<6 | bits.TrailingZeros(v)
			if f(j) {
				return
			}
		}
	}
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (b Bitset) Index0() int {
	for i, v := range b {
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
	}
	return -1
}

// 返回第一个 1 的下标，若不存在则返回-1.
func (b Bitset) Index1() int {
	for i, v := range b {
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
	}
	return -1
}

// 返回下标 >= p 的第一个 1 的下标，若不存在则返回-1.
func (b Bitset) Next1(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i] & (^uint(0) << (p & 63)) // mask off bits below bound
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
		for i++; i < len(b); i++ {
			if b[i] != 0 {
				return i<<6 | bits.TrailingZeros(b[i])
			}
		}
	}
	return -1
}

// 返回下标 >= p 的第一个 0 的下标，若不存在则返回-1.
func (b Bitset) Next0(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i]
		if p&63 != 0 {
			v |= ^(^uint(0) << (p & 63))
		}
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
		for i++; i < len(b); i++ {
			if ^b[i] != 0 {
				return i<<6 | bits.TrailingZeros(^b[i])
			}
		}
	}
	return -1
}

// 返回最后第一个 1 的下标，若不存在则返回 -1
func (b Bitset) LastIndex1() int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i<<6 | (bits.Len(b[i]) - 1) // 如果再 +1，需要改成 i<<6 + bits.Len(b[i])
		}
	}
	return -1
}

// += 1 << i，模拟进位
func (b Bitset) Add(i int) { b.FlipRange(i, b.Next0(i)) }

// -= 1 << i，模拟借位
func (b Bitset) Sub(i int) { b.FlipRange(i, b.Next1(i)) }

// 判断 [l,r) 范围内的数是否全为 0
// https://codeforces.com/contest/1107/problem/D（标准做法是二维前缀和）
func (b Bitset) All0(l, r int) bool {
	i := l >> 6
	if i == r>>6 {
		mask := ^uint(0)<<(l&63) ^ ^uint(0)<<(r&63)
		return b[i]&mask == 0
	}
	if b[i]>>(l&63) != 0 {
		return false
	}
	for i++; i < r>>6; i++ {
		if b[i] != 0 {
			return false
		}
	}
	mask := ^uint(0) << (r & 63)
	return b[r>>6]&^mask == 0
}

// 判断 [l,r) 范围内的数是否全为 1
func (b Bitset) All1(l, r int) bool {
	i := l >> 6
	if i == r>>6 {
		mask := ^uint(0)<<(l&63) ^ ^uint(0)<<(r&63)
		return b[i]&mask == mask
	}
	mask := ^uint(0) << (l & 63)
	if b[i]&mask != mask {
		return false
	}
	for i++; i < r>>6; i++ {
		if ^b[i] != 0 {
			return false
		}
	}
	mask = ^uint(0) << (r & 63)
	return ^(b[r>>6] | mask) == 0
}

// 反转 [l,r) 范围内的比特
// https://codeforces.com/contest/1705/problem/E
func (b Bitset) FlipRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] ^= maskL ^ maskR
		return
	}
	b[i] ^= maskL
	for i++; i < r>>6; i++ {
		b[i] = ^b[i]
	}
	b[i] ^= ^maskR
}

// 将 [l,r) 范围内的比特全部置 1
func (b Bitset) SetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] |= maskL ^ maskR
		return
	}
	b[i] |= maskL
	for i++; i < r>>6; i++ {
		b[i] = ^uint(0)
	}
	b[i] |= ^maskR
}

// 将 [l,r) 范围内的比特全部置 0
func (b Bitset) ResetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] &= ^maskL | maskR
		return
	}
	b[i] &= ^maskL
	for i++; i < r>>6; i++ {
		b[i] = 0
	}
	b[i] &= maskR
}

// 左移 k 位
// LC1981 https://leetcode-cn.com/problems/minimize-the-difference-between-target-and-chosen-elements/
func (b Bitset) Lsh(k int) Bitset {
	if k == 0 {
		return b
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return b
	}
	if offset == 0 {
		// Fast path
		copy(b[shift:], b)
	} else {
		for i := len(b) - 1; i > shift; i-- {
			b[i] = b[i-shift]<<offset | b[i-shift-1]>>(64-offset)
		}
		b[shift] = b[0] << offset
	}
	for i := 0; i < shift; i++ {
		b[i] = 0
	}
	return b
}

// 右移 k 位
func (b Bitset) Rsh(k int) Bitset {
	if k == 0 {
		return b
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return b
	}
	lim := len(b) - 1 - shift
	if offset == 0 {
		// Fast path
		copy(b, b[shift:])
	} else {
		for i := 0; i < lim; i++ {
			b[i] = b[i+shift]>>offset | b[i+shift+1]<<(64-offset)
		}
		// 注意：若前后调用 lsh 和 rsh，需要注意超出 n 的范围的 1 对结果的影响（如果需要，可以把范围开大点）
		b[lim] = b[len(b)-1] >> offset
	}
	for i := lim + 1; i < len(b); i++ {
		b[i] = 0
	}
	return b
}

// 借用 bits 库中的一些方法的名字
func (b Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}
func (b Bitset) TrailingZeros() int { return b.Index1() }
func (b Bitset) Len() int           { return b.LastIndex1() + 1 }

// 下面几个方法均需保证长度相同
func (b Bitset) Equals(c Bitset) bool {
	for i, v := range b {
		if v != c[i] {
			return false
		}
	}
	return true
}

func (b Bitset) HasSubset(c Bitset) bool {
	for i, v := range b {
		if v|c[i] != v {
			return false
		}
	}
	return true
}

// 将 c 的元素合并进 b
func (b Bitset) IOr(c Bitset) Bitset {
	for i, v := range c {
		b[i] |= v
	}
	return b
}

func (b Bitset) Or(c Bitset) Bitset {
	res := make(Bitset, len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

func (b Bitset) IAnd(c Bitset) Bitset {
	for i, v := range c {
		b[i] &= v
	}
	return b
}

func (b Bitset) And(c Bitset) Bitset {
	res := make(Bitset, len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}

func (b Bitset) IXor(c Bitset) Bitset {
	for i, v := range c {
		b[i] ^= v
	}
	return b
}

func (b Bitset) Xor(c Bitset) Bitset {
	res := make(Bitset, len(b))
	for i, v := range b {
		res[i] = v ^ c[i]
	}
	return res
}

func (b Bitset) String() string {
	sb := strings.Builder{}
	sb.WriteString("BitSet{")
	nums := []string{}
	b.ForEach(func(pos int) bool {
		nums = append(nums, fmt.Sprintf("%d", pos))
		return false
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

// https://nyaannyaan.github.io/library/misc/bitset-find-prev.hpp
