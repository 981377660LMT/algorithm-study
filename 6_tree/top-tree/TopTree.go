// https://nyaannyaan.github.io/library/tree/dynamic-rerooting.hpp

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
	yuki768()
}

// No.768 Tapris and Noel play the game on Treeone
// https://yukicoder.me/problems/no/768
// Alice和Bob在一棵树上玩游戏.
// Alice先手，可以选择树上的一个点，然后Bob也可以选择相邻的一个未访问过的点...
// 两人轮流选择，直到没有可选的点则游戏结束.
// 问Alice是否有必胜策略，输出Alice的必胜策略的起始点.
func yuki768() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	D := NewDynamicRerooting(n, make([]Info, n))
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		D.AddEdge(u, v)
	}

	var res []int32
	for i := int32(0); i < n; i++ {
		q := D.QueryAll(i)
		var u int32
		if q.v != 0 {
			u = q.u1
		} else {
			u = q.u0
		}
		if u == 0 {
			res = append(res, i+1)
		}
	}

	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Path = struct{ v, u0, u1 int32 }
type Point = int32
type Info = int32

func vertex(v Info) Path { return Path{0, 0, 1} }
func addEdge(p Path) Point {
	if p.v != 0 {
		return p.u1 ^ 1
	}
	return p.u0 ^ 1
}
func addVertex(x Point, v Info) Path { return Path{x, 0, 1} }
func rake(x, y Point) Point          { return x | y }
func compress(p, c Path) Path {
	res := Path{}
	res.v = c.v
	w0 := p.v | (c.u0 ^ 1)
	if w0 != 0 {
		res.u0 = p.u1
	} else {
		res.u0 = p.u0
	}
	w1 := p.v | (c.u1 ^ 1)
	if w1 != 0 {
		res.u1 = p.u1
	} else {
		res.u1 = p.u0
	}
	return res
}

type DynamicRerooting struct {
	n       int32
	topTree *TopTree
	nodes   []*tNode
}

func NewDynamicRerooting(n int32, info []Info) *DynamicRerooting {
	res := &DynamicRerooting{n: n, topTree: NewTopTree()}
	nodes := make([]*tNode, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = res.topTree.Alloc(info[i])
	}
	res.nodes = nodes
	return res
}

func (dr *DynamicRerooting) AddEdge(u, v int32) {
	dr.topTree.Evert(dr.nodes[u])
	dr.topTree.Link(dr.nodes[u], dr.nodes[v])
}

func (dr *DynamicRerooting) RemoveEdge(u, v int32) {
	dr.topTree.Evert(dr.nodes[u])
	dr.topTree.Cut(dr.nodes[v])
}

func (dr *DynamicRerooting) GetInfo(u int32) Info {
	return dr.nodes[u].info
}

func (dr *DynamicRerooting) SetInfo(u int32, info Info) {
	dr.topTree.SetKey(dr.nodes[u], info)
}

func (dr *DynamicRerooting) QueryAll(root int32) Path {
	return dr.topTree.QueryAll(dr.nodes[root])
}

func (dr *DynamicRerooting) QuerySubtree(root, u int32) Path {
	return dr.topTree.QuerySubtree(dr.nodes[root], dr.nodes[u])
}

type tNode struct {
	l, r, p          *tNode
	info             Info
	key, sum, revSum Path
	light, belong    *sNode
	rev              bool
}

func newTNode(info Info) *tNode {
	return &tNode{info: info}
}

func (tn *tNode) IsRoot() bool {
	return tn.p == nil || (tn.p.l != tn && tn.p.r != tn)
}

type TopTree struct {
	splay *splayTreeforDashedEdge
}

func NewTopTree() *TopTree { return &TopTree{splay: newSplayTreeforDashedEdge()} }

func (tt *TopTree) Link(child, parent *tNode) {
	tt._expose(parent)
	tt._expose(child)
	child.p = parent
	parent.r = child
	tt._update(parent)
}

func (tt *TopTree) Cut(child *tNode) {
	tt._expose(child)
	parent := child.l
	child.l = nil
	parent.p = nil
	tt._update(child)
}

func (tt *TopTree) Evert(t *tNode) {
	tt._expose(t)
	tt._toggle(t)
	tt._push(t)
}

func (tt *TopTree) Alloc(info Info) *tNode {
	t := newTNode(info)
	tt._update(t)
	return t
}

func (tt *TopTree) IsConnected(u, v *tNode) bool {
	tt._expose(u)
	tt._expose(v)
	return u == v || u.p != nil
}

func (tt *TopTree) Lca(u, v *tNode) *tNode {
	if !tt.IsConnected(u, v) {
		return nil
	}
	tt._expose(u)
	return tt._expose(v)
}

func (tt *TopTree) SetKey(t *tNode, info Info) {
	tt._expose(t)
	t.info = info
	tt._update(t)
}

func (tt *TopTree) QueryAll(u *tNode) Path {
	tt.Evert(u)
	return u.sum
}

// TODO: check
func (tt *TopTree) QueryPath(u, v *tNode) Path {
	tt.Evert(u)
	tt._expose(v)
	return v.sum
}

// root为根时,子树u的和.
func (tt *TopTree) QuerySubtree(root, u *tNode) Path {
	tt.Evert(root)
	tt._expose(u)
	l := u.l
	u.l = nil
	tt._update(u)
	res := u.sum
	u.l = l
	tt._update(u)
	return res
}

func (tt *TopTree) _toggle(t *tNode) {
	t.l, t.r = t.r, t.l
	t.sum, t.revSum = t.revSum, t.sum
	t.rev = !t.rev
}

func (tt *TopTree) _rotr(t *tNode) {
	x := t.p
	y := x.p
	tt._push(x)
	tt._push(t)
	x.l = t.r
	if x.l != nil {
		x.l.p = x
	}
	t.r = x
	x.p = t
	tt._update(x)
	tt._update(t)
	t.p = y
	if t.p != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
	}
}

func (tt *TopTree) _rotl(t *tNode) {
	x := t.p
	y := x.p
	tt._push(x)
	tt._push(t)
	x.r = t.l
	if x.r != nil {
		x.r.p = x
	}
	t.l = x
	x.p = t
	tt._update(x)
	tt._update(t)
	t.p = y
	if t.p != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
	}
}

func (tt *TopTree) _push(t *tNode) {
	if t.rev {
		if t.l != nil {
			tt._toggle(t.l)
		}
		if t.r != nil {
			tt._toggle(t.r)
		}
		t.rev = false
	}
}

func (tt *TopTree) _pushRev(t *tNode) {
	if t.rev {
		if t.l != nil {
			tt._toggle(t.l)
		}
		if t.r != nil {
			tt._toggle(t.r)
		}
		t.rev = false
	}
}

func (tt *TopTree) _update(t *tNode) {
	var key Path
	if t.light != nil {
		key = addVertex(t.light.sum, t.info)
	} else {
		key = vertex(t.info)
	}
	sum, revSum := key, key
	if t.l != nil {
		sum = compress(t.l.sum, sum)
		revSum = compress(revSum, t.l.revSum)
	}
	if t.r != nil {
		sum = compress(sum, t.r.sum)
		revSum = compress(t.r.revSum, revSum)
	}
	t.key, t.sum, t.revSum = key, sum, revSum
}

func (tt *TopTree) _splay(t *tNode) {
	tt._push(t)
	{
		rot := t
		for !rot.IsRoot() {
			rot = rot.p
		}
		t.belong = rot.belong
		if t != rot {
			rot.belong = nil
		}
	}
	for !t.IsRoot() {
		q := t.p
		if q.IsRoot() {
			tt._pushRev(q)
			tt._pushRev(t)
			if q.l == t {
				tt._rotr(t)
			} else {
				tt._rotl(t)
			}
		} else {
			r := q.p
			tt._pushRev(r)
			tt._pushRev(q)
			tt._pushRev(t)
			if r.l == q {
				if q.l == t {
					tt._rotr(q)
					tt._rotr(t)
				} else {
					tt._rotl(t)
					tt._rotr(t)
				}
			} else {
				if q.r == t {
					tt._rotl(q)
					tt._rotl(t)
				} else {
					tt._rotr(t)
					tt._rotl(t)
				}
			}
		}
	}
}

func (tt *TopTree) _expose(t *tNode) *tNode {
	var rp *tNode
	for cur := t; cur != nil; cur = cur.p {
		tt._splay(cur)
		if cur.r != nil {
			cur.light = tt.splay._insert(cur.light, addEdge(cur.r.sum))
			cur.r.belong = cur.light
		}
		cur.r = rp
		if cur.r != nil {
			tt.splay._splay(cur.r.belong)
			tt._push(cur.r)
			cur.light = tt.splay._erase(cur.r.belong)
		}
		tt._update(cur)
		rp = cur
	}
	tt._splay(t)
	return rp
}

type sNode struct {
	l, r, p  *sNode
	key, sum Point
}

func newSNode(key Point) *sNode {
	return &sNode{key: key, sum: key}
}

type splayTreeforDashedEdge struct{}

func newSplayTreeforDashedEdge() *splayTreeforDashedEdge {
	return &splayTreeforDashedEdge{}
}

func (st *splayTreeforDashedEdge) _rotr(t *sNode) {
	x := t.p
	y := x.p
	x.l = t.r
	if x.l != nil {
		x.l.p = x
	}
	t.r = x
	x.p = t
	st._update(x)
	st._update(t)
	t.p = y
	if t.p != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
	}
}

func (st *splayTreeforDashedEdge) _rotl(t *sNode) {
	x := t.p
	y := x.p
	x.r = t.l
	if x.r != nil {
		x.r.p = x
	}
	t.l = x
	x.p = t
	st._update(x)
	st._update(t)
	t.p = y
	if t.p != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
	}
}

func (st *splayTreeforDashedEdge) _update(t *sNode) {
	t.sum = t.key
	if t.l != nil {
		t.sum = rake(t.sum, t.l.sum)
	}
	if t.r != nil {
		t.sum = rake(t.sum, t.r.sum)
	}
}

func (st *splayTreeforDashedEdge) _getRight(t *sNode) *sNode {
	for t.r != nil {
		t = t.r
	}
	return t
}

func (st *splayTreeforDashedEdge) _alloc(v Point) *sNode {
	t := newSNode(v)
	st._update(t)
	return t
}

func (st *splayTreeforDashedEdge) _splay(t *sNode) {
	for t.p != nil {
		q := t.p
		if q.p == nil {
			if q.l == t {
				st._rotr(t)
			} else {
				st._rotl(t)
			}
		} else {
			r := q.p
			if r.l == q {
				if q.l == t {
					st._rotr(q)
					st._rotr(t)
				} else {
					st._rotl(t)
					st._rotr(t)
				}
			} else {
				if q.r == t {
					st._rotl(q)
					st._rotl(t)
				} else {
					st._rotr(t)
					st._rotl(t)
				}
			}
		}
	}
}

func (st *splayTreeforDashedEdge) _insert(t *sNode, v Point) *sNode {
	if t == nil {
		t = st._alloc(v)
		return t
	}
	cur := st._getRight(t)
	z := st._alloc(v)
	st._splay(cur)
	z.p = cur
	cur.r = z
	st._update(cur)
	st._splay(z)
	return z
}

func (st *splayTreeforDashedEdge) _erase(t *sNode) *sNode {
	st._splay(t)
	x := t.l
	y := t.r
	if x == nil {
		t = y
		if t != nil {
			t.p = nil
		}
	} else if y == nil {
		t = x
		t.p = nil
	} else {
		x.p = nil
		t = st._getRight(x)
		st._splay(t)
		t.r = y
		y.p = t
		st._update(t)
	}
	return t
}
