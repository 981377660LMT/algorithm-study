// !不同子串的长度之和

package main

import (
	"index/suffixarray"
	"reflect"
	"unsafe"
)

// 枚举每个后缀，计算前缀总数，再减掉重复
func diffSum(s string) int {
	n := len(s)
	_, _, height := suffixArray([]byte(s))
	res := n * (n + 1) * (n + 2) / 6 // 所有子串长度 1到n的平方和
	for _, h := range height {
		res -= h * (h + 1) / 2
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
