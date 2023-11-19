// https://ei1333.github.io/library/structure/trie/binary-trie.hpp

// 带有索引的二进制前缀树
// !添加时可以指定添加的数的id,注意所有数必须为非负整数.
// 如果不都为非负整数,请先加上一个很大的偏移量OFFSET,例如1<<30.
// 模拟树状数组/SortedList

// !注意查询前调用 XorAll!

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	trie := NewBinaryTrieDynamic(1 << 30)
	trie.Add(0, 1)
	trie.Add(1, 3)
	fmt.Println(trie.Max()) // 1
	trie.Remove(1, 2)
}

// 2935. 找出强数对的最大异或值 II
// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	res, left, n := 0, 0, len(nums)
	trie := NewBinaryTrieDynamic(nums[n-1])
	for right, cur := range nums {
		trie.Add(cur, 1)
		for left <= right && cur > 2*nums[left] {
			trie.Remove(nums[left], 1)
			left++
		}
		trie.XorAll(cur)
		max_, _ := trie.Max()
		res = max(res, max_)
		trie.XorAll(cur)
	}
	return res
}

func findMaximumXOR(nums []int) int {
	trie := NewBinaryTrieDynamic(1 << 32)
	res := 0
	for _, num := range nums {
		trie.XorAll(num)
		max_, _ := trie.Max()
		trie.XorAll(num)
		res = max(res, max_)
		trie.Add(num, 1)
	}
	return res
}

func yosupo() {
	// https://judge.yosupo.jp/problem/set_xor_min
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var Q int
	fmt.Fscan(in, &Q)
	trie := NewBinaryTrieDynamic(1 << 29)
	for i := 0; i < Q; i++ {
		var op, x int
		fmt.Fscan(in, &op, &x)
		if op == 0 {
			if trie.Count(x) == 0 {
				trie.Add(x, 1)
			}
		} else if op == 1 {
			if trie.Count(x) != 0 {
				trie.Add(x, -1)
			}
		} else {
			trie.XorAll(x)
			res, _ := trie.Min()
			trie.XorAll(x)
			fmt.Fprintln(out, res)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type BinaryTrieDynamic struct {
	root    *node
	maxLog  int
	xorLazy int
}

type node struct {
	count int // 以该节点为结尾的数的个数
	next  [2]*node
}

func NewBinaryTrieDynamic(max int) *BinaryTrieDynamic {
	return &BinaryTrieDynamic{root: newNode(), maxLog: bits.Len(uint(max))}
}

func (bt *BinaryTrieDynamic) Add(num, count int) {
	bt.root = bt.add(bt.root, num, bt.maxLog, count, true)
}

func (bt *BinaryTrieDynamic) Remove(num int, count int) {
	bt.Add(num, -count)
}

func (bt *BinaryTrieDynamic) Count(num int) int {
	node := bt.find(bt.root, num, bt.maxLog)
	if node == nil {
		return 0
	}
	return node.count
}

// 0<=k<exist.
func (bt *BinaryTrieDynamic) Kth(k int) (res int, ok bool) {
	return bt.kthElement(bt.root, k, bt.maxLog)
}

func (bt *BinaryTrieDynamic) Max() (res int, ok bool) {
	return bt.Kth(bt.root.count - 1)
}

func (bt *BinaryTrieDynamic) Min() (res int, ok bool) {
	return bt.Kth(0)
}

func (bt *BinaryTrieDynamic) CountLess(num int) int {
	return bt.countLess(bt.root, num, bt.maxLog)
}

func (bt *BinaryTrieDynamic) BisectLeft(num int) int {
	return bt.CountLess(num)
}

func (bt *BinaryTrieDynamic) CountLessOrEqual(num int) int {
	return bt.CountLess(num + 1)
}

func (bt *BinaryTrieDynamic) BisectRight(num int) int {
	return bt.CountLessOrEqual(num)
}

func (bt *BinaryTrieDynamic) XorAll(x int) {
	bt.xorLazy ^= x
}

func (bt *BinaryTrieDynamic) Size() int {
	return bt.root.count
}

func (bt *BinaryTrieDynamic) add(t *node, bit, depth, x int, need bool) *node {
	if depth == -1 {
		t.count += x
	} else {
		f := (bt.xorLazy >> depth) & 1
		to := &t.next[f^((bit>>depth)&1)]
		if *to == nil {
			*to = newNode()
			need = false
		}
		*to = bt.add(*to, bit, depth-1, x, false)
		t.count += x
	}
	return t
}

func (bt *BinaryTrieDynamic) find(t *node, bit, depth int) *node {
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

func (bt *BinaryTrieDynamic) kthElement(t *node, k, bitIndex int) (int, bool) {
	if t == nil {
		return 0, false
	}
	if bitIndex == -1 || t == nil {
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

func (bt *BinaryTrieDynamic) countLess(t *node, bit, bitIndex int) int {
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

func newNode() *node {
	return &node{}
}
