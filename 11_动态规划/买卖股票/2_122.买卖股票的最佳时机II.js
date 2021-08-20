/**
 * @param {number[]} prices
 * @return {number}
 * 你可以尽可能地完成更多的交易（多次买卖一支股票）。
 */
var maxProfit = function (prices) {
  const len = prices.length

  // dp[i][0] 表示第i天持有股票所得最多现金 不持有 持有
  const dp = Array.from({ length: len }, () => Array(2).fill(0))
  dp[0] = [0, -prices[0]]
  for (let i = 1; i < len; i++) {
    dp[i][0] = Math.max(dp[i - 1][0], dp[i - 1][1] + prices[i])
    // 第i天买入股票,买入股票的时候，可能会有之前买卖的利润即
    dp[i][1] = Math.max(dp[i - 1][1], dp[i - 1][0] - prices[i])
  }

  // 最后不持有
  return dp[len - 1][0]
}

console.log(maxProfit([1, 2, 3, 4, 5]))

export {}
