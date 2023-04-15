// 1982. 从子集的和还原数组
// https://leetcode.cn/problems/find-array-given-subset-sums/
// n<=15 1e4<=sums[i]<=1e4

package main

import (
	"math/bits"
	"sort"
)

// 先给所有数加上 -min_element 转化成非负问题
// 非负问题中，最小值必为 0，次小值就是原集合最小值，递归求解即可。
// 最后再直接枚举子集，如果正好和为 -min_element，那么这些数就是原集合中的负数。
func recoverArray(n int, sums []int) []int {
	min_ := mins(sums...)
	for i := range sums {
		sums[i] += -min_
	}

	res := recoverArrayFromSubsetSum(sums)
	for state := 0; state < 1<<n; state++ { // 寻找原数组中的负数
		sum_ := 0
		for i := 0; i < n; i++ {
			if state>>i&1 == 1 {
				sum_ += res[i]
			}
		}
		if sum_ == -min_ {
			for i := 0; i < n; i++ {
				if state>>i&1 == 1 {
					res[i] = -res[i]
				}
			}
			break

		}
	}
	return res
}

// 给出由非负整数组成的数组 a 的子集和 sum，返回 a
// https://github.dev/EndlessCheng/codeforces-go/blob/029f576c04914ad4052dbe1073ff644dc219824a/copypasta/common.go
func recoverArrayFromSubsetSum(sums []int) []int {
	sort.Ints(sums)
	n := bits.TrailingZeros(uint(len(sums)))
	skip := map[int]int{}
	res := make([]int, 0, n)
	for j := 0; n > 0; n-- {
		for j++; skip[sums[j]] > 0; j++ {
			skip[sums[j]]--
		}
		s := sums[j]
		_s := make([]int, 1<<len(res))
		for i, v := range res {
			for m, b := 0, 1<<i; m < b; m++ {
				_s[b|m] = _s[m] + v
				skip[_s[b|m]+s]++
			}
		}
		res = append(res, s)
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}
