// 可持久化Dual线段树.
//
//  Api:
//
//  - NewSegmentTreeDynamicDual()
//
//  - NewRoot()
//  - NewNode()
//  - NewNodeFrom()
//
//  - Get()
//  - Update()
//  - CopyInterval()
//  - GetAll()

package main

import "fmt"

func main() {
	demo()
}

func demo() {
	id := func() int { return 0 }
	composition := func(f, g int) int { return f + g }
	st := NewSegmentTreeDynamicDual(id, composition, 0, 10, true)
	root := st.NewRoot()
	root = st.Update(root, 0, 10, 1)
	fmt.Println(st.GetAll(root))
	root = st.Update(root, 0, 5, 2)
	fmt.Println(st.GetAll(root))
	root = st.Update(root, 5, 10, 3)
	fmt.Println(st.GetAll(root))

	newRoot := st.NewRoot()
	newRoot = st.CopyInterval(newRoot, root, 0, 7, 4)
	fmt.Println(st.GetAll(newRoot), st.GetAll(root))
}

type Node[Id comparable] struct {
	l, r *Node[Id]
	x    Id
}

type SegmentTreeDynamicDual[Id comparable] struct {
	id          func() Id
	composition func(f, g Id) Id
	l0, r0      int
	persistent  bool
}

func NewSegmentTreeDynamicDual[Id comparable](
	id func() Id, composition func(f, g Id) Id,
	l0, r0 int,
	persistent bool,
) *SegmentTreeDynamicDual[Id] {
	return &SegmentTreeDynamicDual[Id]{
		id:          id,
		composition: composition,
		l0:          l0,
		r0:          r0,
		persistent:  persistent,
	}
}

func (st *SegmentTreeDynamicDual[Id]) NewRoot() *Node[Id] {
	return st.NewNode(nil, nil, st.id())
}

func (st *SegmentTreeDynamicDual[Id]) NewNode(l, r *Node[Id], x Id) *Node[Id] {
	return &Node[Id]{l: l, r: r, x: x}
}

func (st *SegmentTreeDynamicDual[Id]) NewNodeFrom(n int, f func(int) Id) *Node[Id] {
	if !(st.l0 == 0 && st.r0 == n) {
		panic("invalid range")
	}
	var dfs func(l, r int) *Node[Id]
	dfs = func(l, r int) *Node[Id] {
		if l == r {
			return nil
		}
		if r == l+1 {
			return st.NewNode(nil, nil, f(l))
		}
		m := (l + r) >> 1
		lRoot := dfs(l, m)
		rRoot := dfs(m, r)
		x := st.composition(lRoot.x, rRoot.x)
		return st.NewNode(lRoot, rRoot, x)
	}
	return dfs(0, n)
}

func (st *SegmentTreeDynamicDual[Id]) Get(root *Node[Id], i int) Id {
	if root == nil {
		return st.id()
	}
	x := st.id()
	st.getRec(root, st.l0, st.r0, i, &x)
	return x
}

func (st *SegmentTreeDynamicDual[Id]) Update(root *Node[Id], l, r int, x Id) *Node[Id] {
	if l >= r {
		return root
	}
	root = st.copyNode(root)
	st.updateRec(root, st.l0, st.r0, l, r, x)
	return root
}

// 将other[l:r)的值与x结合，然后覆盖root[l:r)的值.
func (st *SegmentTreeDynamicDual[Id]) CopyInterval(root, other *Node[Id], l, r int, x Id) *Node[Id] {
	if root == other {
		return root
	}
	root = st.copyNode(root)
	st.copyIntervalRec(root, other, st.l0, st.r0, l, r, x)
	return root
}

func (st *SegmentTreeDynamicDual[Id]) GetAll(root *Node[Id]) []Id {
	var res []Id
	var dfs func(c *Node[Id], l, r int, x Id)
	dfs = func(c *Node[Id], l, r int, x Id) {
		if c == nil {
			for i := l; i < r; i++ {
				res = append(res, x)
			}
			return
		}
		x = st.composition(c.x, x)
		if r == l+1 {
			res = append(res, x)
			return
		}
		m := (l + r) >> 1
		dfs(c.l, l, m, x)
		dfs(c.r, m, r, x)
	}
	dfs(root, st.l0, st.r0, st.id())
	return res
}

func (st *SegmentTreeDynamicDual[Id]) copyNode(c *Node[Id]) *Node[Id] {
	if c == nil || !st.persistent {
		return c
	}
	return st.NewNode(c.l, c.r, c.x)
}

func (st *SegmentTreeDynamicDual[Id]) updateRec(c *Node[Id], l, r, ql, qr int, a Id) {
	ql = max(ql, l)
	qr = min(qr, r)
	if a == st.id() || ql >= qr {
		return
	}
	if l == ql && r == qr {
		c.x = st.composition(c.x, a)
		return
	}
	if c.l == nil {
		c.l = st.NewNode(nil, nil, st.id())
	} else {
		c.l = st.copyNode(c.l)
	}
	if c.r == nil {
		c.r = st.NewNode(nil, nil, st.id())
	} else {
		c.r = st.copyNode(c.r)
	}
	c.l.x = st.composition(c.l.x, c.x)
	c.r.x = st.composition(c.r.x, c.x)
	c.x = st.id()
	m := (l + r) >> 1
	st.updateRec(c.l, l, m, ql, qr, a)
	st.updateRec(c.r, m, r, ql, qr, a)
}

func (st *SegmentTreeDynamicDual[Id]) copyIntervalRec(c, d *Node[Id], l, r, ql, qr int, x Id) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr {
		return
	}
	if l == ql && r == qr {
		if d != nil {
			c.x = st.composition(d.x, x)
			c.l = d.l
			c.r = d.r
		} else {
			c.x = x
			c.l = nil
			c.r = nil
		}
		return
	}
	if c.l == nil {
		c.l = st.NewNode(nil, nil, st.id())
	} else {
		c.l = st.copyNode(c.l)
	}
	if c.r == nil {
		c.r = st.NewNode(nil, nil, st.id())
	} else {
		c.r = st.copyNode(c.r)
	}
	c.l.x = st.composition(c.l.x, c.x)
	c.r.x = st.composition(c.r.x, c.x)
	c.x = st.id()
	m := (l + r) >> 1
	if d != nil {
		x = st.composition(d.x, x)
	}
	if d == nil {
		st.copyIntervalRec(c.l, nil, l, m, ql, qr, x)
		st.copyIntervalRec(c.r, nil, m, r, ql, qr, x)
	} else {
		st.copyIntervalRec(c.l, d.l, l, m, ql, qr, x)
		st.copyIntervalRec(c.r, d.r, m, r, ql, qr, x)
	}
}

func (st *SegmentTreeDynamicDual[Id]) getRec(c *Node[Id], l, r, i int, x *Id) {
	if c == nil {
		return
	}
	*x = st.composition(c.x, *x)
	if r == l+1 {
		return
	}
	m := (l + r) >> 1
	if i < m {
		st.getRec(c.l, l, m, i, x)
	} else {
		st.getRec(c.r, m, r, i, x)
	}
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
