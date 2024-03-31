package main

import (
	"fmt"
	"sort"
)

// 给你一个长度为 n 的整数数组 nums 和一个 正 整数 k 。
// 一个子序列的 能量 定义为子序列中 任意 两个元素的差值绝对值的 最小值 。
// 请你返回 nums 中长度 等于 k 的 所有 子序列的 能量和 。
// 由于答案可能会很大，将答案对 109 + 7 取余 后返回。
//
//
// n<=50, 考虑折半枚举.

const MOD int = 1e9 + 7
const INF int = 1e18

func sumOfPowers(nums []int, k int) int {
	n := len(nums)
	mid := n / 2
	sort.Ints(nums)
	left := nums[:mid]
	right := nums[mid:]

	// 处理出每个子集的最大、最小值、内部最小差值.
	getLeft := func(arr []int) [][][2]int {
		n := len(arr)
		groupBySize := make([][][2]int, n+1)

		var dfs func(index int, path []int, minDiff int)
		dfs = func(index int, path []int, minDiff int) {
			if index == n {
				if len(path) > 0 {
					max_ := path[len(path)-1]
					groupBySize[len(path)] = append(groupBySize[len(path)], [2]int{max_, minDiff})
				}
				return
			}

			nextMinDiff := INF
			if len(path) > 0 {
				nextMinDiff = min(minDiff, arr[index]-path[len(path)-1])
			}
			path = append(path, arr[index])
			dfs(index+1, path, nextMinDiff)
			path = path[:len(path)-1]

			dfs(index+1, path, minDiff)
		}

		dfs(0, []int{}, INF)
		return groupBySize
	}

	getRight := func(arr []int) [][][2]int {
		n := len(arr)
		groupBySize := make([][][2]int, n+1)

		var dfs func(index int, path []int, minDiff int)
		dfs = func(index int, path []int, minDiff int) {
			if index == n {
				if len(path) > 0 {
					min_ := path[0]
					groupBySize[len(path)] = append(groupBySize[len(path)], [2]int{min_, minDiff})
				}
				return
			}

			nextMinDiff := INF
			if len(path) > 0 {
				nextMinDiff = min(minDiff, arr[index]-path[len(path)-1])
			}
			path = append(path, arr[index])
			dfs(index+1, path, nextMinDiff)
			path = path[:len(path)-1]

			dfs(index+1, path, minDiff)
		}

		dfs(0, []int{}, INF)
		return groupBySize
	}

	leftInfo := getLeft(left)
	rightInfo := getRight(right)

	res := 0

	// 只选一边的情况.
	if mid == k {
		leftGroup := leftInfo[mid]
		for _, v := range leftGroup {
			res += v[1]
			if res >= MOD {
				res -= MOD
			}
		}
	}

	if n-mid == k {
		rightGroup := rightInfo[n-mid]
		for _, v := range rightGroup {
			res += v[1]
			if res >= MOD {
				res -= MOD
			}
		}
	}

	// 最小差值在左边/右边/右侧最大值-左侧最小值.
	for leftSize := 0; leftSize <= min(k, mid); leftSize++ {
		rightSize := k - leftSize
		if leftSize == 0 || rightSize == 0 {
			continue
		}
		arr1, arr2 := leftInfo[leftSize], rightInfo[rightSize]
		// fmt.Println(arr1)
		// fmt.Println(arr2)

		// 双指针
		for _, v1 := range arr1 {
			for _, v2 := range arr2 {
				res += min(v1[1], min(v2[1], v2[0]-v1[0]))
				if res >= MOD {
					res -= MOD
				}
			}
		}
	}

	res %= MOD
	if res < 0 {
		res += MOD
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// nums = [1,2,3,4], k = 3
	fmt.Println(sumOfPowers([]int{1, 2, 3, 4, 5, 6, 7, 8}, 4)) // 6
}
