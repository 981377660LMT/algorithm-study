package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10}
	BisectInsort(&arr, 5)
	fmt.Println(arr)

	fmt.Println(BisectFind(arr, 5, func(a, b int) bool { return a < b }))
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

// 返回给定元素应插入到此切片的索引位置。如果在该位置已经有相等元素，则返回 found = true。
func BisectFind[T any](s []T, item T, less func(T, T) bool) (index int, found bool) {
	i := sort.Search(len(s), func(i int) bool {
		return less(item, s[i])
	})
	if i > 0 && !less(s[i-1], item) {
		return i - 1, true
	}
	return i, false
}
