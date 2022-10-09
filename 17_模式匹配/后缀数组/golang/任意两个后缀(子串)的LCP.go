// 任意两后缀的 LCP
// !lcp(sa[i], sa[j])=min{height[i＋1..j}
// 求两子串最长公共前缀就转化为了 RMQ 问题。
// 注：若允许离线可以用 Trie+Tarjan 做到线性

// !如果是查询子串[l1:r1],[l2:r2]的LCP 则为min(r1-l1,r2-l2,lcp(l1,l2))

package main

import (
	"index/suffixarray"
	"math"
	"math/bits"
	"reflect"
	"unsafe"
)

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
// 求两个后缀最长公共前缀 O(nlogn)预处理 O(1)查询
func queryLCP(rank, height []int) func(int, int) int {
	n := len(rank)

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// func main() {
// 	_, rank, height := suffixArray([]byte("banana"))
// 	lcp := queryLCP(rank, height)
// 	println(lcp(2, 4))
// }
