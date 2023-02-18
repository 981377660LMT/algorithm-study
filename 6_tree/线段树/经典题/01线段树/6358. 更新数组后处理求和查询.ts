// 给你两个下标从 0 开始的数组 nums1 和 nums2 ，和一个二维数组 queries 表示一些操作。总共有 3 种类型的操作：

// 操作类型 1 为 queries[i] = [1, l, r] 。你需要将 nums1 从下标 l 到下标 r 的所有 0 反转成 1 或将 1 反转成 0 。l 和 r 下标都从 0 开始。
// 操作类型 2 为 queries[i] = [2, p, 0] 。对于 0 <= i < n 中的所有下标，令 nums2[i] = nums2[i] + nums1[i] * p 。
// 操作类型 3 为 queries[i] = [3, 0, 0] 。求 nums2 中所有元素的和。
// 请你返回一个数组，包含所有第三种操作类型的答案。

import { SegmentTree01 } from './SegmentTree01'

function handleQuery(nums1: number[], nums2: number[], queries: number[][]): number[] {
  const n = nums1.length
  const seg01 = new SegmentTree01(nums1)
  let sum = nums2.reduce((a, b) => a + b, 0)
  const res: number[] = []
  for (const [op, a, b] of queries) {
    if (op === 1) {
      seg01.flip(a + 1, b + 1)
    } else if (op === 2) {
      const one = seg01.onesCount(1, n)
      sum += one * a
    } else {
      res.push(sum)
    }
  }
  return res
}

export {}
