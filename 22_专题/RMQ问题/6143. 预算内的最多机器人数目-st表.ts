/* eslint-disable no-shadow */

// cost = max(chargeTimes) + k * sum(runningCosts) ，
// 其中 max(chargeTimes) 是这 k 个机器人中最大充电时间，
// sum(runningCosts) 是这 k 个机器人的运行时间之和。
// !请你返回在 不超过 budget 的前提下，你 最多 可以 连续 运行的机器人数目为多少。
// !二分答案+st表查询区间最值

import { SparseTable } from '../../22_专题/RMQ问题/SparseTable'

function maximumRobots(chargeTimes: number[], runningCosts: number[], budget: number): number {
  const n = chargeTimes.length
  const st = new SparseTable(chargeTimes, Math.max)
  const preSum = Array<number>(n + 1).fill(0)
  for (let i = 1; i <= n; i++) preSum[i] = preSum[i - 1] + runningCosts[i - 1]

  let left = 1
  let right = n
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  function check(mid: number): boolean {
    for (let left = 0; left + mid - 1 < n; left++) {
      const right = left + mid - 1
      const max = st.query(left, right)
      const sum = preSum[right + 1] - preSum[left]
      if (max + mid * sum <= budget) return true
    }
    return false
  }
}

// [11,12,74,67,37,87,42,34,18,90,36,28,34,20]
// [18,98,2,84,7,57,54,65,59,91,7,23,94,20]
// 937
// 预期4

if (require.main === module) {
  console.log(
    maximumRobots(
      [11, 12, 74, 67, 37, 87, 42, 34, 18, 90, 36, 28, 34, 20],
      [18, 98, 2, 84, 7, 57, 54, 65, 59, 91, 7, 23, 94, 20],
      937
    )
  )
}
export {}
