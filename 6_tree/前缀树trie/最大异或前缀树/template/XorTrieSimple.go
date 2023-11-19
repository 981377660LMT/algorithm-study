package main

import (
	"math/bits"
	"sort"
)

// 2935. 找出强数对的最大异或值 II
// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	res, left, n := 0, 0, len(nums)
	trie := NewXORTrie(nums[n-1])
	for right, cur := range nums {
		trie.Insert(cur)
		for left <= right && cur > 2*nums[left] {
			trie.Remove(nums[left])
			left++
		}
		res = max(res, trie.Query(cur))
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Node struct {
	count    int
	children [2]*Node
}

type XORTrieSimple struct {
	bit  int
	root *Node
}

func NewXORTrie(upper int) *XORTrieSimple {
	return &XORTrieSimple{
		bit:  bits.Len(uint(upper)),
		root: &Node{},
	}
}

func (bt *XORTrieSimple) Insert(num int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[bit] == nil {
			root.children[bit] = &Node{}
		}
		root = root.children[bit]
		root.count++
	}
	return
}

// 必须保证num存在于trie中.
func (bt *XORTrieSimple) Remove(num int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		root = root.children[bit]
		root.count--
	}
}

func (bt *XORTrieSimple) Query(num int) (maxXor int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[bit^1] != nil && root.children[bit^1].count > 0 {
			maxXor |= 1 << i
			bit ^= 1
		}
		root = root.children[bit]
	}
	return
}
