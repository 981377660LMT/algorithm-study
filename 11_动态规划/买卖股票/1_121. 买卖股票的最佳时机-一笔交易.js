/**
 * @param {number[]} prices
 * @return {number}
 * 你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。
 * @summary 求波峰浪谷的差值的最大值。
 */
var maxProfit = function (prices) {
  let min = prices[0]
  let profit = 0
  // 7 1 5 3 6 4
  for (let i = 1; i < prices.length; i++) {
    profit = Math.max(profit, prices[i] - min)
    min = Math.min(min, prices[i])
  }

  return profit
}

console.log(maxProfit([1, 2, 3, 4, 5]))

export {}
