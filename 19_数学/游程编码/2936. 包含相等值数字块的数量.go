// !所有相同值的元素都是相邻的.
// https://leetcode.cn/problems/number-of-equal-numbers-blocks/description/

package main

import "sort"

type BigArray interface {
	At(int64) int
	Size() int64
}

func countBlocks(nums BigArray) int {
	res := 0
	n := int(nums.Size())
	ptr := 0
	for ptr < n {
		cur := nums.At(int64(ptr))
		sameCount := sort.Search(n-ptr, func(i int) bool { return nums.At(int64(ptr+i)) != cur })
		ptr += sameCount
		res++
	}
	return res
}
