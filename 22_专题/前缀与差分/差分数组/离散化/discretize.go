package main

import "sort"

// (松)离散化.
//
//  offset: 离散化的起始值偏移量.
//
//	rank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeSparse(nums []int, offset int) (rank func(int) int, count int) {
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
	rank = func(x int) int { return sort.SearchInts(allNums, x) + offset }
	return
}

// (紧)离散化.
//
//  offset: 离散化的起始值偏移量.
//
//	rank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeCompressed(nums []int, offset int) (rank func(int) int, count int) {
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
	rank = func(v int) int { return mp[v] }
	count = len(allNums)
	return
}
