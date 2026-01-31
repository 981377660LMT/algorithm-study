// Shopee-004. Divider
//  https://leetcode.cn/problems/VdG6tT/solution/si-bian-xing-bu-deng-shi-you-hua-dp-by-k-bp1t/

//  N个工程师分成 K 个组(N<=1e4,K<=1e2) 组内噪声
//  noise(l, r) = sum(A[l], A[l + 1], ..., A[r]) * (r - l + 1)
//  最小化总噪声

// 解:
//  记dp[k][j]为前j个人分k组的最小噪声值，
//  则dp[k][j] = min(dp[k-1][j]+(preSum[i]-preSum[j])*(i-j+1))
// n<=1e4,k<=1e2 O(nklogn)
// 满足四边形不等式的代价都满足决策单调性

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

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

	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + nums[i]
	}

	dp := divideAndConquerOptimization(k, n, func(i, j int) int {
		return (preSum[j] - preSum[i]) * (j - i)
	})
	fmt.Fprintln(out, dp[k][n])

}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func divideAndConquerOptimization(H, W int, dist func(i, j int) int) [][]int {
	dp := make([][]int, H+1)
	for i := range dp {
		dp[i] = make([]int, W+1)
		for j := range dp[i] {
			dp[i][j] = INF // !INF if get min
		}
	}
	dp[0][0] = 0

	for i := 1; i <= H; i++ {
		getCost := func(y, x int) int {
			if x >= y {
				return INF
			}
			return dp[i-1][x] + dist(x, y)
		}
		res := monotoneminima(W+1, W+1, getCost)
		for j := 0; j <= W; j++ {
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
