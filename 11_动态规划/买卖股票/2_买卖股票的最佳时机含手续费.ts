/**
 *
 * @param prices 股票的价格
 * @summary 上帝视角：局部最优，见好就收，见差不动 画图
 * 给定一个数组 prices ，其中 prices[i] 是一支给定股票第 i 天的价格。
 * 设计一个算法来计算你所能获取的最大利润。
 */
const maxProfit = (prices: number[], fee: number) => {
  let result = 0

  for (let i = 1; i < prices.length; i++) {
    result += Math.max(prices[i] - prices[i - 1], 0)
  }

  return result
}

console.log(maxProfit([1, 3, 2, 8, 4, 9], 2))

export {}
