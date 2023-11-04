// https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/

package main

type S = struct{ value int }

// 线段树分治的特殊情形.
// 调用 `query` 时，`state` 为对除了 `index` 以外所有点均调用过了 `mutate` 的状态。但不保证调用 `mutate` 的顺序。
// 总计会调用 $O(NlgN)$ 次的 `mutate` 和 `query`, 以及 $O(N)$ 次的 `copy`.
// !将一个不可撤销的数据结构以`修改时拷贝`的方式变成`可撤销`的.
func MutateWithoutOne(
	initState *S,
	start, end int,
	/** 这里的 index 也就是 time. */
	mutate func(state *S, index int),
	/** 通过拷贝实现撤销接口. */
	copy func(state *S) *S,
	query func(state *S, index int),
) {
	var dfs func(state *S, curStart, curEnd int)
	dfs = func(state *S, curStart, curEnd int) {
		if curEnd == curStart+1 {
			query(state, curStart)
			return
		}

		mid := (curStart + curEnd) >> 1
		leftCopy := copy(state)
		for i := curStart; i < mid; i++ {
			mutate(leftCopy, i)
		}
		dfs(leftCopy, mid, curEnd)

		rightCopy := copy(state)
		for i := mid; i < curEnd; i++ {
			mutate(rightCopy, i)
		}
		dfs(rightCopy, curStart, mid)
	}

	dfs(initState, start, end)
}

// 238. 除自身以外数组的乘积
// https://leetcode.cn/problems/product-of-array-except-self/
func productExceptSelf(nums []int) []int {
	res := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		res[i] = 1
	}

	MutateWithoutOne(
		&S{value: 1},
		0, len(nums),
		func(state *S, index int) {
			state.value *= nums[index]
		},
		func(state *S) *S {
			return &S{value: state.value}
		},
		func(state *S, index int) {
			res[index] = state.value
		},
	)

	return res
}
