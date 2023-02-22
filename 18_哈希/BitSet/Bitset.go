// bitset 位集
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e

// API:
// NewBitset(n) 创建一个长度为 n 的位集

// Set(p) 置 1
// Reset(p) 置 0
// Has(p) 判断 p 位是否为 1
// Flip(p) 翻转 p 位

// Foreach(f) 遍历所有 1 的位置

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

import "math/bits"

const _w = bits.UintSize     // 一个 uint 的位数
func NewBitset(n int) Bitset { return make(Bitset, n/_w+1) } // (n+_w-1)/_w

type Bitset []uint

func (b Bitset) Has(p int) bool { return b[p/_w]&(1<<(p%_w)) != 0 } // get
func (b Bitset) Flip(p int)     { b[p/_w] ^= 1 << (p % _w) }
func (b Bitset) Set(p int)      { b[p/_w] |= 1 << (p % _w) }  // 置 1
func (b Bitset) Reset(p int)    { b[p/_w] &^= 1 << (p % _w) } // 置 0

// 遍历所有 1 的位置
// 如果对范围有要求，可在 f 中 return p < n
func (b Bitset) Foreach(f func(p int) (Break bool)) {
	for i, v := range b {
		for ; v > 0; v &= v - 1 {
			j := i*_w | bits.TrailingZeros(v)
			if f(j) {
				return
			}
		}
	}
}

// 返回第一个 0 的下标，若不存在则返回一个不小于 n 的位置
func (b Bitset) Index0() int {
	for i, v := range b {
		if ^v != 0 {
			return i*_w | bits.TrailingZeros(^v)
		}
	}
	return len(b) * _w
}

// 返回第一个 1 的下标，若不存在则返回一个不小于 n 的位置（同 C++ 中的 _Find_first）
func (b Bitset) Index1() int {
	for i, v := range b {
		if v != 0 {
			return i*_w | bits.TrailingZeros(v)
		}
	}
	return len(b) * _w
}

// 返回下标 >= p 的第一个 1 的下标，若不存在则返回一个不小于 n 的位置（类似 C++ 中的 _Find_next，这里是 >=, C++里是 >）
func (b Bitset) Next1(p int) int {
	if i := p / _w; i < len(b) {
		v := b[i] & (^uint(0) << (p % _w)) // mask off bits below bound
		if v != 0 {
			return i*_w | bits.TrailingZeros(v)
		}
		for i++; i < len(b); i++ {
			if b[i] != 0 {
				return i*_w | bits.TrailingZeros(b[i])
			}
		}
	}
	return len(b) * _w
}

// 返回下标 >= p 的第一个 0 的下标，若不存在则返回一个不小于 n 的位置
func (b Bitset) Next0(p int) int {
	if i := p / _w; i < len(b) {
		v := b[i]
		if p%_w > 0 {
			v |= ^(^uint(0) << (p % _w))
		}
		if ^v != 0 {
			return i*_w | bits.TrailingZeros(^v)
		}
		for i++; i < len(b); i++ {
			if ^b[i] != 0 {
				return i*_w | bits.TrailingZeros(^b[i])
			}
		}
	}
	return len(b) * _w
}

// 返回最后第一个 1 的下标，若不存在则返回 -1
func (b Bitset) LastIndex1() int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i*_w | (bits.Len(b[i]) - 1) // 如果再 +1，需要改成 i*_w + bits.Len(b[i])
		}
	}
	return -1
}

// += 1 << i，模拟进位
func (b Bitset) Add(i int) { b.FlipRange(i, b.Next0(i)) }

// -= 1 << i，模拟借位
func (b Bitset) Sub(i int) { b.FlipRange(i, b.Next1(i)) }

// 判断 [l,r] 范围内的数是否全为 0
// https://codeforces.com/contest/1107/problem/D（标准做法是二维前缀和）
func (b Bitset) All0(l, r int) bool {
	i := l / _w
	if i == r/_w {
		mask := ^uint(0)<<(l%_w) ^ ^uint(0)<<(r%_w)
		return b[i]&mask == 0
	}
	if b[i]>>(l%_w) != 0 {
		return false
	}
	for i++; i < r/_w; i++ {
		if b[i] != 0 {
			return false
		}
	}
	mask := ^uint(0) << (r % _w)
	return b[r/_w]&^mask == 0
}

// 判断 [l,r] 范围内的数是否全为 1
func (b Bitset) All1(l, r int) bool {
	i := l / _w
	if i == r/_w {
		mask := ^uint(0)<<(l%_w) ^ ^uint(0)<<(r%_w)
		return b[i]&mask == mask
	}
	mask := ^uint(0) << (l % _w)
	if b[i]&mask != mask {
		return false
	}
	for i++; i < r/_w; i++ {
		if ^b[i] != 0 {
			return false
		}
	}
	mask = ^uint(0) << (r % _w)
	return ^(b[r/_w] | mask) == 0
}

// 反转 [l,r) 范围内的比特
// https://codeforces.com/contest/1705/problem/E
func (b Bitset) FlipRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l%_w), ^uint(0)<<(r%_w)
	i := l / _w
	if i == r/_w {
		b[i] ^= maskL ^ maskR
		return
	}
	b[i] ^= maskL
	for i++; i < r/_w; i++ {
		b[i] = ^b[i]
	}
	b[i] ^= ^maskR
}

// 将 [l,r) 范围内的比特全部置 1
func (b Bitset) SetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l%_w), ^uint(0)<<(r%_w)
	i := l / _w
	if i == r/_w {
		b[i] |= maskL ^ maskR
		return
	}
	b[i] |= maskL
	for i++; i < r/_w; i++ {
		b[i] = ^uint(0)
	}
	b[i] |= ^maskR
}

// 将 [l,r) 范围内的比特全部置 0
func (b Bitset) ResetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l%_w), ^uint(0)<<(r%_w)
	i := l / _w
	if i == r/_w {
		b[i] &= ^maskL | maskR
		return
	}
	b[i] &= ^maskL
	for i++; i < r/_w; i++ {
		b[i] = 0
	}
	b[i] &= maskR
}

// 左移 k 位
// LC1981 https://leetcode-cn.com/problems/minimize-the-difference-between-target-and-chosen-elements/
func (b Bitset) Lsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k/_w, k%_w
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return
	}
	if offset == 0 {
		// Fast path
		copy(b[shift:], b)
	} else {
		for i := len(b) - 1; i > shift; i-- {
			b[i] = b[i-shift]<<offset | b[i-shift-1]>>(_w-offset)
		}
		b[shift] = b[0] << offset
	}
	for i := 0; i < shift; i++ {
		b[i] = 0
	}
}

// 右移 k 位
func (b Bitset) Rsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k/_w, k%_w
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return
	}
	lim := len(b) - 1 - shift
	if offset == 0 {
		// Fast path
		copy(b, b[shift:])
	} else {
		for i := 0; i < lim; i++ {
			b[i] = b[i+shift]>>offset | b[i+shift+1]<<(_w-offset)
		}
		// 注意：若前后调用 lsh 和 rsh，需要注意超出 n 的范围的 1 对结果的影响（如果需要，可以把范围开大点）
		b[lim] = b[len(b)-1] >> offset
	}
	for i := lim + 1; i < len(b); i++ {
		b[i] = 0
	}
}

// 借用 bits 库中的一些方法的名字
func (b Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}
func (b Bitset) OnesCountRange(start, end int) int {
	pos1, pos2 := start/_w, end/_w
	if pos1 == pos2 {
		return bits.OnesCount(b[pos1] & (^uint(0) << (start % _w)) & ((1 << (end % _w)) - 1))
	}
	c := 0
	if start%_w > 0 {
		c += bits.OnesCount(b[pos1] & (^uint(0) << (start % _w)))
	}
	for i := pos1 + 1; i < pos2; i++ {
		c += bits.OnesCount(b[i])
	}
	if end%_w > 0 {
		c += bits.OnesCount(b[pos2] & ((1 << (end % _w)) - 1))
	}
	return c
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
func (b Bitset) IOr(c Bitset) {
	for i, v := range c {
		b[i] |= v
	}
}

func (b Bitset) Or(c Bitset) Bitset {
	res := NewBitset(len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

func (b Bitset) IAnd(c Bitset) {
	for i, v := range c {
		b[i] &= v
	}
}

func (b Bitset) And(c Bitset) Bitset {
	res := NewBitset(len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}
