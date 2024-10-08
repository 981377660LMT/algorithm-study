/* eslint-disable max-len */

// 3181. 执行操作可获得的最大总奖励 II
// leetcode.cn/problems/maximum-total-reward-using-operations-ii/solutions/2805488/typescript-bitsetyou-hua-ke-xing-xing-01-mss1/
// TypeScript bitset优化可行性01背包问题

// BigInt 大数模拟
const maxTotalReward = (rewardValues: number[]): number => {
  rewardValues = [...new Set(rewardValues)].sort((a, b) => a - b)
  const Big1 = BigInt(1)
  let res = Big1
  rewardValues.forEach(v => {
    const BigV = BigInt(v)
    const low = ((Big1 << BigV) - Big1) & res
    res |= low << BigV
  })
  return res.toString(2).length - 1
}

export {}
