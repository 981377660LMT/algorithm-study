// https://www.luogu.com.cn/problem/P1438
// 在线区间一次函数叠加求和/区间等差数列求和

// AddLinear(start, end, k, b):
//   O(logn) 为数组 [start, end) 的每个位置 i 加上 k*(i-start)+b.
// Get(index):
//   O(logn) 返回数组的累加和.

// OnlineRangeLinearAddRangeSum
//
// 原序列：0 0 0 0 0 0
// 差分序列：0 0 0 0 0 0
// 等差序列：1 3 5 7 9
// 加上等差数列后的序列：1 3 5 7 9 0
// 然后差分：1 2 2 2 2 -9
//
// nums[start]+=b
// nums[start+1,end)+=k
// nums[end]-=b+k*(end-start-1)

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	P1438()
}

// P1438 无聊的数列
// https://www.luogu.com.cn/problem/P1438
func P1438() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	S := NewRangeLinearAddPointGetOnline(n, func(i int) int { return nums[i] })
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var start, end int
			var first, diff int
			fmt.Fscan(in, &start, &end, &first, &diff)
			start--
			S.AddLinear(start, end, diff, first)
		} else {
			var index int
			fmt.Fscan(in, &index)
			index--
			fmt.Fprintln(out, S.Get(index))
		}
	}
}

type RangeLinearAddPointGetOnline struct {
	n   int
	bit *BITRangeAddRangeSumArray
}

func NewRangeLinearAddPointGetOnline(n int, f func(i int) int) *RangeLinearAddPointGetOnline {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = f(i)
	}
	for i := n - 1; i > 0; i-- {
		nums[i] -= nums[i-1]
	}
	bit := NewBITRangeAddRangeSumFrom(n, func(i int) int { return nums[i] })
	return &RangeLinearAddPointGetOnline{n: n, bit: bit}
}

// Add `k*(x-start)+b` to A[x] for x in [start, end).
//
//	0 <= start <= end <= n.
func (s *RangeLinearAddPointGetOnline) AddLinear(start, end, k, b int) {
	s.bit.AddRange(start, start+1, b)
	s.bit.AddRange(start+1, end, k)
	s.bit.AddRange(end, end+1, -b-k*(end-start-1))
}

func (s *RangeLinearAddPointGetOnline) Get(index int) int {
	return s.bit.QueryRange(0, index+1)
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
