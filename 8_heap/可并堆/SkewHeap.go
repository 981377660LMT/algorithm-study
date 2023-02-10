// https://ei1333.github.io/library/structure/heap/skew-heap.hpp

package main

import "fmt"

func main() {
	manager := NewSkewHeap(true)
	heaps := make([]*SkewHeapNode, 10) // 这些堆
	for i := range heaps {
		newHeap := manager.Push(nil, E(i), i)
		heaps[i] = newHeap
	}

	fmt.Println(heaps[1].value)
	manager.Add(heaps[1], 10)
	fmt.Println(heaps[1].value)
	newRoot := manager.Meld(heaps[2], heaps[1])
	fmt.Println(newRoot.value)

}

type E = int

type SkewHeapNode struct {
	value, lazy E
	left, right *SkewHeapNode
	index       int
}

type SkewHeap struct {
	isMin bool
}

func NewSkewHeap(isMin bool) *SkewHeap {
	return &SkewHeap{isMin: isMin}
}

// 将(key,index)插入堆中，返回插入后的堆
func (sk *SkewHeap) Push(t *SkewHeapNode, key E, index int) *SkewHeapNode {
	return sk.Meld(t, newNode(key, index))
}

// 删除堆顶元素,返回删除后的堆
func (sk *SkewHeap) Pop(t *SkewHeapNode) *SkewHeapNode {
	return sk.Meld(t.left, t.right)
}

func (sk *SkewHeap) Top(t *SkewHeapNode) E {
	return t.value
}

// 将x与y合并，返回合并后的堆
func (sk *SkewHeap) Meld(x, y *SkewHeapNode) *SkewHeapNode {
	sk.propagate(x)
	sk.propagate(y)
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}
	if (x.value < y.value) != sk.isMin {
		x, y = y, x
	}
	x.right = sk.Meld(y, x.right)
	x.left, x.right = x.right, x.left
	return x
}

func (sk *SkewHeap) Add(t *SkewHeapNode, lazy E) *SkewHeapNode {
	if t == nil {
		return t
	}
	t.lazy += lazy
	sk.propagate(t)
	return t
}

func (sk *SkewHeap) propagate(t *SkewHeapNode) *SkewHeapNode {
	if t != nil && t.lazy != 0 {
		if t.left != nil {
			t.left.lazy += t.lazy
		}
		if t.right != nil {
			t.right.lazy += t.lazy
		}
		t.value += t.lazy
		t.lazy = 0
	}
	return t
}

func newNode(key E, index int) *SkewHeapNode {
	return &SkewHeapNode{value: key, index: index}
}
