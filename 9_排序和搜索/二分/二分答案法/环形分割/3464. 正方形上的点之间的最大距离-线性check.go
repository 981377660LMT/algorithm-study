// 3464. 正方形上的点之间的最大距离
// https://leetcode.cn/problems/maximize-the-distance-between-points-on-a-square/description/

package main

import (
	"slices"
)

// side = 2, points = [[0,0],[1,2],[2,0],[2,2],[2,1]], k = 4
func main() {
	{

		side := 2
		points := [][]int{{0, 0}, {1, 2}, {2, 0}, {2, 2}, {2, 1}}
		k := 4
		println(maxDistance(side, points, k)) // 1
	}

	{
		side := 9
		points := [][]int{{8, 0}, {5, 9}, {2, 0}, {4, 9}, {0, 1}}
		k := 4
		println(maxDistance(side, points, k)) // 3
	}
}

// !在环形数组中选出 k个点，最大化相邻点的最小距离.
// !环形数组的长度为cycleLen, 各个点的位置为positions.
// !0<=positions[i]<cycleLen.
func MaximizeMinDistOnCycle(cycleLen int, positions []int, k int) int {
	n := len(positions)
	if k > n {
		panic("k must be not greater than n")
	}

	positions = append(positions, cycleLen+positions[0])

	// [start, end) 的代价.
	cost := func(start, end int) int {
		return positions[end] - positions[start]
	}

	check := func(mid int) bool {
		{
			// 先求解链上的问题(剪枝)
			count := 0
			left := 0
			for right := 0; right < n; right++ {
				if cost(left, right+1) >= mid {
					count++
					left = right + 1
				}
			}
			if count >= k {
				return true
			}
			if count <= k-2 {
				return false
			}
		}

		next := make([]int, n)
		right := 0
		for left := 0; left < n; left++ {
			for right < n && cost(left, right) < mid {
				right++
			}
			if cost(left, right) >= mid {
				next[left] = right
			} else {
				next[left] = -1
			}
		}

		type dpItem struct{ count, next int }
		dp := make([]dpItem, n+1)
		dp[n] = dpItem{next: n}
		for i := n - 1; i >= 0; i-- {
			if next[i] == -1 {
				dp[i] = dpItem{next: i}
			} else {
				dp[i] = dp[next[i]]
				dp[i].count++
			}
		}

		for i := 0; i < n; i++ {
			count := dp[i].count
			if count <= k-2 {
				break
			}
			end := dp[i].next
			if cost(0, i)+cost(end, n) >= mid {
				count++
			}
			if count >= k {
				return true
			}
		}

		return false
	}

	left, right := 1, cycleLen>>1
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return right
}

func maxDistance(side int, points [][]int, k int) int {
	trans := func(x, y int) int {
		if x == 0 {
			return y
		}
		if y == side {
			return side + x
		}
		if x == side {
			return 3*side - y
		}
		return 4*side - x
	}

	nums := make([]int, len(points))
	for i, p := range points {
		nums[i] = trans(p[0], p[1])
	}
	slices.Sort(nums)

	return MaximizeMinDistOnCycle(4*side, nums, k)
}
