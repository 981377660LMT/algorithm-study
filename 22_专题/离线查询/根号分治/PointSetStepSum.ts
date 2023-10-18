/* eslint-disable no-inner-declarations */

// P3396 哈希冲突
// https://www.luogu.com.cn/problem/P3396
// https://www.luogu.com.cn/blog/danieljiang/ha-xi-chong-tu-ti-xie-gen-hao-ke-ji

class PointSetStepSum {
  private readonly _nums: Float64Array
  private readonly _stepThreshold: number

  /**
   * dp[step][start] 表示步长为step,起点为start的所有元素的和.
   * `dp[step][start] = dp[step][start+step] + nums[start]`.
   */
  private readonly _dp: Float64Array[]

  /**
   * @param stepThreshold 步长阈值,当步长小于等于该值时,使用dp数组预处理答案,否则直接遍历.
   * 预处理时间空间复杂度均为`O(n*stepThreshold)`.
   * 单次遍历时间复杂度为`O(n/stepThreshold)`.
   */
  constructor(arr: ArrayLike<number>, stepThreshold = 80) {
    this._nums = new Float64Array(arr)
    this._stepThreshold = stepThreshold
    this._dp = Array(stepThreshold)
    const n = arr.length
    for (let step = 1; step <= stepThreshold; step++) {
      const curSum = new Float64Array(n + 1)
      for (let start = n - 1; ~start; start--) {
        curSum[start] = curSum[Math.min(n, start + step)] + arr[start]
      }
      this._dp[step - 1] = curSum
    }
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._nums.length) return
    const pre = this._nums[index]
    if (pre === value) return
    this._nums[index] = value
    const delta = value - pre
    for (let step = 1; step <= this._stepThreshold; step++) {
      this._dp[step - 1][index % step] += delta
    }
  }

  query(start: number, step: number): number {
    if (start < 0) start = 0
    if (step <= this._stepThreshold) {
      return this._dp[step - 1][start]
    }
    let sum = 0
    for (let i = start; i < this._nums.length; i += step) {
      sum += this._nums[i]
    }
    return sum
  }

  toString(): string {
    return `PointSetStepSum{${this._nums.join(',')}}`
  }
}

export { PointSetStepSum }

if (require.main === module) {
  // 1714. 数组中特殊等间距元素的和
  // https://leetcode.cn/problems/sum-of-special-evenly-spaced-elements-in-array/
  function solve(nums: number[], queries: number[][]): number[] {
    const MOD = 1e9 + 7
    const R = new PointSetStepSum(nums)
    const res: number[] = Array(queries.length)
    for (let i = 0; i < queries.length; i++) {
      const { 0: start, 1: step } = queries[i]
      res[i] = R.query(start, step) % MOD
    }
    return res
  }
}
