package linkedlist

import (
	"cmnx/src/deque"
	"fmt"
)

// Assert Deque implementation
var _ deque.Deque = (*LinkedList)(nil)

func demo() {
	linkedList := NewLinkedList()
	linkedList.Append(1)
	linkedList.Append(2)
	linkedList.Append(3)
	fmt.Println(linkedList.At(0))
	linkedList.Pop()
	linkedList.ForEach(func(value interface{}, index int) {
		fmt.Println(value, index)
	})
}

func NewLinkedList() *LinkedList {
	root := NewLinkedListNode(nil, nil, nil)
	root.Pre = root
	root.Next = root
	return &LinkedList{root, 0}
}

type LinkedList struct {
	root *LinkedListNode
	size int
}

func (list *LinkedList) Append(value interface{}) {
	node := NewLinkedListNode(value, nil, nil)
	list.root.InsertBefore(node)
	list.size++
}

func (list *LinkedList) Pop() (value interface{}) {
	if list.size == 0 {
		return
	}

	value = list.root.Pre.Remove().Value
	list.size--
	return

}

func (list *LinkedList) AppendLeft(value interface{}) {
	node := NewLinkedListNode(value, nil, nil)
	list.root.InsertAfter(node)
	list.size++
}

func (list *LinkedList) PopLeft() (value interface{}) {
	if list.size == 0 {
		return
	}

	value = list.root.Next.Remove().Value
	list.size--
	return
}

func (list *LinkedList) Len() int {
	return list.size
}

func (list *LinkedList) At(index int) (value interface{}) {
	if index < 0 {
		index += list.size
	}

	if index < 0 || index >= list.size {
		return
	}

	node := list.root.Next
	for i := 0; i < index; i++ {
		node = node.Next
	}
	value = node.Value
	return
}

func (list *LinkedList) ForEach(f func(value interface{}, index int)) {
	node := list.root.Next
	for i := 0; i < list.size; i++ {
		f(node.Value, i)
		node = node.Next
	}
}

func NewLinkedListNode(value interface{}, pre, next *LinkedListNode) *LinkedListNode {
	return &LinkedListNode{
		Value: value,
	}
}

type LinkedListNode struct {
	Value interface{}
	Pre   *LinkedListNode
	Next  *LinkedListNode
}

// 在other之后插入新节点 并返回新节点
func (cur *LinkedListNode) InsertAfter(other *LinkedListNode) *LinkedListNode {
	other.Pre = cur
	other.Next = cur.Next
	other.Pre.Next = other
	if other.Next != nil {
		other.Next.Pre = other
	}
	return other
}

// 在other之前插入新节点 并返回新节点
func (cur *LinkedListNode) InsertBefore(other *LinkedListNode) *LinkedListNode {
	other.Pre = cur.Pre
	other.Next = cur
	other.Next.Pre = other
	if other.Pre != nil {
		other.Pre.Next = other
	}
	return other
}

// 从链表里移除自身
func (cur *LinkedListNode) Remove() *LinkedListNode {
	if cur.Pre != nil {
		cur.Pre.Next = cur.Next
	}
	if cur.Next != nil {
		cur.Next.Pre = cur.Pre
	}
	cur.Pre = nil
	cur.Next = nil
	return cur
}
