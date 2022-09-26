// 考虑 nums 中进行 按位与（bitwise AND）运算得到的值 最大 的 非空 子数组。
// !返回满足要求的 最长 子数组的长度。

import assert from 'assert'
import { SparseTable } from './SparseTable'

function longestSubarray(nums: number[]): number {
  const n = nums.length
  const max = Math.max(...nums)
  const st = new SparseTable(nums, (a, b) => a & b)
  let res = 1

  for (let start = 0; start < n; start++) {
    let left = start
    let right = n - 1

    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      if (st.query(start, mid) === max) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    res = Math.max(res, right - start + 1)
  }

  return res
}

// !其实遍历一遍 找最大值最长连续个数就行了
function longestSubarray2(nums: number[]): number {
  const max = Math.max(...nums)

  let res = 0
  let dp = 0
  for (const num of nums) {
    if (num === max) {
      dp++
      res = Math.max(res, dp)
    } else {
      dp = 0
    }
  }

  return res
}

if (require.main === module) {
  assert.strictEqual(longestSubarray([1, 2, 3, 3, 2, 2]), 2)
  assert.strictEqual(longestSubarray([1, 2, 3, 4]), 1)
}
