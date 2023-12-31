package main

import "sort"

const INF int = 1e18

// 986. 区间列表的交集
// https://leetcode.cn/problems/interval-list-intersections/description/
func intervalIntersection(firstList [][]int, secondList [][]int) [][]int {
	intervals1 := make([]Interval, len(firstList))
	for i, interval := range firstList {
		intervals1[i] = Interval{start: interval[0], end: interval[1] + 1, value: 1}
	}
	intervals2 := make([]Interval, len(secondList))
	for i, interval := range secondList {
		intervals2[i] = Interval{start: interval[0], end: interval[1] + 1, value: 2}
	}

	res := [][]int{}
	enumerate := EnumerateInterval(intervals1, intervals2)
	enumerate(-INF, INF, func(kind int, start, end int, value1, value2 Value) {
		if kind == 3 {
			res = append(res, []int{start, end - 1})
		}
	})
	return res
}

type Value = int

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
	f func(kind int, start, end int, value1, value2 Value),
) {
	intervals1 = append(intervals1[:0:0], intervals1...)
	intervals2 = append(intervals2[:0:0], intervals2...)
	sort.Slice(intervals1, func(i, j int) bool { return intervals1[i].start < intervals1[j].start })
	sort.Slice(intervals2, func(i, j int) bool { return intervals2[i].start < intervals2[j].start })

	return func(allStart, allEnd int, f func(kind int, start, end int, value1, value2 Value)) {
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
				f(3, curStart, minEnd, intervals1[ptr1].value, intervals2[ptr2].value)
				curStart = minEnd
				if end1 == minEnd {
					ptr1++
				}
				if end2 == minEnd {
					ptr2++
				}
			} else if intersect1 {
				curEnd := min(end1, start2)
				f(1, curStart, curEnd, intervals1[ptr1].value, NoneValue)
				curStart = curEnd
				if end1 == curEnd {
					ptr1++
				}
			} else if intersect2 {
				curEnd := min(end2, start1)
				f(2, curStart, curEnd, NoneValue, intervals2[ptr2].value)
				curStart = curEnd
				if end2 == curEnd {
					ptr2++
				}
			} else {
				minStart := min(start1, start2)
				f(0, curStart, minStart, NoneValue, NoneValue)
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
