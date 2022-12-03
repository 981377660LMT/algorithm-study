/* eslint-disable no-param-reassign */

/**
 * @param {number[]} nums
 * @return {number[]}
 * 集合 s 包含从 1 到 n 的整数
 * 集合 丢失了一个数字 并且 有一个数字重复 。
 */
function findErrorNums(nums: number[]): number[] {
  const res = []

  for (let i = 0; i < nums.length; i++) {
    const mapped = Math.abs(nums[i]) - 1
    if (nums[mapped] < 0) res.push(mapped + 1)
    else nums[mapped] *= -1
  }

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] > 0) res.push(i + 1)
  }

  return res
}

export {}
