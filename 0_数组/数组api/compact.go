package main

// https://leetcode.cn/problems/remove-duplicates-from-sorted-array-ii/description/
func removeDuplicates(nums []int) int {
	nums = CompactK(nums, 2)
	return len(nums)
}

// !Like Rust's Vec::dedup.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	for k := 1; k < len(s); k++ {
		if s[k] != s[k-1] {
			if i != k {
				s[i] = s[k]
			}
			i++
		}
	}
	clear(s[i:])
	return s[:i]
}

// !Like Rust's Vec::dedup_by.
func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	for k := 1; k < len(s); k++ {
		if !eq(s[k], s[k-1]) {
			if i != k {
				s[i] = s[k]
			}
			i++
		}
	}
	clear(s[i:])
	return s[:i]
}

func CompactK[S ~[]E, E comparable](s S, k int) S {
	if len(s) <= k {
		return s
	}
	slow := k
	for fast := k; fast < len(s); fast++ {
		if s[fast] != s[slow-k] {
			s[slow] = s[fast]
			slow++
		}
	}
	clear(s[slow:])
	return s[:slow]
}
