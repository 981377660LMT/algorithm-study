package main

import (
	"fmt"
)

func main() {
	// check with sort
	nums := []int{1, 2}
	fmt.Println(bisectLeft(nums, 0, 0, len(nums)-1))
}

func bisectLeft(nums []int, x int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
