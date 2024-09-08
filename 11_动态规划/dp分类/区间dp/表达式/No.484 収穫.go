// No.484 収穫
// https://yukicoder.me/problems/no/484
// 给定n个果实，每个果实在时间A[i]可以收获.
// 起点任意，每次可以移动到相邻位置或者停留在原地.
// 问收获所有果实的最短时间.
//
// !dp[isRight][l][r] 表示区间 [l,r] 未收获，当前位置在 isRight (isRight=0 时在 l−1，1 时在 r+1) 的最短时间
// 每次考虑移动到l-1和r+1的情况，取最小值

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

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	memo := make([]int, n*n*2)
	for i := range memo {
		memo[i] = -1
	}
	var dfs func(l, r, isRight int) int
	dfs = func(l, r, dir int) int {
		if l == 0 && r == n-1 {
			if dir == 0 {
				return nums[0]
			} else {
				return nums[n-1]
			}
		}

		hash := l*n*2 + r*2 + dir
		if memo[hash] != -1 {
			return memo[hash]
		}

		res := INF
		if dir == 0 {
			if l > 0 {
				res = min(res, dfs(l-1, r, 0)+1)
			}
			if r < n-1 {
				res = min(res, dfs(l, r+1, 1)+r-l+1)
			}
			res = max(res, nums[l])
		} else {
			if l > 0 {
				res = min(res, dfs(l-1, r, 0)+r-l+1)
			}
			if r < n-1 {
				res = min(res, dfs(l, r+1, 1)+1)
			}
			res = max(res, nums[r])
		}
		memo[hash] = res
		return res
	}

	res := INF
	for i := 0; i < n; i++ {
		res = min(res, dfs(i, i, 0))
	}
	fmt.Fprintln(out, res)
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
