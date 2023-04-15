// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/data_structure/order_maintenance.cc#L1
//
// Order Maintenance
//
// - create_node(): return new node x
// - insert(x, y): insert node y after node x
// - erase(x): erase node x from the list
// - order(x, y): return true if x is before y
//
// Running Time:
//   worst case O(1) for create_node, erase, and order.
//   amortized O(log n) for insert; very small constant.
//
// Reference:
//   P. Dietz and D. Sleator (1987):
//   "Two algorithms for maintaining order in a list".
//   STOC.
//

// API:
//  Alloc() *Node
//  InsertAfter(x, y *Node) // 将y插入到x后面.
//  IsBefore(x, y *Node) bool // 判断x是否在y前面.
//  Erase(x *Node) // 删除x.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	M := NewOrderMaintenace()
	n := 10000

	nodes := make([]*_Node, n)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		nodes[i] = M.Alloc()
		order[i] = i
	}
	rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

	time1 := time.Now()
	M.InsertAfter(M.Root, nodes[order[0]])
	for i := 1; i < n; i++ {
		M.InsertAfter(nodes[order[i-1]], nodes[order[i]])
	}
	fmt.Println(time.Since(time1))
}

const M uint64 = ^(^uint64(0) >> 1)

type _Node struct {
	prev, next *_Node
	label      uint64
}

type OrderMaintenace struct {
	Root *_Node
}

func NewOrderMaintenace() *OrderMaintenace {
	res := &OrderMaintenace{}
	head := &_Node{}
	head.next = head
	head.prev = head
	res.Root = head
	return res
}

func (om *OrderMaintenace) Alloc() *_Node {
	return &_Node{}
}

// 将y插入到x后面.
func (om *OrderMaintenace) InsertAfter(x, y *_Node) {
	label := x.label
	if om._width(x, x.next) <= 1 {
		mid := x.next
		end := mid.next
		required := uint64(3)
		for om._width(x, end) <= 4*om._width(x, mid) && end != x {
			required++
			end = end.next
			if end == x {
				break
			}
			required++
			end = end.next
			mid = mid.next
		}
		var gap uint64
		if x == end {
			gap = M / required
		} else {
			gap = om._width(x, end) / required
		}
		val := end.label
		for {
			if end == om.Root {
				val += M
			}
			end = end.prev
			if end == x {
				break
			}
			val -= gap
			end.label = val
		}
	}
	y.label = label + om._width(x, x.next)/2
	y.next = x.next
	y.prev = x
	y.next.prev = y
	y.prev.next = y
}

// 从list中删除x.
func (om *OrderMaintenace) Erase(x *_Node) {
	x.prev.next = x.next.prev
	x.next.prev = x.prev.next
}

// 判断x是否在y前面.
func (om *OrderMaintenace) IsBefore(x, y *_Node) bool {
	return x.label < y.label
}

func (om *OrderMaintenace) _width(x, y *_Node) uint64 {
	res := y.label - x.label
	if res-1 >= M {
		res += M
	}
	return res
}
