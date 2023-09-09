package main

type S = struct{ value int }

// 线段树分治的特殊情形.
func MutateWithoutOne(
	state S,
	start, end int,
	mutate func(state *S, index int),
	query func(state *S, index int),
) {
	var dfs func(state *S, curStart, curEnd int)
	dfs = func(state *S, curStart, curEnd int) {
		if curEnd == curStart+1 {
			query(state, curStart)
			return
		}

		mid := (curStart + curEnd) >> 1

		leftCopy := &S{}
		*leftCopy = *state
		for i := curStart; i < mid; i++ {
			mutate(leftCopy, i)
		}
		dfs(leftCopy, mid, curEnd)

		rightCopy := &S{}
		*rightCopy = *state
		for i := mid; i < curEnd; i++ {
			mutate(rightCopy, i)
		}
		dfs(rightCopy, curStart, mid)
	}

	dfs(&state, start, end)
}

func productExceptSelf(nums []int) []int {
	res := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		res[i] = 1
	}

	MutateWithoutOne(
		S{value: 1},
		0, len(nums),
		func(state *S, index int) {
			state.value *= nums[index]
		},
		func(state *S, index int) {
			res[index] = state.value
		},
	)

	return res
}
