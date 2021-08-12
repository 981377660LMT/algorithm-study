/**
 * @param {number} amount
 * @param {number[]} coins
 * @return {number}
 * @description 计算并返回可以凑成总金额的硬币组合数:不考虑排列顺序的完全背包问题
 */
var change = function (amount, coins) {
  const dp = Array(amount + 1).fill(0)
  dp[0] = 1
  for (const coin of coins) {
    for (let i = 0; i <= amount; i++) {
      if (i - coin >= 0) dp[i] = dp[i] + dp[i - coin]
    }
  }
  return dp[amount]
}

console.log(change(5, [1, 2, 5]))
// 输出：4
// 解释：有四种方式可以凑成总金额：
// 5=5
// 5=2+2+1
// 5=2+1+1+1
// 5=1+1+1+1+1
