// 区间修改, 区间查询

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	persistentRangeAffineRangeSum()
}

// https://leetcode.cn/problems/range-module/
type RangeModule struct {
	segmentTree *DynamicSegTreeLazyG[int, int]
	root        *SegNodeG[int, int]
}

func Constructor2() RangeModule {
	res := RangeModule{}

	// RangeAssignRangeSum
	e1 := func() int { return 0 }
	e2 := func(start, end int) int { return 0 }
	id := func() int { return -1 }
	op := func(a, b int) int { return a + b }
	mapping := func(f int, g int, size int) int {
		if f == -1 {
			return g
		}
		return f * size
	}
	composition := func(f int, g int) int {
		if f == -1 {
			return g
		}
		return f
	}

	res.segmentTree = NewDynamicSegTreeLazyG(e1, e2, id, op, mapping, composition, 0, 1e9+10, false)
	res.root = res.segmentTree.NewRoot()
	return res
}

func (this *RangeModule) AddRange(left int, right int) {
	this.segmentTree.UpdateRange(this.root, left, right, 1)
}

func (this *RangeModule) QueryRange(left int, right int) bool {
	return this.segmentTree.Query(this.root, left, right) == right-left
}

func (this *RangeModule) RemoveRange(left int, right int) {
	this.segmentTree.UpdateRange(this.root, left, right, 0)
}

// https://judge.yosupo.jp/problem/persistent_range_affine_range_sum
func persistentRangeAffineRangeSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353
	type E = int
	type Id = struct{ mul, add int }

	var n, q int
	fmt.Fscan(in, &n, &q)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}

	e1 := func() E { return 0 }
	e2 := func(start, end int) E { return 0 }
	id := func() Id { return Id{1, 0} }
	op := func(a, b E) E { return (a + b) % MOD }
	mapping := func(f Id, g E, size int) E {
		return (f.mul*g + f.add*size) % MOD
	}
	composition := func(f, g Id) Id {
		return Id{f.mul * g.mul % MOD, (f.mul*g.add + f.add) % MOD}
	}

	seg := NewDynamicSegTreeLazyG[E, Id](e1, e2, id, op, mapping, composition, 0, n, true)
	roots := make([]*SegNodeG[E, Id], q+1)
	roots[0] = seg.Build(A)

	for i := 1; i <= q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 0 {
			var k, l, r, b, c int
			fmt.Fscan(in, &k, &l, &r, &b, &c)
			k++
			roots[i] = seg.UpdateRange(roots[k], l, r, Id{b, c})
		}
		if t == 1 {
			var k, s, l, r int
			fmt.Fscan(in, &k, &s, &l, &r)
			k++
			s++
			roots[i] = seg.CopyInterval(roots[k], roots[s], l, r, Id{1, 0})
		}
		if t == 2 {
			var k, l, r int
			fmt.Fscan(in, &k, &l, &r)
			k++
			roots[i] = roots[k]
			res := seg.Query(roots[k], l, r)
			fmt.Fprintln(out, res)
		}
	}
}

type DynamicSegTreeLazyG[E any, Id comparable] struct {
	e1          func() E
	e2          func(start, end int) E
	id          func() Id
	op          func(a, b E) E
	mapping     func(f Id, g E, size int) E
	composition func(f, g Id) Id

	L, R       int
	persistent bool
	dataUnit   E
	lazyUnit   Id
}

type SegNodeG[E any, Id comparable] struct {
	x    E
	lazy Id
	l, r *SegNodeG[E, Id]
}

func NewDynamicSegTreeLazyG[E any, Id comparable](
	e1 func() E,
	e2 func(start, end int) E,
	id func() Id,
	op func(a, b E) E,
	mapping func(f Id, g E, size int) E,
	composition func(f, g Id) Id,
	start, end int, persistent bool,
) *DynamicSegTreeLazyG[E, Id] {
	return &DynamicSegTreeLazyG[E, Id]{
		e1:          e1,
		e2:          e2,
		id:          id,
		op:          op,
		mapping:     mapping,
		composition: composition,
		L:           start,
		R:           end,
		persistent:  persistent,
		dataUnit:    e1(),
		lazyUnit:    id(),
	}
}

func (ds *DynamicSegTreeLazyG[E, Id]) NewRoot() *SegNodeG[E, Id] {
	return &SegNodeG[E, Id]{x: ds.e2(ds.L, ds.R), lazy: ds.lazyUnit}
}

func (ds *DynamicSegTreeLazyG[E, Id]) Build(nums []E) *SegNodeG[E, Id] {
	if !(ds.L == 0 && ds.R == len(nums)) {
		panic("invalid range")
	}
	return ds._buildRec(0, len(nums), nums)
}

func (ds *DynamicSegTreeLazyG[E, Id]) Query(root *SegNodeG[E, Id], start, end int) E {
	if start < ds.L {
		start = ds.L
	}
	if end > ds.R {
		end = ds.R
	}
	if start >= end {
		return ds.dataUnit
	}
	x := ds.dataUnit
	ds._queryRec(root, ds.L, ds.R, start, end, &x, ds.lazyUnit)
	return x
}

func (ds *DynamicSegTreeLazyG[E, Id]) QueryAll(root *SegNodeG[E, Id]) E {
	return ds.Query(root, ds.L, ds.R)
}

// L<=index<R
func (ds *DynamicSegTreeLazyG[E, Id]) Set(root *SegNodeG[E, Id], index int, value E) *SegNodeG[E, Id] {
	return ds._setRec(root, ds.L, ds.R, index, value)
}

// L<=index<R
func (ds *DynamicSegTreeLazyG[E, Id]) Update(root *SegNodeG[E, Id], index int, value E) *SegNodeG[E, Id] {
	return ds._updateRec(root, ds.L, ds.R, index, value)
}

// L<=start<=end<=R
func (ds *DynamicSegTreeLazyG[E, Id]) UpdateRange(root *SegNodeG[E, Id], start, end int, lazy Id) *SegNodeG[E, Id] {
	if start == end {
		return root
	}
	return ds._updateRangeRec(root, ds.L, ds.R, start, end, lazy)
}

// 二分查询最大的 right 使得切片 [left:right) 内的值满足 check.
// L<=right<=R
func (ds *DynamicSegTreeLazyG[E, Id]) MinLeft(root *SegNodeG[E, Id], end int, check func(E) bool) int {
	x := ds.dataUnit
	return ds._minLeftRec(root, ds.L, ds.R, end, check, &x)
}

// 二分查询最小的 left 使得切片 [left:right) 内的值满足 check.
// L<=left<=R
func (ds *DynamicSegTreeLazyG[E, Id]) MaxRight(root *SegNodeG[E, Id], start int, check func(E) bool) int {
	x := ds.dataUnit
	return ds._maxRightRec(root, ds.L, ds.R, start, check, &x)
}

func (ds *DynamicSegTreeLazyG[E, Id]) GetAll(root *SegNodeG[E, Id]) []E {
	res := make([]E, 0)
	ds._getAllRec(root, ds.L, ds.R, &res, ds.lazyUnit)
	return res
}

func (ds *DynamicSegTreeLazyG[E, Id]) EnumerateAll(root *SegNodeG[E, Id], f func(index int, value E)) {
	var dfs func(c *SegNodeG[E, Id], l, r int, lazy Id)
	dfs = func(c *SegNodeG[E, Id], l, r int, lazy Id) {
		if c == nil {
			return
		}
		if r-l == 1 {
			f(l, ds.mapping(lazy, c.x, 1))
			return
		}
		m := (l + r) >> 1
		lazy = ds.composition(lazy, c.lazy)
		dfs(c.l, l, m, lazy)
		dfs(c.r, m, r, lazy)
	}
	dfs(root, ds.L, ds.R, ds.lazyUnit)
}

func (ds *DynamicSegTreeLazyG[E, Id]) Copy(node *SegNodeG[E, Id]) *SegNodeG[E, Id] {
	if node == nil || !ds.persistent {
		return node
	}
	return &SegNodeG[E, Id]{l: node.l, r: node.r, x: node.x, lazy: node.lazy}
}

// 将 root[l:r) 用 apply(other[l:r),a) 的值覆盖, 返回新的 root.
func (ds *DynamicSegTreeLazyG[E, Id]) CopyInterval(root *SegNodeG[E, Id], other *SegNodeG[E, Id], l, r int, lazy Id) *SegNodeG[E, Id] {
	if root == other {
		return root
	}
	root = ds.Copy(root)
	ds._copyIntervalRec(root, other, ds.L, ds.R, l, r, lazy)
	return root
}

func (ds *DynamicSegTreeLazyG[E, Id]) _copyIntervalRec(c, d *SegNodeG[E, Id], l, r, ql, qr int, lazy Id) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr {
		return
	}
	if l == ql && r == qr {
		if d != nil {
			c.x = ds.mapping(lazy, d.x, r-l)
			c.lazy = ds.composition(lazy, d.lazy)
			c.l = d.l
			c.r = d.r
		} else {
			c.x = ds.mapping(lazy, ds.e2(l, r), r-l)
			c.lazy = lazy
			c.l = nil
			c.r = nil
		}
		return
	}

	m := (l + r) >> 1
	if c.l == nil {
		c.l = ds._newNode(ds.L, ds.R)
	} else {
		c.l = ds.Copy(c.l)
	}
	if c.r == nil {
		c.r = ds._newNode(ds.L, ds.R)
	} else {
		c.r = ds.Copy(c.r)
	}
	c.l.x = ds.mapping(c.lazy, c.l.x, m-l)
	c.l.lazy = ds.composition(c.lazy, c.l.lazy)
	c.r.x = ds.mapping(c.lazy, c.r.x, r-m)
	c.r.lazy = ds.composition(c.lazy, c.r.lazy)
	c.lazy = ds.lazyUnit
	if d != nil {
		lazy = ds.composition(lazy, d.lazy)
	}
	if d != nil {
		ds._copyIntervalRec(c.l, d.l, l, m, ql, qr, lazy)
		ds._copyIntervalRec(c.r, d.r, m, r, ql, qr, lazy)
	} else {
		ds._copyIntervalRec(c.l, nil, l, m, ql, qr, lazy)
		ds._copyIntervalRec(c.r, nil, m, r, ql, qr, lazy)
	}
	c.x = ds.op(c.l.x, c.r.x)
}

func (ds *DynamicSegTreeLazyG[E, Id]) _newNode(left, right int) *SegNodeG[E, Id] {
	return &SegNodeG[E, Id]{x: ds.e2(left, right), lazy: ds.lazyUnit}
}

func (ds *DynamicSegTreeLazyG[E, Id]) _newNodeWithValue(x E) *SegNodeG[E, Id] {
	return &SegNodeG[E, Id]{x: x, lazy: ds.lazyUnit}
}

func (ds *DynamicSegTreeLazyG[E, Id]) _pushDown(node *SegNodeG[E, Id], l, r int) {
	if node.lazy == ds.lazyUnit {
		return
	}
	m := (l + r) >> 1
	if node.l == nil {
		node.l = ds._newNode(l, m)
	} else {
		node.l = ds.Copy(node.l)
	}
	node.l.x = ds.mapping(node.lazy, node.l.x, m-l)
	node.l.lazy = ds.composition(node.lazy, node.l.lazy)
	if node.r == nil {
		node.r = ds._newNode(m, r)
	} else {
		node.r = ds.Copy(node.r)
	}
	node.r.x = ds.mapping(node.lazy, node.r.x, r-m)
	node.r.lazy = ds.composition(node.lazy, node.r.lazy)
	node.lazy = ds.lazyUnit
}

func (ds *DynamicSegTreeLazyG[E, Id]) _buildRec(left, right int, nums []E) *SegNodeG[E, Id] {
	if left == right {
		return nil
	}
	if right == left+1 {
		return ds._newNodeWithValue(nums[left])
	}
	mid := (left + right) >> 1
	lRoot := ds._buildRec(left, mid, nums)
	rRoot := ds._buildRec(mid, right, nums)
	x := ds.op(lRoot.x, rRoot.x)
	root := ds._newNodeWithValue(x)
	root.l = lRoot
	root.r = rRoot
	return root
}

func (ds *DynamicSegTreeLazyG[E, Id]) _setRec(root *SegNodeG[E, Id], l, r, i int, x E) *SegNodeG[E, Id] {
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
	root.x = ds.op(root.l.x, root.r.x)
	return root
}

func (ds *DynamicSegTreeLazyG[E, Id]) _updateRec(root *SegNodeG[E, Id], l, r, i int, x E) *SegNodeG[E, Id] {
	if l == r-1 {
		root = ds.Copy(root)
		root.x = ds.op(root.x, x)
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
	root.x = ds.op(root.l.x, root.r.x)
	return root
}

func (ds *DynamicSegTreeLazyG[E, Id]) _queryRec(root *SegNodeG[E, Id], l, r, ql, qr int, x *E, lazy Id) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr {
		return
	}
	if root == nil {
		*x = ds.op(*x, ds.mapping(lazy, ds.e2(ql, qr), qr-ql))
		return
	}
	if l == ql && r == qr {
		*x = ds.op(*x, ds.mapping(lazy, root.x, r-l))
		return
	}
	m := (l + r) >> 1
	lazy = ds.composition(lazy, root.lazy)
	ds._queryRec(root.l, l, m, ql, qr, x, lazy)
	ds._queryRec(root.r, m, r, ql, qr, x, lazy)
}

func (ds *DynamicSegTreeLazyG[E, Id]) _updateRangeRec(root *SegNodeG[E, Id], l, r, ql, qr int, lazy Id) *SegNodeG[E, Id] {
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
		root.x = ds.mapping(lazy, root.x, r-l)
		root.lazy = ds.composition(lazy, root.lazy)
		return root
	}
	ds._pushDown(root, l, r)
	m := (l + r) >> 1
	root = ds.Copy(root)
	root.l = ds._updateRangeRec(root.l, l, m, ql, qr, lazy)
	root.r = ds._updateRangeRec(root.r, m, r, ql, qr, lazy)
	root.x = ds.op(root.l.x, root.r.x)
	return root
}

func (ds *DynamicSegTreeLazyG[E, Id]) _minLeftRec(root *SegNodeG[E, Id], l, r, qr int, check func(E) bool, x *E) int {
	if qr <= l {
		return l
	}
	if root == nil {
		root = ds._newNode(l, r)
	}
	qr = min(qr, r)
	if r == qr && check(ds.op(root.x, *x)) {
		*x = ds.op(root.x, *x)
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

func (ds *DynamicSegTreeLazyG[E, Id]) _maxRightRec(root *SegNodeG[E, Id], l, r, ql int, check func(E) bool, x *E) int {
	if r <= ql {
		return r
	}
	if root == nil {
		root = ds._newNode(l, r)
	}
	ql = max(ql, l)
	if l == ql && check(ds.op(*x, root.x)) {
		*x = ds.op(*x, root.x)
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

func (ds *DynamicSegTreeLazyG[E, Id]) _getAllRec(root *SegNodeG[E, Id], l, r int, res *[]E, lazy Id) {
	if root == nil {
		return
	}
	if r-l == 1 {
		*res = append(*res, ds.mapping(lazy, root.x, 1))
		return
	}
	m := (l + r) >> 1
	lazy = ds.composition(lazy, root.lazy)
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
