import { bisectLeft } from './7_二分搜索寻找最左插入位置'
import { bisectRight } from './7_二分搜索寻找最插右入位置'

/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number[]}
 */
var searchRange = (nums: number[], target: number): number[] => {
  const l = bisectLeft(nums, target)
  const r = bisectRight(nums, target) - 1
  return l <= r ? [l, r] : [-1, -1]
}
console.log(searchRange([5, 7, 7, 8, 8, 10], 8))
console.log(searchRange([2, 2], 2))
console.log(searchRange([1, 2, 3], 4))

export default 1
