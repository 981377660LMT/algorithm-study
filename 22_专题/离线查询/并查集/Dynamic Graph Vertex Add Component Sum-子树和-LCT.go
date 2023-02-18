// https://judge.yosupo.jp/problem/dynamic_tree_vertex_add_subtree_sum
// Dynamic Graph Vertex Add Component Sum
// 连接断开边/单点修改权值/查询子树和

// 0 root1 root2 root3 root4 断开(root1-root2) 连接(root3-root4)
// 1 root x 将root的值加上x
// 2 root1 root2 输出root1所在子树的和,其中root2是root1的父亲节点

// 在线查询子树和
// n<=2e5

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
	sums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &sums[i])
	}

	lct := NewLinkCutTreeSubTree(false)
	nodes := lct.Build(sums)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		lct.LinkEdge(nodes[u], nodes[v])
	}

	for i := 0; i < q; i++ {
		var op, root1, root2, root3, root4, add, parent int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &root1, &root2, &root3, &root4)
			lct.CutEdge(nodes[root1], nodes[root2])
			lct.LinkEdge(nodes[root3], nodes[root4])
		} else if op == 1 {
			fmt.Fscan(in, &root1, &add)
			lct.Set(nodes[root1], lct.Get(nodes[root1])+add)
		} else {
			fmt.Fscan(in, &root1, &parent)
			lct.Evert(nodes[parent])
			fmt.Fprintln(out, lct.QuerySubTree(nodes[root1]))
		}
	}
}

type E = int // 子树和

func (*TreeNode) e() E                { return 0 }
func (*TreeNode) op(this, other E) E  { return this + other }
func (*TreeNode) inv(this, other E) E { return this - other }

type LinkCutTreeSubTree struct {
	nodeId int
	edges  map[struct{ u, v int }]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTreeSubTree(check bool) *LinkCutTreeSubTree {
	return &LinkCutTreeSubTree{edges: make(map[struct{ u, v int }]struct{}), check: check}
}

// 各要素の値を vs[i] としたノードを生成し, その配列を返す.
func (lct *LinkCutTreeSubTree) Build(vs []E) []*TreeNode {
	nodes := make([]*TreeNode, len(vs))
	for i, v := range vs {
		nodes[i] = lct.Alloc(v)
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeSubTree) Alloc(key E) *TreeNode {
	res := newTreeNode(key, lct.nodeId)
	lct.nodeId++
	lct.update(res)
	return res
}

// t を根に変更する.
func (lct *LinkCutTreeSubTree) Evert(t *TreeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

func (lct *LinkCutTreeSubTree) LinkEdge(child, parent *TreeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := struct{ u, v int }{id1, id2}
		lct.edges[tuple] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

func (lct *LinkCutTreeSubTree) CutEdge(u, v *TreeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := struct{ u, v int }{id1, id2}
		if _, has := lct.edges[tuple]; !has {
			return
		}
		delete(lct.edges, tuple)
	}

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
func (lct *LinkCutTreeSubTree) QueryLCA(u, v *TreeNode) *TreeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeSubTree) QueryKthAncestor(x *TreeNode, k int) *TreeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.cnt > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.cnt
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

// t を根とする部分木の要素の値の和を返す.
//  !Evert を忘れない！
func (lct *LinkCutTreeSubTree) QuerySubTree(t *TreeNode) E {
	lct.expose(t)
	return t.op(t.key, t.sub)
}

// t の値を v に変更する.
func (lct *LinkCutTreeSubTree) Set(t *TreeNode, key E) *TreeNode {
	lct.expose(t)
	t.key = key
	lct.update(t)
	return t
}

// t の値を返す.
func (lct *LinkCutTreeSubTree) Get(t *TreeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeSubTree) IsConnected(u, v *TreeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeSubTree) expose(t *TreeNode) *TreeNode {
	rp := (*TreeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		if cur.r != nil {
			cur.Add(cur.r)
		}
		cur.r = rp
		if cur.r != nil {
			cur.Erase(cur.r)
		}
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeSubTree) update(t *TreeNode) {
	t.cnt = 1
	if t.l != nil {
		t.cnt += t.l.cnt
	}
	if t.r != nil {
		t.cnt += t.r.cnt
	}

	t.Merge(t.l, t.r)
}

func (lct *LinkCutTreeSubTree) rotr(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) rotl(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) toggle(t *TreeNode) {
	t.l, t.r = t.r, t.l
	t.rev = !t.rev
}

func (lct *LinkCutTreeSubTree) push(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) splay(t *TreeNode) {
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

func (lct *LinkCutTreeSubTree) GetRoot(t *TreeNode) *TreeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

type TreeNode struct {
	key, sum, sub E
	rev           bool
	cnt           int
	id            int
	l, r, p       *TreeNode
}

func newTreeNode(key E, id int) *TreeNode {
	res := &TreeNode{key: key, sum: key, cnt: 1, id: id}
	res.sub = res.e()
	return res
}

func (n *TreeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *TreeNode) Add(other *TreeNode)   { n.sub = n.op(n.sub, other.sum) }
func (n *TreeNode) Erase(other *TreeNode) { n.sub = n.inv(n.sub, other.sum) }
func (n *TreeNode) Merge(n1, n2 *TreeNode) {
	var tmp1, tmp2 E
	if n1 != nil {
		tmp1 = n1.sum
	} else {
		tmp1 = n.e()
	}

	if n2 != nil {
		tmp2 = n2.sum
	} else {
		tmp2 = n.e()
	}

	n.sum = n.op(n.op(tmp1, n.key), n.op(n.sub, tmp2))
}
