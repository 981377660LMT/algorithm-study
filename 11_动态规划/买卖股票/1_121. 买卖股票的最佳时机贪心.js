/**
 * @param {number[]} prices
 * @return {number}
 * 你只能选择 某一天 买入这只股票，并选择在 未来的某一个不同的日子 卖出该股票。
 */
var maxProfit = function (prices) {
  let result = 0

  for (let i = 1; i < prices.length; i++) {
    result += Math.max(prices[i] - prices[i - 1], 0)
  }

  return result
}

console.log(maxProfit([1, 2, 3, 4, 5]))

export {}
