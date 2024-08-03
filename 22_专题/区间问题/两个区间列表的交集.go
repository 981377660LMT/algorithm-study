package main

import "sort"

// 有序区间列表交集长度.
func IntervalsIntersectionLen[T int | int32](intervals1, intervals2 [][2]T) T {
	n1, n2 := len(intervals1), len(intervals2)
	res := T(0)
	left, right := 0, 0
	for left < n1 && right < n2 {
		s1, e1, s2, e2 := intervals1[left][0], intervals1[left][1], intervals2[right][0], intervals2[right][1]
		if (s1 <= e2 && e2 <= e1) || (s2 <= e1 && e1 <= e2) {
			res += min(e1, e2) - max(s1, s2)
		}
		if e1 < e2 {
			left++
		} else {
			right++
		}
	}
	return res
}

func EnumerateIntervalsIntersection[T int | int32](intervals1, intervals2 [][2]T, f func(T, T)) {
	n1, n2 := len(intervals1), len(intervals2)
	left, right := 0, 0
	for left < n1 && right < n2 {
		s1, e1, s2, e2 := intervals1[left][0], intervals1[left][1], intervals2[right][0], intervals2[right][1]
		if (s1 <= e2 && e2 <= e1) || (s2 <= e1 && e1 <= e2) {
			f(max(s1, s2), min(e1, e2))
		}
		if e1 < e2 {
			left++
		} else {
			right++
		}
	}
}

type Interval[V any] struct {
	start, end int
	value      V
}

// 给定两个区间列表，每个区间列表都是成对 `不相交` 的。
// 返回一个函数用于遍历`[allStart, allEnd)`范围内的所有区间。
// start和end表示当前区间的起点和终点，value1和value2表示当前区间在两个区间列表中的值，type表示当前区间的类型。
// - 0: 不在两个区间列表中.
// - 1: 在第一个区间列表中,不在第二个区间列表中.
// - 2: 不在第一个区间列表中,在第二个区间列表中.
// - 3: 在两个区间列表中.
func EnumerateInterval[V any](intervals1 []Interval[V], intervals2 []Interval[V]) func(
	allStart, allEnd int,
	emptyValue V,
	f func(kind int8, start, end int, value1, value2 V) bool,
) {
	intervals1 = append(intervals1[:0:0], intervals1...)
	intervals2 = append(intervals2[:0:0], intervals2...)
	sort.Slice(intervals1, func(i, j int) bool { return intervals1[i].start < intervals1[j].start })
	sort.Slice(intervals2, func(i, j int) bool { return intervals2[i].start < intervals2[j].start })

	return func(allStart, allEnd int, emptyValue V, f func(kind int8, start, end int, value1, value2 V) bool) {
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
				if f(1, curStart, curEnd, intervals1[ptr1].value, emptyValue) {
					return
				}
				curStart = curEnd
				if end1 == curEnd {
					ptr1++
				}
			} else if intersect2 {
				curEnd := min(end2, start1)
				if f(2, curStart, curEnd, emptyValue, intervals2[ptr2].value) {
					return
				}
				curStart = curEnd
				if end2 == curEnd {
					ptr2++
				}
			} else {
				minStart := min(start1, start2)
				if f(0, curStart, minStart, emptyValue, emptyValue) {
					return
				}
				curStart = minStart
			}
		}
	}

}
