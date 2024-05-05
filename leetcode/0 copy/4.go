package main

import "sort"

const MOD int = 1e18

func minCostToEqualizeArray(nums []int, cost1 int, cost2 int) int {
	sort.Ints(nums)
	max_ := nums[len(nums)-1]
	if cost1*2 <= cost2 {
		diffSum := 0
		for _, num := range nums {
			diffSum += max_ - num
		}
		return (diffSum * cost1) % MOD
	}

	distToTarget := func(target int) int {
		diff := make([]int, 0, len(nums))
		diffSum := 0
		for _, num := range nums {
			if num < target {
				diff = append(diff, target-num)
				diffSum += target - num
			}
		}
		if cost1*2 <= cost2 {
			return diffSum * cost1
		}
		if len(diff) == 0 {
			return 0
		}
		if len(diff) == 1 {
			return diff[0] * cost1
		}
		maxDiff := diff[0]
		otherDiff := diffSum - maxDiff
		if maxDiff <= otherDiff {
			if diffSum&1 == 0 {
				return diffSum / 2 * cost2
			}
			return diffSum/2*cost2 + cost1
		}
		return otherDiff*cost2 + (maxDiff-otherDiff)*cost1
	}

	_, y := FibonacciSearch(distToTarget, max_, 1e16, true)
	return y % MOD
}

const INF int = 1e18

// 寻找[left,right]中的一个极值点,不要求单峰性质.
//
//	返回值: (极值点,极值)
func FibonacciSearch(f func(x int) int, left, right int, minimize bool) (int, int) {
	a, b, c, d := left, left+1, left+2, left+3
	step := 0
	for d < right {
		b = c
		c = d
		d = b + c - a
		step++
	}

	get := func(i int) int {
		if right < i {
			return INF
		}
		if minimize {
			return f(i)
		}
		return -f(i)
	}

	ya, yb, yc, yd := get(a), get(b), get(c), get(d)
	for i := 0; i < step; i++ {
		if yb < yc {
			d = c
			c = b
			b = a + d - c
			yd = yc
			yc = yb
			yb = get(b)
		} else {
			a = b
			b = c
			c = a + d - b
			ya = yb
			yb = yc
			yc = get(c)
		}
	}

	x := a
	y := ya
	if yb < y {
		x = b
		y = yb
	}
	if yc < y {
		x = c
		y = yc
	}
	if yd < y {
		x = d
		y = yd
	}

	if minimize {
		return x, y
	}
	return x, -y
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
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

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
