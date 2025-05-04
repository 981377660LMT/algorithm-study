package main

// 求有序数组中位数(向下取整).
func GetMedian(sortedNums []int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(sortedNums) {
		end = len(sortedNums)
	}
	if start >= end {
		return 0
	}
	if (end-start)&1 == 0 {
		return (sortedNums[(end+start)/2-1] + sortedNums[(end+start)/2]) / 2
	}
	return sortedNums[(end+start)/2]
}
