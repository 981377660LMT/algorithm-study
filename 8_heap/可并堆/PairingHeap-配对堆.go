// https://noshi91.github.io/Library/data_structure/pairing_heap.cpp
// https://scrapbox.io/data-structures/Pairing_Heap
// https://oi-wiki.org/ds/pairing-heap/

// 配对堆是一棵满足堆性质的带权多叉树（如下图），即每个节点的权值都小于或等于他的所有儿子
// 一个节点的所有儿子节点形成一个单向链表。每个节点储存第一个儿子的指针，即链表的头节点；和他的右兄弟的指针。
package main

import (
	"fmt"
)

func main() {
	pq1 := NewPairingHeap(func(a, b int) bool { return a < b })
	pq1.Push(1)
	pq1.Push(2)
	pq2 := NewPairingHeap(func(a, b int) bool { return a < b })
	pq2.Push(3)
	pq2.Push(4)
	pq3 := Meld(pq1, pq2)
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Pop())
	fmt.Println(pq3.Empty())
}

type P = int
type PairingHeap struct {
	root *pNode
	less func(P, P) bool
}

type pNode struct {
	value P
	head  *pNode
	next  *pNode // sibling
}

func NewPairingHeap(less func(P, P) bool) *PairingHeap {
	return &PairingHeap{less: less}
}

// 融合两个堆, 返回新的堆, 原来的堆被破坏.
//  !注意两个堆的比较函数必须相同.
func Meld(h1, h2 *PairingHeap) *PairingHeap {
	return &PairingHeap{root: _merge(h1.root, h2.root, h1.less), less: h1.less}
}

func (h *PairingHeap) Empty() bool {
	return h.root == nil
}

func (h *PairingHeap) Top() P {
	if h.Empty() {
		panic("empty heap")
	}
	return h.root.value
}

func (h *PairingHeap) Push(x P) {
	h.root = _merge(h.root, &pNode{value: x}, h.less)
}

func (h *PairingHeap) Pop() P {
	if h.Empty() {
		panic("empty heap")
	}
	res := h.root.value
	h.root = _mergeList(h.root.head, h.less)
	return res
}

func _merge(x, y *pNode, less func(P, P) bool) *pNode {
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}

	if !less(x.value, y.value) {
		x, y = y, x
	}
	y.next = x.head
	x.head = y
	return x
}

func _mergeList(list *pNode, less func(P, P) bool) *pNode {
	if list == nil || list.next == nil {
		return list
	}
	next := list.next
	rem := next.next
	return _merge(_merge(list, next, less), _mergeList(rem, less), less)
}
