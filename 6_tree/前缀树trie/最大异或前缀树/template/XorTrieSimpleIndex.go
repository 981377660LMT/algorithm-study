package main

import (
	"fmt"
	"math/bits"
)

func main() {
	arr := []int{1, 3, 1}
	X := NewXorTrieIndex(3)
	for i, v := range arr {
		fmt.Println(X.Query(v))
		X.Insert(v, int32(i))
	}
}

type Node struct {
	index    int32
	children [2]*Node
}

func NewNode() *Node {
	return &Node{index: -1}
}

type XorTrieIndex struct {
	bit  int32
	root *Node
}

func NewXorTrieIndex(upper int) *XorTrieIndex {
	return &XorTrieIndex{
		bit:  int32(bits.Len(uint(upper))),
		root: NewNode(),
	}
}

func (bt *XorTrieIndex) Insert(num int, index int32) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[bit] == nil {
			root.children[bit] = NewNode()
		}
		root = root.children[bit]
		// 保留最小的index
		if root.index == -1 || index < root.index {
			root.index = index
		}
	}
	return
}

// !查询能获得的最大的异或值时的最小下标.
// !如果trie为空,返回-1.
func (bt *XorTrieIndex) Query(num int) int32 {
	if bt.Empty() {
		return -1
	}
	root := bt.root
	res := int32(-1)
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		if root.children[1^bit] != nil {
			root = root.children[1^bit]
		} else {
			root = root.children[bit]
		}
		res = root.index
	}
	return res
}

func (bt *XorTrieIndex) Empty() bool {
	return bt.root.children[0] == nil && bt.root.children[1] == nil
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
