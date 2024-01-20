// 启发式分治.

package main

import "fmt"

// !给定一个长度为n的数组，问对所有的子数组中的每个元素，是否满足 check 条件.
// n : 数组长度
// check(left, right, curIndex) : 检查[left, right]闭区间内索引为curIndex的元素是否满足条件.
// O(nlogn)
func CheckAllSubarray(n int, check func(left, right int, curIndex int) bool) bool {
	var dfs func(left, right int) bool
	dfs = func(left, right int) bool {
		if left >= right {
			return true
		}
		for x, y := left, right; x <= y; x, y = x+1, y-1 {
			if check(left, right, x) {
				return dfs(left, x-1) && dfs(x+1, right)
			}
			if check(left, right, y) {
				return dfs(left, y-1) && dfs(y+1, right)
			}
		}
		return false
	}
	return dfs(0, n-1)
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(CheckAllSubarray(5, func(left, right int, curIndex int) bool {
		return nums[curIndex] > 4
	}))
}
