// https://noshi91.github.io/Library/data_structure/partially_retroactive_queue.cpp
// PartiallyRetroactiveQueue
// 部分可追溯队列

package main

import (
	"fmt"
	"strings"
)

func main() {
	queue := NewPartiallyRetroactiveQueue[int]()
	fmt.Println(queue.Empty()) // true
	time1 := queue.Now()
	time2 := queue.InsertPush(time1, 1)
	time3 := queue.InsertPush(time2, 2)
	_ = time3
	fmt.Println(queue.Front()) // 2
	fmt.Println(queue.Empty()) // false
	queue.InsertPop()
	fmt.Println(queue.Front()) // 1
	queue.ErasePop()
	fmt.Println(queue.Front()) // 2

	time4 := queue.InsertPush(time1, 3)
	time5 := queue.InsertPush(time4, 4)
	_ = time5
}

// !NotVerified.
type PartiallyRetroactiveQueue[T any] struct {
	list  *List[T]
	front *Node[T]
}

func NewPartiallyRetroactiveQueue[T any]() *PartiallyRetroactiveQueue[T] {
	list := NewList[T]()
	front := list.Root
	return &PartiallyRetroactiveQueue[T]{list: list, front: front}
}

func (q *PartiallyRetroactiveQueue[T]) Now() *Node[T] {
	return q.list.Root
}

func (q *PartiallyRetroactiveQueue[T]) Empty() bool {
	return q.front == q.list.Root
}

func (q *PartiallyRetroactiveQueue[T]) Front() T {
	return q.front.Value
}

func (q *PartiallyRetroactiveQueue[T]) InsertPush(time *Node[T], x T) *Node[T] {
	it := q.list.InsertBefore(x, time)
	if it == q.list.Front() || !it.Prev.inQueue {
		it.inQueue = false
		q.front = q.front.Prev
		q.front.inQueue = true
	} else {
		it.inQueue = true
	}
	return it
}

func (q *PartiallyRetroactiveQueue[T]) ErasePush(time *Node[T]) {
	if time == q.list.Front() || !time.Prev.inQueue {
		q.front.inQueue = false
		q.front = q.front.Next
	}
	q.list.Remove(time)
}

func (q *PartiallyRetroactiveQueue[T]) InsertPop() {
	q.front.inQueue = false
	q.front = q.front.Next
}

func (q *PartiallyRetroactiveQueue[T]) ErasePop() {
	q.front = q.front.Prev
	q.front.inQueue = true
}

type Node[T any] struct {
	inQueue    bool
	Prev, Next *Node[T]
	Value      T
}

func (n *Node[T]) String() string {
	if n == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", n.Value)
}

type List[T any] struct {
	Root *Node[T] // sentinel
	len  int
}

func NewList[T any]() *List[T] {
	root := &Node[T]{}
	root.Prev = root
	root.Next = root
	return &List[T]{Root: root}
}

func (l *List[T]) Len() int { return l.len }

// 返回链表的第一个元素.
// 如果链表为空.则返回nil.
func (l *List[T]) Front() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.Root.Next
}

// 返回链表的最后一个元素.
// 如果链表为空.则返回nil.
func (l *List[T]) Back() *Node[T] {
	if l.len == 0 {
		return nil
	}
	return l.Root.Prev
}

func (l *List[T]) Remove(e *Node[T]) {
	if e == nil || e == l.Root {
		return
	}
	l.remove(e)
}

// 将v插入到链表的第一个元素之前，并返回新插入的元素.
func (l *List[T]) PushFront(v T) *Node[T] {
	return l.insertValueAfter(v, l.Root)
}

// 将v插入到链表的最后一个元素之后，并返回新插入的元素.
func (l *List[T]) PushBack(v T) *Node[T] {
	return l.insertValueAfter(v, l.Root.Prev)
}

func (l *List[T]) PopFront() T {
	if l.len == 0 {
		panic("list is empty")
	}
	front := l.Front()
	l.remove(front)
	return front.Value
}

func (l *List[T]) PopBack() T {
	if l.len == 0 {
		panic("list is empty")
	}
	back := l.Back()
	l.remove(back)
	return back.Value
}

// 将v插入到mark之前，并返回新插入的结点.
func (l *List[T]) InsertBefore(v T, mark *Node[T]) *Node[T] {
	return l.insertValueAfter(v, mark.Prev)
}

// 将v插入到mark之后，并返回新插入的结点.
func (l *List[T]) InsertAfter(v T, mark *Node[T]) *Node[T] {
	return l.insertValueAfter(v, mark)
}

func (l *List[T]) ForEach(f func(value T, index int)) {
	for i, e := 0, l.Front(); e != l.Root; i, e = i+1, e.Next {
		f(e.Value, i)
	}
}

func (l *List[T]) At(index int) *Node[T] {
	if index < 0 || index >= l.len {
		return nil
	}
	if index < l.len/2 {
		e := l.Root.Next
		for i := 0; i < index; i++ {
			e = e.Next
		}
		return e
	} else {
		e := l.Root.Prev
		for i := l.len - 1; i > index; i-- {
			e = e.Prev
		}
		return e
	}
}

func (l *List[T]) Prev(e *Node[T]) *Node[T] {
	if e == nil || e == l.Root {
		return nil
	}
	return e.Prev
}

func (l *List[T]) Next(e *Node[T]) *Node[T] {
	if e == nil || e == l.Root {
		return nil
	}
	return e.Next
}

func (l *List[T]) Clear() {
	l.Root.Next = l.Root
	l.Root.Prev = l.Root
	l.len = 0
}

func (l *List[T]) String() string {
	sb := []string{}
	l.ForEach(func(value T, _ int) {
		sb = append(sb, fmt.Sprintf("%v", value))
	})
	return fmt.Sprintf("List{%s}", strings.Join(sb, ","))
}

func (l *List[T]) insertValueAfter(v T, at *Node[T]) *Node[T] {
	return l.insertAfter(&Node[T]{Value: v}, at)
}

func (l *List[T]) insertAfter(e, at *Node[T]) *Node[T] {
	e.Prev = at
	e.Next = at.Next
	e.Prev.Next = e
	e.Next.Prev = e
	l.len++
	return e
}

func (l *List[T]) remove(e *Node[T]) {
	e.Prev.Next = e.Next
	e.Next.Prev = e.Prev
	e.Next = nil // avoid memory leaks
	e.Prev = nil // avoid memory leaks
	l.len--
}
