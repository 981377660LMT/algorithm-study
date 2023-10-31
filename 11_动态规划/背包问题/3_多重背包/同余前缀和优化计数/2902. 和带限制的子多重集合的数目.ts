// 2902. 和带限制的子多重集合的数目
// https://leetcode.cn/problems/count-of-sub-multisets-with-bounded-sum/description/
// 给你一个下标从 0 开始的非负整数数组 nums 和两个整数 l 和 r 。
// 请你返回 nums 中子多重集合的和在闭区间 [l, r] 之间的 子多重集合的数目 。
// 由于答案可能很大，请你将答案对 109 + 7 取余后返回。
// 1 <= nums.length <= 2e4
// 0 <= nums[i] <= 2e4
// nums 的和不超过 2e4 。

import { boundedKnapsackDpCountWays } from '../BoundedKnapsack'

const MOD = 1e9 + 7
function countSubMultisets(nums: number[], left: number, right: number): number {
  const counter = new Map<number, number>()
  nums.forEach(n => counter.set(n, (counter.get(n) || 0) + 1))

  const values = [...counter.keys()]
  const counts = [...counter.values()]
  const dp = boundedKnapsackDpCountWays(values, counts, right)
  let res = 0
  for (let i = left; i <= Math.min(right, dp.length - 1); i++) {
    res = (res + dp[i]) % MOD
  }
  return res
}

if (require.main === module) {
  // nums = [1,2,2,3], l = 6, r = 6
  console.log(countSubMultisets([1, 2, 2, 3], 6, 6))
}
