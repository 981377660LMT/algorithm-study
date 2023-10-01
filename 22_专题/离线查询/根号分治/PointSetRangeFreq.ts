/**
 * 单点修改，区间频率查询.
 * 单次修改复杂度 `O(sqrt(n))`, 单次查询复杂度 `O(sqrt(n) * log(sqrt(n)))`.
 */
class PointSetRangeFreq {
  private readonly _nums: number[]
  private readonly _blockStart: (bid: number) => number
  private readonly _blockEnd: (bid: number) => number
  private readonly _belong: (index: number) => number
  private readonly _blockSorted: number[][]

  constructor(arr: ArrayLike<number>, blockSize = 2 * ((Math.sqrt(arr.length) + 1) | 0)) {
    const n = arr.length
    const copy: number[] = Array(n)
    for (let i = 0; i < n; i++) copy[i] = arr[i]

    const blockCount = 1 + ((n / blockSize) | 0)
    const blockStart = (bid: number) => bid * blockSize
    const blockEnd = (bid: number) => Math.min((bid + 1) * blockSize, n)
    const belong = (index: number) => (index / blockSize) | 0

    const blockSorted: number[][] = Array(blockCount)
    for (let bid = 0; bid < blockCount; bid++) {
      const curSorted = copy.slice(blockStart(bid), blockEnd(bid)).sort((a, b) => a - b)
      blockSorted[bid] = curSorted
    }

    this._nums = copy
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._belong = belong
    this._blockSorted = blockSorted
  }

  /**
   * 修改下标 `pos` 的值为 `newValue`.
   */
  set(pos: number, newValue: number): void {
    if (this._nums[pos] === newValue) return
    const pre = this._nums[pos]
    this._nums[pos] = newValue
    const bid = this._belong(pos)
    const removeIndex = this._bisectRight(this._blockSorted[bid], pre) - 1
    this._blockSorted[bid].splice(removeIndex, 1)
    const insertIndex = this._bisectRight(this._blockSorted[bid], newValue)
    this._blockSorted[bid].splice(insertIndex, 0, newValue)
  }

  /**
   * 统计 `[start, end)` 中等于 `target` 的元素个数.
   */
  count(start: number, end: number, target: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong(start)
    const bid2 = this._belong(end - 1)
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] === target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd(bid1); i++) res += +(this._nums[i] === target)
    for (let bid = bid1 + 1; bid < bid2; bid++) res += this._count(this._blockSorted[bid], target)
    for (let i = this._blockStart(bid2); i < end; i++) res += +(this._nums[i] === target)
    return res
  }

  /**
   * 统计 `[start, end)` 中严格小于 `target` 的元素个数.
   */
  countLower(start: number, end: number, target: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong(start)
    const bid2 = this._belong(end - 1)
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] < target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd(bid1); i++) res += +(this._nums[i] < target)
    for (let bid = bid1 + 1; bid < bid2; bid++) {
      res += this._bisectLeft(this._blockSorted[bid], target)
    }
    for (let i = this._blockStart(bid2); i < end; i++) res += +(this._nums[i] < target)
    return res
  }

  /**
   * 统计 `[start, end)` 中小于等于 `target` 的元素个数.
   */
  countFloor(start: number, end: number, target: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong(start)
    const bid2 = this._belong(end - 1)
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] <= target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd(bid1); i++) res += +(this._nums[i] <= target)
    for (let bid = bid1 + 1; bid < bid2; bid++) {
      res += this._bisectRight(this._blockSorted[bid], target)
    }
    for (let i = this._blockStart(bid2); i < end; i++) res += +(this._nums[i] <= target)
    return res
  }

  /**
   * 统计 `[start, end)` 中大于等于 `target` 的元素个数.
   */
  countCeiling(start: number, end: number, target: number): number {
    return end - start - this.countLower(start, end, target)
  }

  /**
   * 统计 `[start, end)` 中严格大于 `target` 的元素个数.
   */
  countHigher(start: number, end: number, target: number): number {
    return end - start - this.countFloor(start, end, target)
  }

  // eslint-disable-next-line class-methods-use-this
  private _bisectLeft(nums: number[], target: number, left = 0, right = nums.length - 1): number {
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (nums[mid] >= target) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }

  // eslint-disable-next-line class-methods-use-this
  private _bisectRight(nums: number[], target: number, left = 0, right = nums.length - 1): number {
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (nums[mid] > target) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }

  private _count(nums: number[], target: number, left = 0, right = nums.length - 1): number {
    return (
      this._bisectRight(nums, target, left, right) - this._bisectLeft(nums, target, left, right)
    )
  }
}

export { PointSetRangeFreq }

if (require.main === module) {
  // https://leetcode.cn/problems/range-frequency-queries/description/
  class RangeFreqQuery {
    private readonly _ps: PointSetRangeFreq
    constructor(arr: number[]) {
      this._ps = new PointSetRangeFreq(arr)
    }

    query(left: number, right: number, value: number): number {
      return this._ps.count(left, right + 1, value)
    }
  }
}
