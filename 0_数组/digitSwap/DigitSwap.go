// DigitSwap/SwapDigit

package main

import "sort"

var base = [19]int{
	1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000, 100000000000,
	1000000000000, 10000000000000, 100000000000000, 1000000000000000, 10000000000000000, 100000000000000000, 1000000000000000000,
}

type DigitSwap struct {
	num    int
	digits []int
}

func NewDigitSwap(num int) *DigitSwap {
	if num < 0 {
		panic("num must be non-negative")
	}
	d := &DigitSwap{num: num}
	x := num
	for x > 0 {
		d.digits = append(d.digits, x%10)
		x /= 10
	}
	return d
}

// 交换两个数位.
//  D := NewDigitSwap(1234)
//  D.Swap(0, 1) // 1234 -> 1243
func (d *DigitSwap) Swap(i, j int) int {
	if i == j || d.digits[i] == d.digits[j] {
		return d.num
	}
	return d.num + (d.digits[j]-d.digits[i])*(base[i]-base[j])
}

func (d *DigitSwap) At(i int) int {
	return d.digits[i]
}

func (d *DigitSwap) Len() int {
	return len(d.digits)
}

// 3267. 统计近似相等数对 II
// https://leetcode.cn/problems/count-almost-equal-pairs-ii/
// 给你一个正整数数组 nums 。
// 如果我们执行以下操作 `至多两次` 可以让两个整数 x 和 y 相等，那么我们称这个数对是 近似相等 的：
// 选择 x 或者 y  之一，将这个数字中的两个数位交换。
// 请你返回 nums 中，下标 i 和 j 满足 i < j 且 nums[i] 和 nums[j] 近似相等 的数对数目。
// 注意 ，执行操作后得到的整数可以有前导 0 。
//
// !排序后，一定是换后面那个数
func countPairs(nums []int) int {
	swap2 := func(x int) map[int]struct{} {
		res := map[int]struct{}{x: {}}
		queue := map[int]struct{}{x: {}}
		for t := 0; t < 2; t++ {
			nextQueue := map[int]struct{}{}
			for num := range queue {
				D := NewDigitSwap(num)
				for i := 0; i < D.Len(); i++ {
					for j := i + 1; j < D.Len(); j++ {
						next := D.Swap(i, j)
						nextQueue[next] = struct{}{}
						res[next] = struct{}{}
					}
				}
			}
			queue = nextQueue
		}
		return res
	}

	sort.Ints(nums)
	res := 0
	preCounter := map[int]int{}
	for i := 0; i < len(nums); i++ {
		swapRes := swap2(nums[i])
		for num := range swapRes {
			res += preCounter[num]
		}
		preCounter[nums[i]]++
	}

	return res
}
