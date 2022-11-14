package main

import "sort"

type interval struct{ start, end int }

func maxPalindromes(s string, k int) int {
	n := len(s)

	// 写成指针*interval 会慢一些，但是空间用的少一些
	// 428 ms	106.8 MB 不用指针
	// 580 ms	78.7 MB 用指针
	intervals := make([]interval, 0, n)

	// 中心扩展法求回文子串
	expand := func(left, right int) {
		for left >= 0 && right < n && s[left] == s[right] {
			if right-left+1 >= k {
				intervals = append(intervals, interval{left, right})
			}
			left--
			right++
		}
	}

	for i := 0; i < n; i++ {
		expand(i, i)
		expand(i, i+1)
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].end < intervals[j].end
	})

	res := 0
	preEnd := -1
	for _, interval := range intervals {
		if interval.start > preEnd {
			res++
			preEnd = interval.end
		}
	}
	return res
}
