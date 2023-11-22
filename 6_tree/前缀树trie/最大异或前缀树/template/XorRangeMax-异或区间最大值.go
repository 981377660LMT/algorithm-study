// XorRangeMax-区间最大异或值
// TODO：能不能undo?copy?

package main

import (
	"fmt"
	"math/bits"
	"sort"
)

func main() {
	fmt.Println(maximumStrongPairXor([]int{1, 3, 1}))
}

// https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
// 2935. 找出强数对的最大异或值 II
func maximumStrongPairXor(nums []int) int {
	sort.Ints(nums)
	xm := NewXorRangeMax(nums[len(nums)-1])
	left, res := 0, 0
	for right, cur := range nums {
		xm.Push(cur)
		for left <= right && cur > 2*nums[left] {
			left++
		}
		maxXor, _ := xm.Query(left, right+1, cur)
		res = max(res, maxXor)
	}
	return res
}

// 区间异或最大值.
type XorRangeMax struct {
	trie  *XorTrieSimplePersistent
	roots []*Node
	size  int
}

func NewXorRangeMax(max int) *XorRangeMax {
	trie := NewXorTrieSimplePersistent(max)
	root := trie.NewRoot()
	return &XorRangeMax{
		trie:  trie,
		roots: []*Node{root},
	}
}

// 查询x与区间[start,end)中的数异或的最大值以及最大值对应的下标.
// 如果不存在,返回0,-1.
func (xm *XorRangeMax) Query(start, end int, x int) (maxXor, maxIndex int) {
	if start < 0 {
		start = 0
	}
	if size := xm.Size(); end > size {
		end = size
	}
	if start >= end {
		return 0, -1
	}
	maxXor, node := xm.trie.Query(xm.roots[end], x, start)
	if node == nil {
		return 0, -1
	}
	return maxXor, node.lastIndex
}

func (xm *XorRangeMax) Push(num int) {
	xm.roots = append(xm.roots, xm.trie.Insert(xm.roots[xm.size], num, xm.size))
	xm.size += 1
}

func (xm *XorRangeMax) Size() int {
	return xm.size
}

type Node struct {
	lastIndex int // 最后一次被更新的时间
	chidlren  [2]*Node
}

type XorTrieSimplePersistent struct {
	bit int
}

func NewXorTrieSimplePersistent(upper int) *XorTrieSimplePersistent {
	return &XorTrieSimplePersistent{bit: bits.Len(uint(upper))}
}

func (trie *XorTrieSimplePersistent) NewRoot() *Node {
	return nil
}

func (trie *XorTrieSimplePersistent) Copy(node *Node) *Node {
	if node == nil {
		return node
	}
	return &Node{
		lastIndex: node.lastIndex,
		chidlren:  node.chidlren,
	}
}

func (trie *XorTrieSimplePersistent) Insert(root *Node, num int, lastIndex int) *Node {
	if root == nil {
		root = &Node{}
	}
	return trie._insert(root, num, trie.bit-1, lastIndex)
}

// 查询num与root中的数异或的最大值以及最大值对应的结点.
// !如果root为nil,返回0.
func (trie *XorTrieSimplePersistent) Query(root *Node, num int, leftIndex int) (maxXor int, node *Node) {
	if root == nil {
		return
	}
	for k := trie.bit - 1; k >= 0; k-- {
		bit := (num >> k) & 1
		if root.chidlren[bit^1] != nil && root.chidlren[bit^1].lastIndex >= leftIndex {
			bit ^= 1
			maxXor |= 1 << k
		}
		root = root.chidlren[bit]
	}
	return maxXor, root
}

func (trie *XorTrieSimplePersistent) _insert(root *Node, num int, depth int, lastIndex int) *Node {
	root = trie.Copy(root)
	root.lastIndex = lastIndex
	if depth < 0 {
		return root
	}
	bit := (num >> depth) & 1
	if root.chidlren[bit] == nil {
		root.chidlren[bit] = &Node{}
	}
	root.chidlren[bit] = trie._insert(root.chidlren[bit], num, depth-1, lastIndex)
	return root
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(nums ...int) int {
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}
