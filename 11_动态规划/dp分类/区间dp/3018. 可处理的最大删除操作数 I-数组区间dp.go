// 3018. 可处理的最大删除操作数 I
// https://leetcode.cn/problems/maximum-number-of-removal-queries-that-can-be-processed-i/solutions/2619968/jian-ji-de-dong-tai-gui-hua-by-peaceful-u3fs0/// 给定一个下标 从 0 开始 的数组 nums 和一个下标 从 0 开始 的数组 queries。
// 你可以在开始时执行以下操作 最多一次：
// !用 nums 的子序列替换 nums。
// 我们以给定的顺序开始处理查询；对于每个查询，我们执行以下操作：
// 如果 nums 的第一个 和 最后一个元素 小于 queries[i]，则查询处理 结束。
// 否则，如果 nums 的第一个 或 最后一个元素 大于或等于 queries[i]，则选择它，并从 nums 中 删除 选定的元素。
// 返回通过以最佳方式执行该操作可以处理的 最大 查询数。
// n,q<=1000.
//
// 区间dp
// !dp[l][r]: 表示在数组 nums 中，当数据 [l,r] 还没有被删除时我们所能查询的最大 queries 数
// !枚举每个i没有被删除的情况.

package main

func maximumProcessableQueries(nums []int, queries []int) int {
	n, q := len(nums), len(queries)

	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
	}
	for left := 0; left < n; left++ {
		for right := n - 1; right >= left; right-- {
			cur := 0

			if left > 0 {
				cand := dp[left-1][right]
				if nums[left-1] >= queries[cand] {
					cand++
				}
				cur = max(cur, cand)
			}

			if right < n-1 {
				cand := dp[left][right+1]
				if nums[right+1] >= queries[cand] {
					cand++
				}
				cur = max(cur, cand)
			}
			dp[left][right] = cur

			// !fast break
			if cur == q {
				return q
			}
		}
	}

	res := 0
	for i := 0; i < n; i++ {
		cand := dp[i][i]
		if nums[i] >= queries[cand] {
			cand++
		}
		res = max(res, cand)
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
