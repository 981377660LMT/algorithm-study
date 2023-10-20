/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

/**
 * 数轴上按照顺序分布着n个点,每个点的位置为positions[i],权重为weights[i].
 * 点i到点j的距离定义为 `weights[i]*abs(positions[i]-positions[j])`.
 * 求区间[start,end)内的所有点到点to的距离之和.
 */
function distSumWeighted(positions: ArrayLike<number>, weights: ArrayLike<number>): (start: number, end: number, to: number) => number {
  const preSum: number[] = Array(weights.length + 1)
  const preMul: number[] = Array(weights.length + 1)
  preSum[0] = 0
  preMul[0] = 0
  for (let i = 0; i < weights.length; i++) {
    preSum[i + 1] = preSum[i] + weights[i]
    preMul[i + 1] = preMul[i] + weights[i] * positions[i]
  }

  const cal = (start: number, end: number, to: number, onLeft: boolean): number => {
    if (start >= end) return 0
    const res1 = (preSum[end] - preSum[start]) * positions[to]
    const res2 = preMul[end] - preMul[start]
    return onLeft ? res1 - res2 : res2 - res1
  }

  return (start: number, end: number, to: number): number => {
    const res1 = cal(start, Math.min(end, to), to, true)
    const res2 = cal(Math.max(start, to), end, to, false)
    return res1 + res2
  }
}

export { distSumWeighted }

if (require.main === module) {
  check()
  function check(): void {
    function checkWithBruteForce(positions: number[], weights: number[], start: number, end: number, to: number): number {
      let res = 0
      for (let i = start; i < end; i++) {
        res += weights[i] * Math.abs(positions[i] - positions[to])
      }
      return res
    }

    const n = 1e4
    const positions = Array(n)
      .fill(0)
      .map((_, i) => i)
    const weights = Array(n)
      .fill(0)
      .map(() => (Math.random() * 500) | 0)
    const Q = distSumWeighted(positions, weights)
    for (let i = 0; i < 1e4; i++) {
      const start = (Math.random() * n) | 0
      const end = (Math.random() * n) | 0
      const to = (Math.random() * n) | 0
      const res1 = Q(start, end, to)
      const res2 = checkWithBruteForce(positions, weights, start, end, to)
      if (res1 !== res2) {
        console.error(`res1: ${res1}, res2: ${res2}`)
        throw new Error()
      }
    }

    console.log('pass')
  }
}
