package main

import "fmt"

func main() {
	s := "abab"
	fmt.Println(RangePalindrome(s))
}

// 统计区间内回文串个数
// 返回一个二维数组 dp, dp[i][j] 表示闭区间 [i,j] 内的回文串的个数
// https://codeforces.com/problemset/problem/245/H
func RangePalindrome(s string) [][]int {
	n := len(s)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
		dp[i][i] = 1
		if i+1 < n && s[i] == s[i+1] {
			dp[i][i+1] = 1
		}
	}
	for i := n - 3; i >= 0; i-- {
		for j := i + 2; j < n; j++ {
			if s[i] == s[j] {
				dp[i][j] = dp[i+1][j-1]
			}
		}
	}
	// 到这里为止，dp[i][j] = 1 表示 s[i:j+1] 是回文串
	for i := n - 2; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			dp[i][j] += dp[i][j-1] + dp[i+1][j] - dp[i+1][j-1] // 容斥
		}
	}
	return dp
}
