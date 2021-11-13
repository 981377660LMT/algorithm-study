/**
 *
 * @param n  这 n 个地点从近到远编号为 1 到 n  1 <= n <= 105
 * @param rides  i 位乘客需要从地点 starti 前往 endi ，愿意支付 tipi 元的小费。
 * 1 <= rides.length <= 3 * 104
 * 通过接乘客订单盈利。你只能沿着编号递增的方向前进，不能改变方向。
 * 请你返回在最优接单方案下，你能盈利 最多 多少元。
 * // 1. dp[i]表示选择接乘客i为结尾时所能达到的最大盈利，所以要对rides先排序
   // 2. 初始化dp
   // 3. 类似于LIS问题，对每个乘客，找到之前最右的乘客，使preEnd<=curStart
   //    注意bisectRight是寻找新的插入坐标 找之前的乘客需要再减1
   // 4. 如果pre为0 即之前不能带乘客 那么带或不带当前乘客 dp[i] = Math.max(dp[i], dp[i - 1]) 否则 dp[i] = Math.max(dp[i - 1], dp[pre - 1] + rides[i][1] - rides[i][0] + rides[i][2])
 */
function maxTaxiEarnings(n: number, rides: number[][]): number {
  rides.sort((a, b) => a[1] - b[1])

  // dps[i]表示选择接乘客i为结尾时所能达到的最大盈利
  const dp = Array(rides.length).fill(0)
  rides.forEach(([start, end, tip], index) => {
    dp[index] = end - start + tip
  })

  for (let i = 1; i < rides.length; i++) {
    const preIndex = bisectRight(rides, i)
    if (preIndex === 0) {
      // 带或不带
      dp[i] = Math.max(dp[i], dp[i - 1])
    } else {
      dp[i] = Math.max(dp[i - 1], dp[preIndex - 1] + rides[i][1] - rides[i][0] + rides[i][2])
    }
  }

  return Math.max(...dp)

  function bisectRight(rides: number[][], rideIndex: number) {
    let l = 0
    let r = rideIndex
    const curStart = rides[rideIndex][0]

    while (l <= r) {
      const mid = (l + r) >> 1
      const preEnd = rides[mid][1]
      if (preEnd === curStart) l++
      else if (preEnd > curStart) r = mid - 1
      else l = mid + 1
    }

    return l
  }
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
      [2, 3, 1],
    ])
  )
}
// 输出：33

export { maxTaxiEarnings }
