// !递归版线段树(适合维护区间信息复杂的情况)

package main

import "math/bits"

type LazySegmentTree struct {
	n    int
	data []int
	lazy []int
	// 别的一些信息
}

func NewLazySegmentTree(leaves []int) *LazySegmentTree {
	cap := 1 << (bits.Len(uint(len(leaves)-1)) + 1)
	// !初始化data和lazy数组 然后建树
	tree := &LazySegmentTree{
		n:    len(leaves),
		data: make([]int, cap),
		lazy: make([]int, cap),
	}
	tree._build(1, 1, tree.n, leaves)
	return tree
}

// TODO
func (t *LazySegmentTree) _build(root, left, right int, leaves []int) {
	if left == right {
		// !初始化叶子结点信息 例如data和lazy的monoid
		t.data[root] = leaves[left-1]
		return
	}
	mid := (left + right) >> 1
	t._build(root<<1, left, mid, leaves)
	t._build(root<<1|1, mid+1, right, leaves)
	t._pushUp(root, left, right)
}

func (t *LazySegmentTree) _pushUp(root, left, right int) {
	// !op操作更新root结点的data信息
	t.data[root] = t.data[root<<1] + t.data[root<<1|1]
}

func (t *LazySegmentTree) _pushDown(root, left, right int) {
	// !传播lazy信息(可以判断根的lazy不为monoid时才传播,传播后将根的lazy置为monoid)
	if t.lazy[root] != 0 { // monoid
		mid := (left + right) >> 1
		t._propagate(root<<1, left, mid, t.lazy[root])
		t._propagate(root<<1|1, mid+1, right, t.lazy[root])
		t.lazy[root] = 0 // monoid
	}
}

func (t *LazySegmentTree) _propagate(root, left, right, lazy int) {
	// !mapping + composition 来更新子节点data和lazy信息
	t.data[root] += lazy * (right - left + 1)
	t.lazy[root] += lazy
}

func (t *LazySegmentTree) _query(root, L, R, l, r int) int {
	if L <= l && r <= R {
		return t.data[root]
	}

	t._pushDown(root, l, r)
	mid := (l + r) >> 1
	res := 0 // monoid
	if L <= mid {
		res += t._query(root<<1, L, R, l, mid) // op
	}
	if R > mid {
		res += t._query(root<<1|1, L, R, mid+1, r) // op
	}
	return res
}

func (t *LazySegmentTree) _update(root, L, R, l, r, val int) {
	if L <= l && r <= R {
		t._propagate(root, l, r, val)
		return
	}

	t._pushDown(root, l, r)
	mid := (l + r) >> 1
	if L <= mid {
		t._update(root<<1, L, R, l, mid, val)
	}
	if R > mid {
		t._update(root<<1|1, L, R, mid+1, r, val)
	}
	t._pushUp(root, l, r)
}

// public api
func (t *LazySegmentTree) Query(left, right int) int    { return t._query(1, left, right, 1, t.n) }
func (t *LazySegmentTree) Update(left, right, lazy int) { t._update(1, left, right, 1, t.n, lazy) }
func (t *LazySegmentTree) QueryAll() int                { return t.data[1] }
