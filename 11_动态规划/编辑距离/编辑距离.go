package main

// https://leetcode.cn/problems/edit-distance/
func minDistance(word1 string, word2 string) int {
	return EditDistance(word1, word2)
}

type Str = string

func EditDistance(s1, s2 Str) int {
	n1, n2 := int32(len(s1)), int32(len(s2))
	dp := make([][]int32, n1+1)
	for i := int32(0); i < n1+1; i++ {
		row := make([]int32, n2+1)
		row[0] = i
		dp[i] = row
	}
	for i := int32(0); i < n2+1; i++ {
		dp[0][i] = i
	}

	for i := int32(1); i < n1+1; i++ {
		for j := int32(1); j < n2+1; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min3(
					dp[i-1][j],
					dp[i][j-1],
					dp[i-1][j-1],
				) + 1
			}
		}
	}

	return int(dp[n1][n2])
}

func min3(a, b, c int32) int32 {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
