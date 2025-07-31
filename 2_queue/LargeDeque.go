// https://atcoder.jp/contests/abc413/tasks/abc413_c
//
// at(k) -> 返回队列中第 k 个元素的值 (0-indexed)
// append(v, c) -> 在队列尾部添加 c 个元素 v
// appendleft(v, c) -> 在队列头部添加 c 个元素 v
// pop(c) -> 从队列尾部删除 c 个元素，返回删除的元素之和
// popleft(c) -> 从队列头部删除 c 个元素，返回删除的元素之和
// __len__() -> 返回队列的长度
// __sum__() -> 返回队列中所有元素的和

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
			dq.Append(x, c)
		} else if t == 2 {
			var k int
			fmt.Fscan(in, &k)
			fmt.Fprintln(out, dq.Popleft(k))
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
	value, count, size int
	priority           uint64
	left, right        *Node
}

func newNode(v, c int) *Node {
	return &Node{value: v, count: c, size: c, priority: nextRand()}
}

func update(node *Node) {
	node.size = node.count
	if node.left != nil {
		node.size += node.left.size
	}
	if node.right != nil {
		node.size += node.right.size
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
	// split [0,total-c) + [total-c,total)
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
