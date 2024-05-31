package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(LIS([]int{1, 2, 1, 1, 2}, false))
	fmt.Println(LISDp([]int{1, 2, 1, 1, 2}, false))
	fmt.Println(GetLIS([]int{1, 2, 1, 1, 2}, false))
}

// 求LIS长度
func LIS(nums []int, strict bool) int {
	n := len(nums)
	if n <= 1 {
		return n
	}
	lis := []int{} // lis[i] 表示长度为 i 的上升子序列的最小末尾值
	var f func([]int, int) int
	if strict {
		f = sort.SearchInts
	} else {
		f = func(a []int, x int) int {
			return sort.SearchInts(a, x+1)
		}
	}
	for i := 0; i < n; i++ {
		pos := f(lis, nums[i])
		if pos == len(lis) {
			lis = append(lis, nums[i])
		} else {
			lis[pos] = nums[i]
		}
	}
	return len(lis)
}

const INF int = 2e18

// 返回每个位置为结尾的LIS长度(包括自身)
func LISDp(nums []int, strict bool) (int32, []int32) {
	n := int32(len(nums))
	dp := make([]int, n)
	for i := range dp {
		dp[i] = INF
	}
	lis := int32(0)
	lisRank := make([]int32, n)
	var f func([]int, int) int
	if strict {
		f = sort.SearchInts
	} else {
		f = func(a []int, x int) int {
			return sort.SearchInts(a, x+1)
		}
	}
	for i := int32(0); i < n; i++ {
		pos := int32(f(dp, nums[i]))
		dp[pos] = nums[i]
		lisRank[i] = pos + 1
		if lis < pos+1 {
			lis = pos + 1
		}
	}
	return lis, lisRank
}

// 求LI，返回(LIS,LIS的组成下标)
func GetLIS(nums []int, strict bool) ([]int, []int) {
	n := len(nums)
	lis := []int{}            // lis[i] 表示长度为 i 的上升子序列的最小末尾值
	dpIndex := make([]int, n) // 每个元素对应的LIS长度
	var f func([]int, int) int
	if strict {
		f = sort.SearchInts
	} else {
		f = func(a []int, x int) int {
			return sort.SearchInts(a, x+1)
		}
	}
	for i := 0; i < n; i++ {
		pos := f(lis, nums[i])
		if pos == len(lis) {
			lis = append(lis, nums[i])
		} else {
			lis[pos] = nums[i]
		}
		dpIndex[i] = pos
	}

	res, resIndex := []int{}, []int{}
	j := len(lis) - 1
	for i := n - 1; i >= 0; i-- {
		if dpIndex[i] == j {
			res = append(res, nums[i])
			resIndex = append(resIndex, i)
			j -= 1
		}
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
		resIndex[i], resIndex[j] = resIndex[j], resIndex[i]
	}
	return res, resIndex
}
