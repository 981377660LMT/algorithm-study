// package main

// import "sort"

// const INF int32 = 1e9 + 10

// func findPermutation(nums []int) []int {
// 	if sort.IntsAreSorted(nums) {
// 		return append([]int(nil), nums...)
// 	}
// 	n := int32(len(nums))
// 	newNums := make([]int32, n)
// 	for i := range nums {
// 		newNums[i] = int32(nums[i])
// 	}

// 	resCost, res := INF, []int32{INF}
// 	for i := int32(0); i < n; i++ {
// 		memo := make([]int32, n*(1<<n)*n)
// 		for i := range memo {
// 			memo[i] = -1
// 		}
// 		next_ := make([]int32, n*(1<<n)*n)
// 		hash := func(index, visited, pre int32) int32 {
// 			return index*(1<<n)*n + visited*n + pre
// 		}

// 		first := i
// 		var dfs func(index, visited, pre int32) int32
// 		dfs = func(index, visited, pre int32) int32 {
// 			if index == n {
// 				return abs(pre - newNums[first])
// 			}
// 			hash_ := hash(index, visited, pre)
// 			if memo[hash_] != -1 {
// 				return memo[hash_]
// 			}

// 			resCost := INF
// 			for next := int32(0); next < n; next++ {
// 				if visited&(1<<next) > 0 {
// 					continue
// 				}
// 				nextCost := dfs(index+1, visited|(1<<next), next) + abs(pre-newNums[next])
// 				if nextCost < resCost {
// 					resCost = nextCost
// 					curHash := hash(index, visited, pre)
// 					nextHash := hash(index+1, visited|(1<<next), next)
// 					next_[curHash] = nextHash
// 				}
// 			}
// 			memo[hash_] = resCost
// 			return resCost
// 		}
// 		tmp := dfs(1, 1<<i, i)

// 		if tmp < resCost {
// 			resCost = tmp
// 			curRes := []int32{i}
// 			curState := hash(1, 1<<i, i)
// 			for i := int32(1); i < n; i++ {
// 				curState = next_[curState]
// 				curRes = append(curRes, curState%n)
// 			}
// 			res = curRes
// 		}
// 	}

// 	newRes := make([]int, n)
// 	for i := range res {
// 		newRes[i] = int(res[i])
// 	}
// 	return newRes
// }

// func abs(a int32) int32 {
// 	if a < 0 {
// 		return -a
// 	}
// 	return a
// }

impl Solution {
    pub fn find_permutation(nums: Vec<i32>) -> Vec<i32> {}
}
