/* eslint-disable eqeqeq */

// 2863. 最长半递减数组
// https://leetcode.cn/problems/maximum-length-of-semi-decreasing-subarrays/description/
// !返回最长的子数组，子数组的第一个元素严格大于最后一个元素
// 分块加速查找.分块的好处是，可以支持修改操作.

import { useBlock } from './SqrtDecomposition/useBlock'

function maxSubarrayLength(nums: number[]): number {
  const { belong, blockStart, blockEnd, blockCount } = useBlock(nums)
  const blockMin = Array(blockCount).fill(Infinity)
  nums.forEach((num, i) => {
    const bid = belong[i]
    blockMin[bid] = Math.min(blockMin[bid], num)
  })

  let res = 0
  for (let i = 0; i < nums.length; i++) {
    const rightMostLower = queryRightMostLower(i)
    if (rightMostLower != undefined) {
      res = Math.max(res, rightMostLower - i + 1)
    }
  }
  return res

  /**
   * 查询下标 `pos` 右侧最远的比 `nums[pos]` 严格小的元素的下标.
   */
  function queryRightMostLower(pos: number): number | undefined {
    const bid = belong[pos]
    const target = nums[pos]
    for (let i = blockCount - 1; i > bid; i--) {
      if (blockMin[i] >= target) continue
      for (let j = blockEnd[i] - 1; j >= blockStart[i]; j--) {
        if (nums[j] < target) return j
      }
    }

    for (let i = blockEnd[bid] - 1; i > pos; i--) {
      if (nums[i] < target) return i
    }

    return undefined
  }
}

if (require.main === module) {
  // [7,6,5,4,3,2,1,6,10,11]
  console.log(maxSubarrayLength([7, 6, 5, 4, 3, 2, 1, 6, 10, 11]))
}
