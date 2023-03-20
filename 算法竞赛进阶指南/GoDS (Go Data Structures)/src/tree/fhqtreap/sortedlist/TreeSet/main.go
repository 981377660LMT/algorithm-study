package main

import (
	"fmt"
	"strings"
)

func main() {
	set := NewTreeSet(func(a, b Key) int {
		return a[0] - b[0]
	})
	for i := 10; i > 0; i-- {
		set.Add(Key{i, i})
	}
	fmt.Println(set)

	it := set.Iterator()
	it.Next()
	it.Next()

	// it := set.LowerBound(Key{5, 0})
	fmt.Println(it.Key())
	set.Discard(Key{5, 0}) // 迭代器失效
	fmt.Println(set, it.Key())
	// it = set.UpperBound(Key{15, 0})

	if it.Next() { // 不等于end,并向后移动
		fmt.Println(it.Key())
	}
	if it.Prev() { // 不等于begin,并向前移动
		fmt.Println(it.Key())
	}
	fmt.Println(it.Key(), set)
	set.Erase(it)
	fmt.Println(it.Key(), set)

	it2, ok := set.LowerBound(Key{3, 0})
	fmt.Println(it2.Key(), ok) // 3 true
	it2, ok = set.UpperBound(Key{3, 0})
	fmt.Println(it2.Key(), ok) // 4 true

	fmt.Println(it2.Key(), set)
	it2 = set.Erase(it2)
	fmt.Println(it2.Key(), set)
}

type Key = [2]int
type Value = struct{}

type TreeSet struct {
	tree *_Tree
}

var _itemExists = struct{}{}

func NewTreeSet(comparator func(a, b Key) int, keys ...Key) *TreeSet {
	set := &TreeSet{tree: _NewRedBlackTree(comparator)}
	if len(keys) > 0 {
		for _, value := range keys {
			set.Add(value)
		}
	}
	return set
}

func (set *TreeSet) Iterator() Iterator {
	return set.tree.Iterator()
}

// 返回删除元素的后继元素的迭代器，如果删除的是最后一个元素，则返回end()迭代器。
func (set *TreeSet) Erase(it Iterator) Iterator {
	node := it.node
	it.Next()
	set.tree.DiscardNode(node)
	return it
}

// 返回一个迭代器，指向键值>= key的第一个元素。
func (set *TreeSet) LowerBound(key Key) (Iterator, bool) {
	ceiling, ok := set.tree.Ceiling(key)
	if !ok {
		return set.tree.Iterator(), false
	}
	return set.tree.IteratorAt(ceiling), true
}

// 返回一个迭代器，指向键值> key的第一个元素。
func (set *TreeSet) UpperBound(key Key) (Iterator, bool) {
	higher, ok := set.tree.Higher(key)
	if !ok {
		it := set.tree.Iterator()
		it.End()
		return it, false
	}
	return set.tree.IteratorAt(higher), true
}

func (set *TreeSet) Add(key Key) {
	set.tree.Put(key, _itemExists)
}

func (set *TreeSet) Discard(key Key) bool {
	return set.tree.Discard(key)
}

func (set *TreeSet) Has(key Key) bool {
	if _, contains := set.tree.Get(key); !contains {
		return false
	}
	return true
}

func (set *TreeSet) ForEach(f func(key Key)) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Key())
	}
}

func (set *TreeSet) Empty() bool {
	return set.tree.Size() == 0
}

func (set *TreeSet) Size() int {
	return set.tree.Size()
}

func (set *TreeSet) Clear() {
	set.tree.Clear()
}

func (set *TreeSet) Keys() []Key {
	return set.tree.Keys()
}

// String returns a string representation of container
func (set *TreeSet) String() string {
	items := []string{}
	for _, v := range set.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("TreeSet{%s}", strings.Join(items, ", "))
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

func (tree *_Tree) Get(key Key) (value Value, found bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.Value, true
	}
	return
}

func (tree *_Tree) GetNode(key Key) *_RBNode {
	return tree.lookup(key)
}

func (tree *_Tree) Discard(key Key) bool {
	node := tree.lookup(key)
	return tree.DiscardNode(node)
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

func (tree *_Tree) Empty() bool {
	return tree.size == 0
}

func (tree *_Tree) Size() int {
	return tree.size
}

func (node *_RBNode) Size() int {
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

func (tree *_Tree) Keys() []Key {
	keys := make([]Key, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

func (tree *_Tree) Values() []Value {
	values := make([]Value, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
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

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *_Tree) Right() *_RBNode {
	var parent *_RBNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

func (tree *_Tree) Floor(key Key) (floor *_RBNode, found bool) {
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

func (tree *_Tree) Lower(key Key) (lower *_RBNode, found bool) {
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

func (tree *_Tree) Clear() {
	tree.Root = nil
	tree.size = 0
}

func (tree *_Tree) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *_RBNode) String() string {
	return fmt.Sprintf("%v", node.Key)
}

func output(node *_RBNode, prefix string, isTail bool, str *string) {
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

func (iterator *Iterator) Node() *_RBNode {
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
