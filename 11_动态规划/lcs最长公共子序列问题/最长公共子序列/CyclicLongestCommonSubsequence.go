// 循环位移最长公共子序列

package main

import "fmt"

func main() {
	s1 := "abcde"
	s2 := "axb"
	fmt.Println(CyclicLongestCommonSubsequence(s1, s2))
}

func CyclicLongestCommonSubsequence(s, t string) int {
	n, m := len(s), len(t)
	dp := make([][]int, n*2+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	from := make([][]int, n*2+1)
	for i := range from {
		from[i] = make([]int, m+1)
	}
	eq := func(a, b int) bool {
		return s[(a-1)%n] == t[(b-1)%m]
	}

	for i := 0; i <= n*2; i++ {
		for j := 0; j <= m; j++ {
			dp[i][j] = 0
			if j > 0 && dp[i][j-1] > dp[i][j] {
				dp[i][j] = dp[i][j-1]
				from[i][j] = 0
			}
			if i > 0 && j > 0 && eq(i, j) && dp[i-1][j-1]+1 > dp[i][j] {
				dp[i][j] = dp[i-1][j-1] + 1
				from[i][j] = 1
			}
			if i > 0 && dp[i-1][j] > dp[i][j] {
				dp[i][j] = dp[i-1][j]
				from[i][j] = 2
			}
		}
	}

	res := 0
	for i := 0; i < n; i++ {
		res = max(res, dp[i+n][m])
		x, y := i+1, 0
		for y <= m && from[x][y] == 0 {
			y++
		}
		for ; y <= m && x <= n*2; x++ {
			from[x][y] = 0
			dp[x][m]--
			if x == n*2 {
				break
			}
			for ; y <= m; y++ {
				if from[x+1][y] == 2 {
					break
				}
				if y+1 <= m && from[x+1][y+1] == 1 {
					y++
					break
				}
			}
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
