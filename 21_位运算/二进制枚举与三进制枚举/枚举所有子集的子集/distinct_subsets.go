package main

import (
	"cmp"
	"slices"
)

// 不重复地枚举所有子集.
func DistinctSubsets[S ~[]E, E cmp.Ordered](arr S, f func(subsetView S)) {
	arr = slices.Clone(arr)
	slices.Sort(arr)

	n := len(arr)
	subset := make(S, 0, n)

	var dfs func(start int)
	dfs = func(start int) {
		f(subset)
		for i := start; i < n; i++ {
			if i > start && arr[i] == arr[i-1] {
				continue
			}
			subset = append(subset, arr[i])
			dfs(i + 1)
			subset = subset[:len(subset)-1]
		}
	}
	dfs(0)
}

// https://leetcode.cn/problems/subsets-ii/
func subsetsWithDup(nums []int) [][]int {
	var res [][]int
	DistinctSubsets(nums, func(subset []int) {
		res = append(res, slices.Clone(subset))
	})
	return res
}
