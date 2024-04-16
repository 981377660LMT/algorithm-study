// 按位与最大的二元组
// !给定一个数组，要求找到两个不同的下标i!=j使得A[i]&A[j]最大.

package main

import (
	"fmt"
	"math/bits"
)

// https://www.geeksforgeeks.org/maximum-value-pair-array/
func Solve(nums []int) int {
	countValid := func(pattern int) int {
		count := 0
		for _, num := range nums {
			// 这个数num是之前看过的数.
			if (pattern & num) == pattern {
				count++
			}
		}
		return count
	}

	res := 0
	max_ := 0
	for _, num := range nums {
		if num > max_ {
			max_ = num
		}
	}
	maxBit := bits.Len(uint(max_))
	for bit := maxBit; bit >= 0; bit-- {
		count := countValid(res | (1 << bit))
		if count >= 2 {
			res |= 1 << bit
		}
	}
	return res
}

func main() {
	solve1 := func(nums []int) int {
		trie := NewAndTrie(maxs(nums))
		res := 0
		for _, cur := range nums {
			res = max(res, trie.Query(cur))
			trie.Insert(cur)
		}
		return res
	}

	solve2 := func(nums []int) int {
		res := 0
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				res = max(res, nums[i]&nums[j])
			}
		}
		return res
	}

	// for i := 0; i < 1000; i++ {
	// 	nums := make([]int, 2)
	// 	for j := range nums {
	// 		nums[j] = rand.Intn(10)
	// 	}
	// 	res1, res2 := solve1(nums), solve2(nums)
	// 	res3 := Solve(nums)
	// 	if res2 != res3 {
	// 		fmt.Println(res2, res3, nums)
	// 		panic("not equal1")
	// 	}
	// 	if res1 != res2 {
	// 		fmt.Println(res1, res2, nums)
	// 		panic("not equal2")
	// 	}
	// }
	fmt.Println(solve1([]int{3, 9}))
	fmt.Println(solve2([]int{3, 9}))
}

// 从高往低 如果1超过2个 就往1走 否则就往0走
type Node struct {
	count    int32
	children [2]*Node
}

type AndTrieSimple struct {
	bit  int32
	root *Node
}

func NewAndTrie(upper int) *AndTrieSimple {
	return &AndTrieSimple{
		bit:  int32(bits.Len(uint(upper))),
		root: &Node{},
	}
}

func (bt *AndTrieSimple) Insert(num int) {
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
func (bt *AndTrieSimple) Remove(num int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		bit := (num >> i) & 1
		root = root.children[bit]
		root.count--
	}
}

func (bt *AndTrieSimple) Query(num int) (maxAnd int) {
	root := bt.root
	for i := bt.bit - 1; i >= 0; i-- {
		if root == nil {
			return
		}
		bit := (num >> i) & 1
		if bit == 1 {
			if c := root.children[1]; c != nil && c.count > 0 {
				maxAnd |= 1 << i
				root = root.children[1]
			} else {
				root = root.children[0]
			}
		} else {
			root = root.children[0]
		}
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
