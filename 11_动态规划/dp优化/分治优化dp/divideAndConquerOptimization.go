// 単調最小値DP (aka. 分割統治DP) 优化 offlineDp
// https://ei1333.github.io/library/dp/divide-and-conquer-optimization.hpp
// !用于高速化 dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n) => !将区间[0,n)分成k组的最小代价
//  如果f满足决策单调性 那么对转移的每一行，可以采用 monotoneminima 寻找最值点
//  O(kn^2)优化到O(knlogn)

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

// CF833B-The Bakery
// https://www.luogu.com.cn/problem/CF833B
// 将一个数组分为k段，使得总价值最大。
// 一段区间的价值表示为区间内不同数字的个数。
// n=3e4,k<=50
//
// dp[i][j]=max{dp[i-1][k]+cost(k+1,j) 1<=k<j
// dp[i][j]意为前j个数被分成i段时的最大总价值.
//
// 线段树维护dp[i-1][k]的最大值.
// !如果一个值在区间中出现了很多次，那这些重复的数都不会产生贡献了
// O(nklogn)
//
// ps: 还可以决策单调性+主席树求cost
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

}

// !f(i,j,step): 左闭右开区间[i,j)的代价(0<=i<j<=n)
//
//	可选:step表示当前在第几组(1<=step<=k)
func divideAndConquerOptimization(k, n int, f func(i, j, step int) int) [][]int {
	dp := make([][]int, k+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0

	for k_ := 1; k_ <= k; k_++ {
		getCost := func(y, x int) int {
			if x >= y {
				return INF
			}
			return dp[k_-1][x] + f(x, y, k_)
		}
		res := monotoneminima(n+1, n+1, getCost)
		for j := 0; j <= n; j++ {
			dp[k_][j] = res[j][1]
		}
	}

	return dp
}

// 对每个 0<=i<H 求出 f(i,j) 取得最小值的 (j, f(i,j)) (0<=j<W)
func monotoneminima(H, W int, f func(i, j int) int) [][2]int {
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
			tmp := f(mid, i)
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
