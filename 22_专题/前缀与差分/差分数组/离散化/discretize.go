package main

import "sort"

// (松)离散化.
//
//	rank: 给定一个数，返回它的排名`(0-count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeSparse(nums []int) (rank func(int) int, count int) {
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
	rank = func(x int) int { return sort.SearchInts(allNums, x) }
	return
}

// (紧)离散化.
//  rank: 给定一个在 nums 中的值,返回它的排名(0~len(rank)-1).
//  newNums: 离散化后的数组.
func discretize(nums []int) (rank map[int]int, newNums []int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	rank = make(map[int]int, len(allNums))
	for i, v := range allNums {
		rank[v] = i
	}
	newNums = make([]int, len(nums))
	for i, v := range nums {
		newNums[i] = rank[v]
	}
	return
}
