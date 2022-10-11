import { bisectLeft, bisectRight } from './bisect'

/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number[]}
 */
let searchRange = (nums: number[], target: number): number[] => {
  const pos1 = bisectLeft(nums, target)
  const pos2 = bisectRight(nums, target) - 1
  return pos1 <= pos2 ? [pos1, pos2] : [-1, -1]
}

console.log(searchRange([5, 7, 7, 8, 8, 10], 8))
console.log(searchRange([2, 2], 2))
console.log(searchRange([1, 2, 3], 4))

export default 1
