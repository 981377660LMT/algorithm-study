/* eslint-disable no-inner-declarations */
/* eslint-disable eqeqeq */

import {
  ISortedList,
  SortedListFast
} from '../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

class MaxPrefixQueryWithInsertionsOnly {
  private readonly _queryMax: boolean
  private readonly _slPrefix: ISortedList<{ x: number; y: number }>

  /**
   * @param queryMax 是否查询最大值.默认为true.
   */
  constructor(queryMax = true) {
    this._queryMax = queryMax
    this._slPrefix = new SortedListFast((p1, p2) => p2.x - p1.x)
  }

  /**
   * 将点`(x, y)`加入数据结构中.
   * @complexity O(logn)
   */
  add(x: number, y: number): void {
    this._addPrefix(x, y)
  }

  /**
   * 对`x<=rightX`的所有点`(x, y)`，找到最大的`y`值.
   * @complexity O(logn)
   */
  queryPrefix(rightX: number): number | undefined {
    if (!this._slPrefix.length) return undefined
    const pos = this._slPrefix.bisectLeft({ x: rightX, y: 0 })
    const point = this._slPrefix.at(pos)
    if (point == undefined) return undefined
    return this._queryMax ? point.y : -point.y
  }

  private _addPrefix(x: number, y: number): void {
    if (!this._queryMax) y = -y
    const newPoint = { x, y }
    let pos = this._slPrefix.bisectLeft(newPoint)
    const point = this._slPrefix.at(pos)
    if (point != undefined && point.y >= y) return
    this._slPrefix.add(newPoint)
    this._slPrefix.at(pos)!.y = y
    while (pos > 0) {
      pos--
      if (this._slPrefix.at(pos)!.y > y) break
      this._slPrefix.pop(pos)
    }
  }
}

/**
 * 维护二维平面上的点集，支持以下操作：
 * - `add(x, y)`：将点`(x, y)`加入数据结构中.
 * - `querySuffix(leftX)`：对`x>=leftX`的所有点`(x, y)`，找到最大(小)的`y`值.
 * @summary
 * 对于点对 `(a1,b1)` 和 `(a2,b2)`，如果 `a1<=a2,b1<=b2`，则可以将 `(a1,b1)` 删除.
 * 因此map中的点 `(ai,bi)` 必定满足 `ai<ai+1,bi>bi+1`.
 * 查询时只需二分找到不小于 `leftX` 的点，然后返回其 `y` 值即可.
 * @link https://usaco.guide/adv/springboards?lang=cpp
 */
class MaxSuffixQueryWithInsertionsOnly {
  private readonly _queryMax: boolean
  private readonly _slSuffix: ISortedList<{ x: number; y: number }>

  /**
   * @param queryMax 是否查询最大值.默认为true.
   */
  constructor(queryMax = true) {
    this._queryMax = queryMax
    this._slSuffix = new SortedListFast((p1, p2) => p1.x - p2.x)
  }

  /**
   * 将点`(x, y)`加入数据结构中.
   * @complexity O(logn)
   */
  add(x: number, y: number): void {
    this._addSuffix(x, y)
  }

  /**
   * 对`x>=leftX`的所有点`(x, y)`，找到最大的`y`值.
   * @complexity O(logn)
   */
  querySuffix(leftX: number): number | undefined {
    if (!this._slSuffix.length) return undefined
    const pos = this._slSuffix.bisectLeft({ x: leftX, y: 0 })
    const point = this._slSuffix.at(pos)
    if (point == undefined) return undefined
    return this._queryMax ? point.y : -point.y
  }

  private _addSuffix(x: number, y: number): void {
    if (!this._queryMax) y = -y
    const newPoint = { x, y }
    let pos = this._slSuffix.bisectLeft(newPoint)
    const point = this._slSuffix.at(pos)
    if (point != undefined && point.y >= y) return
    this._slSuffix.add(newPoint)
    this._slSuffix.at(pos)!.y = y
    while (pos > 0) {
      pos--
      if (this._slSuffix.at(pos)!.y > y) break
      this._slSuffix.pop(pos)
    }
  }
}

export { MaxPrefixQueryWithInsertionsOnly, MaxSuffixQueryWithInsertionsOnly }

if (require.main === module) {
  // 2907.最大递增三元组的和
  // https://leetcode.cn/problems/maximum-profitable-triplets-with-increasing-prices-i/description/
  // 找到三个下标 i, j, k,使得 i < j < k 且 prices[i] < prices[j] < prices[k],
  // 并且 profits[i] + profits[j] + profits[k] 最大。
  // 如果无法找到则返回 -1。
  // !三元组:枚举中间.
  function maxProfit(prices: number[], profits: number[]): number {
    const n = prices.length

    const pre = new MaxPrefixQueryWithInsertionsOnly()
    const leftMax = Array<number>(n).fill(0)
    for (let i = 0; i < n; i++) {
      const curX = prices[i]
      const curY = profits[i]
      pre.add(curX, curY)
      const tmp = pre.queryPrefix(curX - 1)
      if (tmp == undefined) continue
      leftMax[i] = tmp
    }

    const suf = new MaxSuffixQueryWithInsertionsOnly()
    const rightMax = Array<number>(n).fill(0)
    for (let i = n - 1; i >= 0; i--) {
      const curX = prices[i]
      const curY = profits[i]
      suf.add(curX, curY)
      const tmp = suf.querySuffix(curX + 1)
      if (tmp == undefined) continue
      rightMax[i] = tmp
    }

    let res = -1
    for (let i = 0; i < n; i++) {
      if (leftMax[i] == 0 || rightMax[i] == 0) continue
      res = Math.max(res, profits[i] + leftMax[i] + rightMax[i])
    }

    return res
  }

  // console.log(maxProfit([10, 2, 3, 4], [100, 2, 7, 10]))
  // console.log(maxProfit([1, 2, 3, 4, 5], [1, 5, 3, 4, 6]))

  test()
  testTime()
  function test(): void {
    class Mocker {
      private readonly _points: { x: number; y: number }[] = []
      private readonly _qeuryMax: boolean
      constructor(queryMax: boolean) {
        this._qeuryMax = queryMax
      }
      add(x: number, y: number): void {
        this._points.push({ x, y })
      }

      querySuffix(leftX: number): number | undefined {
        let res: number | undefined
        if (this._qeuryMax) {
          for (const p of this._points) {
            if (p.x >= leftX && (res == undefined || p.y > res)) res = p.y
          }
        } else {
          for (const p of this._points) {
            if (p.x >= leftX && (res == undefined || p.y < res)) res = p.y
          }
        }
        return res
      }

      queryPrefix(rightX: number): number | undefined {
        let res: number | undefined
        if (this._qeuryMax) {
          for (const p of this._points) {
            if (p.x <= rightX && (res == undefined || p.y > res)) res = p.y
          }
        } else {
          for (const p of this._points) {
            if (p.x <= rightX && (res == undefined || p.y < res)) res = p.y
          }
        }
        return res
      }

      toString(): string {
        return this._points
          .map(({ x, y }) => `(${x},${y})`)
          .join(' ')
          .trim()
      }
    }

    const queryMax = Math.random() > 0.5
    console.log(`queryMax=${queryMax}`)
    const mocker = new Mocker(queryMax)
    const msq = new MaxSuffixQueryWithInsertionsOnly(queryMax)
    const mpq = new MaxPrefixQueryWithInsertionsOnly(queryMax)
    const n = 1e4
    for (let i = 0; i < n; ++i) {
      const x = Math.floor(Math.random() * n)
      const y = Math.floor(Math.random() * n)
      mocker.add(x, y)
      msq.add(x, y)
      mpq.add(x, y)

      const leftX = Math.floor(Math.random() * n)
      const res1 = mocker.querySuffix(leftX)
      const res2 = msq.querySuffix(leftX)
      if (res1 != res2) {
        console.error(`Error: res1=${res1}, res2=${res2}`)
        console.error(`mocker=${mocker.toString()}`)
        console.error(`msq=${msq.toString()}`)
        throw new Error()
      }
      const rightX = Math.floor(Math.random() * n)
      const res3 = mocker.queryPrefix(rightX)
      const res4 = mpq.queryPrefix(rightX)
      if (res3 != res4) {
        console.error(`Error: res3=${res3}, res4=${res4}`)
        console.error(`mocker=${mocker.toString()}`)
        console.error(`msq=${msq.toString()}`)
        throw new Error()
      }
    }

    console.log('pass')
  }

  function testTime(): void {
    const msq = new MaxSuffixQueryWithInsertionsOnly()
    const n = 2e5
    const m = 2e5
    const points: [x: number, y: number][] = []
    const leftXs: number[] = []
    for (let i = 0; i < n; ++i) {
      const x = Math.floor(Math.random() * m)
      const y = Math.floor(Math.random() * m)
      points.push([x, y])
      leftXs.push(Math.floor(Math.random() * m))
    }
    points.sort((p1, p2) => p2[0] - p1[0])

    console.time('querySuffix')
    for (let i = 0; i < n; ++i) {
      msq.add(points[i][0], points[i][1])
      msq.querySuffix(leftXs[i])
    }
    console.timeEnd('querySuffix')
  }
}
