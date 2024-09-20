package main

import (
	"strconv"
)

// https://leetcode.cn/problems/summary-ranges/description/
func summaryRanges(nums []int) []string {
	var res []string
	EnumerateConsecutiveIntervals(
		int32(len(nums)), func(i int32) int { return nums[i] },
		func(min, max int, isIn bool) {
			if !isIn {
				return
			}
			if min == max {
				res = append(res, strconv.Itoa(min))
			} else {
				res = append(res, strconv.Itoa(min)+"->"+strconv.Itoa(max))
			}
		},
	)
	return res
}

// 遍历连续区间/合并连续区间.
func EnumerateConsecutiveIntervals(
	n int32, supplier func(i int32) int,
	consumer func(min, max int, isIn bool),
) {
	if n == 0 {
		return
	}
	i := int32(0)
	for i < n {
		start := i
		for i < n-1 && supplier(i)+1 == supplier(i+1) {
			i++
		}
		consumer(supplier(start), supplier(i), true)
		if i+1 < n {
			consumer(supplier(i)+1, supplier(i+1)-1, false)
		}
		i++
	}
}
