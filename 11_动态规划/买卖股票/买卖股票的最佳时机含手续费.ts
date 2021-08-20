/**
 *
 * @param prices 股票的价格
 * @summary dp[i][0]不持有 dp[i][1]持有
 */
const maxProfit = (prices: number[], fee: number) => {
  const len = prices.length

  // dp[i][0] 表示第i天持有股票所得最多现金 不持有 持有
  const dp = Array.from({ length: len }, () => Array(2).fill(0))
  dp[0] = [0, -prices[0]]
  for (let i = 1; i < len; i++) {
    dp[i][0] = Math.max(dp[i - 1][0], dp[i - 1][1] + prices[i] - fee)
    dp[i][1] = Math.max(dp[i - 1][1], dp[i - 1][0] - prices[i])
  }

  // 最后不持有
  return dp[len - 1][0]
}

console.log(maxProfit([1, 3, 2, 8, 4, 9], 2))

export {}
