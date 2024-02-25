// !两个字符串的通用思路：遍历后缀树，每个节点处计算belong个数.
// !多个字符串的通用思路：拼接，二分答案，对 height 分组，判定组内元素是否满足条件

package main

import (
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func main() {
	SP220()
}

// TODO http://poj.org/problem?id=3415
// https://www.hankcs.com/program/algorithm/poj-3415-common-substrings.html
func 长度不小于k的公共子串的个数() {}

// !后缀树已解决
// 唯一性可以用 height[i] 与前后相邻值的大小来判定
// https://www.luogu.com.cn/problem/CF427D
func 最短公共唯一子串() {}

// TODO
// https://www.hankcs.com/program/algorithm/aoj-2292-common-palindromes.html#h2-0
// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2292
func 公共回文子串() {}

// PHRASES - Relevant Phrases of Annihilation
// https://www.luogu.com.cn/problem/SP220
// 在每个字符串中至少出现两次且不重叠的最长子串
// 给出 n 个字符串，求一个最长字符串，满足其在每一个字符串都互不重叠地出现至少两次，输出其长度。
// 拼接，二分答案，对 height 分组，判定组内元素在每个字符串中至少出现两次且 sa 的最大最小之差不小于二分值（用于判定是否重叠）
func SP220() {}

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
