// SubarrayXorMax - 子数组异或最大值

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
)

func main() {

	{
		bruteForce := func(nums []int) (xor int, start, end int32) {
			for i := 0; i < len(nums); i++ {
				for j := i; j < len(nums); j++ {
					tmp := 0
					for k := i; k <= j; k++ {
						tmp ^= nums[k]
					}
					if tmp > xor {
						xor, start, end = tmp, int32(i), int32(j+1)
					}
				}
			}
			return
		}

		for i := 0; i < 100; i++ {
			nums := make([]int, 100)
			for j := 0; j < 100; j++ {
				nums[j] = rand.Intn(100)
			}
			res1, start1, end1 := SubarrayXorMax1(nums)
			res2, start2, end2 := bruteForce(nums)
			if res1 != res2 {
				fmt.Println(res1, res2)
				panic("error")
			}
			xor1, xor2 := 0, 0
			for i := int(start1); i < int(end1); i++ {
				xor1 ^= nums[i]
			}
			for i := int(start2); i < int(end2); i++ {
				xor2 ^= nums[i]
			}
			if xor1 != res1 || xor2 != res2 {
				fmt.Println(xor1, xor2, res1, res2)
				fmt.Println(start1, end1, start2, end2)
				panic("error")
			}
		}

		fmt.Println("ok")
	}

	{
		bruteForce := func(nums []int) []int32 {
			res := make([]int32, len(nums))
			for i := 0; i < len(nums); i++ {
				res[i] = int32(i) + 1
			}
			for i := 0; i < len(nums); i++ {
				xor, curMax := 0, 0
				for j := i; j < len(nums); j++ {
					xor ^= nums[j]
					if xor > curMax {
						curMax = xor
						res[i] = int32(j + 1)
					}
				}
			}
			return res
		}

		for i := 0; i < 100; i++ {
			nums := make([]int, 100)
			for j := 0; j < 100; j++ {
				nums[j] = rand.Intn(100)
			}
			res1 := SubarrayXorMax2(nums)
			res2 := bruteForce(nums)
			for i := 0; i < len(res1); i++ {
				if res1[i] != res2[i] {
					fmt.Println(res1[i], res2[i], res1, res2)
					panic("error")
				}
			}
		}

		fmt.Println("ok")
	}

}

// nums中子数组的异或最大值.
// 返回任意一个答案.
func SubarrayXorMax1(nums []int) (xor int, start, end int32) {
	max_ := 1
	preXor := make([]int, len(nums)+1)
	for i, num := range nums {
		max_ = max(max_, num)
		preXor[i+1] = preXor[i] ^ num
	}
	trie := NewXorTrieIndex(max_)
	for i := int32(0); i < int32(len(preXor)); i++ {
		trie.Insert(preXor[i], int32(i))
		best := trie.Query(preXor[i])
		if tmp := preXor[i] ^ preXor[best]; tmp > xor {
			xor, start, end = tmp, best, i
		}
	}
	return
}

// 对左端点为i的子数组，找到最近的右端点end，使得这个子数组的异或值最大.
func SubarrayXorMax2(nums []int) (ends []int32) {
	max_ := 1
	for _, num := range nums {
		max_ = max(max_, num)
	}
	trie := NewXorTrieIndex(max_)
	ends = make([]int32, len(nums))
	sufXor := 0
	for i := int32(len(nums) - 1); i >= 0; i-- {
		trie.Insert(sufXor, i)
		sufXor ^= nums[i]
		ends[i] = trie.Query(sufXor) + 1
	}
	return ends
}

type Node struct {
	index    int32
	children [2]*Node
}

func NewNode() *Node {
	return &Node{index: -1}
}

type XorTrieSimpleIndex struct {
	bit  int32
	root *Node
}

func NewXorTrieIndex(upper int) *XorTrieSimpleIndex {
	return &XorTrieSimpleIndex{
		bit:  int32(bits.Len(uint(upper))),
		root: NewNode(),
	}
}

func (bt *XorTrieSimpleIndex) Insert(num int, index int32) {
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
func (bt *XorTrieSimpleIndex) Query(num int) int32 {
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

func (bt *XorTrieSimpleIndex) Empty() bool {
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
