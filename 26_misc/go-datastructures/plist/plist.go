/*
Package list provides list implementations. Currently, this includes a
persistent, immutable linked list.
*/
package main

import (
	"errors"
	"fmt"
)

func main() {
	// 1) 从空链表开始
	pl := Empty

	// 2) Add
	pl = pl.Add("A") // 头插入 "A"
	pl = pl.Add("B") // 头插入 "B"
	pl = pl.Add("C") // 头插入 "C"

	// 链表现在顺序: [C, B, A]
	headVal, _ := pl.Head()
	fmt.Println("Head:", headVal) // C

	// 3) Insert
	// 在位置1插入 "D" => 位置0是 "C", 位置1是 "B"
	pl, err := pl.Insert("D", 1)
	if err != nil {
		panic(err)
	}
	// 新链表: [C, D, B, A]

	// 4) Get
	val, ok := pl.Get(2)
	fmt.Println("Index2:", val, "OK:", ok) // Index2: B OK: true

	// 5) Remove
	pl2, err := pl.Remove(1) // 移除位置1 => "D"
	// pl2: [C, B, A]
	headVal2, _ := pl2.Head()
	fmt.Println("After remove, head:", headVal2) // C

	// pl 仍然是 [C, D, B, A], pl2 是 [C, B, A]
	// 不可变 => 原来的列表 pl 未改变

	// 7) Find / FindIndex
	foundVal, foundOk := pl.Find(func(x interface{}) bool {
		return x == "B"
	})
	idx := pl.FindIndex(func(x interface{}) bool {
		return x == "B"
	})
	fmt.Println("Found B:", foundVal, foundOk, "Index:", idx)

	// 8) IsEmpty, Length
	fmt.Println("IsEmpty:", pl.IsEmpty()) // false
	fmt.Println("Length:", pl.Length())   // 4

	// 9) ForEach
	pl.ForEach(func(x interface{}) { fmt.Println("ForEach:", x) })
}

var (
	// Empty is an empty PersistentList.
	Empty PersistentList = &emptyList{}

	// ErrEmptyList is returned when an invalid operation is performed on an
	// empty list.
	ErrEmptyList = errors.New("Empty list")
)

// PersistentList is an immutable, persistent linked list.
type PersistentList interface {
	// Head returns the head of the list. The bool will be false if the list is
	// empty.
	Head() (any, bool)

	// Tail returns the tail of the list. The bool will be false if the list is
	// empty.
	Tail() (PersistentList, bool)

	// IsEmpty indicates if the list is empty.
	IsEmpty() bool

	// Length returns the number of items in the list.
	Length() uint

	// Add will add the item to the list, returning the new list.
	Add(head any) PersistentList

	// Insert will insert the item at the given position, returning the new
	// list or an error if the position is invalid.
	Insert(val any, pos uint) (PersistentList, error)

	// Get returns the item at the given position or an error if the position
	// is invalid.
	Get(pos uint) (any, bool)

	// Remove will remove the item at the given position, returning the new
	// list or an error if the position is invalid.
	Remove(pos uint) (PersistentList, error)

	// Find applies the predicate function to the list and returns the first
	// item which matches.
	Find(func(any) bool) (any, bool)

	// FindIndex applies the predicate function to the list and returns the
	// index of the first item which matches or -1 if there is no match.
	FindIndex(func(any) bool) int

	ForEach(func(any))
}

type emptyList list

func (e *emptyList) Head() (any, bool)            { return nil, false }
func (e *emptyList) Tail() (PersistentList, bool) { return nil, false }
func (e *emptyList) IsEmpty() bool                { return true }
func (e *emptyList) Length() uint                 { return 0 }
func (e *emptyList) Add(head any) PersistentList  { return &list{head, e} }
func (e *emptyList) Insert(val any, pos uint) (PersistentList, error) {
	if pos == 0 {
		return e.Add(val), nil
	}
	return nil, ErrEmptyList
}
func (e *emptyList) Get(pos uint) (any, bool)                { return nil, false }
func (e *emptyList) Remove(pos uint) (PersistentList, error) { return nil, ErrEmptyList }
func (e *emptyList) Find(func(any) bool) (any, bool)         { return nil, false }
func (e *emptyList) FindIndex(func(any) bool) int            { return -1 }
func (e *emptyList) ForEach(f func(any))                     {}

type list struct {
	head any
	tail PersistentList
}

// Head returns the head of the list. The bool will be false if the list is
// empty.
func (l *list) Head() (any, bool) {
	return l.head, true
}

// Tail returns the tail of the list. The bool will be false if the list is
// empty.
func (l *list) Tail() (PersistentList, bool) {
	return l.tail, true
}

// IsEmpty indicates if the list is empty.
func (l *list) IsEmpty() bool {
	return false
}

// Length returns the number of items in the list.
func (l *list) Length() uint {
	curr := l
	length := uint(0)
	for {
		length += 1
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return length
		}
		curr = tail.(*list)
	}
}

func (l *list) ForEach(f func(any)) {
	curr := l
	for {
		f(curr.head)
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return
		}
		curr = tail.(*list)
	}
}

// Add will add the item to the list, returning the new list.
func (l *list) Add(head any) PersistentList {
	return &list{head, l}
}

// Insert will insert the item at the given position, returning the new list or
// an error if the position is invalid.
func (l *list) Insert(val any, pos uint) (PersistentList, error) {
	if pos == 0 {
		return l.Add(val), nil
	}
	nl, err := l.tail.Insert(val, pos-1)
	if err != nil {
		return nil, err
	}
	return nl.Add(l.head), nil
}

// Get returns the item at the given position or an error if the position is
// invalid.
func (l *list) Get(pos uint) (any, bool) {
	if pos == 0 {
		return l.head, true
	}
	return l.tail.Get(pos - 1)
}

// Remove will remove the item at the given position, returning the new list or
// an error if the position is invalid.
func (l *list) Remove(pos uint) (PersistentList, error) {
	if pos == 0 {
		nl, _ := l.Tail()
		return nl, nil
	}

	nl, err := l.tail.Remove(pos - 1)
	if err != nil {
		return nil, err
	}
	return &list{l.head, nl}, nil
}

// Find applies the predicate function to the list and returns the first item
// which matches.
func (l *list) Find(pred func(any) bool) (any, bool) {
	if pred(l.head) {
		return l.head, true
	}
	return l.tail.Find(pred)
}

// FindIndex applies the predicate function to the list and returns the index
// of the first item which matches or -1 if there is no match.
func (l *list) FindIndex(pred func(any) bool) int {
	curr := l
	idx := 0
	for {
		if pred(curr.head) {
			return idx
		}
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return -1
		}
		curr = tail.(*list)
		idx += 1
	}
}
