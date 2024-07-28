/* eslint-disable curly */
/* eslint-disable nonblock-statement-body-position */
/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

// 动态区间频率查询
// 1.区间加，查询区间某元素出现次数(RangeAddRangeFreq)
// 2.区间赋值，查询区间某元素出现次数(RangeAssignRangeFreq)

import assert from 'assert'

import { useBlock } from './SqrtDecomposition/useBlock'

const INF = 2e9 // !超过int32使用2e15

const TYPEARRAY_RECORD = {
  int8: Int8Array,
  uint8: Uint8Array,
  int16: Int16Array,
  uint16: Uint16Array,
  int32: Int32Array,
  uint32: Uint32Array,
  float32: Float32Array,
  float64: Float64Array
}

type ArrayType = InstanceType<(typeof TYPEARRAY_RECORD)[keyof typeof TYPEARRAY_RECORD]>

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
   * 用于指定`arr`的类型.
   * !如果数组中元素始终为int32，可指定为`int32`来加速内部排序.
   * 默认为`float64`.
   */
  constructor(arr: ArrayLike<number>, type: keyof typeof TYPEARRAY_RECORD = 'float64') {
    const n = arr.length
    const blockSize = type === 'float64' ? (Math.sqrt(n) + 1) | 0 : (1.2 * Math.sqrt(n) + 1) | 0
    const { belong, blockStart, blockEnd, blockCount } = useBlock(n, blockSize)
    const ArrayConstructor = TYPEARRAY_RECORD[type]
    this._nums = new ArrayConstructor(arr)
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._blockLazy = Array(blockCount).fill(0)
    this._blockSorted = new ArrayConstructor(n)
    for (let i = 0; i < blockCount; i++) this._rebuild(i)
  }

  get(index: number): number {
    return this._nums[index] + this._blockLazy[this._belong[index]]
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._nums.length) return
    const bid = this._belong[index]
    const target = value - this._blockLazy[bid]
    if (this._nums[index] === target) return
    this._nums[index] = target
    this._rebuild(bid)
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
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this._nums[i] + this._blockLazy[bid1] === target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd[bid1]; i++)
      res += +(this._nums[i] + this._blockLazy[bid1] === target)
    for (let bid = bid1 + 1; bid < bid2; bid++) {
      res += RangeAddRangeFreq._count(
        this._blockSorted,
        target - this._blockLazy[bid],
        this._blockStart[bid],
        this._blockEnd[bid] - 1
      )
    }
    for (let i = this._blockStart[bid2]; i < end; i++)
      res += +(this._nums[i] + this._blockLazy[bid2] === target)
    return res
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

  toString(): string {
    const res = Array(this._nums.length)
    for (let i = 0; i < this._nums.length; i++) {
      res[i] = this.get(i)
    }
    return `RangeAddRangeFreq{${res.join(',')}}`
  }

  private _rebuild(bid: number): void {
    this._blockSorted.set(
      this._nums.slice(this._blockStart[bid], this._blockEnd[bid]).sort(),
      this._blockStart[bid]
    )
  }

  private static _bisectLeft(
    arr: ArrayLike<number>,
    value: number,
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

  private static _bisectRight(
    arr: ArrayLike<number>,
    value: number,
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

  private static _count(
    arr: ArrayLike<number>,
    value: number,
    left = 0,
    right = arr.length - 1
  ): number {
    return this._bisectRight(arr, value, left, right) - this._bisectLeft(arr, value, left, right)
  }
}

/**
 * 区间赋值，区间频率查询.
 * 单次修改、查询时间复杂度`O(sqrt(n))`，空间复杂度`O(n)`.
 */
class RangeAssignRangeFreq {
  private readonly _color: ArrayType
  private readonly _blockSorted: ArrayType
  private readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _blockColor: ArrayType // 整块赋值后的颜色.INF表示未赋值过.
  private readonly _updateTime: Int32Array // 单点赋值时间戳.-1表示未赋值过.
  private readonly _blockUpdateTime: Int32Array // 整块赋值时间戳.-1表示未赋值过.
  private _time = 0 // 每次赋值操作时间戳+1

  /**
   * @param arr 初始数组.
   * @param arrayType 初始数组类型.
   * 用于指定`arr`的类型.
   * !如果数组中元素始终为int32，可指定为`int32`来加速内部排序.
   * 默认为`float64`.
   */
  constructor(arr: ArrayLike<number>, type: keyof typeof TYPEARRAY_RECORD = 'float64') {
    const n = arr.length
    const blockSize =
      type === 'float64' ? (0.75 * Math.sqrt(n) + 1) | 0 : (1.2 * Math.sqrt(n) + 1) | 0
    const { belong, blockStart, blockEnd, blockCount } = useBlock(n, blockSize)
    const ArrayConstructor = TYPEARRAY_RECORD[type]
    this._color = new ArrayConstructor(arr)
    this._blockSorted = new ArrayConstructor(arr)
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._blockColor = new ArrayConstructor(blockCount).fill(INF)
    this._updateTime = new Int32Array(n).fill(-1)
    this._blockUpdateTime = new Int32Array(blockCount).fill(-1)
    for (let i = 0; i < blockCount; i++) this._rebuild(i)
  }

  get(index: number): number {
    const bid = this._belong[index]
    return this._blockUpdateTime[bid] > this._updateTime[index]
      ? this._blockColor[bid]
      : this._color[index]
  }

  set(index: number, value: number): void {
    this.assign(index, index + 1, value)
  }

  /**
   * 区间`[start, end)`每个元素赋值为`value`.
   */
  assign(start: number, end: number, value: number): void {
    if (start < 0) start = 0
    if (end > this._color.length) end = this._color.length
    if (start >= end) return
    const time = this._time++
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) {
        this._color[i] = value
        this._updateTime[i] = time
      }
      this._rebuild(bid1)
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) {
        this._color[i] = value
        this._updateTime[i] = time
      }
      this._rebuild(bid1)
      for (let bid = bid1 + 1; bid < bid2; bid++) {
        this._blockColor[bid] = value
        this._blockUpdateTime[bid] = time
      }
      for (let i = this._blockStart[bid2]; i < end; i++) {
        this._color[i] = value
        this._updateTime[i] = time
      }
      this._rebuild(bid2)
    }
  }

  /**
   * 查询区间`[start, end)`中元素`target`出现的次数.
   */
  rangeFreq(start: number, end: number, target: number): number {
    if (start < 0) start = 0
    if (end > this._color.length) end = this._color.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      let res = 0
      for (let i = start; i < end; i++) res += +(this.get(i) === target)
      return res
    }
    let res = 0
    for (let i = start; i < this._blockEnd[bid1]; i++) res += +(this.get(i) === target)
    for (let bid = bid1 + 1; bid < bid2; bid++) {
      if (~this._blockUpdateTime[bid]) {
        res += this._blockColor[bid] === target ? this._blockEnd[bid] - this._blockStart[bid] : 0
      } else {
        res += RangeAssignRangeFreq._count(
          this._blockSorted,
          target,
          this._blockStart[bid],
          this._blockEnd[bid] - 1
        )
      }
    }
    for (let i = this._blockStart[bid2]; i < end; i++) res += +(this.get(i) === target)
    return res
  }

  /**
   * 查询区间`[start, end)`中严格大于`lower`的元素出现的次数.
   */
  rangeFreqLower(start: number, end: number, lower: number): number {
    if (start < 0) start = 0
    if (end > this._color.length) end = this._color.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    let res = 0
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) {
        res += +(this.get(i) > lower)
      }
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) {
        res += +(this.get(i) > lower)
      }
      for (let bid = bid1 + 1; bid < bid2; bid++) {
        if (~this._blockUpdateTime[bid]) {
          res += this._blockColor[bid] > lower ? this._blockEnd[bid] - this._blockStart[bid] : 0
        } else {
          const bEnd = this._blockEnd[bid]
          res +=
            bEnd -
            RangeAssignRangeFreq._bisectRight(
              this._blockSorted,
              lower,
              this._blockStart[bid],
              bEnd - 1
            )
        }
      }
      for (let i = this._blockStart[bid2]; i < end; i++) {
        res += +(this.get(i) > lower)
      }
    }
    return res
  }

  /**
   * 查询区间`[start, end)`中大于等于`floor`的元素出现的次数.
   */
  rangeFreqFloor(start: number, end: number, floor: number): number {
    if (start < 0) start = 0
    if (end > this._color.length) end = this._color.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    let res = 0
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) {
        res += +(this.get(i) >= floor)
      }
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) {
        res += +(this.get(i) >= floor)
      }
      for (let bid = bid1 + 1; bid < bid2; bid++) {
        if (~this._blockUpdateTime[bid]) {
          res += this._blockColor[bid] >= floor ? this._blockEnd[bid] - this._blockStart[bid] : 0
        } else {
          const bEnd = this._blockEnd[bid]
          res +=
            bEnd -
            RangeAssignRangeFreq._bisectLeft(
              this._blockSorted,
              floor,
              this._blockStart[bid],
              bEnd - 1
            )
        }
      }
      for (let i = this._blockStart[bid2]; i < end; i++) {
        res += +(this.get(i) >= floor)
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

  toString(): string {
    const res = Array(this._color.length)
    for (let i = 0; i < this._color.length; i++) {
      res[i] = this.get(i)
    }
    return `RangeAssignRangeFreq{${res.join(',')}}`
  }

  private _rebuild(bid: number): void {
    for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
      this._color[i] = this.get(i)
    }
    this._blockColor[bid] = INF
    this._blockUpdateTime[bid] = -1
    this._blockSorted.set(
      this._color.slice(this._blockStart[bid], this._blockEnd[bid]).sort(),
      this._blockStart[bid]
    )
  }

  private static _bisectLeft(
    arr: ArrayLike<number>,
    value: number,
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

  private static _bisectRight(
    arr: ArrayLike<number>,
    value: number,
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

  private static _count(
    arr: ArrayLike<number>,
    value: number,
    left = 0,
    right = arr.length - 1
  ): number {
    return this._bisectRight(arr, value, left, right) - this._bisectLeft(arr, value, left, right)
  }
}

export { RangeAddRangeFreq, RangeAssignRangeFreq }

if (require.main === module) {
  // 2080. 区间内查询数字的频率
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

  class Mocker {
    readonly nums: number[]

    constructor(num: number[]) {
      this.nums = num.slice()
    }

    set(index: number, value: number): void {
      this.nums[index] = value
    }

    add(start: number, end: number, delta: number): void {
      for (let i = start; i < end; i++) this.nums[i] += delta
    }

    assign(start: number, end: number, value: number): void {
      for (let i = start; i < end; i++) this.nums[i] = value
    }

    rangeFreq(start: number, end: number, target: number): number {
      let res = 0
      for (let i = start; i < end; i++) res += +(this.nums[i] === target)
      return res
    }

    rangeFreqLower(start: number, end: number, lower: number): number {
      let res = 0
      for (let i = start; i < end; i++) res += +(this.nums[i] > lower)
      return res
    }

    rangeFreqFloor(start: number, end: number, floor: number): number {
      let res = 0
      for (let i = start; i < end; i++) res += +(this.nums[i] >= floor)
      return res
    }

    rangeFreqCeiling(start: number, end: number, ceiling: number): number {
      let res = 0
      for (let i = start; i < end; i++) res += +(this.nums[i] <= ceiling)
      return res
    }

    rangeFreqUpper(start: number, end: number, upper: number): number {
      let res = 0
      for (let i = start; i < end; i++) res += +(this.nums[i] < upper)
      return res
    }
  }

  testTime1()
  testTime2()

  testRangeAddRangeFreq()
  testRangeAssignRangeFreq()
  testRangePointSetRangeFreq()

  function testTime1(): void {
    const N = 1e5
    const arr = Array.from({ length: N }, (_, i) => i)
    const rf = new RangeAddRangeFreq(arr, 'float64')
    console.time('time1')
    for (let i = 0; i < N; i++) {
      rf.add(0, N, i)
      rf.rangeFreq(0, N, i)
    }
    console.timeEnd('time1') // time1: 2.1s
  }

  function testTime2(): void {
    const N = 1e5
    const arr = Array.from({ length: N }, (_, i) => i)
    const rf = new RangeAssignRangeFreq(arr, 'float64')
    console.time('time1')
    for (let i = 0; i < N; i++) {
      rf.assign(i, N, i)
      rf.rangeFreq(0, N, i)
    }
    console.timeEnd('time1') // time1: 2.1s
  }

  function testRangeAddRangeFreq(): void {
    const N = 2e4
    const arr = Array.from({ length: N }, (_, i) => i)
    const real = new RangeAddRangeFreq(arr)
    const mock = new Mocker(arr)
    const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
    for (let i = 0; i < 2e4; i++) {
      const op = randint(0, 6)
      if (op === 0) {
        // set
        const index = randint(0, N - 1)
        const value = randint(-10, 10)
        real.set(index, value)
        mock.set(index, value)
      } else if (op === 1) {
        // add
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const delta = randint(-5, 5)
        real.add(start, end, delta)
        mock.add(start, end, delta)
      } else if (op === 2) {
        // rangeFreqFloor
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const floor = randint(-10, 10)
        const res = real.rangeFreqFloor(start, end, floor)
        const ans = mock.rangeFreqFloor(start, end, floor)
        assert.strictEqual(res, ans)
      } else if (op === 3) {
        // rangeFreq
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const target = randint(-10, 10)
        const res = real.rangeFreq(start, end, target)
        const ans = mock.rangeFreq(start, end, target)
        assert.strictEqual(res, ans)
      } else if (op === 4) {
        // rangeFreqLower
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const lower = randint(-10, 10)
        const res = real.rangeFreqLower(start, end, lower)
        const ans = mock.rangeFreqLower(start, end, lower)
        assert.strictEqual(res, ans)
      } else if (op === 5) {
        // rangeFreqCeiling
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const ceiling = randint(-10, 10)
        const res = real.rangeFreqCeiling(start, end, ceiling)
        const ans = mock.rangeFreqCeiling(start, end, ceiling)
        assert.strictEqual(res, ans)
      } else if (op === 6) {
        // rangeFreqUpper
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const upper = randint(-10, 10)
        const res = real.rangeFreqUpper(start, end, upper)
        const ans = mock.rangeFreqUpper(start, end, upper)
        assert.strictEqual(res, ans)
      }
    }

    console.log('test 1 done')
  }

  function testRangeAssignRangeFreq(): void {
    const N = 2e4
    const arr = Array.from({ length: N }, (_, i) => i)
    const real = new RangeAssignRangeFreq(arr)
    const mock = new Mocker(arr)
    const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
    for (let i = 0; i < 2e4; i++) {
      const op = randint(0, 6)
      if (op === 0) {
        // set
        const index = randint(0, N - 1)
        const value = randint(-100, 100)
        real.set(index, value)
        mock.set(index, value)
      } else if (op === 1) {
        // assign
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const delta = randint(-10, 10)
        real.assign(start, end, delta)
        mock.assign(start, end, delta)
      } else if (op === 2) {
        // rangeFreqFloor
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const floor = randint(-100, 100)
        const res = real.rangeFreqFloor(start, end, floor)
        const ans = mock.rangeFreqFloor(start, end, floor)
        assert.strictEqual(res, ans)
      } else if (op === 3) {
        // rangeFreq
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const target = randint(-100, 100)
        const res = real.rangeFreq(start, end, target)
        const ans = mock.rangeFreq(start, end, target)
        assert.strictEqual(res, ans)
      } else if (op === 4) {
        // rangeFreqLower
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const lower = randint(-100, 100)
        const res = real.rangeFreqLower(start, end, lower)
        const ans = mock.rangeFreqLower(start, end, lower)
        assert.strictEqual(res, ans)
      } else if (op === 5) {
        // rangeFreqCeiling
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const ceiling = randint(-100, 100)
        const res = real.rangeFreqCeiling(start, end, ceiling)
        const ans = mock.rangeFreqCeiling(start, end, ceiling)
        assert.strictEqual(res, ans)
      } else if (op === 6) {
        // rangeFreqUpper
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const upper = randint(-100, 100)
        const res = real.rangeFreqUpper(start, end, upper)
        const ans = mock.rangeFreqUpper(start, end, upper)
        assert.strictEqual(res, ans)
      }
    }

    console.log('test 2 done')
  }

  function testRangePointSetRangeFreq(): void {
    const N = 2e4
    const arr = Array.from({ length: N }, (_, i) => i)
    const real = new RangeAssignRangeFreq(arr)
    const mock = new Mocker(arr)
    const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
    for (let i = 0; i < 2e4; i++) {
      const op = randint(0, 6)
      if (op === 0) {
        // set
        const index = randint(0, N - 1)
        const value = randint(-100, 100)
        real.set(index, value)
        mock.set(index, value)
      } else if (op === 1) {
        // assign
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const delta = randint(-10, 10)
        real.assign(start, end, delta)
        mock.assign(start, end, delta)
      } else if (op === 2) {
        // rangeFreqFloor
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const floor = randint(-100, 100)
        const res = real.rangeFreqFloor(start, end, floor)
        const ans = mock.rangeFreqFloor(start, end, floor)
        assert.strictEqual(res, ans)
      } else if (op === 3) {
        // rangeFreq
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const target = randint(-100, 100)
        const res = real.rangeFreq(start, end, target)
        const ans = mock.rangeFreq(start, end, target)
        assert.strictEqual(res, ans)
      } else if (op === 4) {
        // rangeFreqLower
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const lower = randint(-100, 100)
        const res = real.rangeFreqLower(start, end, lower)
        const ans = mock.rangeFreqLower(start, end, lower)
        assert.strictEqual(res, ans)
      } else if (op === 5) {
        // rangeFreqCeiling
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const ceiling = randint(-100, 100)
        const res = real.rangeFreqCeiling(start, end, ceiling)
        const ans = mock.rangeFreqCeiling(start, end, ceiling)
        assert.strictEqual(res, ans)
      } else if (op === 6) {
        // rangeFreqUpper
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const upper = randint(-100, 100)
        const res = real.rangeFreqUpper(start, end, upper)
        const ans = mock.rangeFreqUpper(start, end, upper)
        assert.strictEqual(res, ans)
      }
    }

    console.log('test 3 done')
  }

  /**
   * Your RangeFreqQuery object will be instantiated and called as such:
   * var obj = new RangeFreqQuery(arr)
   * var param_1 = obj.query(left,right,value)
   */
}
