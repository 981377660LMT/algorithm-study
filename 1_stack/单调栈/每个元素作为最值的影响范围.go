package main

const MOD int = 1e9 + 7

// 求所有子数组的最小值之和
func sumSubarrayMins(arr []int) int {
	left, right := GetRange(arr, false, true, false)
	res := 0
	for i, num := range arr {
		count := (right[i] - i + 1) * (i - left[i] + 1) // 出现在了多少个子数组里
		res = (res + num*count) % MOD
	}
	return res
}

// 求每个元素作为最值的影响范围(区间)
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
	leftMost, rightMost = make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
}

// 求每个元素作为最值的影响范围(区间)
func GetRange32(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int32) {
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

	n := int32(len(nums))
	leftMost, rightMost = make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int32{}
	for i := int32(0); i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
}
