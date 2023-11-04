package main

import "sort"

// (松)离散化.
//
//  offset: 离散化的起始值偏移量.
//
//	getRank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeSparse(nums []int, offset int) (getRank func(int) int, count int) {
	allNums := append(nums[:0:0], nums...)
	sort.Ints(allNums)
	slow := 0
	for fast := 1; fast < len(allNums); fast++ {
		if allNums[fast] != allNums[slow] {
			slow++
			allNums[slow] = allNums[fast]
		}
	}
	allNums = allNums[:slow+1]
	count = len(allNums)
	getRank = func(x int) int { return sort.SearchInts(allNums, x) + offset }
	return
}

// (紧)离散化.
//
//  offset: 离散化的起始值偏移量.
//
//	getRank: 给定一个数，返回它的排名`(offset ~ offset + count)`.
//	count: 离散化(去重)后的元素个数.
func DiscretizeCompressed(nums []int, offset int) (getRank func(int) int, count int) {
	allNums := append(nums[:0:0], nums...)
	sort.Ints(allNums)
	slow := 0
	for fast := 1; fast < len(allNums); fast++ {
		if allNums[fast] != allNums[slow] {
			slow++
			allNums[slow] = allNums[fast]
		}
	}
	allNums = allNums[:slow+1]
	mp := make(map[int]int, len(allNums))
	for i, v := range allNums {
		mp[v] = i + offset
	}
	getRank = func(v int) int { return mp[v] }
	count = len(allNums)
	return
}
