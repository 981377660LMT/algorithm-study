package main

import (
	"fmt"
	"strings"
)

// 计算字符串 s 的最小周期。如果没有找到，返回 s 的长度。
func minimalPeriod(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	searchZone := (s + s)[1 : 2*n-1]
	idx := strings.Index(searchZone, s)
	if idx != -1 {
		return idx + 1 // 因为 searchZone 从原串 dbl 的下标 1 开始
	}
	return n
}

func main() {
	tests := []string{"abab", "aaa", "aba", "a", "abcab", ""}
	for _, s := range tests {
		fmt.Printf("%q -> %d\n", s, minimalPeriod(s))
	}
}
