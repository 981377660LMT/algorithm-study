// 区间最大不相交集合
// q1:
// 最多不相交区间数量.
// 给定n个区间，每个区间形如[start, end).
// 处理q个查询，每个查询形如[start, end)，求查询区间内可以选出的最多不相交区间数量。
// q2:
// 拓扑图最长路径问题.
// 给定一副包含n个顶点的有向拓扑图，以及m条边(l_i,r_i)，每条边的长度为1，且顶点x到顶点x+1有一条长度为0的边。
// 处理q个查询，每个查询形如[l, r]，求查询区间内的最长距离。

package main

import (
	"fmt"
	"math/bits"
	"sort"
)

func main() {
	intervals := [][2]int{{1, 3}, {2, 4}, {3, 5}, {4, 6}, {5, 7}}
	f := MaxNonIntersectingIntervals(intervals)
	fmt.Println(f(1, 5)) // 2
	fmt.Println(f(2, 6)) // 2
	fmt.Println(f(1, 7)) // 3
}

// interval: [start, end).
func MaxNonIntersectingIntervals(intervals [][2]int) func(start, end int) int32 {
	intervals = append(intervals[:0:0], intervals...)
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })
	pos := int32(1)
	for i := 1; i < len(intervals); i++ {
		for pos > 0 && intervals[pos-1][1] >= intervals[i][1] {
			pos--
		}
		if pos == 0 || intervals[pos-1][0] < intervals[i][0] {
			intervals[pos] = intervals[i]
			pos++
		}
	}
	m := pos
	log := int32(floorLog32(uint32(m)))
	jump := make([][]int32, log+1)
	for i := int32(0); i <= log; i++ {
		jump[i] = make([]int32, m)
	}
	for i, r := m-1, m; i >= 0; i-- {
		for r-1 >= i && intervals[r-1][0] >= intervals[i][1] {
			r--
		}
		jump[0][i] = r
	}
	for i := int32(0); i+1 <= log; i++ {
		for j := int32(0); j < m; j++ {
			if tmp := jump[i][j]; tmp == m {
				jump[i+1][j] = m
			} else {
				jump[i+1][j] = jump[i][tmp]
			}
		}
	}

	f := func(start, end int) int32 {
		if start >= end {
			return 0
		}
		if m == 0 {
			return 0
		}
		left, right := int32(0), m-1
		for left < right {
			mid := (left + right) >> 1
			if intervals[mid][0] < start {
				left = mid + 1
			} else {
				right = mid
			}
		}
		if intervals[left][0] < start || intervals[left][1] > end {
			return 0
		}
		cur := left
		res := int32(1)
		for i := len(jump) - 1; i >= 0; i-- {
			next := jump[i][cur]
			if next == m || intervals[next][1] > end {
				continue
			}
			res += 1 << i
			cur = next
		}
		return res
	}
	return f
}

func floorLog32(x uint32) int {
	if x <= 0 {
		panic("IllegalArgumentException")
	}
	return 31 - bits.LeadingZeros32(x)
}
