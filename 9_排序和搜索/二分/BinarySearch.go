package main

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// 查询非递减数组中目标值的第一个或最后一个位置.
// arr: 非递减数组
// target: 目标值
// findFirst: 是否查询第一个位置. 默认为 true.
// 返回目标值的第一个或最后一个位置. 如果目标值不存在, 返回-1.
func BinarySearch[T Int](arr []T, target T, findFirst bool) int {
	if len(arr) == 0 || arr[0] > target || arr[len(arr)-1] < target {
		return -1
	}
	if findFirst {
		left, right := 0, len(arr)-1
		for left <= right {
			mid := left + (right-left)>>1
			if arr[mid] < target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		if left < len(arr) && arr[left] == target {
			return left
		}
		return -1
	} else {
		left, right := 0, len(arr)-1
		for left <= right {
			mid := left + (right-left)>>1
			if arr[mid] <= target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		if left > 0 && arr[left-1] == target {
			return left - 1
		}
		return -1
	}
}

// https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/submissions/504998244/
func searchRange(nums []int, target int) []int {
	return []int{BinarySearch(nums, target, true), BinarySearch(nums, target, false)}
}
