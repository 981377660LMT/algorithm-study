// https://maspypy.github.io/library/graph/shortest_path/bfs_bitset.hpp
// bitset优化bfs O(V^2>>6)
// 应用:稠密图无权最短路 例如2000*2000的邻接矩阵
// 密グラフの重みなし最短路問題
// 01 行列を vc<bitset> の形で渡す
// O(N^2/w)
// 参考：(4000,4000) を 4000 回で 2 秒以内？

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1400
	// V<=2000,D<=1e18
	// 如果当前在第i行，Matrix[i][j]为1,就可以移动到第j列的任意一行
	// 问是否能做到:从任意一个点出发，经过D个回合后，可以到达任意一个点
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var V, D int
	fmt.Fscan(in, &V, &D)
	graph := make([]_Bitset, V+V)
	for i := range graph {
		graph[i] = _NewBitset(V + V)
	}

	for i := 0; i < V; i++ {
		var S string
		fmt.Fscan(in, &S)
		for j := 0; j < V; j++ {
			if S[j] == '0' {
				continue
			}
			graph[i].Set(j + V)
			graph[i+V].Set(j)
		}
	}
	for s := 0; s < V; s++ {
		dist := BfsBitset(graph, s)

		if D%2 == 0 {
			for t := 0; t < V; t++ {
				if dist[t] == -1 || dist[t] > D {
					fmt.Fprintln(out, "No")
					return
				}
			}
		}

		if D%2 == 1 {
			for t := V; t < V+V; t++ {
				if dist[t] == -1 || dist[t] > D {
					fmt.Fprintln(out, "No")
					return
				}
			}
		}
	}

	fmt.Fprintln(out, "Yes")
}

// O(n^2/w)
func BfsBitset(graph []_Bitset, start int) []int {
	n := len(graph)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}

	unused := _NewBitset(n)
	for i := 0; i < n; i++ {
		unused.Set(i)
	}
	queue := _NewBitset(n)
	queue.Set(start)

	d := 0
	for {
		p := queue.Index1()
		if p >= n {
			break
		}
		next := _NewBitset(n)
		for p < n {
			dist[p] = d
			unused.Reset(p)
			next.IOr(graph[p])
			p = queue.Next1(p + 1)
		}
		queue = next.And(unused)
		d++
	}

	return dist
}

type _Bitset []uint

func _NewBitset(n int) _Bitset { return make(_Bitset, n>>6+1) } // (n+64-1)>>6

func (b _Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b _Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b _Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b _Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b _Bitset) Copy() _Bitset {
	res := make(_Bitset, len(b))
	copy(res, b)
	return res
}

// 遍历所有 1 的位置
// 如果对范围有要求，可在 f 中 return p < n
func (b _Bitset) Foreach(f func(p int) (Break bool)) {
	for i, v := range b {
		for ; v > 0; v &= v - 1 {
			j := i<<6 | bits.TrailingZeros(v)
			if f(j) {
				return
			}
		}
	}
}

// 返回第一个 0 的下标，若不存在则返回一个不小于 n 的位置
func (b _Bitset) Index0() int {
	for i, v := range b {
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
	}
	return len(b) << 6
}

// 返回第一个 1 的下标，若不存在则返回一个不小于 n 的位置（同 C++ 中的 _Find_first）
func (b _Bitset) Index1() int {
	for i, v := range b {
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
	}
	return len(b) << 6
}

// 返回下标 >= p 的第一个 1 的下标，若不存在则返回一个不小于 n 的位置（类似 C++ 中的 _Find_next，这里是 >=, C++里是 >）
func (b _Bitset) Next1(p int) int {
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
	return len(b) << 6
}

// 返回下标 >= p 的第一个 0 的下标，若不存在则返回一个不小于 n 的位置
func (b _Bitset) Next0(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i]
		if p&63 > 0 {
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
	return len(b) << 6
}

// 返回最后第一个 1 的下标，若不存在则返回 -1
func (b _Bitset) LastIndex1() int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i<<6 | (bits.Len(b[i]) - 1) // 如果再 +1，需要改成 i<<6 + bits.Len(b[i])
		}
	}
	return -1
}

// += 1 << i，模拟进位
func (b _Bitset) Add(i int) { b.FlipRange(i, b.Next0(i)) }

// -= 1 << i，模拟借位
func (b _Bitset) Sub(i int) { b.FlipRange(i, b.Next1(i)) }

// 判断 [l,r] 范围内的数是否全为 0
// https://codeforces.com/contest/1107/problem/D（标准做法是二维前缀和）
func (b _Bitset) All0(l, r int) bool {
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

// 判断 [l,r] 范围内的数是否全为 1
func (b _Bitset) All1(l, r int) bool {
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
func (b _Bitset) FlipRange(l, r int) {
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
func (b _Bitset) SetRange(l, r int) {
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
func (b _Bitset) ResetRange(l, r int) {
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
func (b _Bitset) Lsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
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
			b[i] = b[i-shift]<<offset | b[i-shift-1]>>(64-offset)
		}
		b[shift] = b[0] << offset
	}
	for i := 0; i < shift; i++ {
		b[i] = 0
	}
}

// 右移 k 位
func (b _Bitset) Rsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
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
			b[i] = b[i+shift]>>offset | b[i+shift+1]<<(64-offset)
		}
		// 注意：若前后调用 lsh 和 rsh，需要注意超出 n 的范围的 1 对结果的影响（如果需要，可以把范围开大点）
		b[lim] = b[len(b)-1] >> offset
	}
	for i := lim + 1; i < len(b); i++ {
		b[i] = 0
	}
}

// 借用 bits 库中的一些方法的名字
func (b _Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}
func (b _Bitset) TrailingZeros() int { return b.Index1() }
func (b _Bitset) Len() int           { return b.LastIndex1() + 1 }

// 下面几个方法均需保证长度相同
func (b _Bitset) Equals(c _Bitset) bool {
	for i, v := range b {
		if v != c[i] {
			return false
		}
	}
	return true
}

func (b _Bitset) HasSubset(c _Bitset) bool {
	for i, v := range b {
		if v|c[i] != v {
			return false
		}
	}
	return true
}

// 将 c 的元素合并进 b
func (b _Bitset) IOr(c _Bitset) {
	for i, v := range c {
		b[i] |= v
	}
}

func (b _Bitset) Or(c _Bitset) _Bitset {
	res := make(_Bitset, len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

func (b _Bitset) IAnd(c _Bitset) {
	for i, v := range c {
		b[i] &= v
	}
}

func (b _Bitset) And(c _Bitset) _Bitset {
	res := make(_Bitset, len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}
