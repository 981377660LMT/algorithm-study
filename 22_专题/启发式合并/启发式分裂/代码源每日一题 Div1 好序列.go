// http://oj.daimayuan.top/problem/613
// https://zhuanlan.zhihu.com/p/487136369
// https://blog.csdn.net/weixin_51216553/article/details/128007971
// !给你一个数组a，问是否满足他的每个子区间[l,r]满足，至少存在一个元素x仅出现了一次
// 启发式分治

package main

import (
	"fmt"
)

func main() {
	// 5
	// 1 2 3 4 5
	// 5
	// 1 1 1 1 1
	// 5
	// 1 2 3 2 1
	// 5
	// 1 1 2 1 1
	fmt.Println(Solve([]int{1, 2, 3, 4, 5}))
	fmt.Println(Solve([]int{1, 1, 1, 1, 1}))
	fmt.Println(Solve([]int{1, 2, 3, 2, 1}))
	fmt.Println(Solve([]int{1, 1, 2, 1, 1}))

	count := 1
	CheckAllSubarray(1e5, func(left, right int, curPos int) bool {
		count++
		return true
	})
	fmt.Println(count)
}

func Solve(nums []int) bool {
	n := len(nums)
	mp := make(map[int][]int)
	for i, v := range nums {
		mp[v] = append(mp[v], i)
	}
	pre, next := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		pre[i] = -1
		next[i] = n
	}
	for _, pos := range mp {
		for i, v := range pos {
			if i > 0 {
				pre[v] = pos[i-1]
			}
			if i < len(pos)-1 {
				next[v] = pos[i+1]
			}
		}
	}

	return CheckAllSubarray(n, func(left, right int, curPos int) bool {
		return pre[curPos] < left && right < next[curPos]
	})
}

// n : 数组长度
// check(left, right, curPos) : 检查[left, right]区间内索引为curPos的元素是否满足条件.
func CheckAllSubarray(n int, check func(left, right int, curPos int) bool) bool {
	var dfs func(left, right int) bool
	dfs = func(left, right int) bool {
		if left >= right {
			return true
		}
		x, y := left, right
		for x <= y {
			if check(left, right, x) {
				return dfs(left, x-1) && dfs(x+1, right)
			}
			if check(left, right, y) {
				return dfs(left, y-1) && dfs(y+1, right)
			}
			x++
			y--
		}
		return false
	}
	return dfs(0, n-1)
}
