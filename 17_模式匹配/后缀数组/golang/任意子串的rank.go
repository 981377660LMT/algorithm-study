// 任意子串的rank(子串排序)

package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func main() {
	s := "banana"
	dp := SortSubstrings(s, -1)
	fmt.Println(dp)
}

// dp[i][j]：S[i:i+j) 的 rank, -1 表示空串
//  maxLen: 子串最大长度，-1 表示最长为len(s)
func SortSubstrings(s string, maxLen int) (dp [][]int) {
	n := len(s)
	if maxLen == -1 {
		maxLen = n
	}
	sa, _, lcp := suffixArray(s)
	next := 0
	dp = make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, maxLen+1)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}

	for i := range sa {
		l := int(sa[i])
		for k := 1; k < 1+min(n-l, maxLen); k++ {
			r := l + k
			if i > 0 && lcp[i-1] >= k {
				dp[l][r-l] = dp[sa[i-1]][k]
			} else {
				dp[l][r-l] = next
				next++
			}
		}
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s string) ([]int32, []int, []int) {
	n := len(s)

	sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New([]byte(s))).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

	rank := make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	height := make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && s[i+h] == s[j+h]; h++ {
			}
		}
		height[rk] = h
	}

	return sa, rank, height
}
