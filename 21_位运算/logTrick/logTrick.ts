/* eslint-disable no-inner-declarations */

import { gcd, lcm } from '../../19_数学/数论/扩展欧几里得/gcd'

type Interval = { leftStart: number; leftEnd: number; value: number }

/**
 * 将 `arr` 的所有非空子数组的元素进行 `op` 操作，返回所有不同的结果和其出现次数.
 * @param arr 1 <= arr.length <= 1e5.
 * @param op 与/或/gcd/lcm 中的一种操作，具有单调性.
 * @param f
 * 子数组的右端点为right.
 * interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
 * interval 的 value 表示该子数组 arr[left,right] 的 op 结果.
 */
function logTrick(arr: ArrayLike<number>, op: (a: number, b: number) => number, f?: (left: Interval[], right: number) => void): Map<number, number> {
  const res: Map<number, number> = new Map()
  const dp: Interval[] = []
  for (let pos = 0; pos < arr.length; pos++) {
    const cur = arr[pos]
    for (let i = 0; i < dp.length; i++) {
      dp[i].value = op(dp[i].value, cur)
    }
    dp.push({ leftStart: pos, leftEnd: pos + 1, value: cur })

    // 去重
    let ptr = 0
    for (let i = 1; i < dp.length; i++) {
      if (dp[i].value !== dp[ptr].value) {
        ptr++
        dp[ptr] = dp[i]
      } else {
        dp[ptr].leftEnd = dp[i].leftEnd
      }
    }
    dp.length = ptr + 1

    // 将区间[0,pos]分成了dp.length个左闭右开区间.
    // 每一段区间的左端点left范围 在 [dp[i].leftStart,dp[i].leftEnd) 中。
    // 对应子数组 arr[left:pos+1] 的 op 值为 dp[i].value.
    for (let i = 0; i < dp.length; i++) {
      const { leftStart, leftEnd, value } = dp[i]
      res.set(value, (res.get(value) || 0) + leftEnd - leftStart)
    }
    f && f(dp, pos)
  }

  return res
}

export { logTrick }

if (require.main === module) {
  // https://leetcode.cn/problems/bitwise-ors-of-subarrays/
  // 898. 子数组按位或操作
  function subarrayBitwiseORs(arr: number[]): number {
    return logTrick(arr, (a, b) => a | b).size
  }

  // https://leetcode.cn/problems/find-a-value-of-a-mysterious-function-closest-to-target/solutions/343107/zhao-dao-zui-jie-jin-mu-biao-zhi-de-han-shu-zhi-by/
  // 1521. 找到最接近目标值的函数值
  function closestToTarget(arr: number[], target: number): number {
    const counter = logTrick(arr, (a, b) => a & b)
    let res = Infinity
    counter.forEach((_, key) => {
      res = Math.min(res, Math.abs(target - key))
    })
    return res
  }

  // 2447. 最大公因数等于 K 的子数组数目
  function subarrayGCD(nums: number[], k: number): number {
    const counter = logTrick(nums, gcd)
    return counter.get(k) || 0
  }

  // 2470. 最小公倍数为 K 的子数组数目
  // https://leetcode.cn/problems/number-of-subarrays-with-lcm-equal-to-k/
  function subarrayLCM(nums: number[], k: number): number {
    const counter = logTrick(nums, lcm)
    return counter.get(k) || 0
  }

  // 2941. 子数组的最大 GCD-Sum
  // https://leetcode.cn/problems/maximum-gcd-sum-of-a-subarray/description/
  function maxGcdSum(nums: number[], k: number): number {
    const preSum = Array(nums.length + 1).fill(0)
    for (let i = 0; i < nums.length; i++) preSum[i + 1] = preSum[i] + nums[i]
    let res = 0
    logTrick(nums, gcd, (leftIntervals, right) => {
      for (let j = 0; j < leftIntervals.length; j++) {
        const { leftStart, value } = leftIntervals[j]
        if (right - leftStart + 1 >= k) {
          res = Math.max(res, value * (preSum[right + 1] - preSum[leftStart]))
        }
      }
    })
    return res
  }
}
