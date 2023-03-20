package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// func init() {
// 	debug.SetGCPercent(-1)
// }

func UnionOfInterval() {
	// https://atcoder.jp/contests/abc256/tasks/abc256_d
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	odt := NewIntervals(-INF)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		odt.Set(l, r, 1)
	}
	odt.EnumerateAll(func(l, r int, x int) {
		if x == 1 {
			fmt.Fprintln(out, l, r)
		}
	})
}

func main() {
	UnionOfInterval()
}

const INF int = 1e18

type Value = int

type Intervals struct {
	Len       int // 区间数
	Count     int // 区间元素个数之和
	noneValue Value
	mp        *TreeMap
}

func NewIntervals(noneValue Value) *Intervals {
	res := &Intervals{
		noneValue: noneValue,
		mp:        NewTreeMap(),
	}
	res.mp.Set(-INF, noneValue)
	res.mp.Set(INF, noneValue)
	return res
}

// 返回包含 x 的区间的信息.
func (odt *Intervals) Get(x int, erase bool) (start, end int, value Value) {
	iter, _ := odt.mp.UpperBound(x)
	end, _ = iter.Key(), iter.Value()
	iter.Prev()
	start, value = iter.Key(), iter.Value()
	if erase && value != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
		odt.mp.Set(start, odt.noneValue)
		odt.mergeAt(start)
		odt.mergeAt(end)
	}
	return
}

func (odt *Intervals) Set(start, end int, value Value) {
	odt.EnumerateRange(start, end, func(l, r int, x Value) {}, true)
	odt.mp.Set(start, value)
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}

	odt.mergeAt(start)
	odt.mergeAt(end)
}

func (odt *Intervals) EnumerateAll(f func(start, end int, value Value)) {
	odt.EnumerateRange(-INF, INF, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *Intervals) EnumerateRange(start, end int, f func(start, end int, value Value), erase bool) {
	if !(-INF <= start && start <= end && end <= INF) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue

	if !erase {
		it1, _ := odt.mp.UpperBound(start)
		it1.Prev()
		for {
			key1, val1 := it1.Key(), it1.Value()
			if key1 >= end {
				break
			}
			it1.Next()
			key2 := it1.Key()
			f(max(key1, start), min(key2, end), val1)
		}
		return
	}

	iter1, _ := odt.mp.UpperBound(start)
	iter1.Prev()
	if key1, val1 := iter1.Key(), iter1.Value(); key1 < start {
		odt.mp.Set(start, val1)
		if val1 != NONE {
			odt.Len++
		}
	}

	// 分割区间
	iter1, _ = odt.mp.LowerBound(end)
	if key1 := iter1.Key(); key1 > end {
		iter1.Prev()
		val2 := iter1.Value()
		odt.mp.Set(end, val2)
		if val2 != NONE {
			odt.Len++
		}
	}

	iter1, _ = odt.mp.LowerBound(start)
	for {
		key1, val1 := iter1.Key(), iter1.Value()
		if key1 >= end {
			break
		}
		iter1 = odt.mp.Erase(iter1)
		key2 := iter1.Key()
		f(key1, key2, val1)
		if val1 != NONE {
			odt.Len--
			odt.Count -= key2 - key1
		}
	}

	odt.mp.Set(start, NONE)
}

func (odt *Intervals) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int, value Value) {
		var v interface{} = value
		if value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
}

func (odt *Intervals) mergeAt(p int) {
	if p == -INF || p == INF {
		return
	}
	iter1, _ := odt.mp.LowerBound(p)
	val1 := iter1.Value()
	iter1.Prev()
	val2 := iter1.Value()
	if val1 == val2 {
		if val1 != odt.noneValue {
			odt.Len--
		}
		iter1.Next()
		odt.mp.Erase(iter1)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type TreeMap struct {
	tree *_Tree
}

type TreeMapIterator struct {
	iterator Iterator
}

func NewTreeMap() *TreeMap {
	return &TreeMap{tree: _NewRedBlackTree()}
}

func (m *TreeMap) Iterator() Iterator {
	return m.tree.Iterator()
}

func (m *TreeMap) Set(key int, value Value) {
	m.tree.Put(key, value)
}

func (m *TreeMap) Get(key int) (value Value, found bool) {
	return m.tree.Get(key)
}

func (m *TreeMap) Discard(key int) bool {
	return m.tree.Discard(key)
}

func (m *TreeMap) Size() int {
	return m.tree.Size()
}

func (m *TreeMap) Floor(key int) (foundkey int, foundvalue Value, ok bool) {
	node, found := m.tree.Floor(key)
	if found {
		return node.Key, node.Value, true
	}
	return
}

func (m *TreeMap) Ceiling(key int) (foundkey int, foundvalue Value, ok bool) {
	node, found := m.tree.Ceiling(key)
	if found {
		return node.Key, node.Value, true
	}
	return
}

func (m *TreeMap) Lower(key int) (foundkey int, foundvalue Value, ok bool) {
	node, found := m.tree.Lower(key)
	if found {
		return node.Key, node.Value, true
	}
	return
}

func (m *TreeMap) Higher(key int) (foundkey int, foundvalue Value, ok bool) {
	node, found := m.tree.Higher(key)
	if found {
		return node.Key, node.Value, true
	}
	return
}

// 返回删除元素的后继元素的迭代器，如果删除的是最后一个元素，则返回end()迭代器。
func (m *TreeMap) Erase(it Iterator) Iterator {
	node := it.node
	it.Next()
	m.tree.DiscardNode(node)
	return it
}

// 返回一个迭代器，指向键值>= key的第一个元素。
func (m *TreeMap) LowerBound(key int) (Iterator, bool) {
	lower, ok := m.tree.Ceiling(key)
	if !ok {
		return m.tree.Iterator(), false
	}
	return m.tree.IteratorAt(lower), true
}

// 返回一个迭代器，指向键值> key的第一个元素。
func (m *TreeMap) UpperBound(key int) (Iterator, bool) {
	upper, ok := m.tree.Higher(key)
	if !ok {
		it := m.tree.Iterator()
		it.End()
		return it, false
	}
	return m.tree.IteratorAt(upper), true
}

func (m *TreeMap) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}

// !region RedBlackTree
//

type _color bool

const (
	_black, _red _color = true, false
)

type _Tree struct {
	Root *_RBNode
	size int
}

type _RBNode struct {
	Key    int
	Value  Value
	color  _color
	Left   *_RBNode
	Right  *_RBNode
	Parent *_RBNode
}

func _NewRedBlackTree() *_Tree {
	return &_Tree{}
}

func (tree *_Tree) Put(key int, value Value) {
	var insertedNode *_RBNode
	if tree.Root == nil {
		tree.Root = &_RBNode{Key: key, Value: value, color: _red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		for loop {
			compare := key - node.Key
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

func (tree *_Tree) Get(key int) (value Value, found bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.Value, true
	}
	return
}

func (tree *_Tree) GetNode(key int) *_RBNode {
	return tree.lookup(key)
}

func (tree *_Tree) Discard(key int) bool {
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

func (tree *_Tree) Keys() []int {
	keys := make([]int, tree.size)
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

func (tree *_Tree) Floor(key int) (floor *_RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := key - node.Key
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

func (tree *_Tree) Lower(key int) (lower *_RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := key - node.Key
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

func (tree *_Tree) Ceiling(key int) (ceiling *_RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := key - node.Key
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

func (tree *_Tree) Higher(key int) (higher *_RBNode, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := key - node.Key
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

func (tree *_Tree) lookup(key int) *_RBNode {
	node := tree.Root
	for node != nil {
		compare := key - node.Key
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

func (iterator *Iterator) Key() int {
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

func (iterator *Iterator) NextTo(f func(key int, value Value) bool) bool {
	for iterator.Next() {
		key, value := iterator.Key(), iterator.Value()
		if f(key, value) {
			return true
		}
	}
	return false
}

func (iterator *Iterator) PrevTo(f func(key int, value Value) bool) bool {
	for iterator.Prev() {
		key, value := iterator.Key(), iterator.Value()
		if f(key, value) {
			return true
		}
	}
	return false
}
