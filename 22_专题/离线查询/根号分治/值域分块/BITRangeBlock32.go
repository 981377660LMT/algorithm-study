package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	bit := NewBITRangeBlockFrom(10, func(i int32) int { return 0 })
	fmt.Println(bit)
	bit.Add(0, 1)
	fmt.Println(bit)
	bit.Add(2, 2)
	fmt.Println(bit, bit.QueryRange(0, 2))
}

// TLE
// https://leetcode.cn/problems/maximize-the-minimum-powered-city/description/
func maxPower(stations []int, r int, k int) int64 {
	n := len(stations)
	check := func(mid int) bool {
		bit := NewBITRangeBlockFrom(int32(len(stations)), func(i int32) int { return stations[i] })
		curK := k
		for i := 0; i < n; i++ {
			cur := bit.QueryRange(int32(max(0, i-r)), int32(min(i+r+1, n)))
			if cur < mid {
				diff := mid - cur
				bit.Add(min32(int32(i+r), int32(n-1)), diff)
				curK -= diff
				if curK < 0 {
					return false
				}
			}
		}
		return true
	}

	left := 1
	right := int(2e15)
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return int64(right)
}

// 基于分块实现的`树状数组`.
// `O(1)`单点加，`O(sqrt(n))`区间和查询.
// 一般配合莫队算法使用.
type BITRangeBlock struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       []int
	_blockSum   []int
}

func NewBITRangeBlock(n int32) *BITRangeBlock {
	blockSize := int32(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int32, n)
	for i := range belong {
		belong[i] = int32(i) / blockSize
	}
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	for i := range blockStart {
		blockStart[i] = int32(i) * blockSize
		tmp := (int32(i) + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	nums := make([]int, n)
	blockSum := make([]int, blockCount)
	res := &BITRangeBlock{
		_n:          n,
		_belong:     belong,
		_blockStart: blockStart,
		_blockEnd:   blockEnd,
		_nums:       nums,
		_blockSum:   blockSum,
	}
	return res
}

func NewBITRangeBlockFrom(n int32, f func(i int32) int) *BITRangeBlock {
	res := NewBITRangeBlock(n)
	res.Build(n, f)
	return res
}

func (b *BITRangeBlock) Add(index int32, delta int) {
	if index < 0 || index >= b._n {
		panic("index out of range")
	}
	b._nums[index] += delta
	b._blockSum[b._belong[index]] += delta
}

func (b *BITRangeBlock) QueryRange(start, end int32) int {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return 0
	}
	res := 0
	bid1 := b._belong[start]
	bid2 := b._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			res += b._nums[i]
		}
		return res
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		res += b._nums[i]
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += b._blockSum[bid]
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		res += b._nums[i]
	}
	return res
}

func (b *BITRangeBlock) Build(n int32, f func(i int32) int) {
	if n != b._n {
		panic("array length mismatch n")
	}
	for i := range b._nums {
		b._nums[i] = 0
	}
	for i := range b._blockSum {
		b._blockSum[i] = 0
	}
	for i := int32(0); i < n; i++ {
		b.Add(i, f(i))
	}
}

func (b *BITRangeBlock) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := 0; i < int(b._n); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(int32(i), int32(i+1))))
	}
	sb.WriteString("}")
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
