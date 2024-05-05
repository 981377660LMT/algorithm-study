package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	// G := NewNeighborMantainer(10, -1)
	// G.AddDirectedEdge(1, 2)
	// G.AddDirectedEdge(2, 1)
	// G.AddDirectedEdge(2, 3)
	// G.AddDirectedEdge(3, 2)
	// fmt.Println(G.Intersection(1, 3))
	abc350_g()
}

// [ABC350G] Mediator
// https://www.luogu.com.cn/problem/AT_abc350_g
// 初始时，有 N 个点，编号为 0 到 N-1,没有边存在.
// 有 Q 次操作，每次操作有三个整数 a, b, c.
// 1 u v: 在 u 和 v 之间添加一条边,保证 u 和 v 不在同一个连通块中.
// 2 u v: 询问是否存在和 u,v 都相邻的点，若存在输出编号，若不存在输出0.
func abc350_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	const MOD int = 998244353

	var n, q int
	fmt.Fscan(in, &n, &q)

	G := NewNeighborMantainer(int32(n), -1)

	// 添加边.
	addEdge := func(a, b int32) {
		G.AddDirectedEdge(a, b)
		G.AddDirectedEdge(b, a)
	}

	// 查询a和b是否有共同的邻居.
	query := func(a, b int32) (int32, bool) {
		intersection := G.Intersection(a, b)
		if intersection == -1 {
			return 0, false
		}
		return intersection, true
	}

	preRes := 0
	for i := 0; i < q; i++ {
		var A, B, C int
		fmt.Fscan(in, &A, &B, &C)
		A = 1 + (((A * (1 + preRes)) % MOD) % 2)
		B = 1 + (((B * (1 + preRes)) % MOD) % n)
		C = 1 + (((C * (1 + preRes)) % MOD) % n)
		B--
		C--
		if A == 1 {
			addEdge(int32(B), int32(C))
		} else {
			res, ok := query(int32(B), int32(C))
			if !ok {
				preRes = 0
				fmt.Fprintln(out, preRes)
			} else {
				preRes = int(res + 1)
				fmt.Fprintln(out, preRes)
			}
		}
	}
}

// 支持快速查询邻居的交集元素.
type NeighborMantainer struct {
	n         int32
	adjMatrix []*BitSetDynamic32   // 邻接矩阵
	adjSet    []map[int32]struct{} // 邻接表
	threshold int32
}

// threshold: 邻居个数到达 threshold 时，邻接表转换为邻接矩阵. 负数表示使用默认值150.
func NewNeighborMantainer(n int32, threshold int32) *NeighborMantainer {
	if threshold < 0 {
		threshold = 150
	}
	adjMatrix := make([]*BitSetDynamic32, n)
	adjSet := make([]map[int32]struct{}, n)
	for i := range adjSet {
		adjSet[i] = make(map[int32]struct{})
	}
	return &NeighborMantainer{n: n, adjMatrix: adjMatrix, adjSet: adjSet, threshold: threshold}
}

// 添加有向边 from->to.
func (nm *NeighborMantainer) AddDirectedEdge(from, to int32) {
	nm.adjSet[from][to] = struct{}{}
	if int32(len(nm.adjSet[from])) == nm.threshold && nm.adjMatrix[from] == nil {
		nm.adjMatrix[from] = NewBitsetDynamic32(nm.n, 0)
		for to := range nm.adjSet[from] {
			nm.adjMatrix[from].Add(to)
		}
	}
	if nm.adjMatrix[from] != nil {
		nm.adjMatrix[from].Add(to)
	}
}

// 求a和b的一个公共邻居.如果不存在返回-1.
func (nm *NeighborMantainer) Intersection(a, b int32) int32 {
	if len(nm.adjSet[a]) > len(nm.adjSet[b]) {
		a, b = b, a
	}

	b1, b2 := nm.adjMatrix[a], nm.adjMatrix[b]
	// a大b大
	if b1 != nil && b2 != nil {
		return firstIntersection(b1, b2)
	}
	// a小b大
	if b1 == nil && b2 != nil {
		for v := range nm.adjSet[a] {
			if b2.Has(v) {
				return v
			}
		}
		return -1
	}

	// a小b小
	if b1 == nil && b2 == nil {
		s1, s2 := nm.adjSet[a], nm.adjSet[b]
		for v := range s1 {
			if _, ok := s2[v]; ok {
				return v
			}
		}
		return -1
	}
	return -1
}

func firstIntersection(a, b *BitSetDynamic32) int32 {
	// for i, v := range a.data {
	// 	if v&b.data[i] != 0 {
	// 		return int32(i<<6 | bits.TrailingZeros64(v&b.data[i]))
	// 	}
	// }
	// return -1
	tmp := a.And(b)
	return tmp.IndexOfOne(0)
}

// 动态bitset，支持切片操作.
type BitSetDynamic32 struct {
	n    int32
	data []uint64
}

// 建立一个大小为 n 的 bitset，初始值为 filledValue.
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
