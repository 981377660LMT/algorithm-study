// https://yukicoder.me/problems/no/2292
// No.2292 Interval Union Find
// 给定一个n个点，0条边的无向图，有Q个操作，操作有以下四种：
// 1 L R ：连接[L,R)区间内的所有点对之间的边；
// 2 L R ：删除[L,R)区间内的所有点对之间的边；
// 3 u v ：判断u和v是否连通；
// 4 x ：求x的连通块大小；

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	// !x与x+1相连：data[x] = 0
	seg := NewDynamicSegTreeLazy(0, n+10, false)
	root := seg.NewRoot()
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if l == r {
				continue
			}
			root = seg.UpdateRange(root, l, r, 0) // 所有点连通
		} else if op == 2 {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if l == r {
				continue
			}
			root = seg.UpdateRange(root, l, r, 1) // 所有点不连通
		} else if op == 3 {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if l > r {
				l, r = r, l
			}
			var x int
			if l == r {
				x = 0
			} else {
				x = seg.Query(root, l, r)
			}
			if x == 0 {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		} else if op == 4 {
			var v int
			fmt.Fscan(in, &v)
			r := seg.MaxRight(root, v, func(e E) bool { return e == 0 })
			l := seg.MinLeft(root, v, func(e E) bool { return e == 0 })
			fmt.Fprintln(out, r-l+1)
		}
	}
}

// RangeAssignRangeSum
type E = int
type Id = int

func e1() E               { return 0 }
func e2(start, end int) E { return end - start } // 区间[start,end)的初始值.
func id() Id              { return -1 }
func op(a, b E) E         { return a + b }
func mapping(f Id, g E, size int) E {
	if f == -1 {
		return g
	}
	return f * E(size)
}
func composition(f, g Id) Id {
	if f == -1 {
		return g
	}
	return f
}

type DynamicSegTreeLazy struct {
	L, R       int
	persistent bool
	dataUnit   E
	lazyUnit   Id
}

type SegNode struct {
	x    E
	lazy Id
	l, r *SegNode
}

func NewDynamicSegTreeLazy(start, end int, persistent bool) *DynamicSegTreeLazy {
	return &DynamicSegTreeLazy{
		L:          start,
		R:          end + 5,
		persistent: persistent,
		dataUnit:   e1(),
		lazyUnit:   id(),
	}
}

func (ds *DynamicSegTreeLazy) NewRoot() *SegNode {
	return &SegNode{x: e2(ds.L, ds.R), lazy: ds.lazyUnit}
}

func (ds *DynamicSegTreeLazy) Build(nums []E) *SegNode {
	return ds._buildRec(0, len(nums), nums)
}

// L<=start<=end<=R
func (ds *DynamicSegTreeLazy) Query(root *SegNode, start, end int) E {
	if start == end {
		return ds.dataUnit
	}
	x := ds.dataUnit
	ds._queryRec(root, ds.L, ds.R, start, end, &x, ds.lazyUnit)
	return x
}

func (ds *DynamicSegTreeLazy) QueryAll(root *SegNode) E {
	return ds.Query(root, ds.L, ds.R)
}

// L<=index<R
func (ds *DynamicSegTreeLazy) Set(root *SegNode, index int, value E) *SegNode {
	return ds._setRec(root, ds.L, ds.R, index, value)
}

// L<=index<R
func (ds *DynamicSegTreeLazy) Update(root *SegNode, index int, value E) *SegNode {
	return ds._updateRec(root, ds.L, ds.R, index, value)
}

// L<=start<=end<=R
func (ds *DynamicSegTreeLazy) UpdateRange(root *SegNode, start, end int, lazy Id) *SegNode {
	if start == end {
		return root
	}
	return ds._updateRangeRec(root, ds.L, ds.R, start, end, lazy)
}

// 二分查询最大的 right 使得切片 [left:right) 内的值满足 check.
// L<=right<=R
func (ds *DynamicSegTreeLazy) MinLeft(root *SegNode, end int, check func(E) bool) int {
	x := ds.dataUnit
	return ds._minLeftRec(root, ds.L, ds.R, end, check, &x)
}

// 二分查询最小的 left 使得切片 [left:right) 内的值满足 check.
// L<=left<=R
func (ds *DynamicSegTreeLazy) MaxRight(root *SegNode, start int, check func(E) bool) int {
	x := ds.dataUnit
	return ds._maxRightRec(root, ds.L, ds.R, start, check, &x)
}

func (ds *DynamicSegTreeLazy) GetAll(root *SegNode) []E {
	res := make([]E, 0)
	ds._getAllRec(root, ds.L, ds.R, &res, ds.lazyUnit)
	return res
}

func (ds *DynamicSegTreeLazy) EnumerateAll(root *SegNode, f func(index int, value E)) {
	var dfs func(c *SegNode, l, r int, lazy Id)
	dfs = func(c *SegNode, l, r int, lazy Id) {
		if c == nil {
			return
		}
		if r-l == 1 {
			f(l, mapping(lazy, c.x, 1))
			return
		}
		m := (l + r) >> 1
		lazy = composition(lazy, c.lazy)
		dfs(c.l, l, m, lazy)
		dfs(c.r, m, r, lazy)
	}
	dfs(root, ds.L, ds.R, ds.lazyUnit)
}

func (ds *DynamicSegTreeLazy) Copy(node *SegNode) *SegNode {
	if node == nil || !ds.persistent {
		return node
	}
	return &SegNode{l: node.l, r: node.r, x: node.x, lazy: node.lazy}
}

func (ds *DynamicSegTreeLazy) _newNode(left, right int) *SegNode {
	return &SegNode{x: e2(left, right), lazy: ds.lazyUnit}
}

func (ds *DynamicSegTreeLazy) _newNodeWithValue(x E) *SegNode {
	return &SegNode{x: x, lazy: ds.lazyUnit}
}

func (ds *DynamicSegTreeLazy) _pushDown(node *SegNode, l, r int) {
	if node.lazy == ds.lazyUnit {
		return
	}
	m := (l + r) >> 1
	if node.l == nil {
		node.l = ds._newNode(l, m)
	} else {
		node.l = ds.Copy(node.l)
	}
	node.l.x = mapping(node.lazy, node.l.x, m-l)
	node.l.lazy = composition(node.lazy, node.l.lazy)
	if node.r == nil {
		node.r = ds._newNode(m, r)
	} else {
		node.r = ds.Copy(node.r)
	}
	node.r.x = mapping(node.lazy, node.r.x, r-m)
	node.r.lazy = composition(node.lazy, node.r.lazy)
	node.lazy = ds.lazyUnit
}

func (ds *DynamicSegTreeLazy) _buildRec(left, right int, nums []E) *SegNode {
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

func (ds *DynamicSegTreeLazy) _setRec(root *SegNode, l, r, i int, x E) *SegNode {
	if l == r-1 {
		root = ds.Copy(root)
		root.x = x
		root.lazy = ds.lazyUnit
		return root
	}
	ds._pushDown(root, l, r)
	m := (l + r) >> 1
	if root.l == nil {
		root.l = ds._newNode(l, m)
	}
	if root.r == nil {
		root.r = ds._newNode(m, r)
	}
	root = ds.Copy(root)
	if i < m {
		root.l = ds._setRec(root.l, l, m, i, x)
	} else {
		root.r = ds._setRec(root.r, m, r, i, x)
	}
	root.x = op(root.l.x, root.r.x)
	return root
}

func (ds *DynamicSegTreeLazy) _updateRec(root *SegNode, l, r, i int, x E) *SegNode {
	if l == r-1 {
		root = ds.Copy(root)
		root.x = op(root.x, x)
		root.lazy = ds.lazyUnit
		return root
	}
	ds._pushDown(root, l, r)
	m := (l + r) >> 1
	if root.l == nil {
		root.l = ds._newNode(l, m)
	}
	if root.r == nil {
		root.r = ds._newNode(m, r)
	}
	root = ds.Copy(root)
	if i < m {
		root.l = ds._updateRec(root.l, l, m, i, x)
	} else {
		root.r = ds._updateRec(root.r, m, r, i, x)
	}
	root.x = op(root.l.x, root.r.x)
	return root
}

func (ds *DynamicSegTreeLazy) _queryRec(root *SegNode, l, r, ql, qr int, x *E, lazy Id) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr {
		return
	}
	if root == nil {
		*x = op(*x, mapping(lazy, e2(ql, qr), qr-ql))
		return
	}
	if l == ql && r == qr {
		*x = op(*x, mapping(lazy, root.x, r-l))
		return
	}
	m := (l + r) >> 1
	lazy = composition(lazy, root.lazy)
	ds._queryRec(root.l, l, m, ql, qr, x, lazy)
	ds._queryRec(root.r, m, r, ql, qr, x, lazy)
}

func (ds *DynamicSegTreeLazy) _updateRangeRec(root *SegNode, l, r, ql, qr int, lazy Id) *SegNode {
	if root == nil {
		root = ds._newNode(l, r)
	}
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr {
		return root
	}
	if l == ql && r == qr {
		root = ds.Copy(root)
		root.x = mapping(lazy, root.x, r-l)
		root.lazy = composition(lazy, root.lazy)
		return root
	}
	ds._pushDown(root, l, r)
	m := (l + r) >> 1
	root = ds.Copy(root)
	root.l = ds._updateRangeRec(root.l, l, m, ql, qr, lazy)
	root.r = ds._updateRangeRec(root.r, m, r, ql, qr, lazy)
	root.x = op(root.l.x, root.r.x)
	return root
}

func (ds *DynamicSegTreeLazy) _minLeftRec(root *SegNode, l, r, qr int, check func(E) bool, x *E) int {
	if qr <= l {
		return l
	}
	if root == nil {
		root = ds._newNode(l, r)
	}
	qr = min(qr, r)
	if r == qr && check(op(root.x, *x)) {
		*x = op(root.x, *x)
		return l
	}
	if r == l+1 {
		return r
	}
	ds._pushDown(root, l, r)
	m := (l + r) >> 1
	k := ds._minLeftRec(root.r, m, r, qr, check, x)
	if m < k {
		return k
	}
	return ds._minLeftRec(root.l, l, m, qr, check, x)
}

func (ds *DynamicSegTreeLazy) _maxRightRec(root *SegNode, l, r, ql int, check func(E) bool, x *E) int {
	if r <= ql {
		return r
	}
	if root == nil {
		root = ds._newNode(l, r)
	}
	ql = max(ql, l)
	if l == ql && check(op(*x, root.x)) {
		*x = op(*x, root.x)
		return r
	}
	if r == l+1 {
		return l
	}
	ds._pushDown(root, l, r)
	m := (l + r) >> 1
	k := ds._maxRightRec(root.l, l, m, ql, check, x)
	if m > k {
		return k
	}
	return ds._maxRightRec(root.r, m, r, ql, check, x)
}

func (ds *DynamicSegTreeLazy) _getAllRec(root *SegNode, l, r int, res *[]E, lazy Id) {
	if root == nil {
		root = ds._newNode(l, r)
	}
	if r-l == 1 {
		*res = append(*res, mapping(lazy, root.x, 1))
		return
	}
	m := (l + r) >> 1
	lazy = composition(lazy, root.lazy)
	ds._getAllRec(root.l, l, m, res, lazy)
	ds._getAllRec(root.r, m, r, res, lazy)
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
