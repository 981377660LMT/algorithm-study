package main

import (
	"fmt"
	"math/bits"
	"sort"
)

// nums = [2,8,4,32,16,1], queries = [[0,2],[1,4],[0,5]]

func main() {
	// fmt.Println(maximumSubarrayXor([]int{2, 8, 4, 32, 16, 1}, [][]int{{0, 2}, {1, 4}, {0, 5}}))
	// nums = [0,7,3,2,8,5,1], queries = [[0,3],[1,5],[2,4],[2,6],[5,6]]
	// fmt.Println(maximumSubarrayXor([]int{0, 7, 3, 2, 8, 5, 1}, [][]int{{0, 3}, {1, 5}, {2, 4}, {2, 6}, {5, 6}}))
	// nums = [2,8,4,32,16,1], queries = [[0,2],[1,4],[0,5]]
	fmt.Println(maximumSubarrayXor([]int{2, 8, 4, 32, 16, 1}, [][]int{{0, 2}, {1, 4}, {0, 5}}))
	fmt.Println(2 ^ 8 ^ 4)
}

// 给你一个由 n 个整数组成的数组 nums，以及一个大小为 q 的二维整数数组 queries，其中 queries[i] = [li, ri]。

// 对于每一个查询，你需要找出 nums[li..ri] 中任意 子数组 的 最大异或值。

// 数组的异或值 需要对数组 a 反复执行以下操作，直到只剩一个元素，剩下的那个元素就是 异或值：

// 对于除最后一个下标以外的所有下标 i，同时将 a[i] 替换为 a[i] XOR a[i + 1] 。
// 移除数组的最后一个元素。
// 返回一个大小为 q 的数组 answer，其中 answer[i] 表示查询 i 的答案。

// 就是开头结尾两个数
func maximumSubarrayXor(nums []int, queries [][]int) []int {
	n := int32(len(nums))
	groupByLeft := make([][]int32, n)
	for qi, query := range queries {
		groupByLeft[query[0]] = append(groupByLeft[query[0]], int32(qi))
	}

	// 奇数长度，首尾
	// 偶数长度，整个
	res := make([]int, len(queries))
	upper := maxs(nums...) + 1
	for left, qids := range groupByLeft {
		sort.Slice(qids, func(i, j int) bool { return queries[qids[i]][1] < queries[qids[j]][1] })
		X0 := [2]*XORTrieSimple{NewXORTrie(upper), NewXORTrie(upper)}
		X1 := [2]*XORTrieSimple{NewXORTrie(upper), NewXORTrie(upper)}
		_ = X1
		ptr := left
		curMax := 0
		preXor := 0
		for _, qid := range qids {
			for ptr <= queries[qid][1] {
				preXor ^= nums[ptr]
				X0[ptr&1].Insert(nums[ptr])
				X1[ptr&1].Insert(preXor)
				curMax = max(curMax, X0[ptr&1].Query(nums[ptr]))
				curMax = max(curMax, X1[ptr&1].Query(preXor))
				curMax = max(curMax, nums[ptr])
				if ptr&1 != left&1 {
					curMax = max(curMax, preXor)
				}
				ptr++
			}
			res[qid] = curMax
		}
	}
	return res

}

type Node struct {
	count    int32
	children [2]*Node // 数组比 left,right 更快
}

type XORTrieSimple struct {
	bit  int32
	root *Node
}

func NewXORTrie(upper int) *XORTrieSimple {
	return &XORTrieSimple{
		bit:  int32(bits.Len(uint(upper))),
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
		if root == nil {
			return
		}
		bit := (num >> i) & 1
		if root.children[bit^1] != nil && root.children[bit^1].count > 0 {
			maxXor |= 1 << i
			bit ^= 1
		}
		root = root.children[bit]
	}
	return
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
