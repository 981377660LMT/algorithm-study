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
	CF1288E()
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
