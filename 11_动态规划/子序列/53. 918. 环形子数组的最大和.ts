import { maxSubArray } from './53. 最大子序和'

/**
 * @param {number[]} nums
 * @return {number}
 * 53. 最大子序和
 */
const maxSubarraySumCircular = function (nums: number[]): number {
  // 最大值max在原数组中
  // 最大值max在环形数组中
  const sum = nums.reduce((pre, cur) => pre + cur, 0)
  const kanade = (nums: number[], getMax: boolean) => {
    if (nums.length === 0) return 0
    if (nums.length === 1) return nums[0]

    let res = getMax ? -Infinity : Infinity
    let sum = 0
    for (const num of nums) {
      if (getMax) {
        sum = Math.max(sum + num, num)
        res = Math.max(res, sum)
      } else {
        sum = Math.min(sum + num, num)
        res = Math.min(res, sum)
      }
    }

    return res
  }

  const maxInRawArray = kanade(nums, true)
  const maxInCircularArray = sum - kanade(nums.slice(1, -1), false)

  return Math.max(maxInRawArray, maxInCircularArray)
}
