package main

import (
	"index/suffixarray"
	"reflect"
	"unsafe"
)

// 返回 s 的不同子字符串的个数
// https://oi-wiki.org/string/sa/#_13

// 用所有子串的个数，减去相同子串的个数，就可以得到不同子串的个数。
// !子串就是后缀的前缀 按后缀排序的顺序枚举后缀，每次新增的子串就是除了与上一个后缀的 LCP 剩下的前缀
// !计算后缀数组和高度数组。根据高度数组的定义，所有高度之和就是相同子串的个数。(每一对相同子串在高度数组产生1贡献)
func countDistinct(s string) int {
	n := len(s)
	res := n * (n + 1) / 2
	_, _, height := suffixArray([]byte(s))
	for _, h := range height {
		res -= h
	}
	return res
}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s []byte) ([]int32, []int, []int) {
	n := len(s)

	sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

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
