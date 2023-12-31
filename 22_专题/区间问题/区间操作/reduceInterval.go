// 删除重叠区间

package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(ReduceIntervals([][]int{{1, 2}, {2, 3}, {1, 3}}, true))  // [2]
	fmt.Println(ReduceIntervals([][]int{{1, 2}, {2, 3}, {1, 3}}, false)) // [0, 1]

}

// https://leetcode.cn/problems/remove-covered-intervals/
// 1288. 删除被覆盖区间
func removeCoveredIntervals(intervals [][]int) int {
	return len(ReduceIntervals(intervals, true))
}

// 删除重叠覆盖区间.
//
//	intervals 左闭右开区间.
//	removeIncluded 是删除包含的区间还是被包含的区间.默认为删除被包含的区间.
//	返回：按照区间的起点排序的剩余的区间索引(相同的区间会保留).
func ReduceIntervals(intervals [][]int, removeIncluded bool) []int {
	n := len(intervals)
	res := []int{}
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}

	if removeIncluded {
		sort.Slice(order, func(i, j int) bool {
			a, b := order[i], order[j]
			if intervals[a][0] != intervals[b][0] {
				return intervals[a][0] < intervals[b][0]
			}
			return intervals[a][1] > intervals[b][1]
		})
		for i := 0; i < n; i++ {
			cur := order[i]
			if len(res) > 0 {
				pre := res[len(res)-1]
				curStart, curEnd := intervals[cur][0], intervals[cur][1]
				preStart, preEnd := intervals[pre][0], intervals[pre][1]
				if curEnd <= preEnd && curEnd-curStart < preEnd-preStart {
					continue
				}
			}
			res = append(res, cur)
		}
	} else {
		sort.Slice(order, func(i, j int) bool {
			a, b := order[i], order[j]
			if intervals[a][1] != intervals[b][1] {
				return intervals[a][1] < intervals[b][1]
			}
			return intervals[a][0] > intervals[b][0]
		})
		for i := 0; i < n; i++ {
			cur := order[i]
			if len(res) > 0 {
				pre := res[len(res)-1]
				curStart, curEnd := intervals[cur][0], intervals[cur][1]
				preStart, preEnd := intervals[pre][0], intervals[pre][1]
				if curStart <= preStart && curEnd-curStart > preEnd-preStart {
					continue
				}
			}
			res = append(res, cur)
		}
	}
	return res
}
