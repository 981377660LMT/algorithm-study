/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { SortedDictFast } from '../../22_专题/离线查询/根号分治/SortedList/SortedDictFast'

const INF = 2e15

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `SortedList`实现.
 * 初始时，默认区间为`[-INF,INF)`，值为`noneValue`.
 * @deprecated 使用`ODTVan`更快,但是空间占用更多.
 */
class ODTMap<S> {
  private _count = 0
  private _len = 0
  private readonly _leftLimit: number
  private readonly _rightLimit: number
  private readonly _data: SortedDictFast<number, S> = new SortedDictFast()
  private readonly _noneValue: S

  /**
   * 指定哨兵值建立一个ODTMap.
   * @param noneValue 表示空值的哨兵值.
   * @param leftLimit 区间左端点(包含).
   * @param rightLimit 区间右端点(不包含).
   */
  constructor(noneValue: S, leftLimit = -INF, rightLimit = INF) {
    this._noneValue = noneValue
    this._leftLimit = leftLimit
    this._rightLimit = rightLimit
    this._data.set(this._leftLimit, noneValue)
    this._data.set(this._rightLimit, noneValue)
  }

  /**
   * 返回包含`x`的区间的信息.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] {
    const pos2 = this._data.bisectRight(x)
    const pos1 = pos2 - 1
    const [l, vl] = this._data.peekItem(pos1)!
    const r = this._data.peekItem(pos2)![0]
    if (vl !== this._noneValue && erase) {
      this._len--
      this._count -= r - l
      this._data.set(l, this._noneValue)
      this._mergeAt(l)
      this._mergeAt(r)
    }
    return [l, r, vl]
  }

  /**
   * 将区间`[start,end)`的值置为`value`.
   */
  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
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
    if (start >= end) return

    if (!erase) {
      let pos = this._data.bisectRight(start) - 1
      let [k1, v1] = this._data.peekItem(pos)!
      while (k1 < end) {
        pos++
        const [k2, v2] = this._data.peekItem(pos)!
        f(Math.max(k1, start), Math.min(k2, end), v1)
        k1 = k2
        v1 = v2
      }
      return
    }

    let pos = this._data.bisectRight(start) - 1
    let [k, v] = this._data.peekItem(pos)!
    if (k < start) {
      this._data.set(start, v)
      if (v !== this._noneValue) this._len++
    }

    pos = this._data.bisectLeft(end)
    // eslint-disable-next-line semi-style
    ;[k, v] = this._data.peekItem(pos)!
    if (k > end) {
      const v2 = this._data.peekItem(pos - 1)![1]
      this._data.set(end, v2)
      if (v2 !== this._noneValue) this._len++
    }

    pos = this._data.bisectLeft(start)
    let [k1, v1] = this._data.peekItem(pos)!
    while (k1 < end) {
      const [k2, v2] = this._data.peekItem(pos + 1)!
      f(k1, k2, v1)
      if (v1 !== this._noneValue) {
        this._len--
        this._count -= k2 - k1
      }
      this._data.popItem(pos)
      k1 = k2
      v1 = v2
    }

    this._data.set(start, this._noneValue)
  }

  toString(): string {
    const sb: string[] = [`ODTMap(${this.length}) {`]
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
   * 区间内元素个数总和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p === this._leftLimit || p === this._rightLimit) return
    const pos1 = this._data.bisectLeft(p)
    const pos2 = pos1 - 1
    const v1 = this._data.peekItem(pos1)![1]
    const v2 = this._data.peekItem(pos2)![1]
    if (v1 === v2) {
      if (v1 !== this._noneValue) this._len--
      this._data.popItem(pos1)
    }
  }
}

export { ODTMap }

if (require.main === module) {
  const odtMap = new ODTMap<number>(-1)
  odtMap.set(0, 1, 1)
  odtMap.set(-1, 0, 1)
  console.log(odtMap.toString())

  // https://leetcode.cn/problems/count-integers-in-intervals/
  class CountIntervals {
    private readonly _odtMap = new ODTMap<number>(-1)

    add(left: number, right: number): void {
      this._odtMap.set(left, right + 1, 1)
    }

    count(): number {
      return this._odtMap.count
    }
  }

  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  // 352. 将数据流变为多个不相交区间
  class SummaryRanges {
    private readonly _odt = new ODTMap(-1)

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

  // https://leetcode.cn/problems/range-module/
  class RangeModule {
    private readonly _odt = new ODTMap(-1)

    addRange(left: number, right: number): void {
      this._odt.set(left, right, 0)
    }

    queryRange(left: number, right: number): boolean {
      const [start, end, value] = this._odt.get(left)
      return start <= left && right <= end && value === 0
    }

    removeRange(left: number, right: number): void {
      this._odt.set(left, right, -1)
    }
  }

  // 2655. 寻找最大长度的未覆盖区间
  // https://leetcode.cn/problems/find-maximal-uncovered-ranges/
  function findMaximalUncoveredRanges(n: number, ranges: number[][]): number[][] {
    const odt = new ODTMap(-1)
    ranges.forEach(([left, right]) => {
      odt.set(left, right + 1, 1)
    })
    const res: number[][] = []
    odt.enumerateRange(0, n, (s, e, v) => {
      if (v === -1) res.push([s, e - 1])
    })
    return res
  }

  // !56. 合并区间自动合并的方式与这里不同,不能使用
  // https://leetcode.cn/problems/merge-intervals/
  //

  // 1272. 删除区间
  // https://leetcode.cn/problems/remove-interval/
  function removeInterval(intervals: number[][], toBeRemoved: number[]): number[][] {
    const odt = new ODTMap(-1)
    intervals.forEach(([start, end]) => {
      odt.set(start, end, 1)
    })
    odt.set(toBeRemoved[0], toBeRemoved[1], -1)
    const res: number[][] = []
    odt.enumerateAll((start, end, value) => {
      if (value === 1) res.push([start, end])
    })
    return res
  }

  // 128. 最长连续序列
  // https://leetcode.cn/problems/longest-consecutive-sequence/
  // 给定一个未排序的整数数组 nums ，找出数字连续的最长序列（不要求序列元素在原数组中连续）的长度。
  function longestConsecutive(nums: number[]): number {
    const odt = new ODTMap(-1)
    nums.forEach(v => odt.set(v, v + 1, 1))
    let res = 0
    odt.enumerateAll((start, end, value) => {
      if (value === 1) {
        res = Math.max(res, end - start)
      }
    })
    return res
  }
}
