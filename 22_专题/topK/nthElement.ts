/* eslint-disable no-shadow */
/* eslint-disable semi-style */
/* eslint-disable no-param-reassign */

import assert from 'assert'

interface Options<E> {
  lower?: number
  upper?: number
  key?: (e: E) => number
}

/**
 * 快速选择算法找到数组中第`nth`小的数 (`nth`从0开始)
 *
 * @complexity O(n)
 */
function nthElement<E>(nums: E[], nth: number, options?: Options<E>): E {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const { lower = 0, upper = nums.length - 1, key = (e: any) => e } = options || {}
  return quickSelect(nums, nth, lower, upper, key)
}

function quickSelect<E>(
  nums: E[],
  nth: number,
  lower: number,
  upper: number,
  key: (e: E) => number
): E {
  const pos = partition(nums, lower, upper, key)
  if (pos === nth) {
    return nums[pos]
  }

  if (pos < nth) {
    return quickSelect(nums, nth, pos + 1, upper, key)
  }

  return quickSelect(nums, nth, lower, pos - 1, key)
}

function partition<E>(nums: E[], left: number, right: number, key: (e: E) => number): number {
  const randIndex = randInt(left, right)
  ;[nums[left], nums[randIndex]] = [nums[randIndex], nums[left]]
  let pivotIndex = left
  const pivot = nums[left]
  for (let i = left + 1; i <= right; i++) {
    if (key(nums[i]) < key(pivot)) {
      pivotIndex++
      ;[nums[i], nums[pivotIndex]] = [nums[pivotIndex], nums[i]]
    }
  }

  ;[nums[left], nums[pivotIndex]] = [nums[pivotIndex], nums[left]]
  return pivotIndex
}

function randInt(left: number, right: number): number {
  return Math.floor(Math.random() * (right - left + 1)) + left
}

if (require.main === module) {
  const nums = Array.from({ length: 100 }, () => Math.floor(Math.random() * 100))
  nums.sort((a, b) => a - b)
  for (let _ = 0; _ < 1000; _++) {
    const nth = randInt(0, nums.length - 1)
    assert.strictEqual(nthElement(nums, nth), nums[nth])
  }
}

export { nthElement }
