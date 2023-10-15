// 动态区间频率查询
// 区间加，查询区间某元素出现次数
// RangeAddRangeFreq

import assert from 'assert'

import { useBlock } from './SqrtDecomposition/useBlock'

const INF = 2e9 // !超过int32使用2e15

const ARRAYTYPE_RECORD = {
  int8: Int8Array,
  uint8: Uint8Array,
  int16: Int16Array,
  uint16: Uint16Array,
  int32: Int32Array,
  uint32: Uint32Array,
  float32: Float32Array,
  float64: Float64Array
}

type ArrayType = InstanceType<(typeof ARRAYTYPE_RECORD)[keyof typeof ARRAYTYPE_RECORD]>

/**
 * 区间加，区间频率查询.
 * 单次修改、查询时间复杂度`O(sqrt(n)logn)`，空间复杂度`O(n)`.
 */
class RangeAddRangeFreq {
  private readonly _nums: ArrayType
  private readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _blockLazy: number[]
  private readonly _blockSorted: ArrayType

  /**
   * @param arr 初始数组.
   * @param arrayType 初始数组类型.
   * 用于指定`arr`的类型.例如，如果始终为int32，可指定为`int32`来加速内部排序.
   * 默认为`float64`.
   */
  constructor(arr: ArrayLike<number>, type: keyof typeof ARRAYTYPE_RECORD = 'float64') {
    const n = arr.length
    const blockSize = type === 'float64' ? Math.sqrt(n + 1) | 0 : 2 * (Math.sqrt(n + 1) | 0)
    const { belong, blockStart, blockEnd, blockCount } = useBlock(n, blockSize)
    this._nums = new ARRAYTYPE_RECORD[type](arr)
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._blockLazy = Array(blockCount).fill(0)
    this._blockSorted = new ARRAYTYPE_RECORD[type](n)
    for (let i = 0; i < blockCount; i++) this._rebuild(i)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._nums.length || this._nums[index] === value) return
    this._nums[index] = value
    this._rebuild(this._belong[index])
  }

  /**
   * 区间`[start, end)`每个元素加上`delta`.
   */
  add(start: number, end: number, delta: number): void {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) this._nums[i] += delta
      this._rebuild(bid1)
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) this._nums[i] += delta
      this._rebuild(bid1)
      for (let bid = bid1 + 1; bid < bid2; bid++) this._blockLazy[bid] += delta
      for (let i = this._blockStart[bid2]; i < end; i++) this._nums[i] += delta
      this._rebuild(bid2)
    }
  }

  /**
   * 查询区间`[start, end)`中元素`target`出现的次数.
   */
  rangeFreq(start: number, end: number, target: number): number {
    return this.rangeFreqLower(start, end, target - 1) - this.rangeFreqLower(start, end, target)
  }

  /**
   * 查询区间`[start, end)`中严格大于`lower`的元素出现的次数.
   */
  rangeFreqLower(start: number, end: number, lower: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    let res = 0
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) {
        res += +(this._nums[i] + this._blockLazy[bid1] > lower)
      }
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) {
        res += +(this._nums[i] + this._blockLazy[bid1] > lower)
      }
      for (let bid = bid1 + 1; bid < bid2; bid++) {
        const bEnd = this._blockEnd[bid]
        res +=
          bEnd -
          RangeAddRangeFreq._bisectRight(
            this._blockSorted,
            lower - this._blockLazy[bid],
            this._blockStart[bid],
            bEnd - 1
          )
      }
      for (let i = this._blockStart[bid2]; i < end; i++) {
        res += +(this._nums[i] + this._blockLazy[bid2] > lower)
      }
    }
    return res
  }

  /**
   * 查询区间`[start, end)`中大于等于`floor`的元素出现的次数.
   */
  rangeFreqFloor(start: number, end: number, floor: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    let res = 0
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) {
        res += +(this._nums[i] + this._blockLazy[bid1] >= floor)
      }
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) {
        res += +(this._nums[i] + this._blockLazy[bid1] >= floor)
      }
      for (let bid = bid1 + 1; bid < bid2; bid++) {
        const bEnd = this._blockEnd[bid]
        res +=
          bEnd -
          RangeAddRangeFreq._bisectLeft(
            this._blockSorted,
            floor - this._blockLazy[bid],
            this._blockStart[bid],
            bEnd - 1
          )
      }
      for (let i = this._blockStart[bid2]; i < end; i++) {
        res += +(this._nums[i] + this._blockLazy[bid2] >= floor)
      }
    }
    return res
  }

  /**
   * 查询区间`[start, end)`中小于等于`ceiling`的元素出现的次数.
   */
  rangeFreqCeiling(start: number, end: number, ceiling: number): number {
    return end - start - this.rangeFreqLower(start, end, ceiling)
  }

  /**
   * 查询区间`[start, end)`中严格小于`upper`的元素出现的次数.
   */
  rangeFreqUpper(start: number, end: number, upper: number): number {
    return end - start - this.rangeFreqFloor(start, end, upper)
  }

  private _rebuild(bid: number): void {
    this._blockSorted.set(
      this._nums.slice(this._blockStart[bid], this._blockEnd[bid]).sort(),
      this._blockStart[bid]
    )
  }

  private static _bisectLeft<T>(
    arr: ArrayLike<T>,
    value: T,
    left = 0,
    right = arr.length - 1
  ): number {
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] < value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private static _bisectRight<T>(
    arr: ArrayLike<T>,
    value: T,
    left = 0,
    right = arr.length - 1
  ): number {
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] <= value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }
}

/**
 * 区间赋值，区间频率查询.
 * 单次修改时间复杂度 `O(sqrt(n))`，单次查询时间复杂度`O(sqrt(n)logn)`，空间复杂度`O(n)`.
 */
class RangeAssignRangeFreq {}

export { RangeAddRangeFreq, RangeAssignRangeFreq }

if (require.main === module) {
  let rf = new RangeAddRangeFreq([1, 2, 2, 4, 5, 6, 7, 8, 9, 10])
  rf.add(0, 10, 1)
  assert.strictEqual(rf.rangeFreq(0, 10, 5), 1)
  rf.add(0, 10, 2)
  assert.strictEqual(rf.rangeFreq(0, 10, 5), 2)
  console.log(rf.rangeFreqLower(0, 10, 5))
  assert.strictEqual(rf.rangeFreqFloor(0, 10, 5), 9)

  const N = 1e5
  const arr = Array.from({ length: N }, (_, i) => i)
  rf = new RangeAddRangeFreq(arr, 'int32')
  console.time('time1')
  for (let i = 0; i < N; i++) {
    rf.add(0, N, i)
    rf.rangeFreq(0, N, i)
  }
  console.timeEnd('time1') // time1: 2.1s

  // https://leetcode.cn/problems/range-frequency-queries/
  class RangeFreqQuery {
    private readonly _rf: RangeAddRangeFreq

    constructor(arr: number[]) {
      this._rf = new RangeAddRangeFreq(arr)
    }

    query(left: number, right: number, value: number): number {
      return this._rf.rangeFreq(left, right + 1, value)
    }
  }

  /**
   * Your RangeFreqQuery object will be instantiated and called as such:
   * var obj = new RangeFreqQuery(arr)
   * var param_1 = obj.query(left,right,value)
   */
}
