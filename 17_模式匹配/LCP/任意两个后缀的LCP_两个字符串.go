package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// s1, s2 := "banana", "ana"
	// dp := AllPairLCP(s1, s2)
	// fmt.Println(dp)
	abc141()
}

// E - Who Says a Pun?
// https://atcoder.jp/contests/abc141/tasks/abc141_e
// 寻找出现次数>=2且不重叠的子串的最大长度.
// n<=5e3
func abc141() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	var s Str
	fmt.Fscan(in, &s)

	dp := AllPairLCP(s, s)
	res := int32(0)
	for i := int32(0); i < n; i++ {
		for j := i + 1; j < n; j++ {
			res = max32(res, min32(j-i, dp[i][j]))
		}
	}
	fmt.Fprintln(out, res)
}

type Str = string

// 给定两个字符串 s1 和 s2，求它们的所有后缀的最长公共前缀.
// dp[i][j] 表示 s1[i:] 和 s2[j:] 的最长公共前缀.
func AllPairLCP(s1 Str, s2 Str) [][]int32 {
	n, m := len(s1), len(s2)
	dp := make([][]int32, n)
	for i := range dp {
		dp[i] = make([]int32, m)
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
