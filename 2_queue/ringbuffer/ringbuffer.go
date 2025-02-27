// https://github.com/emirpasic/gods/blob/master/queues/circularbuffer/circularbuffer.go
// ringbuffer 是一种队列的实现方式，它是一个环形的缓冲区，可以用来实现队列、栈等数据结构
// 经过工业场景验证过的 ringbuffer 库：disruptor
// go的 buffered chan 是一个 ringbuffer
//
// RingCycle/RingBuffer

package main

import (
	"fmt"
	"strings"
)

func main() {
	q := NewRingBuffer[int](3)
	q.Append(1)
	q.Append(2)
	q.Append(3)
	q.Append(4)
	fmt.Println(q) // RingBuffer\n1, 2, 3
	fmt.Println(q.Size())

	fmt.Println(q.Popleft()) // 2
	fmt.Println(q)           // RingBuffer\n3, 4
	fmt.Println(q.Head())    // 3
	fmt.Println(q.At(0))     // 3
	fmt.Println(q.At(1))     // 4
	fmt.Println(q.At(-1))    // 4
}

// 固定大小的环形缓冲区, 可用于双端队列、固定大小滑动窗口等场景.
type RingBuffer[T comparable] struct {
	values  []T
	start   int
	end     int
	maxSize int
	size    int
}

func NewRingBuffer[T comparable](maxSize int) *RingBuffer[T] {
	if maxSize < 1 {
		panic("Invalid maxSize, should be at least 1")
	}
	queue := &RingBuffer[T]{values: make([]T, maxSize, maxSize), maxSize: maxSize}
	return queue
}

func (queue *RingBuffer[T]) Append(value T) {
	if queue.Full() {
		queue.Popleft()
	}
	queue.values[queue.end] = value
	queue.end++
	if queue.end >= queue.maxSize {
		queue.end = 0
	}
	queue.size++
}

func (queue *RingBuffer[T]) AppendLeft(value T) {
	if queue.Full() {
		queue.Pop()
	}
	queue.start--
	if queue.start < 0 {
		queue.start = queue.maxSize - 1
	}
	queue.values[queue.start] = value
	queue.size++
}

func (queue *RingBuffer[T]) Pop() T {
	if queue.Empty() {
		panic("RingBuffer is empty")
	}
	queue.size--
	queue.end--
	if queue.end < 0 {
		queue.end = queue.maxSize - 1
	}
	value := queue.values[queue.end]
	return value
}

func (queue *RingBuffer[T]) Popleft() T {
	if queue.Empty() {
		panic("RingBuffer is empty")
	}
	queue.size--
	value := queue.values[queue.start]
	queue.start++
	if queue.start >= queue.maxSize {
		queue.start = 0
	}
	return value
}

func (queue *RingBuffer[T]) Head() T {
	if queue.Empty() {
		panic("RingBuffer is empty")
	}
	return queue.values[queue.start]
}

func (queue *RingBuffer[T]) Tail() T {
	if queue.Empty() {
		panic("RingBuffer is empty")
	}
	index := queue.end - 1
	if index < 0 {
		index = queue.maxSize - 1
	}
	return queue.values[index]
}

func (queue *RingBuffer[T]) At(index int) T {
	size := queue.Size()
	if index < 0 {
		index += size
	}
	if index < 0 || index >= size {
		panic("Index out of range")
	}
	index += queue.start
	if index >= queue.maxSize {
		index -= queue.maxSize
	}
	return queue.values[index]
}

func (queue *RingBuffer[T]) Empty() bool {
	return queue.Size() == 0
}

func (queue *RingBuffer[T]) Full() bool {
	return queue.Size() == queue.maxSize
}

func (queue *RingBuffer[T]) Size() int {
	return queue.size
}

func (queue *RingBuffer[T]) Clear() {
	queue.start = 0
	queue.end = 0
	queue.size = 0
}

func (queue *RingBuffer[T]) ForEach(f func(value T)) {
	ptr := queue.start
	for i := 0; i < queue.Size(); i++ {
		f(queue.values[ptr])
		ptr++
		if ptr >= queue.maxSize {
			ptr = 0
		}
	}
}

func (queue *RingBuffer[T]) String() string {
	str := "RingBuffer\n"
	var values []string
	queue.ForEach(func(value T) { values = append(values, fmt.Sprintf("%v", value)) })
	str += strings.Join(values, ", ")
	return str
}

// 641. 设计循环双端队列
// https://leetcode.cn/problems/design-circular-deque/description/
type MyCircularDeque struct {
	k     int
	queue *RingBuffer[int]
}

func Constructor(k int) MyCircularDeque {
	return MyCircularDeque{k: k, queue: NewRingBuffer[int](k)}
}

func (this *MyCircularDeque) InsertFront(value int) bool {
	if this.queue.Size() == this.k {
		return false
	}
	this.queue.AppendLeft(value)
	return true
}

func (this *MyCircularDeque) InsertLast(value int) bool {
	if this.queue.Size() == this.k {
		return false
	}
	this.queue.Append(value)
	return true
}

func (this *MyCircularDeque) DeleteFront() bool {
	if this.queue.Size() == 0 {
		return false
	}
	this.queue.Popleft()
	return true
}

func (this *MyCircularDeque) DeleteLast() bool {
	if this.queue.Size() == 0 {
		return false
	}
	this.queue.Pop()
	return true
}

func (this *MyCircularDeque) GetFront() int {
	if this.queue.Size() == 0 {
		return -1
	}
	return this.queue.Head()
}

func (this *MyCircularDeque) GetRear() int {
	if this.queue.Size() == 0 {
		return -1
	}
	return this.queue.Tail()
}

func (this *MyCircularDeque) IsEmpty() bool {
	return this.queue.Size() == 0
}

func (this *MyCircularDeque) IsFull() bool {
	return this.queue.Size() == this.k
}
