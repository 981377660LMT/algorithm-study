const INF = 2e15

/**
 * @param {number[]} prices
 * @return {number}
 * 你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。
 * @summary 求波峰浪谷的差值的最大值。
 */
function maxProfit(prices) {
  let res = 0
  let preMin = INF
  const n = prices.length

  for (let i = 0; i < n; i++) {
    res = Math.max(res, prices[i] - preMin)
    preMin = Math.min(preMin, prices[i])
  }

  return res
}

console.log(maxProfit([1, 2, 3, 4, 5]))
