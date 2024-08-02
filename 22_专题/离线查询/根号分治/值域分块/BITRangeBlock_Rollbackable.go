package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	bit := NewBITRangeBlock(e, op)
	bit.Build(3, func(i int32) int { return 0 })
	fmt.Println(bit)
	bit.Update(0, 1)
	fmt.Println(bit)
	bit.Update(2, 2)
	fmt.Println(bit, bit.QueryRange(0, 2))
	bit.Update(2, 4)
	fmt.Println(bit, bit.QueryRange(0, 3))
}

// TLE
// https://leetcode.cn/problems/maximize-the-minimum-powered-city/description/
func maxPower(stations []int, r int, k int) int64 {
	n := len(stations)
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	bit := NewBITRangeBlock(e, op)
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
// !`O(1)`单点加，`O(sqrt(n))`区间和查询.
// 一般配合莫队算法使用.
type BITRangeBlock[E comparable] struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       *RollbackArray[E]
	_blockSum   *RollbackArray[E]
	e           func() E
	op          func(a, b E) E
}

func NewBITRangeBlock[E comparable](e func() E, op func(a, b E) E) *BITRangeBlock[E] {
	return &BITRangeBlock[E]{e: e, op: op}
}

func (b *BITRangeBlock[E]) GetTime() (time0, time1 int32) {
	return b._nums.GetTime(), b._blockSum.GetTime()
}

func (b *BITRangeBlock[E]) Rollback(time0, time1 int32) {
	b._nums.Rollback(time0)
	b._blockSum.Rollback(time1)
}

func (b *BITRangeBlock[E]) Build(n int32, f func(i int32) E) {
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
	nums := make([]E, n)
	for i := int32(0); i < n; i++ {
		nums[i] = f(i)
	}
	blockSum := make([]E, blockCount)
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
	b._nums = NewRollbackArrayFrom(nums)
	b._blockSum = NewRollbackArrayFrom(blockSum)
}

func (b *BITRangeBlock[E]) Update(index int32, delta E) {
	b._nums.Set(index, b.op(b._nums.Get(index), delta))
	bid := b._belong[index]
	b._blockSum.Set(bid, b.op(b._blockSum.Get(bid), delta))
}

func (b *BITRangeBlock[E]) QueryRange(start, end int32) E {
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

func (b *BITRangeBlock[E]) String() string {
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

type HistoryItem[V comparable] struct {
	index int32
	value V
}

type RollbackArray[V comparable] struct {
	n       int32
	data    []V
	history []HistoryItem[V]
}

func NewRollbackArray[V comparable](n int32, f func(index int32) V) *RollbackArray[V] {
	data := make([]V, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray[V]{
		n:    n,
		data: data,
	}
}

func NewRollbackArrayFrom[V comparable](data []V) *RollbackArray[V] {
	return &RollbackArray[V]{n: int32(len(data)), data: data}
}

func (r *RollbackArray[V]) GetTime() int32 {
	return int32(len(r.history))
}

func (r *RollbackArray[V]) Rollback(time int32) {
	for int32(len(r.history)) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *RollbackArray[V]) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray[V]) Get(index int32) V {
	return r.data[index]
}

func (r *RollbackArray[V]) Set(index int32, value V) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, HistoryItem[V]{index: index, value: r.data[index]})
	r.data[index] = value
	return true
}

func (r *RollbackArray[V]) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray[V]) Len() int32 {
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
