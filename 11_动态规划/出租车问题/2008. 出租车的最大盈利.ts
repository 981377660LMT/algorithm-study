/* eslint-disable no-param-reassign */

import { bisectRight } from '../../9_排序和搜索/二分/bisect'

/**
 *
 * @param n  这 n 个地点从近到远编号为 1 到 n  1 <= n <= 105
 * @param rides  i 位乘客需要从地点 starti 前往 endi ，愿意支付 tipi 元的小费。
 * 1 <= rides.length <= 3e4
 * 通过接乘客订单盈利。你只能沿着编号递增的方向前进，不能改变方向。
 * 请你返回在最优接单方案下，你能盈利 最多 多少元。
 */
function maxTaxiEarnings(n: number, rides: number[][]): number {
  n = rides.length
  rides.sort((a, b) => a[1] - b[1])
  const dp = Array<number>(n + 1).fill(0)
  for (let i = 0; i < n; i++) {
    dp[i + 1] = dp[i] // to jump
    const [start, end, tip] = rides[i] // not to jump
    const score = end - start + tip
    const prePos = bisectRight(rides, start, { key: e => e[1] })
    dp[i + 1] = Math.max(dp[i + 1], dp[prePos] + score)
  }

  return dp[n]
}

if (require.main === module) {
  console.log(
    maxTaxiEarnings(10, [
      [2, 3, 6],
      [8, 9, 8],
      [5, 9, 7],
      [8, 9, 1],
      [2, 9, 2],
      [9, 10, 6],
      [7, 10, 10],
      [6, 7, 9],
      [4, 9, 7],
      [2, 3, 1]
    ])
  )
}
// 输出：33

export { maxTaxiEarnings }
