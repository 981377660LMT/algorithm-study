// 单点修改, 区间查询
// !在稀疏时比较慢,应该使用sparse版本

package main

import (
	"bufio"
	"fmt"
	"os"
)

func demo() {
	seg := NewDynamicSegTree(0, 10, false)
	root := seg.Build([]E{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	root = seg.Set(root, 1, 1)
	fmt.Println(seg.Query(root, 0, 1), seg.GetAll(root))
}

func main() {
	// https://yukicoder.me/problems/no/789
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	seg := NewDynamicSegTree(0, 1e9+10, false)
	root := seg.NewRoot()
	var n int
	fmt.Fscan(in, &n)

	res := 0
	for i := 0; i < n; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			root = seg.Update(root, a, b)
		} else {
			var left, right int
			fmt.Fscan(in, &left, &right)
			res += seg.Query(root, left, right+1)
		}
	}
	fmt.Fprintln(out, res)
}

type E = int

func e1() E                { return 0 }
func e2(left, right int) E { return 0 } // 区间[left,right)的初始值.
func op(a, b E) E          { return a + b }

type DynamicSegTree struct {
	L, R       int
	persistent bool
	unit       E
}

type SegNode struct {
	l, r *SegNode
	x    E
}

func NewDynamicSegTree(left, right int, persistent bool) *DynamicSegTree {
	return &DynamicSegTree{
		L:          left,
		R:          right + 5,
		persistent: persistent,
		unit:       e1(),
	}
}

func (ds *DynamicSegTree) NewRoot() *SegNode {
	return &SegNode{x: e2(ds.L, ds.R)}
}

func (ds *DynamicSegTree) Build(nums []E) *SegNode {
	return ds._buildRec(0, len(nums), nums)
}

// L<=left<=right<=R
func (ds *DynamicSegTree) Query(root *SegNode, left, right int) E {
	if left == right {
		return ds.unit
	}
	x := ds.unit
	ds._queryRec(root, ds.L, ds.R, left, right, &x)
	return x
}

// L<=index<R
func (ds *DynamicSegTree) Set(root *SegNode, index int, value E) *SegNode {
	return ds._setRec(root, ds.L, ds.R, index, value)
}

// L<=left<R
func (ds *DynamicSegTree) Update(root *SegNode, index int, value E) *SegNode {
	return ds._updateRec(root, ds.L, ds.R, index, value)
}

// L<=right<=R
func (ds *DynamicSegTree) MinLeft(root *SegNode, right int, check func(E) bool) int {
	x := ds.unit
	return ds._minLeftRec(root, ds.L, ds.R, right, check, &x)
}

// L<=left<=R
func (ds *DynamicSegTree) MaxRight(root *SegNode, left int, check func(E) bool) int {
	x := ds.unit
	return ds._maxRightRec(root, ds.L, ds.R, left, check, &x)
}

func (ds *DynamicSegTree) GetAll(root *SegNode) []struct {
	index int
	value E
} {
	res := make([]struct {
		index int
		value E
	}, 0, ds.R-ds.L)
	ds._getAllRec(root, ds.L, ds.R, &res)
	return res
}

func (ds *DynamicSegTree) _newNode(left, right int) *SegNode {
	return &SegNode{x: e2(left, right)}
}

func (ds *DynamicSegTree) _newNodeWithValue(x E) *SegNode {
	return &SegNode{x: x}
}

func (ds *DynamicSegTree) _copyNode(node *SegNode) *SegNode {
	if node == nil || !ds.persistent {
		return node
	}
	return &SegNode{l: node.l, r: node.r, x: node.x}
}

func (ds *DynamicSegTree) _buildRec(left, right int, nums []E) *SegNode {
	if left == right {
		return nil
	}
	if right == left+1 {
		return ds._newNodeWithValue(nums[left])
	}
	mid := (left + right) >> 1
	lRoot := ds._buildRec(left, mid, nums)
	rRoot := ds._buildRec(mid, right, nums)
	x := op(lRoot.x, rRoot.x)
	root := ds._newNodeWithValue(x)
	root.l = lRoot
	root.r = rRoot
	return root
}

func (ds *DynamicSegTree) _setRec(root *SegNode, l, r, i int, x E) *SegNode {
	if l == r-1 {
		root = ds._copyNode(root)
		root.x = x
		return root
	}
	m := (l + r) >> 1
	root = ds._copyNode(root)
	if i < m {
		if root.l == nil {
			root.l = ds._newNode(l, m)
		}
		root.l = ds._setRec(root.l, l, m, i, x)
	} else {
		if root.r == nil {
			root.r = ds._newNode(m, r)
		}
		root.r = ds._setRec(root.r, m, r, i, x)
	}
	var xl, xr E
	if root.l == nil {
		xl = ds.unit
	} else {
		xl = root.l.x
	}
	if root.r == nil {
		xr = ds.unit
	} else {
		xr = root.r.x
	}
	root.x = op(xl, xr)
	return root
}

func (ds *DynamicSegTree) _updateRec(root *SegNode, l, r, i int, x E) *SegNode {
	if l == r-1 {
		root = ds._copyNode(root)
		root.x = op(root.x, x)
		return root
	}
	m := (l + r) >> 1
	root = ds._copyNode(root)
	if i < m {
		if root.l == nil {
			root.l = ds._newNode(l, m)
		}
		root.l = ds._updateRec(root.l, l, m, i, x)
	} else {
		if root.r == nil {
			root.r = ds._newNode(m, r)
		}
		root.r = ds._updateRec(root.r, m, r, i, x)
	}
	var xl, xr E
	if root.l == nil {
		xl = ds.unit
	} else {
		xl = root.l.x
	}
	if root.r == nil {
		xr = ds.unit
	} else {
		xr = root.r.x
	}
	root.x = op(xl, xr)
	return root
}

func (ds *DynamicSegTree) _queryRec(root *SegNode, l, r, ql, qr int, x *E) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr {
		return
	}
	if root == nil {
		*x = op(*x, e2(ql, qr))
		return
	}
	if l == ql && r == qr {
		*x = op(*x, root.x)
		return
	}
	m := (l + r) >> 1
	ds._queryRec(root.l, l, m, ql, qr, x)
	ds._queryRec(root.r, m, r, ql, qr, x)
}

func (ds *DynamicSegTree) _minLeftRec(root *SegNode, l, r, qr int, check func(E) bool, x *E) int {
	if qr <= l {
		return ds.L
	}
	if r <= qr && check(op(root.x, *x)) {
		*x = op(*x, root.x)
		return ds.L
	}
	if r == l+1 {
		return r
	}
	m := (l + r) >> 1
	if root.r == nil {
		root.r = ds._newNode(m, r)
	}
	k := ds._minLeftRec(root.r, m, r, qr, check, x)
	if k != ds.L {
		return k
	}
	if root.l == nil {
		root.l = ds._newNode(l, m)
	}
	return ds._minLeftRec(root.l, l, m, qr, check, x)

}

func (ds *DynamicSegTree) _maxRightRec(root *SegNode, l, r, ql int, check func(E) bool, x *E) int {
	if r <= ql {
		return ds.R
	}
	if ql <= l && check(op(*x, root.x)) {
		*x = op(*x, root.x)
		return ds.R
	}
	if r == l+1 {
		return l
	}
	m := (l + r) >> 1
	if root.l == nil {
		root.l = ds._newNode(l, m)
	}
	k := ds._maxRightRec(root.l, l, m, ql, check, x)
	if k != ds.R {
		return k
	}
	if root.r == nil {
		root.r = ds._newNode(m, r)
	}
	return ds._maxRightRec(root.r, m, r, ql, check, x)
}

func (ds *DynamicSegTree) _getAllRec(root *SegNode, left, right int, res *[]struct {
	index int
	value E
}) {
	if root == nil {
		return
	}
	if right-left == 1 {
		*res = append(*res, struct {
			index int
			value E
		}{left, root.x})
		return
	}
	mid := (left + right) >> 1
	ds._getAllRec(root.l, left, mid, res)
	ds._getAllRec(root.r, mid, right, res)
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
