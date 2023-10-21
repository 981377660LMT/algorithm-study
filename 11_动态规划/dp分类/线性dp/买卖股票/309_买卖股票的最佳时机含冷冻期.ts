/**
 * @param prices 股票的价格
 * @description 你不能同时参与多笔交易;卖出股票后，你无法在第二天买入股票 (即冷冻期为 1 天)。
 * 定义二维dp[i][n]数组: i代表天数，n代表是股票情况，0代表不持有， 1代表持有，2代表冷冻期；
 * 初始化分以下三种情况：
 * dp[0][0] = 0: 第一天不持有股票
 * dp[0][1] = -prices[0]: 第一天持有股票, 利润为-prices[0];
 * dp[0][2] = 0: 第一天为冷冻期，
 * 状态转移方程分三种情况：
 * dp[i][0]：第i天不持有股票， 比较前一天不持有和前一天过渡期情况中最大值；
 * dp[i][1]：第i天持有股票， 比较前一天持有/前一天不持有情况中最大值；(冷冻期下一天不能买)
 * dp[i][2]：第i天冷冻期只有一种情况，当天持有股票且卖出；
 */
const maxProfit = (prices: number[]) => {
  const len = prices.length
  if (len < 2) {
    return 0
  }

  const dp = Array.from({ length: len }, () => Array(3).fill(0))
  dp[0][0] = 0
  dp[0][1] = -prices[0]
  dp[0][2] = -Infinity
  for (let i = 1; i < len; i++) {
    dp[i][0] = Math.max(dp[i - 1][0], dp[i - 1][2])
    dp[i][1] = Math.max(dp[i - 1][0] - prices[i], dp[i - 1][1])
    dp[i][2] = Math.max(dp[i - 1][1] + prices[i])
  }
  console.table(dp)

  return Math.max(dp[len - 1][0], dp[len - 1][2])
}

console.log(maxProfit([4, 3, 2]))
console.log(maxProfit([1, 2, 3, 0, 2]))

export {}
