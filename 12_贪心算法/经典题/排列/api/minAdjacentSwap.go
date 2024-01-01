// 邻位交换的最小次数
// https://leetcode.cn/problems/minimum-adjacent-swaps-to-reach-the-kth-smallest-number/

package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(MinAdjacentSwapNaive([]int{1, 3, 5, 4}, []int{1, 5, 4, 3}))
	fmt.Println(MinAdjacentSwap([]int{1, 3, 5, 4}, []int{1, 5, 4, 3}))
}

// 求使两个数组相等的最少邻位交换次数.如果无法完成，返回-1.
// 对每个数，贪心找到对应的最近位置交换.
// 时间复杂度`O(n^2)`.
func MinAdjacentSwapNaive(nums1, nums2 []int) int {
	indexOf := func(arr []int, target int) int {
		for index, num := range arr {
			if num == target {
				return index
			}
		}
		return -1
	}

	res := 0
	for _, num := range nums1 {
		index := indexOf(nums2, num)
		if index == -1 {
			return -1
		}
		res += index
		nums2 = append(nums2[:index], nums2[index+1:]...)
	}
	return res
}

// 求使两个数组相等的最少邻位交换次数.如果无法完成，返回-1.
// 映射+求逆序对.
// 时间复杂度`O(nlogn)`.
func MinAdjacentSwap(nums1, nums2 []int) int {
	mapping := make(map[int][]int)
	for index, num := range nums2 {
		mapping[num] = append(mapping[num], index)
	}
	for index, num := range nums1 {
		mapped := mapping[num][0]
		nums1[index] = mapped
		mapping[num] = mapping[num][1:]
	}
	return CountInversion(nums1)
}

func CountInversion(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	max_ := maxs(nums) + 1
	if max_ > 2e6 {
		nums = Discretize(nums)
		max_ = maxs(nums) + 1
	}
	res := 0
	bit := NewB(max_)
	for _, v := range nums {
		res += bit.QueryRange(v+1, max_)
		bit.Add(v, 1)
	}
	return res
}

func Discretize(nums []int) []int {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	mp := make(map[int]int, len(allNums))
	for i, v := range allNums {
		mp[v] = i
	}
	newNums := make([]int, len(nums))
	for i, v := range nums {
		newNums[i] = mp[v]
	}
	return newNums
}

func maxs(nums []int) int {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

type B struct {
	n    int
	data []int
}

func NewB(n int) *B {
	res := &B{n: n, data: make([]int, n)}
	return res
}

func (b *B) Add(index int, v int) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *B) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *B) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}
