package main

import (
	"fmt"
	"math"
	"strings"
)

// https://leetcode.cn/problems/maximize-the-minimum-powered-city/description/
func maxPower(stations []int, r int, k int) int64 {
	n := len(stations)
	check := func(mid int) bool {
		bit := NewBITRangeBlockFastQuery(stations)
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
// `O(sqrt(n))`单点加，`O(1)`查询区间和.
// 一般配合莫队算法使用.
type BITRangeBlockFastQuery struct {
	_n           int
	_belong      []int
	_blockStart  []int
	_blockEnd    []int
	_blockCount  int
	_partPreSum  []int
	_blockPreSum []int
}

func NewBITRangeBlockFastQuery(lengthOrArray interface{}) *BITRangeBlockFastQuery {
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
	partPreSum := make([]int, n)
	blockPreSum := make([]int, blockCount)
	res := &BITRangeBlockFastQuery{
		_n:           n,
		_belong:      belong,
		_blockStart:  blockStart,
		_blockEnd:    blockEnd,
		_blockCount:  blockCount,
		_partPreSum:  partPreSum,
		_blockPreSum: blockPreSum,
	}
	if isArray {
		res.Build(lengthOrArray.([]int))
	}
	return res
}

func (b *BITRangeBlockFastQuery) Add(index int, delta int) {
	if index < 0 || index >= b._n {
		panic("index out of range")
	}
	bid := b._belong[index]
	for i := index; i < b._blockEnd[bid]; i++ {
		b._partPreSum[i] += delta
	}
	for id := bid + 1; id < b._blockCount; id++ {
		b._blockPreSum[id] += delta
	}
}

func (b *BITRangeBlockFastQuery) QueryRange(start int, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return 0
	}
	return b._query(end) - b._query(start)
}

func (b *BITRangeBlockFastQuery) Build(arr []int) {
	if len(arr) != b._n {
		panic("array length mismatch n")
	}
	curBlockSum := 0
	for bid := 0; bid < b._blockCount; bid++ {
		curPartSum := 0
		for i := b._blockStart[bid]; i < b._blockEnd[bid]; i++ {
			curPartSum += arr[i]
			b._partPreSum[i] = curPartSum
		}
		b._blockPreSum[bid] = curBlockSum
		curBlockSum += curPartSum
	}
}

func (b *BITRangeBlockFastQuery) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := range b._partPreSum {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	sb.WriteString("}")
	return sb.String()
}

func (b *BITRangeBlockFastQuery) _query(end int) int {
	if end <= 0 {
		return 0
	}
	return b._partPreSum[end-1] + b._blockPreSum[b._belong[end-1]]
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
