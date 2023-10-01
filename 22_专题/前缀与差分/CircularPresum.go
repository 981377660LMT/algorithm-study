package main

import "sort"

// 环形前缀和.
func CircularPresum(nums []int) func(start int, end int) int {
	n := len(nums)
	preSum := make([]int, n+1)
	for i, v := range nums {
		preSum[i+1] = preSum[i] + v
	}

	cal := func(r int) int {
		return preSum[n]*(r/n) + preSum[r%n]
	}

	query := func(start int, end int) int {
		if start >= end {
			return 0
		}
		return cal(end) - cal(start)
	}

	return query
}

const INF int = 1e18

// 100076. 无限数组的最短子数组
// https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
// 求循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
// 1 <= nums.length <= 1e5
// 1 <= nums[i] <= 1e5
// 1 <= target <= 1e9
func minSizeSubarray(nums []int, target int) int {
	Q := CircularPresum(nums)
	res := INF
	for start := range nums {
		cand := sort.Search(1e9+10, func(mid int) bool {
			return Q(start, start+mid) >= target
		})
		if Q(start, start+cand) == target {
			res = min(res, cand)
		}
	}

	if res == INF {
		return -1
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
