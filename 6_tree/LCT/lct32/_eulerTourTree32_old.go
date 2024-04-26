// // 欧拉回路树(Euler Tour Tree)，可删除的并查集.
// // LCT 其实更适用于维护树链的信息，而 ETT 更加适用于维护 子树 的信息。例如，ETT 可以维护子树最小值而 LCT 不能。
// //
// // NewEulerTourTree32
// // Link
// // Cut
// // Get
// // Set
// // UpdateSubTree
// // QuerySubTree
// //
// // TODO: 非常慢

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// )

// func main() {
// 	DynamicTreeVertexAddSubtreeSum()
// }

// func demo() {
// 	ett := NewEulerTourTree(10)
// 	fmt.Println(ett.IsConnected(1, 2))
// 	ett.Link(1, 2)
// 	fmt.Println(ett.IsConnected(1, 2))
// 	ett.Cut(1, 2)
// 	fmt.Println(ett.IsConnected(1, 2))
// 	fmt.Println(ett.QuerySubTree(-1, 1))
// 	// ett.Set(1, 1)
// 	// ett.Set(1, 1)
// 	// ett.Set(1, 1)
// 	// ett.Set(1, 1)
// 	ett.UpdateSubTree(-1, 1, 1)
// 	fmt.Println(ett.QuerySubTree(-1, 1))
// }

// // 动态单点加子树和
// // Dynamic Tree Vertex Add Subtree Sum
// // https://judge.yosupo.jp/problem/dynamic_tree_vertex_add_subtree_sum
// // 0 u v w x  删除 u-v 边, 添加 w-x 边
// // 1 p x  将 p 节点的值加上 x
// // 2 u p  对于边(u,p) 查询结点v的子树的和,p为u的父节点
// func DynamicTreeVertexAddSubtreeSum() {
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, q int32
// 	fmt.Fscan(in, &n, &q)
// 	weights := make([]int32, n)
// 	for i := int32(0); i < n; i++ {
// 		fmt.Fscan(in, &weights[i])
// 	}
// 	edges := make([][2]int32, n-1)
// 	for i := int32(0); i < n-1; i++ {
// 		var v1, v2 int32
// 		fmt.Fscan(in, &v1, &v2)
// 		edges[i] = [2]int32{v1, v2}
// 	}
// 	operations := make([][5]int32, q)
// 	for i := int32(0); i < q; i++ {
// 		var kind int32
// 		fmt.Fscan(in, &kind)
// 		if kind == 0 {
// 			var u, v, w, x int32
// 			fmt.Fscan(in, &u, &v, &w, &x)
// 			operations[i] = [5]int32{kind, u, v, w, x}
// 		} else if kind == 1 {
// 			var p, x int32
// 			fmt.Fscan(in, &p, &x)
// 			operations[i] = [5]int32{kind, p, x, 0, 0}
// 		} else {
// 			var u, p int32
// 			fmt.Fscan(in, &u, &p)
// 			operations[i] = [5]int32{kind, u, p, 0, 0}
// 		}
// 	}

// 	ett := NewEulerTourTree(n)
// 	for i := int32(0); i < n; i++ {
// 		ett.Set(i, int(weights[i]))
// 	}
// 	for i := int32(0); i < n-1; i++ {
// 		u, v := edges[i][0], edges[i][1]
// 		ett.Link(u, v)
// 	}

// 	for i := int32(0); i < q; i++ {
// 		kind := operations[i][0]
// 		if kind == 0 {
// 			u1, v1, u2, v2 := operations[i][1], operations[i][2], operations[i][3], operations[i][4]
// 			ett.Cut(u1, v1)
// 			ett.Link(u2, v2)
// 		} else if kind == 1 {
// 			node, delta := operations[i][1], operations[i][2]
// 			ett.Set(node, ett.Get(node)+int(delta))
// 		} else {
// 			child, parent := operations[i][1], operations[i][2]
// 			fmt.Fprintln(out, ett.QuerySubTree(parent, child))
// 		}
// 	}

// }

// type E = int
// type Id = int32

// func e() E                   { return 0 }
// func id() Id                 { return 0 }
// func op(a, b E) E            { return a + b }
// func mapping(f Id, g E) E    { return int(f) + g }
// func composition(f, g Id) Id { return f + g }

// type EulerTourTree struct {
// 	ptr []map[int32]*Node
// }

// func NewEulerTourTree(n int32) *EulerTourTree {
// 	ptr := make([]map[int32]*Node, n)
// 	for i := int32(0); i < n; i++ {
// 		ptr[i] = make(map[int32]*Node)
// 		ptr[i][i] = &Node{from: i, to: i, size: 1}
// 	}
// 	return &EulerTourTree{ptr: ptr}
// }

// // 连接前必须保证不存在u-v的边.
// func (ett *EulerTourTree) Link(u, v int32) {
// 	tu := Reroot(ett.getNode(u, u))
// 	tv := Reroot(ett.getNode(v, v))
// 	Join(Join(tu, ett.getNode(u, v)), Join(tv, ett.getNode(v, u)))
// }

// // 断开前必须保证存在u-v的边.
// func (ett *EulerTourTree) Cut(u, v int32) {
// 	a, _, c := SplitNode(ett.getNode(u, v), ett.getNode(v, u))
// 	Join(a, c)
// 	delete(ett.ptr[u], v)
// 	delete(ett.ptr[v], u)
// }

// func (ett *EulerTourTree) IsConnected(u, v int32) bool {
// 	return Same(ett.getNode(u, u), ett.getNode(v, v))
// }

// func (ett *EulerTourTree) Get(v int32) E {
// 	t := ett.getNode(v, v)
// 	Splay(t)
// 	return t.value
// }

// func (ett *EulerTourTree) Set(v int32, x E) {
// 	t := ett.getNode(v, v)
// 	Splay(t)
// 	t.value = x
// 	Recalc(t)
// }

// // 更新子树信息.
// // parent为-1时，更新整棵树.
// func (ett *EulerTourTree) UpdateSubTree(parent, child int32, lazy Id) {
// 	if parent != -1 {
// 		ett.Cut(parent, child)
// 	}
// 	t := ett.getNode(child, child)
// 	Splay(t)
// 	t.lazy = composition(lazy, t.lazy)
// 	if parent != -1 {
// 		ett.Link(parent, child)
// 	}
// }

// // 查询子树信息.
// // parent为-1时，查询整棵树.
// func (ett *EulerTourTree) QuerySubTree(parent, child int32) E {
// 	if parent != -1 {
// 		ett.Cut(parent, child)
// 	}
// 	t := ett.getNode(child, child)
// 	Splay(t)
// 	res := t.sum
// 	if parent != -1 {
// 		ett.Link(parent, child)
// 	}
// 	return res
// }

// func (ett *EulerTourTree) getNode(u, v int32) *Node {
// 	nexts := ett.ptr[u]
// 	if t, ok := nexts[v]; ok {
// 		return t
// 	} else {
// 		t = NewNode(u, v)
// 		nexts[v] = t
// 		return t
// 	}
// }

// type Node struct {
// 	children [2]*Node
// 	parent   *Node
// 	from     int32
// 	to       int32
// 	size     int32
// 	value    E
// 	sum      E
// 	lazy     Id
// }

// func NewNode(from, to int32) *Node {
// 	var size int32
// 	if from == to {
// 		size = 1
// 	}
// 	return &Node{
// 		from:  from,
// 		to:    to,
// 		size:  size,
// 		value: e(),
// 		sum:   e(),
// 		lazy:  id(),
// 	}
// }

// func GetRoot(t *Node) *Node {
// 	if t == nil {
// 		return nil
// 	}
// 	for t.parent != nil {
// 		t = t.parent
// 	}
// 	return t
// }

// func Same(t, s *Node) bool {
// 	if t != nil {
// 		Splay(t)
// 	}
// 	if s != nil {
// 		Splay(s)
// 	}
// 	return GetRoot(t) == GetRoot(s)
// }

// func Reroot(t *Node) *Node {
// 	a, b := Split(t)
// 	return Join(b, a)
// }

// func Size(t *Node) int32 {
// 	if t == nil {
// 		return 0
// 	}
// 	return t.size
// }

// func Recalc(t *Node) *Node {
// 	if t == nil {
// 		return t
// 	}
// 	tmp := int32(0)
// 	if t.from == t.to {
// 		tmp = 1
// 	}
// 	t.size = tmp + Size(t.children[0]) + Size(t.children[1])
// 	t.sum = t.value
// 	if t.children[0] != nil {
// 		t.sum = op(t.children[0].sum, t.sum)
// 	}
// 	if t.children[1] != nil {
// 		t.sum = op(t.sum, t.children[1].sum)
// 	}
// 	return t
// }

// func PushDown(t *Node) {
// 	if t.lazy != id() {
// 		t.value = mapping(t.lazy, t.value)
// 		if left := t.children[0]; left != nil {
// 			left.lazy = composition(t.lazy, left.lazy)
// 			left.sum = mapping(t.lazy, left.sum)
// 		}
// 		if right := t.children[1]; right != nil {
// 			right.lazy = composition(t.lazy, right.lazy)
// 			right.sum = mapping(t.lazy, right.sum)
// 		}
// 		t.lazy = id()
// 	}
// 	Recalc(t)
// }

// func Join(l, r *Node) *Node {
// 	if l == nil {
// 		return r
// 	}
// 	if r == nil {
// 		return l
// 	}
// 	for l.children[1] != nil {
// 		l = l.children[1]
// 	}
// 	Splay(l)
// 	l.children[1] = r
// 	r.parent = l
// 	return Recalc(l)
// }

// func Split(t *Node) (*Node, *Node) {
// 	Splay(t)
// 	s := t.children[0]
// 	t.children[0] = nil
// 	if s != nil {
// 		s.parent = nil
// 	}
// 	return s, Recalc(t)
// }

// func Split2(t *Node) (*Node, *Node) {
// 	Splay(t)
// 	l := t.children[0]
// 	r := t.children[1]
// 	t.children[0] = nil
// 	if l != nil {
// 		l.parent = nil
// 	}
// 	t.children[1] = nil
// 	if r != nil {
// 		r.parent = nil
// 	}
// 	return l, r
// }

// func SplitNode(s, t *Node) (a, b, c *Node) {
// 	a, b = Split2(s)
// 	if Same(a, t) {
// 		c, d := Split2(t)
// 		return c, d, b
// 	} else {
// 		c, d := Split2(t)
// 		return a, c, d
// 	}
// }

// func Rotate(t *Node, b bool) {
// 	var v uint8
// 	if b {
// 		v = 1
// 	}
// 	p := t.parent
// 	g := p.parent
// 	p.children[1^v] = t.children[v]
// 	if p.children[1^v] != nil {
// 		t.children[v].parent = p
// 	}
// 	t.children[v] = p
// 	p.parent = t
// 	Recalc(p)
// 	Recalc(t)
// 	t.parent = g
// 	if t.parent != nil {
// 		if g.children[0] == p {
// 			g.children[0] = t
// 		} else {
// 			g.children[1] = t
// 		}
// 		Recalc(g)
// 	}
// }

// func Splay(t *Node) {
// 	PushDown(t)
// 	for t.parent != nil {
// 		p := t.parent
// 		g := p.parent
// 		if g == nil {
// 			PushDown(p)
// 			PushDown(t)
// 			Rotate(t, p.children[0] == t)
// 		} else {
// 			PushDown(g)
// 			PushDown(p)
// 			PushDown(t)
// 			var b uint8
// 			var f bool
// 			if g.children[0] == p {
// 				b = 1
// 				f = true
// 			}
// 			if p.children[1^b] == t {
// 				Rotate(p, f)
// 				Rotate(t, f)
// 			} else {
// 				Rotate(t, !f)
// 				Rotate(t, f)
// 			}
// 		}
// 	}
// }

package main

func main() {

}
