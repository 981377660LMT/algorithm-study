// bitset子串匹配

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	// cf914f()
	cf963d()
}

// F. Substrings in a String
// https://www.luogu.com.cn/problem/CF914F
// 给定一个字符串s，有两种操作：
// 1 pos c，将s[pos]改为c
// 2 start end word，求字符串word在s[start:end]中出现的次数
//
// 由于字符种类少，因此可以 bitset 维护邻接表
// 对每次询问匹配t，处理出一个bitset.对每个t[i]，让 res &= indexes[t[i]]>>i 即可.
// 最后查询 [start,end-len(t)+1) 中合法的匹配起点个数。
func cf914f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)

	bytes := []byte(s)
	indexes := make([]*BitSetDynamic, 26) // 邻接表
	for i := range indexes {
		indexes[i] = NewBitsetDynamic(n, 0)
	}
	for i := 0; i < n; i++ {
		indexes[s[i]-'a'].Add(i)
	}

	validStarts := NewBitsetDynamic(n, 1)

	// 将 s[pos] 修改为 c.
	update := func(pos int, c byte) {
		if bytes[pos] == c {
			return
		}
		indexes[bytes[pos]-'a'].Discard(pos)
		indexes[c-'a'].Add(pos)
		bytes[pos] = c
	}

	// 查询 s[start:end) 中 word 的出现次数.
	query := func(start, end int, word []byte) int {
		if len(word) > end-start {
			return 0
		}

		validStarts.Fill(1)
		for i, c := range word {
			tmp := indexes[c-'a'].Copy()
			tmp.Rsh(i)
			validStarts.IAnd(tmp)
		}

		starts := validStarts.Slice(start, end-len(word)+1) // 合法的匹配起点
		return starts.OnesCount(0, starts.Size())
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var pos int
			fmt.Fscan(in, &pos)
			pos--
			var c string
			fmt.Fscan(in, &c)
			update(pos, c[0])
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			var word string
			fmt.Fscan(in, &word)
			fmt.Fprintln(out, query(start, end, []byte(word)))
		}
	}
}

// Frequency of String
// https://www.luogu.com.cn/problem/CF963D
// 给定一个字符串s和q个询问, 每个询问形如 (fi,wi) 表示次数和字符串.
// !要求找到s的一个最短的子串t使得wi在t中至少出现fi次.求t的最短长度.
func cf963d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var q int
	fmt.Fscan(in, &q)

	n := len(s)
	// 用于存储字符出现的位置
	indexes := make([]*BitSetDynamic, 26)
	for i := range indexes {
		indexes[i] = NewBitsetDynamic(n, 0)
	}
	for i := 0; i < n; i++ {
		indexes[s[i]-'a'].Add(i)
	}

	validStarts := NewBitsetDynamic(n, 1)
	kthValidStart := make([]int, n) // 第(i+1)个合法起点的位置(0-based)

	// 找到s的一个最短的子串t使得word在t中至少出现count次.
	// 不存在则返回-1.
	//
	// !求出合法起点，然后向后找即可.
	query := func(word []byte, count int) int {
		validStarts.Fill(1)
		for i, c := range word {
			tmp := indexes[c-'a'].Copy()
			tmp.Rsh(i)
			validStarts.IAnd(tmp)
		}
		ptr := 0
		validStarts.ForEach(func(pos int) bool {
			kthValidStart[ptr] = pos
			ptr++
			return false
		})
		res := n + 1
		for i := count - 1; i < ptr; i++ {
			res = min(res, kthValidStart[i]-kthValidStart[i-(count-1)])
		}
		if res == n+1 {
			return -1
		}
		return res + len(word)
	}

	for i := 0; i < q; i++ {
		var count int
		var word string
		fmt.Fscan(in, &count, &word)
		fmt.Fprintln(out, query([]byte(word), count))
	}
}

// 动态bitset，支持切片操作.
type BitSetDynamic struct {
	n    int
	data []uint64
}

// 建立一个大小为 n 的 bitset，初始值为 filledValue.
func NewBitsetDynamic(n int, filledValue int) *BitSetDynamic {
	if !(filledValue == 0 || filledValue == 1) {
		panic("filledValue should be 0 or 1")
	}
	data := make([]uint64, n>>6+1)
	if filledValue == 1 {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= (len(data) << 6) - n
		}
	}
	return &BitSetDynamic{n: n, data: data}
}

func (bs *BitSetDynamic) Add(i int) *BitSetDynamic {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BitSetDynamic) Has(i int) bool {
	return bs.data[i>>6]>>(i&63)&1 == 1
}

func (bs *BitSetDynamic) Discard(i int) {
	bs.data[i>>6] &^= 1 << (i & 63)
}

func (bs *BitSetDynamic) Flip(i int) {
	bs.data[i>>6] ^= 1 << (i & 63)
}

func (bs *BitSetDynamic) AddRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] |= maskL ^ maskR
		return
	}
	bs.data[i] |= maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = ^uint64(0)
	}
	bs.data[i] |= ^maskR
}

func (bs *BitSetDynamic) DiscardRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] &= ^maskL | maskR
		return
	}
	bs.data[i] &= ^maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = 0
	}
	bs.data[i] &= maskR
}

func (bs *BitSetDynamic) FlipRange(start, end int) {
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

// 左移 k 位 (<<k).
func (b *BitSetDynamic) Lsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
	if shift >= len(b.data) {
		for i := range b.data {
			b.data[i] = 0
		}
		return
	}

	if offset == 0 {
		copy(b.data[shift:], b.data)
	} else {
		for i := len(b.data) - 1; i > shift; i-- {
			b.data[i] = b.data[i-shift]<<offset | b.data[i-shift-1]>>(64-offset)
		}
		b.data[shift] = b.data[0] << offset
	}

	for i := 0; i < shift; i++ {
		b.data[i] = 0
	}
}

// 右移 k 位 (>>k).
func (b *BitSetDynamic) Rsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
	if shift >= len(b.data) {
		for i := range b.data {
			b.data[i] = 0
		}
		return
	}
	lim := len(b.data) - 1 - shift
	if offset == 0 {
		copy(b.data, b.data[shift:])
	} else {
		for i := 0; i < lim; i++ {
			b.data[i] = b.data[i+shift]>>offset | b.data[i+shift+1]<<(64-offset)
		}
		b.data[lim] = b.data[len(b.data)-1] >> offset
	}
	for i := lim + 1; i < len(b.data); i++ {
		b.data[i] = 0
	}
}

func (bs *BitSetDynamic) Fill(zeroOrOne int) {
	if zeroOrOne == 0 {
		for i := range bs.data {
			bs.data[i] = 0
		}
	} else {
		for i := range bs.data {
			bs.data[i] = ^uint64(0)
		}
		if bs.n != 0 {
			bs.data[len(bs.data)-1] >>= (len(bs.data) << 6) - bs.n
		}
	}
}

func (bs *BitSetDynamic) Clear() {
	for i := range bs.data {
		bs.data[i] = 0
	}
}

func (bs *BitSetDynamic) OnesCount(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > bs.n {
		end = bs.n
	}
	if start == 0 && end == bs.n {
		res := 0
		for _, v := range bs.data {
			res += bits.OnesCount64(v)
		}
		return res
	}
	pos1 := start >> 6
	pos2 := end >> 6
	if pos1 == pos2 {
		return bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)) & ((1 << (end & 63)) - 1))
	}
	count := 0
	if (start & 63) > 0 {
		count += bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)))
		pos1++
	}
	for i := pos1; i < pos2; i++ {
		count += bits.OnesCount64(bs.data[i])
	}
	if (end & 63) > 0 {
		count += bits.OnesCount64(bs.data[pos2] & ((1 << (end & 63)) - 1))
	}
	return count
}

func (bs *BitSetDynamic) AllOne(start, end int) bool {
	i := start >> 6
	if i == end>>6 {
		mask := ^uint64(0)<<(start&63) ^ ^uint64(0)<<(end&63)
		return (bs.data[i] & mask) == mask
	}
	mask := ^uint64(0) << (start & 63)
	if (bs.data[i] & mask) != mask {
		return false
	}
	for i++; i < end>>6; i++ {
		if bs.data[i] != ^uint64(0) {
			return false
		}
	}
	mask = ^uint64(0) << (end & 63)
	return ^(bs.data[end>>6] | mask) == 0
}

func (bs *BitSetDynamic) AllZero(start, end int) bool {
	i := start >> 6
	if i == end>>6 {
		mask := ^uint64(0)<<(start&63) ^ ^uint64(0)<<(end&63)
		return (bs.data[i] & mask) == 0
	}
	if (bs.data[i] >> (start & 63)) != 0 {
		return false
	}
	for i++; i < end>>6; i++ {
		if bs.data[i] != 0 {
			return false
		}
	}
	mask := ^uint64(0) << (end & 63)
	return (bs.data[end>>6] & ^mask) == 0
}

// 返回第一个 1 的下标，若不存在则返回-1.
func (bs *BitSetDynamic) IndexOfOne(position int) int {
	if position == 0 {
		for i, v := range bs.data {
			if v != 0 {
				return i<<6 | bs._lowbit(v)
			}
		}
		return -1
	}
	for i := position >> 6; i < len(bs.data); i++ {
		v := bs.data[i] & (^uint64(0) << (position & 63))
		if v != 0 {
			return i<<6 | bs._lowbit(v)
		}
		for i++; i < len(bs.data); i++ {
			if bs.data[i] != 0 {
				return i<<6 | bs._lowbit(bs.data[i])
			}
		}
	}
	return -1
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (bs *BitSetDynamic) IndexOfZero(position int) int {
	if position == 0 {
		for i, v := range bs.data {
			if v != ^uint64(0) {
				return i<<6 | bs._lowbit(^v)
			}
		}
		return -1
	}
	i := position >> 6
	if i < len(bs.data) {
		v := bs.data[i]
		if position&63 != 0 {
			v |= ^((^uint64(0)) << (position & 63))
		}
		if ^v != 0 {
			res := i<<6 | bs._lowbit(^v)
			if res < bs.n {
				return res
			}
			return -1
		}
		for i++; i < len(bs.data); i++ {
			if ^bs.data[i] != 0 {
				res := i<<6 | bs._lowbit(^bs.data[i])
				if res < bs.n {
					return res
				}
				return -1
			}
		}
	}
	return -1
}

// 返回右侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 n.
func (bs *BitSetDynamic) Next(index int) int {
	if index < 0 {
		index = 0
	}
	if index >= bs.n {
		return bs.n
	}
	k := index >> 6
	x := bs.data[k]
	s := index & 63
	x = (x >> s) << s
	if x != 0 {
		return (k << 6) | bs._lowbit(x)
	}
	for i := k + 1; i < len(bs.data); i++ {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._lowbit(bs.data[i])
	}
	return bs.n
}

// 返回左侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 -1.
func (bs *BitSetDynamic) Prev(index int) int {
	if index >= bs.n-1 {
		index = bs.n - 1
	}
	if index < 0 {
		return -1
	}
	k := index >> 6
	if (index & 63) < 63 {
		x := bs.data[k]
		x &= (1 << ((index & 63) + 1)) - 1
		if x != 0 {
			return (k << 6) | bs._topbit(x)
		}
		k--
	}
	for i := k; i >= 0; i-- {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._topbit(bs.data[i])
	}
	return -1
}

func (bs *BitSetDynamic) Equals(other *BitSetDynamic) bool {
	if len(bs.data) != len(other.data) {
		return false
	}
	for i := range bs.data {
		if bs.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IsSubset(other *BitSetDynamic) bool {
	if bs.n > other.n {
		return false
	}
	for i, v := range bs.data {
		if (v & other.data[i]) != v {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IsSuperset(other *BitSetDynamic) bool {
	if bs.n < other.n {
		return false
	}
	for i, v := range other.data {
		if (v & bs.data[i]) != v {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IOr(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] |= v
	}
	return bs
}

func (bs *BitSetDynamic) IAnd(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] &= v
	}
	return bs
}

func (bs *BitSetDynamic) IXor(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] ^= v
	}
	return bs
}

func (bs *BitSetDynamic) Or(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] | v
	}
	return res
}

func (bs *BitSetDynamic) And(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] & v
	}
	return res
}

func (bs *BitSetDynamic) Xor(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] ^ v
	}
	return res
}

func (bs *BitSetDynamic) IOrRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		bs.data[start>>6] |= other._get(a) << (start & 63)
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		bs.data[end>>6] |= other._get(b) << (end & 63)
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] |= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] |= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) IAndRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		if other._get(a) == 0 {
			bs.data[start>>6] &^= 1 << (start & 63)
		}
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		if other._get(b) == 0 {
			bs.data[end>>6] &^= 1 << (end & 63)
		}
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] &= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] &= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}

}

func (bs *BitSetDynamic) IXorRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		bs.data[start>>6] ^= other._get(a) << (start & 63)
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		bs.data[end>>6] ^= other._get(b) << (end & 63)
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] ^= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] ^= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

// 类似js中类型数组的set操作.如果超出赋值范围，抛出异常.
//
//	other: 要赋值的bitset.
//	offset: 赋值的起始元素下标.
func (bs *BitSetDynamic) Set(other *BitSetDynamic, offset int) {
	left, right := offset, offset+other.n
	if right > bs.n {
		panic("out of range")
	}
	a, b := 0, other.n
	for left < right && (left&63) != 0 {
		if other.Has(a) {
			bs.Add(left)
		} else {
			bs.Discard(left)
		}
		a++
		left++
	}
	for left < right && (right&63) != 0 {
		right--
		b--
		if other.Has(b) {
			bs.Add(right)
		} else {
			bs.Discard(right)
		}
	}

	// other[a:b] -> this[start:end]
	l, r := left>>6, right>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] = other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] = (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) Slice(start, end int) *BitSetDynamic {
	if start < 0 {
		start += bs.n
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += bs.n
	}
	if end > bs.n {
		end = bs.n
	}
	if start >= end {
		return NewBitsetDynamic(0, 0)
	}
	if start == 0 && end == bs.n {
		return bs.Copy()
	}

	res := NewBitsetDynamic(end-start, 0)
	remain := (end - start) & 63
	for i := 0; i < remain; i++ {
		if bs.Has(end - 1) {
			res.Add(end - start - 1)
		}
		end--
	}

	n := (end - start) >> 6
	hi := start & 63
	lo := 64 - hi
	s := start >> 6
	if hi == 0 {
		for i := 0; i < n; i++ {
			res.data[i] ^= bs.data[s+i]
		}
	} else {
		for i := 0; i < n; i++ {
			res.data[i] ^= (bs.data[s+i] >> hi) ^ (bs.data[s+i+1] << lo)
		}
	}

	return res
}

func (bs *BitSetDynamic) Copy() *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	copy(res.data, bs.data)
	return res
}

func (bs *BitSetDynamic) CopyAndResize(size int) *BitSetDynamic {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (1 << remainingBits) - 1
		newBits[len(newBits)-1] &= uint64(mask)
	}
	return &BitSetDynamic{data: newBits, n: size}
}

func (bs *BitSetDynamic) Resize(size int) {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (1 << remainingBits) - 1
		newBits[len(newBits)-1] &= uint64(mask)
	}
	bs.data = newBits
	bs.n = size
}

func (bs *BitSetDynamic) Expand(size int) {
	if size <= bs.n {
		return
	}
	bs.Resize(size)
}

func (bs *BitSetDynamic) BitLength() int {
	return bs._lastIndexOfOne() + 1
}

// 遍历所有 1 的位置.
func (bs *BitSetDynamic) ForEach(f func(pos int) (shouldBreak bool)) {
	for i, v := range bs.data {
		for ; v != 0; v &= v - 1 {
			j := (i << 6) | bs._lowbit(v)
			if f(j) {
				return
			}
		}
	}
}

func (bs *BitSetDynamic) Size() int {
	return bs.n
}

func (bs *BitSetDynamic) String() string {
	sb := strings.Builder{}
	sb.WriteString("BitSetDynamic{")
	nums := []string{}
	bs.ForEach(func(pos int) bool {
		nums = append(nums, fmt.Sprintf("%d", pos))
		return false
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 1, 2)
func (bs *BitSetDynamic) _topbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(x)
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 0, 2)
func (bs *BitSetDynamic) _lowbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return bits.TrailingZeros64(x)
}

func (bs *BitSetDynamic) _get(i int) uint64 {
	return bs.data[i>>6] >> (i & 63) & 1
}

func (bs *BitSetDynamic) _lastIndexOfOne() int {
	for i := len(bs.data) - 1; i >= 0; i-- {
		x := bs.data[i]
		if x != 0 {
			return (i << 6) | (bs._topbit(x))
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
