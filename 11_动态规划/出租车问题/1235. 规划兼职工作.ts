import { zip } from '../../0_数组/zip'

function jobScheduling(startTime: number[], endTime: number[], profit: number[]): number {
  const rides = [...zip(startTime, endTime, profit)]
  rides.sort((a, b) => a[1] - b[1])

  // dp[i]表示选择接乘客i为结尾时所能达到的最大盈利
  const dp = Array(rides.length).fill(0)
  rides.forEach(([_start, _end, profit], index) => {
    dp[index] = profit
  })

  for (let i = 1; i < rides.length; i++) {
    const pre = bisectRight(rides, i)
    // 带或不带
    if (pre === 0) {
      dp[i] = Math.max(dp[i - 1], rides[i][2])
    } else {
      dp[i] = Math.max(dp[i - 1], dp[pre - 1] + rides[i][2])
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
      if (preEnd > curStart) r = mid - 1
      else l = mid + 1
    }

    return l
  }
}

console.log(jobScheduling([1, 2, 3, 3], [3, 4, 5, 6], [50, 10, 40, 70]))
