// 滚动hash+二分答案

import { BigIntHasher } from '../../BigIntHasher'

/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 */
function findLength(nums1: number[], nums2: number[]): number {
  const hasher1 = new BigIntHasher(nums1.map(String))
  const hasher2 = new BigIntHasher(nums2.map(String))
  let left = 0
  let right = Math.min(nums1.length, nums2.length)

  while (left <= right) {
    const mid = (left + right) >> 1
    if (isExist(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  function isExist(len: number): boolean {
    if (len === 0) return true
    const visited = new Set<bigint>()

    for (let left = 1; left + len - 1 <= nums1.length; left++) {
      const hash = hasher1.getHashOfRange(left, left + len - 1)
      visited.add(hash)
    }

    for (let left = 1; left + len - 1 <= nums2.length; left++) {
      const hash = hasher2.getHashOfRange(left, left + len - 1)
      if (visited.has(hash)) return true
    }

    return false
  }
}

console.log(findLength([1, 2, 3, 2, 1], [3, 2, 1, 4, 7]))
console.log(findLength([70, 39, 25, 40, 7], [52, 20, 67, 5, 31])) // 0
console.log(findLength([1, 0, 1, 0, 0, 0, 0, 0, 1, 1], [1, 1, 0, 1, 1, 0, 0, 0, 0, 0])) // 6
console.log(findLength([0, 1, 1, 0, 1, 1, 1, 0, 1, 0], [1, 0, 0, 0, 1, 0, 0, 1, 1, 0])) // 4
export {}
