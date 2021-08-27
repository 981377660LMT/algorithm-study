/**
 * @param {number[]} nums
 * @return {number}
 * 1、如果全是正数就是最后三个数相乘
 * 2、如果有负数最大的乘机要么是最后三个数相乘，
 * 要么是两个最小的负数相乘再乘以最大的正数
 *
 */
var maximumProduct = function (nums) {
  const len = nums.length
  nums.sort((a, b) => a - b)
  return Math.max(nums[len - 3] * nums[len - 2] * nums[len - 1], nums[0] * nums[1] * nums[len - 1])
}
