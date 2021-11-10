/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * 我们最多可以将 K 个值从 0 变成 1 。
 * 返回仅包含 1 的最长（连续）数组的长度。
 */
const longestOnes = function (nums: number[], k: number): number {
  let l = 0
  let res = 0

  for (let r = 0; r < nums.length; r++) {
    if (nums[r] === 0) k--
    while (k < 0) {
      if (nums[l] === 0) k++
      l++
    }

    res = Math.max(res, r - l + 1)
  }

  return res
}
