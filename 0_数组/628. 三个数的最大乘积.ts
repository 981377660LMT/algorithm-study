// 628. 三个数的最大乘积

/**
 * @param {number[]} nums
 * @return {number}
 * 1、全正 => 最后三个数相乘
 * 2、如果有负数最大的乘机要么是最后三个数相乘，
 * 要么是两个最小的负数相乘再乘以最大的正数
 */
function maximumProduct(nums: number[]): number {
  const n = nums.length
  nums.sort((a, b) => a - b)
  return Math.max(nums[n - 3] * nums[n - 2] * nums[n - 1], nums[0] * nums[1] * nums[n - 1])
}

export { maximumProduct }
