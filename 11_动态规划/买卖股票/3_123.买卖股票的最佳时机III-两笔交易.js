/**
 * @param {number[]} prices
 * @return {number}
 * 你最多可以完成 两笔 交易。
 * @summary 一天一共就有五个状态 
 * 没有操作
   第一次买入
   第一次卖出
   第二次买入
   第二次卖出
 */
function maxProfit(prices) {
  const len = prices.length

  const dp = Array.from({ length: len }, () => Array(5).fill(0))
  dp[0][0] = 0
  dp[0][1] = -prices[0]
  dp[0][2] = -Infinity
  dp[0][3] = -Infinity
  dp[0][4] = -Infinity

  for (let i = 1; i < len; i++) {
    dp[i][0] = dp[i - 1][0]
    dp[i][1] = Math.max(dp[i - 1][1], dp[i - 1][0] - prices[i])
    dp[i][2] = Math.max(dp[i - 1][2], dp[i - 1][1] + prices[i])
    dp[i][3] = Math.max(dp[i - 1][3], dp[i - 1][2] - prices[i])
    dp[i][4] = Math.max(dp[i - 1][4], dp[i - 1][3] + prices[i])
  }

  return Math.max(dp[len - 1][0], dp[len - 1][2], dp[len - 1][4])
}

console.log(maxProfit([1, 2, 3, 4, 5]))
