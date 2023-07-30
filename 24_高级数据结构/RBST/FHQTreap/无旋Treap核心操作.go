// RBST:  https://nyaannyaan.github.io/library/rbst/lazy-reversible-rbst.hpp
// Treap: https://nyaannyaan.github.io/library/rbst/treap.hpp
//
// https://hitonanode.github.io/cplib-cpp/data_structure/lazy_rbst.hpp

// 分裂/拼接api:
//  1. Merge(left, right) -> root
//  2. Add(root, node) -> root
//  3. SplitByRank(root, k) -> [0,k) and [k,n)
//  4. SplitByValue(root, v) -> （-inf,v) and [v,inf)
//
// 查询/更新api:
//  1. AllApply(node, lazy) -> node
//  2. Toggle(node)
//  3. Query(node, start, end) -> res
//  4. Update(node, start, end, lazy)
//
// 操作api:
//  1. Reverse(node, start, end)
//  2. Size(node) -> size
//
// 构建api:
//  1. NewRoot() -> root
//  2. NewNode(v) -> node

package main

import (
	"fmt"
	"time"
)

func main() {
	assert := func(cur, expect interface{}) {
		if cur != expect {
			panic(fmt.Sprintf("cur: %v, expect: %v", cur, expect))
		}
	}
	assert(maxIncreasingGroups([]int{1, 2, 5}), 3) // 3
	assert(maxIncreasingGroups([]int{2, 2, 2}), 3)
	assert(maxIncreasingGroups([]int{1, 1}), 1)
	fmt.Println("OK")
}

// 用一个平衡树存储所有数字的频率，创建长度为 res 的数组时，
// 选取频率最大的 res 个数，将频率减 1 后放回平衡树中。
// !难点在于怎么放回平衡树
// 2790. 长度递增组的最大数目
// https://leetcode.cn/problems/maximum-number-of-groups-with-increasing-length/solution/bao-li-mo-ni-fa-by-vclip-wcxi/
func maxIncreasingGroups(usageLimits []int) int {
	n := len(usageLimits)
	tree := NewRoot()
	for i := 0; i < n; i++ {
		node := NewNode(usageLimits[i])
		tree = Add(tree, node)
	}

	for i := 1; i <= n; i++ {
		if Size(tree) < i {
			return i - 1
		}
		big, small := SplitByRank(tree, i)
		UpdateAll(big, -1) // 取出频率最大的 res 个数, 频率减 1
		max_ := QueryAll(small)
		notLess, less := SplitByValue(big, max_)
		nonZero, _ := SplitByValue(less, 1)
		tree = Merge(notLess, small) // 左右拼接后是有序的
		tree = Add(tree, nonZero)    // 顺序插入
	}

	return n
}

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

func Split3ByRank(root *Node, l, r int) (*Node, *Node, *Node) {
	if root == nil {
		return nil, nil, nil
	}
	root, right := SplitByRank(root, r)
	left, mid := SplitByRank(root, l)
	return left, mid, right
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

func Toggle(node *Node) {
	node.left, node.right = node.right, node.left
	node.data = rev(node.data)
	node.isReversed = !node.isReversed
}

// Fold.
func Query(node **Node, start, end int) (res E) {
	if start >= end || node == nil {
		return e()
	}
	left1, right1 := SplitByRank(*node, start)
	left2, right2 := SplitByRank(right1, end-start)
	if left2 != nil {
		res = left2.data
	} else {
		res = e()
	}
	*node = Merge(left1, Merge(left2, right2))
	return
}

func QueryAll(node *Node) E {
	if node == nil {
		return e()
	}
	return node.data
}

// Apply.
func Update(node **Node, start, end int, f Id) {
	if start >= end {
		return
	}
	left1, right1 := SplitByRank(*node, start)
	left2, right2 := SplitByRank(right1, end-start)
	_allApply(left2, f)
	*node = Merge(left1, Merge(left2, right2))
}

// AllApply.
func UpdateAll(node *Node, f Id) {
	node.value = mapping(f, node.value)
	node.data = mapping(f, node.data)
	node.lazy = composition(f, node.lazy)
}

func GetAll(node *Node) []E {
	res := make([]E, 0, Size(node))
	var dfs func(node *Node)
	dfs = func(node *Node) {
		if node == nil {
			return
		}
		_pushDown(node)
		dfs(node.left)
		res = append(res, node.value)
		dfs(node.right)
	}
	dfs(node)
	return res
}

func _allApply(node *Node, f Id) *Node {
	node.value = mapping(f, node.value)
	node.data = mapping(f, node.data)
	node.lazy = composition(f, node.lazy)
	return node
}

func Reverse(node **Node, start, end int) {
	if start >= end {
		return
	}
	left1, right1 := SplitByRank(*node, start)
	left2, right2 := SplitByRank(right1, end-start)
	Toggle(left2)
	*node = Merge(left1, Merge(left2, right2))
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
			_allApply(node.left, node.lazy)
		}
		if node.right != nil {
			_allApply(node.right, node.lazy)
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
	if left := node.left; left != nil {
		node.size += left.size
		node.data = op(left.data, node.data)
	}
	if right := node.right; right != nil {
		node.size += right.size
		node.data = op(node.data, right.data)
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

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------

// 按照leaves的顺序(不是值的顺序)构建一棵树.
func Build(leaves []E) *Node {
	if len(leaves) == 0 {
		return nil
	}
	var dfs func(l, r int) *Node
	dfs = func(l, r int) *Node {
		if r-l == 1 {
			node := NewNode(leaves[l])
			return _pushUp(node)
		}
		mid := (l + r) >> 1
		return Merge(dfs(l, mid), dfs(mid, r))
	}
	return dfs(0, len(leaves))
}

// 将value插入到树node的index位置之前.
func Insert(node **Node, index int, value E) {
	left, right := SplitByRank(*node, index)
	*node = Merge(left, Merge(NewNode(value), right))
}

func Pop(root **Node, index int) (res E) {
	n := Size(*root)
	if index < 0 {
		index += n
	}

	x, y, z := Split3ByRank(*root, index, index+1)
	res = y.value
	*root = Merge(x, z)
	return
}

// Remove [start, stop) from list.
func Erase(root **Node, start, stop int) {
	x, _, z := Split3ByRank(*root, start, stop)
	*root = Merge(x, z)
}

// Rotate [start, stop) to the right `k` times.
func RotateRight(root **Node, start, stop, k int) {
	start++
	n := stop - start + 1 - k%(stop-start+1)
	x, y := SplitByRank(*root, start-1)
	y, z := SplitByRank(y, n)
	z, p := SplitByRank(z, stop-start+1-n)
	*root = Merge(Merge(Merge(x, z), y), p)
}

// Rotate [start, stop) to the left `k` times.
func RotateLeft(root **Node, start, stop, k int) {
	start++
	k %= (stop - start + 1)
	x, y := SplitByRank(*root, start-1)
	y, z := SplitByRank(y, k)
	z, p := SplitByRank(z, stop-start+1-k)
	*root = Merge(Merge(Merge(x, z), y), p)
}

// rbst.Query(0, i) が true となるような最大の i を返す．
//
//	i := rbst.MaxRight(e, func(v E) bool { return v.sum <= k })
//	単調性を仮定．atcoder::lazy_segtree と同じ．
//	e は単位元．
func MaxRight(root *Node, e E, f func(E) bool) int {
	if root == nil {
		return 0
	}
	_pushDown(root)
	now := root
	prodNow := e
	res := 0
	for {
		if now.left != nil {
			_pushDown(now.left)
			pl := op(prodNow, now.left.data)
			if f(pl) {
				prodNow = pl
				res += now.left.size
			} else {
				now = now.left
				continue
			}
		}
		pl := op(prodNow, now.value)
		if !f(pl) {
			return res
		}
		prodNow = pl
		res++
		if now.right == nil {
			return res
		}
		_pushDown(now.right)
		now = now.right
	}
}

// rbst.Query(i, rbst.Size()) が true となるような最小の i を返す．
//
//	i := rbst.MinLeft(e, func(v E) bool { return v.sum >= k })
//	単調性を仮定．atcoder::lazy_segtree と同じ．
//	e は単位元．
func MinLeft(root *Node, e E, f func(E) bool) int {
	if root == nil {
		return 0
	}
	_pushDown(root)
	now := root
	prodNow := e
	res := Size(root)
	for {
		if now.right != nil {
			_pushDown(now.right)
			pr := op(now.right.data, prodNow)
			if f(pr) {
				prodNow = pr
				res -= now.right.size
			} else {
				now = now.right
				continue
			}
		}
		pr := op(now.value, prodNow)
		if !f(pr) {
			return res
		}
		prodNow = pr
		res--
		if now.left == nil {
			return res
		}
		_pushDown(now.left)
		now = now.left
	}

}
