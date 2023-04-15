package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://atcoder.jp/contests/code-thanks-festival-2017/tasks/code_thanks_festival_2017_h
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	uf := NewRetroactiveUnionFind(n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		uf.Union(a-1, b-1, i+1)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		l, r := 0, m+1
		for r-l > 1 {
			mid := (l + r) / 2
			if uf.IsConnected(x, y, mid) {
				r = mid
			} else {
				l = mid
			}
		}
		if r == m+1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, r)
		}
	}
}

// 完全可追溯化并查集.
//  LinkCutTree 维护最小生成树.
type RetroactiveUnionFind struct {
	v      int
	lc     *LinkCutTree
	edge   []*Node
	nodes  [][2]int
	curPos int
}

func NewRetroactiveUnionFind(nodeSize int) *RetroactiveUnionFind {
	return &RetroactiveUnionFind{
		v:     nodeSize,
		lc:    NewLinkCutTree(nodeSize),
		edge:  make([]*Node, nodeSize-1),
		nodes: make([][2]int, nodeSize-1),
	}
}

func (ruf *RetroactiveUnionFind) Union(id1, id2 int, time int) bool {
	pos := ruf.curPos
	if ruf.curPos == ruf.v-1 || ruf.lc.Connected(id1, id2) {
		first, second := ruf.lc.Query(id1, id2)
		if first <= time {
			return false
		}
		pos = second
		p := ruf.nodes[pos]
		ruf.lc.Cut(p[0], ruf.edge[pos])
		ruf.lc.Cut(p[1], ruf.edge[pos])
	} else {
		ruf.curPos++
	}
	ruf.edge[pos] = NewNode(pos, time)
	ruf.lc.Link(id1, ruf.edge[pos])
	ruf.lc.Link(id2, ruf.edge[pos])
	ruf.nodes[pos] = [2]int{id1, id2}
	return true
}

func (ruf *RetroactiveUnionFind) IsConnected(id1, id2 int, time int) bool {
	if !ruf.lc.Connected(id1, id2) {
		return false
	}
	first, _ := ruf.lc.Query(id1, id2)
	return first <= time
}

type Node struct {
	id, alid         int
	val, al          int
	left, right, par *Node
	rev              bool
}

func NewNode(id int, val int) *Node {
	return &Node{id: id, alid: id, val: val, al: val}
}

func (n *Node) isRoot() bool {
	return n.par == nil || (n.par.left != n && n.par.right != n)
}

func (n *Node) push() {
	if !n.rev {
		return
	}
	n.left, n.right = n.right, n.left
	if n.left != nil {
		n.left.rev = !n.left.rev
	}
	if n.right != nil {
		n.right.rev = !n.right.rev
	}
	n.rev = false
}

func (n *Node) eval() {
	n.alid = n.id
	n.al = n.val
	if n.left != nil {
		n.left.push()
		if n.al < n.left.al {
			n.al = n.left.al
			n.alid = n.left.alid
		}
	}
	if n.right != nil {
		n.right.push()
		if n.al < n.right.al {
			n.al = n.right.al
			n.alid = n.right.alid
		}
	}
}

type LinkCutTree struct {
	arr []*Node
}

const INF int = 1e18

func NewLinkCutTree(nodeSize int) *LinkCutTree {
	arr := make([]*Node, nodeSize)
	for i := 0; i < nodeSize; i++ {
		arr[i] = NewNode(-1, -INF)
	}
	return &LinkCutTree{arr: arr}
}

func (lct *LinkCutTree) Connected(id1 int, id2 int) bool {
	return lct.connected(lct.arr[id1], lct.arr[id2])
}

func (lct *LinkCutTree) Link(verId int, edge *Node) {
	lct.link(lct.arr[verId], edge)
}

func (lct *LinkCutTree) Cut(verId int, edge *Node) {
	lct.cut2(lct.arr[verId], edge)
}

func (lct *LinkCutTree) Query(id1 int, id2 int) (int, int) {
	return lct.query(lct.arr[id1], lct.arr[id2])
}

func _MakeNode() *Node {
	return &Node{}
}

func (lct *LinkCutTree) rotate(u *Node, right bool) {
	p := u.par
	g := p.par
	if right {
		p.left = u.right
		if p.left != nil {
			u.right.par = p
		}
		u.right = p
		p.par = u
	} else {
		p.right = u.left
		if p.right != nil {
			u.left.par = p
		}
		u.left = p
		p.par = u
	}
	p.eval()
	u.eval()
	u.par = g
	if g == nil {
		return
	}
	if g.left == p {
		g.left = u
	}
	if g.right == p {
		g.right = u
	}
	g.eval()
}

func (lct *LinkCutTree) splay(u *Node) {
	for !u.isRoot() {
		p := u.par
		gp := p.par
		if p.isRoot() {
			p.push()
			u.push()
			lct.rotate(u, u == p.left)
		} else {
			gp.push()
			p.push()
			u.push()
			flag := u == p.left
			if flag == (p == gp.left) {
				lct.rotate(p, flag)
				lct.rotate(u, flag)
			} else {
				lct.rotate(u, flag)
				lct.rotate(u, !flag)
			}
		}
	}
	u.push()
}

func (lct *LinkCutTree) access(u *Node) *Node {
	var last *Node
	for v := u; v != nil; v = v.par {
		lct.splay(v)
		v.right = last
		v.eval()
		last = v
	}
	lct.splay(u)
	return last
}

func (lct *LinkCutTree) evert(u *Node) {
	lct.access(u)
	u.rev = !u.rev
	u.push()
}

func (lct *LinkCutTree) connected(u, v *Node) bool {
	lct.access(u)
	lct.access(v)
	return u.par != nil
}

func (lct *LinkCutTree) link(u, v *Node) {
	lct.evert(u)
	u.par = v
}

func (lct *LinkCutTree) cut(u *Node) {
	lct.access(u)
	u.left.par = nil
	u.left = nil
	u.eval()
}

func (lct *LinkCutTree) cut2(u, v *Node) {
	lct.access(u)
	lct.access(v)
	if u.isRoot() {
		u.par = nil
	} else {
		v.left.par = nil
		v.left = nil
		v.eval()
	}
}

func (lct *LinkCutTree) query(u, v *Node) (int, int) {
	lct.evert(u)
	lct.access(v)
	return v.al, v.alid
}
