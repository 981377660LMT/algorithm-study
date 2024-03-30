// https://taodaling.github.io/blog/2019/09/12/%E4%BA%8C%E5%88%86/
// 现在很多输出浮点数的题目都会提供两种 AC 条件，
// 一种是输出与真实结果的绝对误差不超过阈值，
// 一种是输出与真实结果的相对误差不超过阈值。
// 相对误差在输出结果很大的时候会发挥巨大的作用。
// 众所周知双精度浮点型共 64 位，其中 1 位用于表示符号，
// 11 位表示指数，其余的 52 位用于表示有效数字。
// 简单换算就可以知道双精度浮点型可以精确表示大概 15 位十进制整数。
// 现在考虑一个问题，最终结果为 10^8 ，但是要求绝对误差小于 10^−8，这现实吗。
// 事实上尾部的数值由于有效数值不足会被舍去。
// 这就会导致二分的时候，(l+r)/2 可能会等于 l 或等于 r，从而导致二分进入死循环。
// 但是有了相对误差，情况就会大为不同，当输出为 10^8 时，我们可以不需要保留任意小数。

package main

import "fmt"

func main() {
	fmt.Println(FirstTrueInt(func(mid int) bool { return mid >= 7 }, 0, 100))
	fmt.Println(LastTrueInt(func(mid int) bool { return mid <= 76 }, 0, 100))
	fmt.Println(LastTrueFloat64(func(mid float64) bool { return mid >= 7 }, 0, 100, 1e-8, -1))  // 绝对误差不超过1e-8
	fmt.Println(FirstTrueFloat64(func(mid float64) bool { return mid >= 7 }, 0, 100, -1, 0.01)) // 相对误差不超过1%
	fmt.Println(LowerBound([]int{1, 2, 3, 4, 5}, 3, 0, 5))
	fmt.Println(UpperBound([]int{1, 2, 3, 4, 5}, 3, 0, 5))
}

func FirstTrueInt(predicate func(mid int) bool, left, right int) int {
	for left != right {
		mid := (left + right) >> 1
		if predicate(mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	if !predicate(left) {
		left++
	}
	return left
}

func LastTrueInt(predicate func(mid int) bool, left, right int) int {
	for left != right {
		mid := (left + right + 1) >> 1
		if predicate(mid) {
			left = mid
		} else {
			right = mid - 1
		}
	}
	if !predicate(left) {
		left--
	}
	return left
}

// 浮点数二分，在区间 [left, right] 上找到最小的 mid 使得 predicate(mid) 为 true.
//
//	absError: 绝对误差.-1表示不使用绝对误差.
//	relativeError: 相对误差.-1表示不使用相对误差.
func FirstTrueFloat64(predicate func(mid float64) bool, left, right, absError, relativeError float64) float64 {
	if absError == -1 && relativeError == -1 {
		panic("absError and relativeError can't be both -1")
	}

	if absError == -1 {
		for right > left {
			if (right < 0 && (right-left) < -right*relativeError) || (left > 0 && (right-left) < left*relativeError) {
				break
			}
			mid := (left + right) / 2
			if predicate(mid) {
				right = mid
			} else {
				left = mid
			}
		}
		return (left + right) / 2
	}

	if relativeError == -1 {
		for right-left > absError {
			mid := (left + right) / 2
			if predicate(mid) {
				right = mid
			} else {
				left = mid
			}
		}
		return (left + right) / 2
	}

	for right-left > absError {
		if (right < 0 && (right-left) < -right*relativeError) || (left > 0 && (right-left) < left*relativeError) {
			break
		}
		mid := (left + right) / 2
		if predicate(mid) {
			right = mid
		} else {
			left = mid
		}
	}
	return (left + right) / 2
}

func LastTrueFloat64(predicate func(mid float64) bool, left, right, absError, relativeError float64) float64 {
	return FirstTrueFloat64(predicate, left, right, absError, relativeError)
}

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

func LowerBound[T Int](arr []T, target T, left, right int) int {
	for left < right {
		mid := (left + right) >> 1
		if arr[mid] >= target {
			right = mid
		} else {
			left = mid + 1
		}
	}
	if arr[left] < target {
		left++
	}
	return left
}

func UpperBound[T Int](arr []T, target T, left, right int) int {
	for left < right {
		mid := (left + right) >> 1
		if arr[mid] > target {
			right = mid
		} else {
			left = mid + 1
		}
	}
	if arr[left] <= target {
		left++
	}
	return left
}
