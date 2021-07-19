/**
 *
 * @param prices 股票的价格
 * @summary 上帝视角：局部最优，见好就收，见差不动
 */
const maxProfit = (prices: number[]) => {
  if (prices.length === 0) return 0
  let allProfit = 0
  let currentPrice = prices[0]
  279
  for (let index = 1; index < prices.length; index++) {
    const priceToday = prices[index]
    if (priceToday > currentPrice) {
      allProfit += priceToday - currentPrice
      currentPrice = priceToday
    }
  }

  return allProfit
}

console.log(maxProfit([1, 2, 3, 4, 5]))

export {}
