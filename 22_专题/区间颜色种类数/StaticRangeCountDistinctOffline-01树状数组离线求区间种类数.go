// StaticRangeCountDistinctOffline-01树状数组离线求区间种类数

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	P1972()
}

// P1972 [SDOI2009] HH的项链
// https://www.luogu.com.cn/problem/P1972
// 区间颜色种类数
func P1972() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	S := NewStaticRangeCountDistinct(nums)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		start--
		S.AddQuery(start, end)
	}
	rangeUniqueCount := S.Run()

	for _, v := range rangeUniqueCount {
		fmt.Fprintln(out, v)
	}
}

// 01树状数组离线求区间种类数.
// O((N + Q)logN)
type StaticRangeCountDistinctOffline struct {
	n     int
	nums  []int
	query [][2]int32
}

func NewStaticRangeCountDistinct(nums []int) *StaticRangeCountDistinctOffline {
	return &StaticRangeCountDistinctOffline{
		n:    len(nums),
		nums: nums,
	}
}

func (sr *StaticRangeCountDistinctOffline) AddQuery(start, end int) {
	sr.query = append(sr.query, [2]int32{int32(start), int32(end)})
}

func (sr *StaticRangeCountDistinctOffline) Run() []int {
	n := sr.n
	q := int32(len(sr.query))
	res := make([]int, q)
	groupByEnd := make([][]int32, sr.n+1)
	for qi := int32(0); qi < q; qi++ {
		groupByEnd[sr.query[qi][1]] = append(groupByEnd[sr.query[qi][1]], qi)
	}
	bit := NewBITArray01(n)
	last := make(map[int]int, n)
	for i, x := range sr.nums {
		if pre, ok := last[x]; ok {
			bit.Remove(pre)
		}
		last[x] = i
		bit.Add(i)
		for _, qi := range groupByEnd[i+1] {
			start, end := sr.query[qi][0], sr.query[qi][1]
			res[qi] = bit.QueryRange(int(start), int(end))
		}
	}

	return res
}

// 01树状数组.

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

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
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
func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}
