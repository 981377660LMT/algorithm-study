// https://ei1333.github.io/library/structure/segment-tree/persistent-segment-tree.hpp
// 可持久化线段树-单点修改区间查询
// Build(leaves []E) *Node : 建树
// Set(root *Node, index int, value E) *Node : 单点修改
// Query(root *Node, left, right int) E : 区间查询

package main

import "fmt"

func main() {
	tree := NewPersistentSegmentTree()
	leaves := []E{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	root := tree.Build(leaves)
	fmt.Println(tree.Query(root, 0, 10)) // 55
	fmt.Print(root)
	root2 := tree.Set(root, 0, 10)
	fmt.Println(tree.Query(root2, 0, 10)) // 65
	fmt.Print(root2)
	root3 := tree.Set(root2, 0, 1)
	fmt.Println(tree.Query(root3, 0, 10)) // 55
	fmt.Print(root)
}

type E = int

func (*PersistentSegmentTree) e() E          { return 0 }
func (*PersistentSegmentTree) op(e1, e2 E) E { return e1 + e2 }

func NewPersistentSegmentTree() *PersistentSegmentTree {
	return &PersistentSegmentTree{}
}

type PersistentSegmentTree struct {
	size int
}

func (t *PersistentSegmentTree) Build(leaves []E) *Node {
	t.size = len(leaves)
	return t.build(0, len(leaves), leaves)
}

func (t *PersistentSegmentTree) Set(root *Node, index int, value E) *Node {
	return t.set(index, value, root, 0, t.size)
}

func (t *PersistentSegmentTree) Query(root *Node, start, end int) E {
	return t.query(start, end, root, 0, t.size)
}

func (t *PersistentSegmentTree) build(l, r int, leaves []E) *Node {
	if l+1 >= r {
		return &Node{data: leaves[l]}
	}
	return t.merge(t.build(l, (l+r)>>1, leaves), t.build((l+r)>>1, r, leaves))
}

func (t *PersistentSegmentTree) set(index int, value E, node *Node, l, r int) *Node {
	if r <= index || index+1 <= l {
		return node
	} else if index <= l && r <= index+1 {
		return &Node{data: value}
	} else {
		return t.merge(t.set(index, value, node.l, l, (l+r)>>1), t.set(index, value, node.r, (l+r)>>1, r))
	}
}

func (t *PersistentSegmentTree) query(start, end int, node *Node, l, r int) E {
	if r <= start || end <= l {
		return t.e()
	} else if start <= l && r <= end {
		return node.data
	} else {
		return t.op(t.query(start, end, node.l, l, (l+r)>>1), t.query(start, end, node.r, (l+r)>>1, r))
	}
}

func (t *PersistentSegmentTree) merge(l, r *Node) *Node {
	return &Node{data: t.op(l.data, r.data), l: l, r: r}
}

type Node struct {
	l, r *Node
	data E
}

func (n *Node) String() string {
	sb := []string{}

	var dfs func(*Node)
	dfs = func(node *Node) {
		if node == nil {
			return
		}
		if node.l == nil && node.r == nil { // leaf
			sb = append(sb, fmt.Sprintf("%v", node.data))
			return
		}
		dfs(node.l)
		dfs(node.r)
	}
	dfs(n)

	return fmt.Sprintf("%v", sb)
}
