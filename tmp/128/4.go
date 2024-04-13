package main

import "math/bits"

// 给你一个 正 整数数组 nums 。

// 请你求出 nums 中有多少个子数组，满足子数组中 第一个 和 最后一个 元素都是这个子数组中的 最大 值。

func numberOfSubarrays(nums []int) int64 {
	mp := make(map[int][]int)
	for i, v := range nums {
		mp[v] = append(mp[v], i)
	}
	st := NewSparseTable(nums, max)
	res := 0
	for _, v := range mp {

	}
}

func NewSparseTable(nums []int, op func(int, int) int) (query func(int, int) int) {
	n := len(nums)
	size := bits.Len(uint(n))
	dp := make([][]int, size)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		dp[0][i] = nums[i]
	}

	for i := 1; i < size; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			dp[i][j] = op(dp[i-1][j], dp[i-1][j+(1<<(i-1))])
		}
	}

	query = func(left, right int) int {
		k := bits.Len(uint(right-left+1)) - 1
		return op(dp[k][left], dp[k][right-(1<<k)+1])
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
