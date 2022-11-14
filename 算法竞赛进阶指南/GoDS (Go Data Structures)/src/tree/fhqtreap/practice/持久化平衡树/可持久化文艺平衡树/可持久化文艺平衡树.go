// https://www.luogu.com.cn/problem/P5055
// P5055 【模板】可持久化文艺平衡树

package main

import (
	"fmt"
	"time"
)

// 您需要写一种数据结构，来维护一个序列，
// 其中需要提供以下操作（对于各个以往的历史版本）：

// 1.在第 p 个数后插入数 x 。
// 2.删除第 p 个数。
// 3.翻转区间 [l,r]，例如原序列是 {5,4,3,2,1}，翻转区间 [2,4]后，
// 结果是 {5,2,3,4,1}。
// 4.查询区间 [l,r] 中所有数的和。
// 和原本平衡树不同的一点是，每一次的任何操作都是基于某一个历史版本，
// 同时生成一个新的版本（操作 4 即保持原版本无变化），
// 新版本即编号为此次操作的序号。

// 本题强制在线。
// 第一行包含一个整数n，表示操作的总数。
// 接下来n行，每行前两个整数vi, opt，vi表示基于的过去版本号(0<vi<i) , opt,表示操作的序号(1≤opti≤ 4)。
// 若opti= 1，则接下来两个整数pi,xi，表示操作为在第pi个数后插入数x。
// 若opti= 2，则接下来一个整数pi，表示操作为删除第pi个数。
// 若opti= 3，则接下来两个整数li,ri，表示操作为翻转区间[li, ri]。
// 若opt,= 4，则接下来两个整数li,ri，表示操作为查询区间[li, ri]的和。
// 强制在线规则:
// 令当前操作之前的最后一次4操作的答案为lastRes(如果之前没有4操作，则lastRes = 0)。
// 则此次操作的pi, xi或li, ri均按位异或上lastRes即可得到真实的pi,xi或li, ri。

// !哪里有问题
func main() {
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// lastRes := 0
	// var n int
	// fmt.Fscan(in, &n)
	// versions := make([]int, n+10)

	// var version, opt, pos, value, left, right int
	// // tree:= NewPersistentFHQTreap(func(a, b int) int {
	// tree := NewPersistentFHQTreap(n)
	// for i := 1; i <= n; i++ {
	// 	fmt.Fscan(in, &version, &opt)
	// 	switch opt {
	// 	case 1:
	// 		fmt.Fscan(in, &pos, &value)
	// 		pos ^= lastRes
	// 		value ^= lastRes
	// 		versions[i] = tree.Insert(&versions[version], pos, value)
	// 	case 2:
	// 		fmt.Fscan(in, &pos)
	// 		pos ^= lastRes
	// 		versions[i] = tree.Pop(&versions[version], pos)
	// 	case 3:
	// 		fmt.Fscan(in, &left, &right)
	// 		left ^= lastRes
	// 		right ^= lastRes
	// 		versions[i] = tree.Reverse(&versions[version], left, right)
	// 	case 4:
	// 		fmt.Fscan(in, &left, &right)
	// 		left ^= lastRes
	// 		right ^= lastRes
	// 		lastRes, root := tree.Query(&versions[version], left, right)
	// 		versions[i] = root
	// 		fmt.Fprintln(out, lastRes)
	// 	}
	// }

	tree := NewPersistentFHQTreap([]int{1, 2, 3, 4, 5})
	versions := make([]int, 100)
	fmt.Println(tree.Root)
	versions[0] = tree.Root
	fmt.Println(tree.Query(versions[0], 0, 3))
	versions[1] = tree.Insert(versions[0], 5, 6)

}

type Node struct {
	// !Raw value
	element int

	// !Data and lazy tag maintained by segment tree
	sum     int
	lazyAdd int

	left, right int
	size        int
	priority    uint
	isReversed  uint8
}

// 需要开 50倍长度的数组 (动态开点)
type PersistentFHQTreap struct {
	Root  int
	seed  uint
	nodes []Node
}

func NewPersistentFHQTreap(nums []int) *PersistentFHQTreap {
	treap := &PersistentFHQTreap{
		seed:  uint(time.Now().UnixNano()/2 + 1),
		nodes: make([]Node, 0, 50*len(nums)+10),
	}

	dummy := &Node{size: 0, priority: treap.fastRand()} // 0 is dummy
	treap.nodes = append(treap.nodes, *dummy)
	treap.Root = treap.build(1, len(nums), nums)
	return treap
}

func (t *PersistentFHQTreap) Insert(root int, index int, value int) int {
	var x, y, z int
	index += 1 // dummy offset
	t.splitByRank(root, index-1, &x, &y)
	z = t.newNode(value)
	t.Root = t.merge(t.merge(x, z), y)
	return t.Root
}

func (t *PersistentFHQTreap) Pop(root int, index int) int {
	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(root, index, &x, &z)
	t.splitByRank(x, index-1, &x, &y)
	t.Root = t.merge(x, z)
	return t.Root
}

// Reverse [start, stop) in place.
func (t *PersistentFHQTreap) Reverse(root int, start, stop int) int {
	start++
	var x, y, z int
	t.splitByRank(root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	t.toggle(y)
	t.Root = t.merge(t.merge(x, y), z)
	return t.Root
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func (t *PersistentFHQTreap) Query(root int, start, stop int) (res int) {
	start++
	var x, y, z int
	t.splitByRank(root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	res = t.nodes[y].sum
	t.Root = t.merge(t.merge(x, y), z)
	return
}

// Build a treap from a slice and return the root nodeId. O(n).
func (t *PersistentFHQTreap) build(left, right int, nums []int) int {
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

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (t *PersistentFHQTreap) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	t.pushDown(root)
	if k <= t.nodes[t.nodes[root].left].size {
		*y = t.copyNode(root)
		t.splitByRank(t.nodes[root].left, k, x, &t.nodes[root].left)
		t.pushUp(*y)
	} else {
		*x = t.copyNode(root)
		t.splitByRank(t.nodes[root].right, k-t.nodes[t.nodes[root].left].size-1, &t.nodes[root].right, y)
		t.pushUp(*x)
	}
}

// Make sure that the height of the resulting tree is at most O(log n).
// A random priority is introduced to determine who is the root after merge operation.
// If left subtree is smaller, merge right subtree with the right child of the left subtree.
// Otherwise, merge left subtree with the left child of the right subtree.
func (t *PersistentFHQTreap) merge(x, y int) int {
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

func (t *PersistentFHQTreap) Len() int {
	return t.nodes[t.Root].size
}

func (t *PersistentFHQTreap) pushUp(root int) {
	node := &t.nodes[root]
	// If left or right is 0(dummy), it will update with monoid.
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
	node.sum = t.nodes[node.left].sum + t.nodes[node.right].sum + node.element
}

// Reverse first and then push down the lazy tag.
func (t *PersistentFHQTreap) pushDown(root int) {
	node := &t.nodes[root]
	if node.left != 0 {
		node.left = t.copyNode(node.left)
	}
	if node.right != 0 {
		node.right = t.copyNode(node.right)
	}

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

// !mapping + composition
func (t *PersistentFHQTreap) propagate(root int, delta int) {
	node := t.nodes[root]
	node.element += delta // need to update raw value (differs from segment tree)
	node.sum += delta * node.size
	node.lazyAdd += delta
}

func (t *PersistentFHQTreap) toggle(root int) {
	t.nodes[root].left, t.nodes[root].right = t.nodes[root].right, t.nodes[root].left
	t.nodes[root].isReversed ^= 1
}

func (t *PersistentFHQTreap) newNode(value int) int {
	node := &Node{
		size:     1,
		priority: t.fastRand(),
		element:  value,
		sum:      value,
	}
	t.nodes = append(t.nodes, *node)
	return len(t.nodes) - 1
}

func (t *PersistentFHQTreap) copyNode(root int) int {
	res := t.newNode(0)
	t.nodes[res] = t.nodes[root]
	return res
}

func (t *PersistentFHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}
