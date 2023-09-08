// 1755. 最接近目标值的子序列和
// https://leetcode.cn/problems/closest-subsequence-sum/description/
// 给你一个整数数组 nums 和一个目标值 goal 。
// 你需要从 nums 中选出一个子序列，使子序列元素总和最接近 goal
// 1 <= nums.length <= 40
// -107 <= nums[i] <= 107
// -1e9 <= goal <= 1e9

import { subsetSumSorted } from '../subsetSum/subsetSum'
import { twoSum } from '../twoSum'

function minAbsDifference(nums: number[], goal: number): number {
  const mid = nums.length >>> 1
  const leftSum = subsetSumSorted(nums.slice(0, mid))
  const rightSum = subsetSumSorted(nums.slice(mid))
  return twoSum(leftSum, rightSum, goal)
}

if (require.main === module) {
  console.log(minAbsDifference([5, -7, 3, 5], 6))
}
