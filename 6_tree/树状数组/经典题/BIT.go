// BITArray
// BITMap
// BIT2Array
// BIT2Map
// BIT2D

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n    int
	log  int
	data []int
}

func NewBitArray(n int) *BITArray {
	return &BITArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func NewBitArrayFrom(arr []int) *BITArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BITArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

func (b *BITArray) Add(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r)
func (b *BITArray) Query(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r).
func (b *BITArray) QueryRange(l, r int) int {
	return b.Query(r) - b.Query(l)
}

// 返回闭区间[0,k]的总和>=x的最小k.要求序列单调增加.
func (b *BITArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 返回闭区间[0,k]的总和>x的最小k.要求序列单调增加.
func (b *BITArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
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
	n    int
	data map[int]int
}

func NewBITMap(n int) *BITMap {
	return &BITMap{n: n + 5, data: make(map[int]int, 1<<10)}
}

func (b *BITMap) Add(i, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r)
func (b *BITMap) Query(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r).
func (b *BITMap) QueryRange(l, r int) int {
	return b.Query(r) - b.Query(l)
}

// !Range Add Range Sum, 0-based.
type BIT2Array struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITArray2(n int) *BIT2Array {
	return &BIT2Array{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// 切片内[start, end)的每个元素加上delta.
//
//	0<=start<=end<=n
func (b *BIT2Array) Add(start, end, delta int) {
	end--
	b.add(start, delta)
	b.add(end+1, -delta)
}

// 求切片内[start, end)的和.
//
//	0<=start<=end<=n
func (b *BIT2Array) Query(start, end int) int {
	end--
	return b.query(end) - b.query(start-1)
}

func (b *BIT2Array) add(index, delta int) {
	index++
	for i := index; i <= b.n; i += i & -i {
		b.tree1[i] += delta
		b.tree2[i] += (index - 1) * delta
	}
}

func (b *BIT2Array) query(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	for i := index; i > 0; i &= i - 1 {
		res += index*b.tree1[i] - b.tree2[i]
	}
	return res
}

// !Range Add Range Sum, 0-based.
type BIT2Map struct {
	n     int
	tree1 map[int]int
	tree2 map[int]int
}

func NewBIT2Map(n int) *BIT2Map {
	return &BIT2Map{
		n:     n + 5,
		tree1: make(map[int]int, 1<<10),
		tree2: make(map[int]int, 1<<10),
	}
}

func (bit *BIT2Map) Add(left, right, delta int) {
	right--
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BIT2Map) Query(left, right int) int {
	right--
	return bit.query(right) - bit.query(left-1)
}

func (bit *BIT2Map) add(index, delta int) {
	index++
	for i := index; i <= bit.n; i += i & -i {
		bit.tree1[i] += delta
		bit.tree2[i] += (index - 1) * delta
	}
}

func (bit *BIT2Map) query(index int) int {
	index++
	if index > bit.n {
		index = bit.n
	}

	res := 0
	for i := index; i > 0; i &= i - 1 {
		res += index*bit.tree1[i] - bit.tree2[i]
	}
	return res
}

func maximumWhiteTiles(tiles [][]int, carpetLen int) int {
	bit := NewBIT2Map(1e9 + 10)
	for _, tile := range tiles {
		bit.Add(int(tile[0]), 1+int(tile[1]), 1)
	}

	res := 0
	for _, tile := range tiles {
		res = max(res, bit.Query(int(tile[0]), 1+int(tile[0]+carpetLen-1)))
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Fenwick Tree Prefix
// https://suisen-cp.github.io/cp-library-cpp/library/datastructure/fenwick_tree/fenwick_tree_prefix.hpp
// 如果每次都是查询前缀，那么可以使用Fenwick Tree Prefix 维护 monoid.
type S = int

func (*FenwickTreePrefix) e() S        { return 0 }
func (*FenwickTreePrefix) op(a, b S) S { return max(a, b) }
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type FenwickTreePrefix struct {
	n    int
	data []S
}

func NewFenwickTreePrefix(n int) *FenwickTreePrefix {
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 0; i < n+1; i++ {
		res.data[i] = res.e()
	}
	return res
}

func NewFenwickTreePrefixWithSlice(nums []S) *FenwickTreePrefix {
	n := len(nums)
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 1; i < n+1; i++ {
		res.data[i] = nums[i-1]
	}
	for i := 1; i < n+1; i++ {
		if j := i + (i & -i); j <= n {
			res.data[j] = res.op(res.data[j], res.data[i])
		}
	}
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *FenwickTreePrefix) Update(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index] = f.op(f.data[index], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *FenwickTreePrefix) Query(right int) S {
	res := f.e()
	if right > f.n {
		right = f.n
	}
	for ; right > 0; right -= right & -right {
		res = f.op(res, f.data[right])
	}
	return res
}

func main() {
	pf := NewFenwickTreePrefixWithSlice([]S{1, 2, 3, 4, 14, 1, 2, 3})
	fmt.Println(pf.Query(1))
	fmt.Println(pf.Query(100))
	pf.Update(9, 10)
	fmt.Println(pf.Query(100))

	bitArray := NewBitArrayFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(bitArray)
	bitArray.Add(1, 1)
	fmt.Println(bitArray)

}
