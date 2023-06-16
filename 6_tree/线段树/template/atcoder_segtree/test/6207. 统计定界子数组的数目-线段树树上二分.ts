// !6207. 统计定界子数组的数目-线段树树上二分解法
// https://leetcode.cn/problems/count-subarrays-with-fixed-bounds

// 给你一个整数数组 nums 和两个整数 minK 以及 maxK 。
// nums 的定界子数组是满足下述条件的一个子数组：
// 子数组中的 最小值 等于 minK 。
// 子数组中的 最大值 等于 maxK 。
// 返回定界子数组的数目。
// 子数组是数组中的一个连续部分。

// !固定子数组的一端，则子数组的最小（最大）值关于另一端点具有单调性，
// !因此可以使用二分查找、滑动窗口来求出使得最小（最大值）值落在某一范围内的区间

import { SegmentTreePointUpdateRangeQuery } from '../SegmentTreePointUpdateRangeQuery'

const INF = 2e15
function countSubarrays(nums: number[], minK: number, maxK: number): number {
  const n = nums.length
  const minTree = new SegmentTreePointUpdateRangeQuery(
    nums,
    () => INF,
    (a, b) => Math.min(a, b)
  )

  const maxTree = new SegmentTreePointUpdateRangeQuery(
    nums,
    () => 0,
    (a, b) => Math.max(a, b)
  )

  let res = 0
  for (let left = 0; left < n; left++) {
    let max1 = maxTree.maxRight(left, x => x < maxK)
    let max2 = maxTree.maxRight(left, x => x <= maxK)
    let min1 = minTree.maxRight(left, x => x > minK)
    let min2 = minTree.maxRight(min1, x => x >= minK)
    res += Math.max(0, Math.min(max2, min2) - Math.max(min1, max1))
  }
  return res
}

console.log(countSubarrays([1, 3, 5, 2, 7, 5], 1, 5)) // 2
console.log(countSubarrays([1, 1, 1, 1], 1, 1)) // 10
