import { bisectLeft } from './7_二分搜索寻找最左插入位置'
import { bisectRight } from './7_二分搜索寻找最插右入位置'

// 统计一个数字在排序数组中出现的次数。
function search(nums: number[], target: number): number {
  const l = bisectLeft(nums, target)
  const r = bisectRight(nums, target)
  return r - l
}
