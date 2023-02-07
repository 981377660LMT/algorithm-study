// dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n)
// !将区间[0,n)分成k组的最小代价
// !f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
// 时间复杂度O(k*n^2)
// 可用分治/CHT优化到O(k*nlogn)

package main

import "math/bits"

const INF int = 1e18

// !f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
func offlineDpMin(k, n int, f func(i, j int) int) int {
	dp := make([]int, n+1)
	for j := 1; j <= n; j++ {
		dp[j] = f(0, j)
	}
	for k_ := 2; k_ <= k; k_++ {
		ndp := make([]int, n+1)
		for i := 0; i <= n; i++ {
			ndp[i] = INF
		}
		for j := 1; j <= n; j++ {
			for i := 0; i < j; i++ {
				cand := dp[i] + f(i, j)
				if cand < ndp[j] {
					ndp[j] = cand
				}
			}
		}
		dp = ndp
	}
	return dp[n]
}

// !f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
func offlineDpMax(k, n int, f func(i, j int) int) int {
	dp := make([]int, n+1)
	for j := 1; j <= n; j++ {
		dp[j] = f(0, j)
	}
	for k_ := 2; k_ <= k; k_++ {
		ndp := make([]int, n+1)
		for i := 0; i <= n; i++ {
			ndp[i] = -INF
		}
		for j := 1; j <= n; j++ {
			for i := 0; i < j; i++ {
				cand := dp[i] + f(i, j)
				if cand > ndp[j] {
					ndp[j] = cand
				}
			}
		}
		dp = ndp
	}
	return dp[n]
}

// 1959. K 次调整数组大小浪费的最小总空间
// https://leetcode.cn/problems/minimum-total-space-wasted-with-k-resizing-operations/submissions/
func minSpaceWastedKResizing(nums []int, k int) int {
	st := NewSparseTable(nums, max)
	preSum := make([]int, len(nums)+1)
	for i, num := range nums {
		preSum[i+1] = preSum[i] + num
	}
	return offlineDpMin(k+1, len(nums), func(i, j int) int {
		max := st(i, j-1)
		return max*(j-i) - (preSum[j] - preSum[i])
	})
}

func NewSparseTable(nums []int, op func(int, int) int) (query func(int, int) int) {
	n := len(nums)
	size := bits.Len(uint(n))
	dp := make([][]int, size)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		dp[0][i] = nums[i]
	}

	for i := 1; i < size; i++ {
		for j := 0; j+(1<<i) <= n; j++ {
			dp[i][j] = op(dp[i-1][j], dp[i-1][j+(1<<(i-1))])
		}
	}

	query = func(left, right int) int {
		k := bits.Len(uint(right-left+1)) - 1
		return op(dp[k][left], dp[k][right-(1<<k)+1])
	}

	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
