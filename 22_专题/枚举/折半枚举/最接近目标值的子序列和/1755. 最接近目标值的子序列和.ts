// 给你一个整数数组 nums 和一个目标值 goal 。
// 你需要从 nums 中选出一个子序列，使子序列元素总和最接近 goal
// 1 <= nums.length <= 40
// -107 <= nums[i] <= 107
// -1e9 <= goal <= 1e9

import { subsetSum } from '../getSubsetSum'
import { twoSum } from '../twoSum'

function minAbsDifference(nums: number[], goal: number): number {
  const mid = nums.length >>> 1
  const left = subsetSum(nums.slice(0, mid)).sort((a, b) => a - b)
  const right = subsetSum(nums.slice(mid)).sort((a, b) => a - b)
  return twoSum(left, right, goal)
}

if (require.main === module) {
  console.log(minAbsDifference([5, -7, 3, 5], 6))
}
