// 树上最大子数组和
// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2450
// 给定一颗n个顶点的树和q个操作,每个顶点有一个权值wi
// 操作有两种类型
// 1 ai bi ci :将ai到bi的最短路径上的所有点(包含端点)的权值改为ci
// 2 ai bi ci :将ai到bi的最短路径上的所有点(包含端点)的权值升序排序后的数组记为A.
//             求出A中非空子数组的最大和. (ci没用)

// n,q<=2e5,-1e4<=wi<=1e4

// 非常极限 禁用gc后2s超时

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

const INF int = 1e18

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	leaves := make([]Node, n)
	for i := 0; i < n; i++ {
		var num int
		fmt.Fscan(in, &num)
		leaves[i] = makeNode(1, num)
	}

	lct := NewLinkCutTreeLazy()
	vs := lct.Build(leaves)

	for i := 0; i < n-1; i++ { // 连接树边
		var v1, v2 int
		fmt.Fscan(in, &v1, &v2)
		v1, v2 = v1-1, v2-1
		lct.LinkEdge(vs[v1], vs[v2])
	}

	for i := 0; i < q; i++ {
		var op, a, b, v int
		fmt.Fscan(in, &op, &a, &b, &v)
		a--
		b--
		if op == 1 {
			lct.UpdatePath(vs[a], vs[b], v)
		} else {
			fmt.Fprintln(out, lct.QueryPath(vs[a], vs[b]).res)
		}
	}
}

type Node struct {
	res, length int // 最大字段和,区间长度
	all         int // 区间和
	left, right int // 前缀/后缀最大和
}

func makeNode(len, lazy int) Node {
	node := Node{}
	node.length = len
	node.all = lazy * len
	// 大于0就取连续子数组,否则取1个
	cur := lazy
	if lazy > 0 {
		cur = node.all
	}
	node.left = cur
	node.right = cur
	node.res = cur
	return node
}

func merge(a, b Node) Node {
	node := Node{}
	node.length = a.length + b.length
	node.all = a.all + b.all
	node.left = max(a.left, a.all+b.left)
	node.right = max(b.right, b.all+a.right)
	node.res = max(max(a.res, b.res), a.right+b.left)
	return node
}

type E = Node
type Id = int

// 区间反转
func (*LinkCutTreeLazy) rev(e E) E {
	e.left, e.right = e.right, e.left
	return e
}
func (*LinkCutTreeLazy) id() Id                                  { return -INF }
func (*LinkCutTreeLazy) op(a, b E) E                             { return merge(a, b) }
func (*LinkCutTreeLazy) mapping(lazy Id, data E) E               { return makeNode(data.length, lazy) }
func (*LinkCutTreeLazy) composition(parentLazy, childLazy Id) Id { return parentLazy }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//
//
//
//
type LinkCutTreeLazy struct{}

func NewLinkCutTreeLazy() *LinkCutTreeLazy {
	return &LinkCutTreeLazy{}
}

// 各要素の値を vs[i] としたノードを生成し, その配列を返す.
func (lct *LinkCutTreeLazy) Build(vs []E) []*treeNode {
	nodes := make([]*treeNode, len(vs))
	for i, v := range vs {
		nodes[i] = lct.Alloc(v)
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeLazy) Alloc(e E) *treeNode {
	return newTreeNode(e, lct.id())
}

// t を根に変更する.
func (lct *LinkCutTreeLazy) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

// child の親を parent にする.
//  child と parent は別の連結成分で, `child が根であること`を要求する.
func (lct *LinkCutTreeLazy) Link(child, parent *treeNode) (ok bool) {
	// child and parent must belong to different connected components.
	if lct.IsConnected(child, parent) {
		return
	}
	// child must be root.
	if child.l != nil {
		return
	}
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

// 存在していない辺 uv を新たに張る.
//  すでに存在している辺 uv に対しては何もしない.
func (lct *LinkCutTreeLazy) LinkEdge(u, v *treeNode) {
	lct.Evert(u)
	lct.expose(v)
	u.p = v
	v.r = u
	lct.update(v)
	return
}

func (lct *LinkCutTreeLazy) LinkEdgeIfAbsent(u, v *treeNode) (ok bool) {
	if lct.IsConnected(u, v) {
		return
	}
	lct.Evert(u)
	lct.expose(v)
	u.p = v
	v.r = u
	lct.update(v)
	return true
}

// child の親と child を切り離す.
func (lct *LinkCutTreeLazy) Cut(child *treeNode) (ok bool) {
	lct.expose(child)
	parent := child.l
	// child must not be root.
	if parent == nil {
		return
	}
	child.l = nil
	parent.p = nil
	lct.update(child)
	return true
}

// 存在している辺を切り離す.
//  存在していない辺に対しては何もしない.
func (lct *LinkCutTreeLazy) CutEdge(u, v *treeNode) (ok bool) {
	lct.Evert(u)
	return lct.Cut(v)
}

// u と v の lca を返す.
//  u と v が異なる連結成分なら nullptr を返す.
//  !上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTreeLazy) QueryLCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeLazy) QueryKthAncestor(x *treeNode, k int) *treeNode {
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
func (lct *LinkCutTreeLazy) QueryToRoot(u *treeNode) E {
	lct.expose(u)
	return u.sum
}

// u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTreeLazy) QueryPath(u, v *treeNode) E {
	lct.Evert(u)
	return lct.QueryToRoot(v)
}

func (lct *LinkCutTreeLazy) UpdateToRoot(t *treeNode, lazy Id) {
	lct.expose(t)
	lct.propagate(t, lazy)
	lct.push(t)
}

func (lct *LinkCutTreeLazy) UpdatePath(u, v *treeNode, lazy Id) {
	lct.Evert(u)
	lct.UpdateToRoot(v, lazy)
}

// t の値を v に変更する.
func (lct *LinkCutTreeLazy) Set(t *treeNode, v E) {
	lct.expose(t)
	t.key = v
	lct.update(t)
}

// t の値を返す.
func (lct *LinkCutTreeLazy) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeLazy) IsConnected(u, v *treeNode) bool {
	lct.expose(u)
	lct.expose(v)
	return u == v || u.p != nil
}

func (lct *LinkCutTreeLazy) expose(t *treeNode) *treeNode {
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

func (lct *LinkCutTreeLazy) update(t *treeNode) *treeNode {
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

func (lct *LinkCutTreeLazy) rotr(t *treeNode) {
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

func (lct *LinkCutTreeLazy) rotl(t *treeNode) {
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

func (lct *LinkCutTreeLazy) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.sum = lct.rev(t.sum)
	t.rev = !t.rev
}

func (lct *LinkCutTreeLazy) propagate(t *treeNode, lazy Id) {
	t.lazy = lct.composition(lazy, t.lazy)
	t.key = lct.mapping(lazy, t.key)
	t.sum = lct.mapping(lazy, t.sum)
}

func (lct *LinkCutTreeLazy) push(t *treeNode) {
	if t.lazy != lct.id() {
		if t.l != nil {
			lct.propagate(t.l, t.lazy)
		}
		if t.r != nil {
			lct.propagate(t.r, t.lazy)
		}
		t.lazy = lct.id()
	}

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

func (lct *LinkCutTreeLazy) splay(t *treeNode) {
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
	lazy     Id
	rev      bool
	sz       int
}

func newTreeNode(v E, lazy Id) *treeNode {
	return &treeNode{key: v, sum: v, lazy: lazy, sz: 1}
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) String() string {
	return fmt.Sprintf("key: %v, sum: %v, sz: %v, rev: %v", n.key, n.sum, n.sz, n.rev)
}
