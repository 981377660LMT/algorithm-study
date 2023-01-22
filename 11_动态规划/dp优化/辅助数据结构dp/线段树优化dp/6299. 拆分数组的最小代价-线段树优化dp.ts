/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { useAtcoderLazySegmentTree } from '../../../../6_tree/线段树/template/atcoder_segtree/AtcoderLazySegmentTree'

// 给你一个整数数组 nums 和一个整数 k 。
// 将数组拆分成一些非空子数组。拆分的 代价 是每个子数组中的 重要性 之和。
// !令 trimmed(subarray) 作为子数组的一个特征，其中所有仅出现一次的数字将会被移除。
// 例如，trimmed([3,1,2,4,3,4]) = [3,4,3,4] 。
// !子数组的 重要性 定义为 k + trimmed(subarray).length 。
// 例如，如果一个子数组是 [1,2,3,3,3,4,4] ，trimmed([1,2,3,3,3,4,4]) = [3,3,3,4,4] 。
// 这个子数组的重要性就是 k + 5 。
// !找出并返回拆分 nums 的所有可行方案中的最小代价。

// 1 <= nums.length <= 1e5
// 0 <= nums[i] < nums.length
// 1 <= k <= 1e9

// O(n^2)
function minCost1(nums: number[], k: number): number {
  const n = nums.length
  const dp = Array(n + 1).fill(INF)
  dp[0] = 0
  for (let i = 1; i < n + 1; i++) {
    const counter = new Map<number, number>()
    let curSum = 0
    for (let j = i; j < n + 1; j++) {
      counter.set(nums[j - 1], (counter.get(nums[j - 1]) || 0) + 1)
      if (counter.get(nums[j - 1]) === 2) curSum += 2
      else if (counter.get(nums[j - 1])! >= 3) curSum += 1
      dp[j] = Math.min(dp[j], dp[i - 1] + curSum + k)
    }
  }

  return dp[n]
}

// O(nlogn) 线段树优化
// dp转移方程变形:
// 用一个变量unique维护只出现一次的元素个数 首次出现unique++ 第二次出现unique--
// dp[i+1] = min(dp[j]+i-j+1-uniquej+k) (0<=j<=i) 即
// dp[i+1] - (i+1) = min((dp[j] - j) - uniquej) + k (0<=j<=i)
// 记 f[i] = dp[i] - i 则
// !f[i+1] = min(f[j] - uniquej) + k (0<=j<=i)
// 记x上一次出现的位置为pre[x],上上次出现的位置为pre2[x]
// 从左到右枚举 x=nums[i]
// !1.[pre[x]+1,i]里面的uniquej都加1
// !2.[pre2[x]+1,pre[x]]里面的uniquej都减1 (pre2[x]不存在则不用减)

const INF = 2e15
function minCost(nums: number[], k: number): number {
  const n = nums.length
  // rangeAddRangeMin
  const dp = useAtcoderLazySegmentTree(Array(n + 10).fill(0), {
    e: () => INF,
    id: () => 0,
    op: (a, b) => Math.min(a, b),
    mapping: (f, x) => f + x,
    composition: (f, g) => f + g
  })

  let res = 0
  const pre = new Map<number, number>()
  const pre2 = new Map<number, number>()
  for (let i = 1; i < n + 1; i++) {
    dp.update(i, i + 1, res)

    const num = nums[i - 1]
    dp.update((pre.get(num) || 0) + 1, i + 1, -1)
    if (pre2.has(num)) dp.update(pre2.get(num)! + 1, pre.get(num)! + 1, 1)

    res = k + dp.query(1, i + 1)
    pre2.set(num, pre.get(num) || 0)
    pre.set(num, i)
  }

  return n + res
}

if (require.main === module) {
  console.log(minCost([1, 2, 1, 2, 1, 3, 3], 2))
}

export {}
