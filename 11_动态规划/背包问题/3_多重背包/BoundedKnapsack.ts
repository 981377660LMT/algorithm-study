/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

/**
 * 单调队列优化多重背包问题.
 * 选择物品使得价值最大，总重量不超过{@link maxCapacity}.求出最大价值.
 * @param values 物品价值.
 * @param weights 物品重量.
 * @param counts 每种物品的数量.
 * @param maxCapacity 背包容量.
 * @complexity O(n*maxCapacity)
 */
function boundedKnapsackDp(values: ArrayLike<number>, weights: ArrayLike<number>, counts: ArrayLike<number>, maxCapacity: number): number {
  type Item = { max: number; j: number }

  const n = values.length
  const dp = Array(maxCapacity + 1).fill(0)
  const queue = Array<Item>(maxCapacity + 1)
  let head = 0
  let tail = 0
  for (let i = 0; i < n; i++) {
    const count = counts[i]
    const value = values[i]
    const weight = weights[i]
    // 按照 j%weight 的结果`分组转移`.
    for (let remainder = 0; remainder < weight; remainder++) {
      head = 0
      tail = 0
      for (let j = 0; j * weight + remainder <= maxCapacity; j++) {
        const cand = dp[j * weight + remainder] - j * value
        while (head < tail && queue[tail - 1].max <= cand) tail--
        queue[tail++] = { max: cand, j }
        dp[j * weight + remainder] = queue[head].max + j * value // 物品个数为两个 j 的差（前缀和思想）
        if (j - queue[head].j === count) head++ // 及时去掉无用数据
      }
    }
  }

  return dp[maxCapacity]
}

/**
 * 多重背包二进制优化.
 * !由于js语言特性，速度比单调队列优化更快.
 */
function boundedKnapsackDpBinary(values: ArrayLike<number>, weights: ArrayLike<number>, counts: ArrayLike<number>, maxCapacity: number): number {
  const n = values.length
  const dp = new Float64Array(maxCapacity + 1)
  for (let i = 0; i < n; i++) {
    let remain = counts[i]
    const value = values[i]
    const weight = weights[i]
    for (let step = 1; remain > 0; step *= 2) {
      const k = Math.min(step, remain)
      for (let j = maxCapacity; j >= k * weight; j--) {
        dp[j] = Math.max(dp[j], dp[j - k * weight] + k * value)
      }
      remain -= k
    }
  }

  return dp[maxCapacity]
}

/**
 * 多重背包朴素解法.
 */
function boundedKnapsackDpNaive(values: ArrayLike<number>, weights: ArrayLike<number>, counts: ArrayLike<number>, maxCapacity: number): number {
  const n = values.length
  const dp = Array<number>((n + 1) * (maxCapacity + 1)).fill(0)
  for (let i = 0; i < n; i++) {
    const count = counts[i]
    const value = values[i]
    const weight = weights[i]
    for (let j = 0; j < maxCapacity + 1; j++) {
      // 枚举选了 k=0,1,2,...num 个第 i 种物品
      for (let k = 0; k <= count && k * weight <= j; k++) {
        dp[(i + 1) * (maxCapacity + 1) + j] = Math.max(dp[(i + 1) * (maxCapacity + 1) + j], dp[i * (maxCapacity + 1) + j - k * weight] + k * value)
      }
    }
  }

  return dp[n * (maxCapacity + 1) + maxCapacity]
}

const MOD = 1e9 + 7

/**
 * 多重背包求方案数(分组前缀和优化).
 * @param values 物品价值.
 * @param counts 每种物品的数量.
 * @param upper 价值上限.
 * @returns dp[i] 表示总价值为 i 的方案数.
 * @complexity O(n*upper).
 */
function boundedKnapsackDpCountWays(values: ArrayLike<number>, counts: ArrayLike<number>, upper?: number): number[] {
  const hasUpper = upper != undefined

  const n = values.length
  let allSum = 0
  let count0 = 0
  for (let i = 0; i < n; i++) {
    const count = counts[i]
    const value = values[i]
    if (value === 0) {
      count0 += count
      continue
    }
    if (!hasUpper) allSum += count * value
  }
  if (hasUpper) allSum = upper
  const dp = Array<number>(allSum + 1).fill(0)
  dp[0] = count0 + 1

  let maxJ = 0
  for (let i = 0; i < n; i++) {
    const value = values[i]
    if (value === 0) {
      continue
    }
    const count = counts[i]
    maxJ = Math.min(maxJ + count * value, allSum)
    for (let j = value; j <= maxJ; j++) {
      dp[j] = (dp[j] + dp[j - value]) % MOD // 同余前缀和
    }
    for (let j = maxJ; j >= value * (count + 1); j--) {
      dp[j] = (dp[j] - dp[j - value * (count + 1)]) % MOD
    }
  }

  for (let i = 0; i < dp.length; i++) {
    if (dp[i] < 0) dp[i] += MOD
  }
  return dp
}

/**
 * 多重背包求方案数.
 */
class BoundedKnapsack {
  private readonly _dp: number[]
  private readonly _mod?: number | undefined
  private readonly _maxValue: number
  private _maxJ = 0

  constructor(maxValue: number, mod?: number, dp?: number[]) {
    if (dp != undefined) {
      this._dp = dp
    } else {
      this._dp = Array(maxValue + 1).fill(0)
      this._dp[0] = 1
    }
    this._maxValue = maxValue
    this._mod = mod
  }

  /**
   * 加入一个价值为value(value>0)的物品，数量为count.
   * @complexity O(maxValue)
   */
  add(value: number, count: number): void {
    if (value <= 0) throw new Error(`value must be positive, but got ${value}`)
    this._maxJ = Math.min(this._maxJ + count * value, this._maxValue)
    if (this._mod == undefined) {
      for (let j = value; j <= this._maxJ; j++) {
        this._dp[j] += this._dp[j - value]
      }
      for (let j = this._maxJ; j >= value * (count + 1); j--) {
        this._dp[j] -= this._dp[j - value * (count + 1)]
      }
    } else {
      for (let j = value; j <= this._maxJ; j++) {
        this._dp[j] = (this._dp[j] + this._dp[j - value]) % this._mod
      }
      for (let j = this._maxJ; j >= value * (count + 1); j--) {
        this._dp[j] = (this._dp[j] - this._dp[j - value * (count + 1)]) % this._mod
      }
    }
  }

  query(value: number): number {
    if (value < 0 || value > this._maxValue) return 0
    if (this._mod == undefined) return this._dp[value]
    if (this._dp[value] < 0) this._dp[value] += this._mod
    return this._dp[value]
  }

  copy(): BoundedKnapsack {
    const res = new BoundedKnapsack(this._maxValue, this._mod, this._dp.slice())
    res._maxJ = this._maxJ
    return res
  }
}

export { boundedKnapsackDp, boundedKnapsackDpBinary, boundedKnapsackDpNaive, boundedKnapsackDpCountWays, BoundedKnapsack }

if (require.main === module) {
  //   4 20
  // 3 9 3
  // 5 9 1
  // 9 4 2
  // 8 1 3
  const n = 4
  const maxCapacity = 2e7
  const values = [3, 5, 9, 80]
  const weights = [9, 9, 4, 1]
  const counts = [3, 1, 2, 3]
  console.time('boundedKnapsackDp')
  console.log(boundedKnapsackDp(values, weights, counts, maxCapacity))
  console.timeEnd('boundedKnapsackDp')

  function waysToReachTarget(target: number, types: number[][]): number {
    // const counts = types.map(t => t[0])
    // const values = types.map(t => t[1])
    // return boundedKnapsackDpCountWays(values, counts, target)[target]
    const K = new BoundedKnapsack(target, 1e9 + 7)
    types.forEach(v => K.add(v[1], v[0]))
    return K.query(target)
  }

  console.log(
    waysToReachTarget(6, [
      [6, 1],
      [3, 2],
      [2, 3]
    ])
  )
}
