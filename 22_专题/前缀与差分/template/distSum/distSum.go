package main

import (
	"fmt"
	"sort"
)

func main() {
	sortedNums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	Q := DistSum(sortedNums)
	fmt.Println(Q(5))
}

// 有序数组所有点到`x=k`的距离之和.
func DistSum(sortedNums []int) func(k int) int {
	n := len(sortedNums)
	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + sortedNums[i]
	}

	return func(k int) int {
		pos := sort.SearchInts(sortedNums, k+1)
		leftSum := k*pos - preSum[pos]
		rightSum := preSum[n] - preSum[pos] - k*(n-pos)
		return leftSum + rightSum
	}
}

// 有序数组中所有点对两两距离之和.一共有`n*(n-1)//2`对点对.
func DistSumOfAllPairs(sortedNums []int) int {
	n := len(sortedNums)
	res := 0
	preSum := 0
	for i := 0; i < n; i++ {
		res += sortedNums[i]*i - preSum
		preSum += sortedNums[i]
	}
	return res
}
