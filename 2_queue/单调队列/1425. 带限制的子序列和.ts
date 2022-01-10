// 此题与 1696. 跳跃游戏 VI 相同

import { ArrayDeque } from '../Deque/ArrayDeque'

type Sum = number
type Index = number

/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * 请你返回 非空 子序列元素和的最大值，子序列需要满足：
 * 子序列中每两个 相邻 的整数 nums[i] 和 nums[j] ，它们在原数组中的下标 i 和 j 满足 i < j 且 j - i <= k 。
 * @summary
 * 单调队列维护窗口最大值
 */
const constrainedSubsetSum = function (nums: number[], k: number): number {
  const n = nums.length
  const queue = new ArrayDeque<[Sum, Index]>(10 ** 5)
  queue.push([nums[0], 0])
  let res = nums[0]
  const dp = nums.slice()
  dp[0] = nums[0]

  for (let i = 1; i < n; i++) {
    if (i - queue.front()![1] > k) queue.shift()
    // 选货不选
    dp[i] = Math.max(dp[i], queue.front()![0] + nums[i])
    while (queue.length && queue.rear()![0] <= dp[i]) queue.pop()
    queue.push([dp[i], i])
    res = Math.max(res, dp[i])
  }

  return res
}

console.log(constrainedSubsetSum([10, 2, -10, 5, 20], 2))
// 输出：37
// 解释：子序列为 [10, 2, 5, 20] 。
