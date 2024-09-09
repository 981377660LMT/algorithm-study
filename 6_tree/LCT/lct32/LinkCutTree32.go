// !单点修改的LCT
// https://ei1333.github.io/library/structure/lct/link-cut-tree.hpp

// NewLinkCutTree(f, s): コンストラクタ. f は 2 つの要素の値をマージする二項演算, s は要素を反転する演算を指す.
// Alloc(v): 要素の値を v としたノードを生成する.
// Build(vs): 各要素の値を vs[i] としたノードを生成し, その配列を返す.
// Evert(t): t を根に変更する.
// LinkEdge(child, parent): child の親を parent にする.如果已经连通则不进行操作
// CutEdge(u,v) : u と v の間の辺を切り離す.如果边不存在则不进行操作
// QueryToRoot(u): u から根までのパス上の頂点の値を二項演算でまとめた結果を返す.
// QueryPath(u, v): u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
// KthAncestor(x, k): x から根までのパスに出現するノードを並べたとき, 0-indexed で k 番目のノードを返す.
// Dist(u, v): u と v の距離を返す.
// LCA(u, v): u と v の lca を返す. u と v が異なる連結成分なら nullptr を返す.
//  !上記の操作は根を勝手に変えるため、根を固定したい場合は Evert で根を固定してから操作する.
// IsConnected(u, v): u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
// Set(t, v): t の値を v に変更する.
// Get(t): t の値を返す.
// GetRoot(t): t の根を返す.
// expose(t): t と根をつなげて, t を splay Tree の根にする.

package main

import (
	"fmt"
	"runtime/debug"
)

// 单组测试用例时禁用GC
func init() {
	debug.SetGCPercent(-1)
}

func main() {

}

func demo() {
	uf := NewLinkCutTree32(true)
	n := 10
	nodes := uf.Build(int32(n), func(i int32) E { return 0 })
	for i := 0; i < n; i++ {
		uf.Set(nodes[i], 1)
	}
	uf.LinkEdge(nodes[1], nodes[0])
	uf.LinkEdge(nodes[2], nodes[0])
	fmt.Println(uf.GetRoot(nodes[1]) == nodes[0])
	fmt.Println(uf.GetRoot(nodes[2]) == nodes[0])
	fmt.Println(uf.GetRoot(nodes[3]) == nodes[0])
	fmt.Println(uf.QueryPath(nodes[1], nodes[2]))

	T := NewLinkCutTree32(true)
	nodes = T.Build(10, func(i int32) E { return 0 })
	for i := 0; i < 9; i++ {
		T.LinkEdge(nodes[i+1], nodes[i])
	}
	fmt.Println(T.KthAncestor(nodes[9], 1).id)
	fmt.Println(T.GetParent(nodes[9]).id)
	fmt.Println(T.Jump(nodes[9], nodes[5], 4).id)
}

// https://www.luogu.com.cn/problem/P3203
// 每个弹力装置初始系数为nums[i].
// 当绵羊踩到第i个弹力装置时，会被弹到第i+nums[i]个弹力装置上。如果不存在第i+nums[i]个弹力装置，则绵羊会被弹飞。
//
// 1 index : 输出从index出发被弹几次后弹飞
// 2 index newValue : 将index位置的弹力装置系数改为newValue
// n<=2e5 q<=1e5
//
// 方法1：LCT
// 方法2：分块
func 弹飞绵羊(nums []int32, operations [][3]int32) []int32 {
	n := int32(len(nums))
	weights := make([]int, n+1) // n是虚拟结点
	for i := range weights {
		weights[i] = 1
	}
	tree := NewLinkCutTree32(false)
	nodes := tree.Build(int32(len(weights)), func(i int32) E { return E(weights[i]) })
	for i := int32(0); i < n; i++ {
		v := nums[i]
		if i+v < n {
			tree.LinkEdge(nodes[i], nodes[i+v])
		} else {
			tree.LinkEdge(nodes[i], nodes[n])
		}
	}

	res := []int32{}
	for _, op := range operations {
		kind := op[0]
		if kind == 1 {
			// 查询index到虚拟根节点的距离
			cur := op[1]
			dist := tree.QueryPath(nodes[cur], nodes[n]) - 1
			res = append(res, int32(dist))
		} else {
			cur, newValue := op[1], op[2]
			preValue := nums[cur]
			nums[cur] = newValue
			if cur+preValue < n {
				tree.CutEdge(nodes[cur], nodes[cur+preValue])
			} else {
				tree.CutEdge(nodes[cur], nodes[n])
			}
			if cur+newValue < n {
				tree.LinkEdge(nodes[cur], nodes[cur+newValue])
			} else {
				tree.LinkEdge(nodes[cur], nodes[n])
			}
		}
	}

	return res
}

type E = int

func (*LinkCutTree32) rev(e E) E   { return e } // 区间反转
func (*LinkCutTree32) op(a, b E) E { return a + b }

type LinkCutTree32 struct {
	nodeId int32
	edges  map[[2]int32]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTree32(check bool) *LinkCutTree32 {
	return &LinkCutTree32{edges: make(map[[2]int32]struct{}), check: check}
}

func (lct *LinkCutTree32) Build(n int32, f func(i int32) E) []*treeNode {
	nodes := make([]*treeNode, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = lct.Alloc(f(i))
	}
	return nodes
}

func (lct *LinkCutTree32) Alloc(e E) *treeNode {
	res := newTreeNode(e, lct.nodeId)
	lct.nodeId++
	return res
}

// t を根に変更する.
func (lct *LinkCutTree32) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

// 存在していない辺 uv を新たに張る.
//
//	すでに存在している辺 uv に対しては何もしない.
func (lct *LinkCutTree32) LinkEdge(child, parent *treeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		lct.edges[[2]int32{id1, id2}] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

// 存在している辺を切り離す.
//
//	存在していない辺に対しては何もしない.
func (lct *LinkCutTree32) CutEdge(u, v *treeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := [2]int32{id1, id2}
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
//
//	u と v が異なる連結成分なら nullptr を返す.
//	!上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTree32) LCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTree32) KthAncestor(x *treeNode, k int32) *treeNode {
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

func (lct *LinkCutTree32) GetParent(t *treeNode) *treeNode {
	lct.expose(t)
	p := t.l
	if p == nil {
		return nil
	}
	for {
		lct.push(p)
		if p.r == nil {
			return p
		}
		p = p.r
	}
}

func (lct *LinkCutTree32) Jump(from, to *treeNode, k int32) *treeNode {
	lct.Evert(to)
	lct.expose(from)
	for {
		lct.push(from)
		rs := int32(0)
		if from.r != nil {
			rs = from.r.sz
		}
		if k < rs {
			from = from.r
			continue
		}
		if k == rs {
			break
		}
		k -= rs + 1
		from = from.l
	}
	lct.splay(from)
	return from
}

func (lct *LinkCutTree32) Dist(u, v *treeNode) int32 {
	lct.Evert(u)
	lct.expose(v)
	return v.sz - 1
}

// u から根までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTree32) QueryToRoot(u *treeNode) E {
	lct.expose(u)
	return u.sum
}

// u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTree32) QueryPath(u, v *treeNode) E {
	lct.Evert(u)
	return lct.QueryToRoot(v)
}

// t の値を v に変更する.
func (lct *LinkCutTree32) Set(t *treeNode, v E) {
	lct.expose(t)
	t.key = v
	lct.update(t)
}

// t の値を返す.
func (lct *LinkCutTree32) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTree32) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

// t の根を返す.
func (lct *LinkCutTree32) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

func (lct *LinkCutTree32) expose(t *treeNode) *treeNode {
	var rp *treeNode
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		cur.r = rp
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTree32) update(t *treeNode) *treeNode {
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

func (lct *LinkCutTree32) rotr(t *treeNode) {
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

func (lct *LinkCutTree32) rotl(t *treeNode) {
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

func (lct *LinkCutTree32) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.sum = lct.rev(t.sum)
	t.rev = !t.rev
}

func (lct *LinkCutTree32) push(t *treeNode) {
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

func (lct *LinkCutTree32) splay(t *treeNode) {
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
	sz       int32
	id       int32
}

func newTreeNode(v E, id int32) *treeNode {
	return &treeNode{key: v, sum: v, sz: 1, id: id}
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) String() string {
	return fmt.Sprintf("key: %v, sum: %v, sz: %v, rev: %v", n.key, n.sum, n.sz, n.rev)
}
