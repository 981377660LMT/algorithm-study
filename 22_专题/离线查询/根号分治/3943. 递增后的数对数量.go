// 3943. 递增后的数对数量
// https://leetcode.cn/problems/number-of-pairs-after-increment/description/
//
// 给你两个整数数组 nums1 和 nums2，以及一个二维整数数组 queries。
// 每个 queries[i] 都属于以下两种类型之一：
// [1, x, y, val]：将 nums2[x..y] 中的每个元素都 增加 val。
// [2, tot]：计算 满足 nums1[j] + nums2[k] == tot 的数对 (j, k) 的数量。
// 返回一个整数数组 answer，其中 answer[j] 表示第 jth 个类型 2 查询的数对数量。
// 1 <= nums1.length <= 5
// 1 <= nums2.length <= 5e4
//
// 卡哈希

package main

import "math"

func numberOfPairs(nums1 []int, nums2 []int, queries [][]int) []int {
	n1, n2 := len(nums1), len(nums2)
	blocks := useBlock(n2, int(math.Sqrt(float64(n2)))+1)
	belong, blockStart, blockEnd, blockCount := blocks.belong, blocks.blockStart, blocks.blockEnd, blocks.blockCount
	blockLazy := make([]int, blockCount)
	blockCounter := make([]map[int]int, blockCount)
	for i := 0; i < blockCount; i++ {
		blockCounter[i] = make(map[int]int)
		for j := blockStart[i]; j < blockEnd[i]; j++ {
			blockCounter[i][nums2[j]]++
		}
	}

	var res []int
	for _, query := range queries {
		if query[0] == 1 {
			l, r, val := query[1], query[2], query[3]
			bid1, bid2 := belong[l], belong[r]
			if bid1 == bid2 {
				for i := l; i <= r; i++ {
					blockCounter[bid1][nums2[i]]--
					nums2[i] += val
					blockCounter[bid1][nums2[i]]++
				}
			} else {
				for i := l; i < blockEnd[bid1]; i++ {
					blockCounter[bid1][nums2[i]]--
					nums2[i] += val
					blockCounter[bid1][nums2[i]]++
				}
				for i := bid1 + 1; i < bid2; i++ {
					blockLazy[i] += val
				}
				for i := blockStart[bid2]; i <= r; i++ {
					blockCounter[bid2][nums2[i]]--
					nums2[i] += val
					blockCounter[bid2][nums2[i]]++
				}
			}
		} else {
			tot := query[1]
			cur := 0
			for i := 0; i < n1; i++ {
				target := tot - nums1[i]
				for j := 0; j < blockCount; j++ {
					cur += blockCounter[j][target-blockLazy[j]]
				}
			}
			res = append(res, cur)
		}
	}

	return res
}

// blockSize := int(math.Sqrt(float64(len(nums))) + 1)
func useBlock(n int, blockSize int) struct {
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
