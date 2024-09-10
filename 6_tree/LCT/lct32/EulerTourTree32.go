// 欧拉树 Euler Tour Tree
// 专门用于处理动态树的子树问题.
// !非常慢
//
// api:
//  NewEtt(n int32, f func(i int32) E) *Ett
//  Build(adjList [][]int32)
//  Link(u, v int32)
//  LinkSafely(u, v int32) bool
//  Cut(u, v int32)
//  CutSafely(u, v int32) bool
//  Reroot(v int32)
//  Find(v int32) *ettNode
//  IsConnected(u, v int32) bool
//  Get(v int32) E
//  Set(v int32, x E)
//  UpdateSubtree(v, p int32, f Id)
//  QuerySubtree(v, p int32) E
//  Part() int32

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
	// demo()
	// DynamicTreeVertexAddSubtreeSum()

}

func demo() {
	ett := NewEtt(10, func(i int32) E { return int(i) })
	fmt.Println(ett.IsConnected(1, 2))
	ett.Link(1, 2)
	fmt.Println(ett.IsConnected(1, 2))
	ett.Cut(1, 2)
	fmt.Println(ett.IsConnected(1, 2))
	fmt.Println(ett.QuerySubtree(1, -1))

	adjList := [][]int32{{1}, {0, 2}, {1}}
	ett.Build(adjList)
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

	var n, q int32
	fmt.Fscan(in, &n, &q)
	weights := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var v1, v2 int32
		fmt.Fscan(in, &v1, &v2)
		tree[v1] = append(tree[v1], v2)
		tree[v2] = append(tree[v2], v1)
	}

	operations := make([][5]int32, q)
	for i := int32(0); i < q; i++ {
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

	ett := NewEtt(n, func(i int32) E { return int(weights[i]) })
	ett.Build(tree)

	for i := int32(0); i < q; i++ {
		kind := operations[i][0]
		if kind == 0 {
			u1, v1, u2, v2 := operations[i][1], operations[i][2], operations[i][3], operations[i][4]
			ett.Cut(u1, v1)
			ett.Link(u2, v2)
		} else if kind == 1 {
			node, delta := operations[i][1], operations[i][2]
			ett.Set(node, ett.Get(node)+int(delta))
		} else {
			child, parent := operations[i][1], operations[i][2]
			fmt.Fprintln(out, ett.QuerySubtree(child, parent))
		}
	}
}

type E = int
type Id = int32

func e() E                   { return 0 }
func id() Id                 { return 0 }
func op(a, b E) E            { return a + b }
func mapping(f Id, g E) E    { return int(f) + g }
func composition(f, g Id) Id { return f + g }

type ettNode struct {
	key, data        E
	lazy             Id
	par, left, right *ettNode
}

func newEttNode(key E, lazy Id) *ettNode {
	return &ettNode{key: key, data: key, lazy: lazy}
}

type Ett struct {
	n            int32
	groupNumbers int32
	ptrVertex    []*ettNode
	ptrEdge      map[int]*ettNode
}

// 点权.
func NewEtt(n int32, f func(i int32) E) *Ett {
	res := &Ett{n: n, groupNumbers: n, ptrEdge: make(map[int]*ettNode, 2*n)}
	res._build(n, f)
	return res
}

// 按照临接表连接所有边.
func (ett *Ett) Build(adjList [][]int32) {
	n := len(adjList)
	visited := make([]bool, n)
	a := make([]int, 0, 3*n)

	var dfs func(v, p int32)
	dfs = func(v, p int32) {
		a = append(a, int(v)*n+int(v))
		for _, x := range adjList[v] {
			if x == p {
				continue
			}
			a = append(a, int(v)*n+int(x))
			dfs(x, v)
			a = append(a, int(x)*n+int(v))
		}
	}

	var rec func(l, r int32) *ettNode
	rec = func(l, r int32) *ettNode {
		mid := (l + r) >> 1
		u, v := a[mid]/n, a[mid]%n
		var node *ettNode
		if u == v {
			node = ett.ptrVertex[u]
			visited[u] = true
		} else {
			node = newEttNode(e(), id())
			ett.ptrEdge[a[mid]] = node
		}
		if l != mid {
			node.left = rec(l, mid)
			node.left.par = node
		}
		if mid+1 != r {
			node.right = rec(mid+1, r)
			node.right.par = node
		}
		ett._update(node)
		return node
	}

	for root := int32(0); root < int32(n); root++ {
		if visited[root] {
			continue
		}
		a = a[:0]
		dfs(root, -1)
		rec(0, int32(len(a)))
	}
}

// 要保证不存在u-v的边.
func (ett *Ett) Link(u, v int32) {
	ett.Reroot(u)
	ett.Reroot(v)
	uvNode := newEttNode(e(), id())
	vuNode := newEttNode(e(), id())
	n, ui, vi := int(ett.n), int(u), int(v)
	ett.ptrEdge[ui*n+vi] = uvNode
	ett.ptrEdge[vi*n+ui] = vuNode
	uNode, vNode := ett.ptrVertex[u], ett.ptrVertex[v]
	ett._merge(uNode, uvNode)
	ett._merge(uvNode, vNode)
	ett._merge(vNode, vuNode)
	ett.groupNumbers--
}

func (ett *Ett) LinkSafely(u, v int32) bool {
	if ett.IsConnected(u, v) {
		return false
	}
	ett.Link(u, v)
	return true
}

// 要保证存在u-v的边.
func (ett *Ett) Cut(u, v int32) {
	ett.Reroot(v)
	ett.Reroot(u)
	n := int(ett.n)
	ui, vi := int(u), int(v)
	uvNode := ett.ptrEdge[ui*n+vi]
	vuNode := ett.ptrEdge[vi*n+ui]
	delete(ett.ptrEdge, ui*n+vi)
	delete(ett.ptrEdge, vi*n+ui)
	var a, c *ettNode
	a, _ = ett._splitLeft(uvNode)
	_, c = ett._splitRight(vuNode)
	a = ett._pop(a)
	c = ett._popleft(c)
	ett._merge(a, c)
	ett.groupNumbers++
}

func (ett *Ett) CutSafely(u, v int32) bool {
	if _, has := ett.ptrEdge[int(u)*int(ett.n)+int(v)]; !has {
		return false
	}
	ett.Cut(u, v)
	return true
}

// Evert.
func (ett *Ett) Reroot(v int32) {
	node := ett.ptrVertex[v]
	x, y := ett._splitRight(node)
	ett._merge(y, x)
	ett._splay(node)
}

func (ett *Ett) Find(v int32) *ettNode {
	return ett._leftSplay(ett.ptrVertex[v])
}

func (ett *Ett) IsConnected(u, v int32) bool {
	uNode, vNode := ett.ptrVertex[u], ett.ptrVertex[v]
	ett._splay(uNode)
	ett._splay(vNode)
	return uNode.par != nil || uNode == vNode
}

func (ett *Ett) Get(v int32) E {
	node := ett.ptrVertex[v]
	ett._splay(node)
	return node.key
}

func (ett *Ett) Set(v int32, x E) {
	node := ett.ptrVertex[v]
	ett._splay(node)
	node.key = x
	ett._update(node)
}

func (ett *Ett) UpdateSubtree(cur, parent int32, lazy Id) {
	node := ett.ptrVertex[cur]
	ett.Reroot(cur)
	if parent == -1 {
		ett._splay(node)
		node.key = mapping(lazy, node.key)
		node.data = mapping(lazy, node.data)
		node.lazy = composition(lazy, node.lazy)
		return
	}
	ett.Reroot(parent)
	ni, pi, vi := int(ett.n), int(parent), int(cur)
	a, b := ett._splitRight(ett.ptrEdge[vi*ni+pi])
	b, d := ett._splitLeft(ett.ptrEdge[pi*ni+vi])
	ett._splay(node)
	node.key = mapping(lazy, node.key)
	node.data = mapping(lazy, node.data)
	node.lazy = composition(lazy, node.lazy)
	ett._propagate(node)
	ett._merge(a, b)
	ett._merge(b, d)
}

func (ett *Ett) QuerySubtree(cur, parent int32) E {
	node := ett.ptrVertex[cur]
	ett.Reroot(cur)
	if parent == -1 {
		ett._splay(node)
		return node.data
	}
	ett.Reroot(parent)
	ni, pi, vi := int(ett.n), int(parent), int(cur)
	a, b := ett._splitRight(ett.ptrEdge[pi*ni+vi])
	b, d := ett._splitLeft(ett.ptrEdge[vi*ni+pi])
	ett._splay(node)
	res := node.data
	ett._merge(a, b)
	ett._merge(b, d)
	return res
}

func (ett *Ett) Part() int32 {
	return ett.groupNumbers
}

func (ett *Ett) _build(n int32, f func(i int32) E) {
	ett.ptrVertex = make([]*ettNode, n)
	for i := int32(0); i < n; i++ {
		ett.ptrVertex[i] = newEttNode(f(i), id())
	}
}

func (ett *Ett) _popleft(v *ettNode) *ettNode {
	v = ett._leftSplay(v)
	if v.right != nil {
		v.right.par = nil
	}
	return v.right
}

func (ett *Ett) _pop(v *ettNode) *ettNode {
	v = ett._rightSplay(v)
	if v.left != nil {
		v.left.par = nil
	}
	return v.left
}

func (ett *Ett) _splitLeft(v *ettNode) (*ettNode, *ettNode) {
	ett._splay(v)
	x, y := v, v.right
	if y != nil {
		y.par = nil
	}
	x.right = nil
	ett._update(x)
	return x, y
}

func (ett *Ett) _splitRight(v *ettNode) (*ettNode, *ettNode) {
	ett._splay(v)
	x := v.left
	y := v
	if x != nil {
		x.par = nil
	}
	y.left = nil
	ett._update(y)
	return x, y
}

func (ett *Ett) _merge(u, v *ettNode) {
	if u == nil || v == nil {
		return
	}
	u = ett._rightSplay(u)
	ett._splay(v)
	u.right = v
	v.par = u
	ett._update(u)
}

func (ett *Ett) _splay(node *ettNode) {
	ett._propagate(node)
	for node.par != nil && node.par.par != nil {
		pnode := node.par
		gnode := pnode.par
		ett._propagate(gnode)
		ett._propagate(pnode)
		ett._propagate(node)
		node.par = gnode.par
		var tmp1, tmp2 *ettNode
		if (gnode.left == pnode) == (pnode.left == node) {
			if pnode.left == node {
				tmp1 = node.right
				pnode.left = tmp1
				node.right = pnode
				pnode.par = node
				tmp2 = pnode.right
				gnode.left = tmp2
				pnode.right = gnode
				gnode.par = pnode
			} else {
				tmp1 = node.left
				pnode.right = tmp1
				node.left = pnode
				pnode.par = node
				tmp2 = pnode.left
				gnode.right = tmp2
				pnode.left = gnode
				gnode.par = pnode
			}
			if tmp1 != nil {
				tmp1.par = pnode
			}
			if tmp2 != nil {
				tmp2.par = gnode
			}
		} else {
			if pnode.left == node {
				tmp1 = node.right
				pnode.left = tmp1
				node.right = pnode
				tmp2 = node.left
				gnode.right = tmp2
				node.left = gnode
				pnode.par = node
				gnode.par = node
			} else {
				tmp1 = node.left
				pnode.right = tmp1
				node.left = pnode
				tmp2 = node.right
				gnode.left = tmp2
				node.right = gnode
				pnode.par = node
				gnode.par = node
			}
			if tmp1 != nil {
				tmp1.par = pnode
			}
			if tmp2 != nil {
				tmp2.par = gnode
			}
		}
		ett._update(gnode)
		ett._update(pnode)
		ett._update(node)
		if node.par == nil {
			return
		}
		if node.par.left == gnode {
			node.par.left = node
		} else {
			node.par.right = node
		}
	}

	if node.par == nil {
		return
	}
	pnode := node.par
	ett._propagate(pnode)
	ett._propagate(node)
	if pnode.left == node {
		pnode.left = node.right
		if pnode.left != nil {
			pnode.left.par = pnode
		}
		node.right = pnode
	} else {
		pnode.right = node.left
		if pnode.right != nil {
			pnode.right.par = pnode
		}
		node.left = pnode
	}
	node.par = nil
	pnode.par = node
	ett._update(pnode)
	ett._update(node)
}

func (ett *Ett) _leftSplay(node *ettNode) *ettNode {
	ett._splay(node)
	for node.left != nil {
		node = node.left
	}
	ett._splay(node)
	return node
}

func (ett *Ett) _rightSplay(node *ettNode) *ettNode {
	ett._splay(node)
	for node.right != nil {
		node = node.right
	}
	ett._splay(node)
	return node
}

func (ett *Ett) _propagate(node *ettNode) {
	if node == nil || node.lazy == id() {
		return
	}
	if node.left != nil {
		node.left.key = mapping(node.lazy, node.left.key)
		node.left.data = mapping(node.lazy, node.left.data)
		node.left.lazy = composition(node.lazy, node.left.lazy)
	}
	if node.right != nil {
		node.right.key = mapping(node.lazy, node.right.key)
		node.right.data = mapping(node.lazy, node.right.data)
		node.right.lazy = composition(node.lazy, node.right.lazy)
	}
	node.lazy = id()
}

func (ett *Ett) _update(node *ettNode) {
	ett._propagate(node.left)
	ett._propagate(node.right)
	node.data = node.key
	if node.left != nil {
		node.data = op(node.left.data, node.data)
	}
	if node.right != nil {
		node.data = op(node.data, node.right.data)
	}
}
