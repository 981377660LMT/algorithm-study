// 3430. 最多 K 个元素的子数组的最值之和
// https://leetcode.cn/problems/maximum-and-minimum-sums-of-at-most-size-k-subarrays/description/
// 给你一个整数数组 nums 和一个 正 整数 k 。 返回 最多 有 k 个元素的所有子数组的 最大 和 最小 元素之和。

package main

func minMaxSubarraySum(nums []int, k int) int64 {
	// 左侧可选长度为L，右侧可选长度为R，长度不超过k的非空子数组个数.
	// 左侧和右侧都包含当前元素.
	countInRange := func(left, right, k int) int {
		upper := right
		if upper > k {
			upper = k
		}
		if upper <= 0 {
			return 0
		}

		if left > k {
			return (k + k - upper + 1) * upper / 2
		}

		pos := k - left + 1
		if pos > upper {
			return upper * left
		}

		c1 := pos - 1
		res1 := c1 * left
		c2 := upper - pos + 1
		min_ := k - (upper - 1)
		max_ := k - (pos - 1)
		res2 := (min_ + max_) * c2 / 2
		return res1 + res2
	}

	calc := func(nums []int, isMax bool, k int) int {
		n := len(nums)
		leftMost, rightMost := GetRange(nums, isMax, true, false)
		res := 0
		for i := 0; i < n; i++ {
			L := i - leftMost[i] + 1
			R := rightMost[i] - i + 1
			count := countInRange(L, R, k)
			res += nums[i] * count
		}
		return res
	}

	sum1 := calc(nums, false, k)
	sum2 := calc(nums, true, k)
	return int64(sum1 + sum2)
}

func GetRange(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int) {
	compareLeft := func(stackValue, curValue int) bool {
		if isLeftStrict && isMax {
			return stackValue <= curValue
		} else if isLeftStrict && !isMax {
			return stackValue >= curValue
		} else if !isLeftStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	compareRight := func(stackValue, curValue int) bool {
		if isRightStrict && isMax {
			return stackValue <= curValue
		} else if isRightStrict && !isMax {
			return stackValue >= curValue
		} else if !isRightStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	n := len(nums)
	leftMost = make([]int, n)
	rightMost = make([]int, n)

	for i := 0; i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			rightMost[top] = i - 1
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			leftMost[top] = i + 1
		}
		stack = append(stack, i)
	}

	return
}
