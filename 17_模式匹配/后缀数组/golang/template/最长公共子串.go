// 有一个序列A。X,Y是给定的A的两个子串。
// 每次操作可以在X的开头或末尾增添或删除一个数字，且需满足任意时刻X非空且为A的子串，
// 求把X变成Y的最少操作次数。
// 题目保证解存在
// n<=2e5，1<=nums[i]<=n

// 1. 如果最长公共子串>=1，那么可以将X一直删到LCS，然后将X一直加到Y，操作len(X)+len(Y)-2*LCS次
// 2. 如果最长公共子串=0，那么就要求从X中某个顶点出发跑到Y中某个顶点的最短路

package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"unsafe"
)

// https://judge.yosupo.jp/problem/longest_common_substring
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s1, s2 string
	fmt.Fscan(in, &s1, &s2)
	start1, end1, start2, end2 := longestCommonSubstring1(s1, s2)
	fmt.Fprintln(out, start1, end1, start2, end2)
}

// 最长公共子串
func longestCommonSubstring1(s1, s2 string) (start1, end1, start2, end2 int32) {
	ords1, ords2 := make([]rune, len(s1)), make([]rune, len(s2))
	for i, x := range s1 {
		ords1[i] = x
	}
	for i, x := range s2 {
		ords2[i] = x
	}
	return longestCommonSubstring2(ords1, ords2)
}

func longestCommonSubstring2(ords1, ords2 []int32) (start1, end1, start2, end2 int32) {
	if len(ords1) == 0 || len(ords2) == 0 {
		return
	}

	dummy := max32(maxs32(ords1), maxs32(ords2)) + 1
	sb := make([]int32, 0, len(ords1)+len(ords2)+1)
	sb = append(sb, ords1...)
	sb = append(sb, dummy)
	sb = append(sb, ords2...)
	sa, _, lcp := SuffixArray32(int32(len(sb)), func(i int32) int32 { return sb[i] })

	len_ := int32(0)
	len1 := int32(len(ords1))
	for i := 1; i < len(sb); i++ {
		if (sa[i-1] < len1) == (sa[i] < len1) {
			continue
		}
		if lcp[i] <= len_ {
			continue
		}
		len_ = lcp[i]

		// 来自s和t的不同子串
		// 找到了(严格)更长的公共子串,更新答案
		i1, i2 := sa[i-1], sa[i]
		if i1 > i2 {
			i1, i2 = i2, i1
		}

		start1 = i1
		end1 = start1 + len_
		start2 = i2 - len1 - 1
		end2 = start2 + len_
	}

	return
}

func SuffixArray32(n int32, f func(i int32) int32) (sa, rank, height []int32) {
	s := make([]byte, 0, n*4)
	for i := int32(0); i < n; i++ {
		v := f(i)
		s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	sa = make([]int32, 0, n)
	for _, v := range _sa {
		if v&3 == 0 {
			sa = append(sa, v>>2)
		}
	}
	rank = make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[sa[i]] = i
	}
	height = make([]int32, n)
	h := int32(0)
	for i := int32(0); i < n; i++ {
		rk := rank[i]
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && f(i+h) == f(j+h); h++ {
			}
		}
		height[rk] = h
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func maxs32(a []int32) int32 {
	res := a[0]
	for _, v := range a[1:] {
		if v > res {
			res = v
		}
	}
	return res
}
