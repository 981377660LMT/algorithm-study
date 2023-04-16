//
// Euler Tour Tree
//
// Description:
//   Maintain dynamic unrooted tree with supporting
//   - make_node(x)       : return singleton with value x
//   - Link(u,v)          : add Link between u and v
//   - cut(uv)            : remove edge uv
//   - sum_in_component(u): return sum of all values in the component
//
// Algorithm:
//   Maintain euler tours by splay trees.
//   This data structure is originally proposed by Miltersen et al,
//   and independently by Fredman and Henzinger.
//
// Complexity:
//   O(log n)
//
// References:
//   M. L. Fredman and M. R. Henzinger (1998):
//   Lower bounds for fully dynamic connectivity problems in graphs.
//   Algorithmica, vol. 22, no. 3, pp. 351–362.
//
//   P. B. Miltersen, S. Subramanian, J. S. Vitter, and R. Tamassia (1994):
//   Complexity models for incremental computation.
//   Theoretical Computer Science, vol. 130. no. 1, pp. 203–236.

// API:
//  Link(u, v)    : add Link between u and v
//  Cut(uv)       : remove edge uv
//  GetSum(u)     : return sum of all values in the component

package main

import (
	"fmt"
	"time"
)

func main2() {
	ett := NewEulerTourTree()
	a, b, c, d, e, f, g := ett.Alloc(3), ett.Alloc(1), ett.Alloc(4), ett.Alloc(1), ett.Alloc(5), ett.Alloc(9), ett.Alloc(2)

	ab := ett.Link(a, b)
	ett.Link(a, c)
	ett.Link(b, d)
	ett.Link(b, e)
	cf := ett.Link(c, f)
	ett.Link(c, g)
	fmt.Println(ett.GetSum(a), ab, cf)
	ett.Cut(ab)
	fmt.Println(ett.GetSum(a), ab, cf)
}

func main() {
	E := NewEulerTourTree()
	a, b, c, d, e := E.Alloc(1), E.Alloc(2), E.Alloc(3), E.Alloc(4), E.Alloc(5)
	fmt.Println(E.GetSum(a), E.GetSum(b), E.GetSum(c), E.GetSum(d), E.GetSum(e))
	ab := E.Link(a, b)
	fmt.Println(E.GetSum(a), E.GetSum(b), E.GetSum(c), E.GetSum(d), E.GetSum(e))
	E.Cut(ab)
	fmt.Println(E.GetSum(a), E.GetSum(b), E.GetSum(c), E.GetSum(d), E.GetSum(e))

	// 1e5
	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		ab = E.Link(a, b)
		E.Cut(ab)
	}
	fmt.Println(time.Since(time1))
}

type _ENode struct {
	x, s int // value, sum
	ch   [2]*_ENode
	p, r *_ENode
}

type EulerTourTree struct {
}

func NewEulerTourTree() *EulerTourTree {
	return &EulerTourTree{}
}

func (ett *EulerTourTree) Alloc(x int) *_ENode {
	return ett.makeNode(x, nil, nil)
}

func (ett *EulerTourTree) Link(u *_ENode, v *_ENode) *_ENode {
	ett.splay(u)
	ett.splay(v)
	uv := ett.makeNode(0, u, ett.disconnect(v, 1))
	vu := ett.makeNode(0, v, ett.disconnect(u, 1))
	uv.r = vu
	vu.r = uv
	ett.join(uv, vu)
	return uv
}

func (ett *EulerTourTree) Cut(uv *_ENode) {
	ett.splay(uv)
	ett.disconnect(uv, 1)
	ett.splay(uv.r)
	ett.join(ett.disconnect(uv, 0), ett.disconnect(uv.r, 1))
	uv, uv.r = nil, nil // TODO: delete
}

func (ett *EulerTourTree) GetSum(u *_ENode) int {
	ett.splay(u)
	return u.s
}

func (ett *EulerTourTree) sum(t *_ENode) int {
	if t == nil {
		return 0
	}
	return t.s
}

func (ett *EulerTourTree) update(t *_ENode) *_ENode {
	if t != nil {
		t.s = t.x + ett.sum(t.ch[0]) + ett.sum(t.ch[1])
	}
	return t
}

func (ett *EulerTourTree) dir(t *_ENode) int {
	if t != t.p.ch[0] {
		return 1
	}
	return 0
}

func (ett *EulerTourTree) connect(p *_ENode, t *_ENode, d int) {
	p.ch[d] = t
	if t != nil {
		t.p = p
	}
	ett.update(p)
}

func (ett *EulerTourTree) disconnect(t *_ENode, d int) *_ENode {
	c := t.ch[d]
	t.ch[d] = nil
	if c != nil {
		c.p = nil
	}
	ett.update(t)
	return c
}

func (ett *EulerTourTree) rot(t *_ENode) {
	p := t.p
	d := ett.dir(t)
	if p.p != nil {
		ett.connect(p.p, t, ett.dir(p))
	} else {
		t.p = p.p
	}
	ett.connect(p, t.ch[1^d], d)
	ett.connect(t, p, 1^d)
}

func (ett *EulerTourTree) splay(t *_ENode) {
	for ; t.p != nil; ett.rot(t) {
		if t.p.p != nil {
			if ett.dir(t) == ett.dir(t.p) {
				ett.rot(t.p)
			} else {
				ett.rot(t)
			}
		}
	}
}

func (ett *EulerTourTree) join(s *_ENode, t *_ENode) {
	if s == nil || t == nil {
		return
	}
	for ; s.ch[1] != nil; s = s.ch[1] {
	}
	ett.splay(s)
	for ; t.ch[0] != nil; t = t.ch[0] {
	}
	ett.connect(t, s, 0)
}

func (ett *EulerTourTree) makeNode(x int, l *_ENode, r *_ENode) *_ENode {
	t := &_ENode{x: x}
	ett.connect(t, l, 0)
	ett.connect(t, r, 1)
	return t
}
