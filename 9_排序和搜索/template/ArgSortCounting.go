package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{3, 1, 2, 4, 1}
	fmt.Println(ArgSortCounting(nums, 1, 4))   // Output: [1 2 0 3]
	fmt.Println(ArgSort(nums))                 // Output: [1 2 0 3]
	fmt.Println(ReArrage(nums, ArgSort(nums))) // Output: [1 2 3 4]
}

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// 返回一个数字数组的排序后的索引.内部使用计数排序.
func ArgSortCounting[T Int](nums []T, min, max T) []int {
	counter := make([]int32, max-min+1)
	for _, v := range nums {
		counter[v-min]++
	}
	for i := 1; i < len(counter); i++ {
		counter[i] += counter[i-1]
	}
	order := make([]int, len(nums))
	for i := len(nums) - 1; i >= 0; i-- { // 值相等时，按照下标从小到大排序
		v := nums[i] - min
		counter[v]--
		order[counter[v]] = i
	}
	return order
}

func ArgSort[T Int](nums []T) []int {
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func ReArrage[T any](arr []T, order []int) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = arr[order[i]]
	}
	return res
}
