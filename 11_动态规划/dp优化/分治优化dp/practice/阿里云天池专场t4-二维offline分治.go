// 阿里云天池专场
// 221021天池-04. 意外惊喜 https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/

// 比 2218. 从栈中取出K个硬币的最大面值和 https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
// !多了每个礼物包中的礼物价值 非严格递增 的条件，去除了背包总数的限制
// n<=2000 k<=1000 O(nklogn)

package main

func brilliantSurprise(present [][]int, lim int) (res int) {
	dp := make([]int, lim+1)
	var dfs func([][]int, []int)
	dfs = func(grid [][]int, sums []int) {
		if len(grid) == 1 {
			s := 0
			for i, v := range grid[0] {
				if i >= lim {
					break
				}
				s += v
				res = max(res, dp[lim-(i+1)]+s)
			}
			return
		}

		tmp := append([]int{}, dp...)

		m := len(grid) / 2
		for i, r := range grid[:m] {
			for j := lim; j >= len(r); j-- {
				dp[j] = max(dp[j], dp[j-len(r)]+sums[i])
			}
		}
		dfs(grid[m:], sums[m:])

		dp = tmp
		for i, r := range grid[m:] {
			for j := lim; j >= len(r); j-- {
				dp[j] = max(dp[j], dp[j-len(r)]+sums[m+i])
			}
		}
		dfs(grid[:m], sums[:m])
	}

	sums := make([]int, len(present))
	for i, p := range present {
		for _, v := range p {
			sums[i] += v
		}
	}
	dfs(present, sums)
	return
}

// 不优化的O(k*n^2)
// dp[k][j] 表示前k个礼物包中，取出j个礼物的最大价值
// dp[k][j] = max(dp[k-1][i]+sum(present[k][i+1:j+1])) (0<=i<j)
func brilliantSurprise2(present [][]int, limit int) int {
	n := len(present)
	dp := make([]int, limit+1)

	for k := 0; k < n; k++ {
		ndp := make([]int, limit+1)
		preSum := make([]int, len(present[k])+1)
		for i := 1; i <= len(present[k]); i++ {
			preSum[i] = preSum[i-1] + present[k][i-1]
		}

		for j := 0; j <= limit; j++ {
			for len_ := 0; len_ <= len(present[k]); len_++ { // 从第k个礼物包中取出len个礼物
				i := j - len_
				if i < 0 {
					break
				}
				ndp[j] = max(ndp[j], dp[i]+preSum[len_])
			}
		}
		dp = ndp
	}

	return dp[limit]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
