/**
 * @param {number[]} nums
 * @return {number} 两个数的最大乘积
 */
function maximumProduct(nums: number[]): number {
  const n = nums.length
  nums.sort((a, b) => a - b)
  return Math.max(nums[n - 1] * nums[n - 2], nums[0] * nums[1])
}

export { maximumProduct }
