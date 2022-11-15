package main

import (
	"fmt"
	"time"
)

func main() {
	tree := NewFHQTreap(100000)
	fmt.Println(tree.InOrder())
	tree.Build([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(tree.InOrder())
	// api
	fmt.Println(tree.Size(), tree.Query(1, 1))
	tree.Update(1, 2, 1)
	fmt.Println(tree.Size(), tree.Query(1, 1))
	tree.Reverse(0, 4)
	fmt.Println(tree.InOrder())
	tree.Erase(1, 3)

}

type Node struct {
	// !Raw value
	element int

	// !Data and lazy tag maintained by segment tree (defaults to range sum)
	sum     int
	lazyAdd int

	// FHQTreap inner attributes
	left, right int
	size        int
	priority    uint64
	isReversed  uint8
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) pushUp(root int) {
	rootRef := t.nodes[root]
	rootRef.size = 1
	rootRef.sum = rootRef.element
	if rootRef.left != 0 {
		rootRef.size += t.nodes[rootRef.left].size
		rootRef.sum += t.nodes[rootRef.left].sum
	}
	if rootRef.right != 0 {
		rootRef.size += t.nodes[rootRef.right].size
		rootRef.sum += t.nodes[rootRef.right].sum
	}
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) pushDown(root int) {

	node := t.nodes[root]

	if node.isReversed == 1 {
		t.toggle(node.left)
		t.toggle(node.right)
		node.isReversed = 0
	}

	if node.lazyAdd != 0 {
		delta := node.lazyAdd
		// !Not dummy node
		if node.left != 0 {
			t.propagate(node.left, delta)
		}
		if node.right != 0 {
			t.propagate(node.right, delta)
		}
		node.lazyAdd = 0
	}
}

// !Segment tree function.Need to be modified according to actual situation.
func (t *FHQTreap) propagate(root int, delta int) {
	rootRef := t.nodes[root]
	rootRef.element += delta // need to update raw value (differs from segment tree)
	rootRef.sum += delta * rootRef.size
	rootRef.lazyAdd += delta
}

func (t *FHQTreap) toggle(root int) {
	t.nodes[root].left, t.nodes[root].right = t.nodes[root].right, t.nodes[root].left
	t.nodes[root].isReversed ^= 1
}

type FHQTreap struct {
	seed  uint64
	root  int
	nodes []*Node
}

// Need to be modified according to the actual situation to implement a segment tree.
func NewFHQTreap(initCapacity int) *FHQTreap {
	treap := &FHQTreap{
		seed:  uint64(time.Now().UnixNano()/2 + 1),
		nodes: make([]*Node, 0, initCapacity),
	}

	dummy := &Node{size: 0} // 0 is dummy
	treap.nodes = append(treap.nodes, dummy)
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
	res := t.nodes[y].sum
	t.root = t.merge(t.merge(x, y), z) // !顺序影响速度吗
	return res
}

// Query [0, n) (defaults to range sum).
func (t *FHQTreap) QueryAll() int {
	return t.nodes[t.root].sum
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
		sum:      value,
		size:     1,
		priority: t.nextRand(),
	}
	t.nodes = append(t.nodes, node)
	return len(t.nodes) - 1
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
// // [0, 2^64 - 1)
// u64 rng() {
//   static u64 _x =
//       u64(chrono::duration_cast<chrono::nanoseconds>(
//               chrono::high_resolution_clock::now().time_since_epoch())
//               .count()) *
//       10150724397891781847ULL;
//   _x ^= _x << 7;
//   return _x ^= _x >> 9;
// }
func (t *FHQTreap) nextRand() uint64 {
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
