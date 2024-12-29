package main

import "fmt"

// 移除数组中连续重复的元素/有序数组原地去重.
// Compact/dedup.
func UniqueInplace(arr *([]int)) {
	if len(*arr) == 0 {
		return
	}
	nums := *arr
	size := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != nums[size] {
			size++
			nums[size] = nums[fast]
		}
	}
	size++
	*arr = nums[:size]
}

func main() {
	nums := []int{1, 2, 2, 3, 2}
	UniqueInplace(&nums)
	fmt.Println(nums)
}
