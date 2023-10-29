package main

// https://leetcode.cn/problems/maximum-or/
func maximumOr(nums []int, k int) int64 {
	mul := ProductWithoutOne(nums, func() int { return 0 }, func(a, b int) int { return a | b })
	res := int64(0)
	for i := 0; i < len(nums); i++ {
		cur := int64(nums[i]<<k | mul[i])
		if cur > res {
			res = cur
		}
	}
	return res
}

func constructProductMatrix(grid [][]int) [][]int {
	return ProductWithoutOne2D(grid, func() int { return 1 }, func(a, b int) int { return a * b % 12345 })
}

type E = int

// 除自身以外数组的乘积.
func ProductWithoutOne(nums []E, e func() E, op func(E, E) E) []E {
	n := len(nums)
	res := make([]E, n)
	for i := 0; i < n-1; i++ {
		res[i+1] = op(res[i], nums[i])
	}
	x := e()
	for i := n - 1; i >= 0; i-- {
		res[i] = op(res[i], x)
		x = op(nums[i], x)
	}
	return res
}

func ProductWithoutOne2D(nums [][]E, e func() E, op func(E, E) E) [][]E {
	row := len(nums)
	col := len(nums[0])

	unit := e()
	res := make([][]E, row)
	for i := range res {
		res[i] = make([]E, col)
		for j := range res[i] {
			res[i][j] = unit
		}
	}

	for i := row - 1; i >= 0; i-- {
		tmp1 := res[i]
		tmp2 := nums[i]
		for j := col - 1; j >= 0; j-- {
			tmp1[j] = unit
			unit = op(tmp2[j], unit)
		}
	}

	unit = e()
	for i := 0; i < row; i++ {
		tmp1 := res[i]
		tmp2 := nums[i]
		for j := 0; j < col; j++ {
			tmp1[j] = op(unit, tmp1[j])
			unit = op(unit, tmp2[j])
		}
	}

	return res
}
