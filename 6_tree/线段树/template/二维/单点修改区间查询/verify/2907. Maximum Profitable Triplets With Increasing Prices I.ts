// 2907.最大递增三元组的和
// https://leetcode.cn/problems/maximum-profitable-triplets-with-increasing-prices-i/description/

// 找到三个下标 i, j, k,使得 i < j < k 且 prices[i] < prices[j] < prices[k],
// 并且 profits[i] + profits[j] + profits[k] 最大。
// 如果无法找到则返回 -1。

// !三元组:枚举中间.

import { SegmentTree2DSparse } from '../SegmentTree2DSparse'

const INF = 1e9

function maxProfit(prices: number[], profits: number[]): number {
  const n = prices.length
  const xs = Array<number>(n)
  for (let i = 0; i < n; i++) xs[i] = i
  const tree = new SegmentTree2DSparse({
    xs,
    ys: prices,
    ws: profits,
    e: () => 0,
    op: Math.max,
    discretizeX: false
  })

  let res = -1
  for (let i = 0; i < n; i++) {
    const curX = i
    const curY = prices[i]
    const max1 = tree.query(0, curX, 0, curY)
    if (max1 === 0) continue
    const max2 = tree.query(curX + 1, INF, curY + 1, INF)
    if (max2 === 0) continue
    res = Math.max(res, max1 + max2 + profits[i])
  }
  return res
}

if (require.main === module) {
  // prices = [10,2,3,4], profits = [100,2,7,10]
  console.log(maxProfit([10, 2, 3, 4], [100, 2, 7, 10]))
}
