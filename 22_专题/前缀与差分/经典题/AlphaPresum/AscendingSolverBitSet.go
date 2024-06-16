package main

import (
	"cmp"
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
)

func countOfPeaks(nums []int, queries [][]int) []int {
	arr := NewAscendingSolverBitSet(
		int32(len(nums)), func(i int32) int32 { return int32(nums[i]) },
		func(a, b int32) bool { return a < b },
	)
	var res []int
	for _, q := range queries {
		if q[0] == 1 {
			start, end := int32(q[1]), int32(q[2])+1
			if end-start <= 1 {
				res = append(res, 0)
				continue
			}
			res = append(res, int(arr.CountPeek(start, end)))
		} else {
			arr.Set(int32(q[1]), int32(q[2]))
		}
	}
	return res
}

func main() {
	arr := []int{3, 2, 3}
	less := func(a, b int) bool { return a < b }
	solver := NewAscendingSolverBitSet(
		int32(len(arr)), func(i int32) int { return arr[i] },
		less,
	)
	// down在前一个, up在后一个
	fmt.Println(solver.down) // [0] => arr[0]>arr[1]
	fmt.Println(solver.up)   // [2] => arr[2]>arr[1]

	test()

	// nums = [4,1,4,2,1,5], queries = [[2,2,4],[1,0,2],[1,0,4]]
	fmt.Println(countOfPeaks([]int{4, 1, 4, 2, 1, 5}, [][]int{{2, 2, 4}, {1, 0, 2}, {1, 0, 4}}))
}

// BitSet维护区间递增/区间递减.
// 区间元素个数<=1，视为递增/递减.
type AscendingSolverBitSet[V cmp.Ordered] struct {
	n    int32
	arr  []V
	less func(a, b V) bool
	down *BitSetDynamic32 // down[i] = 1 表示 arr[i] > arr[i+1]
	up   *BitSetDynamic32 // up[i] = 1 表示 arr[i] > arr[i-1]
}

func NewAscendingSolverBitSet[V cmp.Ordered](
	n int32, f func(i int32) V, less func(a, b V) bool,
) *AscendingSolverBitSet[V] {
	arr := make([]V, n)
	for i := int32(0); i < n; i++ {
		arr[i] = f(i)
	}
	solver := &AscendingSolverBitSet[V]{n: n, arr: arr, less: less}
	down, up := NewBitsetDynamic32(n, 0), NewBitsetDynamic32(n, 0)
	for i := int32(1); i < n; i++ {
		if less(arr[i-1], arr[i]) {
			up.Add(i)
		}
		if less(arr[i], arr[i-1]) {
			down.Add(i)
		}
	}
	solver.down, solver.up = down, up
	return solver
}

func (solver *AscendingSolverBitSet[V]) Set(i int32, v V) {
	if solver.arr[i] == v {
		return
	}
	if i > 0 {
		if solver.less(v, solver.arr[i-1]) {
			solver.down.Add(i)
		} else {
			solver.down.Discard(i)
		}
		if solver.less(solver.arr[i-1], v) {
			solver.up.Add(i)
		} else {
			solver.up.Discard(i)
		}
	}
	if i+1 < solver.n {
		if solver.less(solver.arr[i+1], v) {
			solver.down.Add(i + 1)
		} else {
			solver.down.Discard(i + 1)
		}
		if solver.less(v, solver.arr[i+1]) {
			solver.up.Add(i + 1)
		} else {
			solver.up.Discard(i + 1)
		}
	}
	solver.arr[i] = v
}

func (solver *AscendingSolverBitSet[V]) Get(i int32) V { return solver.arr[i] }

// 区间元素个数<=1，视为递增.
func (solver *AscendingSolverBitSet[V]) IsAscending(start, end int32) bool {
	if start < 0 {
		start = 0
	}
	if end > solver.n {
		end = solver.n
	}
	if start >= end {
		return true
	}
	return solver.down.OnesCount(start+1, end) == 0
}

// 区间元素个数<=1，视为递减.
func (solver *AscendingSolverBitSet[V]) IsDescending(start, end int32) bool {
	if start < 0 {
		start = 0
	}
	if end > solver.n {
		end = solver.n
	}
	if start >= end {
		return true
	}
	return solver.up.OnesCount(start+1, end) == 0
}

// 求区间内的峰值个数(arr[i] > arr[i-1] && arr[i] > arr[i+1]).
func (solver *AscendingSolverBitSet[V]) CountPeek(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > solver.n {
		end = solver.n
	}
	a := solver.down.Slice(start+1, end)
	b := solver.up.Slice(start+1, end)
	b.Lsh(1)
	res := int32(0)
	for i := range a.data {
		res += int32(bits.OnesCount64(a.data[i] & b.data[i]))
	}
	return res
}

// 动态bitset，支持切片操作.
type BitSetDynamic32 struct {
	n    int32
	data []uint64
}

// 建立一个大小为 n 的 bitset，初始值为 filledValue.
// [0,n).
func NewBitsetDynamic32(n int32, filledValue int32) *BitSetDynamic32 {
	if !(filledValue == 0 || filledValue == 1) {
		panic("filledValue should be 0 or 1")
	}
	data := make([]uint64, n>>6+1)
	if filledValue == 1 {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= int32((len(data) << 6)) - n
		}
	}
	return &BitSetDynamic32{n: n, data: data}
}

func (bs *BitSetDynamic32) Add(i int32) *BitSetDynamic32 {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BitSetDynamic32) Has(i int32) bool {
	return bs.data[i>>6]>>(i&63)&1 == 1
}

func (bs *BitSetDynamic32) Discard(i int32) {
	bs.data[i>>6] &^= 1 << (i & 63)
}

func (bs *BitSetDynamic32) Flip(i int32) {
	bs.data[i>>6] ^= 1 << (i & 63)
}

func (bs *BitSetDynamic32) AddRange(start, end int32) {
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

func (bs *BitSetDynamic32) DiscardRange(start, end int32) {
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

func (bs *BitSetDynamic32) FlipRange(start, end int32) {
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
func (b *BitSetDynamic32) Lsh(k int32) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
	if shift >= int32(len(b.data)) {
		for i := range b.data {
			b.data[i] = 0
		}
		return
	}
	if offset == 0 {
		copy(b.data[shift:], b.data)
	} else {
		for i := int32(len(b.data)) - 1; i > shift; i-- {
			b.data[i] = b.data[i-shift]<<offset | b.data[i-shift-1]>>(64-offset)
		}
		b.data[shift] = b.data[0] << offset
	}
	for i := int32(0); i < shift; i++ {
		b.data[i] = 0
	}
}

// 右移 k 位 (>>k).
func (b *BitSetDynamic32) Rsh(k int32) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
	if shift >= int32(len(b.data)) {
		for i := range b.data {
			b.data[i] = 0
		}
		return
	}
	lim := int32(len(b.data)) - 1 - shift
	if offset == 0 {
		copy(b.data, b.data[shift:])
	} else {
		for i := int32(0); i < lim; i++ {
			b.data[i] = b.data[i+shift]>>offset | b.data[i+shift+1]<<(64-offset)
		}
		b.data[lim] = b.data[int32(len(b.data))-1] >> offset
	}
	for i := lim + 1; i < int32(len(b.data)); i++ {
		b.data[i] = 0
	}
}

func (bs *BitSetDynamic32) Fill(zeroOrOne int32) {
	if zeroOrOne == 0 {
		for i := range bs.data {
			bs.data[i] = 0
		}
	} else {
		for i := range bs.data {
			bs.data[i] = ^uint64(0)
		}
		if bs.n != 0 {
			bs.data[len(bs.data)-1] >>= int32((len(bs.data) << 6)) - bs.n
		}
	}
}

func (bs *BitSetDynamic32) Clear() {
	for i := range bs.data {
		bs.data[i] = 0
	}
}

func (bs *BitSetDynamic32) OnesCount(start, end int32) int32 {
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
		return int32(res)
	}
	pos1 := start >> 6
	pos2 := end >> 6
	if pos1 == pos2 {
		return int32(bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)) & ((1 << (end & 63)) - 1)))
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
	return int32(count)
}

func (bs *BitSetDynamic32) AllOne(start, end int32) bool {
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

func (bs *BitSetDynamic32) AllZero(start, end int32) bool {
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
func (bs *BitSetDynamic32) IndexOfOne(position int32) int32 {
	if position == 0 {
		for i, v := range bs.data {
			if v != 0 {
				return int32(i<<6) | bs._lowbit(v)
			}
		}
		return -1
	}
	for i := int32(position >> 6); i < int32(len(bs.data)); i++ {
		v := bs.data[i] & (^uint64(0) << (position & 63))
		if v != 0 {
			return int32(i<<6 | bs._lowbit(v))
		}
		for i++; i < int32(len(bs.data)); i++ {
			if bs.data[i] != 0 {
				return i<<6 | bs._lowbit(bs.data[i])
			}
		}
	}
	return -1
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (bs *BitSetDynamic32) IndexOfZero(position int32) int32 {
	if position == 0 {
		for i, v := range bs.data {
			if v != ^uint64(0) {
				return int32(i<<6) | bs._lowbit(^v)
			}
		}
		return -1
	}
	i := int32(position >> 6)
	if i < int32(len(bs.data)) {
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
		for i++; i < int32(len(bs.data)); i++ {
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
func (bs *BitSetDynamic32) Next(index int32) int32 {
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
	for i := k + 1; i < int32(len(bs.data)); i++ {
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
func (bs *BitSetDynamic32) Prev(index int32) int32 {
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

func (bs *BitSetDynamic32) Equals(other *BitSetDynamic32) bool {
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

func (bs *BitSetDynamic32) IsSubset(other *BitSetDynamic32) bool {
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

func (bs *BitSetDynamic32) IsSuperset(other *BitSetDynamic32) bool {
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

func (bs *BitSetDynamic32) IOr(other *BitSetDynamic32) *BitSetDynamic32 {
	for i, v := range other.data {
		bs.data[i] |= v
	}
	return bs
}

func (bs *BitSetDynamic32) IAnd(other *BitSetDynamic32) *BitSetDynamic32 {
	for i, v := range other.data {
		bs.data[i] &= v
	}
	return bs
}

func (bs *BitSetDynamic32) IXor(other *BitSetDynamic32) *BitSetDynamic32 {
	for i, v := range other.data {
		bs.data[i] ^= v
	}
	return bs
}

func (bs *BitSetDynamic32) Or(other *BitSetDynamic32) *BitSetDynamic32 {
	res := NewBitsetDynamic32(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] | v
	}
	return res
}

func (bs *BitSetDynamic32) And(other *BitSetDynamic32) *BitSetDynamic32 {
	res := NewBitsetDynamic32(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] & v
	}
	return res
}

func (bs *BitSetDynamic32) Xor(other *BitSetDynamic32) *BitSetDynamic32 {
	res := NewBitsetDynamic32(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] ^ v
	}
	return res
}

func (bs *BitSetDynamic32) IOrRange(start, end int32, other *BitSetDynamic32) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := int32(0), other.n
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
	l, r := int32(start>>6), int32(end>>6)
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := int32(0); i < n; i++ {
			bs.data[l+i] |= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := int32(0); i < n; i++ {
			bs.data[l+i] |= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic32) IAndRange(start, end int32, other *BitSetDynamic32) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := int32(0), other.n
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
		for i := int32(0); i < n; i++ {
			bs.data[l+i] &= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := int32(0); i < n; i++ {
			bs.data[l+i] &= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}

}

func (bs *BitSetDynamic32) IXorRange(start, end int32, other *BitSetDynamic32) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := int32(0), other.n
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
		for i := int32(0); i < n; i++ {
			bs.data[l+i] ^= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := int32(0); i < n; i++ {
			bs.data[l+i] ^= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

// 类似js中类型数组的set操作.如果超出赋值范围，抛出异常.
//
//	other: 要赋值的bitset.
//	offset: 赋值的起始元素下标.
func (bs *BitSetDynamic32) Set(other *BitSetDynamic32, offset int32) {
	left, right := offset, offset+other.n
	if right > bs.n {
		panic("out of range")
	}
	a, b := int32(0), other.n
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
		for i := int32(0); i < n; i++ {
			bs.data[l+i] = other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := int32(0); i < n; i++ {
			bs.data[l+i] = (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic32) Slice(start, end int32) *BitSetDynamic32 {
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
		return NewBitsetDynamic32(0, 0)
	}
	if start == 0 && end == bs.n {
		return bs.Copy()
	}

	res := NewBitsetDynamic32(end-start, 0)
	remain := (end - start) & 63
	for i := int32(0); i < remain; i++ {
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
		for i := int32(0); i < n; i++ {
			res.data[i] ^= bs.data[s+i]
		}
	} else {
		for i := int32(0); i < n; i++ {
			res.data[i] ^= (bs.data[s+i] >> hi) ^ (bs.data[s+i+1] << lo)
		}
	}

	return res
}

func (bs *BitSetDynamic32) Copy() *BitSetDynamic32 {
	res := NewBitsetDynamic32(bs.n, 0)
	copy(res.data, bs.data)
	return res
}

func (bs *BitSetDynamic32) CopyAndResize(size int32) *BitSetDynamic32 {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (1 << remainingBits) - 1
		newBits[len(newBits)-1] &= uint64(mask)
	}
	return &BitSetDynamic32{data: newBits, n: size}
}

func (bs *BitSetDynamic32) Resize(size int32) {
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

func (bs *BitSetDynamic32) Expand(size int32) {
	if size <= bs.n {
		return
	}
	bs.Resize(size)
}

func (bs *BitSetDynamic32) BitLength() int32 {
	return bs._lastIndexOfOne() + 1
}

// 遍历所有 1 的位置.
func (bs *BitSetDynamic32) ForEach(f func(pos int32) (shouldBreak bool)) {
	for i, v := range bs.data {
		for ; v != 0; v &= v - 1 {
			j := int32((i << 6)) | bs._lowbit(v)
			if f(j) {
				return
			}
		}
	}
}

func (bs *BitSetDynamic32) Size() int32 {
	return bs.n
}

func (bs *BitSetDynamic32) String() string {
	sb := strings.Builder{}
	sb.WriteString("BitSetDynamic{")
	nums := []string{}
	bs.ForEach(func(pos int32) bool {
		nums = append(nums, fmt.Sprintf("%d", pos))
		return false
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 1, 2)
func (bs *BitSetDynamic32) _topbit(x uint64) int32 {
	if x == 0 {
		return -1
	}
	return int32(63 - bits.LeadingZeros64(x))
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 0, 2)
func (bs *BitSetDynamic32) _lowbit(x uint64) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.TrailingZeros64(x))
}

func (bs *BitSetDynamic32) _get(i int32) uint64 {
	return bs.data[i>>6] >> (i & 63) & 1
}

func (bs *BitSetDynamic32) _lastIndexOfOne() int32 {
	for i := int32(len(bs.data)) - 1; i >= 0; i-- {
		x := bs.data[i]
		if x != 0 {
			return (i << 6) | (bs._topbit(x))
		}
	}
	return -1
}

func test() {
	for i := 0; i < 20; i++ {
		n := rand.Intn(500) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(100)
		}

		less := func(a, b int) bool { return a < b }
		solver := NewAscendingSolverBitSet(
			int32(n), func(i int32) int { return arr[i] },
			less,
		)

		isAscending := func(start, end int32) bool {
			for i := start + 1; i < end; i++ {
				if less(arr[i], arr[i-1]) {
					return false
				}
			}
			return true
		}
		_ = isAscending

		isDescending := func(start, end int32) bool {
			for i := start + 1; i < end; i++ {
				if less(arr[i-1], arr[i]) {
					return false
				}
			}
			return true
		}
		_ = isDescending

		for s := 0; s < 100; s++ {

			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					res1, res2 := solver.IsAscending(int32(i), int32(j)), isAscending(int32(i), int32(j))
					if res1 != res2 {
						fmt.Println("Error1", i, j, res1, res2)
						fmt.Println(arr)
						panic("error1")
					}
					res1, res2 = solver.IsDescending(int32(i), int32(j)), isDescending(int32(i), int32(j))
					if res1 != res2 {
						fmt.Println("Error2")
						panic("error2")
					}
				}
			}

			index := int32(rand.Intn(n))
			v := rand.Intn(10)
			solver.Set(index, v)
			arr[index] = v

			// get
			for i := 0; i < n; i++ {
				res1, res2 := solver.Get(int32(i)), arr[i]
				if res1 != res2 {
					fmt.Println("Error3", i, res1, res2)
					fmt.Println(arr)
					panic("error3")
				}
			}
		}
	}

	fmt.Println("pass")
}
