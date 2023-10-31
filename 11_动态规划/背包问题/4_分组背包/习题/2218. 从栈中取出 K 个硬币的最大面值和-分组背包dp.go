package main

// https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/description/
func maxValueOfCoins(piles [][]int, k int) int {
	n := len(piles)
	dp := make([]int, k+1)
	curCount := 0
	for i := 0; i < n; i++ {
		pile := piles[i]
		m := len(pile)
		for j := 1; j < m; j++ {
			pile[j] += pile[j-1]
		}
		curCount = min(k, curCount+m)
		for j := curCount; j > 0; j-- {
			max_ := 0
			for w := 0; w < min(m, j); w++ {
				max_ = max(max_, dp[j-w-1]+pile[w])
			}
			dp[j] = max(dp[j], max_)
		}
	}
	return dp[k]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
