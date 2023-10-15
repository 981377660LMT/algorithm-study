package main

import "math"

// 2863. 最长半递减数组
// https://leetcode.cn/problems/maximum-length-of-semi-decreasing-subarrays/description/
// !返回最长的子数组，子数组的第一个元素严格大于最后一个元素
// 分块加速查找.分块的好处是，可以支持修改操作.
func maxSubarrayLength(nums []int) int {
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockMin := make([]int, blockCount)
	for i := range blockMin {
		blockMin[i] = INF
	}
	for i, v := range nums {
		blockMin[belong[i]] = min(blockMin[belong[i]], v)
	}

	queryRightMostLower := func(pos int) (res int, ok bool) {
		bid := belong[pos]
		target := nums[pos]
		for i := blockCount - 1; i > bid; i-- {
			if blockMin[i] >= target {
				continue
			}
			for j := blockEnd[i] - 1; j >= blockStart[i]; j-- {
				if nums[j] < target {
					return j, true
				}
			}
		}

		for i := blockEnd[bid] - 1; i > pos; i-- {
			if nums[i] < target {
				return i, true
			}
		}

		return 0, false
	}

	res := 0
	for i := range nums {
		if rightMostLower, ok := queryRightMostLower(i); ok {
			res = max(res, rightMostLower-i+1)
		}
	}
	return res
}

const INF int = 1e18

// blockSize := int(math.Sqrt(float64(len(nums))) + 1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
