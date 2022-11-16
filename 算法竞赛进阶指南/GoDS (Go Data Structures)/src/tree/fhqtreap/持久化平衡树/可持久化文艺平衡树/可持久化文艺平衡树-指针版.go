// https://www.luogu.com.cn/problem/P5055
// P5055 【模板】可持久化文艺平衡树
// https://www.luogu.com.cn/problem/solution/P5055

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
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

// !注：如果用的是指针写法，必要时禁止 GC，能加速不少
func init() { debug.SetGCPercent(-1) }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	roots := make([]*Node, n+1)

	lastRes := 0
	var version, opt, pos, value, left, right int
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &version, &opt)
		roots[i] = roots[version]
		switch opt {
		case 1:
			fmt.Fscan(in, &pos, &value)
			pos ^= lastRes
			value ^= lastRes
			newNode := Insert(roots[i], pos, value)
			roots[i] = newNode
		case 2:
			fmt.Fscan(in, &pos)
			pos ^= lastRes
			newNode := Pop(roots[i], pos-1)
			roots[i] = newNode
		case 3:
			fmt.Fscan(in, &left, &right)
			left ^= lastRes
			right ^= lastRes
			newNode := Reverse(roots[i], left, right)
			roots[i] = newNode
		case 4:
			fmt.Fscan(in, &left, &right)
			left ^= lastRes
			right ^= lastRes
			lastRes = Query(roots[i], left-1, right)
			fmt.Fprintln(out, lastRes)
		}
	}

}

type Node struct {
	// !Raw value
	element int

	// !Data and lazy tag maintained by segment tree
	sum     int
	lazyAdd int

	left, right *Node
	size        int
	priority    uint
	isReversed  uint8
}

func newNode(value int) *Node {
	return &Node{
		element:  value,
		sum:      value,
		size:     1,
		priority: nextRand(),
	}
}

func copyNode(node *Node) *Node {
	nodeCopy := *node // 赋值浅拷贝
	return &nodeCopy
}

func pushUp(node *Node) {
	if node == nil {
		return
	}
	node.size = 1
	node.sum = node.element
	if node.left != nil {
		node.size += node.left.size
		node.sum += node.left.sum
	}
	if node.right != nil {
		node.size += node.right.size
		node.sum += node.right.sum
	}
}

func pushDown(node *Node) {
	if node == nil {
		return
	}

	if node.isReversed != 0 {
		node.left, node.right = node.right, node.left
		if node.left != nil {
			node.left = copyNode(node.left)
			node.left.isReversed ^= 1
		}
		if node.right != nil {
			node.right = copyNode(node.right)
			node.right.isReversed ^= 1
		}
		node.isReversed = 0
	}
}

func splitByRank(root *Node, k int, left **Node, right **Node) {
	if root == nil {
		*left = nil
		*right = nil
		return
	}

	pushDown(root)

	leftSize := 0
	if root.left != nil {
		leftSize = root.left.size
	}
	if leftSize >= k {
		*right = copyNode(root)
		splitByRank(root.left, k, left, &((*right).left))
		pushUp(*right)
	} else {
		*left = copyNode(root)
		splitByRank(root.right, k-leftSize-1, &((*left).right), right)
		pushUp(*left)
	}
}

func merge(left *Node, right *Node) *Node {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if left.priority <= right.priority {
		pushDown(right)
		right.left = merge(left, right.left)
		pushUp(right)
		return right
	} else {
		pushDown(left)
		left.right = merge(left.right, right)
		pushUp(left)
		return left
	}
}

// Insert value before index.
func Insert(root *Node, k int, value int) *Node {
	var left, right *Node
	splitByRank(root, k, &left, &right)
	return merge(merge(left, newNode(value)), right)
}

func Pop(root *Node, k int) *Node {
	k++
	var a, b, c *Node
	splitByRank(root, k, &b, &c)
	splitByRank(b, k-1, &a, &b)
	return merge(a, c)
}

// Reverse [start, stop) in place.
func Reverse(root *Node, l int, r int) *Node {
	var a, b, c *Node
	splitByRank(root, r, &a, &c)
	splitByRank(a, l-1, &a, &b)
	if b != nil {
		b.isReversed ^= 1
	}
	return merge(merge(a, b), c)
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func Query(root *Node, l int, r int) int {
	l++
	var a, b, c *Node
	splitByRank(root, r, &a, &c)
	splitByRank(a, l-1, &a, &b)
	res := 0
	if b != nil {
		res = b.sum
	}
	root = merge(a, merge(b, c))
	return res
}

var seed = uint(time.Now().UnixNano()/2 + 1)

func nextRand() uint {
	seed ^= seed << 13
	seed ^= seed >> 17
	seed ^= seed << 5
	return seed
}

// Return all elements in index order.
func InOrder(node *Node) []int {
	res := make([]int, 0, node.size)
	inOrder(node, &res)
	return res
}

func inOrder(root *Node, res *[]int) {
	if root == nil {
		return
	}
	pushDown(root)
	inOrder(root.left, res)
	*res = append(*res, root.element)
	inOrder(root.right, res)
}
