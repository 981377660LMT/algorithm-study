// https://ei1333.github.io/library/structure/trie/binary-trie.hpp

// 可持久化的二进制前缀树
// !可以模拟可持久化平衡树

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
	gits = append(gits, trie.Erase(gits[2], 0, 1))
	fmt.Println(trie.Count(gits[3], 0)) // 0
}

type BinaryTrie struct {
	xorLazy int
	maxLog  int
}

type Node struct {
	Ids   []int // 出现过的id
	Count int   // 以该节点为结尾的数的个数
	next  [2]*Node
}

func NewBinaryTrie(maxLog int) *BinaryTrie {
	return &BinaryTrie{maxLog: maxLog}
}

// Build root node.
func (bt *BinaryTrie) Alloc() *Node {
	return &Node{}
}

func (bt *BinaryTrie) AddWithId(root *Node, num, count, id int) *Node {
	return bt.add(root, num, id, bt.maxLog, count, true)
}

func (bt *BinaryTrie) Add(root *Node, num, count int) *Node {
	return bt.AddWithId(root, num, count, -1)
}

func (bt *BinaryTrie) Erase(root *Node, num, count int) *Node {
	return bt.AddWithId(root, num, -count, -1)
}

// Return the number of num in the trie.
// If num is not in the trie, return nil.
//  num>=0.
func (bt *BinaryTrie) Find(root *Node, num int) *Node {
	return bt.find(root, num, bt.maxLog)
}

func (bt *BinaryTrie) IndexOfAll(root *Node, num int) []int {
	node := bt.find(root, num, bt.maxLog)
	if node == nil {
		return nil
	}
	return node.Ids
}

func (bt *BinaryTrie) Count(root *Node, num int) int {
	node := bt.find(root, num, bt.maxLog)
	if node == nil {
		return 0
	}
	return node.Count
}

// 0<=k<exist.
func (bt *BinaryTrie) Kth(root *Node, k int) (int, *Node) {
	return bt.kthElement(root, k, bt.maxLog)
}

func (bt *BinaryTrie) Max(root *Node) (int, *Node) {
	return bt.Kth(root, root.Count-1)
}

func (bt *BinaryTrie) Min(root *Node) (int, *Node) {
	return bt.Kth(root, 0)
}

func (bt *BinaryTrie) CountLess(root *Node, num int) int {
	return bt.countLess(root, num, bt.maxLog)
}

func (bt *BinaryTrie) BisectLeft(root *Node, num int) int {
	return bt.CountLess(root, num)
}

func (bt *BinaryTrie) CountLessOrEqual(root *Node, num int) int {
	if num < 0 {
		return 0
	}
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
	return &Node{Ids: append([]int{}, t.Ids...), Count: t.Count, next: t.next}
}

func (bt *BinaryTrie) add(t *Node, bit, idx, depth, x int, need bool) *Node {
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
		to := &t.next[f^((bit>>depth)&1)]
		if *to == nil {
			*to = &Node{}
			need = false
		}
		*to = bt.add(*to, bit, idx, depth-1, x, false)
		t.Count += x
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

func (bt *BinaryTrie) kthElement(t *Node, k, bitIndex int) (int, *Node) {
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

func (bt *BinaryTrie) countLess(t *Node, bit, bitIndex int) int {
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
