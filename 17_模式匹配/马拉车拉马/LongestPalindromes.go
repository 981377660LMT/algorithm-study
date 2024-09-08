package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// yosupo()
	xmascontest2015noon()
}

func yosupo() {
	// https://judge.yosupo.jp/problem/enumerate_palindromes
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	res := LongestPalindromesLength(int32(len(s)), func(i, j int32) bool { return s[i] == s[j] })
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// C - Colored Tiles
// https://atcoder.jp/contests/xmascontest2015noon/tasks/xmascontest2015_c
func xmascontest2015noon() {}

// 对2*n-1个回文中心, 求出每个中心对应的极大回文子串的长度.
// aaaaa -> 1 2 3 4 5 4 3 2 1 奇偶交替.
func LongestPalindromesLength(n int32, equals func(i, j int32) bool) []int32 {
	res := make([]int32, 2*n-1)
	palindromes := LongestPalindromes(n, equals)
	for _, p := range palindromes {
		s, e := p[0], p[1]
		res[s+e-1] = e - s
	}
	return res
}

// 给定一个字符串，返回极长回文子串的区间.这样的极长回文子串最多有 2n-1 个.
// ManacherSimple.
func LongestPalindromes(n int32, equals func(i, j int32) bool) [][2]int32 {
	f := func(i, j int32) bool {
		if i > j {
			return false
		}
		if i&1 == 1 {
			return true
		}
		return equals(i>>1, j>>1)
	}
	dp := make([]int32, 2*n-1)
	i, j := int32(0), int32(0)
	for i < 2*n-1 {
		for i-j >= 0 && i+j < 2*n-1 && f(i-j, i+j) {
			j++
		}
		dp[i] = j
		k := int32(1)
		for i-k >= 0 && i+k < 2*n-1 && k < j && dp[i-k] != j-k {
			dp[i+k] = min32(j-k, dp[i-k])
			k++
		}
		i += k
		j = max32(j-k, 0)
	}

	res := make([][2]int32, 0, len(dp))
	for i := int32(0); i < int32(len(dp)); i++ {
		if dp[i] == 0 {
			continue
		}
		l := (i - dp[i] + 2) / 2
		r := (i + dp[i] - 1) / 2
		if l <= r {
			res = append(res, [2]int32{l, r + 1})
		}
	}
	res = res[:len(res):len(res)]
	return res
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
