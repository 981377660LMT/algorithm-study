// 単調最小値DP (aka. 分割統治DP) 优化 offlineDp
// https://ei1333.github.io/library/dp/divide-and-conquer-optimization.hpp
// !用于高速化 dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n) => !将区间[0,n)分成k组的最小代价
//  如果f满足决策单调性 那么对转移的每一行，可以采用 monotoneminima 寻找最值点
//  O(kn^2)优化到O(knlogn)

package main

import (
	"sort"
)

//
//
//
// 给你一个房屋数组houses 和一个整数 k ，
// 其中 houses[i] 是第 i 栋房子在一条街上的位置，
// 现需要在这条街上安排 k 个邮筒。
// 请你返回每栋房子与离它最近的邮筒之间的距离的 最小 总和。
// !dp[k][n] 表示将[1,n]分成k段时的最优解
// !dp[k][j]=min(dp[k-1][i]+f(i,j)) (i<j) f是一个Monotone函数
func minDistance(houses []int, k int) int {
	sort.Ints(houses)

	// 求距离的函数
	n := len(houses)
	memo := make([][]int, n+1)
	for i := range memo {
		memo[i] = make([]int, n+1)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}
	var dist func(i, j int) int
	dist = func(i, j int) int {
		if i >= j {
			return 0
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		res := houses[j] - houses[i] + dist(i+1, j-1)
		memo[i][j] = res
		return res
	}

	dist2 := func(i, j int) int {
		return dist(i, j-1)
	}
	dp := divideAndConquerOptimization(k, n, dist2)
	return dp[k][n]
}

const INF int = 1e18

//  !dist(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
func divideAndConquerOptimization(k, n int, dist func(i, j int) int) [][]int {
	dp := make([][]int, k+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		for j := range dp[i] {
			dp[i][j] = INF // !INF if get min
		}
	}
	dp[0][0] = 0

	for i := 1; i <= k; i++ {
		getCost := func(y, x int) int {
			if x >= y {
				return INF
			}
			return dp[i-1][x] + dist(x, y)
		}
		res := monotoneminima(n+1, n+1, getCost)
		for j := 0; j <= n; j++ {
			dp[i][j] = res[j][1]
		}
	}

	return dp
}

func monotoneminima(H, W int, dist func(i, j int) int) [][2]int {
	dp := make([][2]int, H) // dp[i] 表示第i行取到`最小值`的(索引,值)

	var dfs func(top, bottom, left, right int)
	dfs = func(top, bottom, left, right int) {
		if top > bottom {
			return
		}

		mid := (top + bottom) / 2
		index := -1
		res := 0
		for i := left; i <= right; i++ {
			tmp := dist(mid, i)
			if index == -1 || tmp < res { // !less if get min
				index = i
				res = tmp
			}
		}
		dp[mid] = [2]int{index, res}
		dfs(top, mid-1, left, index)
		dfs(mid+1, bottom, index, right)
	}

	dfs(0, H-1, 0, W-1)
	return dp
}
