// z函数
// 也可以求LCP后查询两个后缀的LCP

package main

import (
	"index/suffixarray"
	"math"
	"math/bits"
	"reflect"
	"unsafe"
)

func sumScores(s string) int64 {
	n := len(s)
	lcp := queryLCP(s)
	res := 0
	for i := 0; i < n; i++ {
		res += lcp(0, i)
	}
	return int64(res)

}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
// 求两个后缀最长公共前缀 O(nlogn)预处理 O(1)查询
func queryLCP(s string) func(int, int) int {
	n := len(s)

	_, rank, height := suffixArray([]byte(s))

	max := int(math.Ceil(math.Log2(float64(n)))) + 1
	st := make([][]int, n)
	for i := range st {
		st[i] = make([]int, max)
	}

	for i, v := range height {
		st[i][0] = v
	}
	for j := 1; 1<<j <= n; j++ {
		for i := 0; i+1<<j <= n; i++ {
			st[i][j] = min(st[i][j-1], st[i+1<<(j-1)][j-1])
		}
	}

	_q := func(l, r int) int { k := bits.Len(uint(r-l)) - 1; return min(st[l][k], st[r-1<<k][k]) }
	lcp := func(i, j int) int {
		if i == j {
			return n - i
		}
		// 将 s[i:] 和 s[j:] 通过 rank 数组映射为 height 的下标
		ri, rj := rank[i], rank[j]
		if ri > rj {
			ri, rj = rj, ri
		}
		return _q(ri+1, rj+1)
	}

	return lcp

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
