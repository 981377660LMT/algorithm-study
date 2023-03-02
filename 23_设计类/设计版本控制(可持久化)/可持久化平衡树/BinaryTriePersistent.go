// https://ei1333.github.io/library/structure/trie/binary-trie.hpp

// 可持久化的二进制前缀树
// 可以模拟可持久化平衡树

package main

import "fmt"

func main() {
	trie := NewBinaryTrie(30)
	var gits []*Node // 各个版本的根节点
	gits = append(gits, trie.Alloc())
	gits = append(gits, trie.Add(gits[0], 0, 1))
	fmt.Println(trie.Count(gits[0], 0)) // 0
	fmt.Println(trie.Count(gits[1], 0)) // 1
	gits = append(gits, trie.Add(gits[1], 1, 1))
	fmt.Println(trie.Count(gits[2], 0)) // 1
	fmt.Println(trie.Count(gits[2], 1)) // 1
	gits = append(gits, trie.Add(gits[2], 0, -1))
	fmt.Println(trie.Count(gits[3], 0)) // 0
}

func findMaximumXOR(nums []int) int {
	trie := NewBinaryTrie(32)
	root := trie.Alloc()
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

type BinaryTrie struct {
	xorLazy int
	maxLog  int
}

type Node struct {
	count int
	next  [2]*Node
}

// maxLog: max log of x
// useId: whether to record the id of the number
func NewBinaryTrie(maxLog int) *BinaryTrie {
	return &BinaryTrie{maxLog: maxLog}
}

// Build root node.
func (bt *BinaryTrie) Alloc() *Node {
	return &Node{}
}

func (bt *BinaryTrie) Add(root *Node, num, count int) *Node {
	return bt.add(root, num, bt.maxLog, count, true)
}

func (bt *BinaryTrie) Count(root *Node, num int) int {
	node := bt.find(root, num, bt.maxLog)
	if node == nil {
		return 0
	}
	return node.count
}

// 0<=k<exist.
//  如果不存在,返回(*,false).
func (bt *BinaryTrie) Kth(root *Node, k int) (res int, ok bool) {
	return bt.kthElement(root, k, bt.maxLog)
}

func (bt *BinaryTrie) Max(root *Node) (res int, ok bool) {
	return bt.Kth(root, root.count-1)
}

func (bt *BinaryTrie) Min(root *Node) (res int, ok bool) {
	return bt.Kth(root, 0)
}

func (bt *BinaryTrie) CountLess(root *Node, num int) int {
	return bt.countLess(root, num, bt.maxLog)
}

func (bt *BinaryTrie) BisectLeft(root *Node, num int) int {
	return bt.CountLess(root, num)
}

func (bt *BinaryTrie) CountLessOrEqual(root *Node, num int) int {
	return bt.CountLess(root, num+1)
}

func (bt *BinaryTrie) BisectRight(root *Node, num int) int {
	return bt.CountLessOrEqual(root, num)
}

func (bt *BinaryTrie) XorAll(x int) {
	bt.xorLazy ^= x
}

// 每次克隆时返回新的结点实现可持久化.
func (bt *BinaryTrie) clone(t *Node) *Node {
	return &Node{count: t.count, next: t.next}
}

func (bt *BinaryTrie) add(t *Node, bit, depth, x int, need bool) *Node {
	if need {
		t = bt.clone(t)
	}
	if depth == -1 {
		t.count += x
	} else {
		f := (bt.xorLazy >> depth) & 1
		to := &t.next[f^((bit>>depth)&1)]
		if *to == nil {
			*to = bt.Alloc()
			need = false
		}
		*to = bt.add(*to, bit, depth-1, x, false)
		t.count += x
	}
	return t
}

func (bt *BinaryTrie) find(t *Node, bit, depth int) *Node {
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

func (bt *BinaryTrie) kthElement(t *Node, k, bitIndex int) (int, bool) {
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
		res, ok := bt.kthElement(t.next[f^1], k-count, bitIndex-1)
		res |= 1 << bitIndex
		return res, ok
	}
	return bt.kthElement(t.next[f], k, bitIndex-1)
}

func (bt *BinaryTrie) countLess(t *Node, bit, bitIndex int) int {
	if bitIndex == -1 {
		return 0
	}
	res := 0
	f := (bt.xorLazy >> bitIndex) & 1
	if (bit>>bitIndex)&1 == 1 && t.next[f] != nil {
		res += t.next[f].count
	}
	if t.next[f^(bit>>bitIndex&1)] != nil {
		res += bt.countLess(t.next[f^(bit>>bitIndex&1)], bit, bitIndex-1)
	}
	return res
}
