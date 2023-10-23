import { gcd } from '../数论/扩展欧几里得/gcd'

/**
 * 统计数组所有子数组的 gcd 的不同个数，复杂度 O(n*log^2max)
 */
function countGcdOfAllSubarray(nums: number[]): number {
  return logTrick(nums, gcd).size
}

/**
 * 将 `arr` 的所有非空子数组的元素进行 `op` 操作，返回所有不同的结果和其出现次数.
 * @param arr 1 <= arr.length <= 1e5.
 * @param op 与/或/gcd/lcm 中的一种操作，具有单调性.
 * @param f `arr[:end]` 中所有子数组的结果为 `preCounter`.
 */
function logTrick(
  arr: ArrayLike<number>,
  op: (a: number, b: number) => number,
  f?: (end: number, preCounter: ReadonlyMap<number, number>) => void
): Map<number, number> {
  const res: Map<number, number> = new Map()
  const dp: number[] = []
  for (let pos = 0; pos < arr.length; pos++) {
    const cur = arr[pos]
    for (let i = 0; i < dp.length; i++) {
      dp[i] = op(dp[i], cur)
    }
    dp.push(cur)

    // 去重
    let ptr = 0
    for (let i = 1; i < dp.length; i++) {
      if (dp[i] !== dp[ptr]) {
        ptr++
        dp[ptr] = dp[i]
      }
    }

    dp.length = ptr + 1
    for (let i = 0; i < dp.length; i++) {
      res.set(dp[i], (res.get(dp[i]) || 0) + 1)
    }
    f && f(pos + 1, res)
  }

  return res
}

if (require.main === module) {
  console.log(countGcdOfAllSubarray([6, 10, 15]))
  console.log(countGcdOfAllSubarray([5, 5, 5]))
}
