// multiset中允许存储重复元素的实现方式是在红黑树的节点数上增加`个数`属性对每种元素计数。
// 当多个元素的值相等时，它们在红黑树中被视为相等的节点。

// 在红黑树中插入、删除和查找节点时，需要对相等的节点的个数进行正确的处理。
// 当插入一个新节点时，如果它的键值已经存在于红黑树中，则需要将对应节点的计数器加一；
// 否则，需要创建一个新节点并将其插入到合适的位置。

// 当删除一个节点时，如果该节点的计数器大于1，则只需要将其计数器减一；
// 否则，需要将该节点从红黑树中删除。
// 在查找节点时，需要遍历红黑树并递归比较节点的键值，
// 直到找到对应节点或到达叶子节点为止。当找到相等的节点时，需要返回该节点的计数器。

// 注意cpp里的迭代器:
// !Begin指向第一个元素,
// !End指向最后一个元素的下一个位置,
// 这里的迭代器设计为:
// !Begin指向第一个元素的前一个位置, First指向第一个元素;
// !Last指向最后一个元素, End指向最后一个元素的下一个位置。

// test by:
// https://judge.yosupo.jp/submission/130489
// https://leetcode.cn/problems/stock-price-fluctuation/
// https://leetcode.cn/problems/design-a-leaderboard
// https://atcoder.jp/contests/abc241/tasks/abc241_d
// !和cpp不同，这里删除元素会使得所有迭代器失效

package main

import (
	"fmt"
	"strings"
)

func debugDiscardNode() {
	nums := []int{1, 2, 3}
	sl := NewMultiSet(func(a, b Key) int { return a - b })
	for _, num := range nums {
		sl.Add(num)
	}
	it1 := sl.LowerBound(1)
	it2 := sl.LowerBound(2)
	it2 = sl.Erase(it2) // 此时it1会失效
	// it1.Next()
	fmt.Println(it1.Key(), it2.Key())
}

func main() {
	debugDiscardNode()
}

type Key = int
type Value = struct{}

var _itemExists = struct{}{}

type MultiSet struct {
	tree *_Tree
}

func NewMultiSet(comparator func(a, b Key) int) *MultiSet {
	set := &MultiSet{tree: _NewRedBlackTree(comparator)}
	return set
}

func (set *MultiSet) Iterator() Iterator {
	return set.tree.Iterator()
}

// 返回删除元素的后继元素的迭代器，如果删除的是最后一个元素，则返回end()迭代器。
func (set *MultiSet) Erase(it Iterator) Iterator {
	node := it.node
	it.Next()
	set.tree.DiscardNode(node)
	return it
}

func (set *MultiSet) Insert(key Key) Iterator {
	node := set.tree.Put(key, _itemExists)
	return set.tree.IteratorAt(node)
}

// 返回一个迭代器，指向键值>= key的第一个元素。
//  `find` in C++ Multiset.
func (set *MultiSet) LowerBound(key Key) Iterator {
	ceiling, ok := set.tree.Ceiling(key)
	if !ok {
		it := set.tree.Iterator()
		it.End()
		return it
	}
	return set.tree.IteratorAt(ceiling)
}

// 返回一个迭代器，指向键值> key的第一个元素。
func (set *MultiSet) UpperBound(key Key) Iterator {
	higher, ok := set.tree.Higher(key)
	if !ok {
		it := set.tree.Iterator()
		it.End()
		return it
	}
	return set.tree.IteratorAt(higher)
}

func (set *MultiSet) Add(key Key) {
	set.tree.Put(key, _itemExists)
}

func (set *MultiSet) Discard(key Key) bool {
	return set.tree.Discard(key)
}

func (set *MultiSet) Has(key Key) bool {
	if _, contains := set.tree.Get(key); !contains {
		return false
	}
	return true
}

func (set *MultiSet) ForEach(f func(key Key)) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Key())
	}
}

func (set *MultiSet) Empty() bool {
	return set.tree.Size() == 0
}

func (set *MultiSet) Size() int {
	return set.tree.Size()
}

func (set *MultiSet) Clear() {
	set.tree.Clear()
}

func (set *MultiSet) Keys() []Key {
	return set.tree.Keys()
}

func (set *MultiSet) Min() (key Key, ok bool) {
	if node := set.tree.Left(); node != nil {
		return node.Key, true
	}
	return
}

func (set *MultiSet) Max() (key Key, ok bool) {
	if node := set.tree.Right(); node != nil {
		return node.Key, true
	}
	return
}

func (set *MultiSet) Ceiling(key Key) (ceiling Key, ok bool) {
	if node, contains := set.tree.Ceiling(key); contains {
		return node.Key, true
	}
	return
}

func (set *MultiSet) Floor(key Key) (floor Key, ok bool) {
	if node, contains := set.tree.Floor(key); contains {
		return node.Key, true
	}
	return
}

func (set *MultiSet) Higher(key Key) (higher Key, ok bool) {
	if node, contains := set.tree.Higher(key); contains {
		return node.Key, true
	}
	return
}

func (set *MultiSet) Lower(key Key) (lower Key, ok bool) {
	if node, contains := set.tree.Lower(key); contains {
		return node.Key, true
	}
	return
}

// String returns a string representation of container
func (set *MultiSet) String() string {
	items := []string{}
	for _, v := range set.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("MultiSet{%s}", strings.Join(items, ", "))
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
	Count  int // for multiset
	color  _color
	Left   *_RBNode
	Right  *_RBNode
	Parent *_RBNode
}

// NewWith instantiates a red-black tree with the custom comparator.
func _NewRedBlackTree(comparator func(a, b Key) int) *_Tree {
	return &_Tree{Comparator: comparator}
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

func (tree *_Tree) DiscardNode(node *_RBNode) bool {
	if node == nil {
		return false
	}

	// for multiset
	if node.Count > 1 {
		node.Count--
		tree.size--
		return true
	}

	var child *_RBNode
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node.Value = pred.Value
		node.Count = pred.Count // for multiset
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

func (tree *_Tree) Discard(key Key) bool {
	node := tree.lookup(key)
	if node == nil {
		return false
	}
	return tree.DiscardNode(node)
}

func (node *_RBNode) Size() int {
	if node == nil {
		return 0
	}
	size := node.Count
	if node.Left != nil {
		size += node.Left.Size()
	}
	if node.Right != nil {
		size += node.Right.Size()
	}
	return size
}

func (tree *_Tree) Put(key Key, value Value) *_RBNode {
	var insertedNode *_RBNode
	if tree.Root == nil {
		tree.Root = &_RBNode{Key: key, Value: value, color: _red, Count: 1}
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

				// for multiset
				node.Count++
				tree.size++
				return node
			case compare < 0:
				if node.Left == nil {
					node.Left = &_RBNode{Key: key, Value: value, color: _red, Count: 1}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &_RBNode{Key: key, Value: value, color: _red, Count: 1}
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
	return insertedNode
}

func (tree *_Tree) Keys() []Key {
	keys := make([]Key, tree.size)
	it := tree.Iterator()
	pos := 0
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
		pos++
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

func (tree *_Tree) Empty() bool {
	return tree.size == 0
}

func (tree *_Tree) Size() int {
	return tree.size
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
	return fmt.Sprintf("%v - %v", node.Key, node.Count)
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
	iterator.state = _begin
}

func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.state = _end
}

func (iterator *Iterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

func (iterator *Iterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}

func (iterator *Iterator) IsBegin() bool {
	return iterator.state == _begin
}

func (iterator *Iterator) IsEnd() bool {
	return iterator.state == _end
}

func (iterator *Iterator) IsFirst() bool {
	if iterator.state != _between {
		return false
	}
	_, hasLower := iterator.tree.Lower(iterator.node.Key)
	return !hasLower
}

func (iterator *Iterator) IsLast() bool {
	if iterator.state != _between {
		return false
	}
	_, hasHigher := iterator.tree.Higher(iterator.node.Key)
	return !hasHigher
}

type Iterator struct {
	pos   int // 位于每个结点的第几个元素上, 从0开始. 用于处理 Multiset 重复元素的情况.
	tree  *_Tree
	node  *_RBNode
	state _state
}

func (iterator Iterator) String() string {
	return fmt.Sprintf("Iterator{pos: %d, node: %v, state: %v}", iterator.pos, iterator.node, iterator.state)
}

type _state byte

const (
	_begin, _between, _end _state = 0, 1, 2
)

func (tree *_Tree) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, state: _begin}
}

// 定位到node的第一个元素处.
func (tree *_Tree) IteratorAt(node *_RBNode) Iterator {
	return Iterator{tree: tree, node: node, state: _between}
}

func (iterator *Iterator) Next() bool {
	if iterator.state == _end {
		goto end
	}
	if iterator.state == _begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}

	// for multiset
	if iterator.pos+1 < iterator.node.Count {
		iterator.pos++
		return true
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
	iterator.state = _end
	return false

between:
	iterator.pos = 0
	iterator.state = _between
	return true
}

func (iterator *Iterator) Prev() bool {
	if iterator.state == _begin {
		goto begin
	}

	if iterator.state == _end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	// for multiset
	if iterator.pos-1 >= 0 {
		iterator.pos--
		return true
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
	iterator.state = _begin
	return false

between:
	iterator.pos = iterator.node.Count - 1
	iterator.state = _between
	return true
}
