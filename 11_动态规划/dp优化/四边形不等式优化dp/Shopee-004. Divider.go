package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	isDebug := false

	var n, k int
	var nums []int
	if isDebug {
		n, k = 4, 2
		nums = []int{1, 3, 2, 4}
	} else {
		fmt.Fscan(in, &n, &k)
		nums = make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &nums[i])
		}
	}

	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + nums[i]
	}

	// 记录每个分组对应的当前最优的遍历起始点
	prePos := make([]int, k+5)
	for i := 0; i < k+5; i++ {
		prePos[i] = i - 1
	}

	dp := make([][]int, n+5)
	for i := 0; i < n+5; i++ {
		dp[i] = make([]int, k+5)
		for j := 0; j < k+5; j++ {
			dp[i][j] = math.MaxInt64
		}
	}

	for i := 1; i <= n; i++ {
		maxG := min(k, i)
		for g := 1; g <= maxG; g++ {
			if g == 1 {
				dp[i][g] = preSum[i] * i
			} else {
				for pi := prePos[g]; pi < i; pi++ {
					cand := dp[pi][g-1] + (preSum[i]-preSum[pi])*(i-pi)
					if cand < dp[i][g] {
						dp[i][g] = cand
						prePos[g] = pi
					}
				}
			}
		}
	}

	fmt.Fprintln(out, dp[n][k])
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
