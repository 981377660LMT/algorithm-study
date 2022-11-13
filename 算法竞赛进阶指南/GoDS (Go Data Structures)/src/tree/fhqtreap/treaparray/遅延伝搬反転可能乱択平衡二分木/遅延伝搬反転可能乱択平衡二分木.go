// An arraylist impleted by FHQTreap.
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
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	treap := NewFHQTreap(nums)
	fmt.Println(treap.At(0), treap)
	treap.Reverse(0, 7)
	fmt.Println(treap.InOrder(), treap)
}

type Node struct {
	// !Raw value
	raw int

	// !Data and lazy tag maintained by segment tree
	sum     int
	lazyAdd int

	// FHQTreap inner attributes
	left, right int
	size        int
	priority    uint
	isReversed  uint8
}

// !op
func (t *FHQTreap) pushUp(root int) {
	node := t.nodes[root]
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
	node.sum = t.nodes[node.left].sum + t.nodes[node.right].sum + node.raw
}

// Reverse first and then push down the lazy tag.
func (t *FHQTreap) pushDown(root int) {
	node := t.nodes[root]

	if node.isReversed == 1 {
		t.toggle(node.left)
		t.toggle(node.right)
		node.isReversed = 0
	}

	// func (tree *LazySegTree) pushDown(root int) {
	// 	tree.allApply(2*root, tree.lazy[root])
	// 	tree.allApply(2*root+1, tree.lazy[root])
	// 	tree.lazy[root] = tree.lazyUnit()
	// }
	if node.lazyAdd != 0 {
		delta := node.lazyAdd
		t.propagate(node.left, delta)
		t.propagate(node.right, delta)
		node.lazyAdd = 0
	}
}

// !mapping + composition
func (t *FHQTreap) propagate(root int, delta int) {
	node := t.nodes[root]
	node.raw += delta // need to update raw value (differs from segment tree)
	node.sum += delta * node.size
	node.lazyAdd += delta
}

type FHQTreap struct {
	seed  uint
	root  int
	nodes []*Node
}

// Need to be modified according to the actual situation to implement a segment tree.
func NewFHQTreap(nums []int) *FHQTreap {
	treap := &FHQTreap{
		seed:  uint(time.Now().UnixNano()/2 + 1),
		nodes: make([]*Node, 0, max(len(nums), 16)),
	}

	dummy := &Node{size: 0, priority: treap.fastRand()} // 0 is dummy
	treap.nodes = append(treap.nodes, dummy)
	treap.root = treap.build(1, len(nums), nums)
	return treap
}

// Return the value at the k-th position (0-indexed).
func (t *FHQTreap) At(index int) int {
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
	res := &t.nodes[y].raw
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

// Append value to the end of the list.
func (t *FHQTreap) Append(value int) {
	t.Insert(t.Size(), value)
}

// Insert value before index.
func (t *FHQTreap) Insert(index int, value int) {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index-1, &x, &y)
	z = t.newNode(value)
	t.root = t.merge(t.merge(x, z), y)
}

// Remove and return item at index.
func (t *FHQTreap) Pop(index int) int {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	res := &t.nodes[y].raw
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

// Update [start, stop) with value (defaults to range add).
//  0 <= start <= stop <= n
//  !alias:Apply
func (t *FHQTreap) Update(start, stop int, delta int) {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t.propagate(y, delta)
	t.root = t.merge(t.merge(x, y), z)
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func (t *FHQTreap) Query(start, stop int) int {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	res := t.nodes[y].sum
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// Query [0, n) (defaults to range sum).
func (t *FHQTreap) QueryAll() int {
	return t.nodes[t.root].sum
}

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// Return all elements in index order.
func (t *FHQTreap) InOrder() []int {
	res := make([]int, 0, t.Size())
	t.inOrder(t.root, &res)
	return res
}

func (t *FHQTreap) inOrder(root int, res *[]int) {
	if root == 0 {
		return
	}
	t.pushDown(root) // !pushDown lazy tag
	t.inOrder(t.nodes[root].left, res)
	*res = append(*res, t.nodes[root].raw)
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
func (t *FHQTreap) newNode(data int) int {
	node := &Node{
		size:     1,
		priority: t.fastRand(),
		raw:      data,
		sum:      data,
	}
	t.nodes = append(t.nodes, node)
	return len(t.nodes) - 1
}

// Build a treap from a slice and return the root nodeId. O(n).
func (t *FHQTreap) build(left, right int, nums []int) int {
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
