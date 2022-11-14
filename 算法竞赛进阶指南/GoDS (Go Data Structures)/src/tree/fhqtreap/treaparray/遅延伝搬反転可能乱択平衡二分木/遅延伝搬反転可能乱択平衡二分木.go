// An effective arraylist implemented by FHQTreap.
//
// 「強すぎてAtCoderRated出禁になった最強データ構造・平衡二分木のRBSTによる実装。」 —— nyaan
//
// Author:
// https://github.com/981377660LMT/algorithm-study
//
// Reference:
// https://baobaobear.github.io/post/20191215-fhq-treap/
// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go

package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	const INF int = 1e18
	nums := make([]Element, 1e5)
	for i := range nums {
		nums[i] = Element(i)
	}
	// 区间更新：加上一个数，区间查询：区间最小值
	treap := NewFHQTreap(nums, Operations{
		elementMonoid: func() Element {
			return INF
		},
		dataMonoid: func() Data {
			return INF
		},
		lazyMonoid: func() Lazy {
			return 0
		},
		op: func(a, b Data, e Element) Data {
			return min(min(a, b), e)
		},
		mappingElement: func(lazy Lazy, element Element) Element {
			return element + lazy
		},
		mappingData: func(lazy Lazy, data Data) Data {
			return data + lazy
		},
		composition: func(lazy1 Lazy, lazy2 Lazy) Lazy {
			return lazy1 + lazy2
		},
	})

	// 1e5 reverse
	time1 := time.Now()
	for i := 0; i < 2e4; i++ {
		treap.Reverse(30000, 60000)
		treap.Query(30000, 60000)
		treap.Insert(30000, 0)
		// treap.Erase(200, 300)
		treap.Append(0)
	}
	fmt.Println(time.Since(time1))

}

// !Type and functions to be implemented.
// type Element = interface{}
// type Data = interface{}
// type Lazy = interface{}
type Element = int
type Data = int
type Lazy = int
type Operations struct {
	elementMonoid  func() Element
	dataMonoid     func() Data
	lazyMonoid     func() Lazy
	op             func(leftData, rightData Data, element Element) Data
	mappingData    func(lazy Lazy, data Data) Data
	mappingElement func(lazy Lazy, element Element) Element
	composition    func(parentLazy Lazy, childLazy Lazy) Lazy
}

// !Template
//
func NewFHQTreap(nums []Element, operations Operations) *FHQTreap {
	treap := &FHQTreap{
		seed:       uint(time.Now().UnixNano()/2 + 1),
		nodes:      make([]*Node, 0, max(len(nums), 16)),
		Operations: operations,
	}

	// dummy node 0
	dummy := &Node{
		size: 0, priority: treap.fastRand(),
		element: treap.elementMonoid(), data: treap.dataMonoid(), lazy: treap.lazyMonoid(),
	}
	treap.nodes = append(treap.nodes, dummy)
	treap.root = treap.build(1, len(nums), nums)
	return treap
}

type FHQTreap struct {
	seed  uint
	root  int
	nodes []*Node

	// Segment-tree like operations
	Operations
}

type Node struct {
	// !Raw value
	element Element

	// !Data and lazy tag maintained by segment tree
	data Data
	lazy Lazy

	// FHQTreap inner attributes
	left, right int
	size        int
	priority    uint
	isReversed  uint8
}

// !op
func (t *FHQTreap) pushUp(root int) {
	node := t.nodes[root]
	// If left or right is 0(dummy), it will update with monoid.
	node.data = t.op(t.nodes[node.left].data, t.nodes[node.right].data, node.element)
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
}

// !Reverse first and then push down the lazy tag.
func (t *FHQTreap) pushDown(root int) {
	node := t.nodes[root]

	if node.isReversed == 1 {
		t.toggle(node.left)
		t.toggle(node.right)
		node.isReversed = 0
	}

	// !Not dummy node
	if node.left != 0 {
		t.propagate(node.left, node.lazy)
	}
	if node.right != 0 {
		t.propagate(node.right, node.lazy)
	}
	node.lazy = t.lazyMonoid()
}

// !mapping + composition
func (t *FHQTreap) propagate(root int, lazy Lazy) {
	node := t.nodes[root]
	node.element = t.mappingElement(lazy, node.element)
	node.data = t.mappingData(lazy, node.data)
	node.lazy = t.composition(lazy, node.lazy)
}

// Return the element at the k-th position (0-indexed).
func (t *FHQTreap) At(index int) Element {
	n := t.Size()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		panic(fmt.Sprintf("index %d out of range [0,%d]", index, n-1))
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	res := &t.nodes[y].element
	t.root = t.merge(t.merge(x, y), z)
	return *res
}

// Reverse [start, stop) in place.
func (t *FHQTreap) Reverse(start, stop int) {
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start, &x, &y)
	t.toggle(y)
	t.root = t.merge(t.merge(x, y), z)
}

// Append element to the end of the list.
func (t *FHQTreap) Append(element Element) {
	t.Insert(t.Size(), element)
}

// Insert element before index.
func (t *FHQTreap) Insert(index int, element Element) {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index-1, &x, &y)
	z = t.newNode(element)
	t.root = t.merge(t.merge(x, z), y)
}

// Remove and return element at index.
func (t *FHQTreap) Pop(index int) Element {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	res := &t.nodes[y].element
	t.root = t.merge(x, z)
	return *res
}

// Remove [start, stop) from list.
func (t *FHQTreap) Erase(start, stop int) {
	var x, y, z int
	start++ // dummy offset
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t.root = t.merge(x, z)
}

// Update [start, stop) with value .
//  0 <= start <= stop <= n
//  !alias:Apply
func (t *FHQTreap) Update(start, stop int, lazy Lazy) {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t.propagate(y, lazy)
	t.root = t.merge(t.merge(x, y), z)
}

// Query data in [start, stop).
//  0 <= start <= stop <= n
func (t *FHQTreap) Query(start, stop int) Data {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	res := t.nodes[y].data
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// Query data in [0, n).
func (t *FHQTreap) QueryAll() Data {
	return t.nodes[t.root].data
}

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// Return all elements in index order.
func (t *FHQTreap) InOrder() []Element {
	res := make([]Element, 0, t.Size())
	t.inOrder(t.root, &res)
	return res
}

func (t *FHQTreap) inOrder(root int, res *[]Element) {
	if root == 0 {
		return
	}
	t.pushDown(root) // !pushDown lazy tag
	t.inOrder(t.nodes[root].left, res)
	*res = append(*res, t.nodes[root].element)
	t.inOrder(t.nodes[root].right, res)
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (t *FHQTreap) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	t.pushDown(root)
	if k <= t.nodes[t.nodes[root].left].size {
		*y = root
		t.splitByRank(t.nodes[root].left, k, x, &t.nodes[root].left)
	} else {
		*x = root
		t.splitByRank(t.nodes[root].right, k-t.nodes[t.nodes[root].left].size-1, &t.nodes[root].right, y)
	}

	t.pushUp(root)
}

// Make sure that the height of the resulting tree is at most O(log n).
// A random priority is introduced to determine who is the root after merge operation.
// If left subtree is smaller, merge right subtree with the right child of the left subtree.
// Otherwise, merge left subtree with the left child of the right subtree.
func (t *FHQTreap) merge(x, y int) int {
	if x == 0 {
		return y
	}
	if y == 0 {
		return x
	}

	if t.nodes[x].priority < t.nodes[y].priority {
		t.pushDown(x)
		t.nodes[x].right = t.merge(t.nodes[x].right, y)
		t.pushUp(x)
		return x
	} else {
		t.pushDown(y)
		t.nodes[y].left = t.merge(x, t.nodes[y].left)
		t.pushUp(y)
		return y
	}
}

// Add a new node and return its nodeId.
func (t *FHQTreap) newNode(ele Element) int {
	node := &Node{
		size:     1,
		priority: t.fastRand(),
		element:  ele,
		data:     t.dataMonoid(),
		lazy:     t.lazyMonoid(),
	}
	t.nodes = append(t.nodes, node)
	return len(t.nodes) - 1
}

// Build a treap from a slice and return the root nodeId. O(n).
func (t *FHQTreap) build(left, right int, nums []Element) int {
	if left > right {
		return 0
	}
	mid := (left + right) >> 1
	newNode := t.newNode(nums[mid-1])
	t.nodes[newNode].left = t.build(left, mid-1, nums)
	t.nodes[newNode].right = t.build(mid+1, right, nums)
	t.pushUp(newNode)
	return newNode
}

func (t *FHQTreap) toggle(root int) {
	t.nodes[root].left, t.nodes[root].right = t.nodes[root].right, t.nodes[root].left
	t.nodes[root].isReversed ^= 1
}

func (t *FHQTreap) String() string {
	sb := []string{"TreapArray{"}
	values := []string{}
	for i := 0; i < t.Size(); i++ {
		values = append(values, fmt.Sprintf("%d", t.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")

}

// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go#L31
func (t *FHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
