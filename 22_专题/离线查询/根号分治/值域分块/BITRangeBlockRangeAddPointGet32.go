package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	bit := NewBITRangeBlockRangeAddPointGet32(10)
	fmt.Println(bit)
	bit.AddRange(0, 8, 2)
	fmt.Println(bit)
}

// 基于分块实现的`树状数组`.
// `O(1)`单点查询，`O(sqrt(n))`区间加.
// 一般配合莫队算法使用.
type BITRangeBlockRangeAddPointGet32 struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       []int
	_blockLazy  []int
}

func NewBITRangeBlockRangeAddPointGet32(n int32) *BITRangeBlockRangeAddPointGet32 {
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
	res := &BITRangeBlockRangeAddPointGet32{
		_n:          n,
		_belong:     belong,
		_blockStart: blockStart,
		_blockEnd:   blockEnd,
		_nums:       nums,
		_blockLazy:  blockSum,
	}
	return res
}

func NewBITRangeBlockRangeAddPointGet32From(n int32, f func(i int32) int) *BITRangeBlockRangeAddPointGet32 {
	res := NewBITRangeBlockRangeAddPointGet32(n)
	res.Build(n, f)
	return res
}

func (b *BITRangeBlockRangeAddPointGet32) Get(index int32) int {
	if index < 0 || index >= b._n {
		panic("index out of range")
	}
	return b._nums[index] + b._blockLazy[b._belong[index]]
}

func (b *BITRangeBlockRangeAddPointGet32) AddRange(start, end int32, delta int) {
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
			b._nums[i] += delta
		}
		return
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		b._nums[i] += delta
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		b._blockLazy[bid] += delta
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		b._nums[i] += delta
	}
}

func (b *BITRangeBlockRangeAddPointGet32) Build(n int32, f func(i int32) int) {
	if n != b._n {
		panic("array length mismatch n")
	}
	for i := range b._nums {
		b._nums[i] = f(int32(i))
	}
	for i := range b._blockLazy {
		b._blockLazy[i] = 0
	}
}

func (b *BITRangeBlockRangeAddPointGet32) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := range b._nums {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.Get(int32(i))))
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
