/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable prefer-destructuring */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable no-inner-declarations */
// 珂朵莉树(ODT)/Intervals
// !noneValue不使用symbol,而是自定义的哨兵值,更加灵活.

import { VanEmdeBoasTree } from './VanEmdeBoasTree'

const INF = 2e15

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `VanEmdeBoasTree`实现.
 * @warning 允许插入的元素范围为[0,INF).
 */
class ODTVan<S> {
  private readonly _noneValue: S
  private _len = 0
  private _count = 0
  private readonly _data: Map<number, S> = new Map()
  private readonly _leftLimit = 0
  private readonly _rightLimit = INF
  private readonly _fs: VanEmdeBoasTree = new VanEmdeBoasTree()

  /**
   * 指定哨兵值建立一个ODT.初始时,所有位置的值为 {@link noneValue}.
   * 默认区间为`[0,INF)`，值为`noneValue`.
   * @param noneValue 表示空值的哨兵值.
   * @warning 允许插入的元素范围为[0,INF).
   */
  constructor(noneValue: S) {
    this._noneValue = noneValue
  }

  /**
   * 返回包含`x`的区间的信息.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] {
    const start = this._fs.prev(x)
    const end = this._fs.next(x + 1)
    const value = this._getOrNone(start)
    if (erase && value !== this._noneValue) {
      this._len--
      this._count -= end - start
      this._data.set(start, this._noneValue)
      this._mergeAt(start)
      this._mergeAt(end)
    }
    return [start, end, value]
  }

  /**
   * !start>=0.
   */
  set(start: number, end: number, value: S): void {
    if (start < 0) throw new Error('start must be non-negative!')
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
    this._fs.insert(start)
    this._data.set(start, value)
    if (value !== this._noneValue) {
      this._len++
      this._count += end - start
    }
    this._mergeAt(start)
    this._mergeAt(end)
  }

  enumerateAll(f: (start: number, end: number, value: S) => void): void {
    this.enumerateRange(this._leftLimit, this._rightLimit, f, false)
  }

  /**
   * 遍历范围`[start, end)`内的所有区间.
   */
  enumerateRange(
    start: number,
    end: number,
    f: (start: number, end: number, value: S) => void,
    erase = false
  ): void {
    if (start < this._leftLimit) start = this._leftLimit
    if (end > this._rightLimit) end = this._rightLimit
    if (start >= end) return

    const none = this._noneValue
    if (!erase) {
      let left = this._fs.prev(start)
      while (left < end) {
        const right = this._fs.next(left + 1)
        f(Math.max(left, start), Math.min(right, end), this._getOrNone(left))
        left = right
      }
      return
    }

    let p = this._fs.prev(start)
    if (p < start) {
      this._fs.insert(start)
      const v = this._getOrNone(p)
      this._data.set(start, v)
      if (v !== none) this._len++
    }

    p = this._fs.next(end)
    if (end < p) {
      const v = this._getOrNone(this._fs.prev(end))
      this._data.set(end, v)
      this._fs.insert(end)
      if (v !== none) this._len++
    }

    p = start
    while (p < end) {
      const q = this._fs.next(p + 1)
      const x = this._getOrNone(p)
      f(p, q, x)
      if (x !== none) {
        this._len--
        this._count -= q - p
      }
      this._fs.erase(p)
      p = q
    }

    this._fs.insert(start)
    this._data.set(start, none)
  }

  toString(): string {
    const sb: string[] = [`ODT(${this.length}) {`]
    this.enumerateAll((start, end, value) => {
      const v = value === this._noneValue ? 'null' : value
      sb.push(`  [${start},${end}):${v}`)
    })
    sb.push('}')
    return sb.join('\n')
  }

  /**
   * 区间个数.
   */
  get length(): number {
    return this._len
  }

  /**
   * 区间内元素个数之和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p <= 0 || this._rightLimit <= p) return
    const q = this._fs.prev(p - 1)
    const dataP = this._getOrNone(p)
    const dataQ = this._getOrNone(q)
    if (dataP === dataQ) {
      if (dataP !== this._noneValue) this._len--
      this._fs.erase(p)
    }
  }

  private _getOrNone(x: number): S {
    const res = this._data.get(x)
    return res == undefined ? this._noneValue : res
  }
}

export { ODTVan }

if (require.main === module) {
  const INF = 2e15
  const van = new ODTVan(INF)
  console.log(van.get(-1))
  console.log(van.toString())
  van.set(0, 10, 1)
  van.set(2, 5, 2)
  console.log(van.get(8))
  console.log(van.toString())
  van.enumerateRange(
    -1000,
    7,
    (start, end, value) => {
      console.log(start, end, value)
    },
    false
  )
  console.log(van.toString(), van.length)

  // 352. 将数据流变为多个不相交区间
  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  class SummaryRanges {
    private readonly _odt = new ODTVan(-1)

    addNum(value: number): void {
      this._odt.set(value, value + 1, 0)
    }

    getIntervals(): number[][] {
      const res: number[][] = []
      this._odt.enumerateAll((start, end, value) => {
        if (value === 0) res.push([start, end - 1])
      })
      return res
    }
  }

  // 128. 最长连续序列
  // https://leetcode.cn/problems/longest-consecutive-sequence/
  // 给定一个未排序的整数数组 nums ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。
  function longestConsecutive(nums: number[]): number {
    const odt = new ODTVan(0)
    const OFFSET = 1e9 + 10
    nums.forEach(v => odt.set(v + OFFSET, v + 1 + OFFSET, 1))
    let res = 0
    odt.enumerateAll((start, end, value) => {
      if (value === 1) {
        res = Math.max(res, end - start)
      }
    })
    return res
  }
}
