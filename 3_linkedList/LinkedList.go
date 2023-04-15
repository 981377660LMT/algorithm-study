package main

import (
	"fmt"
	"strings"
)

func main() {
	list := NewList()
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	fmt.Println(list)
}

type Node struct {
	next, prev *Node
	list       *List
	Value      interface{}
}

func (e *Node) Next() *Node {
	if p := e.next; e.list != nil && p != &e.list.Root {
		return p
	}
	return nil
}

func (e *Node) Prev() *Node {
	if p := e.prev; e.list != nil && p != &e.list.Root {
		return p
	}
	return nil
}

func (e *Node) String() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

type List struct {
	Root Node // sentinel list element, only &root, root.prev, and root.next are used
	len  int  // current list length excluding (this) sentinel element
}

func NewList() *List { return new(List).Clear() }

// Init.
func (l *List) Clear() *List {
	l.Root.next = &l.Root
	l.Root.prev = &l.Root
	l.len = 0
	return l
}

func (l *List) Len() int { return l.len }

// 返回链表的第一个元素，如果链表为空，则返回nil.
func (l *List) Front() *Node {
	if l.len == 0 {
		return nil
	}
	return l.Root.next
}

// 返回链表的最后一个元素，如果链表为空，则返回nil.
func (l *List) Back() *Node {
	if l.len == 0 {
		return nil
	}
	return l.Root.prev
}

// 删除链表中的e元素，如果e不在链表中，则不做任何操作.
//  返回e的值.
func (l *List) Remove(e *Node) interface{} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// 将v插入到链表的第一个元素之前，并返回新插入的元素.
func (l *List) PushFront(v interface{}) *Node {
	l.lazyInit()
	return l.insertValue(v, &l.Root)
}

// 将v插入到链表的最后一个元素之后，并返回新插入的元素.
func (l *List) PushBack(v interface{}) *Node {
	l.lazyInit()
	return l.insertValue(v, l.Root.prev)
}

// 将v插入到mark之前，并返回新插入的结点.
func (l *List) InsertBefore(v interface{}, mark *Node) *Node {
	if mark.list != l {
		return nil
	}

	return l.insertValue(v, mark.prev)
}

// 将v插入到mark之后，并返回新插入的结点.
func (l *List) InsertAfter(v interface{}, mark *Node) *Node {
	if mark.list != l {
		return nil
	}
	return l.insertValue(v, mark)
}

// 将e移动到链表的头部.
func (l *List) MoveToFront(e *Node) {
	if e.list != l || l.Root.next == e {
		return
	}
	l.move(e, &l.Root)
}

// 将e移动到链表的尾部.
func (l *List) MoveToBack(e *Node) {
	if e.list != l || l.Root.prev == e {
		return
	}

	l.move(e, l.Root.prev)
}

// 将e移动到mark之前.
func (l *List) MoveBefore(e, mark *Node) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// 将e移动到mark之后.
func (l *List) MoveAfter(e, mark *Node) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.Root.prev)
	}
}

func (l *List) PushFrontList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.Root)
	}
}

func (l *List) ForEach(f func(value interface{}, index int)) {
	for i, e := 0, l.Front(); e != nil; e, i = e.Next(), i+1 {
		f(e.Value, i)
	}
}

func (l *List) String() string {
	sb := []string{}
	l.ForEach(func(value interface{}, _ int) {
		sb = append(sb, fmt.Sprintf("%v", value))
	})
	return fmt.Sprintf("List{%s}", strings.Join(sb, ","))
}

func (l *List) lazyInit() {
	if l.Root.next == nil {
		l.Clear()
	}
}

func (l *List) insert(e, at *Node) *Node {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

func (l *List) insertValue(v interface{}, at *Node) *Node {
	return l.insert(&Node{Value: v}, at)
}

func (l *List) remove(e *Node) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

func (l *List) move(e, at *Node) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}
