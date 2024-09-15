// MinLeftMaxRight

package main

import "fmt"

func main() {
	maxRightBruteForce := func(left int, check func(right int) bool, upper int) int {
		for right := upper; right > left; right-- {
			if check(right) {
				return right
			}
		}
		return left
	}

	minLeftBruteForce := func(right int, check func(left int) bool, lower int) int {
		for left := lower; left < right; left++ {
			if check(left) {
				return left
			}
		}
		return right
	}

	{
		// test MaxRight
		check := func(right int) bool {
			return right >= 5
		}

		upper := 10
		if got, want := MaxRight(0, check, upper), maxRightBruteForce(0, check, upper); got != want {
			panic("MaxRight")
		}
	}

	{
		// test MinLeft
		check := func(left int) bool {
			return left >= 5
		}

		lower := 0
		if got, want := MinLeft(10, check, lower), minLeftBruteForce(10, check, lower); got != want {
			panic("MinLeft")
		}
	}

	{
		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		preSum := make([]int, len(nums)+1)
		for i, v := range nums {
			preSum[i+1] = preSum[i] + v
		}

		fmt.Println(MaxRight(0, func(right int) bool { return preSum[right]-preSum[0] < 6 }, len(nums))) // 2
		fmt.Println(MaxRight(0, func(right int) bool { return preSum[right]-preSum[0] <= 6 }, 2))        // 2
		fmt.Println(MaxRight(0, func(right int) bool { return preSum[right]-preSum[0] <= 6 }, 3))        // 3

		fmt.Println(MinLeft(len(nums), func(left int) bool { return preSum[len(nums)]-preSum[left] < 20 }, 0))  // 8
		fmt.Println(MinLeft(len(nums), func(left int) bool { return preSum[len(nums)]-preSum[left] < 27 }, 0))  // 8
		fmt.Println(MinLeft(len(nums), func(left int) bool { return preSum[len(nums)]-preSum[left] <= 27 }, 0)) // 7
		fmt.Println(MinLeft(len(nums), func(left int) bool { return preSum[len(nums)]-preSum[left] <= 27 }, 8)) // 8
		fmt.Println(MinLeft(len(nums), func(left int) bool { return preSum[len(nums)]-preSum[left] <= 27 }, 7)) // 7

	}
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含，使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
