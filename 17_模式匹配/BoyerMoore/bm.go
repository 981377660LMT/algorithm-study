// https://github.com/spaghetti-source/algorithm/blob/master/string/knuth_morris_pratt.cc
//
// # Boyer-Moore string matching
//
// Description:
//
//	It processes a pattern string to find
//	all occurrence of a given text.
//
// Algorithm:
//
//	It matches a pattern string from the last to the front.
//	If the string is random, this can skip many comparison.
//
// Complexity:
//
//	O(n + |occur|) with O(m) preprocessing.
//	In a random case, it reduced to O(n/m + |occur|).
//
// Verified:
//
//	SPOJ 21524
//
// Comment:
//
//	BM is often considered faster than KMP.
//	However, in the programming contest setting,
//	these are equally fast.

package main

import "fmt"

func main() {
	s1 := []int{200, 97, 97}
	s2 := []int{97, 97}
	bm := NewBm(int32(len(s2)), func(i int32) byte { return byte(s2[i]) + 1 })
	fmt.Println(bm.Match(int32(len(s1)), func(i int32) byte {
		return byte(s1[i]) + 1
	}))
}

// https://leetcode.cn/problems/number-of-subarrays-that-match-a-pattern-ii/description/
func countMatchingSubarrays(nums []int, pattern []int) int {
	arr := make([]int, len(nums)-1)
	for i := range arr {
		if nums[i] > nums[i+1] {
			arr[i] = -1
		} else if nums[i] < nums[i+1] {
			arr[i] = 1
		} else {
			arr[i] = 0
		}
	}
	bm := NewBm(int32(len(pattern)), func(i int32) byte { return byte(pattern[i]) + 1 })
	return len(bm.Match(int32(len(arr)), func(i int32) byte { return byte(arr[i]) + 1 }))
}

// TODO: FIXME
// Boyer-Moore string matching.
type Bm struct {
	n          int32
	pattern    func(int32) byte
	skip, next []int32
}

func NewBm(n int32, pattern func(int32) byte) *Bm {
	skip := make([]int32, 256)
	for i := int32(0); i < n; i++ {
		skip[pattern(i)] = n - i - 1
	}
	g := make([]int32, n)
	next := make([]int32, n)
	for i := int32(0); i < n; i++ {
		g[i] = n
		next[i] = 2*n - i - 1
	}
	for i, j := n-1, n; i >= 0; {
		g[i] = j
		for j < n && pattern(j) != pattern(i) {
			next[j] = min32(next[j], n-i-1)
			j = g[j]
		}
		i--
		j--
	}
	return &Bm{n: n, pattern: pattern, skip: skip, next: next}
}

func (bm *Bm) Match(m int32, text func(int32) byte) (res []int32) {
	for i := bm.n - 1; i < m; {
		j := bm.n - 1
		for j >= 0 && text(i) == bm.pattern(j) {
			i--
			j--
		}
		if j < 0 {
			res = append(res, i+1)
			i += bm.n + 1
		} else {
			i += max32(bm.skip[text(i)], bm.next[j])
		}
	}
	return
}

type Kmp struct {
	n       int32
	pattern func(int32) int32
	fail    []int32
}

func NewKmp(n int32, pattern func(int32) int32) *Kmp {
	res := &Kmp{n: n, pattern: pattern}
	res.fail = make([]int32, n+1)
	for i := range res.fail {
		res.fail[i] = -1
	}
	for i, j := int32(1), int32(-1); i <= n; i++ {
		for j >= 0 && pattern(j) != pattern(i-1) {
			j = res.fail[j]
		}
		j++
		res.fail[i] = j
	}
	return res
}

func (kmp *Kmp) Match(m int32, text func(int32) int32) (res []int32) {
	for i, k := int32(0), int32(0); i < m; i++ {
		for k >= 0 && (k >= kmp.n || text(i) != kmp.pattern(k)) {
			k = kmp.fail[k]
		}
		if k++; k == kmp.n {
			res = append(res, i-kmp.n+1)
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

func max(a, b int) int {
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
