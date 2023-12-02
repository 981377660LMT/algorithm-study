// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/data_structure/order_maintenance.cc#L1
//
// Order Maintenance
//
// - alloc(): return new node x
// - insertAfter(x, y): insert node y after node x
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

// !API:
//  Alloc() *Node
//  Build(nums []int) // 用nums中的元素构建链表.
//  InsertAfter(x, y *Node) // 将y插入到x后面.
//  IsBefore(x, y *Node) bool // 判断x是否在y前面.
//  Erase(x *Node) // 删除x.

// !维护元素顺序的链表/带插入全序集维护.
// 用于维护元素的先后顺序, 以及判断元素是否在另一个元素的前面.

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	n := 1000

	// check with perm
	perm := rand.Perm(n)
	mp := make(map[int]int) // (value->index)
	for i, v := range perm {
		mp[v] = i
	}

	om := NewOrderMaintenace()
	nodes := om.Build(perm) // 每个节点代表0-n-1中每个元素.

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			res1 := om.IsBefore(nodes[i], nodes[j])
			num1, num2 := perm[i], perm[j]
			res2 := mp[num1] < mp[num2]
			if res1 != res2 {
				fmt.Println(perm, res1, res2)
				panic(fmt.Sprintf("error: %d %d", num1, num2))
			}
		}
	}
	fmt.Println("pass")

}

const M uint64 = ^(^uint64(0) >> 1)

type Node struct {
	prev, next *Node
	label      uint64
}

type OrderMaintenace struct {
	Head *Node
}

func NewOrderMaintenace() *OrderMaintenace {
	res := &OrderMaintenace{}
	root := &Node{}
	root.next = root
	root.prev = root
	res.Head = root
	return res
}

func (om *OrderMaintenace) Alloc() *Node {
	return &Node{}
}

func (om *OrderMaintenace) Build(nums []int) []*Node {
	n := len(nums)
	res := make([]*Node, n)
	for i := 0; i < n; i++ {
		res[i] = om.Alloc()
	}
	pre := om.Head
	for _, cur := range res {
		om.InsertAfter(pre, cur)
		cur = pre
	}
	return res
}

// 将y插入到x后面.
func (om *OrderMaintenace) InsertAfter(x, y *Node) {
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
		gap := M
		if x != end {
			gap = om._width(x, end) / required
		}
		val := end.label
		for {
			if end == om.Head {
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
func (om *OrderMaintenace) Erase(x *Node) {
	x.prev.next = x.next.prev
	x.next.prev = x.prev.next
}

// 判断x是否在y前面.
func (om *OrderMaintenace) IsBefore(x, y *Node) bool {
	return x.label > y.label
}

func (om *OrderMaintenace) _width(x, y *Node) uint64 {
	res := y.label - x.label
	if res-1 >= M {
		res += M
	}
	return res
}
