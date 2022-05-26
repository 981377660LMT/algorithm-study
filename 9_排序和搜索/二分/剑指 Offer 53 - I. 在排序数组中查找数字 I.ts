import { bisectLeft, bisectRight } from './bisect'

// 统计一个数字在排序数组中出现的次数。
function search(nums: number[], target: number): number {
  const l = bisectLeft(nums, target)
  const r = bisectRight(nums, target)
  return r - l
}
