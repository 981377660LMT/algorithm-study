// https://usaco.guide/plat/sqrt?lang=cpp#batching
//
// Maintain a "buffer" of the latest updates (up to $\sqrt N$).
// The answer for each sum query can be calculated with prefix sums
// and by examining each update within the buffer.
// When the buffer gets too large ($\ge \sqrt N$), clear it and recalculate prefix sums.

/**
 * 基于批处理的`PointSetRangeSum`实现.适合查询多、更新少的场景.
 */
class PointSetRangeSumBatching {
  private readonly _nums: number[]
  private readonly _preSum: Float64Array
  private readonly _updates: { index: number; delta: number }[] = []
  private readonly _rebuildThreshold: number

  /**
   * @param rebuildThreshold 当更新次数超过此阈值时，重建前缀和数组.默认为1000.
   */
  constructor(nums: number[], rebuildThreshold = 1000) {
    this._nums = nums.slice()
    this._preSum = new Float64Array(nums.length + 1)
    this._rebuildThreshold = rebuildThreshold
    this._build()
  }

  set(index: number, value: number): void {
    this._updates.push({ index, delta: value - this._nums[index] })
    this._nums[index] = value
  }

  add(index: number, delta: number): void {
    this.set(index, this._nums[index] + delta)
  }

  query(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    let res = this._preSum[end] - this._preSum[start]
    for (let i = 0; i < this._updates.length; i++) {
      if (this._updates[i].index >= start && this._updates[i].index < end) {
        res += this._updates[i].delta
      }
    }
    if (this._updates.length >= this._rebuildThreshold) {
      this._updates.length = 0
      this._build()
    }
    return res
  }

  private _build(): void {
    this._preSum[0] = 0
    for (let i = 0; i < this._nums.length; i++) {
      this._preSum[i + 1] = this._preSum[i] + this._nums[i]
    }
  }
}

export { PointSetRangeSumBatching }

if (require.main === module) {
  const sum = new PointSetRangeSumBatching([1, 2, 3, 4, 5])
  console.log(sum.query(0, 5))
  console.log(sum.query(0, 3))
  console.log(sum.query(1, 4))
  console.log(sum.query(2, 5))
}
