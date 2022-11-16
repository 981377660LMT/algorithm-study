// !被特殊数据卡了
// https://leetcode.cn/problems/longest-increasing-subsequence-ii/submissions/

package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	nums := make([]int, 1e5)
	for i := range nums {
		if i <= 5e4 {
			nums[i] = i + 5e4
		} else {
			nums[i] = 5e4
		}
	}
	k := int(5e4)
	// nums := make([]int, 1e5)
	// for i := range nums {
	// 	nums[i] = i + 1
	// }
	// k := 1

	time1 := time.Now()
	fmt.Println(lengthOfLIS(nums, k))
	fmt.Println(time.Since(time1))
}

func lengthOfLIS(nums []int, k int) int {
	// !更新:单点更新,查询:区间最大值
	treap := NewFHQTreap(len(nums))
	treap.Build(make([]int, 1e5+1))

	for _, num := range nums {
		preMax := treap.Query(num-k, num)
		treap.Update(num, num+1, preMax+1)
	}
	return treap.QueryAll()
}

type Node struct {
	// !Raw value
	element int

	// !Data and lazy tag maintained by segment tree (defaults to range max)
	max  int
	lazy int

	// FHQTreap inner attributes
	left, right int
	size        int
	priority    uint64
	isReversed  uint8
}

// !op
func (t *FHQTreap) pushUp(root int) {
	if root == 0 {
		return
	}

	rootRef := &t.nodes[root]
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

	rootRef := &t.nodes[root]
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
	rootRef := &t.nodes[root]
	rootRef.element = max(rootRef.element, delta)
	rootRef.max = max(rootRef.max, delta)
	rootRef.lazy = max(rootRef.lazy, delta)
}

func (t *FHQTreap) toggle(root int) {
	t.nodes[root].left, t.nodes[root].right = t.nodes[root].right, t.nodes[root].left
	t.nodes[root].isReversed ^= 1
}

type FHQTreap struct {
	seed  uint64
	root  int
	nodes []Node
}

// Need to be modified according to the actual situation to implement a segment tree.
func NewFHQTreap(initCapacity int) *FHQTreap {
	treap := &FHQTreap{
		seed:  uint64(time.Now().UnixNano()/2 + 1),
		nodes: make([]Node, 0, initCapacity),
	}

	dummy := &Node{size: 0} // 0 is dummy
	treap.nodes = append(treap.nodes, *dummy)
	return treap
}

// 返回build后的根节点版本编号
func (t *FHQTreap) Build(nums []int) int {
	n := len(nums)
	nodes := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nodes = append(nodes, t.newNode(nums[i]))
	}

	stack := []int{}
	pre := make([]int, n)
	for i := 0; i < n; i++ {
		pre[i] = -1
	}

	for i := 0; i < n; i++ {
		last := -1
		for len(stack) > 0 && t.nodes[stack[len(stack)-1]].priority > t.nodes[nodes[i]].priority {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			pre[i] = stack[len(stack)-1]
		}
		if last != -1 {
			pre[last] = i
		}

		stack = append(stack, i)
	}

	root := -1
	for i := 0; i < n; i++ {
		if pre[i] != -1 {
			if i < pre[i] {
				t.nodes[nodes[pre[i]]].left = nodes[i]
			} else {
				t.nodes[nodes[pre[i]]].right = nodes[i]
			}
		} else {
			root = i
		}
	}

	t.root = nodes[root]
	t.build(nodes[root])
	return nodes[root]
}

func (t *FHQTreap) build(root int) {
	nodeRef := t.nodes[root]
	if nodeRef.left != 0 {
		t.build(nodeRef.left)
	}
	if nodeRef.right != 0 {
		t.build(nodeRef.right)
	}
	t.pushUp(root)
}

// Remove [start, stop) from list.
func (t *FHQTreap) Erase(start, stop int) {
	var x, y, z int
	start++ // dummy offset
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t.root = t.merge(x, z)
}

// Reverse [start, stop) in place.
func (t *FHQTreap) Reverse(start, stop int) {
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start, &x, &y)
	t.toggle(y)
	t.root = t.merge(t.merge(x, y), z)
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
	res := &t.nodes[y].element
	t.root = t.merge(t.merge(x, y), z)
	return *res
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
	res := &t.nodes[y].element
	t.root = t.merge(x, z)
	return *res
}

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// Update [start, stop) with value (defaults to range add).
//  0 <= start <= stop <= n
func (t *FHQTreap) Update(start, stop int, delta int) {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	t.propagate(y, delta)
	t.root = t.merge(t.merge(x, y), z)
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func (t *FHQTreap) Query(start, stop int) int {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	res := t.nodes[y].max
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// Query [0, n) (defaults to range sum).
func (t *FHQTreap) QueryAll() int {
	return t.nodes[t.root].max
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
func (t *FHQTreap) newNode(value int) int {
	node := &Node{
		element:  value,
		max:      value,
		size:     1,
		priority: t.nextRand(),
	}
	t.nodes = append(t.nodes, *node)
	return len(t.nodes) - 1
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

// https://nyaannyaan.github.io/library/rbst/rbst-base.hpp
func (t *FHQTreap) nextRand() uint64 {
	// t.seed *=
	t.seed ^= t.seed << 7
	t.seed ^= t.seed >> 9
	return t.seed
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
