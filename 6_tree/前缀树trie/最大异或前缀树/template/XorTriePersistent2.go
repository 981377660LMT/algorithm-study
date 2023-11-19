// https://maspypy.github.io/library/ds/binary_trie.hpp

// 可持久化的二进制前缀树,可以模拟可持久化平衡树

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	yosupo()
}

func demo() {
	trie := NewBinaryTriePersistent(1e9, true)
	root := trie.NewRoot()
	root = trie.Add(root, 1, 1)
	trie.Enumerate(root, func(value int, count int) { fmt.Println(value, count) })
	for i := 0; i < 10; i++ {
		fmt.Println(_topbit(i))
	}
}

func yosupo() {
	// https://judge.yosupo.jp/problem/set_xor_min
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var Q int
	fmt.Fscan(in, &Q)
	trie := NewBinaryTriePersistent(1<<29, true)
	root := trie.NewRoot()
	for i := 0; i < Q; i++ {
		var op, x int
		fmt.Fscan(in, &op, &x)
		if op == 0 {
			if trie.CountRange(root, x, x+1, 0) == 0 {
				root = trie.Add(root, x, 1)
			}
		} else if op == 1 {
			if trie.CountRange(root, x, x+1, 0) != 0 {
				root = trie.Add(root, x, -1)
			}
		} else {
			res, _ := trie.Min(root, x)
			fmt.Fprintln(out, res)
		}
	}
}

// 2935. 找出强数对的最大异或值 II
// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	res, left, n := 0, 0, len(nums)
	trie := NewBinaryTriePersistent(nums[n-1], false)
	root := trie.NewRoot()
	for right, cur := range nums {
		root = trie.Add(root, cur, 1)
		for left <= right && cur > 2*nums[left] {
			root = trie.Remove(root, nums[left], 1)
			left++
		}
		max_, _ := trie.Max(root, cur)
		res = max(res, max_)
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type BinaryTriePersistent struct {
	maxLog     int
	persistent bool
	xorLazy    int
}

type Node struct {
	width       int
	value       int
	count       int
	left, right *Node
}

func NewBinaryTriePersistent(max int, persistent bool) *BinaryTriePersistent {
	return &BinaryTriePersistent{maxLog: bits.Len(uint(max)), persistent: persistent}
}

func (bt *BinaryTriePersistent) NewRoot() *Node {
	return nil
}

func (bt *BinaryTriePersistent) Add(root *Node, num, count int) *Node {
	if root == nil {
		root = _newNode(0, 0)
	}
	return bt._add(root, bt.maxLog, num, count)
}

func (bt *BinaryTriePersistent) Remove(root *Node, num, count int) *Node {
	return bt.Add(root, num, -count)
}

// 0<=k<exist.
//
//	如果不存在,返回(*,false).
func (bt *BinaryTriePersistent) Kth(root *Node, k int, xor int) (res int, ok bool) {
	if root == nil || k < 0 || k >= root.count {
		return
	}
	return bt._kth(root, 0, k, bt.maxLog, xor) ^ xor, true
}

func (bt *BinaryTriePersistent) Max(root *Node, xor int) (res int, ok bool) {
	return bt.Kth(root, root.count-1, xor)
}

func (bt *BinaryTriePersistent) Min(root *Node, xor int) (res int, ok bool) {
	return bt.Kth(root, 0, xor)
}

func (bt *BinaryTriePersistent) CountLess(root *Node, num int, xor int) int {
	if root == nil {
		return 0
	}
	return bt._prefixCount(root, bt.maxLog, num, xor, 0)
}

func (bt *BinaryTriePersistent) BisectLeft(root *Node, num int, xor int) int {
	return bt.CountLess(root, num, xor)
}

func (bt *BinaryTriePersistent) CountLessOrEqual(root *Node, num int, xor int) int {
	return bt.CountLess(root, num+1, xor)
}

func (bt *BinaryTriePersistent) BisectRight(root *Node, num int, xor int) int {
	return bt.CountLess(root, num+1, xor)
}

func (bt *BinaryTriePersistent) CountRange(root *Node, start, end int, xor int) int {
	return bt.CountLess(root, end, xor) - bt.CountLess(root, start, xor)
}

func (bt *BinaryTriePersistent) Size(root *Node) int {
	if root == nil {
		return 0
	}
	return root.count
}

func (bt *BinaryTriePersistent) Enumerate(root *Node, f func(value int, count int)) {
	if root == nil {
		return
	}
	var dfs func(root *Node, val int, height int)
	dfs = func(root *Node, val int, height int) {
		if height == 0 {
			f(val, root.count)
			return
		}
		if c := root.left; c != nil {
			dfs(c, (val<<c.width)|c.value, height-c.width)
		}
		if c := root.right; c != nil {
			dfs(c, (val<<c.width)|c.value, height-c.width)
		}
	}
	dfs(root, 0, bt.maxLog)
}

func (bt *BinaryTriePersistent) Copy(node *Node) *Node {
	if node == nil || !bt.persistent {
		return node
	}
	return &Node{
		width: node.width,
		value: node.value,
		count: node.count,
		left:  node.left,
		right: node.right,
	}
}

func (bt *BinaryTriePersistent) _add(root *Node, height int, val, count int) *Node {
	root = bt.Copy(root)
	root.count += count
	if height == 0 {
		return root
	}
	goRight := (val>>(height-1))&1 == 1
	var c *Node
	if goRight {
		c = root.right
	} else {
		c = root.left
	}
	if c == nil {
		c = _newNode(height, val)
		c.count = count
		if goRight {
			root.right = c
		} else {
			root.left = c
		}
		return root
	}
	w := c.width
	if (val >> (height - w)) == c.value {
		c = bt._add(c, height-w, val&((1<<(height-w))-1), count)
		if goRight {
			root.right = c
		} else {
			root.left = c
		}
		return root
	}
	same := w - 1 - _topbit((val>>(height-w))^(c.value))
	n := _newNode(same, c.value>>(w-same))
	n.count = c.count + count
	c = bt.Copy(c)
	c.width = w - same
	c.value &= (1 << (w - same)) - 1
	if (val>>(height-same-1))&1 == 1 {
		n.left = c
		n.right = _newNode(height-same, val&((1<<(height-same))-1))
		n.right.count = count
	} else {
		n.right = c
		n.left = _newNode(height-same, val&((1<<(height-same))-1))
		n.left.count = count
	}
	if goRight {
		root.right = n
	} else {
		root.left = n
	}
	return root
}

func (bt *BinaryTriePersistent) _prefixCount(root *Node, height int, limit int, xor int, val int) int {
	now := (val << height) ^ xor
	a, b := limit>>height, now>>height
	if a > b {
		return root.count
	}
	if height == 0 || a < b {
		return 0
	}
	res := 0
	if c := root.left; c != nil {
		w := c.width
		res += bt._prefixCount(c, height-w, limit, xor, (val<<w)|c.value)
	}
	if c := root.right; c != nil {
		w := c.width
		res += bt._prefixCount(c, height-w, limit, xor, (val<<w)|c.value)
	}
	return res
}

func (bt *BinaryTriePersistent) _kth(root *Node, val, k, height, xor int) int {
	if height == 0 {
		return val
	}
	left, right := root.left, root.right
	if (xor>>(height-1))&1 == 1 {
		left, right = right, left
	}
	leftSize := 0
	if left != nil {
		leftSize = left.count
	}
	var c *Node
	if k < leftSize {
		c = left
	} else {
		c = right
		k -= leftSize
	}
	w := c.width
	return bt._kth(c, (val<<w)|c.value, k, height-w, xor)
}

func _newNode(width int, value int) *Node {
	return &Node{width: width, value: value}
}

func _topbit(x int) int {
	if x == 0 {
		return -1
	}
	return 31 - bits.LeadingZeros32(uint32(x))
}
