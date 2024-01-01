// 基于 VectorTrie 实现的 PersistentList/PersistentVector.
// https://github.com/shanzi/algo-ds/tree/master/vector_trie
// !关注tail优化和transient

package main

import (
	"fmt"
	"time"
)

func main() {
	demo()
}

func demo() {
	list0 := NewPersistentList()
	list1 := list0.Push(1)
	list2 := list1.Push(2)
	list3, _ := list2.Set(0, 3)
	list4, _ := list3.Pop()
	fmt.Println(list0, list1, list2, list3, list4)

	time1 := time.Now()
	list5 := NewPersistentList().WithMutations(func(t ITransientList) {
		for i := 0; i < 1e6; i++ {
			t.Push(i)
		}
	})
	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 18ms
	_ = list5
	// fmt.Println(list5)
}

const (
	SHIFT     = 5
	NODE_SIZE = (1 << SHIFT)
	MASK      = NODE_SIZE - 1
)

var globalOwnerId int32

func nextId() int32 {
	globalOwnerId++
	return globalOwnerId
	// return atomic.AddInt32(&id, 1)  // concurrent safe
}

type IPersistentList interface {
	Get(index int) (interface{}, bool)
	Set(index int, value interface{}) (IPersistentList, bool)
	Push(value interface{}) IPersistentList
	Pop() (IPersistentList, interface{})
	Len() int

	// immutable.js style api
	AsMutable() ITransientList
	WithMutations(f func(ITransientList)) IPersistentList
}

// The difference is that Transient has an ID while List does not (0).
type PersistentList TransitentList

var EMPTY_PERSISTENT_LIST = &PersistentList{}

func NewPersistentList() IPersistentList {
	return EMPTY_PERSISTENT_LIST
}

func (head *PersistentList) Get(index int) (interface{}, bool) {
	return (*TransitentList)(head).Get(index)
}

func (head *PersistentList) Set(index int, value interface{}) (IPersistentList, bool) {
	t := head.AsMutable()
	if t.Set(index, value) {
		return t.AsImmutable(), true
	} else {
		return head, false
	}
}

func (head *PersistentList) Push(value interface{}) IPersistentList {
	t := head.AsMutable()
	t.Push(value)
	return t.AsImmutable()
}

func (head *PersistentList) Pop() (IPersistentList, interface{}) {
	if head.len == 1 {
		value, _ := head.Get(0)
		return EMPTY_PERSISTENT_LIST, value
	} else {
		t := head.AsMutable()
		value := t.Pop()
		return t.AsImmutable(), value
	}
}

func (head *PersistentList) AsMutable() ITransientList {
	id := nextId()
	return &TransitentList{id, head.len, head.level, head.offset, head.root, head.tail}
}

func (head *PersistentList) WithMutations(f func(ITransientList)) IPersistentList {
	t := head.AsMutable()
	f(t)
	return t.AsImmutable()
}

func (head *PersistentList) Len() int {
	return int(head.len)
}

func (head *PersistentList) String() string {
	res := make([]string, head.len)
	for i := 0; i < int(head.len); i++ {
		value, _ := head.Get(i)
		res[i] = fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf("%v", res)
}

type ITransientList interface {
	Get(n int) (interface{}, bool)
	Set(n int, value interface{}) bool
	Push(value interface{})
	Pop() interface{}
	Len() int

	// !dont modify previous transient list after AsImmutable
	AsImmutable() IPersistentList
}

type TransitentList struct {
	ownerId int32
	len     int32
	level   int8
	offset  int32 // 列表中在 tail 节点之前的节点当中储存的元素的数量，同时也是 tail 节点中下标0的元素在整个 List 当中的 Index
	root    *trieNode
	tail    *trieNode
}

func (head *TransitentList) Get(n int) (interface{}, bool) {
	n32 := int32(n)
	if n32 < 0 || n32 >= head.len {
		return nil, false
	}

	if n32 >= head.offset {
		return head.tail.getChildValue(n32 - head.offset), true
	}

	root := head.root
	for lv := head.level - 1; ; lv-- {
		index := (n32 >> (lv * SHIFT)) & MASK
		if lv <= 0 {
			// Arrived at leaves node, return value
			return root.getChildValue(index), true
		} else {
			// Update root node
			root = root.getChildNode(index)
		}
	}
}

func (head *TransitentList) Set(n int, value interface{}) bool {
	n32 := int32(n)
	if n32 < 0 || n32 >= head.len {
		panic("Index out of bound")
	}

	ok := false
	if n32 >= head.offset {
		head.tail, ok = setTail(head.ownerId, head.tail, n32-head.offset, value)
	} else {
		head.root, ok = setInNode(head.ownerId, head.root, n32, head.level, value)
	}
	return ok
}

func (head *TransitentList) Push(value interface{}) {
	if head.len-head.offset < NODE_SIZE {
		// Tail node has free space
		head.tail, _ = setTail(head.ownerId, head.tail, head.len-head.offset, value)
	} else {
		// Tail node is full
		n := head.offset
		lv := head.level
		root := head.root

		// Increase the depth of tree while the capacity is not enough
		for lv == 0 || (n>>(lv*SHIFT)) > 0 {
			parent := newTrieNode(head.ownerId)
			parent.children[0] = root
			root = parent
			lv++
		}

		head.root = putTail(head.ownerId, root, head.tail, n, lv)
		head.tail, _ = setTail(head.ownerId, nil, 0, value)

		head.level = lv
		head.offset += NODE_SIZE
	}
	head.len++
}

func (head *TransitentList) Pop() interface{} {
	if head.len == 0 {
		panic("Remove from empty list")
	}

	value := head.tail.getChildValue(head.len - head.offset - 1)
	head.tail, _ = setTail(head.ownerId, head.tail, head.len-head.offset-1, nil) // clear reference to release memory

	head.len--

	if head.len == 0 {
		head.level = 0
		head.offset = 0
		head.root = nil
		head.tail = nil
	} else {
		if head.len <= head.offset {
			// tail is empty, retrieve new tail from root
			head.root, head.tail = getTail(head.ownerId, head.root, head.len-1, head.level)
			head.offset -= NODE_SIZE
		}

		// Reduce the depth of tree if root only have one child
		n := head.offset - 1
		lv := head.level
		root := head.root

		for lv > 1 && (n>>((lv-1)*SHIFT)) == 0 {
			// if root only have one child, it must be zero
			root = root.getChildNode(0)
			lv--
		}

		head.root = root
		head.level = lv
	}

	return value
}

func (head *TransitentList) AsImmutable() IPersistentList {
	perisitHead := (*PersistentList)(head)
	perisitHead.ownerId = 0
	return perisitHead
}

func (head *TransitentList) Len() int {
	return int(head.len)
}

func setInNode(id int32, root *trieNode, n int32, level int8, value interface{}) (*trieNode, bool) {
	index := (n >> ((level - 1) * SHIFT)) & MASK

	if level == 1 {
		if root.getChildValue(index) == value {
			return root, false
		}
		return root.setChildValue(id, index, value), true
	} else {
		child := root.getChildNode(index)
		newChild, wasAltered := setInNode(id, child, n, level-1, value)
		if wasAltered {
			return root.setChildValue(id, index, newChild), true
		}
		return root, false
	}
}

func setTail(id int32, tail *trieNode, index int32, value interface{}) (*trieNode, bool) {
	if tail == nil {
		tail = newTrieNode(id)
	}

	if tail.getChildValue(index) == value {
		return tail, false
	}
	return tail.setChildValue(id, index, value), true
}

func getTail(id int32, root *trieNode, n int32, level int8) (*trieNode, *trieNode) {
	index := int32((n >> ((level - 1) * SHIFT)) & MASK)

	if level == 1 {
		return nil, root
	} else {
		child, tail := getTail(id, root.getChildNode(index), n, level-1)

		if index == 0 && child == nil {
			// The first element has been removed, which means current node
			// becomes empty, remove current node by returning nil
			return nil, tail
		} else {
			// Current node is not empty
			root = root.setChildValue(id, index, child)
			return root, tail
		}
	}
}

// find the tail node and put it into the correct position
func putTail(id int32, root *trieNode, tail *trieNode, n int32, level int8) *trieNode {
	index := int32((n >> ((level - 1) * SHIFT)) & MASK)

	if level == 1 {
		return tail
	} else {
		if root == nil {
			root = newTrieNode(id)
		}
		return root.setChildValue(id, index, putTail(id, root.getChildNode(index), tail, n, level-1))
	}
}

type trieNode struct {
	id       int32
	children []interface{}
}

func newTrieNode(id int32) *trieNode {
	return &trieNode{id, make([]interface{}, NODE_SIZE)}
}

func (node *trieNode) getChildNode(index int32) *trieNode {
	if child := node.children[index]; child != nil {
		return child.(*trieNode)
	} else {
		return nil
	}
}

func (node *trieNode) getChildValue(index int32) interface{} {
	return node.children[index]
}

func (node *trieNode) setChildValue(id int32, index int32, value interface{}) *trieNode {
	if node.id == id {
		node.children[index] = value // transient
		return node
	} else {
		newNode := node.clone(id)
		newNode.children[index] = value
		return newNode
	}
}

func (node *trieNode) clone(id int32) *trieNode {
	children := make([]interface{}, NODE_SIZE)
	copy(children, node.children)
	return &trieNode{id, children}
}
