package main

import "fmt"

func main() {
	s1, s2 := "banana", "ana"
	dp := AllPairLCP(s1, s2)
	fmt.Println(dp)
}

type Str = string

// 给定两个字符串 s1 和 s2，求它们的所有后缀的最长公共前缀.
// dp[i][j] 表示 s1[i:] 和 s2[j:] 的最长公共前缀.
func AllPairLCP(s1 Str, s2 Str) [][]int {
	n, m := len(s1), len(s2)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, m)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if s1[i] == s2[j] {
				if i+1 < n && j+1 < m {
					dp[i][j] = 1 + dp[i+1][j+1]
				} else {
					dp[i][j] = 1
				}
			}
		}
	}
	return dp
}
