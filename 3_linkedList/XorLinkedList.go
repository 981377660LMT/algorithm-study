// 异或链表,利用计算机的的位异或操作(⊕)，来降低双向链表的存储需求。
// https://zh.wikipedia.org/wiki/%E5%BC%82%E6%88%96%E9%93%BE%E8%A1%A8
// https://www.cnblogs.com/lovemdx/p/3257334.html

// API:
// - NewXorLinkedList() *XorLinkedList
// - Append(value X)
// - AppendLeft(value X)
// - Pop()
// - PopLeft()
// - Front() X
// - Back() X
// - Empty() bool
// - Reverse()
// - Concat(left, right *XorLinkedList) *XorLinkedList

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	queue := NewXorLinkedList()
	queue.Append(1)
	queue.Append(2)
	queue.Append(2)
	queue.Append(2)
	queue.Append(10)
	fmt.Println(queue.Front(), queue.Back())
	queue.Reverse()
	fmt.Println(queue.Front(), queue.Back())
}

type X int
type XorLinkedList struct {
	front *xorNode
	back  *xorNode
	size  int
}

func NewXorLinkedList() *XorLinkedList {
	return &XorLinkedList{}
}

// 合并两个链表, 返回合并后的链表，原链表会被破坏.
func Concat(left, right *XorLinkedList) *XorLinkedList {
	if left.Empty() {
		return right
	}
	if right.Empty() {
		return left
	}
	left.back.link ^= ptoi(right.front)
	right.front.link ^= ptoi(left.back)
	res := &XorLinkedList{}
	res.front = left.front
	res.back = right.back
	res.size = left.size + right.size
	left.front = nil
	left.back = nil
	right.front = nil
	right.back = nil
	return res
}

func (list *XorLinkedList) Append(value X) {
	if list.Empty() {
		next := newXorNode(value, 0)
		list.front = next
		list.back = next
	} else {
		next := newXorNode(value, ptoi(list.back))
		list.back.link ^= ptoi(next)
		list.back = next
	}
	list.size++
}

func (list *XorLinkedList) AppendLeft(value X) {
	list.Reverse()
	list.Append(value)
	list.Reverse()
	list.size++
}

func (list *XorLinkedList) Pop() (res X) {
	if list.Empty() {
		panic("empty list")
	}

	res = list.back.value
	if list.front == list.back {
		list.front = nil
		list.back = nil
	} else {
		next := itop(list.back.link)
		next.link ^= ptoi(list.back)
		list.back = next
	}
	list.size--
	return
}

func (list *XorLinkedList) PopLeft() (res X) {
	res = list.front.value
	list.Reverse()
	list.Pop()
	list.Reverse()
	list.size--
	return
}

func (list *XorLinkedList) Empty() bool {
	return list.front == nil
}

func (list *XorLinkedList) Reverse() {
	list.front, list.back = list.back, list.front
}

func (list *XorLinkedList) Front() X {
	return list.front.value
}

func (list *XorLinkedList) Back() X {
	return list.back.value
}

func (list *XorLinkedList) Len() int {
	return list.size
}

type xorNode struct {
	value X
	link  uintptr
}

func newXorNode(value X, link uintptr) *xorNode {
	return &xorNode{value, link}
}

func ptoi(ptr *xorNode) uintptr {
	return uintptr(unsafe.Pointer(ptr))
}

func itop(id uintptr) *xorNode {
	return (*xorNode)(unsafe.Pointer(id))
}
