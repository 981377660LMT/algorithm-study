package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	e := func() int32 { return 0 }
	op := func(a, b int32) int32 { return max32(a, b) }
	bit := NewBITRangeBlock(e, op)
	bit.Build(3, func(i int32) int32 { return 0 })
	fmt.Println(bit)
	bit.Update(0, 1)
	fmt.Println(bit)
	bit.Update(2, 2)
	fmt.Println(bit, bit.QueryRange(0, 2))
	bit.Update(2, 4)
	fmt.Println(bit, bit.QueryRange(0, 3))
	time0, time1 := bit.GetTime()
	bit.Update(1, 3)
	fmt.Println(bit, bit.QueryRange(0, 3))
	bit.Rollback(time0, time1)
	fmt.Println(bit, bit.QueryRange(0, 3))
	bit.Rollback(time0-1, time1-1)
	fmt.Println(bit, bit.QueryRange(0, 3))
}

// 基于分块实现的`树状数组`.
// !`O(1)`单点加，`O(sqrt(n))`区间和查询.
// 一般配合莫队算法使用.
type BITRangeBlock struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       *RollbackArray32
	_blockSum   *RollbackArray32
	e           func() int32
	op          func(a, b int32) int32
}

func NewBITRangeBlock(e func() int32, op func(a, b int32) int32) *BITRangeBlock {
	return &BITRangeBlock{e: e, op: op}
}

func (b *BITRangeBlock) GetTime() (time0, time1 int32) {
	return b._nums.GetTime(), b._blockSum.GetTime()
}

func (b *BITRangeBlock) Rollback(time0, time1 int32) {
	b._nums.Rollback(time0)
	b._blockSum.Rollback(time1)
}

func (b *BITRangeBlock) Build(n int32, f func(i int32) int32) {
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
		tmp := blockStart[i] + blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		nums[i] = f(i)
	}
	blockSum := make([]int32, blockCount)
	for i := int32(0); i < blockCount; i++ {
		blockSum[i] = b.e()
	}
	for i := int32(0); i < n; i++ {
		bid := belong[i]
		blockSum[bid] = b.op(blockSum[bid], f(i))
	}
	b._n = n
	b._belong = belong
	b._blockStart = blockStart
	b._blockEnd = blockEnd
	b._nums = NewRollbackArray32From(nums)
	b._blockSum = NewRollbackArray32From(blockSum)
}

func (b *BITRangeBlock) Update(index int32, delta int32) {
	b._nums.Set(index, b.op(b._nums.Get(index), delta))
	bid := b._belong[index]
	b._blockSum.Set(bid, b.op(b._blockSum.Get(bid), delta))
}

func (b *BITRangeBlock) QueryRange(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return b.e()
	}
	res := b.e()
	bid1 := b._belong[start]
	bid2 := b._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			res = b.op(res, b._nums.Get(i))
		}
		return res
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		res = b.op(res, b._nums.Get(i))
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res = b.op(res, b._blockSum.Get(bid))
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		res = b.op(res, b._nums.Get(i))
	}
	return res
}

func (b *BITRangeBlock) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := 0; i < int(b._n); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", b.QueryRange(int32(i), int32(i+1))))
	}
	sb.WriteString("}")
	return sb.String()
}

const mask int = 1<<32 - 1

type RollbackArray32 struct {
	n       int32
	data    []int32
	history []int // (index, value)
}

func NewRollbackArray32(n int32, f func(index int32) int32) *RollbackArray32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray32{
		n:    n,
		data: data,
	}
}

func NewRollbackArray32From(data []int32) *RollbackArray32 {
	return &RollbackArray32{n: int32(len(data)), data: data}
}

func (r *RollbackArray32) GetTime() int32 {
	return int32(len(r.history))
}

func (r *RollbackArray32) Rollback(time int32) {
	for int32(len(r.history)) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair>>32] = int32(pair & mask)
	}
}

func (r *RollbackArray32) Get(index int32) int32 {
	return r.data[index]
}

func (r *RollbackArray32) Set(index int32, value int32) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, int(index)<<32|int(r.data[index]))
	r.data[index] = value
	return true
}

func (r *RollbackArray32) GetAll() []int32 {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray32) Len() int32 {
	return r.n
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
