package main

import (
	"bufio"
	"fmt"
	"os"
)

// 解决dp`恰好选k个`的问题
// 需要高速化 dp[pos][使用次数] 的 dp 时,
// 如果dp(k+1)-dp(k)<=dp(k)-dp(k-1) ，则可以使用 wqs 二分.
// !问题转化为`每使用一次操作罚款 penalty 元,求最大分数`.
// 对penalty 二分搜索，转化为 dp[pos]一个维度的dp.
//
// k: 恰好选k个
// maxWeight: 最大权重
// !dp: func(penalty int) (cost, count int): 每使用一次操作罚款 penalty 元, 返回 [子问题dp的`最大值`, `最大的`操作使用次数]
//
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iv/solution/yi-chong-ji-yu-wqs-er-fen-de-you-xiu-zuo-x36r/
func AliensDp(k int, maxWeight int, getDp func(penalty int) (cost, count int)) int {
	left, right := 1, maxWeight
	penalty := 0

	for left <= right {
		mid := (left + right) >> 1
		// 如果操作次数大于等于 k，那么可以更新答案
		// 这里即使操作次数严格大于 k，更新答案也没有关系，因为总能二分到等于 k 的
		if _, count := getDp(mid); count >= k {
			penalty = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// penalty 时为0表示操作次数无限
	cost, _ := getDp(penalty)
	return cost + penalty*k
}

func main() {
	P1484()
}

// P1484 种树
// https://www.luogu.com.cn/problem/P1484
// 不相邻选k个点,使得这k个点的和最大.
// TODO: 有问题
func P1484() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	// 20 10
	// 1 -1 1 -1 1 -1 1 -1 1 -1 1 10 1 -1 1 -1 1 -1 1 -1
	// expected:18

	// 打家劫舍dp
	getDp := func(penalty int) (res, count int) {
		v0, c0 := 0, 0       // 不选
		v1, c1 := nums[0], 1 // 选
		for i := 1; i < n; i++ {
			score := nums[i]
			nv0, nc0 := v0, c0
			if v1 >= v0 {
				nv0, nc0 = v1, max(c0, c1)
			}
			nv1, nc1 := v0+score-penalty, max(c0, c1+1)
			v0, c0, v1, c1 = nv0, nc0, nv1, nc1
		}

		if v1 > v0 {
			return v1, c1
		}
		if v1 < v0 {
			return v0, c0
		}
		return v0, max(c0, c1)
	}

	maxWeight := maxs(nums)
	res := AliensDp(k, maxWeight, getDp)
	fmt.Fprintln(out, res)
}

// 714. 买卖股票的最佳时机含手续费
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iv/solution/yi-chong-ji-yu-wqs-er-fen-de-you-xiu-zuo-x36r/
// 给定一个整数数组 prices ，
// 它的第 i 个元素 prices[i] 是一支给定的股票在第 i 天的价格。
// 设计一个算法来计算你所能获取的最大利润。你最多可以完成 k 笔交易。
func maxProfit(k int, prices []int) int {
	getDp := func(penalty int) (res, count int) {
		v0, c0 := 0, 0          // 不持有股票
		v1, c1 := -prices[0], 0 // 持有股票
		for i := 1; i < len(prices); i++ {
			p := prices[i]
			nv0, nc0, nv1, nc1 := v0, c0, v1, c1
			if cand := v1 + p - penalty; cand >= v0 { // !注意取等号，让使用次数最大
				nv0, nc0 = cand, max(c0, c1+1)
			}
			if cand := v0 - p; cand >= v1 { // !注意取等号，让使用次数最大
				nv1, nc1 = cand, max(c0, c1)
			}
			v0, c0, v1, c1 = nv0, nc0, nv1, nc1
		}
		return v0, c0
	}

	maxWeight := int(1e18)
	return AliensDp(k, maxWeight, getDp)
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

func maxs(nums []int) int {
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > res {
			res = nums[i]
		}
	}
	return res
}
