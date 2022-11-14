// An effective arraylist implemented by FHQTreap.
//
// Author:
// https://github.com/981377660LMT/algorithm-study
//
// Reference:
// https://baobaobear.github.io/post/20191215-fhq-treap/
// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go

// 为什么这么慢？
// https://leetcode.cn/problems/longest-increasing-subsequence-ii/submissions/
// https://leetcode.cn/problems/longest-increasing-subsequence-ii/solution/f-by-zzhcjb-0s2r/

// https://leetcode.cn/problems/longest-increasing-subsequence-ii/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// pushUp pushDown 次数太多了 为什么会有2e8次
var pushUpCount int
var pushDownCount int
var newNodeCount int
var splitCount int
var mergeCount int

func main() {
	// 5
	// 5e4
	nums := make([]int, 1e5)
	for i := range nums {
		if i <= 5e4 {
			nums[i] = i + 5e4
		} else {
			nums[i] = 5e4
		}
	}

	time1 := time.Now()
	fmt.Println(lengthOfLIS(nums, int(5e4)))
	fmt.Println(time.Since(time1))
	fmt.Println(pushUpCount, pushDownCount, newNodeCount, splitCount, mergeCount)
}

func lengthOfLIS(nums []int, k int) int {
	initNums := make([]int, 1e5)
	// 更新:单点更新,查询:区间最大值
	treap := NewFHQTreap(initNums)

	// 是不是min写的有问题
	for _, num := range nums {
		preMax := treap.Query(max(0, num-k), num)
		treap.Update(num, num+1, preMax+1)
	}
	return treap.QueryAll()
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
	priority    uint
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) pushUp(root int) {
	pushUpCount++
	node := t.nodes[root]
	// If left or right is 0(dummy), it will update with monoid.
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
	node.max = max(max(t.nodes[node.left].max, t.nodes[node.right].max), node.element)
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) pushDown(root int) {
	pushDownCount++
	node := t.nodes[root]
	if node.lazy != 0 {
		delta := node.lazy
		// !Not dummy
		if node.left != 0 {
			t.propagate(node.left, delta)
		}
		if node.right != 0 {
			t.propagate(node.right, delta)
		}
		node.lazy = 0
	}
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) propagate(root int, delta int) {
	node := t.nodes[root]
	node.element = max(node.element, delta)
	node.max = max(node.max, delta)
	node.lazy = max(node.lazy, delta)
}

type FHQTreap struct {
	seed  uint
	root  int
	nodes []*Node // Use pointer to avoid copying
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

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// // Return the value at the k-th position (0-indexed).
// func (t *FHQTreap) At(index int) int {
// 	n := t.Size()
// 	if index < 0 {
// 		index += n
// 	}

// 	if index < 0 || index >= n {
// 		panic(fmt.Sprintf("index %d out of range [0,%d]", index, n-1))
// 	}

// 	index += 1 // dummy offset
// 	var x, y, z int
// 	t.splitByRank(t.root, index, &y, &z)
// 	t.splitByRank(y, index-1, &x, &y)
// 	res := &t.nodes[y].element
// 	t.root = t.merge(t.merge(x, y), z)
// 	return *res
// }

// Update [start, stop) with value (defaults to range add).
//  0 <= start <= stop <= n
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
	res := t.nodes[y].max
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// Query [0, n) (defaults to range sum).
func (t *FHQTreap) QueryAll() int {
	return t.nodes[t.root].max
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (t *FHQTreap) splitByRank(root, k int, x, y *int) {

	if root == 0 {
		*x, *y = 0, 0
		return
	}
	splitCount++
	t.pushDown(root)

	if k <= t.nodes[t.nodes[root].left].size {
		// fmt.Println(k, t.nodes[t.nodes[root].left].size)
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
	mergeCount++
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
	newNodeCount++
	node := &Node{
		size:     1,
		priority: uint(rand.Uint64()),
		element:  data,
		max:      data,
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

// func (t *FHQTreap) String() string {
// 	sb := []string{"TreapArray{"}
// 	values := []string{}
// 	for i := 0; i < t.Size(); i++ {
// 		values = append(values, fmt.Sprintf("%d", t.At(i)))
// 	}
// 	sb = append(sb, strings.Join(values, ","), "}")
// 	return strings.Join(sb, "")

// }

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
