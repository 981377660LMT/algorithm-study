/**
 * @param {number[]} prices
 * @return {number}
 * 你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。
 */
let maxProfit = function (prices) {
  const len = prices.length

  // dp[i][0] 表示第i天持有股票所得最多现金 不持有 持有
  const dp = Array.from({ length: len }, () => Array(2).fill(0))
  dp[0] = [0, -prices[0]]
  for (let i = 1; i < len; i++) {
    dp[i][0] = Math.max(dp[i - 1][0], dp[i - 1][1] + prices[i])
    // 第i天买入股票，所得现金就是买入今天的股票后所得现金即：-prices[i]
    // 因为股票全程只能买卖一次，所以如果买入股票，那么第i天持有股票即dp[i][0]一定就是 -prices[i]
    dp[i][1] = Math.max(dp[i - 1][1], -prices[i])
  }

  // 最后不持有
  return dp[len - 1][0]
}

console.log(maxProfit([1, 2, 3, 4, 5]))

export {}
