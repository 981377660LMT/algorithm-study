import { kanade } from './53. 最大子数组和'

/**
 * @description
 * 最大的环形子数组和
 */
const maxSubarraySumCircular = function (nums: number[]): number {
  // 最大的环形子数组和 = max(最大子数组和，数组总和-最小子数组和)
  const sum = nums.reduce((pre, cur) => pre + cur, 0)
  const maxInRawArray = kanade(nums, true)
  // 有环：最大和子数组包含首尾元素
  const maxInCircularArray = sum - kanade(nums.slice(1, -1), false)
  return Math.max(maxInRawArray, maxInCircularArray)
}

export { maxSubarraySumCircular }
