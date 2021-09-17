import { fix } from './最长连续 1 模型'

/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @description 
 * 给定一个由若干 0 和 1 组成的数组 A，我们最多可以将 K 个值从 0 变成 1 。
   返回仅包含 1 的最长（连续）子数组的长度。
 */
const longestOnes = function (nums: number[], k: number): number {
  return fix(nums.join(''), '1', k)
}

console.log(longestOnes([0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1], 3))
export default 1
