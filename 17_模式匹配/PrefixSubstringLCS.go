// Title: Prefix-Substring LCS
// !字符串s的前缀与字符串t的子串的最长公共子序列长度

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/prefix_substring_lcs
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	var s, t string
	fmt.Fscan(in, &s, &t)

	LCS := NewPrefixSubstringLCSWithString(s, t)
	for i := 0; i < q; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		fmt.Println(LCS.Query(a, b, c))
	}
}

type PrefixSubstringLCS struct {
	dp1 [][]uint32
}

func NewPrefixSubstringLCS(seq1, seq2 []int) *PrefixSubstringLCS {
	n1, n2 := uint32(len(seq1)), uint32(len(seq2))
	dp1, dp2 := make([][]uint32, n1+1), make([][]uint32, n1+1)
	for i := range dp1 {
		dp1[i] = make([]uint32, n2+1)
		dp2[i] = make([]uint32, n2+1)
	}
	for j := uint32(0); j < n2+1; j++ {
		dp1[0][j] = j
	}

	for i := uint32(1); i < n1+1; i++ {
		for j := uint32(1); j < n2+1; j++ {
			if seq1[i-1] == seq2[j-1] {
				dp1[i][j] = dp2[i][j-1]
				dp2[i][j] = dp1[i-1][j]
			} else {
				dp1[i][j] = maxUint32(dp1[i-1][j], dp2[i][j-1])
				dp2[i][j] = minUint32(dp1[i-1][j], dp2[i][j-1])
			}
		}
	}

	return &PrefixSubstringLCS{dp1: dp1}
}

func NewPrefixSubstringLCSWithString(s, t string) *PrefixSubstringLCS {
	ords1, ords2 := make([]int, len(s)), make([]int, len(t))
	for i := range s {
		ords1[i] = int(s[i])
	}
	for i := range t {
		ords2[i] = int(t[i])
	}
	return NewPrefixSubstringLCS(ords1, ords2)
}

// 求 seq1[:a) 和 seq2[b:c) 的最长公共子序列长度.
func (ps *PrefixSubstringLCS) Query(a int, b, c int) int {
	b32, c32 := uint32(b), uint32(c)
	res := 0
	for i := b32 + 1; i < c32+1; i++ {
		if ps.dp1[a][i] <= b32 {
			res++
		}
	}
	return res
}

func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func maxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
