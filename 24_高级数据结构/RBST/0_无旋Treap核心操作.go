// 1. Merge(left, right) -> root
// 2. SplitByRank(root, k) -> [0,k) and [k,n)
// 3. SplitByValue(root, v) -> [0,v) and [v,n)
// 4. AllApply(node, lazy) -> node
// 5. Toggle(node)

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(maxIncreasingGroups([]int{1, 2, 5}))
	fmt.Println(maxIncreasingGroups([]int{2, 2, 2}))
}

// 用一个平衡树存储所有数字的频率，创建长度为 res 的数组时，
// 选取频率最大的 res 个数，将频率减 1 后放回平衡树中。
func maxIncreasingGroups(usageLimits []int) int {
	n := len(usageLimits)
	root := NewRoot()
	for i := 0; i < n; i++ {
		node := NewNode(usageLimits[i])
		root = Add(root, node)
	}

	for i := 1; i <= n; i++ {
		left, right := SplitByRank(root, i)
		if Size(left) < i {
			return i - 1
		}
		fmt.Println(left.data, left.size, "pre")
		if right != nil {
			fmt.Println(right.data, right.size, "right")
		}
		AllApply(left, -1) // 取出频率最大的 res 个数, 频率减 1
		fmt.Println(left.data, left.size, "after")
		nonZero, zero := SplitByValue(root, 1)
		if zero != nil {
			fmt.Println(zero.data, zero.size, "zero")
		}
		if nonZero != nil {
			fmt.Println(nonZero.data, nonZero.size, "nonZero")
		}
		root = Merge(nonZero, right)
	}
	return n
}

const INF int = 1e18

// RangeAddRangeMin
type E = int
type Id = int

func rev(e E) E              { return e }
func id() Id                 { return 0 }
func op(e1, e2 E) E          { return min(e1, e2) }
func mapping(f Id, e E) E    { return f + e }
func composition(f, g Id) Id { return f + g }
func less(e1, e2 E) bool     { return e1 > e2 } // 维护最大值

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

func NewNode(v E) *Node {
	res := &Node{value: v, data: v, size: 1, lazy: id(), priority: _nextRand()}
	return res
}

func NewRoot() *Node {
	return nil
}

func Add(root, node *Node) *Node {
	if node == nil {
		return root
	}
	left, right := SplitByValue(root, node.value)
	return Merge(Merge(left, node), right)
}

// 拼接两段区间
//
//	Merge left and right, return new root
func Merge(left, right *Node) *Node {
	if left == nil || right == nil {
		if left == nil {
			return right
		}
		return left
	}

	if left.priority < right.priority {
		_pushDown(left)
		left.right = Merge(left.right, right)
		return _pushUp(left)
	} else {
		_pushDown(right)
		right.left = Merge(left, right.left)
		return _pushUp(right)
	}
}

func Toggle(node *Node) {
	tmp := node.left
	node.left = node.right
	node.right = tmp
	node.data = rev(node.data)
	node.isReversed = !node.isReversed
}

func AllApply(node *Node, f Id) *Node {
	node.value = mapping(f, node.value)
	node.data = mapping(f, node.data)
	node.lazy = composition(f, node.lazy)
	return node
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
	if less(root.value, value) {
		first, second := SplitByValue(root.left, value)
		root.left = second
		return first, _pushUp(root)
	} else {
		first, second := SplitByValue(root.right, value)
		root.right = first
		return _pushUp(root), second
	}
}

func Size(node *Node) int {
	if node == nil {
		return 0
	}
	return node.size
}

func _pushDown(node *Node) {
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
