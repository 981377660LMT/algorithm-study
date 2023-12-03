/* eslint-disable no-inner-declarations */

import { sortSearch } from '../../9_排序和搜索/二分/sortSearch'

/**
 * 环形数组前缀和.
 * @param arr 环形数组的循环部分.
 * @returns 返回区间 `[start, end)` 的和.
 */
function circularPreSum(arr: ArrayLike<number>): (start: number, end: number) => number {
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

/**
 * 二维环形数组前缀和.
 * @param grid 二维环形数组的循环部分.
 * @returns 返回区间 `[row1, row2) * [col1, col2)` 的和.
 */
function circularPreSum2D(
  grid: ArrayLike<ArrayLike<number>>
): (row1: number, col1: number, row2: number, col2: number) => number {
  const n = grid.length
  const m = grid[0].length
  const preSum = Array((n + 1) * (m + 1)).fill(0)
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < m; j++) {
      preSum[(i + 1) * (m + 1) + j + 1] =
        preSum[i * (m + 1) + j + 1] +
        preSum[(i + 1) * (m + 1) + j] -
        preSum[i * (m + 1) + j] +
        grid[i][j]
    }
  }

  const cal = (r: number, c: number) => {
    const res1 = preSum[n * (m + 1) + m] * Math.floor(r / n) * Math.floor(c / m)
    const res2 = preSum[(r % n) * (m + 1) + m] * Math.floor(c / m)
    const res3 = preSum[n * (m + 1) + (c % m)] * Math.floor(r / n)
    const res4 = preSum[(r % n) * (m + 1) + (c % m)]
    return res1 + res2 + res3 + res4
  }

  const query = (row1: number, col1: number, row2: number, col2: number) => {
    if (row1 >= row2 || col1 >= col2) return 0
    const res1 = cal(row2, col2)
    const res2 = cal(row1, col2)
    const res3 = cal(row2, col1)
    const res4 = cal(row1, col1)
    return res1 - res2 - res3 + res4
  }

  return query
}

export { circularPreSum }

if (require.main === module) {
  const query = circularPreSum([1, 2, 3, 4, 5])
  console.log(query(20, 11))

  // 100076. 无限数组的最短子数组
  // https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
  // 求循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
  // 1 <= nums.length <= 1e5
  // 1 <= nums[i] <= 1e5
  // 1 <= target <= 1e9

  function minSizeSubarray(nums: number[], target: number): number {
    let res = Infinity
    const sum = circularPreSum(nums)
    for (let start = 0; start < nums.length; start++) {
      const cand = sortSearch(0, 1e9 + 10, mid => sum(start, start + mid) >= target)
      if (sum(start, start + cand) === target) {
        res = Math.min(res, cand)
      }
    }

    return res === Infinity ? -1 : res
  }

  const S = circularPreSum2D([
    [1, 2, 3],
    [4, 5, 6]
  ])
  console.log(S(0, 0, 3, 3))
}
