// !队列适合模拟round robin的场景.

package main

import (
	"fmt"
)

// 281. 锯齿迭代器
// https://leetcode.cn/problems/zigzag-iterator/description/?envType=problem-list-v2&envId=design
type ZigzagIterator struct {
	iter *InterleaveLongestIterator[int]
}

func Constructor(v1, v2 []int) *ZigzagIterator {
	iter := NewInterleaveLongestIterator(NewSliceIterator(v1), NewSliceIterator(v2))
	return &ZigzagIterator{iter: iter}
}

func (this *ZigzagIterator) next() int {
	return this.iter.Next()
}

func (this *ZigzagIterator) hasNext() bool {
	return this.iter.HasNext()
}

func main() {
	// 示例数据
	a := []int{1}
	b := []int{10, 20, 30, 40}
	c := []int{44, 23}

	// 创建迭代器
	{
		iter := NewInterleaveIterator(NewSliceIterator(a), NewSliceIterator(b), NewSliceIterator(c))

		// 遍历迭代器
		for iter.HasNext() {
			fmt.Println(iter.Next())
		}
	}

	fmt.Println("")

	// 创建迭代器
	{
		iter := NewInterleaveLongestIterator(NewSliceIterator(a), NewSliceIterator(b), NewSliceIterator(c))

		// 遍历迭代器
		for iter.HasNext() {
			fmt.Println(iter.Next())
		}
	}
}

type Iterator[T any] interface {
	Next() T
	HasNext() bool
}

type SliceIterator[T any] struct {
	slice []T
	index int
}

func NewSliceIterator[T any](slice []T) *SliceIterator[T] {
	return &SliceIterator[T]{slice: slice}
}

func (s *SliceIterator[T]) Next() T {
	val := s.slice[s.index]
	s.index++
	return val
}

func (s *SliceIterator[T]) HasNext() bool {
	return s.index < len(s.slice)
}

// InterleaveIterator 实现了交错迭代器：
// 每轮从所有迭代器取一个值，任一迭代器耗尽则整个迭代结束.
// 等价于 python 中的 chain.from_iterable(zip(*iterables)).
type InterleaveIterator[T any] struct {
	roundValid   bool // 当前轮是否所有迭代器都有元素
	roundChecked bool // 当前轮是否已经检查过
	ptr          int32
	iters        []Iterator[T]
}

func NewInterleaveIterator[T any](iters ...Iterator[T]) *InterleaveIterator[T] {
	if len(iters) == 0 {
		return &InterleaveIterator[T]{iters: iters, roundValid: false, roundChecked: true}
	}
	return &InterleaveIterator[T]{iters: iters}
}

func (it *InterleaveIterator[T]) Next() T {
	iter := it.iters[it.ptr]
	res := iter.Next()
	it.ptr++
	if it.ptr == int32(len(it.iters)) {
		it.ptr = 0
		it.roundChecked = false
	}
	return res
}

func (it *InterleaveIterator[T]) HasNext() bool {
	if it.roundChecked {
		return it.roundValid
	}
	it.roundValid = true
	for _, iter := range it.iters {
		if !iter.HasNext() {
			it.roundValid = false
			break
		}
	}
	it.roundChecked = true
	return it.roundValid
}

// InterleaveLongestIterator 实现了交错最长迭代器：
// 轮询输出各迭代器的值，已耗尽的迭代器会被跳过，直到所有迭代器均耗尽
type InterleaveLongestIterator[T any] struct {
	iters []Iterator[T]
	queue *cycleQueue[int32]
}

func NewInterleaveLongestIterator[T any](iters ...Iterator[T]) *InterleaveLongestIterator[T] {
	n := len(iters)
	queue := newCycleQueue[int32](n)
	for i := 0; i < n; i++ {
		if iters[i].HasNext() {
			queue.Append(int32(i))
		}
	}
	return &InterleaveLongestIterator[T]{iters: iters, queue: queue}
}

func (it *InterleaveLongestIterator[T]) Next() T {
	iterId := it.queue.Popleft()
	iter := it.iters[iterId]
	res := iter.Next()
	if iter.HasNext() {
		it.queue.Append(iterId)
	}
	return res
}

func (it *InterleaveLongestIterator[T]) HasNext() bool {
	return it.queue.Len() > 0
}

type cycleQueue[T any] struct {
	values  []T
	start   int
	end     int
	maxSize int
	size    int
}

func newCycleQueue[T any](maxSize int) *cycleQueue[T] {
	queue := &cycleQueue[T]{values: make([]T, maxSize, maxSize), maxSize: maxSize}
	return queue
}

func (queue *cycleQueue[T]) Append(value T) {
	queue.values[queue.end] = value
	queue.end++
	if queue.end >= queue.maxSize {
		queue.end = 0
	}
	queue.size++
}

func (queue *cycleQueue[T]) Popleft() T {
	queue.size--
	value := queue.values[queue.start]
	queue.start++
	if queue.start >= queue.maxSize {
		queue.start = 0
	}
	return value
}

func (queue *cycleQueue[T]) Len() int {
	return queue.size
}
