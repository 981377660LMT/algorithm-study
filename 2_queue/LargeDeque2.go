// https://atcoder.jp/contests/abc413/tasks/abc413_c
//
// Insert(pos, v, c) -> 在位置 pos 前插入 c 个 v
// Erase(pos, count) -> 删除位置 pos 开始的 count 个元素，返回它们的和
// RangeSum(l, length) -> 区间求和 [l, l+length)
// Sum() -> 返回队列中所有元素的和
// At(k) -> 返回队列中第 k 个元素的值 (0-indexed)
// Len() -> 返回队列的长度
// Pop(c) -> 从队列尾部删除 c 个元素，返回删除的元素之和
// Popleft(c) -> 从队列头部删除 c 个元素，返回删除的元素之和
// Append(v, c) -> 在队列尾部添加 c 个 v
// AppendLeft(v, c) -> 在队列头部添加 c 个 v

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	dq := NewLargeDeque()
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var c, x int
			fmt.Fscan(in, &c, &x)
			// dq.Append(x, c)
			dq.Insert(dq.Len(), x, c)
		} else if t == 2 {
			var k int
			fmt.Fscan(in, &k)
			// fmt.Fprintln(out, dq.Popleft(k))
			fmt.Fprintln(out, dq.Erase(0, k))
		}
	}
}

var seed = uint64(time.Now().UnixNano()/2 + 1)

func nextRand() uint64 {
	seed ^= seed << 7
	seed ^= seed >> 9
	return seed
}

type Node struct {
	value, count, size, sum int
	priority                uint64
	left, right             *Node
}

func newNode(v, c int) *Node {
	return &Node{
		value: v, count: c,
		size: c, sum: v * c,
		priority: nextRand(),
	}
}

func update(nd *Node) {
	nd.size = nd.count
	nd.sum = nd.value * nd.count
	if nd.left != nil {
		nd.size += nd.left.size
		nd.sum += nd.left.sum
	}
	if nd.right != nil {
		nd.size += nd.right.size
		nd.sum += nd.right.sum
	}
}

// (<=k, >k)
func split(node *Node, k int) (a, b *Node) {
	if node == nil {
		return nil, nil
	}
	leftSize := 0
	if node.left != nil {
		leftSize = node.left.size
	}
	if k < leftSize {
		a, node.left = split(node.left, k)
		update(node)
		return a, node
	}
	if k < leftSize+node.count {
		kept := k - leftSize
		aNode := &Node{value: node.value, count: kept, priority: node.priority}
		bNode := &Node{value: node.value, count: node.count - kept, priority: node.priority}
		aNode.left, bNode.right = node.left, node.right
		update(aNode)
		update(bNode)
		return aNode, bNode
	}
	node.right, b = split(node.right, k-leftSize-node.count)
	update(node)
	return node, b
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority < b.priority {
		a.right = merge(a.right, b)
		update(a)
		return a
	}
	b.left = merge(a, b.left)
	update(b)
	return b
}

type LargeDeque struct {
	root *Node
}

func NewLargeDeque() *LargeDeque {
	return &LargeDeque{}
}

func (dq *LargeDeque) Append(v, c int) {
	if c <= 0 {
		return
	}
	total := dq.Len()
	a, b := split(dq.root, total-1)
	if b != nil && b.value == v {
		b.count += c
		update(b)
	} else {
		b = merge(b, newNode(v, c))
	}
	dq.root = merge(a, b)
}

func (dq *LargeDeque) AppendLeft(v, c int) {
	if c <= 0 {
		return
	}
	a, b := split(dq.root, 0)
	if a != nil && a.value == v {
		a.count += c
		update(a)
	} else {
		a = merge(newNode(v, c), a)
	}
	dq.root = merge(a, b)
}

func (dq *LargeDeque) Pop(c int) int {
	total := dq.Len()
	a, b := split(dq.root, total-c)
	var dfs func(node *Node) int
	dfs = func(node *Node) int {
		if node == nil {
			return 0
		}
		return dfs(node.left) + dfs(node.right) + node.value*node.count
	}
	s := dfs(b)
	dq.root = a
	return s
}

func (dq *LargeDeque) Popleft(c int) int {
	a, b := split(dq.root, c)
	var dfs func(node *Node) int
	dfs = func(node *Node) int {
		if node == nil {
			return 0
		}
		return dfs(node.left) + dfs(node.right) + node.value*node.count
	}
	s := dfs(a)
	dq.root = b
	return s
}

func (dq *LargeDeque) Len() int {
	if dq.root == nil {
		return 0
	}
	return dq.root.size
}

func (dq *LargeDeque) At(k int) int {
	node := dq.root
	for node != nil {
		leftSize := 0
		if node.left != nil {
			leftSize = node.left.size
		}
		if k < leftSize {
			node = node.left
		} else if k < leftSize+node.count {
			return node.value
		} else {
			k -= leftSize + node.count
			node = node.right
		}
	}
	panic("index out of range")
}

// 区间求和 [l, l+length).
func (dq *LargeDeque) RangeSum(l, length int) int {
	a, bc := split(dq.root, l)
	b, c := split(bc, length)
	res := 0
	if b != nil {
		res = b.sum
	}
	dq.root = merge(a, merge(b, c))
	return res
}

// 在位置 pos 前插入 c 个 v.
func (dq *LargeDeque) Insert(pos, v, c int) {
	a, b := split(dq.root, pos)
	mid := newNode(v, c)
	dq.root = merge(a, merge(mid, b))
}

// 删除位置 pos 开始的 c 个元素，返回它们的和.
func (dq *LargeDeque) Erase(pos, count int) int {
	a, bc := split(dq.root, pos)
	b, c := split(bc, count)
	res := 0
	if b != nil {
		res = b.sum
	}
	dq.root = merge(a, c)
	return res
}
