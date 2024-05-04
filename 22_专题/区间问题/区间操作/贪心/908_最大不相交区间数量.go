package main

import "sort"

// 452. 用最少数量的箭引爆气球
// https://leetcode.cn/problems/minimum-number-of-arrows-to-burst-balloons/submissions/528930338/
func findMinArrowShots(points [][]int) int {
	res := MaxNonIntersectingIntervals(
		int32(len(points)), func(i int32) (start, end int) { return points[i][0], points[i][1] },
		false, true,
	)
	return len(res)
}

// 最大不相交区间.
// 给定 n 个区间 [left_i,right_i]，尽可能多地选择区间, 满足选择的区间之间互不相交.
// 返回区间索引列表.如果有多种方案，返回字典序最小的那个.
//  allowOverlapping: 是否允许选择的区间端点重合(相等).
//  endInclusive: 是否包含区间右端点.
func MaxNonIntersectingIntervals(
	n int32, f func(i int32) (start, end int),
	allowOverlapping bool, endInclusive bool,
) []int32 {
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []int32{0}
	}

	starts, ends := make([]int, n), make([]int, n)
	for i := int32(0); i < n; i++ {
		starts[i], ends[i] = f(i)
	}
	order := make([]int32, n)
	for i := int32(0); i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return ends[order[i]] < ends[order[j]] })

	res := []int32{order[0]}
	preEnd := ends[order[0]]
	for _, i := range order[1:] {
		start, end := starts[i], ends[i]
		if !endInclusive {
			end--
		}
		if allowOverlapping {
			if start >= preEnd {
				res = append(res, i)
				preEnd = end
			}
		} else {
			if start > preEnd {
				res = append(res, i)
				preEnd = end
			}
		}
	}
	return res
}
