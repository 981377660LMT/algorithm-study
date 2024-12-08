// 下标从0开始
// 1.BITArray: 单点修改, 区间查询
// 2.BITMap: 单点修改, 区间查询
// 3.BITRangeAddPointGetArray: 区间修改, 单点查询(差分)
// 4.BITRangeAddPointGetMap: 区间修改, 单点查询(差分)
// 5.BITRangeAddRangeSumArray: 区间修改, 区间查询
// 6.BITRangeAddRangeSumMap: 区间修改, 区间查询
// 7.BITPrefixArray: 单点修改, 前缀查询
// 8.BITPrefixMap: 单点修改, 前缀查询

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	yosupo()
	// demo()
	// test01()
	// CF1288E()
	// test01()
}

// https://judge.yosupo.jp/problem/predecessor_problem
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	var t string
	fmt.Fscan(in, &t)
	bit01 := NewBITArray01From(n, func(i int) bool { return t[i] == '1' })
	for i := 0; i < q; i++ {
		var op, k int
		fmt.Fscan(in, &op, &k)
		if op == 0 {
			bit01.Add(k)
		} else if op == 1 {
			bit01.Remove(k)
		} else if op == 2 {
			ok := bit01.Has(k)
			if ok {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		} else if op == 3 {
			// n := bit01.QueryPrefix(k)
			// if n == bit01.QueryAll() {
			// 	fmt.Fprintln(out, -1)
			// } else {
			// 	fmt.Fprintln(out, bit01.Kth(n, 0))
			// }
			fmt.Fprintln(out, bit01.Next(k))
		} else if op == 4 {
			// n := bit01.QueryPrefix(k + 1)
			// if n == 0 {
			// 	fmt.Fprintln(out, -1)
			// } else {
			// 	fmt.Fprintln(out, bit01.Kth(n-1, 0))
			// }
			fmt.Fprintln(out, bit01.Prev(k))
		}
	}
}

func test01() {
	bit01 := NewBITArray01(10)
	fmt.Println(bit01)
	bit01.Add(1)
	bit01.Add(2)
	bit01.Add(4)
	fmt.Println(bit01.Kth(0, 0))
	fmt.Println(bit01.Kth(1, 0))
	fmt.Println(bit01.Kth(2, 0))
	fmt.Println(bit01.Kth(3, 0))
}

// https://leetcode.cn/problems/longest-uploaded-prefix/description/
type LUPrefix struct {
	bit *BITArray
}

func Constructor(n int) LUPrefix {
	return LUPrefix{bit: NewBitArray(n + 1)}
}

func (this *LUPrefix) Upload(video int) {
	this.bit.Add(video-1, 1)
}

func (this *LUPrefix) Longest() int {
	return this.bit.MaxRightWithIndex(0, func(index int, sum int) bool { return sum == index })
}

// https://www.luogu.com.cn/problem/CF1288E
// 给定一个数组，元素为0-n-1，不断进行把一个数组中的一个元素插入数组开头，求出每个元素出现位置的最小值和最大值。
//
// n+q 个位置，每个位置0/1表示这个位置是否有数字.
// !虽然无法很方便的用数组来模拟每个数字向前移动的情况，但是我们可以很方便的用数组来模拟每个数字添加到尾部的操作
// !因为难以整体移动数组，因此只能维护每个数出现的位置
func CF1288E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	moved := make([]int, q)
	for i := range moved {
		fmt.Fscan(in, &moved[i])
		moved[i]--
	}

	res := make([][2]int, n)
	for i := 0; i < n; i++ {
		res[i] = [2]int{i + 1, -1}
	}

	pos := make([]int, n) // 记录每个元素的位置
	for i := 0; i < n; i++ {
		pos[i] = n - 1 - i
	}
	bit := NewBITArray01From(
		n+q,
		func(i int) bool {
			return i < n
		},
	)

	// n-1, n-2, ..., 0, null ,null, null, ..., null 一共n+q个位置
	for to := n; to < n+q; to++ {
		v := moved[to-n]
		from := pos[v]
		res[v][1] = max(res[v][1], bit.QueryRange(from, n+q)) // 移动前的最大值
		pos[v] = to
		bit.Remove(from)
		bit.Add(to)
		res[v][0] = min(res[v][0], 1) // 移动后的最小值
	}
	for i := 0; i < n; i++ {
		res[i][1] = max(res[i][1], bit.QueryRange(pos[i], n+q)) // 更新最大值
	}

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i][0], res[i][1])
	}
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(start int, predicate func(sum int) bool) int {
	s := 0
	i := start
	k := func() int {
		for {
			if i&1 == 1 {
				s -= b.data[i-1]
				i--
			}
			if i == 0 {
				return bits.Len32(uint32(b.n))
			}
			k := bits.TrailingZeros32(uint32(i)) - 1
			if i+(1<<k) > b.n {
				return k
			}
			t := s + b.data[i+(1<<k)-1]
			if !predicate(t) {
				return k
			}
			s -= b.data[i-1]
			i -= i & -i
		}
	}()
	for k > 0 {
		k--
		if i+(1<<k)-1 < b.n {
			t := s + b.data[i+(1<<k)-1]
			if predicate(t) {
				i += 1 << k
				s = t
			}
		}
	}
	return i
}

// MaxRightWithIndex
func (b *BITArray) MaxRightWithIndex(start int, predicate func(index int, sum int) bool) int {
	s := 0
	i := start
	k := func() int {
		for {
			if i&1 == 1 {
				s -= b.data[i-1]
				i--
			}
			if i == 0 {
				return bits.Len32(uint32(b.n))
			}
			k := bits.TrailingZeros32(uint32(i)) - 1
			if i+(1<<k) > b.n {
				return k
			}
			t := s + b.data[i+(1<<k)-1]
			if !predicate(i+(1<<k), t) {
				return k
			}
			s -= b.data[i-1]
			i -= i & -i
		}
	}()
	for k > 0 {
		k--
		if i+(1<<k)-1 < b.n {
			t := s + b.data[i+(1<<k)-1]
			if predicate(i+(1<<k), t) {
				i += 1 << k
				s = t
			}
		}
	}
	return i
}

// verify: https://judge.yosupo.jp/submission/254655
func (b *BITArray) MinLeft(end int, predicate func(sum int) bool) int {
	s := 0
	i := end
	k := 0
	for i > 0 && predicate(s) {
		s += b.data[i-1]
		k = bits.TrailingZeros32(uint32(i))
		i -= i & -i
	}
	if predicate(s) {
		return 0
	}
	for k > 0 {
		k--
		t := s - b.data[i+(1<<k)-1]
		if !predicate(t) {
			i += 1 << k
			s = t
		}
	}
	return i + 1
}

func (b *BITArray) Kth(k int, start int) int {
	return b.MaxRight(start, func(x int) bool { return x <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

// 01树状数组.
type BITArray01 struct {
	n    int
	size int // data、bit的长度
	data []uint64
	bit  *BITArray
}

func NewBITArray01(n int) *BITArray01 {
	return NewBITArray01From(n, func(index int) bool { return false })
}

func NewBITArray01From(n int, f func(index int) bool) *BITArray01 {
	size := n>>6 + 1
	data := make([]uint64, size)
	for i := 0; i < n; i++ {
		if f(i) {
			data[i>>6] |= 1 << (i & 63)
		}
	}
	bit := NewBitArrayFrom(size, func(i int) int { return bits.OnesCount64(data[i]) })
	return &BITArray01{n: n, size: size, data: data, bit: bit}
}

func (bit01 *BITArray01) QueryAll() int {
	return bit01.bit.QueryAll()
}

func (bit01 *BITArray01) QueryPrefix(end int) int {
	i, j := end>>6, end&63
	res := bit01.bit.QueryPrefix(i)
	res += bits.OnesCount64(bit01.data[i] & ((1 << j) - 1))
	return res
}

func (bit01 *BITArray01) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > bit01.n {
		end = bit01.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return bit01.QueryPrefix(end)
	}
	res := 0
	res -= bits.OnesCount64(bit01.data[start>>6] & ((1 << (start & 63)) - 1))
	res += bits.OnesCount64(bit01.data[end>>6] & ((1 << (end & 63)) - 1))
	res += bit01.bit.QueryRange(start>>6, end>>6)
	return res
}

func (bit01 *BITArray01) Add(index int) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 1 {
		return false
	}
	bit01.data[i] |= 1 << j
	bit01.bit.Add(i, 1)
	return true
}

func (bit01 *BITArray01) Remove(index int) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 0 {
		return false
	}
	bit01.data[i] ^= 1 << j
	bit01.bit.Add(i, -1)
	return true
}

func (bit01 *BITArray01) Has(index int) bool {
	i, j := index>>6, index&63
	return (bit01.data[i]>>j)&1 == 1
}

// 0<=k<bit01.QueryAll().
// 如果不存在，返回 -1.
func (bit01 *BITArray01) Kth(k int, start int) int {
	if k >= bit01.QueryAll() {
		return -1
	}
	k += bits.OnesCount64(bit01.data[start>>6] & ((1 << (start & 63)) - 1))
	start >>= 6
	mid := 0
	check := func(preSum int) bool {
		if preSum <= k {
			if preSum > mid {
				mid = preSum
			}
		}
		return preSum <= k
	}
	pos := bit01.bit.MaxRight(start, check)
	if pos == bit01.n {
		return -1
	}
	k -= mid
	x := bit01.data[pos]
	p := bits.OnesCount64(x)
	if p <= k {
		return -1
	}
	k = sort.Search(64, func(n int) bool { return (p - bits.OnesCount64(x>>(n+1))) > k })
	return pos<<6 | k
}

// 如果不存在，返回 -1.
func (bit01 *BITArray01) Next(k int) int {
	if k < 0 {
		k = 0
	}
	if k >= bit01.n {
		return -1
	}
	pos := k >> 6
	k &= 63
	x := bit01.data[pos] & ^((1 << k) - 1)
	if x != 0 {
		return pos<<6 | bits.TrailingZeros64(x)
	}
	pos = bit01.bit.Kth(0, pos+1)
	if pos == bit01.size || bit01.data[pos] == 0 {
		return -1
	}
	return pos<<6 | bits.TrailingZeros64(bit01.data[pos])
}

// 如果不存在，返回 -1.
func (bit01 *BITArray01) Prev(k int) int {
	if k >= bit01.n {
		k = bit01.n - 1
	}
	if k < 0 {
		return -1
	}
	pos := k >> 6
	k &= 63
	x := bit01.data[pos]
	if k < 63 {
		x &= (1 << (k + 1)) - 1
	}
	if x != 0 {
		return pos<<6 | (bits.Len64(x) - 1)
	}
	pos = bit01.bit.MinLeft(pos, func(sum int) bool { return sum <= 0 }) - 1
	if pos == -1 {
		return -1
	}
	return pos<<6 | (bits.Len64(bit01.data[pos]) - 1)
}

func (bit01 *BITArray01) String() string {
	res := []string{}
	for i := 0; i < bit01.n; i++ {
		if bit01.QueryRange(i, i+1) == 1 {
			res = append(res, fmt.Sprintf("%d", i))
		}
	}
	return fmt.Sprintf("BITArray01: [%v]", strings.Join(res, ", "))
}

// !Point Add Range Sum, 0-based.
type BITMap struct {
	n     int
	total int
	data  map[int]int
}

func NewBITMap(n int) *BITMap {
	return &BITMap{n: n + 1, data: make(map[int]int)}
}

func (b *BITMap) Add(i, v int) {
	b.total += v
	for i++; i <= b.n; i += i & -i {
		b.data[i-1] += v
	}
}

// [0, end)
func (b *BITMap) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

func (b *BITMap) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITMap) QueryAll() int {
	return b.total
}

func (b *BITMap) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1 树状数组查找第 k(0-based) 个1的位置.
// UpperBound.
func (b *BITMap) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

// 基于差分实现，区间修改，单点查询.
type BITRangeAddPointGetArray struct {
	bit *BITArray
}

func NewBITRangeAddPointGet(n int) *BITRangeAddPointGetArray {
	return &BITRangeAddPointGetArray{bit: NewBitArray(n)}
}

func NewBITRangeAddPointGetFrom(n int, f func(i int) int) *BITRangeAddPointGetArray {
	return &BITRangeAddPointGetArray{bit: NewBitArrayFrom(n, f)}
}

func (b *BITRangeAddPointGetArray) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.bit.n {
		end = b.bit.n
	}
	if start >= end {
		return
	}
	b.bit.Add(start, delta)
	b.bit.Add(end, -delta)
}

func (b *BITRangeAddPointGetArray) Get(index int) int {
	return b.bit.QueryPrefix(index + 1)
}

func (b *BITRangeAddPointGetArray) String() string {
	res := []string{}
	for i := 0; i < b.bit.n; i++ {
		res = append(res, fmt.Sprintf("%d", b.Get(i)))
	}
	return fmt.Sprintf("BITRangeAddPointGetArray: [%v]", strings.Join(res, ", "))
}

type BITRangeAddPointGetArrayDual struct {
	n    int
	data []int
}

func NewBITRangeAddPointGetArrayDual(n int) *BITRangeAddPointGetArrayDual {
	return &BITRangeAddPointGetArrayDual{n: n, data: make([]int, n)}
}

func (b *BITRangeAddPointGetArrayDual) AddRange(start, end int, delta int) {
	neg := -delta
	for start < end {
		b.data[end-1] += delta
		end -= end & -end
	}
	for end < start {
		b.data[start-1] += neg
		start -= start & -start
	}
}

func (b *BITRangeAddPointGetArrayDual) Get(index int) int {
	res := 0
	for index++; index <= b.n; index += index & -index {
		res += b.data[index-1]
	}
	return res
}

func (b *BITRangeAddPointGetArrayDual) GetAll() []int {
	res := append(b.data[:0:0], b.data...)
	for i := b.n; i >= 1; i-- {
		j := i + (i & -i)
		if j <= b.n {
			res[i-1] += res[j-1]
		}
	}
	return res
}

type BITRangeAddPointGetMap struct {
	bit *BITMap
}

func NewBITRangeAddPointGetMap(n int) *BITRangeAddPointGetMap {
	return &BITRangeAddPointGetMap{bit: NewBITMap(n)}
}

func (b *BITRangeAddPointGetMap) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.bit.n {
		end = b.bit.n
	}
	if start >= end {
		return
	}
	b.bit.Add(start, delta)
	b.bit.Add(end, -delta)
}

func (b *BITRangeAddPointGetMap) Get(end int) int {
	return b.bit.QueryPrefix(end)
}

// !Range Add Range Sum, 0-based.
type BITRangeAddRangeSumArray struct {
	n    int
	bit0 *BITArray
	bit1 *BITArray
}

func NewBITRangeAddRangeSumArray(n int) *BITRangeAddRangeSumArray {
	return &BITRangeAddRangeSumArray{
		n:    n,
		bit0: NewBitArray(n),
		bit1: NewBitArray(n),
	}
}

func NewBITRangeAddRangeSumFrom(n int, f func(index int) int) *BITRangeAddRangeSumArray {
	return &BITRangeAddRangeSumArray{
		n:    n,
		bit0: NewBitArrayFrom(n, f),
		bit1: NewBitArray(n),
	}
}

func (b *BITRangeAddRangeSumArray) Add(index int, delta int) {
	b.bit0.Add(index, delta)
}

func (b *BITRangeAddRangeSumArray) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return
	}
	b.bit0.Add(start, -delta*start)
	b.bit0.Add(end, delta*end)
	b.bit1.Add(start, delta)
	b.bit1.Add(end, -delta)
}

func (b *BITRangeAddRangeSumArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	rightRes := b.bit1.QueryPrefix(end)*end + b.bit0.QueryPrefix(end)
	leftRes := b.bit1.QueryPrefix(start)*start + b.bit0.QueryPrefix(start)
	return rightRes - leftRes
}

func (b *BITRangeAddRangeSumArray) String() string {
	res := []string{}
	for i := 0; i < b.n; i++ {
		res = append(res, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITRangeAddRangeSumArray: [%v]", strings.Join(res, ", "))
}

// !Range Add Range Sum, 0-based.
type BITRangeAddRangeSumMap struct {
	n    int
	bit0 *BITMap
	bit1 *BITMap
}

func NewBITRangeAddRangeSumMap(n int) *BITRangeAddRangeSumMap {
	return &BITRangeAddRangeSumMap{
		n:    n + 5,
		bit0: NewBITMap(n),
		bit1: NewBITMap(n),
	}
}

func (b *BITRangeAddRangeSumMap) Add(index int, delta int) {
	b.bit0.Add(index, delta)
}

func (b *BITRangeAddRangeSumMap) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return
	}
	b.bit0.Add(start, -delta*start)
	b.bit0.Add(end, delta*end)
	b.bit1.Add(start, delta)
	b.bit1.Add(end, -delta)
}

func (b *BITRangeAddRangeSumMap) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	rightRes := b.bit1.QueryPrefix(end)*end + b.bit0.QueryPrefix(end)
	leftRes := b.bit1.QueryPrefix(start)*start + b.bit0.QueryPrefix(start)
	return rightRes - leftRes
}

// Fenwick Tree Prefix
// https://suisen-cp.github.io/cp-library-cpp/library/datastructure/fenwick_tree/fenwick_tree_prefix.hpp
// 如果每次都是查询前缀，那么可以使用Fenwick Tree Prefix 维护 monoid.
type S = int

func (*BITPrefixArray) e() S        { return 0 }
func (*BITPrefixArray) op(a, b S) S { return max(a, b) }
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type BITPrefixArray struct {
	n    int
	data []S
}

func NewBITPrefixArray(n int) *BITPrefixArray {
	res := &BITPrefixArray{}
	data := make([]S, n)
	for i := range data {
		data[i] = res.e()
	}
	res.n = n
	res.data = data
	return res
}

func NewBITPrefixFrom(n int, f func(index int) S) *BITPrefixArray {
	res := &BITPrefixArray{}
	total := res.e()
	data := make([]S, n)
	for i := range data {
		data[i] = f(i)
		total = res.op(total, data[i])
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = res.op(data[j-1], data[i-1])
		}
	}
	res.n = n
	res.data = data
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *BITPrefixArray) Update(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index-1] = f.op(f.data[index-1], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= end <= n
func (f *BITPrefixArray) QueryPrefix(end int) S {
	if end > f.n {
		end = f.n
	}
	res := f.e()
	for ; end > 0; end &= end - 1 {
		res = f.op(res, f.data[end-1])
	}
	return res
}

type V = int

func (*BITPrefixMap) e() V        { return 0 }
func (*BITPrefixMap) op(a, b V) V { return max(a, b) }

type BITPrefixMap struct {
	n    int
	data map[int]V
}

func NewBITPrefixMap(n int) *BITPrefixMap {
	return &BITPrefixMap{n: n, data: make(map[int]V)}
}

func (f *BITPrefixMap) Update(index int, value V) {
	for index++; index <= f.n; index += index & -index {
		f.data[index-1] = f.op(f._get(index-1), value)
	}
}

func (f *BITPrefixMap) QueryPrefix(end int) V {
	if end > f.n {
		end = f.n
	}
	res := f.e()
	for ; end > 0; end &= end - 1 {
		res = f.op(res, f._get(end-1))
	}
	return res
}

func (f *BITPrefixMap) _get(index int) V {
	if v, ok := f.data[index]; ok {
		return v
	}
	return f.e()
}

func maximumWhiteTiles(tiles [][]int, carpetLen int) int {
	bit := NewBITRangeAddRangeSumMap(1e9 + 10)
	for _, tile := range tiles {
		bit.AddRange(int(tile[0]), 1+int(tile[1]), 1)
	}

	res := 0
	for _, tile := range tiles {
		res = max(res, bit.QueryRange(int(tile[0]), 1+int(tile[0]+carpetLen-1)))
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func demo() {

	bitArray := NewBitArrayFrom(10, func(index int) int { return index + 1 })
	fmt.Println(bitArray)
	bitArray.Add(1, 1)
	fmt.Println(bitArray)
	fmt.Println(bitArray.Kth(0, 0))

	bp := NewBITRangeAddPointGetFrom(10, func(index int) int { return index + 1 })
	fmt.Println(bp)
	bp.AddRange(2, 2, 2)
	fmt.Println(bp)
	fmt.Println(bp.Get(3))

	bitArray2 := NewBITRangeAddRangeSumFrom(10, func(index int) int { return index + 1 })
	fmt.Println(bitArray2)
	bitArray2.AddRange(1, 3, 2)
	fmt.Println(bitArray2)
	bitArray2.AddRange(1, 3, -2)
	fmt.Println(bitArray2)
	bitArray2.AddRange(1, 1, -2)
	fmt.Println(bitArray2)

	bitPrefixMap := NewBITPrefixMap(int(1e9))
	bitPrefixMap.Update(1, 1)
	bitPrefixMap.Update(2e9, 2)
	fmt.Println(bitPrefixMap.QueryPrefix(int(1e8)))

}
