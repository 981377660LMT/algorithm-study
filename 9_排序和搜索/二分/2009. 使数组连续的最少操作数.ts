import { bisectLeft } from './7_二分搜索寻找最左插入位置'
import { bisectRight } from './7_二分搜索寻找最插右入位置'

/**
 * @param {number[]} nums
 * @return {number}
 * 连续数组:排序
 */
const minOperations = function (nums: number[]): number {
  const n = nums.length
  const uniqueNums = [...new Set(nums)]
  const un = uniqueNums.length
  uniqueNums.sort((a, b) => a - b)
  let res = un

  // 去重排序后 对每个数 二分找到左右边界 判断要删除多少个
  for (let i = 0; i < uniqueNums.length; i++) {
    const value = uniqueNums[i]
    // 找出区间 [nums[i], nums[i] + n-1] 中的不同的整数数量 cur。这 cur 个数是不需要修改的
    // 而剩余的 n - cur 个数是需要修改的
    const rightShouldKeepNumIndex = bisectRight(uniqueNums, value + n - 1)
    const leftShouldKeepNumIndex = bisectLeft(uniqueNums, value - (n - 1))
    res = Math.min(
      res,
      un - (rightShouldKeepNumIndex - 1 - i + 1), // rightShouldKeepNumCount - 1才是前一个可以保留的数的下标
      un - (i - leftShouldKeepNumIndex + 1)
    )
  }

  return res + (n - un) // 相同元素必须去
}
console.log([4, 2, 5, 3]) // nums 已经是连续的
console.log([1, 2, 3, 5, 6])
// 一个可能的解是将最后一个元素变为 4 。
// 结果数组为 [1,2,3,5,4] ，是连续数组。

export {}
