/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable class-methods-use-this */

// API:
//  insert(left, right)              向区间集合中插入一个区间.
//  discard(left, right)               从区间集合中删除一个区间.
//  nextStart(x)                     返回第一个大于等于x的区间起点.
//  prevStart(x)                     返回最后一个小于等于x的区间起点.
//  ceiling(x)                       返回区间内第一个大于等于x的元素.
//  floor(x)                         返回区间内第一个小于等于x的元素.
//  getRange(x)                      返回包含x的区间.
//  includes(x)                      判断x是否在区间集合中.
//  includesInterval(left, right)    判断[left, right]是否在区间集合中.
//  at(i)                            返回第i个区间.
//  getAll()                         返回所有区间.
//  islice(min,max)                  返回区间集合中包含在[min,max]区间内的所有区间的迭代器.
//  enumerate(min,max,f)             遍历SegmentSet中包含在[min,max]区间内的所有区间范围.
//  length                           SegmentSet中区间的个数.
//  count                            SegmentSet中区间的元素数量.

import { SortedListFast } from '../离线查询/根号分治/SortedList/SortedListFast'

const INF = 2e15

type Interval = { left: number; right: number }

/**
 * 管理区间的数据结构.
 *
 * - 所有区间都是`闭区间` 例如 [1,1] 表示 长为1的区间,起点为1;
 * - 有交集的区间会被合并,例如 [1,2]和[2,3]会被合并为[1,3].
 */
class SegmentSet {
  private readonly _sl = new SortedListFast<Interval>((a, b) => a.left - b.left)
  private readonly _flyWeight = { left: 0, right: 0 }
  private _count = 0

  /**
   * 向区间集合中插入一个闭区间.
   */
  insert(left: number, right: number): void {
    if (left > right) {
      return
    }
    let it1 = this._sl.bisectRight(this.getFlyWeight(left, INF))
    const it2 = this._sl.bisectRight(this.getFlyWeight(right, INF))
    if (it1 > 0 && this._sl.at(it1 - 1)!.right >= left) {
      it1--
    }
    if (it1 !== it2) {
      const tmp1 = this._sl.at(it1)!.left
      left = Math.min(left, tmp1)
      const tmp2 = this._sl.at(it2 - 1)!.right
      right = Math.max(right, tmp2)
      let removed = 0
      this._sl.enumerate(
        it1,
        it2,
        v => {
          removed += v.right - v.left + 1
        },
        true
      )
      this._count -= removed
    }
    this._sl.add({ left, right })
    this._count += right - left + 1
  }

  /**
   * 从区间集合中删除一个闭区间.
   */
  discard(left: number, right: number): boolean {
    if (left > right) {
      return false
    }
    let it1 = this._sl.bisectRight(this.getFlyWeight(left, INF))
    const it2 = this._sl.bisectRight(this.getFlyWeight(right, INF))
    if (it1 > 0 && this._sl.at(it1 - 1)!.right >= left) {
      it1--
    }
    if (it1 === it2) {
      return false
    }
    let nl = Math.min(left, this._sl.at(it1)!.left)
    let nr = Math.max(right, this._sl.at(it2 - 1)!.right)
    let removed = 0
    this._sl.enumerate(
      it1,
      it2,
      v => {
        removed += v.right - v.left + 1
      },
      true
    )
    this._count -= removed
    if (nl < left) {
      this._sl.add({ left: nl, right: left })
      this._count += left - nl + 1
    }
    if (nr > right) {
      this._sl.add({ left: right, right: nr })
      this._count += nr - right + 1
    }
    return true
  }

  /**
   * 返回第一个大于等于x的区间起点.
   */
  nextStart(x: number): number | undefined {
    const it = this._sl.bisectLeft(this.getFlyWeight(x, -INF))
    if (it === this._sl.length) {
      return undefined
    }
    return this._sl.at(it)!.left
  }

  /**
   * 返回最后一个小于等于x的区间起点.
   */
  prevStart(x: number): number | undefined {
    const it = this._sl.bisectRight(this.getFlyWeight(x, INF)) - 1
    if (it < 0) {
      return undefined
    }
    return this._sl.at(it)!.left
  }

  /**
   * 返回区间内第一个大于等于x的元素.
   */
  floor(x: number): number | undefined {
    const it = this._sl.bisectRight(this.getFlyWeight(x, INF))
    if (it === 0) {
      return undefined
    }
    return Math.min(x, this._sl.at(it - 1)!.right)
  }

  /**
   * 返回区间内第一个小于等于x的元素.
   */
  ceiling(x: number): number | undefined {
    const it = this._sl.bisectRight(this.getFlyWeight(x, INF))
    if (it > 0 && this._sl.at(it - 1)!.right >= x) {
      return x
    }
    if (it !== this._sl.length) {
      return this._sl.at(it)!.left
    }
    return undefined
  }

  /**
   * 返回包含x的区间.
   */
  getInterval(x: number): Interval | undefined {
    const it = this._sl.bisectRight(this.getFlyWeight(x, INF))
    if (it === 0 || this._sl.at(it - 1)!.right < x) {
      return undefined
    }
    return this._sl.at(it - 1)!
  }

  /**
   * 判断x是否在区间集合中.
   */
  includes(x: number): boolean {
    const it = this._sl.bisectRight(this.getFlyWeight(x, INF))
    return it > 0 && this._sl.at(it - 1)!.right >= x
  }

  /**
   * 判断区间[left,right]是否在区间集合中.
   */
  includesInterval(left: number, right: number): boolean {
    if (left > right) {
      return false
    }
    const it1 = this._sl.bisectRight(this.getFlyWeight(left, INF))
    if (it1 === 0) {
      return false
    }
    const it2 = this._sl.bisectRight(this.getFlyWeight(right, INF))
    if (it1 !== it2) {
      return false
    }
    return this._sl.at(it1 - 1)!.right >= right
  }

  /**
   * 返回第index个区间.
   */
  at(index: number): Interval | undefined {
    return this._sl.at(index)
  }

  /**
   * 有时需要输出闭区间, insert/remove 时右端点需要加一，getAll 时右端点需要减一.
   * ```ts
   * seg.insert(left, right + 1)
   * seg.remove(left, right + 1)
   * seg.getAll().map(v => [v[0], v[1] - 1])
   * ```
   */
  getAll(): Interval[] {
    const res: Interval[] = Array(this._sl.length)
    this._sl.forEach((v, i) => {
      res[i] = v
    })
    return res
  }

  /**
   * 返回一个迭代器, 遍历 SegmentSet 中在 `[min,max]` 内的所有区间范围.
   */
  *irange(min: number, max: number): IterableIterator<Interval> {
    if (min > max) {
      return
    }
    let it = this._sl.bisectRight(this.getFlyWeight(min, INF)) - 1
    if (it < 0) it++
    const islice = this._sl.slice(it, this._sl.length)
    for (const v of islice) {
      if (v.left > max) {
        return
      }
      yield { left: Math.max(v.left, min), right: Math.min(v.right, max) }
    }
  }

  /**
   * 遍历 SegmentSet 中在 `[min,max]` 内的所有区间范围.
   */
  enumerateRange(min: number, max: number, f: (interval: Interval) => void): void {
    if (min > max) {
      return
    }
    let it = this._sl.bisectRight(this.getFlyWeight(min, INF)) - 1
    if (it < 0) it++
    const islice = this._sl.slice(it, this._sl.length)
    for (const v of islice) {
      if (v.left > max) {
        break
      }
      f({ left: Math.max(v.left, min), right: Math.min(v.right, max) })
    }
  }

  toString(): string {
    return `SegmentSet(${this.getAll().join(', ')})`
  }

  getFlyWeight(left: number, right: number): Interval {
    this._flyWeight.left = left
    this._flyWeight.right = right
    return this._flyWeight
  }

  *[Symbol.iterator](): IterableIterator<Interval> {
    yield* this._sl
  }

  get length(): number {
    return this._sl.length
  }

  get count(): number {
    return this._count
  }
}

if (require.main === module) {
  // https://leetcode.cn/problems/count-integers-in-intervals/description/
  // 2276. 统计区间中的整数数目
  class CountIntervals {
    private readonly _ss = new SegmentSet()

    add(left: number, right: number): void {
      this._ss.insert(left, right)
    }

    count(): number {
      return this._ss.count
    }
  }

  // 352. 将数据流变为多个不相交区间
  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  class SummaryRanges {
    private readonly _ss = new SegmentSet()

    addNum(value: number): void {
      this._ss.insert(value, value + 1)
    }

    getIntervals(): number[][] {
      return this._ss.getAll().map(v => [v.left, v.right - 1])
    }
  }

  // https://leetcode.cn/problems/range-module/submissions/
  // Range模块/Range 模块
  class RangeModule {
    private readonly _ss = new SegmentSet()

    addRange(left: number, right: number): void {
      this._ss.insert(left, right)
    }

    queryRange(left: number, right: number): boolean {
      return this._ss.includesInterval(left, right)
    }

    removeRange(left: number, right: number): void {
      this._ss.discard(left, right)
    }
  }

  // https://leetcode.cn/problems/find-maximal-uncovered-ranges/
  function findMaximalUncoveredRanges(n: number, ranges: number[][]): number[][] {
    const seg = new SegmentSet()
    seg.insert(0, n) // 右边界+1
    ranges.forEach(v => {
      seg.discard(v[0], v[1] + 1) // 右边界+1
    })
    return seg.getAll().map(v => [v.left, v.right - 1]) // 右边界-1
  }

  // 56. 合并区间
  // https://leetcode.cn/problems/merge-intervals/description/
  function merge(intervals: number[][]): number[][] {
    const seg = new SegmentSet()
    intervals.forEach(v => {
      seg.insert(v[0], v[1])
    })
    return seg.getAll().map(v => [v.left, v.right])
  }

  // 57. 插入区间
  // https://leetcode.cn/problems/insert-interval/
  function insert(intervals: number[][], newInterval: number[]): number[][] {
    const seg = new SegmentSet()
    intervals.forEach(v => {
      seg.insert(v[0], v[1])
    })
    seg.insert(newInterval[0], newInterval[1])
    return seg.getAll().map(v => [v.left, v.right])
  }

  // 1272. 删除区间
  // https://leetcode.cn/problems/remove-interval/
  function removeInterval(intervals: number[][], toBeRemoved: number[]): number[][] {
    const seg = new SegmentSet()
    intervals.forEach(v => {
      seg.insert(v[0], v[1])
    })
    seg.discard(toBeRemoved[0], toBeRemoved[1])
    return seg.getAll().map(v => [v.left, v.right])
  }

  const ss = new SegmentSet()
  ss.insert(1, 3)
  ss.insert(2, 4)
  ss.insert(5, 7)
  ss.insert(6, 8)
  ss.enumerateRange(2, 6, v => {
    console.log(v)
  })
  console.log(ss.nextStart(4), ss.toString())
  console.log(ss.prevStart(99), ss.toString())
}

export { SegmentSet }
