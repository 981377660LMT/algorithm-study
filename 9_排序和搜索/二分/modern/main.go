package main

import "fmt"

func main() {
	// 有序数组二分查找
	{
		arr := []int{1, 3, 5, 5, 5, 6, 7, 9}
		target := 5
		index := BinarySearch(func(mid int) bool { return arr[mid] <= target }, -1, len(arr))
		if index == -1 {
			fmt.Println("没有找到小于等于", target, "的元素")
		} else {
			fmt.Printf("最后一个小于等于 %d 的元素是 %d，位于索引 %d\n", target, arr[index], index)
		}
	}

	// x的平方根
	{
		target := 3.0
		res := BinarySearchReal(func(mid float64) bool { return mid*mid <= target }, 0, target, 100)
		fmt.Println("平方根:", res)
	}
}

// 整数二分, 返回搜索区间内最后一个满足 check 函数为 true 的整数.
//
//	check: 函数满足单调性（当 x 从小到大变化时，结果从 true 到 false 或从 false 到 true 仅发生一次变化）.
//	ok: 初始一个满足 check(x) == true 的值.
//	ng: 初始一个满足 check(x) == false 的值.
func BinarySearch(check func(int) bool, ok, ng int) int {
	for abs(ok-ng) > 1 {
		mid := (ok + ng) / 2
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 浮点数二分, 返回搜索区间内最后一个满足 check 函数为 true 的浮点数.
//
//	check: 函数满足单调性（当 x 从小到大变化时，结果从 true 到 false 或从 false 到 true 仅发生一次变化）.
//	ok: 初始一个满足 check(x) == true 的值.
//	ng: 初始一个满足 check(x) == false 的值.
//	iter: 二分迭代次数, 迭代次数越多, 精度越高.
func BinarySearchReal(check func(float64) bool, ok, ng float64, iter int) float64 {
	for i := 0; i < iter; i++ {
		mid := (ok + ng) / 2
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return (ok + ng) / 2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
