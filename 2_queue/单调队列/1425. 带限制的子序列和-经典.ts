// 此题与 1696. 跳跃游戏 VI 相同

import { ArrayDeque } from '../Deque/ArrayDeque'

/**
 * @param {number[]} nums  -10^4 <= nums[i] <= 10^4
 * @param {number} k
 * @return {number}
 * 请你返回 非空 子序列元素和的最大值，子序列需要满足：
 * 子序列中每两个 相邻 的整数 nums[i] 和 nums[j] ，它们在原数组中的下标 i 和 j 满足 i < j 且 j - i <= k 。
 * @summary
 * 单减的单调队列，队首表示之前的子序列元素和最大值
 */
const constrainedSubsetSum = function (nums: number[], k: number): number {
  const n = nums.length
  const queue = new ArrayDeque<[preMaxSum: number, index: number]>(n)
  queue.push([nums[0], 0])

  const dp = nums.slice()
  dp[0] = nums[0]
  let res = nums[0]

  for (let i = 1; i < n; i++) {
    while (queue.length > 0 && i - queue.at(0)![1] > k) queue.shift()
    // 选货不选
    dp[i] = Math.max(dp[i], queue.at(0)![0] + nums[i])
    while (queue.length > 0 && queue.at(-1)![0] <= dp[i]) queue.pop()
    queue.push([dp[i], i])
    res = Math.max(res, dp[i])
  }

  return res
}

console.log(constrainedSubsetSum([10, 2, -10, 5, 20], 2))
// 输出：37
// 解释：子序列为 [10, 2, 5, 20] 。
