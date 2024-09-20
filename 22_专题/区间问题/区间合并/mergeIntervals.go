package main

import "sort"

// 56. 合并区间
// https://leetcode.cn/problems/merge-intervals/description/
func merge(intervals [][]int) (res [][]int) {
	MergeIntervals(
		int32(len(intervals)), func(i int32) (int, int) { return intervals[i][0], intervals[i][1] },
		func(l, r int) { res = append(res, []int{l, r}) },
	)
	return
}

// 合并所有重叠的闭区间，返回一个不重叠的区间列表.
func MergeIntervals(n int32, supplier func(int32) (int, int), consumer func(int, int)) {
	if n == 0 {
		return
	}
	order := make([]int32, n)
	for i := int32(0); i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		l1, _ := supplier(order[i])
		l2, _ := supplier(order[j])
		return l1 < l2
	})
	preL, preR := supplier(order[0])
	for _, i := range order[1:] {
		curL, curR := supplier(i)
		if curL <= preR {
			preR = max(preR, curR)
		} else {
			consumer(preL, preR)
			preL, preR = curL, curR
		}
	}
	consumer(preL, preR)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
