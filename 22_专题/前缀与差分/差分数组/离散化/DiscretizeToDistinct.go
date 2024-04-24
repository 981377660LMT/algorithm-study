package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	newNums, originValue := DiscretizeToDistinct(nums)
	fmt.Println(newNums, originValue)
	for i := range nums {
		if originValue[newNums[i]] != nums[i] {
			panic("error")
		}
	}
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func DiscretizeToDistinct(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for i, v := range order {
		origin[i] = nums[v]
		newNums[v] = int32(i)
	}
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}
