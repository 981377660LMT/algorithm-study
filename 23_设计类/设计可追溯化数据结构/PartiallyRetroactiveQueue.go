// TODO 有错误
package main

import (
	"fmt"
	"strings"
)

func main() {
	queue := NewPartiallyRetroactiveQueue()
	time1 := queue.Now()
	time2 := queue.InsertPush(time1, 1)
	fmt.Println(queue.Front())
	time3 := queue.InsertPush(time2, 2)
	fmt.Println(queue.Front(), queue.list)
	queue.InsertPush(time3, 3)
	fmt.Println(queue.Front(), queue.list)
}

type S int

// 部分可追溯化队列.
//  !可以Insert/Erase任意时间点的`入队`操作, 以及Insert/Erase当前时间点的`出队`操作.
type PartiallyRetroactiveQueue struct {
	list      *SimpleList
	frontIter *Time
}

type NodeType struct {
	queue   *PartiallyRetroactiveQueue
	value   S
	inQueue bool
}

func NewPartiallyRetroactiveQueue() *PartiallyRetroactiveQueue {
	res := &PartiallyRetroactiveQueue{}
	list := NewSimpleList()
	res.list = list
	res.frontIter = &list.Root
	return res
}

func (q *PartiallyRetroactiveQueue) Now() *Time {
	return q.list.Back()
}

func (q *PartiallyRetroactiveQueue) Empty() bool {
	return q.frontIter == &q.list.Root
}

func (q *PartiallyRetroactiveQueue) Front() S {
	if q.Empty() {
		panic("cannot get the front of an empty queue")
	}
	return q.frontIter.Value.value
}

func (q *PartiallyRetroactiveQueue) InsertPush(time *Time, x S) *Time {
	node := &NodeType{queue: q, value: x}
	iter := q.list.InsertBefore(node, time)
	if iter == q.list.Front() || !iter.Prev().Value.inQueue {
		iter.Value.inQueue = false
		q.frontIter = q.frontIter.Prev()
		q.frontIter.Value.inQueue = true
	} else {
		iter.Value.inQueue = true
	}
	return iter
}

func (q *PartiallyRetroactiveQueue) ErasePush(t *Time) {
	if t == q.Now() {
		panic("cannot erase the last element")
	}
	if t == q.list.Front() || !t.Prev().Value.inQueue {
		q.frontIter.Value.inQueue = false
		q.frontIter = q.frontIter.Next()
	}
	q.list.Remove(t)
}

func (q *PartiallyRetroactiveQueue) InsertPop() {
	if q.Empty() {
		panic("cannot pop an empty queue")
	}
	q.frontIter.Value.inQueue = false
	q.frontIter = q.frontIter.Next()
}

func (q *PartiallyRetroactiveQueue) ErasePop() {
	if q.frontIter == q.list.Front() {
		panic("cannot erase the first element")
	}
	q.frontIter = q.frontIter.Prev()
	q.frontIter.Value.inQueue = true
}

func (q *PartiallyRetroactiveQueue) String() string {
	return fmt.Sprintf("PartiallyRetroactiveQueue{%s}", q.list)
}

type Time struct {
	next, prev *Time
	Value      *NodeType
}

func (e *Time) Next() *Time {
	return e.next
}

func (e *Time) Prev() *Time {
	return e.prev
}

func (e *Time) String() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

type SimpleList struct {
	Root Time // sentinel list element, only &root, root.prev, and root.next are used
}

func NewSimpleList() *SimpleList { return new(SimpleList).Clear() }

// Init.
func (l *SimpleList) Clear() *SimpleList {
	l.Root.next = &l.Root
	l.Root.prev = &l.Root
	return l
}

func (l *SimpleList) Front() *Time {
	return l.Root.next
}

func (l *SimpleList) Back() *Time {
	return l.Root.prev
}

func (l *SimpleList) Remove(e *Time) *NodeType {
	l.remove(e)
	return e.Value
}

func (l *SimpleList) InsertBefore(v *NodeType, mark *Time) *Time {
	return l.insertValue(v, mark.prev)
}

func (l *SimpleList) ForEach(f func(value *NodeType, index int)) {
	for i, e := 0, l.Front(); e != &l.Root; e, i = e.Next(), i+1 {
		f(e.Value, i)
	}
}

func (l *SimpleList) String() string {
	sb := []string{}
	l.ForEach(func(value *NodeType, _ int) {
		sb = append(sb, fmt.Sprintf("%v", value.value))
	})
	return fmt.Sprintf("List{%s}", strings.Join(sb, ","))
}

func (l *SimpleList) insert(e, at *Time) *Time {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	return e
}

func (l *SimpleList) insertValue(v *NodeType, at *Time) *Time {
	return l.insert(&Time{Value: v}, at)
}

func (l *SimpleList) remove(e *Time) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
}
