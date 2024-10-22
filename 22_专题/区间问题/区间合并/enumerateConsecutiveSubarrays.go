package main

func main() {
	EnumetateConsecutiveSubarrays([]int{1, 2, 3, 5, 6, 7, 9}, 1, func(start, end int) {
		println(start, end)
	})
}

// 遍历连续的子数组
func EnumetateConsecutiveSubarrays(nums []int, diff int, f func(start, end int)) {
	i, n := 0, len(nums)
	for i < n {
		start := i
		for i < n-1 && nums[i]+diff == nums[i+1] {
			i += 1
		}
		i++
		f(start, i)
	}
}
