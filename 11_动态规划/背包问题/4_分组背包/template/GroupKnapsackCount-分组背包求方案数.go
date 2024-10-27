package main

const MOD int = 1e9 + 7

// 分组背包求方案数.
// 每个物品有count[i]个, 求选择k个物品的方案数.
func GroupKnapsackCount(count []int, k int, mod int) []int {
	n := len(count)
	dp := make([]int, k+1)
	dp[0] = 1
	mx := 0
	for i := 0; i < n; i++ {
		mx += count[i]
		if mx > k {
			mx = k
		}
		for j := 1; j <= mx; j++ {
			dp[j] = (dp[j] + dp[j-1]) % mod
		}
		for j := mx; j >= count[i]; j-- {
			dp[j] = (dp[j] - dp[j-count[i]]) % mod
		}
	}
	for i := 0; i <= k; i++ {
		if dp[i] < 0 {
			dp[i] += mod
		}
	}
	return dp
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// https://leetcode.cn/problems/number-of-dice-rolls-with-target-sum/
func numRollsToTarget(n int, k int, target int) int {
	count := make([]int, n)
	for i := 0; i < n; i++ {
		count[i] = k
	}
	target -= n
	dp := GroupKnapsackCount(count, target, MOD)
	return dp[target]
}

// https://leetcode.cn/problems/find-the-original-typed-string-ii/description/
func possibleStringCount(word string, k int) int {
	var lens []int
	EnumerateGroup([]byte(word), func(_ []byte, start, end int) {
		lens = append(lens, end-start)
	})

	all_ := 1
	for _, l := range lens {
		all_ *= l
		all_ %= MOD
	}

	calcBad := func() int {
		if len(lens) >= k {
			return 0
		}
		target := k - len(lens) - 1
		dp := GroupKnapsackCount(lens, target, MOD)
		sum_ := 0
		for _, v := range dp {
			sum_ += v
			sum_ %= MOD
		}
		return sum_
	}

	bad := calcBad()
	return (all_ - bad + MOD) % MOD
}

// 遍历连续相同元素的分组.相当于python中的`itertools.groupby`.
func EnumerateGroup[T comparable](arr []T, f func(group []T, start, end int)) {
	ptr := 0
	n := len(arr)
	for ptr < n {
		leader := arr[ptr]
		group := []T{leader}
		start := ptr
		ptr++
		for ptr < n && arr[ptr] == leader {
			group = append(group, arr[ptr])
			ptr++
		}
		f(group, start, ptr)
	}
}
