// 比树状数组快一些

package main

import (
	"sort"
)

// https://leetcode.cn/problems/min-cost-to-connect-all-points/
func minCostConnectPoints(points [][]int) int {
	myPoints := make([][2]int, len(points))
	for i, p := range points {
		myPoints[i] = [2]int{p[0], p[1]}
	}
	res, _ := ManhattanMST(myPoints)
	return res
}

// 曼哈顿最小生成树
func ManhattanMST(points [][2]int) (mst int, mstEdges [][2]int) {
	points = append(points[:0:0], points...)
	n := len(points)
	data := make([][3]int, 0, 4*n)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}

	for a := 0; a < 2; a++ {
		for i := range points {
			points[i][0] = -points[i][0]
		}
		for b := 0; b < 2; b++ {
			for i := range points {
				points[i][0], points[i][1] = points[i][1], points[i][0]
			}
			sort.Slice(idx, func(i, j int) bool {
				return points[idx[i]][0]+points[idx[i]][1] < points[idx[j]][0]+points[idx[j]][1]
			})

			mp := NTM(func(a, b int) int { return a - b })
			for _, i := range idx {
				x, y := points[i][0], points[i][1]
				for it := mp.LowerBound(-y); !it.IsEnd(); it = mp.Erase(it) {
					j := it.Value()
					xj, yj := points[j][0], points[j][1]
					dx, dy := x-xj, y-yj
					if dy > dx {
						break
					}
					data = append(data, [3]int{dx + dy, i, j})
				}
				mp.Set(-y, i)
			}
		}
	}

	sort.Slice(data, func(i, j int) bool { return data[i][0] < data[j][0] })
	uf := _nuf(n)
	for _, e := range data {
		cost, i, j := e[0], e[1], e[2]
		if uf.Union(i, j) {
			mst += cost
			mstEdges = append(mstEdges, [2]int{i, j})
		}
	}
	return
}

type Key = int
type Value = int

type TM struct {
	tree *_Tree
}

func NTM(comparator func(a, b Key) int) *TM {
	return &TM{tree: _NewRedBlackTree(comparator)}
}

func (m *TM) Set(key Key, value Value) {
	m.tree.Put(key, value)
}

func (m *TM) Iterator() Iterator {
	return m.tree.Iterator()
}

// 返回删除元素的后继元素的迭代器，如果删除的是最后一个元素，则返回end()迭代器。
func (m *TM) Erase(it Iterator) Iterator {
	node := it.node
	it.Next()
	m.tree.DiscardNode(node)
	return it
}

// 返回一个迭代器，指向键值>= key的第一个元素。
func (m *TM) LowerBound(key Key) Iterator {
	lower, ok := m.tree.Ceiling(key)
	if !ok {
		it := m.tree.Iterator()
		it.End()
		return it
	}
	return m.tree.IteratorAt(lower)
}

// !region RedBlackTree
//

type _color bool

const (
	_black, _red _color = true, false
)

// _Tree holds elements of the red-black tree
type _Tree struct {
	Root       *_RBNode
	size       int
	Comparator func(a, b Key) int
}

// _RBNode is a single element within the tree
type _RBNode struct {
	Key    Key
	Value  Value
	color  _color
	Left   *_RBNode
	Right  *_RBNode
	Parent *_RBNode
}

// NewWith instantiates a red-black tree with the custom comparator.
func _NewRedBlackTree(comparator func(a, b Key) int) *_Tree {
	return &_Tree{Comparator: comparator}
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *_Tree) Put(key Key, value Value) {
	var insertedNode *_RBNode
	if tree.Root == nil {
		tree.Root = &_RBNode{Key: key, Value: value, color: _red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:
				node.Key = key
				node.Value = value
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = &_RBNode{Key: key, Value: value, color: _red}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &_RBNode{Key: key, Value: value, color: _red}
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++
}

func (tree *_Tree) DiscardNode(node *_RBNode) bool {
	if node == nil {
		return false
	}

	var child *_RBNode
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node.Value = pred.Value
		node = pred
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == _black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = _black
		}
	}
	tree.size--
	return true
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *_Tree) Left() *_RBNode {
	var parent *_RBNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

func (tree *_Tree) Ceiling(key Key) (ceiling *_RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

func (tree *_Tree) Higher(key Key) (higher *_RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare < 0:
			higher, found = node, true
			node = node.Left
		case compare >= 0:
			node = node.Right
		}
	}
	if found {
		return higher, true
	}
	return nil, false
}

func (tree *_Tree) lookup(key Key) *_RBNode {
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (node *_RBNode) grandparent() *_RBNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *_RBNode) uncle() *_RBNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *_RBNode) sibling() *_RBNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (tree *_Tree) rotateLeft(node *_RBNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *_Tree) rotateRight(node *_RBNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *_Tree) replaceNode(old *_RBNode, new *_RBNode) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *_Tree) insertCase1(node *_RBNode) {
	if node.Parent == nil {
		node.color = _black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *_Tree) insertCase2(node *_RBNode) {
	if nodeColor(node.Parent) == _black {
		return
	}
	tree.insertCase3(node)
}

func (tree *_Tree) insertCase3(node *_RBNode) {
	uncle := node.uncle()
	if nodeColor(uncle) == _red {
		node.Parent.color = _black
		uncle.color = _black
		node.grandparent().color = _red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *_Tree) insertCase4(node *_RBNode) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *_Tree) insertCase5(node *_RBNode) {
	node.Parent.color = _black
	grandparent := node.grandparent()
	grandparent.color = _red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *_RBNode) maximumNode() *_RBNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *_Tree) deleteCase1(node *_RBNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *_Tree) deleteCase2(node *_RBNode) {
	sibling := node.sibling()
	if nodeColor(sibling) == _red {
		node.Parent.color = _red
		sibling.color = _black
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *_Tree) deleteCase3(node *_RBNode) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == _black &&
		nodeColor(sibling) == _black &&
		nodeColor(sibling.Left) == _black &&
		nodeColor(sibling.Right) == _black {
		sibling.color = _red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *_Tree) deleteCase4(node *_RBNode) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == _red &&
		nodeColor(sibling) == _black &&
		nodeColor(sibling.Left) == _black &&
		nodeColor(sibling.Right) == _black {
		sibling.color = _red
		node.Parent.color = _black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *_Tree) deleteCase5(node *_RBNode) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == _black &&
		nodeColor(sibling.Left) == _red &&
		nodeColor(sibling.Right) == _black {
		sibling.color = _red
		sibling.Left.color = _black
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == _black &&
		nodeColor(sibling.Right) == _red &&
		nodeColor(sibling.Left) == _black {
		sibling.color = _red
		sibling.Right.color = _black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *_Tree) deleteCase6(node *_RBNode) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = _black
	if node == node.Parent.Left && nodeColor(sibling.Right) == _red {
		sibling.Right.color = _black
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == _red {
		sibling.Left.color = _black
		tree.rotateRight(node.Parent)
	}
}

func nodeColor(node *_RBNode) _color {
	if node == nil {
		return _black
	}
	return node.color
}

type Iterator struct {
	tree     *_Tree
	node     *_RBNode
	position _position
}

type _position byte

const (
	begin, between, end _position = 0, 1, 2
)

func (tree *_Tree) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, position: begin}
}

func (tree *_Tree) IteratorAt(node *_RBNode) Iterator {
	return Iterator{tree: tree, node: node, position: between}
}

func (iterator *Iterator) Next() bool {
	if iterator.position == end {
		goto end
	}
	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	for iterator.node.Parent != nil {
		node := iterator.node
		iterator.node = iterator.node.Parent
		if node == iterator.node.Left {
			goto between
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false

between:
	iterator.position = between
	return true
}

func (iterator *Iterator) Value() Value {
	return iterator.node.Value
}

func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
}

func (iterator *Iterator) IsEnd() bool {
	return iterator.position == end
}

func _nuf(n int) *_uf {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_uf{
		Part:   n,
		size:   n,
		rank:   rank,
		parent: parent,
	}
}

type _uf struct {
	Part   int
	size   int
	rank   []int
	parent []int
}

func (ufa *_uf) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_uf) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}
