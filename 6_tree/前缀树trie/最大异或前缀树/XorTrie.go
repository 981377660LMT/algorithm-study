package main

import "math/bits"

func findMaximumXOR(nums []int) int {
	res := 0
	trie := NewXORTrie(1<<32, false)
	for _, num := range nums {
		max_, _ := trie.MaxXor(num)
		res = max(res, max_)
		trie.Insert(num, 0)
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
	Counter  map[int]int
	children [2]*Node
	count    int
}

type XORTrie struct {
	bit   int
	useId bool
	root  *Node
}

func NewXORTrie(upper int, useId bool) *XORTrie {
	res := &XORTrie{bit: bits.Len(uint(upper)), useId: useId}
	res.root = res.alloc()
	return res
}

func (bt *XORTrie) Insert(num int, id int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[bit] == nil {
			root.children[bit] = bt.alloc()
		}
		root = root.children[bit]
		root.count++
	}
	if bt.useId {
		root.Counter[id]++
	}
}

func (bt *XORTrie) MaxXor(num int) (int, *Node) {
	root := bt.root
	res := 0
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[bit^1] != nil && root.children[bit^1].count > 0 {
			res |= 1 << i
			root = root.children[bit^1]
		} else if root.children[bit] != nil && root.children[bit].count > 0 {
			root = root.children[bit]
		} else {
			break
		}
	}
	return res, root
}

func (bt *XORTrie) Remove(num int, id int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		root = root.children[bit]
		root.count--
	}
	if bt.useId {
		root.Counter[id]--
		if root.Counter[id] == 0 {
			delete(root.Counter, id)
		}
	}
}

func (bt *XORTrie) alloc() *Node {
	if !bt.useId {
		return &Node{}
	}
	return &Node{Counter: make(map[int]int)}
}
