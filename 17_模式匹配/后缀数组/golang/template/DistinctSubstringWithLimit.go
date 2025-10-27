// "带限制的本质不同子串个数"问题通用解法.
// !本质不同：SA+LCP 去重

package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"unsafe"
)

func main() {
	// numberofSubstrings()
	cf271d()
}

// 2261. 含最多 K 个可整除元素的子数组(带限制的本质不同子串个数)
// https://leetcode.cn/problems/k-divisible-elements-subarrays/
// 找出并返回满足要求的不同的子数组数，要求子数组中最多 k 个可被 p 整除的元素。
func countDistinct(nums []int, k int, p int) (res int) {
	n := int32(len(nums))
	preSum := make([]int32, n+1)
	for i, v := range nums {
		preSum[i+1] = preSum[i]
		if v%p == 0 {
			preSum[i+1]++
		}
	}
	check := func(start, end int32) bool {
		return int(preSum[end]-preSum[start]) <= k
	}
	return DistinctSubstringWithLimit(n, func(i int32) int32 { return int32(nums[i]) }, check)
}

// Good Substrings
// https://codeforces.com/problemset/problem/271/D
// 若一个字符串中包含不超过k个不好的字母，则这个字符串是“好的”。
// 求出s的所有子串中，本质不同的好的子串的数量。
func cf271d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var t string
	fmt.Fscan(in, &t)
	var k int32
	fmt.Fscan(in, &k)

	n := int32(len(s))
	isBad := [26]bool{}
	for i, b := range t {
		isBad[i] = b == '0'
	}
	preSum := make([]int32, n+1)
	for i, c := range s {
		preSum[i+1] = preSum[i]
		if isBad[c-'a'] {
			preSum[i+1]++
		}
	}

	check := func(start, end int32) bool {
		return preSum[end]-preSum[start] <= k
	}
	res := DistinctSubstringWithLimit(n, func(i int32) int32 { return int32(s[i]) }, check)
	fmt.Fprintln(out, res)
}

// 本质不同子串个数.
// https://judge.yosupo.jp/problem/number_of_substrings
func numberofSubstrings() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	check := func(start, end int32) bool { return true }
	res := DistinctSubstringWithLimit(int32(len(s)), func(i int32) int32 { return int32(s[i]) }, check)
	fmt.Fprintln(out, res)
}

// 带限制的本质不同子串个数.
//
//	predicate: 子串是否满足条件的判定函数,要求具有二段性.
func DistinctSubstringWithLimit(
	n int32, f func(i int32) int32,
	predicate func(start, end int32) bool,
) int {
	// 1.双指针枚举后缀，求出所有满足条件的子串的数量
	validLength := make([]int32, n) // 每个后缀的合法长度
	res, right := 0, int32(0)
	for left := int32(0); left < n; left++ {
		for right < n && predicate(left, right+1) {
			right++
		}
		res += int(right - left)
		validLength[left] = right - left
	}

	// 2.计算子串重复数量:按后缀排序的顺序枚举后缀 lcp(height)去重
	sa, _, height := SuffixArray32(n, f)
	for i := int32(0); i < n-1; i++ {
		suffix1, suffix2 := sa[i], sa[i+1]
		len1, len2 := validLength[suffix1], validLength[suffix2]
		res -= int(min32(height[i+1], min32(len1, len2)))
	}
	return res
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
