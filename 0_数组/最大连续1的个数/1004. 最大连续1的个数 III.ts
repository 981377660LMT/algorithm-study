/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * 我们最多可以将 K 个值从 0 变成 1 。
 * 返回仅包含 1 的最长（连续）数组的长度。
 */
var longestOnes = function (nums: number[], k: number): number {
  let l = 0
  let r = 0
  let res = 0
  while (r < nums.length) {
    if (nums[r] === 0) k--
    while (k < 0) {
      l++
      if (nums[l - 1] === 0) k++
    }
    res = Math.max(res, r - l + 1)
    r++
  }
  return res
}
