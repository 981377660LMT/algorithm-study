/* eslint-disable no-inner-declarations */

/**
 * 给定二分答案的区间 [{@link left}, {@link right}], 求第 {@link kth} 小的答案.
 * @param countNgt 答案不超过mid时, 满足条件的个数.
 * @param kth 从0开始.
 */
function getKth0(
  left: number,
  right: number,
  countNgt: (mid: number) => number,
  kth: number
): number {
  while (left <= right) {
    const mid = left + Math.floor((right - left) / 2)
    if (countNgt(mid) <= kth) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }
  return right + 1
}

/**
 * 给定二分答案的区间 [{@link left}, {@link right}], 求第 {@link kth} 小的答案.
 * @param countNgt 答案不超过mid时, 满足条件的个数.
 * @param kth 从1开始.
 */
function getKth1(
  left: number,
  right: number,
  countNgt: (mid: number) => number,
  kth: number
): number {
  while (left <= right) {
    const mid = left + Math.floor((right - left) / 2)
    if (countNgt(mid) < kth) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }
  return left
}

const EPS = 1e-12

/**
 * 给定二分答案的区间 [{@link left}, {@link right}], 求第 {@link kth} 小的答案.
 * @param countNgt 答案不超过mid时, 满足条件的个数.
 * @param kth 从0开始.
 */
function getKth0Float64(
  left: number,
  right: number,
  countNgt: (mid: number) => number,
  kth: number
): number {
  while (left <= right) {
    const mid = left + (right - left) / 2
    if (countNgt(mid) <= kth) {
      left = mid + EPS
    } else {
      right = mid - EPS
    }
  }
  return right + EPS
}

/**
 * 给定二分答案的区间 [{@link left}, {@link right}], 求第 {@link kth} 小的答案.
 * @param countNgt 答案不超过mid时, 满足条件的个数.
 * @param kth 从1开始.
 */
function getKth1Float64(
  left: number,
  right: number,
  countNgt: (mid: number) => number,
  kth: number
): number {
  while (left <= right) {
    const mid = left + (right - left) / 2
    if (countNgt(mid) < kth) {
      left = mid + EPS
    } else {
      right = mid - EPS
    }
  }
  return left
}

export { getKth0, getKth1, getKth0Float64, getKth1Float64 }

if (require.main === module) {
  // 719. 找出第 K 小的数对距离
  // https://leetcode.cn/problems/find-k-th-smallest-pair-distance/
  function smallestDistancePair(nums: number[], k: number): number {
    const countNgt = (mid: number): number => {
      let res = 0
      let left = 0
      for (let right = 0; right < nums.length; right++) {
        while (left <= right && nums[right] - nums[left] > mid) {
          left++
        }
        res += right - left
      }
      return res
    }

    nums = nums.slice().sort((a, b) => a - b)
    const res1 = getKth1(0, nums[nums.length - 1] - nums[0], countNgt, k)
    const res0 = getKth0(0, nums[nums.length - 1] - nums[0], countNgt, k - 1)
    if (res0 !== res1) {
      throw new Error(`res0=${res0}, res1=${res1}`)
    }
    return res0
  }
}
