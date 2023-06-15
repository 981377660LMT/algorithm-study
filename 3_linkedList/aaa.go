// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2017 - ionous. Modified to create intrusive linked lists.

// Package inlist implements an intrusive doubly linked list.
//
// To iterate over a list (where l is a *List):
//	for e := l.Front(); e != nil; e = inlist.Next(e) {
//		// do something with e
//	}
//
package inlist

// Element mimics golang's container/list/Element; this implements the ilist/inlist/Intrusive interface to allow the storage of generic values in inlist.List.
// see also: NewElement(), and Value()
type Element struct {
	Hook
	Value interface{}
}

// NewElement creates an intrusive Element to store the passed value.
func NewElement(v interface{}) Intrusive {
	return &Element{Value: v}
}

// Value returns the contents of an element created via NewElement.
// This panics if i is nil, or of any type other than Element.
func Value(i Intrusive) interface{} {
	return i.(*Element).Value
}

// Hook provides an implementation of Intrusive with no additional data. Use this as an anonymous member of a struct to create an intrusive element.
// ( Alternatively, you can implement the Intrusive interface yourself. )
//
// For example:
// type MyElement struct {
// 	inlist.Hook
// 	MyData int // or whatever data you need.
// }
// l:= inlist.New()
// l.PushBack(&MyElement{MyData:23})
//
type Hook struct {
	prev, next Intrusive
	list       *List
}

// List implements Intrusive.List: returning the list pointer provided by SetElements().
func (s *Hook) List() *List {
	return s.list
}

// Predecessor implements Intrusive.Predecessor: returning the prev sibling provided by SetElements().
func (s *Hook) Predecessor() Intrusive {
	return s.prev
}

// Successor implements Intrusive.Successor: returning the next sibling provided by SetElements().
func (s *Hook) Successor() Intrusive {
	return s.next
}

// SetElements implements Intrusive.SetElements: storing the passed values without checking or changing them.
func (s *Hook) SetElements(l *List, p Intrusive, n Intrusive) {
	s.list, s.prev, s.next = l, p, n
}

// Intrusive provides an interface for intrusive linked list manipulation.
// Use inlist.Next(), inlist.Prev() for list traversal.
type Intrusive interface {
	// List, return the list given via SetElements.
	List() *List
	// Predecessor, return the prev element given via SetElements.
	Predecessor() Intrusive
	// Successor, return the next element given via SetElements.
	Successor() Intrusive
	// SetElements, record the passed values as given.
	SetElements(list *List, prev Intrusive, next Intrusive)
}

// Next returns the next list element or nil.
// Unlike e.Successor, this accounts for the list's sentinel.
func Next(e Intrusive) (ret Intrusive) {
	if list := e.List(); list != nil {
		if n := e.Successor(); n != &list.root {
			ret = n
		}
	}
	return
}

// Prev returns the previous list element or nil.
// Unlike e.Predecessor, this accounts for the list's sentinel.
func Prev(e Intrusive) (ret Intrusive) {
	if list := e.List(); list != nil {
		if p := e.Predecessor(); p != &list.root {
			ret = p
		}
	}
	return
}

// List is an intrusive doubly linked list, implemented as a ring.
// The zero value for List is empty and ready to use.
type List struct {
	root Hook // marker for start,end of the list
	cnt  int  // current list length excluding the sentinel
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.cnt = 0
	return l
}

// New returns an initialized list.
func New() *List {
	return new(List).Init()
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List) Len() int {
	return l.cnt
}

// Front returns the first element of list l or nil.
func (l *List) Front() (ret Intrusive) {
	if l.cnt != 0 {
		ret = l.root.next
	}
	return
}

// Back returns the last element of list l or nil.
func (l *List) Back() (ret Intrusive) {
	if l.cnt != 0 {
		ret = l.root.prev
	}
	return
}

// lazyInit lazily initializes a zero List value.
// Note: for any element where e.List() == l, lazyInit() has been called.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.cnt, and returns e.
func (l *List) insert(e, at Intrusive) Intrusive {
	// at  -> e -> n (old:at.next)
	n := at.Successor()
	at.SetElements(l, at.Predecessor(), e)
	e.SetElements(l, at, n)
	n.SetElements(l, e, n.Successor())
	l.cnt++
	return e
}

// remove removes e from its list, decrements l.cnt, and returns e.
func (l *List) remove(e Intrusive) Intrusive {
	// eprev -> enext
	eprev, enext := e.Predecessor(), e.Successor()
	eprev.SetElements(l, eprev.Predecessor(), enext)
	enext.SetElements(l, eprev, enext.Successor())
	e.SetElements(nil, nil, nil)
	l.cnt--
	return e
}

// Remove removes e from l if e is an element of list l.
// It returns true if it was removed.
func (l *List) Remove(e Intrusive) (okay bool) {
	if e.List() == l {
		l.remove(e)
		okay = true
	}
	return
}

// PushFront inserts e at the head of list l and returns e.
func (l *List) PushFront(e Intrusive) Intrusive {
	l.lazyInit()
	return l.insert(e, &l.root)
}

// PushBack inserts e at the tail of list l and returns e.
func (l *List) PushBack(e Intrusive) Intrusive {
	l.lazyInit()
	return l.insert(e, l.root.prev)
}

// InsertBefore inserts e directly ahead of mark and returns e.
// If mark is not in this list, no insertion occurs, and this returns nil.
func (l *List) InsertBefore(v Intrusive, mark Intrusive) (ret Intrusive) {
	if mark.List() == l {
		ret = l.insert(v, mark.Predecessor())
	}
	return
}

// InsertAfter inserts e directly following mark and returns e.
// If mark is not in this list, no insertion occurs, and this returns nil.
func (l *List) InsertAfter(v Intrusive, mark Intrusive) (ret Intrusive) {
	if mark.List() == l {
		ret = l.insert(v, mark)
	}
	return
}

// MoveToFront moves element e to the head the list, and returns true.
// If e is not in this list, the list is not modified, and this returns false.
func (l *List) MoveToFront(e Intrusive) (okay bool) {
	wrongList := e.List() != l || l.root.next == e
	if !wrongList {
		l.insert(l.remove(e), &l.root)
		okay = true
	}
	return
}

// MoveToBack puts e at the tail of the list, and returns true.
// If e is not in this list, the list is not modified, and this returns false.
func (l *List) MoveToBack(e Intrusive) (okay bool) {
	wrongList := e.List() != l || l.root.prev == e
	if !wrongList {
		l.insert(l.remove(e), l.root.prev)
		okay = true
	}
	return
}

// MoveBefore puts e directly in front of mark, and returns true.
// If e or mark is not in this list, or e == mark, the list is not modified, and this returns false.
func (l *List) MoveBefore(e, mark Intrusive) (okay bool) {
	wrongList := e.List() != l || e == mark || mark.List() != l
	if !wrongList {
		l.insert(l.remove(e), mark.Predecessor())
		okay = true
	}
	return
}

// MoveAfter puts element e directly behind mark, and returns true.
// If e or mark is not in this list, or e == mark, the list is not modified, and this returns false.
func (l *List) MoveAfter(e, mark Intrusive) (okay bool) {
	wrongList := e.List() != l || e == mark || mark.List() != l
	if !wrongList {
		l.insert(l.remove(e), mark)
		okay = true
	}
	return
}

// MoveBackList moves all elements from other to the end of this list.
func (l *List) MoveBackList(other *List) {
	if l != other {
		l.lazyInit()
		for e := other.Front(); e != nil; {
			n := Next(e)
			l.insert(e, l.root.prev)
			e = n
		}
		other.Init() // empty the other list
	}
}

// MoveFrontList moves all elements from other to the front of this list.
func (l *List) MoveFrontList(other *List) {
	if l != other {
		l.lazyInit()
		for e := other.Back(); e != nil; {
			p := Prev(e)
			l.insert(e, &l.root)
			e = p
		}
		other.Init() // empty the other list
	}
}
