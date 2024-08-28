// api:
// - BinaryTrie(log int32, persistent bool)
// - NewRoot() *trieNode
// - Add(root *trieNode, value int, count int) *trieNode
// - Enuerate(root *trieNode, f func(value, count int))
// - Kth(root *trieNode, k int, xor int) int
// - Min(root *trieNode, xor int) int
// - Max(root *trieNode, xor int) int
// - PrefixCount(root *trieNode, upper int, xor int) int
// - Count(root *trieNode, lo, hi, xor int) int

package main

import (
	"math/bits"
	"sort"
)

// 1803. 统计异或值在范围内的数对有多少
// https://leetcode.cn/problems/count-pairs-with-xor-in-a-range/description/
func countPairs(nums []int, low int, high int) int {
	log := int32(bits.Len(uint(maxs(nums...))))
	bt := NewBinaryTrie(log, false)
	root := bt.NewRoot()
	for _, v := range nums {
		root = bt.Add(root, v, 1)
	}
	res := 0
	for _, v := range nums {
		res += bt.Count(root, low, high+1, v)
	}
	return res / 2
}

// 2935. 找出强数对的最大异或值 II
// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	res, left, n := 0, 0, len(nums)
	log := int32(bits.Len(uint(nums[n-1])))
	trie := NewBinaryTrie(log, false)
	root := trie.NewRoot()
	for right, cur := range nums {
		root = trie.Add(root, cur, 1)
		for left <= right && cur > 2*nums[left] {
			root = trie.Add(root, nums[left], -1)
			left++
		}
		res = max(res, trie.Max(root, cur))
	}
	return res
}

type trieNode struct {
	width int32
	value int
	count int
	l, r  *trieNode
}

type BinaryTrie struct {
	persistent bool
	log        int32
}

func NewBinaryTrie(log int32, persistent bool) *BinaryTrie {
	return &BinaryTrie{persistent: persistent, log: log}
}

func (bt *BinaryTrie) NewRoot() *trieNode {
	return nil
}

func (bt *BinaryTrie) Add(root *trieNode, value int, count int) *trieNode {
	if root == nil {
		root = bt.newNode(0, 0)
	}
	return bt.addRec(root, bt.log, value, count)
}

func (bt *BinaryTrie) Enuerate(root *trieNode, f func(value, count int)) {
	var dfs func(*trieNode, int, int32)
	dfs = func(root *trieNode, val int, ht int32) {
		if ht == 0 {
			f(val, root.count)
			return
		}
		if c := root.l; c != nil {
			dfs(c, val<<c.width|c.value, ht-c.width)
		}
		if c := root.r; c != nil {
			dfs(c, val<<c.width|c.value, ht-c.width)
		}
	}
	if root != nil {
		dfs(root, 0, bt.log)
	}
}

func (bt *BinaryTrie) Kth(root *trieNode, k int, xor int) int {
	return bt.kthRec(root, 0, k, bt.log, xor) ^ xor
}

func (bt *BinaryTrie) Min(root *trieNode, xor int) int {
	return bt.Kth(root, 0, xor)
}

func (bt *BinaryTrie) Max(root *trieNode, xor int) int {
	return bt.Kth(root, root.count-1, xor)
}

// [0, upper)
func (bt *BinaryTrie) PrefixCount(root *trieNode, upper int, xor int) int {
	if root == nil {
		return 0
	}
	return bt.prefixCountRec(root, bt.log, upper, xor, 0)
}

// [lo, hi)
func (bt *BinaryTrie) Count(root *trieNode, lo, hi, xor int) int {
	return bt.PrefixCount(root, hi, xor) - bt.PrefixCount(root, lo, xor)
}

func (bt *BinaryTrie) newNode(width int32, value int) *trieNode {
	return &trieNode{width: width, value: value}
}

func (bt *BinaryTrie) copyNode(c *trieNode) *trieNode {
	if c == nil || !bt.persistent {
		return c
	}
	return &trieNode{width: c.width, value: c.value, count: c.count, l: c.l, r: c.r}
}

func (bt *BinaryTrie) addRec(root *trieNode, ht int32, val int, count int) *trieNode {
	root = bt.copyNode(root)
	root.count += count
	if ht == 0 {
		return root
	}
	goRight := (val>>(ht-1))&1 == 1
	c := root.l
	if goRight {
		c = root.r
	}
	if c == nil {
		c = bt.newNode(ht, val)
		c.count = count
		if goRight {
			root.r = c
		} else {
			root.l = c
		}
		return root
	}
	w := c.width
	if (val >> (ht - w)) == c.value {
		c = bt.addRec(c, ht-w, val&(1<<(ht-w)-1), count)
		if goRight {
			root.r = c
		} else {
			root.l = c
		}
		return root
	}
	same := w - 1 - int32(topbit((val>>(ht-w))^(c.value)))
	n := bt.newNode(same, c.value>>(w-same))
	n.count = c.count + count
	c = bt.copyNode(c)
	c.width = w - same
	c.value &= (1 << (w - same)) - 1
	if (val>>(ht-same-1))&1 == 1 {
		n.l = c
		n.r = bt.newNode(ht-same, val&(1<<(ht-same)-1))
		n.r.count = count
	} else {
		n.r = c
		n.l = bt.newNode(ht-same, val&(1<<(ht-same)-1))
		n.l.count = count
	}
	if goRight {
		root.r = n
	} else {
		root.l = n
	}
	return root
}

func (bt *BinaryTrie) kthRec(root *trieNode, val int, k int, ht int32, xor int) int {
	if ht == 0 {
		return val
	}
	left, right := root.l, root.r
	if (xor>>(ht-1))&1 == 1 {
		left, right = right, left
	}
	sl := 0
	if left != nil {
		sl = left.count
	}
	var c *trieNode
	if k < sl {
		c = left
	} else {
		c = right
		k -= sl
	}
	w := c.width
	return bt.kthRec(c, val<<w|(c.value), k, ht-w, xor)
}

func (bt *BinaryTrie) prefixCountRec(root *trieNode, ht int32, limit int, xor int, val int) int {
	now := (val << ht) ^ (xor)
	if (limit >> ht) > (now >> ht) {
		return root.count
	}
	if ht == 0 || (limit>>ht) < (now>>ht) {
		return 0
	}
	res := 0
	if c := root.l; c != nil {
		w := c.width
		res += bt.prefixCountRec(c, ht-w, limit, xor, val<<w|c.value)
	}
	if c := root.r; c != nil {
		w := c.width
		res += bt.prefixCountRec(c, ht-w, limit, xor, val<<w|c.value)
	}
	return res
}

func topbit(x int) int {
	return bits.Len64(uint64(x)) - 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
