// An effective arraylist implemented by FHQTreap.
//
// Author:
// https://github.com/981377660LMT/algorithm-study
//
// Reference:
// https://baobaobear.github.io/post/20191215-fhq-treap/
// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go

// !不用指针快很多
package main

import (
	"time"
)

func lengthOfLIS(nums []int, k int) int {
	// 更新:单点更新,查询:区间最大值
	treap := NewFHQTreap(len(nums))
	treap.Build(make([]int, 1e5+1))

	for _, num := range nums {
		preMax := treap.Query(num-k, num)
		treap.Update(num, num+1, preMax+1)
	}
	return treap.QueryAll()
}

// static uint64_t x_ = 88172645463325252ULL;
// return x_ ^= x_ << 7, x_ ^= x_ >> 9, x_ & 0xFFFFFFFFull;
// https://nyaannyaan.github.io/library/rbst/rbst-base.hpp
var seed uint64 = uint64(time.Now().UnixNano()/2 + 1)

func nextRand() uint64 {
	seed ^= seed << 7
	seed ^= seed >> 9
	return seed & 0xFFFFFFFF
}

type Node struct {
	// !Raw value
	element int

	// !Data and lazy tag maintained by segment tree (defaults to range sum)
	max  int
	lazy int

	// FHQTreap inner attributes
	left, right int
	size        int
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) pushUp(root int) {
	if root == 0 {
		return
	}

	rootRef := t.nodes[root]
	rootRef.size = 1
	rootRef.max = rootRef.element
	if rootRef.left != 0 {
		rootRef.size += t.nodes[rootRef.left].size
		rootRef.max = max(rootRef.max, t.nodes[rootRef.left].max)
	}

	if rootRef.right != 0 {
		rootRef.size += t.nodes[rootRef.right].size
		rootRef.max = max(rootRef.max, t.nodes[rootRef.right].max)
	}
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) pushDown(root int) {

	rootRef := t.nodes[root]
	if rootRef.lazy != 0 {
		delta := rootRef.lazy
		// !Not dummy
		if rootRef.left != 0 {
			t.propagate(rootRef.left, delta)
		}
		if rootRef.right != 0 {
			t.propagate(rootRef.right, delta)
		}
		rootRef.lazy = 0
	}
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) propagate(root int, delta int) {
	rootRef := t.nodes[root]
	rootRef.element = max(rootRef.element, delta)
	rootRef.max = max(rootRef.max, delta)
	rootRef.lazy = max(rootRef.lazy, delta)
}

type FHQTreap struct {
	seed  uint
	root  int
	nodes []*Node
}

// Need to be modified according to the actual situation to implement a segment tree.
func NewFHQTreap(initCapacity int) *FHQTreap {
	treap := &FHQTreap{
		seed:  uint(time.Now().UnixNano()/2 + 1),
		nodes: make([]*Node, 0, initCapacity),
	}

	dummy := &Node{size: 0} // 0 is dummy
	treap.nodes = append(treap.nodes, dummy)
	return treap
}

// 返回build后的根节点版本编号
func (t *FHQTreap) Build(nums []int) int {
	t.root = t.build(1, len(nums), nums)
	return t.root
}

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

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// Update [start, stop) with value (defaults to range add).
//  0 <= start <= stop <= n
func (t *FHQTreap) Update(start, stop int, delta int) {
	if start >= stop {
		return
	}
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	t.propagate(y, delta)
	t.root = t.merge(t.merge(x, y), z) // !顺序影响速度吗 (Update和Query是否需要不同地merge?)
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func (t *FHQTreap) Query(start, stop int) int {
	if start >= stop {
		return 0
	}
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	res := t.nodes[y].max
	t.root = t.merge(t.merge(x, y), z) // !顺序影响速度吗
	return res
}

// Query [0, n) (defaults to range sum).
func (t *FHQTreap) QueryAll() int {
	return t.nodes[t.root].max
}

func (t *FHQTreap) Insert(index int, value int) {
	var x, y int
	t.splitByRank(t.root, index, &x, &y)
	t.root = t.merge(t.merge(x, t.newNode(value)), y)
}

func (t *FHQTreap) Pop(index int) {
	var x, y, z int
	t.splitByRank(t.root, index, &x, &z)
	t.splitByRank(x, index-1, &x, &y)
	t.root = t.merge(x, z)
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
		t.pushUp(*y)
	} else {
		*x = root
		t.splitByRank(t.nodes[root].right, k-t.nodes[t.nodes[root].left].size-1, &t.nodes[root].right, y)
		t.pushUp(*x)
	}
}

// Make sure that the height of the resulting tree is at most O(log n).
// A random priority is introduced to determine who is the root after merge operation.
// If left subtree is smaller, merge right subtree with the right child of the left subtree.
// Otherwise, merge left subtree with the left child of the right subtree.
func (t *FHQTreap) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x | y
	}

	// https://nyaannyaan.github.io/library/rbst/rbst-base.hpp
	// if (int((rng() * (l->cnt + r->cnt)) >> 32) < l->cnt) {
	if int(nextRand()*(uint64(t.nodes[x].size)+uint64(t.nodes[y].size))>>32) < t.nodes[x].size {
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
func (t *FHQTreap) newNode(value int) int {
	node := &Node{
		element: value,
		max:     value,
		size:    1,
	}
	t.nodes = append(t.nodes, node)
	return len(t.nodes) - 1 // !返回新节点的编号(当前在nodes中的下标)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
