package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	分组背包问题()
}

// https://www.acwing.com/problem/content/9/
func 分组背包问题() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, maxCapacity int
	fmt.Fscan(in, &n, &maxCapacity)
	groups := make([][]item, n)
	for i := 0; i < n; i++ {
		var m int
		fmt.Fscan(in, &m)
		groups[i] = make([]item, m)
		for j := 0; j < m; j++ {
			var weight, value int
			fmt.Fscan(in, &weight, &value)
			groups[i][j] = item{value: value, weight: weight}
		}
	}

	res := GroupKnapsackAtMostOne(groups, maxCapacity)
	fmt.Fprintln(out, res)
}

type item struct{ value, weight int }

// 分组背包，每组最多选一个.最大化价值.
func GroupKnapsackAtMostOne(groups [][]item, maxCapacity int) int {
	dp := make([]int, maxCapacity+1)
	maxJ := 0
	for _, g := range groups {
		curMax := 0
		for _, w := range g {
			if w.weight > curMax {
				curMax = w.weight
			}
		}
		maxJ += curMax
		if maxJ > maxCapacity {
			maxJ = maxCapacity
		}
		// 这里 j 的初始值可以优化成前 i 个组的每组最大重量之和（但不能超过 maxW）
		for j := maxJ; j >= 0; j-- {
			for _, it := range g {
				if v, w := it.value, it.weight; w <= j {
					dp[j] = max(dp[j], dp[j-w]+v) // 如果 it.w 可能为 0 则需要用 dp[2][] 来滚动（或者保证每组至多一个 0 且 0 在该组最前面）
				}
			}
		}
	}
	return dp[maxCapacity]
}

// 分组背包，每组恰好选一个.最大化价值.
// 返回值：dp[j] 表示从每组恰好选一个，能否凑成重量 j.
// LC1981 https://leetcode-cn.com/problems/minimize-the-difference-between-target-and-chosen-elements/
func GroupKnapsackExactOne(groups [][]int, maxCapacity int) []bool {
	dp := make([]bool, maxCapacity+1) // dp[i][j] 表示能否从前 i 组物品中选出重量恰好为 j 的，且每组都恰好选一个物品
	dp[0] = true
	maxJ := 0
	for _, g := range groups {
		curMax := 0
		for _, w := range g {
			if w > curMax {
				curMax = w
			}
		}
		maxJ += curMax
		if maxJ > maxCapacity {
			maxJ = maxCapacity
		}
		// 这里 j 的初始值可以优化成前 i 个组的每组最大重量之和（但不能超过 maxW）
		for j := maxJ; j >= 0; j-- {
			ok := false
			for _, w := range g {
				if w <= j && dp[j-w] {
					dp[j] = true
					ok = true
					break
				}
			}
			if !ok {
				dp[j] = false
			}
		}
	}
	return dp
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

// https://leetcode-cn.com/problems/minimize-the-difference-between-target-and-chosen-elements/
func minimizeTheDifference(mat [][]int, target int) int {
	maxSum := 0
	for _, row := range mat {
		max := 0
		for _, v := range row {
			if v > max {
				max = v
			}
		}
		maxSum += max
	}

	ok := GroupKnapsackExactOne(mat, maxSum)
	res := math.MaxInt
	for cand := range ok {
		if ok[cand] {
			res = min(res, abs(cand-target))
		}
	}
	return res
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
