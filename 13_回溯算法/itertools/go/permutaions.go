package main

import (
	"fmt"
	"slices"

	"time"
)

func main() {

	time1 := time.Now()
	count := 0

	EnumeratePermutations(11, 11, func(perm []int) bool {
		count++
		return false
	})

	EnumeratePermutationsAll(11, func(indicesView []int) bool {
		count++
		return false
	})
	fmt.Println(count)
	fmt.Println(time.Since(time1))

	perm := []int{3, 1, 2}
	for NextPermutation(perm) {
		fmt.Println(perm)
	}
	fmt.Println(perm)
}

// 从一个长度为 n 的数组中选择 r 个元素，按字典序生成所有排列，每个排列用下标表示.
// https://docs.python.org/3/library/itertools.html#itertools.permutations
// 11!(4e7) => 420ms.
func EnumeratePermutations(n, r int, f func(indicesView []int) (shouldBreak bool)) {
	ids := make([]int, n)
	for i := range ids {
		ids[i] = i
	}
	if f(ids[:r]) {
		return
	}
	cycles := make([]int, r)
	for i := range cycles {
		cycles[i] = n - i
	}
	for {
		i := r - 1
		for ; i >= 0; i-- {
			cycles[i]--
			if cycles[i] == 0 {
				tmp := ids[i]
				copy(ids[i:], ids[i+1:])
				ids[n-1] = tmp
				cycles[i] = n - i
			} else {
				j := cycles[i]
				ids[i], ids[n-j] = ids[n-j], ids[i]
				if f(ids[:r]) {
					return
				}
				break
			}
		}
		if i == -1 {
			return
		}
	}
}

// 生成全排列(不保证字典序).
// 11!(4e7) => 380ms.
func EnumeratePermutationsAll(n int, f func(indicesView []int) (shouldBreak bool)) {
	var dfs func([]int, int) bool
	dfs = func(a []int, i int) bool {
		if i == len(a) {
			return f(a)
		}
		dfs(a, i+1)
		for j := i + 1; j < len(a); j++ {
			a[i], a[j] = a[j], a[i]
			if dfs(a, i+1) {
				return true
			}
			a[i], a[j] = a[j], a[i]
		}
		return false
	}
	ids := make([]int, n)
	for i := range ids {
		ids[i] = i
	}
	dfs(ids, 0)
}

// 调用完之后
// 返回 true：a 修改为其下一个排列（即比 a 大且字典序最小的排列）
// 返回 false：a 修改为其字典序最小的排列（即 a 排序后的结果）
// - [31. 下一个排列](https://leetcode.cn/problems/next-permutation/)
// - [1850. 邻位交换的最小次数](https://leetcode.cn/problems/minimum-adjacent-swaps-to-reach-the-kth-smallest-number/) 2073
func NextPermutation(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	defer slices.Reverse(a[i+1:])
	if i < 0 {
		return false
	}
	j := n - 1
	for j >= 0 && a[i] >= a[j] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	return true
}
