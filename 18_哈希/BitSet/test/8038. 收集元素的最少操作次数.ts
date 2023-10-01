// 8038. 收集元素的最少操作次数
// https://leetcode.cn/problems/minimum-operations-to-collect-elements/description/

import { BitSet } from '../BitSet'

// 给你一个正整数数组 nums 和一个整数 k 。
// 一次操作中，你可以将数组的最后一个元素删除，将该元素添加到一个集合中。
// 请你返回收集元素 1, 2, ..., k 需要的 最少操作次数 。
// !其实只要记录每个数字出现的次数就行了，不需要用到BitSet.
function minOperations(nums: number[], k: number): number {
  const target = new BitSet(100)
  for (let i = 1; i <= k; i++) target.add(i)
  const cur = new BitSet(100)
  for (let i = nums.length - 1; ~i; i--) {
    cur.add(nums[i])
    if (cur.and(target).equals(target)) return nums.length - i
  }
  return -1
}
