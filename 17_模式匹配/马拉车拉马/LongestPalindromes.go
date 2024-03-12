package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/enumerate_palindromess
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	res := LongestPalindromesLength(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

const INF int32 = 1e9 + 10

// 对2*n-1个回文中心, 求出每个中心对应的极大回文子串的长度.
func LongestPalindromesLength(n int32, f func(i int32) int32) []int32 {
	res := make([]int32, 2*n-1)
	palindromes := LongestPalindromes(n, f)
	for _, p := range palindromes {
		s, e := p[0], p[1]
		res[s+e-1] = e - s
	}
	return res
}

// 给定一个字符串，返回极长回文子串的区间.这样的极长回文子串最多有 2n-1 个.
func LongestPalindromes(n int32, f func(i int32) int32) [][2]int32 {
	m := n*2 - 1
	sb := make([]int32, m)
	for i := n - 1; i >= 0; i-- {
		sb[2*i] = f(i)
	}
	for i := int32(0); i < n-1; i++ {
		sb[2*i+1] = INF
	}
	dp := make([]int32, m)
	i, j := int32(0), int32(0)
	for i < m {
		for i-j >= 0 && i+j < m && sb[i-j] == sb[i+j] {
			j++
		}
		dp[i] = j
		k := int32(1)
		for i-k >= 0 && i+k < m && k+dp[i-k] < j {
			dp[i+k] = dp[i-k]
			k++
		}
		i += k
		j -= k
	}
	for i := int32(0); i < m; i++ {
		if ((i ^ dp[i]) & 1) == 0 {
			dp[i]--
		}
	}
	res := make([][2]int32, 0, m)
	for i := int32(0); i < m; i++ {
		if dp[i] == 0 {
			continue
		}
		start := (i - dp[i] + 1) / 2
		end := (i + dp[i] + 1) / 2
		res = append(res, [2]int32{start, end})
	}
	return res
}
