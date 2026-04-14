// G - Copy Query
// https://atcoder.jp/contests/abc453/tasks/abc453_g
//
// 有 $N$ 个长度为 $M$ 的数列 $A_1, A_2, \dots, A_N$。初始所有元素均为 $0$。
// 你需要按顺序处理 $Q$ 个操作：
//
// 1.  **复制 (Type 1)**：`1 X Y`
//     将数列 $A_Y$ 的内容完全复制给 $A_X$（即 $A_X = A_Y$）。
// 2.  **单点更新 (Type 2)**：`2 X Y Z`
//     将第 $X$ 个数列的第 $Y$ 个元素修改为 $Z$（即 $A_{X,Y} = Z$）。
// 3.  **区间查询 (Type 3)**：`3 X L R`
//     计算并输出第 $X$ 个数列在区间 $[L, R]$ 内的元素之和。
//
// 数据范围
// - $1 \le N, M, Q \le 2 \times 10^5$
// - $0 \le Z \le 10^9$
// - $1 \le L \le R \le M$
// - 时间限制：2s，内存限制：1024MiB
//
// 这类问题通常涉及**可持久化数据结构**（如可持久化线段树）。
// - 操作 1 类似于将 $A_X$ 的当前根节点指向 $A_Y$ 的根节点。
// - 操作 2 类似于在 $A_X$ 序列对应的版本上进行单点修改，产生新版本。
// - 操作 3 则是对特定版本的数据结构进行区间求和。

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

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	seg := NewDynamicSegTreeSparseG(
		func() int { return 0 },
		func(a, b int) int { return a + b },
		0,
		m,
		true,
	)
	data := make([]*SegNode[int], n)
	for i := range n {
		data[i] = seg.NewRoot()
	}

	op1 := func(x, y int) {
		data[x] = data[y]
	}

	op2 := func(x, y, z int) {
		data[x] = seg.Set(data[x], y, z)
	}

	op3 := func(x, l, r int) {
		res := seg.Query(data[x], l, r)
		fmt.Fprintln(out, res)
	}

	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		switch t {
		case 1:
			var x, y int
			fmt.Fscan(in, &x, &y)
			op1(x-1, y-1)
		case 2:
			var x, y, z int
			fmt.Fscan(in, &x, &y, &z)
			op2(x-1, y-1, z)
		case 3:
			var x, l, r int
			fmt.Fscan(in, &x, &l, &r)
			op3(x-1, l-1, r)
		}
	}
}

type DynamicSegTreeSparseG[E any] struct {
	L, R       int
	persistent bool
	unit       E
	e          func() E
	op         func(E, E) E
}

type SegNode[E any] struct {
	idx       int
	l, r      *SegNode[E]
	data, sum E
}

// 指定 [left,right) 区间建立动态开点线段树.
func NewDynamicSegTreeSparseG[E any](
	e func() E, op func(E, E) E,
	left, right int, persistent bool) *DynamicSegTreeSparseG[E] {
	return &DynamicSegTreeSparseG[E]{
		L:          left,
		R:          right + 1,
		persistent: persistent,
		unit:       e(),
		e:          e,
		op:         op,
	}
}

func (ds *DynamicSegTreeSparseG[E]) NewRoot() *SegNode[E] { return nil }

// 查询区间 [left, right).
// L<=left<=right<=R
func (ds *DynamicSegTreeSparseG[E]) Query(root *SegNode[E], left, right int) E {
	if left < ds.L {
		left = ds.L
	}
	if right > ds.R {
		right = ds.R
	}
	if left >= right {
		return ds.unit
	}
	x := ds.unit
	ds._queryRec(root, ds.L, ds.R, left, right, &x)
	return x
}

func (ds *DynamicSegTreeSparseG[E]) QueryAll(root *SegNode[E]) E {
	return ds.Query(root, ds.L, ds.R)
}

// L<=index<R
func (ds *DynamicSegTreeSparseG[E]) Set(root *SegNode[E], index int, value E) *SegNode[E] {
	return ds._setRec(root, ds.L, ds.R, index, value)
}

func (ds *DynamicSegTreeSparseG[E]) Get(root *SegNode[E], index int) E {
	return ds._getRec(root, index)
}

// L<=left<R
func (ds *DynamicSegTreeSparseG[E]) Update(root *SegNode[E], index int, value E) *SegNode[E] {
	return ds._updateRec(root, ds.L, ds.R, index, value)
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 check.
// L<=right<=R
func (ds *DynamicSegTreeSparseG[E]) MinLeft(root *SegNode[E], right int, check func(E) bool) int {
	x := ds.unit
	return ds._minLeftRec(root, ds.L, ds.R, right, check, &x)
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 check.
// L<=left<=R
func (ds *DynamicSegTreeSparseG[E]) MaxRight(root *SegNode[E], left int, check func(E) bool) int {
	x := ds.unit
	return ds._maxRightRec(root, ds.L, ds.R, left, check, &x)
}

func (ds *DynamicSegTreeSparseG[E]) GetAll(root *SegNode[E]) []struct {
	index int
	value E
} {
	res := make([]struct {
		index int
		value E
	}, 0)
	ds._getAllRec(root, &res)
	return res
}

func (ds *DynamicSegTreeSparseG[E]) Copy(node *SegNode[E]) *SegNode[E] {
	if node == nil || !ds.persistent {
		return node
	}
	return &SegNode[E]{idx: node.idx, l: node.l, r: node.r, data: node.data, sum: node.sum}
}

func (ds *DynamicSegTreeSparseG[E]) _pushUp(node *SegNode[E]) {
	node.sum = node.data
	if node.l != nil {
		node.sum = ds.op(node.l.sum, node.sum)
	}
	if node.r != nil {
		node.sum = ds.op(node.sum, node.r.sum)
	}
}

func (ds *DynamicSegTreeSparseG[E]) _newNode(idx int, x E) *SegNode[E] {
	return &SegNode[E]{idx: idx, data: x, sum: x}
}

func (ds *DynamicSegTreeSparseG[E]) _setRec(root *SegNode[E], l, r, i int, x E) *SegNode[E] {
	if root == nil {
		root = ds._newNode(i, x)
		return root
	}
	root = ds.Copy(root)
	if root.idx == i {
		root.data = x
		ds._pushUp(root)
		return root
	}
	m := (l + r) >> 1
	if i < m {
		if root.idx < i {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.l = ds._setRec(root.l, l, m, i, x)
	} else {
		if i < root.idx {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.r = ds._setRec(root.r, m, r, i, x)
	}
	ds._pushUp(root)
	return root
}

func (ds *DynamicSegTreeSparseG[E]) _updateRec(root *SegNode[E], l, r, i int, x E) *SegNode[E] {
	if root == nil {
		root = ds._newNode(i, x)
		return root
	}
	root = ds.Copy(root)
	if root.idx == i {
		root.data = ds.op(root.data, x)
		ds._pushUp(root)
		return root
	}
	m := (l + r) >> 1
	if i < m {
		if root.idx < i {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.l = ds._updateRec(root.l, l, m, i, x)
	} else {
		if i < root.idx {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.r = ds._updateRec(root.r, m, r, i, x)
	}
	ds._pushUp(root)
	return root
}

func (ds *DynamicSegTreeSparseG[E]) _queryRec(root *SegNode[E], l, r, ql, qr int, x *E) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr || root == nil {
		return
	}
	if l == ql && r == qr {
		*x = ds.op(*x, root.sum)
		return
	}
	m := (l + r) >> 1
	ds._queryRec(root.l, l, m, ql, qr, x)
	if ql <= root.idx && root.idx < qr {
		*x = ds.op(*x, root.data)
	}
	ds._queryRec(root.r, m, r, ql, qr, x)
}

func (ds *DynamicSegTreeSparseG[E]) _minLeftRec(root *SegNode[E], l, r, qr int, check func(E) bool, x *E) int {
	if root == nil || qr <= l {
		return ds.L
	}
	if check(ds.op(root.sum, *x)) {
		*x = ds.op(root.sum, *x)
		return ds.L
	}
	m := (l + r) >> 1
	k := ds._minLeftRec(root.r, m, r, qr, check, x)
	if k != ds.L {
		return k
	}
	if root.idx < qr {
		*x = ds.op(root.data, *x)
		if !check(*x) {
			return root.idx + 1
		}
	}
	return ds._minLeftRec(root.l, l, m, qr, check, x)
}

func (ds *DynamicSegTreeSparseG[E]) _maxRightRec(root *SegNode[E], l, r, ql int, check func(E) bool, x *E) int {
	if root == nil || r <= ql {
		return ds.R
	}
	if check(ds.op(*x, root.sum)) {
		*x = ds.op(*x, root.sum)
		return ds.R
	}
	m := (l + r) >> 1
	k := ds._maxRightRec(root.l, l, m, ql, check, x)
	if k != ds.R {
		return k
	}
	if ql <= root.idx {
		*x = ds.op(*x, root.data)
		if !check(*x) {
			return root.idx
		}
	}
	return ds._maxRightRec(root.r, m, r, ql, check, x)
}

func (ds *DynamicSegTreeSparseG[E]) _getAllRec(root *SegNode[E], res *[]struct {
	index int
	value E
}) {
	if root == nil {
		return
	}
	ds._getAllRec(root.l, res)
	*res = append(*res, struct {
		index int
		value E
	}{root.idx, root.data})
	ds._getAllRec(root.r, res)
}

func (ds *DynamicSegTreeSparseG[E]) _getRec(root *SegNode[E], idx int) E {
	if root == nil {
		return ds.unit
	}
	if idx == root.idx {
		return root.data
	}
	if idx < root.idx {
		return ds._getRec(root.l, idx)
	}
	return ds._getRec(root.r, idx)
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
