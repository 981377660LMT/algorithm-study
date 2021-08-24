// 如果数组是非负的，那么前缀和就是一个单调不递减数组，我们有时候可以基于它来做二分。
// 数组 A 有多少个连续的子数组，其元素只和在 [lower, upper]的范围内。
// 1 <= nums.length <= 10**5
// 思路：遍历数组中的每一个元素，求得到当前元素的前缀和

import { bisectLeft } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'
import { bisectRight } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最插右入位置'
import { bisectInsort } from '../../../9_排序和搜索/二分api/7_二分搜索插入元素'

// lower<=Si-Sj<=upper 等价于 Si-upper<=Sj<=Si-lower
function countRangeSum(nums: number[], lower: number, upper: number): number {
  const pre = [0]
  let res = 0
  let sum = 0

  for (const num of nums) {
    sum += num
    res += bisectRight(pre, sum - lower) - bisectLeft(pre, sum - upper)
    bisectInsort(pre, sum)
  }

  return res
}

console.log(countRangeSum([-2, 5, -1], -2, 2))

export {}
