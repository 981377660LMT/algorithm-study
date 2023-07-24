// Reference:
//  https://maspypy.github.io/library/ds/randomized_bst/rbst_acted_monoid.hpp
//
//
// !幺半群上的RBST
//
//  分裂/拼接api:
//   1. Merge(left, right) -> root
//   2. Add(root, node) -> root
//   3. SplitByRank(root, k) -> [0,k) and [k,n)
//   4. SplitByValue(root, v) -> （-inf,v) and [v,inf)
//
//  查询/更新api:
//   1. Query(node, start, end) -> res
//   2. QueryAll(node) -> res
//   3. Update(node, start, end, lazy)
//   4. AllApply(node, lazy) -> node
//
//  构建api:
//   1. NewRoot() -> root
//   2. NewNode(v) -> node
//
//  操作api:
//   1. Toggle(node)
//   2. Reverse(node, start, end)
//   3. Size(node) -> size
//   Pop/Erase/At/BisectLeft/BisectRight... 都是基于分裂/拼接实现的

package main

import (
	"time"
)

const INF int = 1e18

// RangeAddRangeMax
type E = int
type Id = int

func rev(e E) E              { return e }
func e() E                   { return 0 }
func id() Id                 { return 0 }
func op(e1, e2 E) E          { return max(e1, e2) }
func mapping(f Id, e E) E    { return f + e }
func composition(f, g Id) Id { return f + g }
func less(e1, e2 E) bool     { return e1 > e2 } // !维护最大值

//
//
//

// 每个结点代表一段区间
type Node struct {
	left, right *Node
	value       E
	data        E
	lazy        Id
	size        int
	priority    uint64
	isReversed  bool
}

func NewRoot() *Node {
	return nil
}

func NewNode(v E) *Node {
	res := &Node{value: v, data: v, size: 1, lazy: id(), priority: _nextRand()}
	return res
}

// 合并两棵树, 保证Value有序.
func Add(root, node *Node) *Node {
	if node == nil {
		return root
	}
	left, right := SplitByValue(root, node.value)
	return Merge(Merge(left, node), right)
}

// 合并`左右`两棵树，保证Rank有序.
func Merge(left, right *Node) *Node {
	if left == nil || right == nil {
		if left == nil {
			return right
		}
		return left
	}

	if left.priority > right.priority {
		_pushDown(left)
		left.right = Merge(left.right, right)
		return _pushUp(left)
	} else {
		_pushDown(right)
		right.left = Merge(left, right.left)
		return _pushUp(right)
	}
}

// split root to [0,k) and [k,n)
func SplitByRank(root *Node, k int) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	_pushDown(root)
	leftSize := Size(root.left)
	if k <= leftSize {
		first, second := SplitByRank(root.left, k)
		root.left = second
		return first, _pushUp(root)
	} else {
		first, second := SplitByRank(root.right, k-leftSize-1)
		root.right = first
		return _pushUp(root), second
	}
}

// split root to `less than value` and `greater than or equal to value`
func SplitByValue(root *Node, value E) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	_pushDown(root)
	if less(value, root.value) {
		first, second := SplitByValue(root.left, value)
		root.left = second
		return first, _pushUp(root)
	} else {
		first, second := SplitByValue(root.right, value)
		root.right = first
		return _pushUp(root), second
	}
}

// UpdateAll.
func AllApply(node *Node, f Id) *Node {
	node.value = mapping(f, node.value)
	node.data = mapping(f, node.data)
	node.lazy = composition(f, node.lazy)
	return node
}

func Toggle(node *Node) {
	node.left, node.right = node.right, node.left
	node.data = rev(node.data)
	node.isReversed = !node.isReversed
}

// Fold.
func Query(node *Node, start, end int) (res E) {
	if start >= end {
		return e()
	}
	left1, right1 := SplitByRank(node, start)
	left2, right2 := SplitByRank(right1, end-start)
	if left2 != nil {
		res = left2.data
	} else {
		res = e()
	}
	*node = *Merge(left1, Merge(left2, right2))
	return
}

func QueryAll(node *Node) E {
	if node == nil {
		return e()
	}
	return node.data
}

// Apply.
func Update(node *Node, start, end int, f Id) {
	if start >= end {
		return
	}
	left1, right1 := SplitByRank(node, start)
	left2, right2 := SplitByRank(right1, end-start)
	AllApply(left2, f)
	*node = *Merge(left1, Merge(left2, right2))
}

func Reverse(node *Node, start, end int) {
	if start >= end {
		return
	}
	left1, right1 := SplitByRank(node, start)
	left2, right2 := SplitByRank(right1, end-start)
	Toggle(left2)
	*node = *Merge(left1, Merge(left2, right2))
}

func Size(node *Node) int {
	if node == nil {
		return 0
	}
	return node.size
}

func _pushDown(node *Node) {
	if node == nil {
		return
	}
	if node.lazy != id() {
		if node.left != nil {
			AllApply(node.left, node.lazy)
		}
		if node.right != nil {
			AllApply(node.right, node.lazy)
		}
		node.lazy = id()
	}
	if node.isReversed {
		if node.left != nil {
			Toggle(node.left)
		}
		if node.right != nil {
			Toggle(node.right)
		}
		node.isReversed = false
	}
}

func _pushUp(node *Node) *Node {
	node.size = 1
	node.data = node.value
	if node.left != nil {
		node.size += node.left.size
		node.data = op(node.left.data, node.data)
	}
	if node.right != nil {
		node.size += node.right.size
		node.data = op(node.data, node.right.data)
	}
	return node
}

var _seed = uint64(time.Now().UnixNano()/2 + 1)

func _nextRand() uint64 {
	_seed ^= _seed << 7
	_seed ^= _seed >> 9
	return _seed
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
