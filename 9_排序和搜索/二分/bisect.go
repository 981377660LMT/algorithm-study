package main

import (
	"fmt"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10}
	BisectInsort(&arr, 5)
	fmt.Println(arr)
}

// sort.SearchInts in go.
//
// sort.Search(n, func(i int) bool { return nums[i] >= target }) in go.
func BisectLeft(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

// sort.Search(n, func(i int) bool { return nums[i] > target }) in go.
func BisectRight(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

// bisect.insort/bisect.insortRight in python.
func BisectInsort(nums *[]int, insertion int) {
	pos := BisectRight(*nums, insertion)
	*nums = append(*nums, 0)
	copy((*nums)[pos+1:], (*nums)[pos:])
	(*nums)[pos] = insertion
}

var InsortRight = BisectInsort
