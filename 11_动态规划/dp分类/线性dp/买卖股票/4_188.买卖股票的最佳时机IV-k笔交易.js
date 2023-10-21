/**
 * @param {number} k
 * @param {number[]} prices
 * @return {number}
 * 你最多可以完成  k 笔 交易。
 * 2 * k + 1 个状态
 * @summary 类比123题
 */
function maxProfit(k, prices) {
  if (prices == null || prices.length < 2 || k == 0) {
    return 0
  }

  const len = prices.length
  const dp = Array.from({ length: len }, () => Array(2 * k + 1).fill(0))
  dp[0][0] = 0
  dp[0][1] = -prices[0]
  for (let i = 2; i < 2 * k; i++) {
    dp[0][i] = -Infinity
  }

  for (let i = 1; i < len; i++) {
    for (let j = 0; j < 2 * k; j += 2) {
      dp[i][j + 1] = Math.max(dp[i - 1][j + 1], dp[i - 1][j] - prices[i])
      dp[i][j + 2] = Math.max(dp[i - 1][j + 2], dp[i - 1][j + 1] + prices[i])
    }
  }

  const res = dp[len - 1].filter((_, i) => i % 2 === 0)
  return Math.max.apply(null, res)
}

console.log(maxProfit(2, [3, 2, 6, 5, 0, 3]))
