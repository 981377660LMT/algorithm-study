// !dp[i]=min(dp[j]+(preSum[i]-preSum[j])**2)+M
// 诗人小G/打印文章

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

	for {
		var n, k int
		fmt.Fscan(in, &n, &k)
		if n == 0 && k == 0 {
			break
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &nums[i])
		}

		// !dp[i]=min(dp[j]+(preSum[i]-preSum[j])**2)+M
		preSum := make([]int, n+1)
		for i := 0; i < n; i++ {
			preSum[i+1] = preSum[i] + nums[i]
		}

		res := offlineOnlineDp(n, func(i, j int) int {
			// 0<=i<j<=n
			return (preSum[j]-preSum[i])*(preSum[j]-preSum[i]) + k
		})
		fmt.Fprintln(out, res)
	}

}

// !dist(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
func offlineOnlineDp(n int, dist func(i, j int) int) int {
	dp := make([]int, n+1)
	used := make([]bool, n+1)
	used[n] = true

	update := func(k, val int) {
		if !used[k] {
			dp[k] = val
		}
		dp[k] = min(dp[k], val) // min if get min
		used[k] = true
	}

	var dfs func(top, bottom, left, right int) // induce
	dfs = func(top, bottom, left, right int) {
		if top == bottom {
			return
		}
		mid := (top + bottom) / 2
		index := left
		res := dist(mid, index) + dp[index]
		for i := left; i <= right; i++ {
			tmp := dist(mid, i) + dp[i]
			if tmp < res { // !less if get min
				res = tmp
				index = i
			}
		}

		update(mid, res)
		dfs(top, mid, left, index)
		dfs(mid+1, bottom, index, right)
	}

	var solve func(left, right int)
	solve = func(left, right int) {
		if left+1 == right {
			update(left, dist(left, right)+dp[right])
			return
		}
		mid := (left + right) / 2
		solve(mid, right)
		dfs(left, mid, mid, right)
		solve(left, mid)
	}

	solve(0, n)
	return dp[0]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
