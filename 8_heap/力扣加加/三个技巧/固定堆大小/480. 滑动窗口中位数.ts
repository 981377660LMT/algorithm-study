// 维护一个有序数组即可

import { bisectLeft } from '../../../../9_排序和搜索/二分/7_二分搜索寻找最左插入位置'
import { bisectInsort } from '../../../../9_排序和搜索/二分/7_二分搜索插入元素'

// 用二分方法维护
function medianSlidingWindow(nums: number[], k: number): number[] {
  const sortedList: number[] = []
  const res: number[] = []
  let l = 0

  for (let r = 0; r < nums.length; r++) {
    bisectInsort(sortedList, nums[r])

    while (sortedList.length > k) {
      sortedList.splice(bisectLeft(sortedList, nums[l]), 1)
      l++
    }

    if (sortedList.length === k) {
      res.push((sortedList[~~(k / 2)] + sortedList[~~((k - 1) / 2)]) / 2)
    }
  }

  return res
}

console.log(medianSlidingWindow([1, 3, -1, -3, 5, 3, 6, 7], 3))
