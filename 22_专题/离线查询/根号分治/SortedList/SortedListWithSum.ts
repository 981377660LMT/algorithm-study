/* eslint-disable brace-style */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable no-inner-declarations */
/* eslint-disable generator-star-spacing */
/* eslint-disable prefer-destructuring */

// class  SortedListWithSum

import { SortedListFast } from './SortedListFast'

interface IOptions<V> {
  values?: Iterable<V>
  compareFn?: (a: V, b: V) => number
  abelGroup?: {
    e: () => V
    op: (a: V, b: V) => V
    inv: (a: V) => V
  }
}

interface ISortedListFastWithSum<V> {
  sumSlice(start?: number, end?: number): V
  sumRange(min: V, max: V): V
  sumAll(): V
}

/**
 * 支持区间求和的有序列表.
 * {@link sumSlice} 和 {@link sumRange} 的时间复杂度为 `O(sqrt(n))`.
 * {@link sumAll} 的时间复杂度为 `O(1)`.
 */
class SortedListFastWithSum<V = number>
  extends SortedListFast<V>
  implements ISortedListFastWithSum<V>
{
  private readonly _e: () => V
  private readonly _op: (a: V, b: V) => V
  private readonly _inv: (a: V) => V
  private _allSum: V
  private _sums: V[] = []

  constructor(options?: IOptions<V>) {
    super()
    const {
      values = [],
      compareFn = (a: any, b: any) => a - b,
      abelGroup = {
        e: () => 0 as any,
        op: (a: any, b: any) => a + b,
        inv: (a: any) => -a as any
      }
    } = options ?? {}
    this._e = abelGroup.e
    this._op = abelGroup.op
    this._inv = abelGroup.inv
    this._allSum = this._e()
    this._reBuild([...values], compareFn)
  }

  /**
   * 返回区间 `[start, end)` 的和.
   */
  sumSlice(start = 0, end = this.length): V {
    if (start < 0) start += this._len
    if (start < 0) start = 0
    if (end < 0) end += this._len
    if (end > this._len) end = this._len
    if (start >= end) return this._e()

    let res = this._e()
    let { pos, index } = this._findKth(start)
    let count = end - start
    for (; count && pos < this._blocks.length; pos++) {
      const block = this._blocks[pos]
      const endPos = Math.min(block.length, index + count)
      const curCount = endPos - index
      if (curCount === block.length) {
        res = this._op(res, this._sums[pos])
      } else {
        for (let j = index; j < endPos; j++) res = this._op(res, block[j])
      }
      count -= curCount
      index = 0
    }

    return res
  }

  /**
   * 返回范围 `[min, max]` 的和.
   */
  sumRange(min: V, max: V): V {
    if (this._compareFn(min, max) > 0) return this._e()

    let res = this._e()
    let { pos, index: start } = this._locLeft(min)
    for (let i = pos; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(max, block[0]) < 0) break
      if (start === 0 && this._compareFn(block[block.length - 1], max) <= 0) {
        res = this._op(res, this._sums[i])
      } else {
        for (let j = start; j < block.length; j++) {
          const cur = block[j]
          if (this._compareFn(cur, max) > 0) break
          res = this._op(res, cur)
        }
      }
      start = 0
    }
    return res
  }

  sumAll(): V {
    return this._allSum
  }

  // eslint-disable-next-line class-methods-use-this
  protected override _init(): void {
    return undefined
  }

  protected override _reBuild(data: V[], compareFn: (a: V, b: V) => number): void {
    data.sort(compareFn)
    const n = data.length
    const blocks = []
    const sums = []
    let allSum = this._e()
    for (let i = 0; i < n; i += SortedListFast._LOAD) {
      const newBlock = data.slice(i, i + SortedListFast._LOAD)
      blocks.push(newBlock)
      let cur = this._e()
      for (let j = 0; j < newBlock.length; j++) cur = this._op(cur, newBlock[j])
      sums.push(cur)
      allSum = this._op(allSum, cur)
    }
    const mins = Array(blocks.length)
    for (let i = 0; i < blocks.length; i++) {
      const cur = blocks[i]
      mins[i] = cur[0]
    }

    this._compareFn = compareFn
    this._len = n
    this._blocks = blocks
    this._mins = mins
    this._tree = []
    this._shouldRebuildTree = true
    this._sums = sums
    this._allSum = allSum
  }

  override add(value: V): this {
    const { _blocks, _mins, _sums } = this
    this._len++
    this._allSum = this._op(this._allSum, value)
    if (!_blocks.length) {
      _blocks.push([value])
      _mins.push(value)
      _sums.push(value)
      this._shouldRebuildTree = true
      return this
    }

    const load = SortedListFast._LOAD
    const { pos, index } = this._locRight(value)
    this._updateTree(pos, 1)
    const block = _blocks[pos]
    block.splice(index, 0, value)
    _mins[pos] = block[0]
    _sums[pos] = this._op(_sums[pos], value)

    // !-> [x]*load + [x]*(block.length - load)
    if (load + load < block.length) {
      const oldSum = _sums[pos]
      _blocks.splice(pos + 1, 0, block.slice(load))
      _mins.splice(pos + 1, 0, block[load])
      block.splice(load)
      this._shouldRebuildTree = true
      this._rebuildSum(pos)
      this._sums.splice(pos + 1, 0, this._op(oldSum, this._inv(this._sums[pos])))
    }

    return this
  }

  override enumerate(start: number, end: number, f?: (value: V) => void, erase?: boolean): void {
    if (start < 0) start = 0
    if (end > this._len) end = this._len
    if (start >= end) return

    let { pos, index: startIndex } = this._findKth(start)
    let count = end - start
    for (; count && pos < this._blocks.length; pos++) {
      const block = this._blocks[pos]
      const endIndex = Math.min(block.length, startIndex + count)
      if (f) {
        for (let j = startIndex; j < endIndex; j++) f(block[j])
      }
      const deleted = endIndex - startIndex

      if (erase) {
        if (deleted === block.length) {
          // !delete block
          this._blocks.splice(pos, 1)
          this._mins.splice(pos, 1)
          this._allSum = this._op(this._allSum, this._inv(this._sums[pos]))
          this._sums.splice(pos, 1)
          this._shouldRebuildTree = true
          pos--
        } else {
          // !delete [index, end)
          for (let i = startIndex; i < endIndex; i++) {
            this._updateTree(pos, -1)
            const inv = this._inv(block[i])
            this._allSum = this._op(this._allSum, inv)
            this._sums[pos] = this._op(this._sums[pos], inv)
          }
          block.splice(startIndex, deleted)
          this._mins[pos] = block[0]
        }
        this._len -= deleted
      }

      count -= deleted
      startIndex = 0
    }
  }

  override clear(): void {
    super.clear()
    this._sums = []
    this._allSum = this._e()
  }

  protected override _delete(pos: number, index: number): void {
    const { _blocks, _mins, _sums } = this

    // !delete element
    this._len--
    this._updateTree(pos, -1)
    const block = _blocks[pos]
    const deleted = block[index]
    block.splice(index, 1)
    this._allSum = this._op(this._allSum, this._inv(deleted))
    if (block.length) {
      _mins[pos] = block[0]
      _sums[pos] = this._op(_sums[pos], this._inv(deleted))
      return
    }

    // !delete block
    _blocks.splice(pos, 1)
    _mins.splice(pos, 1)
    _sums.splice(pos, 1)
    this._shouldRebuildTree = true
  }

  private _rebuildSum(pos: number): void {
    let cur = this._e()
    const block = this._blocks[pos]
    for (let i = 0; i < block.length; i++) cur = this._op(cur, block[i])
    this._sums[pos] = cur
  }
}

export { SortedListFastWithSum }

if (require.main === module) {
  // 可用于维护topKSum
  const sl = new SortedListFastWithSum<number>({ values: [1] })
  const nums = [11, 2, 3, 4, 5, 6, 7, 8, 9]
  for (const num of nums) sl.add(num)
  console.log(sl.sumSlice(0, 3))
  sl.erase(0, 3)
  console.log(sl.sumSlice(0, 3))
  console.log(sl.sumRange(0, 1000))

  console.log('-'.repeat(20))

  // https://leetcode.cn/problems/minimum-difference-in-sums-after-removal-of-elements/
  function minimumDifference(nums: number[]): number {
    const n = nums.length / 3
    const pre = new SortedListFastWithSum({ values: nums.slice(0, n) })
    const suf = new SortedListFastWithSum({ values: nums.slice(n) })
    let res = pre.sumSlice(0, n) - suf.sumSlice(suf.length - n, suf.length)
    for (let i = n; i < 2 * n; i++) {
      pre.add(nums[i])
      suf.discard(nums[i])
      res = Math.min(res, pre.sumSlice(0, n) - suf.sumSlice(suf.length - n, suf.length))
    }
    return res
  }

  // 1825. 求出 MK 平均值
  // https://leetcode.cn/problems/finding-mk-average/
  class MKAverage {
    private readonly _m: number
    private readonly _k: number
    private readonly _queue: number[] = []
    private _head = 0
    private _lastK = new SortedListFastWithSum<number>()

    constructor(m: number, k: number) {
      this._m = m
      this._k = k
    }

    addElement(num: number): void {
      this._queue.push(num)
      if (this._queue.length === this._m) {
        this._lastK = new SortedListFastWithSum({ values: this._queue.slice(-this._m) })
        return
      }
      if (this._queue.length > this._m) {
        this._lastK.add(num)
        this._lastK.discard(this._queue[this._head])
        this._head++
      }
    }

    calculateMKAverage(): number {
      if (this._queue.length < this._m) return -1
      const midSum = this._lastK.sumSlice(this._k, -this._k)
      return Math.floor(midSum / (this._m - 2 * this._k))
    }
  }

  testSumSlice()
  testSumRange()
  function testSumSlice() {
    const sl = new SortedListFastWithSum<number>()

    let sortedNums: number[] = []
    for (let i = 0; i < 10000; i++) {
      const num = Math.floor(Math.random() * 10000)
      sl.add(num)
      sortedNums.push(num)
      sortedNums.sort((a, b) => a - b)
      const start = ~~(sl.length * Math.random())
      const end = ~~(sl.length * Math.random())
      const res1 = sl.sumSlice(start, end)
      const res2 = sortedNums.slice(start, end).reduce((a, b) => a + b, 0)
      if (res1 !== res2) {
        console.log(res1, res2)
        console.log(sl.slice(0, 10))
        console.log(sortedNums.slice(0, 10))
        throw new Error()
      }

      // discard/add
      const willDiscard = Math.random() > 0.5
      if (willDiscard) {
        const discard = sortedNums[~~(Math.random() * sortedNums.length)]
        sl.discard(discard)
        const index = sortedNums.findIndex(num => num === discard)
        sortedNums.splice(index, 1)
      } else {
        const add = Math.floor(Math.random() * 100)
        sl.add(add)
        sortedNums.push(add)
        sortedNums.sort((a, b) => a - b)
      }
    }
  }
  function testSumRange() {
    // SortedListFastWithSum.setLoadFactor(5)
    const sl = new SortedListFastWithSum<number>()

    const sortedNums = []
    for (let i = 0; i < 10000; i++) {
      const num = Math.floor(Math.random() * 10000)
      sl.add(num)
      sortedNums.push(num)
      sortedNums.sort((a, b) => a - b)
      const min = ~~(Math.random() * 10000)
      const max = ~~(Math.random() * 10000)
      const res1 = sl.sumRange(min, max)
      const res2 = sortedNums.filter(num => num >= min && num <= max).reduce((a, b) => a + b, 0)
      if (res1 !== res2) {
        console.log(res1, res2)
        console.log(sl.slice(0, 10))
        console.log(sortedNums.slice(0, 10))
        throw new Error()
      }
      const willDiscard = Math.random() > 0.5
      if (willDiscard) {
        const randint = sortedNums[~~(Math.random() * sortedNums.length)]
        const index = sortedNums.findIndex(num => num === randint)
        sl.discard(randint)
        sortedNums.splice(index, 1)
      } else {
        const randint = Math.floor(Math.random() * 100)
        sl.add(randint)
        sortedNums.push(randint)
        sortedNums.sort((a, b) => a - b)
      }
    }
  }

  console.log('ok')
}
