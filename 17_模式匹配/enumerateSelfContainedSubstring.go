package main

import "sort"

// 3104. Find Longest Self-Contained Substring (查找最长的自包含子字符串)
// https://leetcode.cn/problems/find-longest-self-contained-substring/description/
// 给定字符串s，找到最长的自包含子字符串，返回其长度，没有则返回-1。
func maxSubstringLength(s string) int {
	n := int32(len(s))
	intervals := EnumerateSelfContainedSubstring(s)
	res := -1
	for _, interval := range intervals {
		start, end := interval[0], interval[1]
		if m := end - start; m != n {
			res = max(res, int(m))
		}
	}
	return res
}

// 返回所有自包含子串.
// 自包含子串是指，该子串中所有的字符，均未在子串以外的部分出现.
// 时间复杂度: O(n + ∑^2).
func EnumerateSelfContainedSubstring(s string) [][2]int32 {
	first, last, counter := map[rune]int32{}, map[rune]int32{}, map[rune]int32{}
	for i, c := range s {
		if _, ok := first[c]; !ok {
			first[c] = int32(i)
		}
		last[c] = int32(i)
		counter[c]++
	}
	chars := make([]rune, 0, len(first))
	for c := range first {
		chars = append(chars, c)
	}
	sort.Slice(chars, func(i, j int) bool { return first[chars[i]] < first[chars[j]] })

	m := int32(len(chars))
	var res [][2]int32
	for i, c1 := range chars {
		left, right, count := first[c1], int32(0), int32(0)
		for j := i; j < int(m); j++ {
			c2 := chars[j]
			right = max32(right, last[c2])
			count += counter[c2]
			if count == right-left+1 {
				res = append(res, [2]int32{left, right + 1})
			}
		}
	}
	return res
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
