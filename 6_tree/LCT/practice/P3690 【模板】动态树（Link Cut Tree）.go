// https://www.luogu.com.cn/problem/P3690
// 给定n个点以及每个点的权值，要你处理接下来的m个操作。
// 操作有四种，操作从0到3编号。点从1到n编号。
// 0 x y 代表询问从x到y的路径上的点的权值的xor和。保证x到y是联通的。
// 1 x y 代表连接x到y，若x到y已经联通则无需连接。
// 2 x y 代表删除边(x, y)，不保证边(a, y)存在。
// 3 x y 代表将点x上的权值变成y。

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
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	lct := NewLinkCutTree()
	vs := lct.Build(nums)

	for i := 0; i < q; i++ {
		var op, x, y int
		fmt.Fscan(in, &op, &x, &y)
		x, y = x-1, y-1
		if op == 0 {
			fmt.Fprintln(out, lct.QueryPath(vs[x], vs[y]))
		} else if op == 1 {
			lct.LinkEdge(vs[x], vs[y])
		} else if op == 2 {
			lct.CutEdge(vs[x], vs[y])
		} else {
			lct.Set(vs[x], y+1)
		}
	}
}

type E = int

func (*LinkCutTree) rev(e E) E   { return e } // 区间反转
func (*LinkCutTree) op(a, b E) E { return a ^ b }

type LinkCutTree struct {
	nodeId int
	edges  map[struct{ u, v int }]struct{}
}

func NewLinkCutTree() *LinkCutTree {
	return &LinkCutTree{edges: make(map[struct{ u, v int }]struct{})}
}

// 各要素の値を vs[i] としたノードを生成し, その配列を返す.
func (lct *LinkCutTree) Build(vs []E) []*treeNode {
	nodes := make([]*treeNode, len(vs))
	for i, v := range vs {
		nodes[i] = lct.Alloc(v)
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTree) Alloc(e E) *treeNode {
	res := newTreeNode(e, lct.nodeId)
	lct.nodeId++
	return res
}

// t を根に変更する.
func (lct *LinkCutTree) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

// 存在していない辺 uv を新たに張る.
//  すでに存在している辺 uv に対しては何もしない.
func (lct *LinkCutTree) LinkEdge(child, parent *treeNode) (ok bool) {
	if lct.IsConnected(child, parent) {
		return
	}
	id1, id2 := child.id, parent.id
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	tuple := struct{ u, v int }{id1, id2}
	lct.edges[tuple] = struct{}{}
	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

// 存在している辺を切り離す.
//  存在していない辺に対しては何もしない.
func (lct *LinkCutTree) CutEdge(u, v *treeNode) (ok bool) {
	id1, id2 := u.id, v.id
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	tuple := struct{ u, v int }{id1, id2}
	if _, has := lct.edges[tuple]; !has {
		return
	}
	delete(lct.edges, tuple)
	lct.Evert(u)
	lct.expose(v)
	parent := v.l
	v.l = nil
	lct.update(v)
	parent.p = nil
	return true
}

// u と v の lca を返す.
//  u と v が異なる連結成分なら nullptr を返す.
//  !上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTree) QueryLCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTree) QueryKthAncestor(x *treeNode, k int) *treeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.sz > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.sz
			}
			if k == 0 {
				return x
			}
			k--
			x = x.l
		}
	}
	return nil
}

// u から根までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTree) QueryToRoot(u *treeNode) E {
	lct.expose(u)
	return u.sum
}

// u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTree) QueryPath(u, v *treeNode) E {
	lct.Evert(u)
	return lct.QueryToRoot(v)
}

// t の値を v に変更する.
func (lct *LinkCutTree) Set(t *treeNode, v E) {
	lct.expose(t)
	t.key = v
	lct.update(t)
}

// t の値を返す.
func (lct *LinkCutTree) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTree) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

// t の根を返す.
func (lct *LinkCutTree) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

func (lct *LinkCutTree) expose(t *treeNode) *treeNode {
	rp := (*treeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		cur.r = rp
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTree) update(t *treeNode) *treeNode {
	t.sz = 1
	t.sum = t.key
	if t.l != nil {
		t.sz += t.l.sz
		t.sum = lct.op(t.l.sum, t.sum)
	}
	if t.r != nil {
		t.sz += t.r.sz
		t.sum = lct.op(t.sum, t.r.sum)
	}
	return t
}

func (lct *LinkCutTree) rotr(t *treeNode) {
	x := t.p
	y := x.p
	x.l = t.r
	if t.r != nil {
		t.r.p = x
	}
	t.r = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTree) rotl(t *treeNode) {
	x := t.p
	y := x.p
	x.r = t.l
	if t.l != nil {
		t.l.p = x
	}
	t.l = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTree) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.sum = lct.rev(t.sum)
	t.rev = !t.rev
}

func (lct *LinkCutTree) push(t *treeNode) {
	if t.rev {
		if t.l != nil {
			lct.toggle(t.l)
		}
		if t.r != nil {
			lct.toggle(t.r)
		}
		t.rev = false
	}
}

func (lct *LinkCutTree) splay(t *treeNode) {
	lct.push(t)
	for !t.IsRoot() {
		q := t.p
		if q.IsRoot() {
			lct.push(q)
			lct.push(t)
			if q.l == t {
				lct.rotr(t)
			} else {
				lct.rotl(t)
			}
		} else {
			r := q.p
			lct.push(r)
			lct.push(q)
			lct.push(t)
			if r.l == q {
				if q.l == t {
					lct.rotr(q)
					lct.rotr(t)
				} else {
					lct.rotl(t)
					lct.rotr(t)
				}
			} else {
				if q.r == t {
					lct.rotl(q)
					lct.rotl(t)
				} else {
					lct.rotr(t)
					lct.rotl(t)
				}
			}
		}
	}
}

type treeNode struct {
	l, r, p  *treeNode
	key, sum E
	rev      bool
	sz       int
	id       int
}

func newTreeNode(v E, id int) *treeNode {
	return &treeNode{key: v, sum: v, sz: 1, id: id}
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) String() string {
	return fmt.Sprintf("key: %v, sum: %v, sz: %v, rev: %v", n.key, n.sum, n.sz, n.rev)
}
