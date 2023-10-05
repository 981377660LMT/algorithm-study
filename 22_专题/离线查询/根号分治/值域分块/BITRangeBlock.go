package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	bit := NewBITRangeBlock(10)
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
		bit := NewBITRangeBlock(stations)
		curK := k
		for i := 0; i < n; i++ {
			cur := bit.QueryRange(max(0, i-r), min(i+r+1, n))
			if cur < mid {
				diff := mid - cur
				bit.Add(min(i+r, n-1), diff)
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
// `O(1)`单点加，`O(sqrt(n))`查询区间和.
// 一般配合莫队算法使用.
type BITRangeBlock struct {
	_n          int
	_belong     []int
	_blockStart []int
	_blockEnd   []int
	_nums       []int
	_blockSum   []int
}

func NewBITRangeBlock(lengthOrArray interface{}) *BITRangeBlock {
	var n int
	var isArray bool
	if length, ok := lengthOrArray.(int); ok {
		n = length
	} else {
		n = len(lengthOrArray.([]int))
		isArray = true
	}
	blockSize := int(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int, n)
	for i := range belong {
		belong[i] = i / blockSize
	}
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	for i := range blockStart {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
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
	if isArray {
		res.Build(lengthOrArray.([]int))
	}
	return res
}

func (b *BITRangeBlock) Add(index int, delta int) {
	if index < 0 || index >= b._n {
		panic("index out of range")
	}
	b._nums[index] += delta
	b._blockSum[b._belong[index]] += delta
}

func (b *BITRangeBlock) QueryRange(start int, end int) int {
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

func (b *BITRangeBlock) Build(arr []int) {
	if len(arr) != b._n {
		panic("array length mismatch n")
	}
	for i := range b._nums {
		b._nums[i] = 0
	}
	for i := range b._blockSum {
		b._blockSum[i] = 0
	}
	for i, v := range arr {
		b.Add(i, v)
	}
}

func (b *BITRangeBlock) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i, v := range b._nums {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
