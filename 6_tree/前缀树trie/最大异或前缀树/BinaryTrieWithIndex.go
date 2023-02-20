// https://ei1333.github.io/library/structure/trie/binary-trie.hpp

// 带有索引的二进制前缀树
// !添加时可以指定添加的数的id,注意所有数必须为非负整数.
// 如果不都为非负整数,请先加上一个很大的偏移量OFFSET,例如1<<30.
// 模拟树状数组/SortedList

package main

import "fmt"

// https://leetcode.cn/problems/count-the-number-of-fair-pairs/submissions/
func countFairPairs(nums []int, lower int, upper int) int64 {
	res := 0
	sl := NewBinaryTrie(34)
	const OFFSET = 1 << 30
	for _, num := range nums {
		sl.Add(num+OFFSET, 1)
		res += sl.CountLessOrEqual(upper-num+OFFSET) - sl.CountLess(lower-num+OFFSET)
	}
	return int64(res)
}

func main() {
	trie := NewBinaryTrie(30)
	for i := 0; i < 10; i++ {
		trie.AddWithId(i, 1, i)
	}
	res, node := trie.Kth(5)
	fmt.Println(res, node.Ids)
	trie.XorAll(1)
	res, node = trie.Kth(5)
	fmt.Println(res, node.Ids)
	trie.XorAll(1)
}

type BinaryTrie struct {
	xorLazy int
	root    *node
	maxLog  int
}

type node struct {
	Ids   []int // 出现过的id
	Count int   // 以该节点为结尾的数的个数
	next  [2]*node
}

func NewBinaryTrie(maxLog int) *BinaryTrie {
	return &BinaryTrie{root: &node{}, maxLog: maxLog}
}

func (bt *BinaryTrie) AddWithId(num, count, id int) {
	bt.root = bt.add(bt.root, num, id, bt.maxLog, count, true)
}

// num>=0
func (bt *BinaryTrie) Add(num, count int) {
	bt.AddWithId(num, count, -1)
}

// num>=0
func (bt *BinaryTrie) Erase(num, count int) {
	bt.AddWithId(num, -count, -1)
}

// Return the number of num in the trie.
// If num is not in the trie, return nil.
//  num>=0.
func (bt *BinaryTrie) Find(num int) *node {
	return bt.find(bt.root, num, bt.maxLog)
}

func (bt *BinaryTrie) IndexOfAll(num int) []int {
	node := bt.find(bt.root, num, bt.maxLog)
	if node == nil {
		return nil
	}
	return node.Ids
}

func (bt *BinaryTrie) Count(num int) int {
	node := bt.find(bt.root, num, bt.maxLog)
	if node == nil {
		return 0
	}
	return node.Count
}

// 0<=k<exist.
func (bt *BinaryTrie) Kth(k int) (int, *node) {
	return bt.kthElement(bt.root, k, bt.maxLog)
}

func (bt *BinaryTrie) Max() (int, *node) {
	return bt.Kth(bt.root.Count - 1)
}

func (bt *BinaryTrie) Min() (int, *node) {
	return bt.Kth(0)
}

func (bt *BinaryTrie) CountLess(num int) int {
	return bt.countLess(bt.root, num, bt.maxLog)
}

func (bt *BinaryTrie) BisectLeft(num int) int {
	return bt.CountLess(num)
}

func (bt *BinaryTrie) CountLessOrEqual(num int) int {
	return bt.CountLess(num + 1)
}

func (bt *BinaryTrie) BisectRight(num int) int {
	return bt.CountLessOrEqual(num)
}

func (bt *BinaryTrie) XorAll(x int) {
	bt.xorLazy ^= x
}

// 修改克隆逻辑可以可持久化.
func (bt *BinaryTrie) clone(t *node) *node {
	return t
	// return &node{Ids: append([]int{}, t.Ids...), Count: t.Count, next: t.next}
}

func (bt *BinaryTrie) add(t *node, bit, idx, depth, x int, need bool) *node {
	if need {
		t = bt.clone(t)
	}
	if depth == -1 {
		t.Count += x
		if idx >= 0 {
			t.Ids = append(t.Ids, idx)
		}
	} else {
		f := (bt.xorLazy >> depth) & 1
		to := &t.next[f^((bit>>depth)&1)] // TODO
		if *to == nil {
			*to = &node{}
			need = false
		}
		*to = bt.add(*to, bit, idx, depth-1, x, false)
		t.Count += x
	}
	return t
}

func (bt *BinaryTrie) find(t *node, bit, depth int) *node {
	if depth == -1 {
		return t
	}
	f := (bt.xorLazy >> depth) & 1
	to := t.next[f^((bit>>depth)&1)]
	if to == nil {
		return nil
	}
	return bt.find(to, bit, depth-1)
}

func (bt *BinaryTrie) kthElement(t *node, k, bitIndex int) (int, *node) {
	if bitIndex == -1 {
		return 0, t
	}
	f := (bt.xorLazy >> bitIndex) & 1
	count := 0
	if t.next[f] != nil {
		count = t.next[f].Count
	}
	if count <= k {
		res, node := bt.kthElement(t.next[f^1], k-count, bitIndex-1)
		res |= 1 << bitIndex
		return res, node
	}
	return bt.kthElement(t.next[f], k, bitIndex-1)
}

func (bt *BinaryTrie) countLess(t *node, bit, bitIndex int) int {
	if bitIndex == -1 {
		return 0
	}
	res := 0
	f := (bt.xorLazy >> bitIndex) & 1
	if (bit>>bitIndex)&1 == 1 && t.next[f] != nil {
		res += t.next[f].Count
	}
	if t.next[f^(bit>>bitIndex&1)] != nil {
		res += bt.countLess(t.next[f^(bit>>bitIndex&1)], bit, bitIndex-1)
	}
	return res
}
