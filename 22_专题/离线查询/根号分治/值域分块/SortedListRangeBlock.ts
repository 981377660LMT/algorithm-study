/* eslint-disable generator-star-spacing */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// SortedListRangeBlock

import { ISortedList, ISortedListIterator } from '../SortedList/SortedListFast'

/**
 * 值域分块模拟SortedList.
 * `O(1)`add/remove，`O(sqrt(n))`查询.
 * 一般配合莫队算法使用.
 */
class SortedListRangeBlock implements ISortedList<number> {
  /** 每个块的大小. */
  private readonly _blockSize: number

  /** 每个数出现的次数. */
  private readonly _counter: Uint32Array

  /** 每个块的数的个数. */
  private readonly _blockCount: Uint32Array

  /** 每个块的和. */
  private readonly _blockSum: number[]

  /** 每个数所在的块. */
  private readonly _belong: Uint16Array

  /** 所有数的个数. */
  private _len = 0

  /**
   * @param max 值域的最大值.
   * 0 <= max <= 1e6.
   */
  constructor(max: number) {
    max += 5
    const size = Math.sqrt(max) | 0
    const count = 1 + ((max / size) | 0)
    this._blockSize = size
    this._counter = new Uint32Array(max)
    this._blockCount = new Uint32Array(count)
    this._blockSum = Array(count).fill(0)
    const belong = new Uint16Array(max)
    for (let i = 0; i < max; i++) {
      belong[i] = (i / size) | 0
    }
    this._belong = belong
  }

  /** O(1). */
  add(value: number): void {
    this._counter[value]++
    const pos = this._belong[value]
    this._blockCount[pos]++
    this._blockSum[pos] += value
    this._len++
  }

  /** O(1). */
  remove(value: number): void {
    this._counter[value]--
    const pos = this._belong[value]
    this._blockCount[pos]--
    this._blockSum[pos] -= value
    this._len--
  }

  /** O(1). */
  discard(value: number): boolean {
    if (!this.has(value)) return false
    this.remove(value)
    return true
  }

  /** O(1). */
  has(value: number): boolean {
    return !!this._counter[value]
  }

  /** O(sqrt(n)). */
  at(index: number): number | undefined {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) return void 0
    for (let i = 0; i < this._blockCount.length; i++) {
      const count = this._blockCount[i]
      if (index < count) {
        let num = i * this._blockSize
        while (true) {
          const numCount = this._counter[num]
          if (index < numCount) return num
          index -= numCount
          num++
        }
      }
      index -= count
    }
    return void 0
  }

  /**
   * 严格小于 value 的元素个数.
   * 也即第一个大于等于 value 的元素的下标.
   * `O(sqrt(n))`.
   */
  bisectLeft(value: number): number {
    const pos = this._belong[value]
    let res = 0
    for (let i = 0; i < pos; i++) {
      res += this._blockCount[i]
    }
    for (let v = pos * this._blockSize; v < value; v++) {
      res += this._counter[v]
    }
    return res
  }

  /**
   * 小于等于 value 的元素个数.
   * 也即第一个大于 value 的元素的下标.
   * `O(sqrt(n))`.
   */
  bisectRight(value: number): number {
    return this.bisectLeft(value + 1)
  }

  /** O(1). */
  count(value: number): number {
    return this._counter[value]
  }

  /**
   * 返回范围 `[min, max]` 内数的个数.
   * `O(sqrt(n))`.
   */
  countRange(min: number, max: number): number {
    if (min > max) return 0

    const minPos = this._belong[min]
    const maxPos = this._belong[max]
    if (minPos === maxPos) {
      let res = 0
      for (let i = min; i <= max; i++) {
        res += this._counter[i]
      }
      return res
    }

    let res = 0
    const minEnd = (minPos + 1) * this._blockSize
    for (let v = min; v < minEnd; v++) {
      res += this._counter[v]
    }
    for (let i = minPos + 1; i < maxPos; i++) {
      res += this._blockCount[i]
    }
    const maxStart = maxPos * this._blockSize
    for (let v = maxStart; v <= max; v++) {
      res += this._counter[v]
    }
    return res
  }

  /** O(sqrt(n)). */
  lower(value: number): number | undefined {
    // 当前块内寻找
    const pos = this._belong[value]
    const start = pos * this._blockSize
    for (let v = value - 1; v >= start; v--) {
      if (this._counter[v]) return v
    }

    // 按照块的顺序寻找
    for (let i = pos - 1; ~i; i--) {
      if (!this._blockCount[i]) continue
      let num = (i + 1) * this._blockSize - 1
      while (true) {
        if (this._counter[num]) return num
        num--
      }
    }

    return void 0
  }

  /** O(sqrt(n)). */
  higher(value: number): number | undefined {
    // 当前块内寻找
    const pos = this._belong[value]
    const end = (pos + 1) * this._blockSize
    for (let v = value + 1; v < end; v++) {
      if (this._counter[v]) return v
    }

    // 按照块的顺序寻找
    for (let i = pos + 1; i < this._blockCount.length; i++) {
      if (!this._blockCount[i]) continue
      let num = i * this._blockSize
      while (true) {
        if (this._counter[num]) return num
        num++
      }
    }

    return void 0
  }

  /** O(sqrt(n)). */
  floor(value: number): number | undefined {
    if (this.has(value)) return value
    return this.lower(value)
  }

  /** O(sqrt(n)). */
  ceiling(value: number): number | undefined {
    if (this.has(value)) return value
    return this.higher(value)
  }

  /**
   * 返回区间 `[start, end)` 的和.
   * `O(sqrt(n))`.
   */
  sumSlice(start = 0, end = this.length): number {
    if (start < 0) start += this._len
    if (start < 0) start = 0
    if (end < 0) end += this._len
    if (end > this._len) end = this._len
    if (start >= end) return 0

    let res = 0
    let remain = end - start
    let [cur, index] = this._findKth(start)
    const sufCount = this._counter[cur] - index
    if (sufCount >= remain) return remain * cur

    res += sufCount * cur
    remain -= sufCount
    cur++

    // 当前块内的和
    const blockEnd = (this._belong[cur] + 1) * this._blockSize
    while (remain > 0 && cur < blockEnd) {
      const count = this._counter[cur]
      const real = count < remain ? count : remain
      res += real * cur
      remain -= real
      cur++
    }

    // 以块为单位消耗remain
    let pos = this._belong[cur]
    while (remain >= this._blockCount[pos]) {
      res += this._blockSum[pos]
      remain -= this._blockCount[pos]
      pos++
      cur += this._blockSize
    }

    // 剩余的
    while (remain > 0) {
      const count = this._counter[cur]
      const real = count < remain ? count : remain
      res += real * cur
      remain -= real
      cur++
    }

    return res
  }

  /**
   * 返回范围 `[min, max]` 的和.
   * `O(sqrt(n))`.
   */
  sumRange(min: number, max: number): number {
    const minPos = this._belong[min]
    const maxPos = this._belong[max]
    if (minPos === maxPos) {
      let res = 0
      for (let i = min; i <= max; i++) {
        res += this._counter[i] * i
      }
      return res
    }

    let res = 0
    const minEnd = (minPos + 1) * this._blockSize
    for (let v = min; v < minEnd; v++) {
      res += this._counter[v] * v
    }
    for (let i = minPos + 1; i < maxPos; i++) {
      res += this._blockSum[i]
    }
    const maxStart = maxPos * this._blockSize
    for (let v = maxStart; v <= max; v++) {
      res += this._counter[v] * v
    }
    return res
  }

  forEach(callbackfn: (value: number, index: number) => void, reverse = false): void {
    if (reverse) {
      let ptr = 0
      for (let i = this._counter.length - 1; ~i; i--) {
        const count = this._counter[i]
        for (let _ = 0; _ < count; _++) callbackfn(i, ptr++)
      }
    } else {
      let ptr = 0
      for (let i = 0; i < this._counter.length; i++) {
        const count = this._counter[i]
        for (let _ = 0; _ < count; _++) callbackfn(i, ptr++)
      }
    }
  }

  toString(): string {
    const res: number[] = []
    this.forEach(value => res.push(value))
    return `SortedListRangeBlock{${res.join(', ')}}`
  }

  /** O(sqrt(n)). */
  pop(index = -1): number | undefined {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) return void 0
    const value = this.at(index)!
    this.remove(value)
    return value
  }

  slice(start = 0, end = this.length): number[] {
    if (start < 0) start += this._len
    if (start < 0) start = 0
    if (end < 0) end += this._len
    if (end > this._len) end = this._len
    if (start >= end) return []
    const res = Array(end - start)
    let count = 0
    this.enumerate(
      start,
      end,
      value => {
        res[count++] = value
      },
      false
    )
    return res
  }

  /** O(sqrt(n)). */
  erase(start = 0, end = this.length): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerate(start, end, undefined, true)
  }

  enumerate(start: number, end: number, f?: (value: number) => void, erase?: boolean): void {
    if (start < 0) start = 0
    if (end > this._len) end = this._len
    if (start >= end) return

    let remain = end - start
    let [cur, index] = this._findKth(start)
    const sufCount = this._counter[cur] - index
    const real = sufCount < remain ? sufCount : remain
    if (f) {
      for (let _ = 0; _ < real; _++) f(cur)
    }
    if (erase) {
      for (let _ = 0; _ < real; _++) this.remove(cur)
    }
    remain -= sufCount
    cur++

    // 当前块内
    const blockEnd = (this._belong[cur] + 1) * this._blockSize
    while (remain > 0 && cur < blockEnd) {
      const count = this._counter[cur]
      const real = count < remain ? count : remain
      remain -= real
      if (f) {
        for (let _ = 0; _ < real; _++) f(cur)
      }
      if (erase) {
        for (let _ = 0; _ < real; _++) this.remove(cur)
      }
      cur++
    }

    // 以块为单位消耗remain
    let pos = this._belong[cur]
    while (remain >= this._blockCount[pos]) {
      remain -= this._blockCount[pos]
      if (f) {
        for (let v = cur; v < cur + this._blockSize; v++) {
          const count = this._counter[v]
          for (let _ = 0; _ < count; _++) f(v)
        }
      }
      if (erase) {
        this._counter.fill(0, cur, cur + this._blockSize)
        this._len -= this._blockCount[pos]
        this._blockCount[pos] = 0
        this._blockSum[pos] = 0
      }
      pos++
      cur += this._blockSize
    }

    // 剩余的
    while (remain > 0) {
      const count = this._counter[cur]
      const real = count < remain ? count : remain
      remain -= real
      if (f) {
        for (let _ = 0; _ < real; _++) f(cur)
      }
      if (erase) {
        for (let _ = 0; _ < real; _++) this.remove(cur)
      }
      cur++
    }
  }

  clear(): void {
    this._counter.fill(0)
    this._blockCount.fill(0)
    this._blockSum.fill(0)
    this._len = 0
  }

  update(...values: number[]): void {
    values.forEach(value => {
      this.add(value)
    })
  }

  merge(other: ISortedList<number>): void {
    other.forEach(value => {
      this.add(value)
    })
  }

  *entries(): IterableIterator<[index: number, value: number]> {
    let ptr = 0
    for (let v = 0; v < this._counter.length; v++) {
      const count = this._counter[v]
      for (let _ = 0; _ < count; _++) {
        yield [ptr++, v]
      }
    }
  }

  *[Symbol.iterator](): IterableIterator<number> {
    for (let v = 0; v < this._counter.length; v++) {
      const count = this._counter[v]
      for (let _ = 0; _ < count; _++) {
        yield v
      }
    }
  }

  get length(): number {
    return this._len
  }

  get min(): number | undefined {
    return this.at(0)
  }

  get max(): number | undefined {
    if (!this._len) return void 0
    for (let i = this._blockCount.length - 1; ~i; i--) {
      if (!this._blockCount[i]) continue
      let num = (i + 1) * this._blockSize - 1
      while (true) {
        if (this._counter[num]) return num
        num--
      }
    }
    return void 0
  }

  iteratorAt(index: number): ISortedListIterator<number> {
    throw new Error('Method not implemented.')
  }

  lowerBound(value: number): ISortedListIterator<number> {
    throw new Error('Method not implemented.')
  }

  upperBound(value: number): ISortedListIterator<number> {
    throw new Error('Method not implemented.')
  }

  iRange(min: number, max: number, reverse?: boolean | undefined): IterableIterator<number> {
    throw new Error('Method not implemented.')
  }

  iSlice(start: number, end: number, reverse?: boolean | undefined): IterableIterator<number> {
    throw new Error('Method not implemented.')
  }

  /**
   * 返回索引在`kth`处的元素的`value`,以及该元素是`value`中的第几个(`index`).
   */
  private _findKth(kth: number): [value: number, index: number] {
    for (let i = 0; i < this._blockCount.length; i++) {
      const count = this._blockCount[i]
      if (kth < count) {
        let num = i * this._blockSize
        while (true) {
          const numCount = this._counter[num]
          if (kth < numCount) return [num, kth]
          kth -= numCount
          num++
        }
      }
      kth -= count
    }

    throw new Error('unreachable')
  }
}

export { SortedListRangeBlock }

if (require.main === module) {
  debugErase()

  function debugErase(): void {
    const assertEqual = (a: unknown, b: unknown): void => {
      if (a !== b) {
        console.log(a, b)
        throw new Error('assertEqual')
      }
    }

    const n = 10
    const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 10))
    const sl = new SortedListRangeBlock(100)
    sl.update(...nums)
    const sortedNums = nums.sort((a, b) => a - b)

    for (let i = 0; i < 20; i++) {
      let start = Math.floor(Math.random() * sl.length)
      let end = Math.floor(Math.random() * sl.length)

      sl.erase(start, end)
      sortedNums.splice(start, end - start)

      // expect(sl.length).toBe(sortedNums.length)
      assertEqual(sl.length, sortedNums.length)
      assertEqual(sortedNums.join(','), [...sl].join(','))
    }
  }

  // https://leetcode.cn/problems/sliding-subarray-beauty/ 2200ms
  function getSubarrayBeauty(nums: number[], k: number, x: number): number[] {
    const res: number[] = []
    const sl = new SortedListRangeBlock(200)
    const OFFSET = 100
    const n = nums.length
    for (let right = 0; right < n; right++) {
      sl.add(nums[right] + OFFSET)
      if (right >= k) {
        sl.discard(nums[right - k] + OFFSET)
      }
      if (right >= k - 1) {
        const xth = sl.at(x - 1)! - OFFSET
        res.push(xth < 0 ? xth : 0)
      }
    }
    return res
  }
}
