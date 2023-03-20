// # 给你一个由正整数组成的数组 nums 和一个 正 整数 k 。

// # 如果 nums 的子集中，任意两个整数的绝对差均不等于 k ，则认为该子数组是一个 美丽 子集。

// # 返回数组 nums 中 非空 且 美丽 的子集数目。

// # nums 的子集定义为：可以经由 nums 删除某些元素（也可能不删除）得到的一个数组。只有在删除元素时选择的索引不同的情况下，两个子集才会被视作是不同的子集。

package main

import "runtime/debug"

func init() {
	debug.SetGCPercent(-1)
}

func beautifulSubsets(nums []int, k int) int {
	res := 0
	for s := 1; s < 1<<len(nums); s++ {
		cur := []int{}
		for i, v := range nums {
			if s>>i&1 == 1 {
				cur = append(cur, v)
			}
		}

		flag := true
		for i := 0; i < len(cur); i++ {
			for j := i + 1; j < len(cur); j++ {
				if abs(cur[i]-cur[j]) == k {
					flag = false
					break
				}
			}
			if !flag {
				break
			}
		}

		if flag {
			res++
		}
	}

	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
