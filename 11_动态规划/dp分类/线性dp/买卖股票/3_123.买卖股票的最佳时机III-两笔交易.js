/* eslint-disable prefer-destructuring */

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
  const n = prices.length
  let dp = new Int32Array([0, -prices[0], -1 << 30, -1 << 30, -1 << 30])

  for (let i = 1; i < n; i++) {
    const ndp = new Int32Array(5)
    ndp[0] = dp[0]
    ndp[1] = Math.max(dp[1], dp[0] - prices[i])
    ndp[2] = Math.max(dp[2], dp[1] + prices[i])
    ndp[3] = Math.max(dp[3], dp[2] - prices[i])
    ndp[4] = Math.max(dp[4], dp[3] + prices[i])
    dp = ndp
  }

  return Math.max(...dp)
}

console.log(maxProfit([1, 2, 3, 4, 5]))
