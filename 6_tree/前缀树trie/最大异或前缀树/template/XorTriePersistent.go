// https://ei1333.github.io/library/structure/trie/binary-trie.hpp

// 可持久化的二进制前缀树
// 可以模拟可持久化平衡树

package main

import (
	"fmt"
	"math/bits"
	"sort"
)

func main() {
	trie := NewBinaryTriePersistent(1 << 30)
	var gits []*Node // 各个版本的根节点
	gits = append(gits, trie.NewRoot())
	gits = append(gits, trie.Add(gits[0], 0, 1))
	fmt.Println(trie.Count(gits[0], 0)) // 0
	fmt.Println(trie.Count(gits[1], 0)) // 1
	gits = append(gits, trie.Add(gits[1], 1, 1))
	fmt.Println(trie.Count(gits[2], 0)) // 1
	fmt.Println(trie.Count(gits[2], 1)) // 1
	gits = append(gits, trie.Add(gits[2], 0, -1))
	fmt.Println(trie.Count(gits[3], 0)) // 0
}

// 2935. 找出强数对的最大异或值 II
// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	res, left, n := 0, 0, len(nums)
	trie := NewBinaryTriePersistent(nums[n-1])
	root := trie.NewRoot()
	for right, cur := range nums {
		root = trie.Add(root, cur, 1)
		for left <= right && cur > 2*nums[left] {
			root = trie.Remove(root, nums[left], 1)
			left++
		}
		trie.XorAll(cur)
		max_, _ := trie.Max(root)
		res = max(res, max_)
		trie.XorAll(cur)
	}
	return res
}

func findMaximumXOR(nums []int) int {
	trie := NewBinaryTriePersistent(1 << 32)
	root := trie.NewRoot()
	res := 0
	for _, num := range nums {
		trie.XorAll(num)
		max_, _ := trie.Max(root)
		trie.XorAll(num)
		res = max(res, max_)
		root = trie.Add(root, num, 1)
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
	maxLog  int
	xorLazy int
}

type Node struct {
	count int
	next  [2]*Node
}

// maxLog: max log of x
func NewBinaryTriePersistent(max int) *BinaryTriePersistent {
	return &BinaryTriePersistent{maxLog: bits.Len(uint(max))}
}

// Build root node.
func (bt *BinaryTriePersistent) NewRoot() *Node {
	return &Node{}
}

func (bt *BinaryTriePersistent) Add(root *Node, num, count int) *Node {
	return bt._add(root, num, bt.maxLog, count, true)
}

func (bt *BinaryTriePersistent) Remove(root *Node, num, count int) *Node {
	return bt.Add(root, num, -count)
}

func (bt *BinaryTriePersistent) Count(root *Node, num int) int {
	node := bt._find(root, num, bt.maxLog)
	if node == nil {
		return 0
	}
	return node.count
}

// 0<=k<exist.
//
//	如果不存在,返回(*,false).
func (bt *BinaryTriePersistent) Kth(root *Node, k int) (res int, ok bool) {
	return bt._kthElement(root, k, bt.maxLog)
}

func (bt *BinaryTriePersistent) Max(root *Node) (res int, ok bool) {
	return bt.Kth(root, root.count-1)
}

func (bt *BinaryTriePersistent) Min(root *Node) (res int, ok bool) {
	return bt.Kth(root, 0)
}

func (bt *BinaryTriePersistent) CountLess(root *Node, num int) int {
	return bt._countLess(root, num, bt.maxLog)
}

func (bt *BinaryTriePersistent) BisectLeft(root *Node, num int) int {
	return bt.CountLess(root, num)
}

func (bt *BinaryTriePersistent) CountLessOrEqual(root *Node, num int) int {
	return bt.CountLess(root, num+1)
}

func (bt *BinaryTriePersistent) BisectRight(root *Node, num int) int {
	return bt.CountLessOrEqual(root, num)
}

func (bt *BinaryTriePersistent) XorAll(x int) {
	bt.xorLazy ^= x
}

func (bt *BinaryTriePersistent) Size(root *Node) int {
	if root == nil {
		return 0
	}
	return root.count
}

// 每次克隆时返回新的结点实现可持久化.
func (bt *BinaryTriePersistent) Copy(t *Node) *Node {
	return &Node{count: t.count, next: t.next}
}

func (bt *BinaryTriePersistent) _add(t *Node, bit, depth, x int, need bool) *Node {
	if need {
		t = bt.Copy(t)
	}
	if depth == -1 {
		t.count += x
	} else {
		f := (bt.xorLazy >> depth) & 1
		to := &t.next[f^((bit>>depth)&1)]
		if *to == nil {
			*to = bt.NewRoot()
			need = false
		}
		*to = bt._add(*to, bit, depth-1, x, false)
		t.count += x
	}
	return t
}

func (bt *BinaryTriePersistent) _find(t *Node, bit, depth int) *Node {
	if depth == -1 {
		return t
	}
	f := (bt.xorLazy >> depth) & 1
	to := t.next[f^((bit>>depth)&1)]
	if to == nil {
		return nil
	}
	return bt._find(to, bit, depth-1)
}

func (bt *BinaryTriePersistent) _kthElement(t *Node, k, bitIndex int) (int, bool) {
	if t == nil {
		return 0, false
	}
	if bitIndex == -1 {
		return 0, true
	}
	f := (bt.xorLazy >> bitIndex) & 1
	count := 0
	if t.next[f] != nil {
		count = t.next[f].count
	}
	if count <= k {
		res, ok := bt._kthElement(t.next[f^1], k-count, bitIndex-1)
		res |= 1 << bitIndex
		return res, ok
	}
	return bt._kthElement(t.next[f], k, bitIndex-1)
}

func (bt *BinaryTriePersistent) _countLess(t *Node, bit, bitIndex int) int {
	if bitIndex == -1 {
		return 0
	}
	res := 0
	f := (bt.xorLazy >> bitIndex) & 1
	if (bit>>bitIndex)&1 == 1 && t.next[f] != nil {
		res += t.next[f].count
	}
	if t.next[f^(bit>>bitIndex&1)] != nil {
		res += bt._countLess(t.next[f^(bit>>bitIndex&1)], bit, bitIndex-1)
	}
	return res
}
