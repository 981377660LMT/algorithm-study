// 下标从0开始
// 1.BITArray: 区间修改, 单点查询
// 2.BITMap: 区间修改, 单点查询
// 3.BITRangeAddRangeSumArray: 区间修改, 区间查询
// 4.BITRangeAddRangeSumMap: 区间修改, 区间查询
// 5.BITPrefixArray: 单点修改, 区间查询
// 6.BITPrefixMap: 单点修改, 区间查询

package main

import (
	"fmt"
	"strings"
)

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

func NewBitArrayFrom(n int, f func(index int) int) *BITArray {
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

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
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
func (f *BITPrefixArray) Query(end int) S {
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

func (f *BITPrefixMap) Query(end int) V {
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

func main() {

	bitArray := NewBitArrayFrom(10, func(index int) int { return index + 1 })
	fmt.Println(bitArray)
	bitArray.Add(1, 1)
	fmt.Println(bitArray)
	fmt.Println(bitArray.Kth(0))

	bitArray2 := NewBITRangeAddRangeSumFrom(10, func(index int) int { return index + 1 })
	fmt.Println(bitArray2)
	bitArray2.AddRange(1, 3, 1)
	fmt.Println(bitArray2)
	fmt.Println(bitArray2.QueryRange(1, 3))
	fmt.Println(bitArray2.QueryRange(1, 4))
	bitArray2.Add(1, 1)
	fmt.Println(bitArray2)

	bitPrefixMap := NewBITPrefixMap(int(1e9))
	bitPrefixMap.Update(1, 1)
	bitPrefixMap.Update(2e9, 2)
	fmt.Println(bitPrefixMap.Query(int(1e8)))

}
