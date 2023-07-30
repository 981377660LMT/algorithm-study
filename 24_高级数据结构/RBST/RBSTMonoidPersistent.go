// Reference:
//
//	RBST:  https://nyaannyaan.github.io/library/rbst/lazy-reversible-hpp
//	       https://hitonanode.github.io/cplib-cpp/data_structure/lazy_hpp
//	Treap: https://nyaannyaan.github.io/library/rbst/treap.hpp
//
// !幺半群上的可持久化RBST, 注意翻转
//
//		分裂/拼接api:
//		 1. Merge(left, right) -> root
//		 2. Add(root, node) -> root
//		 3. SplitByRank(root, k) -> [0,k) and [k,n)
//		 4. SplitByValue(root, v) -> （-inf,v) and [v,inf)
//
//		查询/更新api:
//		 1. Query(node, start, end) -> res
//		 2. QueryAll(node) -> res
//		 3. UpdateRange(node, start, end, lazy) -> node
//		 4. UpdateAll(node, lazy) -> node
//	   5. Set(node, k, v) -> node
//		 6. Get(node, k) -> v
//		 7. GetAll(node) -> []v
//
//		构建api:
//		 1. NewRoot() -> root
//		 2. NewNode(v) -> node
//	   3. Build(leaves) -> root
//
//		操作api:
//		 1. Reverse(node, start, end) -> node
//		 2. CopyWithin(node, start, end, to) -> node (持久化为true时)
//	   3. KthNode(node, k) -> node
//		 4. At(node, k) -> v
//		 5. Size(node) -> size
//		 6. Pop(node, k) -> node, v
//		 7. Erase(node, start, end) -> node
//		 8. Insert(node, k, v) -> node
//		 9. RotateRight(node, start, end, k) -> node
//		10. RotateLeft(node, start, end, k) -> node
//
//		!因为支持可持久化，所有修改操作都必须返回新的root，所有非修改操作都传入指向root的指针

package main

import (
	"bufio"
	"fmt"
	"os"
)

func demo() {
	// tree1 := Build([]E{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}})
	// fmt.Println(QueryAll(tree1))
	// tree2 := Update(tree1, 1, 3, 1)
	// fmt.Println(QueryAll(tree2))
	// fmt.Println(QueryAll(tree1))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		leaves[i] = E{sum: x, size: 1}
	}
	root := Build(leaves)

	for i := 0; i < q; i++ {
		var t, a, b, c, d, v int
		fmt.Fscan(in, &t)
		if t == 1 {
			fmt.Fscan(in, &a, &b, &v)
			root = Update(root, a-1, b, v)
		} else if t == 2 {
			fmt.Fscan(in, &a, &b, &c, &d)
			root = CopyWithin(root, a-1, c-1, d)
		} else if t == 3 {
			fmt.Fscan(in, &a, &b)
			fmt.Fprintln(out, Query(&root, a-1, b).sum)
		}
	}
}

const INF int = 1e18

const _PERSISTENT = true // !是否启用持久化

// RangeAddRangeMax
type E = struct{ sum, size int }
type Id = int

func rev(e E) E     { return e }
func e() E          { return E{} }
func id() Id        { return 0 }
func op(e1, e2 E) E { return E{e1.sum + e2.sum, e1.size + e2.size} }

func mapping(f Id, e E) E {
	if f == 0 {
		return e
	}
	return E{f*e.size + e.sum, e.size}
}

func composition(f, g Id) Id { return f + g }
func less(e1, e2 E) bool     { return e1.sum < e2.sum }

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
	isReversed  bool
}

func NewRoot() *Node {
	return nil
}

func NewNode(v E) *Node {
	res := &Node{value: v, data: v, size: 1, lazy: id()}
	return res
}

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
	leftSize, rightSize := uint32(Size(left)), uint32(Size(right))
	rand := _nextRand()
	if rand%(leftSize+rightSize) < leftSize {
		_pushDown(left)
		left = _copyNode(left)
		left.right = Merge(left.right, right)
		return _pushUp(left)
	} else {
		_pushDown(right)
		right = _copyNode(right)
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
		root = _copyNode(root)
		root.left = second
		_pushUp(root)
		return first, root
	} else {
		first, second := SplitByRank(root.right, k-leftSize-1)
		root = _copyNode(root)
		root.right = first
		_pushUp(root)
		return root, second
	}
}

// 左中右子树:[0, l) and [l, r) and [r, n).
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
		root = _copyNode(root)
		root.left = second
		return first, _pushUp(root)
	} else {
		first, second := SplitByValue(root.right, value)
		root = _copyNode(root)
		root.right = first
		return _pushUp(root), second
	}
}

// Fold.
func Query(node **Node, start, end int) (res E) {
	if start >= end {
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
func Update(node *Node, start, end int, f Id) *Node {
	if start >= end {
		return node
	}
	left1, right1 := SplitByRank(node, start)
	left2, right2 := SplitByRank(right1, end-start)
	_allApply(left2, f)
	return Merge(left1, Merge(left2, right2))
}

func UpdateAll(node *Node, f Id) *Node {
	if node == nil {
		return nil
	}
	return Update(node, 0, node.size, f)
}

func Reverse(node *Node, start, end int) *Node {
	if start >= end {
		return node
	}
	left1, right1 := SplitByRank(node, start)
	left2, right2 := SplitByRank(right1, end-start)
	_toggle(left2)
	return Merge(left1, Merge(left2, right2))
}

func ReverseAll(node *Node) *Node {
	node = _copyNode(node)
	_toggle(node)
	return node
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
	if node.isReversed || node.lazy != id() {
		node.left, node.right = _copyNode(node.left), _copyNode(node.right)
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
			_toggle(node.left)
		}
		if node.right != nil {
			_toggle(node.right)
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

func _allApply(node *Node, f Id) *Node {
	node.value = mapping(f, node.value)
	node.data = mapping(f, node.data)
	node.lazy = composition(f, node.lazy)
	return node
}

func _toggle(node *Node) {
	node.left, node.right = node.right, node.left
	node.data = rev(node.data)
	node.isReversed = !node.isReversed
}

var _x uint32 = 123456789
var _y uint32 = 362436069
var _z uint32 = 521288629
var _w uint32 = 88675123

func _nextRand() uint32 {
	t := _x ^ (_x << 11)
	_x, _y, _z = _y, _z, _w
	_w = (_w ^ (_w >> 19)) ^ (t ^ (t >> 8))
	return _w
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

func _copyNode(node *Node) *Node {
	if node == nil || !_PERSISTENT {
		return node
	}
	return &Node{
		left:       node.left,
		right:      node.right,
		value:      node.value,
		data:       node.data,
		lazy:       node.lazy,
		size:       node.size,
		isReversed: node.isReversed,
	}
}

func CopyWithin(root *Node, target int, start, end int) *Node {
	if !_PERSISTENT {
		panic("CopyWithin should be used in persistent mode")
	}
	len := end - start
	p1Left, p1Right := SplitByRank(root, start)
	p2Left, p2Right := SplitByRank(p1Right, len)
	root = Merge(p1Left, Merge(p2Left, p2Right))
	p3Left, p3Right := SplitByRank(root, target)
	_, p4Right := SplitByRank(p3Right, len)
	root = Merge(p3Left, Merge(p2Left, p4Right))
	return root
}

func KthNode(root *Node, k int) *Node {
	cur := root
	for cur != nil {
		leftSize := Size(cur.left)
		if leftSize+1 == k {
			break
		} else if leftSize >= k {
			cur = cur.left
		} else {
			k -= leftSize + 1
			cur = cur.right
		}
	}
	return cur
}

func At(root *Node, index int) E {
	n := Size(root)
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		return e()
	}
	node := KthNode(root, index)
	if node == nil {
		return e()
	}
	return node.value
}

func Pop(root *Node, index int) (newRoot *Node, res E) {
	n := Size(root)
	if index < 0 {
		index += n
	}

	x, y, z := Split3ByRank(root, index, index+1)
	res = y.value
	newRoot = Merge(x, z)
	return
}

// Remove [start, stop) from list.
func Erase(root *Node, start, stop int) *Node {
	x, _, z := Split3ByRank(root, start, stop)
	return Merge(x, z)
}

// Insert node before pos.
func Insert(root *Node, pos int, node *Node) *Node {
	n := Size(root)
	if pos < 0 {
		pos += n
	}
	if pos < 0 {
		pos = 0
	}
	if pos > n {
		pos = n
	}
	left, right := SplitByRank(root, pos)
	return Merge(left, Merge(node, right))
}

// Rotate [start, stop) to the right `k` times.
func RotateRight(root *Node, start, stop, k int) *Node {
	start++
	n := stop - start + 1 - k%(stop-start+1)

	x, y := SplitByRank(root, start-1)
	y, z := SplitByRank(y, n)
	z, p := SplitByRank(z, stop-start+1-n)
	return Merge(Merge(Merge(x, z), y), p)
}

// Rotate [start, stop) to the left `k` times.
func RotateLeft(root *Node, start, stop, k int) *Node {
	start++
	k %= (stop - start + 1)

	x, y := SplitByRank(root, start-1)
	y, z := SplitByRank(y, k)
	z, p := SplitByRank(z, stop-start+1-k)
	return Merge(Merge(Merge(x, z), y), p)
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
