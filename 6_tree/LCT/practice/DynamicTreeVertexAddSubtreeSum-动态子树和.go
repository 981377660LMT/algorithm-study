// 动态子树和查询
// n,q=2e5
// 事件监听:
// 构造函数传入一个监听器(emitter)在合适的时机调用监听器的回调,就可以拿到左右子树/父亲上的metadata了.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/dynamic_tree_vertex_add_subtree_sum
	// 0 u v w x  删除 u-v 边, 添加 w-x 边
	// 1 p x  将 p 节点的值加上 x
	// 2 u p  对于边(u,p) 查询结点v的子树的和,p为u的父节点

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	lct := NewLinkCutTreeSubTree(&ei1333Listener{})
	vs := lct.Build(nums)

	for i := 0; i < n-1; i++ { // 连接树边
		var v1, v2 int
		fmt.Fscan(in, &v1, &v2)
		lct.LinkEdge(vs[v1], vs[v2])
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var u, v, w, x int
			fmt.Fscan(in, &u, &v, &w, &x)
			lct.CutEdge(vs[u], vs[v])
			lct.LinkEdge(vs[w], vs[x])
		} else if op == 1 {
			var root, delta int
			fmt.Fscan(in, &root, &delta)
			lct.Set(vs[root], vs[root].key+delta)
		} else {
			var root, parent int
			fmt.Fscan(in, &root, &parent)
			lct.Evert(vs[parent])
			fmt.Fprintln(out, lct.QuerySubTree(vs[root]).csum)
		}
	}
}

type E = int

type Metadata struct {
	csum int // 子树和
	sum  int // 整棵树的和
	psum int // 父亲子树和
	lsum int // 左儿子子树和
}

type ei1333Listener struct{}

func (*ei1333Listener) OnToggle(cur *Metadata) { cur.psum, cur.csum = cur.csum, cur.psum }

func (*ei1333Listener) OnMerge(cost E, cur, parent, child *Metadata) {
	cur.sum = parent.sum + child.sum + cost + cur.lsum
	cur.psum = parent.psum + cost + cur.lsum
	cur.csum = child.csum + cost + cur.lsum
}

func (*ei1333Listener) OnAdd(cur, child *Metadata) { cur.lsum += child.sum }

func (*ei1333Listener) OnErase(cur, child *Metadata) { cur.lsum -= child.sum }

type nodeListener interface {
	OnToggle(cur *Metadata)
	OnMerge(cost E, cur, parent, child *Metadata)
	OnAdd(cur, child *Metadata)
	OnErase(cur, child *Metadata)
}

type LinkCutTreeSubTree struct {
	sm       Metadata // !ident
	listener nodeListener
	nodeId   int
	edges    map[struct{ u, v int }]struct{}
}

func NewLinkCutTreeSubTree(listener nodeListener) *LinkCutTreeSubTree {
	return &LinkCutTreeSubTree{listener: listener, edges: make(map[struct{ u, v int }]struct{})}
}

// 各要素の値を vs[i] としたノードを生成し, その配列を返す.
func (lct *LinkCutTreeSubTree) Build(vs []E) []*treeNode {
	nodes := make([]*treeNode, len(vs))
	for i, v := range vs {
		nodes[i] = lct.Alloc(v)
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeSubTree) Alloc(key E) *treeNode {
	res := newTreeNode(key, lct.sm, lct.nodeId)
	lct.nodeId++
	lct.update(res)
	return res
}

// t を根に変更する.
func (lct *LinkCutTreeSubTree) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

func (lct *LinkCutTreeSubTree) LinkEdge(child, parent *treeNode) (ok bool) {
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

func (lct *LinkCutTreeSubTree) CutEdge(u, v *treeNode) (ok bool) {
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
func (lct *LinkCutTreeSubTree) QueryLCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeSubTree) QueryKthAncestor(x *treeNode, k int) *treeNode {
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

func (lct *LinkCutTreeSubTree) QuerySubTree(t *treeNode) Metadata {
	lct.expose(t)
	return t.md
}

// t の値を v に変更する.
func (lct *LinkCutTreeSubTree) Set(t *treeNode, v E) *treeNode {
	lct.expose(t)
	t.key = v
	lct.update(t)
	return t
}

// t の値を返す.
func (lct *LinkCutTreeSubTree) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeSubTree) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeSubTree) expose(t *treeNode) *treeNode {
	rp := (*treeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		if cur.r != nil {
			lct.listener.OnAdd(&cur.md, &cur.r.md)
		}
		cur.r = rp
		if cur.r != nil {
			lct.listener.OnErase(&cur.md, &cur.r.md)
		}
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeSubTree) update(t *treeNode) {
	t.sz = 1
	if t.l != nil {
		t.sz += t.l.sz
	}
	if t.r != nil {
		t.sz += t.r.sz
	}

	tmp1, tmp2 := &lct.sm, &lct.sm
	if t.l != nil {
		tmp1 = &t.l.md
	}
	if t.r != nil {
		tmp2 = &t.r.md
	}

	lct.listener.OnMerge(t.key, &t.md, tmp1, tmp2)
}

func (lct *LinkCutTreeSubTree) rotr(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) rotl(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	lct.listener.OnToggle(&t.md)
	t.rev = !t.rev
}

func (lct *LinkCutTreeSubTree) push(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) splay(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

type treeNode struct {
	key     E
	md      Metadata
	l, r, p *treeNode
	rev     bool
	sz      int
	id      int
}

func newTreeNode(key E, md Metadata, id int) *treeNode {
	return &treeNode{key: key, md: md, sz: 1, id: id}
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) String() string {
	return fmt.Sprintf("key: %v, sum: %v, sz: %v, rev: %v", n.key, n.md, n.sz, n.rev)
}
