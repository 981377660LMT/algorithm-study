package main

import "sort"

const INF int32 = 1e9 + 10

// def min2(a: int, b: int) -> int:
//     return a if a < b else b

// class Solution:
//     def findPermutation(self, nums: List[int]) -> List[int]:
//         @lru_cache(None)
//         def dfs(index: int, visited: int, pre: int, first: int) -> int:
//             if index == n:
//                 return abs(pre - nums[first])
//             nonlocal next_
//             resCost = INF
//             for next in range(n):
//                 if visited & (1 << next):
//                     continue
//                 nextCost = dfs(index + 1, visited | (1 << next), next, first) + abs(
//                     pre - nums[next]
//                 )
//                 if nextCost < resCost:
//                     resCost = nextCost
//                     next_[(index, visited, pre)] = (
//                         index + 1,
//                         visited | (1 << next),
//                         next,
//                     )
//             return resCost

// resCost = INF
// res = [INF]
// n = len(nums)
// for i in range(n):
//
//	next_ = dict()
//	tmp = dfs(1, 1 << i, i, i)
//	# print(tmp, resCost, first)
//	if tmp < resCost:
//	    resCost = tmp
//	    curRes = [i]
//	    curState = (1, 1 << i, i)
//	    for _ in range(1, n):
//	        curState = next_[curState]
//	        curRes.append(curState[2])
//	    res = curRes
//
// dfs.cache_clear()
// return res

// func main() {
// 	// nums = [1,0,2]
// 	fmt.Println(findPermutation([]int{1, 0, 2}))
// }

// 1. index 维度可以不用，优化模板
// 2. 可以只计算0开头

func findPermutation(nums []int) []int {
	if sort.IntsAreSorted(nums) {
		return append([]int(nil), nums...)
	}

	n := int32(len(nums))
	newNums := make([]uint8, n)
	for i := range nums {
		newNums[i] = uint8(nums[i])
	}

	resCost, res := uint8(255), []int32{}
	for i := int32(0); i < n; i++ {
		memo := make([]uint8, n*(1<<n)*n)
		for i := range memo {
			memo[i] = 255
		}
		transfer := make([]int32, n*(1<<n)*n)
		hash := func(index, visited, pre int32) int32 {
			return index*(1<<n)*n + visited*n + pre
		}

		first := i
		var dfs func(index, visited, pre int32) uint8
		dfs = func(index, visited, pre int32) uint8 {
			if index == n {
				return abs(uint8(pre) - newNums[first])
			}
			hash_ := hash(index, visited, pre)
			if memo[hash_] != 255 {
				return memo[hash_]
			}

			resCost := uint8(255)
			for next := int32(0); next < n; next++ {
				if visited&(1<<next) > 0 {
					continue
				}
				nextCost := dfs(index+1, visited|(1<<next), next) + abs(uint8(pre)-newNums[next])
				if nextCost < resCost {
					resCost = nextCost
					curHash := hash(index, visited, pre)
					nextHash := hash(index+1, visited|(1<<next), next)
					transfer[curHash] = nextHash
				}
			}
			memo[hash_] = resCost
			return resCost
		}
		tmp := dfs(1, 1<<i, i)

		if tmp < resCost {
			resCost = tmp
			curRes := []int32{i}
			curState := hash(1, 1<<i, i)
			for i := int32(1); i < n; i++ {
				curState = transfer[curState]
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

func abs(a uint8) uint8 {
	if a < 0 {
		return -a
	}
	return a
}
