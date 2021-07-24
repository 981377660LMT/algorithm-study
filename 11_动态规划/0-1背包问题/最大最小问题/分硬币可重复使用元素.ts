/**
 * @description 计算可以凑成总金额所需的最少的硬币个数
 * @summary 不考虑排列顺序的完全背包问题
 */
const canPartition = function (nums: number[], amount: number) {
  // dp[i]表示面值为i时需要硬币的最少数量
  const dp = Array(amount + 1).fill(Infinity)
  dp[0] = 0
  for (const coin of nums) {
    for (let i = 0; i <= amount; i++) {
      if (i - coin >= 0) {
        dp[i] = Math.min(dp[i], dp[i - coin] + 1)
      }
    }
  }
  return dp[amount] === Infinity ? -1 : dp[amount]
}

console.dir(canPartition([1, 2, 5], 11), { depth: null })

export {}
