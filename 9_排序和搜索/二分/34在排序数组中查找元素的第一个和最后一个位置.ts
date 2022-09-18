import { bisectLeft, bisectRight } from './bisect'

/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number[]}
 */
let searchRange = (nums: number[], target: number): number[] => {
  const l = bisectLeft(nums, target)
  const r = bisectRight(nums, target) - 1
  return l <= r ? [l, r] : [-1, -1]
}

console.log(searchRange([5, 7, 7, 8, 8, 10], 8))
console.log(searchRange([2, 2], 2))
console.log(searchRange([1, 2, 3], 4))

export default 1
