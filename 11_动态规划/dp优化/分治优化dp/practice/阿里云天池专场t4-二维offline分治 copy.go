// 阿里云天池专场
// 221021天池-04. 意外惊喜 https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/

// 比 2218. 从栈中取出K个硬币的最大面值和 https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
// !多了每个礼物包中的礼物价值 非严格递增 的条件，去除了背包总数的限制
// n<=2000 k<=1000 O(nklogn)
// TODO
package main

import "fmt"

func main() {
	// [[1,2],[2,3],[3,4]], limit = 3
	fmt.Println(brilliantSurprise([][]int{{1, 2}, {2, 3}, {3, 4}}, 3))
}

func brilliantSurprise(present [][]int, lim int) int {
	preSums := make([][]int, len(present))
	for i := range preSums {
		preSums[i] = make([]int, len(present[i])+1)
		for j := range present[i] {
			preSums[i][j+1] = preSums[i][j] + present[i][j]
		}
	}

	var divideAndConquerOptimization func(k, n int, dist func(i, j, step int) int) [][]int
	//  !dist(i,j,step): 左闭右开区间[i,j)的代价(0<=i<j<=n) step表示当前在第几组(1<=step<=k)
	divideAndConquerOptimization = func(k, n int, dist func(i, j, step int) int) [][]int {
		dp := make([][]int, k+1)
		for i := range dp {
			dp[i] = make([]int, n+1)
			for j := range dp[i] {
				dp[i][j] = INF // !INF if get min
			}
		}
		dp[0][0] = 0

		for k_ := 1; k_ <= k; k_++ {
			getCost := func(y, x int) int {
				if x >= y {
					return INF
				}
				return dp[k_-1][x] + dist(x, y, k_)
			}
			res := monotoneminima(n+1, n+1, getCost)
			for j := 0; j <= n; j++ {
				dp[k_][j] = res[j][1]
			}
		}

		return dp
	}
	res := divideAndConquerOptimization(len(present), lim, func(i, j, step int) int {
		return preSums[step-1][j] - preSums[step-1][i]
	})
	return res[len(present)][lim]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const INF int = 1e18

// 对每个 0<=i<H 求出 dist(i,j) 取得最小值的 (j, dist(i,j)) (0<=j<W)
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
