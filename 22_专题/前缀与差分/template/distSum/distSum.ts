/* eslint-disable no-inner-declarations */

import { getMedian } from '../../../../19_数学/中位数/getMedian'

/**
 * **有序数组**所有点到`x=k`的距离之和.
 */
function distSum(sortedNums: ArrayLike<number>): (k: number) => number {
  const bisectRight = (nums: ArrayLike<number>, target: number) => {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = left + ((right - left) >>> 1)
      if (nums[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  const preSum = Array(sortedNums.length + 1)
  preSum[0] = 0
  for (let i = 0; i < sortedNums.length; i++) {
    preSum[i + 1] = preSum[i] + sortedNums[i]
  }

  return (k: number): number => {
    const pos = bisectRight(sortedNums, k)
    const leftSum = k * pos - preSum[pos]
    const rightSum = preSum[preSum.length - 1] - preSum[pos] - k * (sortedNums.length - pos)
    return leftSum + rightSum
  }
}

/**
 * **有序数组**中`[start,end)`范围内所有点到`x=k`的距离之和.
 */
function distSumRange(
  sortedNums: ArrayLike<number>
): (k: number, start: number, end: number) => number {
  const bisectLeft = (nums: ArrayLike<number>, target: number) => {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = left + ((right - left) >>> 1)
      if (nums[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  const preSum = Array(sortedNums.length + 1)
  preSum[0] = 0
  for (let i = 0; i < sortedNums.length; i++) {
    preSum[i + 1] = preSum[i] + sortedNums[i]
  }

  return (k: number, start: number, end: number): number => {
    if (start < 0) start = 0
    if (end > sortedNums.length) end = sortedNums.length
    if (start >= end) return 0
    const pos = bisectLeft(sortedNums, k)
    if (pos <= start) {
      return preSum[end] - preSum[start] - k * (end - start)
    }
    if (pos >= end) {
      return k * (end - start) - (preSum[end] - preSum[start])
    }
    const leftSum = k * (pos - start) - (preSum[pos] - preSum[start])
    const rightSum = preSum[end] - preSum[pos] - k * (end - pos)
    return leftSum + rightSum
  }
}

/**
 * **有序数组**中所有点对两两距离之和.一共有`n*(n-1)//2`对点对.
 */
function distSumOfAllPairs(sortedNums: ArrayLike<number>, start = 0, end = sortedNums.length) {
  let res = 0
  let preSum = 0
  for (let i = start; i < end; i++) {
    res += sortedNums[i] * i - preSum
    preSum += sortedNums[i]
  }
  return res
}

export { distSum, distSumRange, distSumOfAllPairs }

if (require.main === module) {
  // 100123. 执行操作使频率分数最大
  function maxFrequencyScore(nums: number[], k: number): number {
    nums.sort((a, b) => a - b)
    const D = distSumRange(nums)
    let res = 0
    let left = 0
    for (let right = 0; right < nums.length; right++) {
      while (left <= right) {
        const median = getMedian(nums, left, right + 1)
        if (D(median, left, right + 1) <= k) {
          break
        }
        left++
      }
      res = Math.max(res, right - left + 1)
    }
    return res
  }
}
