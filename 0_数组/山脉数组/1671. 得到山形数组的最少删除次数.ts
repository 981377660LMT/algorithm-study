// 问题转化为 左右上升子序列的最长长度
// 3 <= nums.length <= 1000

import { lengthOfLIS } from '../../11_动态规划/lis最长上升子序列问题/1_最长上升子序列dp加二分解法'

// 枚举山顶 O(n^2logn)
function minimumMountainRemovals(nums: number[]): number {
  let res = nums.length

  for (let peakIndex = 1; peakIndex < nums.length - 1; peakIndex++) {
    const left = lengthOfLIS(nums.slice(0, peakIndex + 1))
    const right = lengthOfLIS(nums.slice(peakIndex).reverse())
    if (left === 1 || right === 1) continue // 没有
    res = Math.min(res, nums.length + 1 - left - right)
    if (res === 0) return 0
  }

  return res
}

// console.log(minimumMountainRemovals([1, 3, 1])) // 0
// console.log(minimumMountainRemovals([2, 1, 1, 5, 6, 2, 3, 1])) // 3
console.log(minimumMountainRemovals([100, 92, 89, 77, 74, 66, 64, 66, 64])) // 6
