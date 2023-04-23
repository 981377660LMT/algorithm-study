//
// Burkhard-Keller Tree (metric tree)
//
// Description:
//   Let V be a (finite) set, and d: V x V -> R be a metric.
//   BK tree supports the following operations:
// !    - insert(p): insert a point p, O((log n)^2)
// !    - traverse(p,d): enumerate all q with d(p,q) <= d
//
// Remark:
// !  To delete elements and/or rebalance the tree,
// !  we can use the same technique as the scapegoat tree(替罪羊树暴力重构).
//
// Reference
//   W. Burkhard and R. Keller (1973):
//   Some approaches to best-match file searching,
//   Communications of the ACM, vol. 16, issue. 4, pp. 230--236.
//
// https://zhuanlan.zhihu.com/p/360925212
// https://www.twblogs.net/a/5c22542bbd9eee16b4a778a2
// https://yfsyfs.github.io/2019/06/25/%E7%BA%A0%E9%94%99%E5%88%A9%E5%99%A8-BK%E6%A0%91/
// BKTree用于拼写纠错

package main

import "fmt"

func main() {
	bk := NewBKTree(func(a, b P) int {
		return abs(a[0]-b[0]) + abs(a[1]-b[1])
	})
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			bk.Insert(P{i, j})
		}
	}

	bk.Enumerate(P{0, 0}, 2, func(p P) {
		fmt.Println(p)
	})
}

type P [2]int // [x,y]

type BKTree struct {
	root    *BNode
	calDist func(a, b P) int
}

func NewBKTree(calDist func(a, b P) int) *BKTree {
	return &BKTree{calDist: calDist}
}

// Inserts a point, O((log n)^2).
func (t *BKTree) Insert(point P) {
	t.root = t._insert(t.root, point)
}

// Enumerates all q with d(p,q) <= maxDist.
func (t *BKTree) Enumerate(point P, maxDist int, f func(point P)) {
	t._traverse(t.root, point, maxDist, f)
}

func (t *BKTree) _insert(node *BNode, point P) *BNode {
	if node == nil {
		return &BNode{p: point, ch: map[int]*BNode{}}
	}
	d := t.calDist(node.p, point)
	node.ch[d] = t._insert(node.ch[d], point)
	return node
}

func (t *BKTree) _traverse(node *BNode, point P, maxDist int, f func(point P)) {
	if node == nil {
		return
	}
	d := t.calDist(node.p, point)
	if d <= maxDist {
		f(node.p)
	}
	for k := range node.ch {
		if -maxDist <= k-d && k-d <= maxDist {
			t._traverse(node.ch[k], point, maxDist, f)
		}
	}
}

type BNode struct {
	p  P
	ch map[int]*BNode
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
