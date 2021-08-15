import { betweenK } from '../api'

/**
 * @param {number[]} nums
 * @param {number} left
 * @param {number} right
 * @return {number}
 * 求连续、非空且其中最大元素满足大于等于L 小于等于R的子数组个数。
 */
const numSubarrayBoundedMax = function (nums: number[], left: number, right: number): number {
  return betweenK(left, right, nums)
}

console.log(numSubarrayBoundedMax([2, 1, 4, 3], 2, 3))
