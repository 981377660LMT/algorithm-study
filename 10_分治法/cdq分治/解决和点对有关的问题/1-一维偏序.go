// cdq分治：基于归并排序的分治

package main

import "fmt"

func main() {
	nums := []int{3, 2, 1, 4, 5}
	fmt.Println(CountInv(nums))
	MergeSort(nums)
	fmt.Println(nums)
}

// https://leetcode.cn/problems/shu-zu-zhong-de-ni-xu-dui-lcof/
func MergeSort(nums []int) {
	n := len(nums)
	tmp := make([]int, n)

	var f func(start, end int)
	f = func(start, end int) {
		if start+1 >= end {
			return
		}
		mid := (start + end) >> 1
		f(start, mid)
		f(mid, end)

		i, j, k := start, mid, 0
		for i < mid && j < end {
			if nums[i] <= nums[j] {
				tmp[k] = nums[i]
				i++
			} else {
				tmp[k] = nums[j]
				j++
			}
			k++
		}
		for i < mid {
			tmp[k] = nums[i]
			i++
			k++
		}
		for j < end {
			tmp[k] = nums[j]
			j++
			k++
		}
		copy(nums[start:end], tmp[:k])
	}
	f(0, n)
}

func CountInv(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	n := len(nums)
	tmp := make([]int, n)
	res := 0

	var f func(start, end int)
	f = func(start, end int) {
		if start+1 >= end {
			return
		}
		mid := (start + end) >> 1
		f(start, mid)
		f(mid, end)

		i, j, k := start, mid, 0
		for i < mid && j < end {
			if nums[i] <= nums[j] {
				tmp[k] = nums[i]
				i++
			} else {
				tmp[k] = nums[j]
				j++
				res += mid - i // !注意这里
			}
			k++
		}
		for i < mid {
			tmp[k] = nums[i]
			i++
			k++
		}
		for j < end {
			tmp[k] = nums[j]
			j++
			k++
		}
		copy(nums[start:end], tmp[:k])
	}
	f(0, n)
	return res
}
