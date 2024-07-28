package main

// 100348. 统计 1 显著的字符串的数量-根号值域
// https://leetcode.cn/problems/count-the-number-of-substrings-with-dominant-ones/solutions/2860181/mei-ju-by-tsreaper-x830/
// 给你一个二进制字符串 s。请你统计并返回其中 1 显著 的子字符串的数量。
// !如果字符串中 1 的数量 大于或等于 0 的数量的 平方，则认为该字符串是一个 1 显著 的字符串
// !n<=4e4
func numberOfSubstrings(s string) int {
	n := int32(len(s))
	nexts := make([]int32, n+1)
	nexts[n] = n
	for i := n - 1; i >= 0; i-- {
		if s[i] == '0' {
			nexts[i] = i
		} else {
			nexts[i] = nexts[i+1]
		}
	}

	upper := int32(0)
	for upper*upper <= n {
		upper++
	}
	res := 0
	for i := int32(0); i < n; i++ {
		j, zeros := i, int32(0)
		if s[i] == '0' {
			zeros = 1
		}
		for j != n && zeros < upper {
			ones := (nexts[j+1] - i) - zeros
			okCount := max(0, ones-zeros*zeros+1)
			len_ := nexts[j+1] - j
			res += int(min(okCount, len_))
			j = nexts[j+1]
			zeros++
		}
	}
	return res
}

// 按照target的数量分组枚举子串
// snippet:
func EnumerateSubstringGroupByTargetCount[T comparable](arr []T, target T) {
	n := int32(len(arr))
	nexts := make([]int32, n+1)
	nexts[n] = n
	for i := n - 1; i >= 0; i-- {
		if arr[i] == target {
			nexts[i] = i
		} else {
			nexts[i] = nexts[i+1]
		}
	}
	for i := int32(0); i < n; i++ {
		j, count := i, 1
		if arr[i] == target {
			count = 1
		} else {
			count = 0
		}
		for j != n {
			// do something using (i, j, count)
			//
			j = nexts[j+1]
			count++
		}
	}
}
