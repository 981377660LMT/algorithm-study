// https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/

package main

// 线段树分治的特殊情形.
// 调用 `query` 时，`state` 为对除了 `index` 以外所有点均调用过了 `mutate` 的状态。但不保证调用 `mutate` 的顺序。
// 总计会调用 $O(NlgN)$ 次的 `mutate` 和 `undo`, 以及 $O(N)$ 次的 `query`.
func MutateWithoutOneUndo(
	start, end int,
	/** 这里的 index 也就是 time. */
	mutate func(index int),
	undo func(),
	query func(index int),
) {
	var dfs func(curStart, curEnd int)
	dfs = func(curStart, curEnd int) {
		if curEnd == curStart+1 {
			query(curStart)
			return
		}

		mid := (curStart + curEnd) >> 1
		for i := curStart; i < mid; i++ {
			mutate(i)
		}
		dfs(mid, curEnd)
		for i := curStart; i < mid; i++ {
			undo()
		}

		for i := mid; i < curEnd; i++ {
			mutate(i)
		}
		dfs(curStart, mid)
		for i := mid; i < curEnd; i++ {
			undo()
		}
	}

	dfs(start, end)
}

// 238. 除自身以外数组的乘积
// https://leetcode.cn/problems/product-of-array-except-self/
func productExceptSelf(nums []int) []int {
	n := len(nums)
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = 1
	}

	cur := 1
	history := make([]int, 0, n)
	MutateWithoutOneUndo(
		0, n,
		func(index int) {
			history = append(history, cur)
			cur *= nums[index]
		},
		func() {
			cur = history[len(history)-1]
			history = history[:len(history)-1]
		},
		func(index int) {
			res[index] = cur
		},
	)

	return res
}
