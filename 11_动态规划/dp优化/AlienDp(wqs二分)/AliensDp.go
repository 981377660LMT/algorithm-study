package main

// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iv/solution/yi-chong-ji-yu-wqs-er-fen-de-you-xiu-zuo-x36r/
func maxProfit(k int, prices []int) int {
	n := len(prices)

	// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-with-transaction-fee/
	// 714. 买卖股票的最佳时机含手续费
	getDp := func(penalty int) [2]int {
		dp0, dp1 := [2]int{0, 0}, [2]int{-prices[0], 0} // dp1: 不持有股票, dp2: 持有股票
		for i := 1; i < n; i++ {
			ndp0, ndp1 := dp0, dp1
			if cand := dp1[0] + prices[i] - penalty; cand >= dp0[0] {
				ndp0 = [2]int{cand, max(dp0[1], dp1[1]+1)} // !注意让使用次数最大
			}
			if cand := dp0[0] - prices[i]; cand >= dp1[0] {
				ndp1 = [2]int{cand, max(dp0[1], dp1[1])} // !注意让使用次数最大
			}
			dp0, dp1 = ndp0, ndp1
		}

		return dp0
	}

	return AliensDp(n, k, getDp)
}

// 需要高速化 dp[pos][使用次数] 的 dp 时,
// 如果dp(k+1)-dp(k)<=dp(k)-dp(k-1) ，则可以使用 wqs 二分.
// !问题转化为`每使用一次操作罚款 penalty 元,求最大分数`.
// 对penalty 二分搜索，转化为 dp[pos]一个维度的dp.
//  !dp: func(penalty int) [2]int: 每使用一次操作罚款 penalty 元, 返回 [子问题dp的`最大值`, `最大的`操作使用次数]
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iv/solution/yi-chong-ji-yu-wqs-er-fen-de-you-xiu-zuo-x36r/
func AliensDp(n, k int, getDp func(penalty int) [2]int) int {
	if n == 0 {
		return 0
	}

	left, right := 1, int(1e18)
	penalty := 0

	for left <= right {
		mid := (left + right) >> 1
		// 如果操作次数大于等于 k，那么可以更新答案
		// 这里即使操作次数严格大于 k，更新答案也没有关系，因为总能二分到等于 k 的
		if cand := getDp(mid); cand[1] >= k {
			penalty = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// penalty 时为0表示操作次数无限
	res := getDp(penalty)
	return res[0] + penalty*k
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
