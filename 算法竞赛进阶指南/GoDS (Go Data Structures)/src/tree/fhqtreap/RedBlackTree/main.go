// 红黑树

package main

import "fmt"

func main() {
	tree := NewRedBlackTree(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})
	tree.Put(4, "d")
	tree.Put(2, "b")
	tree.Put(1, "a")
	tree.Put(3, "c")
	fmt.Println(tree.Get(1)) // a true
	fmt.Println(tree.Floor(1))
	fmt.Println(tree.Higher(1))
}

type Key = interface{}
type Value = interface{}

type rbColor bool

const (
	rbBlack, rbRed rbColor = true, false
)

// RBTree holds elements of the red-black tree
type RBTree struct {
	Root       *RBNode
	size       int
	Comparator func(a, b Key) int
}

// RBNode is a single element within the tree
type RBNode struct {
	Key    Key
	Value  Value
	color  rbColor
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
}

// NewWith instantiates a red-black tree with the custom comparator.
func NewRedBlackTree(comparator func(a, b Key) int) *RBTree {
	return &RBTree{Comparator: comparator}
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *RBTree) Put(key Key, value Value) {
	var insertedNode *RBNode
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		tree.Comparator(key, key)
		tree.Root = &RBNode{Key: key, Value: value, color: rbRed}
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
					node.Left = &RBNode{Key: key, Value: value, color: rbRed}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &RBNode{Key: key, Value: value, color: rbRed}
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

func (tree *RBTree) Get(key Key) (value Value, found bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.Value, true
	}
	return
}

func (tree *RBTree) GetNode(key Key) *RBNode {
	return tree.lookup(key)
}

func (tree *RBTree) Discard(key Key) bool {
	node := tree.lookup(key)
	return tree.DiscardNode(node)
}

func (tree *RBTree) DiscardNode(node *RBNode) bool {
	if node == nil {
		return false
	}

	var child *RBNode
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
		if node.color == rbBlack {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = rbBlack
		}
	}
	tree.size--
	return true
}

func (tree *RBTree) Empty() bool {
	return tree.size == 0
}

func (tree *RBTree) Size() int {
	return tree.size
}

func (node *RBNode) Size() int {
	if node == nil {
		return 0
	}
	size := 1
	if node.Left != nil {
		size += node.Left.Size()
	}
	if node.Right != nil {
		size += node.Right.Size()
	}
	return size
}

func (tree *RBTree) Keys() []Key {
	keys := make([]Key, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

func (tree *RBTree) Values() []Value {
	values := make([]Value, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *RBTree) Left() *RBNode {
	var parent *RBNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *RBTree) Right() *RBNode {
	var parent *RBNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

func (tree *RBTree) Floor(key Key) (floor *RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left
		case compare > 0:
			floor, found = node, true
			node = node.Right
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

func (tree *RBTree) Lower(key Key) (lower *RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare <= 0:
			node = node.Left
		case compare > 0:
			lower, found = node, true
			node = node.Right
		}
	}
	if found {
		return lower, true
	}
	return nil, false
}

func (tree *RBTree) Ceiling(key Key) (ceiling *RBNode, found bool) {
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

func (tree *RBTree) Higher(key Key) (higher *RBNode, found bool) {
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

func (tree *RBTree) Clear() {
	tree.Root = nil
	tree.size = 0
}

func (tree *RBTree) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *RBNode) String() string {
	return fmt.Sprintf("%v", node.Key)
}

func output(node *RBNode, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}

func (tree *RBTree) lookup(key Key) *RBNode {
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

func (node *RBNode) grandparent() *RBNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *RBNode) uncle() *RBNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *RBNode) sibling() *RBNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (tree *RBTree) rotateLeft(node *RBNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *RBTree) rotateRight(node *RBNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *RBTree) replaceNode(old *RBNode, new *RBNode) {
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

func (tree *RBTree) insertCase1(node *RBNode) {
	if node.Parent == nil {
		node.color = rbBlack
	} else {
		tree.insertCase2(node)
	}
}

func (tree *RBTree) insertCase2(node *RBNode) {
	if nodeColor(node.Parent) == rbBlack {
		return
	}
	tree.insertCase3(node)
}

func (tree *RBTree) insertCase3(node *RBNode) {
	uncle := node.uncle()
	if nodeColor(uncle) == rbRed {
		node.Parent.color = rbBlack
		uncle.color = rbBlack
		node.grandparent().color = rbRed
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *RBTree) insertCase4(node *RBNode) {
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

func (tree *RBTree) insertCase5(node *RBNode) {
	node.Parent.color = rbBlack
	grandparent := node.grandparent()
	grandparent.color = rbRed
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *RBNode) maximumNode() *RBNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *RBTree) deleteCase1(node *RBNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *RBTree) deleteCase2(node *RBNode) {
	sibling := node.sibling()
	if nodeColor(sibling) == rbRed {
		node.Parent.color = rbRed
		sibling.color = rbBlack
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *RBTree) deleteCase3(node *RBNode) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == rbBlack &&
		nodeColor(sibling) == rbBlack &&
		nodeColor(sibling.Left) == rbBlack &&
		nodeColor(sibling.Right) == rbBlack {
		sibling.color = rbRed
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *RBTree) deleteCase4(node *RBNode) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == rbRed &&
		nodeColor(sibling) == rbBlack &&
		nodeColor(sibling.Left) == rbBlack &&
		nodeColor(sibling.Right) == rbBlack {
		sibling.color = rbRed
		node.Parent.color = rbBlack
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *RBTree) deleteCase5(node *RBNode) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == rbBlack &&
		nodeColor(sibling.Left) == rbRed &&
		nodeColor(sibling.Right) == rbBlack {
		sibling.color = rbRed
		sibling.Left.color = rbBlack
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == rbBlack &&
		nodeColor(sibling.Right) == rbRed &&
		nodeColor(sibling.Left) == rbBlack {
		sibling.color = rbRed
		sibling.Right.color = rbBlack
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *RBTree) deleteCase6(node *RBNode) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = rbBlack
	if node == node.Parent.Left && nodeColor(sibling.Right) == rbRed {
		sibling.Right.color = rbBlack
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == rbRed {
		sibling.Left.color = rbBlack
		tree.rotateRight(node.Parent)
	}
}

func nodeColor(node *RBNode) rbColor {
	if node == nil {
		return rbBlack
	}
	return node.color
}

type Iterator struct {
	tree     *RBTree
	node     *RBNode
	position _position
}

type _position byte

const (
	begin, between, end _position = 0, 1, 2
)

func (tree *RBTree) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, position: begin}
}

func (tree *RBTree) IteratorAt(node *RBNode) Iterator {
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

func (iterator *Iterator) Prev() bool {
	if iterator.position == begin {
		goto begin
	}
	if iterator.position == end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	for iterator.node.Parent != nil {
		node := iterator.node
		iterator.node = iterator.node.Parent
		if node == iterator.node.Right {
			goto between
		}
	}

begin:
	iterator.node = nil
	iterator.position = begin
	return false

between:
	iterator.position = between
	return true
}

func (iterator *Iterator) Value() Value {
	return iterator.node.Value
}

func (iterator *Iterator) Key() Key {
	return iterator.node.Key
}

func (iterator *Iterator) Node() *RBNode {
	return iterator.node
}

func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = begin
}

func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
}

func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}

func (iterator *Iterator) NextTo(f func(key Key, value interface{}) bool) bool {
	for iterator.Next() {
		key, value := iterator.Key(), iterator.Value()
		if f(key, value) {
			return true
		}
	}
	return false
}

func (iterator *Iterator) PrevTo(f func(key Key, value interface{}) bool) bool {
	for iterator.Prev() {
		key, value := iterator.Key(), iterator.Value()
		if f(key, value) {
			return true
		}
	}
	return false
}
