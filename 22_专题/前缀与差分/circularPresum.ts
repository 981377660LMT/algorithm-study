/* eslint-disable no-inner-declarations */

import { sortSearch } from '../../9_排序和搜索/二分/sortSearch'

/**
 * 环形数组前缀和.
 * @param arr 环形数组的循环部分.
 * @returns 返回区间 `[start, end)` 的和.
 */
function cicularPreSum(arr: ArrayLike<number>): (start: number, end: number) => number {
  const n = arr.length
  const preSum = Array(n + 1).fill(0)
  for (let i = 0; i < n; i++) preSum[i + 1] = preSum[i] + arr[i]
  const cal = (r: number) => preSum[n] * Math.floor(r / n) + preSum[r % n]
  const query = (start: number, end: number) => {
    if (start >= end) return 0
    return cal(end) - cal(start)
  }
  return query
}

export { cicularPreSum }

if (require.main === module) {
  const query = cicularPreSum([1, 2, 3, 4, 5])
  console.log(query(20, 11))

  // 100076. 无限数组的最短子数组
  // https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
  // 求循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
  // 1 <= nums.length <= 1e5
  // 1 <= nums[i] <= 1e5
  // 1 <= target <= 1e9

  function minSizeSubarray(nums: number[], target: number): number {
    let res = Infinity
    const sum = cicularPreSum(nums)
    for (let start = 0; start < nums.length; start++) {
      const cand = sortSearch(0, 1e9 + 10, mid => sum(start, start + mid) >= target)
      if (sum(start, start + cand) === target) {
        res = Math.min(res, cand)
      }
    }

    return res === Infinity ? -1 : res
  }
}
