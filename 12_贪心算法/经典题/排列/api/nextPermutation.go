package main

// 返回下一个字典序的排列.
func nextPermutation(nums []int, inPlace bool) (res []int, ok bool) {
	if !inPlace {
		nums = append(nums[:0:0], nums...)
	}
	left, right := len(nums)-1, len(nums)-1
	for left > 0 && nums[left-1] >= nums[left] {
		left--
	}
	if left == 0 {
		return
	}
	last := left - 1
	for nums[right] <= nums[last] {
		right--
	}
	nums[last], nums[right] = nums[right], nums[last]
	for i, j := last+1, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return nums, true
}

// 返回上一个字典序的排列.
func prePermutation(nums []int, inPlace bool) (res []int, ok bool) {
	if !inPlace {
		nums = append(nums[:0:0], nums...)
	}
	left, right := len(nums)-1, len(nums)-1
	for left > 0 && nums[left-1] <= nums[left] {
		left--
	}
	if left == 0 {
		return
	}
	last := left - 1
	for nums[right] >= nums[last] {
		right--
	}
	nums[last], nums[right] = nums[right], nums[last]
	for i, j := last+1, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return nums, true
}
