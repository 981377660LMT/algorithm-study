/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable class-methods-use-this */

// API:
//  insert(left, right)      向区间集合中插入一个区间.
//  erase(left, right)       从区间集合中删除一个区间.
//  nextStart(x)             返回第一个大于等于x的区间起点.
//  prevStart(x)             返回最后一个小于等于x的区间起点.
//  ceiling(x)               返回区间内第一个大于等于x的元素.
//  floor(x)                 返回区间内第一个小于等于x的元素.
//  getRange(x)              返回包含x的区间.
//  has(x)                   判断x是否在区间集合中.
//  hasRange(left, right)    判断[left, right]是否在区间集合中.
//  at(i)                    返回第i个区间.
//  getAll()                 返回所有区间.
//  islice(min,max)          返回区间集合中包含在[min,max]区间内的所有区间的迭代器.
//  enumerate(min,max,f)     遍历SegmentSet中包含在[min,max]区间内的所有区间范围.
//  length                   SegmentSet中区间的个数.
//  count                    SegmentSet中区间的元素数量.

import { SortedList } from '../离线查询/根号分治/SortedList/SortedList'

const INF = 2e15

/**
 * 管理区间的数据结构.
 * @description
 * 1.所有区间都是`闭区间` 例如 [1,1] 表示 长为1的区间,起点为1;
 * !2.有交集的区间会被合并,例如 [1,2]和[2,3]会被合并为[1,3].
 */
class SegmentSet {
  private readonly _sl = new SortedList<[left: number, right: number]>((a, b) => a[0] - b[0])
  private _count = 0

  /**
   * 向区间集合中插入一个闭区间.
   */
  insert(left: number, right: number): void {
    if (left > right) {
      return
    }
    let it1 = this._sl.bisectRight([left, INF])
    const it2 = this._sl.bisectRight([right, INF])
    if (it1 > 0 && this._sl.at(it1 - 1)![1] >= left) {
      it1--
    }
    if (it1 !== it2) {
      const tmp1 = this._sl.at(it1)![0]
      left = Math.min(left, tmp1)
      const tmp2 = this._sl.at(it2 - 1)![1]
      right = Math.max(right, tmp2)
      let removed = 0
      this._sl.enumerate(
        it1,
        it2,
        v => {
          removed += v[1] - v[0] + 1
        },
        true
      )
      this._count -= removed
    }
    this._sl.add([left, right])
    this._count += right - left + 1
  }

  /**
   * 从区间集合中删除一个闭区间.
   */
  erase(left: number, right: number): void {
    if (left > right) {
      return
    }
    let it1 = this._sl.bisectRight([left, INF])
    const it2 = this._sl.bisectRight([right, INF])
    if (it1 > 0 && this._sl.at(it1 - 1)![1] >= left) {
      it1--
    }
    if (it1 === it2) {
      return
    }
    let nl = Math.min(left, this._sl.at(it1)![0])
    let nr = Math.max(right, this._sl.at(it2 - 1)![1])
    let removed = 0
    this._sl.enumerate(
      it1,
      it2,
      v => {
        removed += v[1] - v[0] + 1
      },
      true
    )
    this._count -= removed
    if (nl < left) {
      this._sl.add([nl, left])
      this._count += left - nl + 1
    }
    if (nr > right) {
      this._sl.add([right, nr])
      this._count += nr - right + 1
    }
  }

  /**
   * 返回第一个大于等于x的区间起点.
   */
  nextStart(x: number): number | undefined {
    const it = this._sl.bisectLeft([x, -INF])
    if (it === this._sl.length) {
      return undefined
    }
    return this._sl.at(it)![0]
  }

  /**
   * 返回最后一个小于等于x的区间起点.
   */
  prevStart(x: number): number | undefined {
    const it = this._sl.bisectRight([x, INF]) - 1
    if (it < 0) {
      return undefined
    }
    return this._sl.at(it)![0]
  }

  /**
   * 返回区间内第一个大于等于x的元素.
   */
  floor(x: number): number | undefined {
    const it = this._sl.bisectRight([x, INF])
    if (it === 0) {
      return undefined
    }
    return Math.min(x, this._sl.at(it - 1)![1])
  }

  /**
   * 返回区间内第一个小于等于x的元素.
   */
  ceiling(x: number): number | undefined {
    const it = this._sl.bisectRight([x, INF])
    if (it > 0 && this._sl.at(it - 1)![1] >= x) {
      return x
    }
    if (it !== this._sl.length) {
      return this._sl.at(it)![0]
    }
    return undefined
  }

  /**
   * 返回包含x的区间.
   */
  getInterval(x: number): [left: number, right: number] | undefined {
    const it = this._sl.bisectRight([x, INF])
    if (it === 0 || this._sl.at(it - 1)![1] < x) {
      return undefined
    }
    return this._sl.at(it - 1)!
  }

  /**
   * 判断x是否在区间集合中.
   */
  includes(x: number): boolean {
    const it = this._sl.bisectRight([x, INF])
    return it > 0 && this._sl.at(it - 1)![1] >= x
  }

  /**
   * 判断区间[left,right]是否在区间集合中.
   */
  includesInterval(left: number, right: number): boolean {
    if (left > right) {
      return false
    }
    const it1 = this._sl.bisectRight([left, INF])
    if (it1 === 0) {
      return false
    }
    const it2 = this._sl.bisectRight([right, INF])
    if (it1 !== it2) {
      return false
    }
    return this._sl.at(it1 - 1)![1] >= right
  }

  /**
   * 返回第index个区间.
   */
  at(index: number): [left: number, right: number] | undefined {
    return this._sl.at(index)
  }

  getAll(): [left: number, right: number][] {
    const res: [left: number, right: number][] = Array(this._sl.length).fill(0)
    this._sl.forEach((v, i) => {
      res[i] = v
    })
    return res
  }

  /**
   * 返回一个迭代器, 遍历 SegmentSet 中在 `[min,max]` 内的所有区间范围.
   */
  *irange(min: number, max: number): IterableIterator<[left: number, right: number]> {
    if (min > max) {
      return
    }
    let it = this._sl.bisectRight([min, INF]) - 1
    if (it < 0) it++
    const islice = this._sl.slice(it, this._sl.length)
    for (const v of islice) {
      if (v[0] > max) {
        return
      }
      yield [Math.max(v[0], min), Math.min(v[1], max)]
    }
  }

  /**
   * 遍历 SegmentSet 中在 `[min,max]` 内的所有区间范围.
   */
  enumerateRange(
    min: number,
    max: number,
    f: (interval: [left: number, right: number]) => void
  ): void {
    if (min > max) {
      return
    }
    let it = this._sl.bisectRight([min, INF]) - 1
    if (it < 0) it++
    const islice = this._sl.slice(it, this._sl.length)
    for (const v of islice) {
      if (v[0] > max) {
        break
      }
      f([Math.max(v[0], min), Math.min(v[1], max)])
    }
  }

  toString(): string {
    return `SegmentSet(${this.getAll().join(', ')})`
  }

  *[Symbol.iterator](): IterableIterator<[left: number, right: number]> {
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
  class CountIntervals {
    private readonly _ss = new SegmentSet()

    add(left: number, right: number): void {
      this._ss.insert(left, right)
    }

    count(): number {
      return this._ss.count
    }
  }

  class SummaryRanges {
    private readonly _ss = new SegmentSet()

    addNum(value: number): void {
      this._ss.insert(value, value + 1)
    }

    getIntervals(): number[][] {
      return this._ss.getAll().map(v => [v[0], v[1] - 1])
    }
  }

  // https://leetcode.cn/problems/range-module/submissions/
  class RangeModule {
    private readonly _ss = new SegmentSet()

    addRange(left: number, right: number): void {
      this._ss.insert(left, right)
    }

    queryRange(left: number, right: number): boolean {
      const intervals: any[] = []
      this._ss.enumerateRange(left, right, v => {
        intervals.push(v)
      })
      return intervals.length === 1 && intervals[0][0] === left && intervals[0][1] === right
    }

    removeRange(left: number, right: number): void {
      this._ss.erase(left, right)
    }
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
