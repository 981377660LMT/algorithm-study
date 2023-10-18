/* eslint-disable no-inner-declarations */

class RangeStepSum {
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

  query(start: number, stop: number, step: number): number {
    if (start < 0) start = 0
    if (stop > this._nums.length) stop = this._nums.length
    if (start >= stop) return 0
    if (step <= this._stepThreshold) {
      const curDp = this._dp[step - 1]
      // 找到 >=stop 的第一个形为start+k*step的数
      const div = Math.ceil((stop - start) / step)
      const nextStart = Math.min(start + div * step, this._nums.length)
      return curDp[start] - curDp[nextStart]
    }
    let sum = 0
    for (let i = start; i < stop; i += step) {
      sum += this._nums[i]
    }
    return sum
  }

  toString(): string {
    return `RangeStepSum{${this._nums.join(',')}}`
  }
}

export { RangeStepSum }

if (require.main === module) {
  // 1714. 数组中特殊等间距元素的和
  // https://leetcode.cn/problems/sum-of-special-evenly-spaced-elements-in-array/
  function solve(nums: number[], queries: number[][]): number[] {
    const MOD = 1e9 + 7
    const n = nums.length
    const R = new RangeStepSum(nums)
    const res: number[] = Array(queries.length)
    for (let i = 0; i < queries.length; i++) {
      const { 0: start, 1: step } = queries[i]
      const stop = n
      res[i] = R.query(start, stop, step) % MOD
    }
    return res
  }

  testTime()
  test()
  function testTime(): void {
    const n = 2e5
    const arr = Array(n)
      .fill(0)
      .map(() => (Math.random() * 5) | 0)
    console.time('pointSetRangeStepSum')
    const pointSetRangeStepSum = new RangeStepSum(arr)
    for (let i = 0; i < n; i++) {
      pointSetRangeStepSum.query(0, n, i + 1)
    }
    console.timeEnd('pointSetRangeStepSum')
  }

  function test(): void {
    class Mocker {
      private readonly _arr: number[]
      constructor(arr: number[]) {
        this._arr = arr.slice()
      }
      set(index: number, value: number): void {
        this._arr[index] = value
      }
      query(start: number, stop: number, step: number): number {
        let sum = 0
        for (let i = start; i < stop; i += step) {
          sum += this._arr[i]
        }
        return sum
      }
      toString(): string {
        return this._arr.toString()
      }
    }

    const n = 1e4
    const arr = Array(n)
      .fill(0)
      .map(() => (Math.random() * 5) | 0)
    const mocker = new Mocker(arr)
    const pointSetRangeStepSum = new RangeStepSum(arr)
    for (let i = 0; i < 1e5; i++) {
      const op = Math.random() < 0 ? 'set' : 'query'
      const start = (Math.random() * n) | 0
      const stop = (Math.random() * n) | 0
      const step = (1 + Math.random() * n) | 0
      if (op === 'set') {
        // const value = (Math.random() * 5) | 0
        // mocker.set(start, value)
        // pointSetRangeStepSum.set(start, value)
      } else {
        const res1 = mocker.query(start, stop, step)
        const res2 = pointSetRangeStepSum.query(start, stop, step)
        if (res1 !== res2) {
          console.error('error', res1, res2, start, stop, step)
          console.log("mocker's arr", mocker.toString())
          console.log('pointSetRangeStepSum', pointSetRangeStepSum.toString())
          throw new Error()
        }
      }
    }

    console.log('success')
  }
}
