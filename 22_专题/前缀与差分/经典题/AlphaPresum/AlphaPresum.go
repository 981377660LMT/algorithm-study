package main

// https://leetcode.cn/problems/number-of-same-end-substrings/description/
// 2955. Number of Same-End Substrings
func sameEndSubstringCount(s string, queries [][]int) []int {
	S := AlphaPresum(s, 26, 97)
	res := make([]int, len(queries))
	for i, q := range queries {
		start, end := q[0], q[1]+1
		cur := 0
		for i := 0; i < 26; i++ {
			freq := S(start, end, i+97)
			cur += freq * (freq + 1) / 2
		}
		res[i] = cur
	}
	return res
}

type Str = string // []byte

// 给定字符集信息和字符s，返回一个查询函数.该函数可以查询s[start:end]间ord的个数.
// 当字符种类很少时，可以用一个counter数组实现区间哈希值的快速计算.
func AlphaPresum(s Str, sigma int, offset int) func(start, end int, ord int) int {
	preSum := make([][]int32, len(s)+1)
	for i := range preSum {
		preSum[i] = make([]int32, sigma)
	}
	for i := 1; i <= len(s); i++ {
		copy(preSum[i], preSum[i-1])
		preSum[i][int(s[i-1])-offset]++
	}

	return func(start, end int, ord int) int {
		if start < 0 {
			start = 0
		}
		if end > len(s) {
			end = len(s)
		}
		if start >= end {
			return 0
		}
		return int(preSum[end][ord-offset] - preSum[start][ord-offset])
	}
}
