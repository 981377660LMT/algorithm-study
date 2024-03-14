// 动态子树加子树和查询
// n,q=2e5

// api:
//  Alloc
//  !Evert: 将node作为根节点.
//  Link
//  Cut
//  _pushUp
//  Toggle
//  Get
//  Set
//  SubtreeAdd
//  SubtreeSum
//  GetRoot
//  IsRoot
//  Lca
//  KthAncestor
//  GetParent

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1891F()
	// DynamicTreeVertexAddSubtreeSum()
	// DynamicTreeSubtreeAddSubtreeSum()
}

// A Growing Tree
// https://www.luogu.com.cn/problem/CF1891F
// 给定一棵树，初始只有一个权值为0的结点.
// q次操作:
// 1 x: 添加一个新结点，父亲为x.
// 2 x delta: 将x的子树中所有结点的权值加上delta.
// 操作完后，输出每个结点的权值.
func CF1891F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var q int32
		fmt.Fscan(in, &q)
		lct := NewLinkCutTreeSubTreeAdd()
		nodes := make([]*Node, 0, q+1)
		nodes = append(nodes, lct.Alloc(0))

		for i := int32(0); i < q; i++ {
			var kind int32
			fmt.Fscan(in, &kind)
			if kind == 1 {
				var x int32
				fmt.Fscan(in, &x)
				x--
				nodes = append(nodes, lct.Alloc(0))
				lct.Link(nodes[x], nodes[len(nodes)-1])
			} else {
				var x, delta int32
				fmt.Fscan(in, &x, &delta)
				x--
				lct.Evert(nodes[0])
				lct.SubtreeAdd(nodes[x], int(delta))
			}
		}

		for i := 0; i < len(nodes); i++ {
			fmt.Fprint(out, lct.Get(nodes[i]), " ")
		}
		fmt.Fprintln(out)
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		solve()
	}
}

// 动态单点加子树和
// Dynamic Tree Vertex Add Subtree Sum
// https://judge.yosupo.jp/problem/dynamic_tree_vertex_add_subtree_sum
// 0 u v w x  删除 u-v 边, 添加 w-x 边
// 1 p x  将 p 节点的值加上 x
// 2 u p  对于边(u,p) 查询结点v的子树的和,p为u的父节点
func DynamicTreeVertexAddSubtreeSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	weights := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}
	edges := make([][2]int32, n-1)
	for i := 0; i < n-1; i++ {
		var v1, v2 int32
		fmt.Fscan(in, &v1, &v2)
		edges[i] = [2]int32{v1, v2}
	}
	operations := make([][5]int32, q)
	for i := 0; i < q; i++ {
		var kind int32
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var u, v, w, x int32
			fmt.Fscan(in, &u, &v, &w, &x)
			operations[i] = [5]int32{kind, u, v, w, x}
		} else if kind == 1 {
			var p, x int32
			fmt.Fscan(in, &p, &x)
			operations[i] = [5]int32{kind, p, x, 0, 0}
		} else {
			var u, p int32
			fmt.Fscan(in, &u, &p)
			operations[i] = [5]int32{kind, u, p, 0, 0}
		}
	}

	lct := NewLinkCutTreeSubTreeAdd()
	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = lct.Alloc(int(weights[i]))
	}
	for i := 0; i < n-1; i++ {
		u, v := edges[i][0], edges[i][1]
		lct.Link(nodes[u], nodes[v])
	}

	for i := 0; i < q; i++ {
		kind := operations[i][0]
		if kind == 0 {
			u1, v1, u2, v2 := operations[i][1], operations[i][2], operations[i][3], operations[i][4]
			lct.Cut(nodes[u1], nodes[v1])
			lct.Link(nodes[u2], nodes[v2])
		} else if kind == 1 {
			node, delta := operations[i][1], operations[i][2]
			lct.Set(nodes[node], lct.Get(nodes[node])+int(delta))
		} else {
			child, parent := operations[i][1], operations[i][2]
			lct.Evert(nodes[parent])
			fmt.Fprintln(out, lct.SubtreeSum(nodes[child]))
		}
	}

}

// Dynamic Tree Subtree Add Subtree Sum
// https://judge.yosupo.jp/problem/dynamic_tree_subtree_add_subtree_sum
// 0 u v w x  删除 u-v 边, 添加 w-x 边
// 1 child parent x
// 2 child parent
func DynamicTreeSubtreeAddSubtreeSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	weights := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}
	edges := make([][2]int32, n-1)
	for i := 0; i < n-1; i++ {
		var v1, v2 int32
		fmt.Fscan(in, &v1, &v2)
		edges[i] = [2]int32{v1, v2}
	}
	operations := make([][5]int32, q)
	for i := 0; i < q; i++ {
		var kind int32
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var u, v, w, x int32
			fmt.Fscan(in, &u, &v, &w, &x)
			operations[i] = [5]int32{kind, u, v, w, x}
		} else if kind == 1 {
			var child, parent, delta int32
			fmt.Fscan(in, &child, &parent, &delta)
			operations[i] = [5]int32{kind, child, parent, delta, 0}
		} else {
			var child, parent int32
			fmt.Fscan(in, &child, &parent)
			operations[i] = [5]int32{kind, child, parent, 0, 0}
		}
	}

	lct := NewLinkCutTreeSubTreeAdd()
	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = lct.Alloc(int(weights[i]))
	}
	for i := 0; i < n-1; i++ {
		u, v := edges[i][0], edges[i][1]
		lct.Link(nodes[u], nodes[v])
	}

	for i := 0; i < q; i++ {
		kind := operations[i][0]
		if kind == 0 {
			u1, v1, u2, v2 := operations[i][1], operations[i][2], operations[i][3], operations[i][4]
			lct.Cut(nodes[u1], nodes[v1])
			lct.Link(nodes[u2], nodes[v2])
		} else if kind == 1 {
			child, parent, delta := operations[i][1], operations[i][2], operations[i][3]
			lct.Evert(nodes[parent])
			lct.SubtreeAdd(nodes[child], int(delta))
		} else {
			child, parent := operations[i][1], operations[i][2]
			lct.Evert(nodes[parent])
			fmt.Fprintln(out, lct.SubtreeSum(nodes[child]))
		}
	}

}

type E = int // 子树和

func e() E               { return 0 }
func add(a, b E) E       { return a + b }
func sub(a, b E) E       { return a - b }
func mul(a E, x int32) E { return a * int(x) }

type Node struct {
	left, right, parent            *Node
	key, sum, lazy, cancel, subsum E
	cnt, subcnt                    int32
	rev                            bool
}

type LinkCutTreeSubTreeAdd struct{}

func NewLinkCutTreeSubTreeAdd() *LinkCutTreeSubTreeAdd {
	return &LinkCutTreeSubTreeAdd{}
}

func (lct *LinkCutTreeSubTreeAdd) Alloc(value E) *Node {
	return &Node{
		key:    value,
		sum:    value,
		lazy:   e(),
		cancel: e(),
		subsum: e(),
		cnt:    1,
	}
}

// 将node作为根节点.
func (lct *LinkCutTreeSubTreeAdd) Evert(node *Node) {
	lct._expose(node)
	lct.Toggle(node)
	lct._pushDown(node)
}

func (lct *LinkCutTreeSubTreeAdd) Link(u, v *Node) {
	lct.Evert(u)
	lct._expose(v)
	u.parent = v
	v.right = u
	lct._pushUp(v)
}

func (lct *LinkCutTreeSubTreeAdd) Cut(u, v *Node) {
	lct.Evert(u)
	lct._expose(v)
	v.left = nil
	u.parent = nil
	lct._pushUp(v)
}

func (lct *LinkCutTreeSubTreeAdd) _pushUp(node *Node) *Node {
	if node == nil {
		return node
	}
	_merge(node, node.left, node.right)
	return node
}

func (lct *LinkCutTreeSubTreeAdd) Toggle(node *Node) {
	if node == nil {
		return
	}
	node.left, node.right = node.right, node.left
	node.rev = !node.rev
}

func (lct *LinkCutTreeSubTreeAdd) Get(node *Node) E {
	lct._expose(node)
	return node.key
}

func (lct *LinkCutTreeSubTreeAdd) Set(node *Node, key E) {
	lct._expose(node)
	node.key = key
	lct._pushUp(node)
}

func (lct *LinkCutTreeSubTreeAdd) SubtreeAdd(node *Node, addVal E) {
	lct._expose(node)
	l := node.left
	if l != nil {
		node.left = nil
		lct._pushUp(node)
	}
	_apply(node, addVal)
	if l != nil {
		node.left = l
		lct._pushUp(node)
	}
}

func (lct *LinkCutTreeSubTreeAdd) SubtreeSum(node *Node) E {
	lct._expose(node)
	return add(node.key, node.subsum)
}

// 返回u,v的最近公共祖先.如果u,v不在同一棵树上,返回nil.
func (lct *LinkCutTreeSubTreeAdd) Lca(u, v *Node) *Node {
	if lct.GetRoot(u) != lct.GetRoot(v) {
		return nil
	}
	lct._expose(u)
	return lct._expose(v)
}

func (lct *LinkCutTreeSubTreeAdd) KthAncestor(node *Node, k int32) *Node {
	lct._expose(node)
	for node != nil {
		lct._pushDown(node)
		if node.right != nil && node.right.cnt > k {
			node = node.right
		} else {
			if node.right != nil {
				k -= node.right.cnt
			}
			if k == 0 {
				return node
			}
			k--
			node = node.left
		}
	}
	return nil
}

// 获取x的根节点.
func (lct *LinkCutTreeSubTreeAdd) GetRoot(x *Node) *Node {
	lct._expose(x)
	for x.left != nil {
		lct._pushDown(x)
		x = x.left
	}
	return x
}

func (lct *LinkCutTreeSubTreeAdd) GetParent(node *Node) *Node {
	lct._expose(node)
	p := node.left
	if p == nil {
		return nil
	}
	for {
		lct._pushDown(p)
		if p.right == nil {
			return p
		}
		p = p.right
	}
}

func (lct *LinkCutTreeSubTreeAdd) IsRoot(t *Node) bool {
	return t.parent == nil || (t.parent.left != t && t.parent.right != t)
}

func (lct *LinkCutTreeSubTreeAdd) _expose(node *Node) *Node {
	var rp *Node
	for cur := node; cur != nil; cur = cur.parent {
		lct._splay(cur)
		lct._pushDown(cur)
		if cur.right != nil {
			_makeNormal(cur, cur.right)
		}
		if rp != nil {
			_fetch(rp)
			_makePrefer(cur, rp)
		}
		cur.right = rp
		rp = cur
	}
	lct._splay(node)
	return rp
}

func (lct *LinkCutTreeSubTreeAdd) _pushRev(t *Node) {
	if t == nil {
		return
	}
	if t.rev {
		if t.left != nil {
			lct.Toggle(t.left)
		}
		if t.right != nil {
			lct.Toggle(t.right)
		}
		t.rev = false
	}
}

func (lct *LinkCutTreeSubTreeAdd) _splay(t *Node) {
	lct._pushDown(t)
	for !lct.IsRoot(t) {
		q := t.parent
		if lct.IsRoot(q) {
			lct._pushRev(q)
			lct._pushRev(t)
			lct._rotate(t)
		} else {
			r := q.parent
			lct._pushRev(r)
			lct._pushRev(q)
			lct._pushRev(t)
			if lct._pos(q) == lct._pos(t) {
				lct._rotate(q)
				lct._rotate(t)
			} else {
				lct._rotate(t)
				lct._rotate(t)
			}
		}
	}
}

func (lct *LinkCutTreeSubTreeAdd) _pushDown(t *Node) {
	if t == nil {
		return
	}
	if t.rev {
		if t.left != nil {
			lct.Toggle(t.left)
		}
		if t.right != nil {
			lct.Toggle(t.right)
		}
		t.rev = false
	}
	if t.left != nil {
		_fetch(t.left)
	}
	if t.right != nil {
		_fetch(t.right)
	}
}

func (lct *LinkCutTreeSubTreeAdd) _pos(t *Node) int32 {
	if t.parent != nil {
		if t.parent.left == t {
			return -1
		}
		if t.parent.right == t {
			return 1
		}
	}
	return 0
}

func (lct *LinkCutTreeSubTreeAdd) _rotate(t *Node) {
	x := t.parent
	y := x.parent
	lct._pushDown(x)
	lct._pushDown(t)
	if lct._pos(t) == -1 {
		if x.left = t.right; t.right != nil {
			t.right.parent = x
		}
		t.right = x
		x.parent = t
	} else {
		if x.right = t.left; t.left != nil {
			t.left.parent = x
		}
		t.left = x
		x.parent = t
	}
	xc := x.cancel
	lct._pushUp(x)
	lct._pushUp(t)
	t.cancel = xc
	t.parent = y
	if t.parent = y; y != nil {
		if y.left == x {
			y.left = t
		}
		if y.right == x {
			y.right = t
		}
	}
}

func _makeNormal(t, other *Node) {
	t.subsum = add(t.subsum, other.sum)
	t.subcnt += other.cnt
}

func _makePrefer(t, other *Node) {
	t.subsum = sub(t.subsum, other.sum)
	t.subcnt -= other.cnt
}

func _merge(t, l, r *Node) {
	var left, right E
	var leftCnt, rightCnt int32
	if l != nil {
		left = l.sum
		leftCnt = l.cnt
	} else {
		left = e()
	}
	if r != nil {
		right = r.sum
		rightCnt = r.cnt
	} else {
		right = e()
	}
	t.sum = add(add(left, t.key), add(t.subsum, right))
	t.cnt = 1 + leftCnt + rightCnt + t.subcnt
	if l != nil {
		l.cancel = t.lazy
	}
	if r != nil {
		r.cancel = t.lazy
	}
}

func _apply(t *Node, addVal E) {
	t.key = add(t.key, addVal)
	t.sum = add(t.sum, mul(addVal, t.cnt))
	t.lazy = add(t.lazy, addVal)
	t.subsum = add(t.subsum, mul(addVal, t.subcnt))
}

func _fetch(t *Node) {
	if t.parent == nil {
		return
	}
	_apply(t, sub(t.parent.lazy, t.cancel))
	t.cancel = t.parent.lazy
}
