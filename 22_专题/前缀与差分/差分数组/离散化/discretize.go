package main

import (
	"fmt"
	"sort"
)

// (松)离散化.
//
//	 offset: 离散化的起始值偏移量.
//
//		getRank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//		count: 离散化(去重)后的元素个数.
func DiscretizeSparse(nums []int, offset int) (getRank func(int) int, count int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	allNums := make([]int, 0, count)
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	getRank = func(x int) int { return sort.SearchInts(allNums, x) + offset }
	return
}

// (紧)离散化.
//
//	 offset: 离散化的起始值偏移量.
//
//		getRank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//		count: 离散化(去重)后的元素个数.
func DiscretizeCompressed(nums []int, offset int) (getRank func(int) int, count int) {
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
		mp[v] = i + offset
	}
	getRank = func(v int) int { return mp[v] }
	count = len(allNums)
	return
}

// 不带相同值的离散化，转换为 0-n-1.
// rank: 离散化后的排名.
// keys: keys[ranks[i]] = nums[i].
func DiscretizeUnique(nums []int) (rank []int, keys []int) {
	rank = argSort(nums)
	keys = reArrage(nums, rank)
	rank = argSort(rank)
	return
}

func argSort(nums []int) []int {
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reArrage(nums []int, order []int) []int {
	res := make([]int, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func main() {
	fmt.Println(DiscretizeUnique([]int{3, 2, 1, 3, 2, 1}))
}
