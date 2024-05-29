package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
)

func main() {
	bf := func(nums []int, k int) int {
		pairs := make([]int, 0, len(nums)*len(nums))
		for _, a := range nums {
			for _, b := range nums {
				pairs = append(pairs, a^b)
			}
		}
		sort.Ints(pairs)
		return pairs[k-1]
	}

	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rand.Intn(100)
		}
		k := rand.Intn(n*n) + 1
		if KthXorPair(nums, k) != bf(nums, k) {
			fmt.Println("error")
		}
	}
	fmt.Println("pass")
}

// 异或第k小的二元组(一共有n^2个).
// 时间复杂度为O(n(logM+logn)，其中M是最大的数，空间复杂度为O(n).
func KthXorPair(nums []int, k int) int {
	nums = append(nums[:0:0], nums...)
	n := int32(len(nums))
	buffer := NewBufferWithCleanerAndCapacity[*Interval](
		NewInterval,
		func(_ *Interval) {},
		2*n,
	)
	sort.Ints(nums)
	lastLevel := make([]*Interval, 0, n)
	curLevel := make([]*Interval, 0, n)
	lastLevel = append(lastLevel, NewIntervalWithBuffer(buffer, 0, n-1))
	level := int32(FloorLog64(uint64(nums[n-1])))
	mask := 0
	for ; level >= 0; level-- {
		curLevel = curLevel[:0]
		for _, interval := range lastLevel {
			l, r := interval.l, interval.r
			m := r
			for m >= l && (nums[m]>>level)&1 == 1 {
				m--
			}
			interval.m = m
		}
		total := 0
		for _, inter := range lastLevel {
			total += int(inter.m-inter.l+1) * int(inter.relative.m-inter.relative.l+1)
			total += int(inter.r-inter.m) * int(inter.relative.r-inter.relative.m)
		}
		if total < k {
			k -= total
			mask |= 1 << level
			for _, inter := range lastLevel {
				if inter.relative == inter {
					if inter.l <= inter.m && inter.m < inter.r {
						a := NewIntervalWithBuffer(buffer, inter.l, inter.m)
						b := NewIntervalWithBuffer(buffer, inter.m+1, inter.r)
						a.relative = b
						b.relative = a
						curLevel = append(curLevel, a, b)
					}
				} else if inter.r >= inter.relative.r {
					if inter.l <= inter.m && inter.relative.r > inter.relative.m {
						a := NewIntervalWithBuffer(buffer, inter.l, inter.m)
						b := NewIntervalWithBuffer(buffer, inter.relative.m+1, inter.relative.r)
						a.relative = b
						b.relative = a
						curLevel = append(curLevel, a, b)
					}
					if inter.m < inter.r && inter.relative.m >= inter.relative.l {
						a := NewIntervalWithBuffer(buffer, inter.m+1, inter.r)
						b := NewIntervalWithBuffer(buffer, inter.relative.l, inter.relative.m)
						a.relative = b
						b.relative = a
						curLevel = append(curLevel, a, b)
					}
				}
			}
		} else {
			for _, inter := range lastLevel {
				if inter.relative == inter {
					if inter.l <= inter.m {
						a := NewIntervalWithBuffer(buffer, inter.l, inter.m)
						a.relative = a
						curLevel = append(curLevel, a)
					}
					if inter.m < inter.r {
						a := NewIntervalWithBuffer(buffer, inter.m+1, inter.r)
						a.relative = a
						curLevel = append(curLevel, a)
					}
				} else if inter.r >= inter.relative.r {
					if inter.l <= inter.m && inter.relative.l <= inter.relative.m {
						a := NewIntervalWithBuffer(buffer, inter.l, inter.m)
						b := NewIntervalWithBuffer(buffer, inter.relative.l, inter.relative.m)
						a.relative = b
						b.relative = a
						curLevel = append(curLevel, a, b)
					}
					if inter.m < inter.r && inter.relative.m < inter.relative.r {
						a := NewIntervalWithBuffer(buffer, inter.m+1, inter.r)
						b := NewIntervalWithBuffer(buffer, inter.relative.m+1, inter.relative.r)
						a.relative = b
						b.relative = a
						curLevel = append(curLevel, a, b)
					}
				}
			}
		}

		for _, inter := range lastLevel {
			buffer.Release(inter)
		}

		lastLevel, curLevel = curLevel, lastLevel
	}

	return mask
}

func FloorLog64(x uint64) int {
	if x <= 0 {
		panic("IllegalArgumentException")
	}
	return 63 - bits.LeadingZeros64(x)
}

type Interval struct {
	l, r, m  int32
	relative *Interval
}

func NewInterval() *Interval {
	res := &Interval{}
	res.relative = res
	return res
}

func NewIntervalWithBuffer(buffer *Buffer[*Interval], l, r int32) *Interval {
	res := buffer.Alloc()
	res.l = l
	res.r = r
	return res
}

// A buffer that recycles objects.
type Buffer[T any] struct {
	recycles    []T
	supplier    func() T // Supplier is a function that returns a new object.
	cleaner     func(T)  // Cleaner is a function that cleans up / initializes an object.
	allocTime   int32
	releaseTime int32
}

func NewBuffer[T any](supplier func() T) *Buffer[T] {
	return NewBufferWithCleaner(supplier, func(T) {})
}

func NewBufferWithCleaner[T any](supplier func() T, cleaner func(T)) *Buffer[T] {
	return NewBufferWithCleanerAndCapacity(supplier, cleaner, 0)
}

func NewBufferWithCleanerAndCapacity[T any](supplier func() T, cleaner func(T), capacity int32) *Buffer[T] {
	return &Buffer[T]{
		recycles: make([]T, 0, capacity),
		supplier: supplier,
		cleaner:  cleaner,
	}
}

func (b *Buffer[T]) Alloc() T {
	b.allocTime++
	if len(b.recycles) == 0 {
		res := b.supplier()
		b.cleaner(res)
		return res
	} else {
		res := b.recycles[len(b.recycles)-1]
		b.recycles = b.recycles[:len(b.recycles)-1]
		return res
	}
}

func (b *Buffer[T]) Release(e T) {
	b.releaseTime++
	b.cleaner(e)
	b.recycles = append(b.recycles, e)
}

func (b *Buffer[T]) Check() {
	if b.allocTime != b.releaseTime {
		panic(fmt.Sprintf("Buffer alloc %d but release %d", b.allocTime, b.releaseTime))
	}
}
