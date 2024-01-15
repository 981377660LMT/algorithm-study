// F - Hop Sugoroku
// !根号分治 + 懒更新

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc335f()
}

// F - Hop Sugoroku
// https://atcoder.jp/contests/abc335/tasks/abc335_f
// !加速:

// for i in range(1, n+1):
//
//	for j in range(i+A[i],n+1,A[i]):
//		dp[j] += dp[i]
func abc335f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	modAdd := func(a, b int) int {
		a += b
		if a >= MOD {
			a -= MOD
		}
		return a
	}

	THRESHOLD := 80 // !n=2e5, 80时较快
	res := 0
	dp := make([]int, n)
	dp[0] = 1
	lazy := make([][]int, THRESHOLD+1) // lazy[step][mod] 表示步长为step时，模为mod的所有数之和(分组前缀和).
	for i := range lazy {
		lazy[i] = make([]int, THRESHOLD+1)
	}

	for i := 0; i < n; i++ {
		for step := 1; step <= THRESHOLD; step++ {
			dp[i] = modAdd(dp[i], lazy[step][i%step])
		}

		res = modAdd(res, dp[i])

		step := nums[i]
		if step > THRESHOLD {
			for j := i + step; j < n; j += step {
				dp[j] = modAdd(dp[j], dp[i])
			}
		} else {
			lazy[step][i%step] = modAdd(lazy[step][i%step], dp[i])
		}
	}

	fmt.Fprintln(out, res)
}
