package main

import (
	"fmt"
	"math"
	"strings"
)

// https://leetcode.cn/problems/maximize-the-minimum-powered-city/description/
func maxPower(stations []int, r int, k int) int64 {
	n := len(stations)
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	bit := NewBITRangeBlockFastQuery32(e, op, inv)
	check := func(mid int) bool {
		bit.Build(int32(len(stations)), func(i int32) int { return stations[i] })
		curK := k
		for i := 0; i < n; i++ {
			cur := bit.QueryRange(int32(max(0, i-r)), int32(min(i+r+1, n)))
			if cur < mid {
				diff := mid - cur
				bit.Update(min32(int32(i+r), int32(n-1)), diff)
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
// !`O(sqrt(n))`单点加，`O(1)`查询区间和，要求有逆元.
// 一般配合莫队算法使用.
type BITRangeBlockFastQuery32[E any] struct {
	_n           int32
	_belong      []int32
	_blockStart  []int32
	_blockEnd    []int32
	_blockCount  int32
	_partPreSum  []E
	_blockPreSum []E
	e            func() E
	op           func(a, b E) E
	inv          func(a E) E
}

func NewBITRangeBlockFastQuery32[E any](e func() E, op func(a, b E) E, inv func(a E) E) *BITRangeBlockFastQuery32[E] {
	return &BITRangeBlockFastQuery32[E]{e: e, op: op, inv: inv}
}

func (b *BITRangeBlockFastQuery32[E]) Build(n int32, f func(i int32) E) {
	blockSize := int32(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int32, n)
	for i := int32(0); i < n; i++ {
		belong[i] = i / blockSize
	}
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	for i := int32(0); i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	partPreSum := make([]E, n)
	blockPreSum := make([]E, blockCount)

	curBlockSum := b.e()
	for bid := int32(0); bid < blockCount; bid++ {
		curPartSum := b.e()
		for i := blockStart[bid]; i < blockEnd[bid]; i++ {
			curPartSum = b.op(curPartSum, f(i))
			partPreSum[i] = curPartSum
		}
		blockPreSum[bid] = curBlockSum
		curBlockSum = b.op(curBlockSum, curPartSum)
	}
	b._n = n
	b._belong = belong
	b._blockStart = blockStart
	b._blockEnd = blockEnd
	b._blockCount = blockCount
	b._partPreSum = partPreSum
	b._blockPreSum = blockPreSum
}

func (b *BITRangeBlockFastQuery32[E]) Update(index int32, delta E) {
	bid := b._belong[index]
	for i := index; i < b._blockEnd[bid]; i++ {
		b._partPreSum[i] = b.op(b._partPreSum[i], delta)
	}
	for id := bid + 1; id < b._blockCount; id++ {
		b._blockPreSum[id] = b.op(b._blockPreSum[id], delta)
	}
}

func (b *BITRangeBlockFastQuery32[E]) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return b.e()
	}
	return b.op(b._query(end), b.inv(b._query(start)))
}

func (b *BITRangeBlockFastQuery32[E]) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := int32(0); i < int32(len(b._partPreSum)); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", b.QueryRange(i, i+1)))
	}
	sb.WriteString("}")
	return sb.String()
}

func (b *BITRangeBlockFastQuery32[E]) _query(end int32) E {
	if end <= 0 {
		return b.e()
	}
	return b.op(b._partPreSum[end-1], b._blockPreSum[b._belong[end-1]])
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
