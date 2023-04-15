// 可持久化堆(随机堆)

//
// Persistent Heap (based on randomized meldable heap)
//
// Description:
//   Meldable heap with O(1) copy.
//
// Algorithm:
//   It is a persistence version of randomized meldable heap.
//
// Complexity:
//   O(log n) time/space for each operations.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	heapq := NewPersistentHeap(func(v1, v2 H) bool { return v1 < v2 })
	nodes := heapq.Build([]H{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	time1 := time.Now()
	for i := 0; i < 200000; i++ {
		nodes[0] = heapq.Push(nodes[0], i)
	}
	fmt.Println(time.Since(time1))
	fmt.Println(heapq.Top(nodes[0]))
	nodes[0] = heapq.Pop(nodes[0])
	fmt.Println(heapq.Top(nodes[0]))
}

type H = int
type PersistentHeap struct {
	less func(H, H) bool
}

func NewPersistentHeap(less func(H, H) bool) *PersistentHeap {
	return &PersistentHeap{less: less}
}

type RNode struct {
	value       H
	left, right *RNode
}

func (heap *PersistentHeap) Alloc() *RNode {
	return &RNode{}
}

func (heap *PersistentHeap) Build(nums []H) []*RNode {
	res := make([]*RNode, len(nums))
	for i, num := range nums {
		res[i] = &RNode{value: num}
	}
	return res
}

func (heap *PersistentHeap) Meld(node1, node2 *RNode) *RNode {
	if node1 == nil || node2 == nil {
		if node1 == nil {
			return node2
		}
		return node1
	}
	if heap.less(node2.value, node1.value) {
		node1, node2 = node2, node1
	}
	if rand.Intn(2)&1 == 0 {
		return &RNode{value: node1.value, left: heap.Meld(node1.left, node2), right: node1.right}
	}
	return &RNode{value: node1.value, left: node1.left, right: heap.Meld(node1.right, node2)}
}

func (heap *PersistentHeap) Push(node *RNode, v H) *RNode {
	return heap.Meld(node, &RNode{value: v})
}

func (heap *PersistentHeap) Pop(node *RNode) *RNode {
	return heap.Meld(node.left, node.right)
}

func (heap *PersistentHeap) Top(node *RNode) H {
	return node.value
}
