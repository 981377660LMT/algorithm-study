// 返回nums的各个子集的元素和.
// !比求出所有的子集的元素和再排序要快很多

package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	nums := make([]int, 25)
	for i := range nums {
		nums[i] = i
	}
	time1 := time.Now()
	res := SubsetSum(nums, false)
	sort.Ints(res)
	time2 := time.Now()
	SubsetSum(nums, true)
	time3 := time.Now()
	SubsetSumSortedWithState(nums)
	time4 := time.Now()

	fmt.Println(time2.Sub(time1)) // 1.5s
	fmt.Println(time3.Sub(time2)) // 300ms
	fmt.Println(time4.Sub(time3)) // 500ms
}

// O(2^n)返回nums的各个子集的元素和.
//
//	sorted: 是否返回排序后的结果.
func SubsetSum(nums []int, sorted bool) []int {
	if sorted {
		return _subsetSumSorted(nums)
	}
	return _subsetSumUnsorted(nums)
}

func _subsetSumUnsorted(nums []int) []int {
	n := len(nums)
	dp := make([]int, 1<<n)
	for i := 0; i < n; i++ {
		for pre := 0; pre < 1<<i; pre++ {
			dp[pre|(1<<i)] = dp[pre] + nums[i]
		}
	}
	return dp
}

// O(2^n)返回nums的各个子集的元素和的排序后的结果.
// !比求出所有的子集的元素和再排序要快很多
func _subsetSumSorted(nums []int) []int {
	dp := []int{0}
	for _, v := range nums {
		ndp := make([]int, len(dp))
		for i, w := range dp {
			ndp[i] = w + v
		}
		dp = merge(dp, ndp)
	}
	return dp
}

func merge(a, b []int) []int {
	n1, n2 := len(a), len(b)
	res := make([]int, n1+n2)
	i, j, k := 0, 0, 0
	for i < n1 && j < n2 {
		if a[i] < b[j] {
			res[k] = a[i]
			i++
		} else {
			res[k] = b[j]
			j++
		}
		k++
	}
	for i < n1 {
		res[k] = a[i]
		i++
		k++
	}
	for j < n2 {
		res[k] = b[j]
		j++
		k++
	}
	return res
}

// O(2^n)返回nums的各个子集的元素和的排序后的结果, 并且记录状态.
func SubsetSumSortedWithState(nums []int) [][2]int {
	merge := func(a, b [][2]int) [][2]int {
		n1, n2 := len(a), len(b)
		res := make([][2]int, n1+n2)
		i, j, k := 0, 0, 0
		for i < n1 && j < n2 {
			if a[i][0] < b[j][0] {
				res[k] = a[i]
				i++
			} else {
				res[k] = b[j]
				j++
			}
			k++
		}
		for i < n1 {
			res[k] = a[i]
			i++
			k++
		}
		for j < n2 {
			res[k] = b[j]
			j++
			k++
		}
		return res
	}

	dp := [][2]int{{0, 0}}
	for i, x := range nums {
		ndp := make([][2]int, len(dp))
		for j, p := range dp {
			ndp[j][0] = p[0] + x
			ndp[j][1] = p[1] | 1<<i
		}
		dp = merge(dp, ndp)
	}

	return dp
}
