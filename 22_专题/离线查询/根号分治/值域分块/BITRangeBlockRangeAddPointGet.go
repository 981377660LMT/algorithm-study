package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	e, op := func() int { return 0 }, func(a, b int) int { return max(a, b) }
	bit := NewBITRangeBlockRangeAddPointGet32(10, e, op)
	fmt.Println(bit)
	bit.UpdateRange(0, 8, 2)
	fmt.Println(bit)
	bit.UpdateRange(2, 4, 3)
	fmt.Println(bit, bit.Get(3))
}

// 基于分块实现的`树状数组`.
// !`O(1)`单点查询，`O(sqrt(n))`区间加.
// 一般配合莫队算法使用.
type BITRangeBlockRangeAddPointGet32[E any] struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       []E
	_blockLazy  []E
	e           func() E
	op          func(a, b E) E
}

func NewBITRangeBlockRangeAddPointGet32[E any](n int32, e func() E, op func(a, b E) E) *BITRangeBlockRangeAddPointGet32[E] {
	b := &BITRangeBlockRangeAddPointGet32[E]{e: e, op: op}
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
	nums := make([]E, n)
	for i := int32(0); i < n; i++ {
		nums[i] = e()
	}
	blockLazy := make([]E, blockCount)
	for i := int32(0); i < blockCount; i++ {
		blockLazy[i] = e()
	}
	b._n = n
	b._belong = belong
	b._blockStart = blockStart
	b._blockEnd = blockEnd
	b._nums = nums
	b._blockLazy = blockLazy
	return b
}

func (b *BITRangeBlockRangeAddPointGet32[E]) Get(index int32) E {
	return b.op(b._nums[index], b._blockLazy[b._belong[index]])
}

func (b *BITRangeBlockRangeAddPointGet32[E]) UpdateRange(start, end int32, delta E) {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return
	}
	bid1 := b._belong[start]
	bid2 := b._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			b._nums[i] = b.op(b._nums[i], delta)
		}
		return
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		b._nums[i] = b.op(b._nums[i], delta)
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		b._blockLazy[bid] = b.op(b._blockLazy[bid], delta)
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		b._nums[i] = b.op(b._nums[i], delta)
	}
}

func (b *BITRangeBlockRangeAddPointGet32[E]) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := range b._nums {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", b.Get(int32(i))))
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
