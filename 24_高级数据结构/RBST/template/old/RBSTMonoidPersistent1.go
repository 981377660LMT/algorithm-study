// Reference:
//
//	https://maspypy.github.io/library/ds/randomized_bst/rbst_acted_monoid.hpp
//
// !阿贝尔群上的可持久化RBST,元素必须满足交换律(commutative).
//	!例如不能将 {left:xxx, right:xx} 作为Abel群.
//
//	分裂/拼接api:
//	 1. Merge(left, right) -> root
//	 2. Add(root, node) -> root
//	 3. SplitByRank(root, k) -> [0,k) and [k,n)
//	 4. SplitByValue(root, v) -> （-inf,v) and [v,inf)
//
//	查询/更新api:
//	 1. Query(node, start, end) -> res
//	 2. QueryAll(node) -> res
//	 3. UpdateRange(node, start, end, lazy) -> node
//	 4. UpdateAll(node, lazy) -> node
//   5. Set(node, k, v) -> node
//	 6. Get(node, k) -> v
//	 7. GetAll(node) -> []v
//   8. SplitMaxRight(node,check) -> left,right
//
//	构建api:
//	 1. NewRoot() -> root
//	 2. NewNode(v) -> node
//   3. Build(leaves) -> root
//
//	操作api:
//	 1. Reverse(node, start, end) -> node
//	 2. CopyWithin(node, start, end, to) -> node (持久化为true时)
//	 3. Size(node) -> size
//	 4. Pop(node, k) -> node, v
//	 5. Erase(node, start, end) -> node
//	 6. Insert(node, k, v) -> node
//	 7. RotateRight(node, start, end, k) -> node
//	 8. RotateLeft(node, start, end, k) -> node

//	 Pop/Erase/At/BisectLeft/BisectRight... 都是基于分裂/拼接实现的
//
//	!因为支持可持久化，所有修改操作都必须返回新的root.

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	main2()
}

// https://atcoder.jp/contests/arc030/tasks/arc030_4
func main2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		leaves[i] = x
	}
	R := NewRBSTAbelGroup(true)
	root := R.Build(leaves)

	for i := 0; i < q; i++ {
		var t, a, b, c, d, v int32
		fmt.Fscan(in, &t)
		if t == 1 {
			fmt.Fscan(in, &a, &b, &v)
			root = R.UpdateRange(root, a-1, b, int(v))
		} else if t == 2 {
			fmt.Fscan(in, &a, &b, &c, &d)
			root = R.CopyWithin(root, a-1, c-1, d)
		} else if t == 3 {
			fmt.Fscan(in, &a, &b)
			fmt.Fprintln(out, R.Query(root, a-1, b))
		}
	}
}

func demo() {
	R := NewRBSTAbelGroup(false)
	root := R.Build([]E{1, 2, 3, 4, 5})
	fmt.Println(R.GetAll(root)) // [1 2 3 4 5]
	left, right := R.SplitMaxRight(root, func(e E) bool { return e <= 10 })
	fmt.Println(R.GetAll(left))  // [1 2 3 4]
	fmt.Println(R.GetAll(right)) // [5]
	left, a := R.Pop(left, 0)
	fmt.Println(a)
	fmt.Println(R.GetAll(left)) // [2 3 4]
	left = R.Set(left, 0, 3)
	fmt.Println(R.GetAll(left)) // [3 3 4]
	left = R.Erase(left, 0, 2)
	fmt.Println(R.GetAll(left)) // [4]
	left = R.Insert(left, 0, R.NewNode(1))
	fmt.Println(R.GetAll(left)) // [1 4]
	fmt.Println(R.GetAll(left)) // [1 4]
	left = R.Reverse(left, 0, 2)
	fmt.Println(R.GetAll(left)) // [4 1]
	left = R.UpdateRange(left, 0, 2, 1)
	fmt.Println(R.GetAll(left)) // [5 2]
	left = R.Insert(left, 1, R.NewNode(3))
	left = R.Insert(left, 1, R.NewNode(3))
	fmt.Println(R.GetAll(left)) // [5 3 3 2]
	left = R.RotateRight(left, 0, 2, 1)
	fmt.Println(R.GetAll(left)) // [3 5 3 2]
}

// RangeAddRangeSum
type E = int
type Id = int

func e() E                            { return 0 }
func id() Id                          { return 0 }
func op(e1, e2 E) E                   { return e1 + e2 }
func mapping(f Id, e E, size int32) E { return f*int(size) + e }
func composition(f, g Id) Id          { return f + g }
func less(e1, e2 E) bool              { return e1 < e2 }

//
//
//

// 每个结点代表一段区间
type Node struct {
	left, right *Node
	value, data E
	lazy        Id
	size        int32
	isReversed  bool
}

type RBSTAbelGroup struct {
	persistent bool
	x, y, z, w uint32
}

func NewRBSTAbelGroup(persistent bool) *RBSTAbelGroup {
	return &RBSTAbelGroup{
		persistent: persistent,
		x:          123456789,
		y:          362436069,
		z:          521288629,
		w:          88675123,
	}
}

func (rbst *RBSTAbelGroup) NewRoot() *Node {
	return nil
}

func (rbst *RBSTAbelGroup) NewNode(v E) *Node {
	return &Node{value: v, data: v, lazy: id(), size: 1}
}

// 按照leaves的顺序构建一棵树.(不保证Value有序)
func (rbst *RBSTAbelGroup) Build(leaves []E) *Node {
	var dfs func(l, r int) *Node
	dfs = func(l, r int) *Node {
		if l == r {
			return nil
		}
		if l+1 == r {
			return rbst.NewNode(leaves[l])
		}
		mid := (l + r) >> 1
		left, right := dfs(l, mid), dfs(mid+1, r)
		root := rbst.NewNode(leaves[mid])
		root.left, root.right = left, right
		rbst._pushUp(root)
		return root
	}
	return dfs(0, len(leaves))
}

// 合并两棵树, 保证Value有序.
func (rbst *RBSTAbelGroup) Add(root, node *Node) *Node {
	if node == nil {
		return root
	}
	left, right := rbst.SplitByValue(root, node.value)
	return rbst.Merge(rbst.Merge(left, node), right)
}

// 合并`左右`两棵树，保证Rank有序.
func (rbst *RBSTAbelGroup) Merge(left, right *Node) *Node {
	return rbst._mergeRec(left, right)
}

func (rbst *RBSTAbelGroup) Merge3(a, b, c *Node) *Node {
	return rbst.Merge(rbst.Merge(a, b), c)
}

func (rbst *RBSTAbelGroup) Merge4(a, b, c, d *Node) *Node {
	return rbst.Merge(rbst.Merge(rbst.Merge(a, b), c), d)
}

// 左右子树:[0, k) and [k, n).
func (rbst *RBSTAbelGroup) SplitByRank(root *Node, k int32) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	return rbst._splitRec(root, k)
}

// 左中右子树:[0, l) and [l, r) and [r, n).
func (rbst *RBSTAbelGroup) Split3ByRank(root *Node, l, r int32) (*Node, *Node, *Node) {
	if root == nil {
		return nil, nil, nil
	}
	root, right := rbst.SplitByRank(root, r)
	left, mid := rbst.SplitByRank(root, l)
	return left, mid, right
}

// 四个子树:[0, i) and [i, j) and [j, k) and [k, n).
func (rbst *RBSTAbelGroup) Split4ByRank(root *Node, i, j, k int32) (*Node, *Node, *Node, *Node) {
	if root == nil {
		return nil, nil, nil, nil
	}
	root, d := rbst.SplitByRank(root, k)
	a, b, c := rbst.Split3ByRank(root, i, j)
	return a, b, c, d
}

// 小子树和大子树:[-inf,value) and [value,inf)
func (rbst *RBSTAbelGroup) SplitByValue(root *Node, value E) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	rbst._pushDown(root)
	if less(value, root.value) {
		first, second := rbst.SplitByValue(root.left, value)
		root = rbst._copyNode(root)
		root.left = second
		return first, rbst._pushUp(root)
	} else {
		first, second := rbst.SplitByValue(root.right, value)
		root = rbst._copyNode(root)
		root.right = first
		return rbst._pushUp(root), second
	}
}

func (rbst *RBSTAbelGroup) Query(root *Node, start, end int32) E {
	if start >= end || root == nil {
		return e()
	}
	return rbst._queryRec(root, start, end, false)
}

func (rbst *RBSTAbelGroup) QueryAll(root *Node) E {
	if root == nil {
		return e()
	}
	return root.data
}

func (rbst *RBSTAbelGroup) Reverse(root *Node, start, end int32) *Node {
	if end-start <= 1 || root == nil {
		return root
	}
	left, mid, right := rbst.Split3ByRank(root, start, end)
	mid.isReversed = !mid.isReversed
	mid.left, mid.right = mid.right, mid.left
	return rbst.Merge3(left, mid, right)
}

func (rbst *RBSTAbelGroup) UpdateRange(root *Node, start, end int32, f Id) *Node {
	return rbst._updateRangeRec(root, start, end, f)
}

func (rbst *RBSTAbelGroup) CopyWithin(root *Node, target int32, start, end int32) *Node {
	if !rbst.persistent {
		panic("CopyWithin only works on persistent RBST")
	}
	len := end - start
	p1Left, p1Right := rbst.SplitByRank(root, start)
	p2Left, p2Right := rbst.SplitByRank(p1Right, len)
	root = rbst.Merge(p1Left, rbst.Merge(p2Left, p2Right))
	p3Left, p3Right := rbst.SplitByRank(root, target)
	_, p4Right := rbst.SplitByRank(p3Right, len)
	root = rbst.Merge(p3Left, rbst.Merge(p2Left, p4Right))
	return root
}

func (rbst *RBSTAbelGroup) Set(root *Node, k int32, v E) *Node {
	return rbst._setRec(root, k, v)
}

func (rbst *RBSTAbelGroup) Get(root *Node, k int32) E {
	return rbst._getRec(root, k, false, id())
}

func (rbst *RBSTAbelGroup) GetAll(root *Node) []E {
	res := make([]E, 0, rbst.Size(root))
	var dfs func(root *Node, rev bool, lazy Id)
	dfs = func(root *Node, rev bool, lazy Id) {
		if root == nil {
			return
		}
		me := mapping(lazy, root.value, 1)
		lazy = composition(lazy, root.lazy)
		left, right := root.left, root.right
		if rev {
			left, right = right, left
		}
		nextRev := rev != root.isReversed
		dfs(left, nextRev, lazy)
		res = append(res, me)
		dfs(right, nextRev, lazy)
	}
	dfs(root, false, id())
	return res
}

func (rbst *RBSTAbelGroup) SplitMaxRight(root *Node, check func(x E) bool) (*Node, *Node) {
	x := e()
	return rbst._splitMaxRightRec(root, &x, check)
}

func (rbst *RBSTAbelGroup) _splitMaxRightRec(root *Node, x *E, check func(v E) bool) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	rbst._pushDown(root)
	root = rbst._copyNode(root)
	y := op(*x, root.data)
	if check(y) {
		*x = y
		return root, nil
	}
	left, right := root.left, root.right
	if left != nil {
		y := op(*x, left.data)
		if !check(y) {
			n1, n2 := rbst._splitMaxRightRec(left, x, check)
			root.left = n2
			rbst._pushUp(root)
			return n1, root
		}
		*x = y
	}
	y = op(*x, root.value)
	if !check(y) {
		root.left = nil
		rbst._pushUp(root)
		return left, root
	}
	*x = y
	n1, n2 := rbst._splitMaxRightRec(right, x, check)
	root.right = n1
	rbst._pushUp(root)
	return root, n2
}

func (rbst *RBSTAbelGroup) Size(root *Node) int32 {
	if root == nil {
		return 0
	}
	return root.size
}

func (rbst *RBSTAbelGroup) _pushUp(node *Node) *Node {
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

func (rbst *RBSTAbelGroup) _pushDown(node *Node) {
	if node.lazy != id() || node.isReversed {
		node.left, node.right = rbst._copyNode(node.left), rbst._copyNode(node.right)
	}
	if node.isReversed {
		if left := node.left; left != nil {
			left.isReversed = !left.isReversed
			left.left, left.right = left.right, left.left
		}
		if right := node.right; right != nil {
			right.isReversed = !right.isReversed
			right.left, right.right = right.right, right.left
		}
		node.isReversed = false
	}
	if node.lazy != id() {
		if left := node.left; left != nil {
			left.value = mapping(node.lazy, left.value, 1)
			left.data = mapping(node.lazy, left.data, left.size)
			left.lazy = composition(node.lazy, left.lazy)
		}
		if right := node.right; right != nil {
			right.value = mapping(node.lazy, right.value, 1)
			right.data = mapping(node.lazy, right.data, right.size)
			right.lazy = composition(node.lazy, right.lazy)
		}
		node.lazy = id()
	}
}

func (rbst *RBSTAbelGroup) _copyNode(node *Node) *Node {
	if node == nil || !rbst.persistent {
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

func (rbst *RBSTAbelGroup) _mergeRec(left, right *Node) *Node {
	if left == nil || right == nil {
		if left == nil {
			return right
		}
		return left
	}
	leftSize, rightSize := uint32(left.size), uint32(right.size)
	rand := rbst._nextRand()
	if rand%(leftSize+rightSize) < leftSize {
		rbst._pushDown(left)
		left = rbst._copyNode(left)
		left.right = rbst._mergeRec(left.right, right)
		rbst._pushUp(left)
		return left
	} else {
		rbst._pushDown(right)
		right = rbst._copyNode(right)
		right.left = rbst._mergeRec(left, right.left)
		rbst._pushUp(right)
		return right
	}
}

func (rbst *RBSTAbelGroup) _splitRec(root *Node, k int32) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	rbst._pushDown(root)
	leftSize := rbst.Size(root.left)
	if k <= leftSize {
		left, right := rbst._splitRec(root.left, k)
		root = rbst._copyNode(root)
		root.left = right
		rbst._pushUp(root)
		return left, root
	} else {
		left, right := rbst._splitRec(root.right, k-leftSize-1)
		root = rbst._copyNode(root)
		root.right = left
		rbst._pushUp(root)
		return root, right
	}
}

func (rbst *RBSTAbelGroup) _setRec(root *Node, k int32, v E) *Node {
	if root == nil {
		return nil
	}
	rbst._pushDown(root)
	leftSize := rbst.Size(root.left)
	if k < leftSize {
		root = rbst._copyNode(root)
		root.left = rbst._setRec(root.left, k, v)
		rbst._pushUp(root)
		return root
	} else if k == leftSize {
		root = rbst._copyNode(root)
		root.value = v
		rbst._pushUp(root)
		return root
	} else {
		root = rbst._copyNode(root)
		root.right = rbst._setRec(root.right, k-leftSize-1, v)
		rbst._pushUp(root)
		return root
	}
}

func (rbst *RBSTAbelGroup) _updateRangeRec(root *Node, l, r int32, lazy Id) *Node {
	rbst._pushDown(root)
	root = rbst._copyNode(root)
	if l == 0 && r == root.size {
		root.value = mapping(lazy, root.value, 1)
		root.data = mapping(lazy, root.data, root.size)
		root.lazy = lazy
		return root
	}
	leftSize := rbst.Size(root.left)
	if l < leftSize {
		root.left = rbst._updateRangeRec(root.left, l, min32(r, leftSize), lazy)
	}
	if l <= leftSize && leftSize < r {
		root.value = mapping(lazy, root.value, 1)
	}
	k := 1 + leftSize
	if k < r {
		root.right = rbst._updateRangeRec(root.right, max32(k, l)-k, r-k, lazy)
	}
	rbst._pushUp(root)
	return root
}

func (rbst *RBSTAbelGroup) _getRec(root *Node, k int32, rev bool, lazy Id) E {
	left, right := root.left, root.right
	if rev {
		left, right = right, left
	}
	leftSize := rbst.Size(left)
	if k == leftSize {
		return mapping(lazy, root.value, 1)
	}
	lazy = composition(lazy, root.lazy)
	nextRev := rev != root.isReversed
	if k < leftSize {
		return rbst._getRec(left, k, nextRev, lazy)
	} else {
		return rbst._getRec(right, k-leftSize-1, nextRev, lazy)
	}
}

func (rbst *RBSTAbelGroup) _queryRec(root *Node, l, r int32, rev bool) E {
	if l == 0 && r == root.size {
		return root.data
	}
	left, right := root.left, root.right
	if rev {
		left, right = right, left
	}
	leftSize := rbst.Size(left)
	nextRev := rev != root.isReversed
	res := e()
	if l < leftSize {
		y := rbst._queryRec(left, l, min32(r, leftSize), nextRev)
		res = op(res, mapping(root.lazy, y, min32(r, leftSize)-l))
	}
	if l <= leftSize && leftSize < r {
		res = op(res, root.value)
	}
	k := 1 + leftSize
	if k < r {
		y := rbst._queryRec(right, max32(k, l)-k, r-k, nextRev)
		res = op(res, mapping(root.lazy, y, r-max32(k, l)))
	}
	return res
}

func (rbst *RBSTAbelGroup) _nextRand() uint32 {
	t := rbst.x ^ (rbst.x << 11)
	rbst.x, rbst.y, rbst.z = rbst.y, rbst.z, rbst.w
	rbst.w = (rbst.w ^ (rbst.w >> 19)) ^ (t ^ (t >> 8))
	return rbst.w
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

func (rbst *RBSTAbelGroup) Pop(root *Node, index int32) (newRoot *Node, res E) {
	n := rbst.Size(root)
	if index < 0 {
		index += n
	}

	x, y, z := rbst.Split3ByRank(root, index, index+1)
	res = y.value
	newRoot = rbst.Merge(x, z)
	return
}

// Remove [start, stop) from list.
func (rbst *RBSTAbelGroup) Erase(root *Node, start, stop int32) *Node {
	x, _, z := rbst.Split3ByRank(root, start, stop)
	return rbst.Merge(x, z)
}

// Insert node before pos.
func (rbst *RBSTAbelGroup) Insert(root *Node, pos int32, node *Node) *Node {
	n := rbst.Size(root)
	if pos < 0 {
		pos += n
	}
	if pos < 0 {
		pos = 0
	}
	if pos > n {
		pos = n
	}
	left, right := rbst.SplitByRank(root, pos)
	return rbst.Merge(left, rbst.Merge(node, right))
}

// Rotate [start, stop) to the right `k` times.
func (rbst *RBSTAbelGroup) RotateRight(root *Node, start, stop, k int32) *Node {
	start++
	n := stop - start + 1 - k%(stop-start+1)

	x, y := rbst.SplitByRank(root, start-1)
	y, z := rbst.SplitByRank(y, n)
	z, p := rbst.SplitByRank(z, stop-start+1-n)
	return rbst.Merge(rbst.Merge(rbst.Merge(x, z), y), p)
}

// Rotate [start, stop) to the left `k` times.
func (rbst *RBSTAbelGroup) RotateLeft(root *Node, start, stop, k int32) *Node {
	start++
	k %= (stop - start + 1)

	x, y := rbst.SplitByRank(root, start-1)
	y, z := rbst.SplitByRank(y, k)
	z, p := rbst.SplitByRank(z, stop-start+1-k)
	return rbst.Merge(rbst.Merge(rbst.Merge(x, z), y), p)
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
