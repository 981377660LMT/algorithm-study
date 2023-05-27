// 珂朵莉树(ODT)/Intervals

import { FastSet } from './FastSet'

const INF = 2e15

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `FastSet`实现.
 */
class ODT<S> {
  private _len = 0
  private _count = 0
  private readonly _leftLimit: number
  private readonly _rightLimit: number
  private readonly _noneValue: S
  private readonly _data: S[]
  private readonly _fs: FastSet

  /**
   * 指定区间长度和哨兵值建立一个ODT.初始时,所有位置的值为 {@link noneValue}.
   * @param n 区间范围为`[0, n)`.
   * @param noneValue 表示空值的哨兵值.
   */
  constructor(n: number, noneValue: S) {
    const data = Array(n)
    for (let i = 0; i < n; i++) data[i] = noneValue
    const fs = new FastSet(n)
    fs.insert(0)

    this._leftLimit = 0
    this._rightLimit = n
    this._noneValue = noneValue
    this._data = data
    this._fs = fs
  }

  /**
   * 返回包含`x`的区间的信息.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] {
    const start = this._fs.prev(x)
    const end = this._fs.next(x + 1)
    const value = this._data[start]
    if (erase && value !== this._noneValue) {
      this._len--
      this._count -= end - start
      this._data[start] = this._noneValue
      this._mergeAt(start)
      this._mergeAt(end)
    }
    return [start, end, value]
  }

  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
    this._fs.insert(start)
    this._data[start] = value
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
        f(Math.max(left, start), Math.min(right, end), this._data[left])
        left = right
      }
      return
    }

    let p = this._fs.prev(start)
    if (p < start) {
      this._fs.insert(start)
      this._data[start] = this._data[p]
      if (this._data[start] !== none) {
        this._len++
      }
    }

    p = this._fs.next(end)
    if (end < p) {
      this._data[end] = this._data[this._fs.prev(end)]
      this._fs.insert(end)
      if (this._data[end] !== none) {
        this._len++
      }
    }

    p = start
    while (p < end) {
      const q = this._fs.next(p + 1)
      const x = this._data[p]
      f(p, q, x)
      if (this._data[p] !== none) {
        this._len--
        this._count -= q - p
      }
      this._fs.erase(p)
      p = q
    }

    this._fs.insert(start)
    this._data[start] = none
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateAll((start, end, value) => {
      const v = value === this._noneValue ? 'null' : value
      sb.push(`[${start},${end}):${v}`)
    })
    return `ODT{${sb.join(', ')}}`
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
    if (this._data[p] === this._data[q]) {
      if (this._data[p] !== this._noneValue) this._len--
      this._fs.erase(p)
    }
  }
}

export { ODT }

if (require.main === module) {
  const odt = new ODT(10, INF)
  console.log(odt.toString())
  odt.set(0, 10, 1)
  odt.set(2, 5, 2)
  console.log(odt.get(8))
  console.log(odt.toString())
  odt.enumerateRange(
    1,
    7,
    (start, end, value) => {
      console.log(start, end, value)
    },
    true
  )
  console.log(odt.toString(), odt.length)

  // 352. 将数据流变为多个不相交区间
  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  class SummaryRanges {
    private readonly _odt = new ODT(1e4 + 10, -1)

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

  /**
   * Your SummaryRanges object will be instantiated and called as such:
   * var obj = new SummaryRanges()
   * obj.addNum(value)
   * var param_2 = obj.getIntervals()
   */
}
