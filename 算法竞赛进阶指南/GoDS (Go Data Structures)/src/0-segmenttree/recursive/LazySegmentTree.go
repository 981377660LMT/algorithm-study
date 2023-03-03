// !递归版线段树(适合维护区间信息复杂的情况)

package main

import "math/bits"

type E = struct{ zero, one, inv int }
type Id = bool

func (*LazySegmentTree) e() E   { return E{} }
func (*LazySegmentTree) id() Id { return false }
func (*LazySegmentTree) op(left, right E) E {
	return E{
		zero: left.zero + right.zero,
		one:  left.one + right.one,
		inv:  left.inv + right.inv + left.one*right.zero,
	}
}

func (*LazySegmentTree) mapping(lazy Id, data E) E {
	if !lazy {
		return data
	}
	return E{
		zero: data.one,
		one:  data.zero,
		inv:  data.one*data.zero - data.inv,
	}
}
func (*LazySegmentTree) composition(parentLazy, childLazy Id) Id {
	return (parentLazy && !childLazy) || (!parentLazy && childLazy)
}

type LazySegmentTree struct {
	n    int
	data []E
	lazy []Id
}

func NewLazySegmentTree(leaves []E) *LazySegmentTree {
	tree := &LazySegmentTree{}
	n := int(len(leaves))
	tree.n = n
	log := int(bits.Len(uint(n - 1)))
	size := 1 << log
	tree.data = make([]E, 2*size)
	tree.lazy = make([]Id, size) // !叶子结点不需要更新lazy (composition)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	// !初始化data和lazy数组 然后建树
	tree.build(1, 1, n, leaves)
	return tree
}

func (t *LazySegmentTree) build(root, left, right int, leaves []E) {
	if left == right {
		// !初始化叶子结点信息 例如data和lazy的monoid
		t.data[root] = leaves[left-1]
		return
	}
	mid := (left + right) >> 1
	t.build(root<<1, left, mid, leaves)
	t.build(root<<1|1, mid+1, right, leaves)
	t.pushUp(root)
}

func (t *LazySegmentTree) pushUp(root int) {
	t.data[root] = t.op(t.data[root<<1], t.data[root<<1|1])
}

func (t *LazySegmentTree) pushDown(root int) {
	if t.lazy[root] != t.id() {
		t.propagate(root<<1, t.lazy[root])
		t.propagate(root<<1|1, t.lazy[root])
		t.lazy[root] = t.id()
	}
}

func (t *LazySegmentTree) propagate(root int, lazy Id) {
	t.data[root] = t.mapping(lazy, t.data[root])
	if root < len(t.lazy) { // !叶子结点不需要更新lazy (composition)
		t.lazy[root] = t.composition(lazy, t.lazy[root])
	}
}

func (t *LazySegmentTree) query(root, L, R, l, r int) E {
	if L <= l && r <= R {
		return t.data[root]
	}

	t.pushDown(root)
	mid := (l + r) >> 1
	res := t.e()
	if L <= mid {
		res = t.op(res, t.query(root<<1, L, R, l, mid))
	}
	if R > mid {
		res = t.op(res, t.query(root<<1|1, L, R, mid+1, r))
	}
	return res
}

func (t *LazySegmentTree) update(root, L, R, l, r int, val Id) {
	if L <= l && r <= R {
		t.propagate(root, val)
		return
	}

	t.pushDown(root)
	mid := (l + r) >> 1
	if L <= mid {
		t.update(root<<1, L, R, l, mid, val)
	}
	if R > mid {
		t.update(root<<1|1, L, R, mid+1, r, val)
	}
	t.pushUp(root)
}

// 查询闭区间[left,right]的值
//  1<=left<=right<=n
func (t *LazySegmentTree) Query(left, right int) E { return t.query(1, left, right, 1, t.n) }

// 更新闭区间[left,right]的值
//  1<=left<=right<=n
func (t *LazySegmentTree) Update(left, right int, lazy Id) { t.update(1, left, right, 1, t.n, lazy) }
func (t *LazySegmentTree) QueryAll() E                     { return t.data[1] }
