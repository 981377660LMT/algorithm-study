package main

import "sort"

// 查询区间[start,end)内等于value的元素个数.
func RangeFreq(nums []int) func(start, end int, value int) int {
	mp := map[int][]int{}
	for i, v := range nums {
		mp[v] = append(mp[v], i)
	}
	return func(start, end int, value int) int {
		pos := mp[value]
		return sort.SearchInts(pos, end) - sort.SearchInts(pos, start)
	}
}
