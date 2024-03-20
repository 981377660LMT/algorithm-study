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
		bit := NewBITRangeBlockFastQuery32From(int32(len(stations)), func(i int32) int { return stations[i] })
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
// `O(sqrt(n))`单点加，`O(1)`查询区间和.
// 一般配合莫队算法使用.
type BITRangeBlockFastQuery32 struct {
	_n           int32
	_belong      []int32
	_blockStart  []int32
	_blockEnd    []int32
	_blockCount  int32
	_partPreSum  []int
	_blockPreSum []int
}

func NewBITRangeBlockFastQuery32(n int32) *BITRangeBlockFastQuery32 {
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
	partPreSum := make([]int, n)
	blockPreSum := make([]int, blockCount)
	res := &BITRangeBlockFastQuery32{
		_n:           n,
		_belong:      belong,
		_blockStart:  blockStart,
		_blockEnd:    blockEnd,
		_blockCount:  blockCount,
		_partPreSum:  partPreSum,
		_blockPreSum: blockPreSum,
	}
	return res
}

func NewBITRangeBlockFastQuery32From(n int32, f func(i int32) int) *BITRangeBlockFastQuery32 {
	res := NewBITRangeBlockFastQuery32(n)
	res.Build(n, f)
	return res
}

func (b *BITRangeBlockFastQuery32) Add(index int32, delta int) {
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

func (b *BITRangeBlockFastQuery32) QueryRange(start, end int32) int {
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

func (b *BITRangeBlockFastQuery32) Build(n int32, f func(i int32) int) {
	if n != b._n {
		panic("array length mismatch n")
	}
	curBlockSum := 0
	for bid := int32(0); bid < b._blockCount; bid++ {
		curPartSum := 0
		for i := b._blockStart[bid]; i < b._blockEnd[bid]; i++ {
			curPartSum += f(i)
			b._partPreSum[i] = curPartSum
		}
		b._blockPreSum[bid] = curBlockSum
		curBlockSum += curPartSum
	}
}

func (b *BITRangeBlockFastQuery32) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := int32(0); i < int32(len(b._partPreSum)); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	sb.WriteString("}")
	return sb.String()
}

func (b *BITRangeBlockFastQuery32) _query(end int32) int {
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
