// https://leetcode.cn/problems/palindrome-rearrangement-queries/
// 100129. 回文串重新排列查询
// 给你一个长度为 偶数 n ，下标从 0 开始的字符串 s 。
// 同时给你一个下标从 0 开始的二维整数数组 queries ，其中 queries[i] = [ai, bi, ci, di] 。
// 对于每个查询 i ，你需要执行以下操作：
// 将下标在范围 0 <= ai <= bi < n / 2 内的 子字符串 s[ai:bi] 中的字符重新排列。
// 将下标在范围 n / 2 <= ci <= di < n 内的 子字符串 s[ci:di] 中的字符重新排列。
// 对于每个查询，你的任务是判断执行操作后能否让 s 变成一个 回文串 。
// 每个查询与其他查询都是 独立的 。
// !请你返回一个下标从 0 开始的数组 answer ，如果第 i 个查询执行操作后，可以将 s 变为一个回文串，那么 answer[i] = true，否则为 false 。
// 子字符串 指的是一个字符串中一段连续的字符序列。
// !s[x:y] 表示 s 中从下标 x 到 y 且两个端点 都包含 的子字符串。
//
// !将后半部分区间反转，变为两个区间问题.

package main

import "sort"

func canMakePalindromeQueries(s string, queries [][]int) []bool {
	n := len(s) / 2
	arr1 := make([]byte, n)
	for i := 0; i < n; i++ {
		arr1[i] = s[i] - 'a'
	}
	arr2 := make([]byte, n)
	for i := 0; i < n; i++ {
		arr2[i] = s[2*n-i-1] - 'a'
	}
	diffPreSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		diffPreSum[i] = diffPreSum[i-1]
		if arr1[i-1] != arr2[i-1] {
			diffPreSum[i]++
		}
	}

	C1 := AlphaPresum(arr1, 26, 0)
	C2 := AlphaPresum(arr2, 26, 0)

	res := make([]bool, len(queries))
	for qi, query := range queries {
		l1, r1, l2, r2 := query[0], query[1], query[2], query[3]
		start1, end1 := l1, r1+1
		start2, end2 := 2*n-r2-1, 2*n-l2

		hash1, hash2 := [26]int{}, [26]int{}
		for i := 0; i < 26; i++ {
			hash1[i] = C1(start1, end1, i)
			hash2[i] = C2(start2, end2, i)
		}

		ok := true
		intervals1 := []Interval{{start: start1, end: end1}}
		intervals2 := []Interval{{start: start2, end: end2}}
		enumerate := EnumerateInterval(intervals1, intervals2)
		enumerate(0, n, func(kind int, start, end int, value1, value2 Value) bool {
			if kind == 0 {
				diff := diffPreSum[end] - diffPreSum[start]
				if diff > 0 {
					ok = false
					return true
				}
				return false
			} else if kind == 1 {
				for i := 0; i < 26; i++ {
					count := C2(start, end, i)
					if count > hash1[i] {
						ok = false
						return true
					} else {
						hash1[i] -= count
					}
				}
				return false
			} else if kind == 2 {
				for i := 0; i < 26; i++ {
					count := C1(start, end, i)
					if count > hash2[i] {
						ok = false
						return true
					} else {
						hash2[i] -= count
					}
				}
				return false
			} else {
				return false
			}
		})

		if !ok {
			res[qi] = false
			continue
		}

		for i := 0; i < 26; i++ {
			if hash1[i] != hash2[i] {
				ok = false
				break
			}
		}
		res[qi] = ok
	}

	return res
}

type Str = []byte

// 给定字符集信息和字符s，返回一个查询函数.该函数可以查询s[start:end]间ord的个数.
// 当字符种类很少时，可以用一个counter数组实现区间哈希值的快速计算.
func AlphaPresum(s Str, sigma int, offset int) func(start, end int, ord int) int {
	preSum := make([][]int, len(s)+1)
	for i := range preSum {
		preSum[i] = make([]int, sigma)
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
		return preSum[end][ord-offset] - preSum[start][ord-offset]
	}
}

type Value = int

const INF Value = 1e18
const NoneValue Value = -INF

type Interval = struct {
	start, end int
	value      Value
}

// 给定两个区间列表，每个区间列表都是成对 `不相交` 的。
// 返回一个函数用于遍历`[allStart, allEnd)`范围内的所有区间。
// start和end表示当前区间的起点和终点，value1和value2表示当前区间在两个区间列表中的值，type表示当前区间的类型。
// - 0: 不在两个区间列表中.
// - 1: 在第一个区间列表中,不在第二个区间列表中.
// - 2: 不在第一个区间列表中,在第二个区间列表中.
// - 3: 在两个区间列表中.
func EnumerateInterval(intervals1 []Interval, intervals2 []Interval) func(
	allStart, allEnd int,
	f func(kind int, start, end int, value1, value2 Value) bool,
) {
	intervals1 = append(intervals1[:0:0], intervals1...)
	intervals2 = append(intervals2[:0:0], intervals2...)
	sort.Slice(intervals1, func(i, j int) bool { return intervals1[i].start < intervals1[j].start })
	sort.Slice(intervals2, func(i, j int) bool { return intervals2[i].start < intervals2[j].start })

	return func(allStart, allEnd int, f func(kind int, start, end int, value1, value2 Value) bool) {
		ptr1, ptr2 := 0, 0
		curStart := allStart
		for ptr1 < len(intervals1) && intervals1[ptr1].end <= curStart {
			ptr1++
		}
		for ptr2 < len(intervals2) && intervals2[ptr2].end <= curStart {
			ptr2++
		}

		for curStart < allEnd {
			var start1, end1, start2, end2 int
			if ptr1 < len(intervals1) {
				start1 = min(intervals1[ptr1].start, allEnd)
				end1 = min(intervals1[ptr1].end, allEnd)
			} else {
				start1, end1 = allEnd, allEnd
			}
			if ptr2 < len(intervals2) {
				start2 = min(intervals2[ptr2].start, allEnd)
				end2 = min(intervals2[ptr2].end, allEnd)
			} else {
				start2, end2 = allEnd, allEnd
			}

			// x = curStart 与两个区间相交的清况
			intersect1 := start1 <= curStart && curStart < end1
			intersect2 := start2 <= curStart && curStart < end2

			if intersect1 && intersect2 {
				minEnd := min(end1, end2)
				if f(3, curStart, minEnd, intervals1[ptr1].value, intervals2[ptr2].value) {
					return
				}
				curStart = minEnd
				if end1 == minEnd {
					ptr1++
				}
				if end2 == minEnd {
					ptr2++
				}
			} else if intersect1 {
				curEnd := min(end1, start2)
				if f(1, curStart, curEnd, intervals1[ptr1].value, NoneValue) {
					return
				}
				curStart = curEnd
				if end1 == curEnd {
					ptr1++
				}
			} else if intersect2 {
				curEnd := min(end2, start1)
				if f(2, curStart, curEnd, NoneValue, intervals2[ptr2].value) {
					return
				}
				curStart = curEnd
				if end2 == curEnd {
					ptr2++
				}
			} else {
				minStart := min(start1, start2)
				if f(0, curStart, minStart, NoneValue, NoneValue) {
					return
				}
				curStart = minStart
			}
		}
	}

}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
