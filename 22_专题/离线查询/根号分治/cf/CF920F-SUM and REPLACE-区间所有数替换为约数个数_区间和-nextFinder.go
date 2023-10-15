// CF920F-SUM and REPLACE-区间所有数替换为约数个数_区间和

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

var C []int

func init() {
	C = CountFactorsOfAll(1e6 + 10)
}

// 0 start end : 将区间所有数nums[i]替换为nums[i]的约数个数.
// 1 start end : 求区间所有数的和.
// n,q<=3e5, nums[i]<=1e6
//
// !类似区间开方, 有限次操作内nums[i]会不发生变化, 之后就不会变了.
func RangeReplaceRangeSum(nums []int, operations [][3]int) []int {
	nums = append(nums[:0:0], nums...)
	n := len(nums)
	nextFinder := NewNextFinder(n)
	bit := NewBitArrayFrom(nums)
	res := []int{}

	for _, op := range operations {
		op, start, end := op[0], op[1], op[2]
		if op == 0 {
			i := start
			for i < end {
				pre := nums[i]
				nums[i] = C[nums[i]]
				bit.Add(i, nums[i]-pre)
				if nums[i] == pre { // 不变
					nextFinder.Erase(i)
				}
				i = nextFinder.Next(i + 1)
			}
		} else {
			res = append(res, bit.QueryRange(start, end))
		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	operations := make([][3]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			operations[i] = [3]int{0, start, end}
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			operations[i] = [3]int{1, start, end}
		}
	}

	res := RangeReplaceRangeSum(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// LinearSequenceUnionFind 线性序列并查集(NextFinder).
type NextFinder struct {
	n     int
	right []int
	data  []uint64
}

func NewNextFinder(n int) *NextFinder {
	len := (n >> 6) + 1
	f := &NextFinder{
		n:     n,
		right: make([]int, len),
		data:  make([]uint64, len),
	}
	MASK := uint64(1<<64 - 1)
	for i := range f.right {
		f.right[i] = i
		f.data[i] = MASK
	}
	return f
}

// Next 下一个
//
//	如果不存在，返回n.
func (f *NextFinder) Next(x int) int {
	if x < 0 {
		x = 0
	}
	n := f.n
	if x >= n {
		return n
	}
	div := x >> 6
	mod := x & 63
	mask := f.data[div] >> mod
	if mask != 0 {
		return ((div << 6) | mod) + bits.TrailingZeros64(mask)
	}
	div = f.findNext(div + 1)
	return (div << 6) + bits.TrailingZeros64(f.data[div])
}

// Erase 删除
func (f *NextFinder) Erase(x int) {
	div := x >> 6
	mod := x & 63
	if (f.data[div]>>mod)&1 != 0 { // flip
		f.data[div] ^= 1 << mod
	}
	if f.data[div] == 0 {
		f.right[div] = div + 1 // union to right
	}
}

func (f *NextFinder) Has(x int) bool {
	if x < 0 || x >= f.n {
		return false
	}
	return (f.data[x>>6]>>(x&63))&1 != 0
}

func (f *NextFinder) String() string {
	sb := []string{}
	for i := 0; i < f.n; i++ {
		if f.Has(i) {
			sb = append(sb, fmt.Sprintf("%d", i))
		}
	}
	return "Finder(" + strings.Join(sb, ",") + ")"
}

func (f *NextFinder) findNext(x int) int {
	if f.right[x] == x {
		return x
	}
	f.right[x] = f.findNext(f.right[x])
	return f.right[x]
}

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 返回[0,upper]的所有数的约数个数.
func CountFactorsOfAll(upper int) []int {
	res := make([]int, upper+1)
	for i := 1; i <= upper; i++ {
		for j := i; j <= upper; j += i {
			res[j]++
		}
	}
	return res
}
