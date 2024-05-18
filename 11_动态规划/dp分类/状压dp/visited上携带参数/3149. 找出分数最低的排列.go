// 给你一个数组 nums ，它是 [0, 1, 2, ..., n - 1] 的一个排列。
// 对于任意一个 [0, 1, 2, ..., n - 1] 的排列 perm ，其 分数 定义为：
// score(perm) = |perm[0] - nums[perm[1]]| + |perm[1] - nums[perm[2]]| + ... + |perm[n - 1] - nums[perm[0]]|
// 返回具有 最低 分数的排列 perm 。如果存在多个满足题意且分数相等的排列，则返回其中字典序最小的一个。
//
// 字典序最小：按照字典序从小到大搜索.
// 求方案：记录状态转移的前驱或者后继.

package main

import (
	"math/bits"
	"sort"
)

const INF32 int32 = 1e9 + 10

func findPermutation(nums []int) []int {
	if sort.IntsAreSorted(nums) {
		return nums
	}

	n := int32(len(nums))
	newNums := make([]int32, n)
	for i := range nums {
		newNums[i] = int32(nums[i])
	}

	resCost, res := INF32, []int32{INF32}
	memo := make([]int32, (1<<n)*n)
	next_ := make([]int32, (1<<n)*n)
	for i := int32(0); i < n; i++ {
		for i := range memo {
			memo[i] = -1
		}
		hash := func(visited, pre int32) int32 {
			return visited*n + pre
		}

		first := i
		var dfs func(visited, pre int32) int32
		dfs = func(visited, pre int32) int32 {
			index := int32(bits.OnesCount32(uint32(visited)))
			if index == n {
				return abs(pre - newNums[first])
			}
			hash_ := hash(visited, pre)
			if memo[hash_] != -1 {
				return memo[hash_]
			}

			resCost := INF32
			for next := int32(0); next < n; next++ {
				if visited&(1<<next) > 0 {
					continue
				}
				nextCost := dfs(visited|(1<<next), next) + abs(pre-newNums[next])
				if nextCost < resCost {
					resCost = nextCost
					curHash := hash(visited, pre)
					nextHash := hash(visited|(1<<next), next)
					next_[curHash] = nextHash
				}
			}
			memo[hash_] = resCost
			return resCost
		}
		tmp := dfs(1<<i, i)

		if tmp < resCost {
			resCost = tmp
			curRes := []int32{i}
			curState := hash(1<<i, i)
			for i := int32(1); i < n; i++ {
				curState = next_[curState]
				curRes = append(curRes, curState%n)
			}
			res = curRes
		}
	}

	newRes := make([]int, n)
	for i := range res {
		newRes[i] = int(res[i])
	}
	return newRes
}

func abs(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}
