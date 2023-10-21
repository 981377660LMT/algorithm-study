/**
 * @param {number[]} prices
 * @return {number}
 * 你可以尽可能地完成更多的交易（多次买卖一支股票）。
 */
function maxProfit(prices) {
  let res = 0

  for (let i = 1; i < prices.length; i++) {
    res += Math.max(prices[i] - prices[i - 1], 0)
  }

  return res
}

console.log(maxProfit([1, 2, 3, 4, 5]))

export {}
