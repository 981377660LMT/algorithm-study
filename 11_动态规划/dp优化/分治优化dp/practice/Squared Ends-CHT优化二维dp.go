// Squared Ends
// https://csacademy.com/contest/round-70/task/squared-ends/
// 给定一个数组,分成k个子数组,
// 子数组[A[l],...,A[r]]的代价为(A[l]-A[r])^2
// 最小化所有子数组的代价和
// n<=1e4 k<=100
// CHT或者分治优化dp都可以解决这个问题

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	dp := divideAndConquerOptimization(k, n, func(i, j int) int {
		return (nums[j-1] - nums[i]) * (nums[j-1] - nums[i])
	})
	fmt.Fprintln(out, dp[k][n])
}

const INF int = 1e18

// !dist(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
func divideAndConquerOptimization(k, n int, dist func(i, j int) int) [][]int {
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
			return dp[k_-1][x] + dist(x, y)
		}
		res := monotoneminima(n+1, n+1, getCost)
		for j := 0; j <= n; j++ {
			dp[k_][j] = res[j][1]
		}
	}

	return dp
}

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
