/* eslint-disable no-inner-declarations */
/* eslint-disable no-cond-assign */
/* eslint-disable no-param-reassign */

import { discretizeSparse } from '../../../../22_专题/前缀与差分/差分数组/离散化/discretize'

// !单点修改+区间查询

const INF = 2e15

class SegmentTreePointUpdateRangeQuery<E = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _data: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  /**
   * 单点更新,区间查询的线段树.
   * @param nOrLeaves 大小或叶子节点的值.
   * @param e 幺元.
   * @param op 结合律.
   */
  constructor(nOrLeaves: number | ArrayLike<E>, e: () => E, op: (a: E, b: E) => E) {
    const n = typeof nOrLeaves === 'number' ? nOrLeaves : nOrLeaves.length
    let size = 1
    while (size < n) size <<= 1
    const data = Array(size << 1)
    for (let i = 0; i < data.length; i++) data[i] = e()

    this._n = n
    this._size = size
    this._data = data
    this._e = e
    this._op = op

    if (typeof nOrLeaves !== 'number') this.build(nOrLeaves)
  }

  set(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._data[index] = value
    while ((index >>= 1)) {
      this._data[index] = this._op(this._data[index << 1], this._data[(index << 1) | 1])
    }
  }

  get(index: number): E {
    if (index < 0 || index >= this._n) return this._e()
    return this._data[index + this._size]
  }

  /**
   * 将`index`处的值与作用素`value`结合.
   */
  update(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._data[index] = this._op(this._data[index], value)
    while ((index >>= 1)) {
      this._data[index] = this._op(this._data[index << 1], this._data[(index << 1) | 1])
    }
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return this._e()

    let leftRes = this._e()
    let rightRes = this._e()
    for (start += this._size, end += this._size; start < end; start >>= 1, end >>= 1) {
      if (start & 1) leftRes = this._op(leftRes, this._data[start++])
      if (end & 1) rightRes = this._op(this._data[--end], rightRes)
    }
    return this._op(leftRes, rightRes)
  }

  queryAll(): E {
    return this._data[1]
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (value: E) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    let res = this._e()
    while (true) {
      while (!(start & 1)) start >>= 1
      if (!predicate(this._op(res, this._data[start]))) {
        while (start < this._size) {
          start <<= 1
          const tmp = this._op(res, this._data[start])
          if (predicate(tmp)) {
            res = tmp
            start++
          }
        }
        return start - this._size
      }
      res = this._op(res, this._data[start])
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (value: E) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    let res = this._e()
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      if (!predicate(this._op(this._data[end], res))) {
        while (end < this._size) {
          end = (end << 1) | 1
          const tmp = this._op(this._data[end], res)
          if (predicate(tmp)) {
            res = tmp
            end--
          }
        }
        return end + 1 - this._size
      }
      res = this._op(this._data[end], res)
      if ((end & -end) === end) break
    }
    return 0
  }

  build(arr: ArrayLike<E>): void {
    if (arr.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < arr.length; i++) {
      this._data[i + this._size] = arr[i] // 叶子结点
    }
    for (let i = this._size - 1; i > 0; i--) {
      this._data[i] = this._op(this._data[i << 1], this._data[(i << 1) | 1])
    }
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreePointUpdateRangeQuery(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }
}

export { SegmentTreePointUpdateRangeQuery }

if (require.main === module) {
  const seg = new SegmentTreePointUpdateRangeQuery(
    10,
    () => 0,
    (a, b) => a + b
  )
  console.log(seg.toString())
  seg.set(0, 1)
  seg.set(1, 2)
  console.log(seg.toString())
  seg.update(3, 4)
  console.log(seg.toString())
  console.log(seg.query(0, 4))
  seg.build([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(seg.toString())
  console.log(seg.minLeft(10, x => x < 15))
  console.log(seg.maxRight(0, x => x <= 15))
  console.log(seg.queryAll())

  benchMark()
  function benchMark(): void {
    const n = 2e5
    const seg = new SegmentTreePointUpdateRangeQuery<number>(
      n,
      () => 0,
      (parent, child) => parent + child
    )
    console.time('update')
    for (let i = 0; i < n; i++) {
      seg.update(i, i)
      seg.query(0, i)
    }
    console.timeEnd('update')
  }

  // https://leetcode.cn/problems/maximum-sum-queries/
  // 2736. 最大和查询 (二维偏序+离线查询)
  // 对于第 i 个查询，在所有满足 nums1[j] >= xi 且 nums2[j] >= yi 的下标 j (0 <= j < n) 中，
  // 找出 nums1[j] + nums2[j] 的 最大值 ，
  // 如果不存在满足条件的 j 则返回 -1 。
  // 返回数组 answer ，其中 answer[i] 是第 i 个查询的答案。
  //
  // !即:对每个查询(x,y),求出右上角的点的`横坐标+纵坐标`的最大值
  function maximumSumQueries(nums1: number[], nums2: number[], queries: number[][]): number[] {
    const points = nums1.map((v, i) => [v, nums2[i]]).sort((a, b) => a[0] - b[0] || a[1] - b[1])
    const qWithId = queries.map((q, i) => [q[0], q[1], i]).sort((a, b) => a[0] - b[0] || a[1] - b[1])

    const allY = new Set(nums2)
    queries.forEach(q => allY.add(q[1]))
    const [rank, count] = discretize([...allY])

    const seg = new SegmentTreePointUpdateRangeQuery<number>(count, () => -INF, Math.max)
    const res = Array(queries.length).fill(-1)
    let pi = points.length - 1
    for (let i = qWithId.length - 1; i >= 0; i--) {
      const [qx, qy, qid] = qWithId[i]
      while (pi >= 0 && points[pi][0] >= qx) {
        seg.update(rank(points[pi][1])!, points[pi][0] + points[pi][1])
        pi--
      }
      const curMax = seg.query(rank(qy)!, count)
      res[qid] = curMax === -INF ? -1 : curMax
    }

    return res
  }

  /**
   * (松)离散化.
   * @returns
   * rank: 给定一个数,返回它的排名`(0-count)`.
   * count: 离散化(去重)后的元素个数.
   */
  function discretize(nums: number[]): [rank: (num: number) => number, count: number] {
    const allNums = [...new Set(nums)].sort((a, b) => a - b)
    const rank = (num: number) => {
      let left = 0
      let right = allNums.length - 1
      while (left <= right) {
        const mid = (left + right) >>> 1
        if (allNums[mid] >= num) {
          right = mid - 1
        } else {
          left = mid + 1
        }
      }
      return left
    }
    return [rank, allNums.length]
  }

  // 2907. Maximum Profitable Triplets With Increasing Prices I
  // 找到三个下标 i, j, k,使得 i < j < k 且 prices[i] < prices[j] < prices[k],
  // 并且 profits[i] + profits[j] + profits[k] 最大。
  // 如果无法找到则返回 -1。
  //
  // !三元组:枚举中间.
  // 用树状数组更新和查询各节点左边和右边各个价格的最高利润。如果计算得到的利润为 0 则表明没有符合要求的利润，因此可以忽略该节点。
  function maxProfit(prices: number[], profits: number[]): number {
    const n = prices.length
    const [getRank, count] = discretizeSparse(prices)

    const leftMax = Array<number>(count).fill(0)
    const tree1 = new SegmentTreePointUpdateRangeQuery(count, () => 0, Math.max)
    for (let i = 0; i < n; i++) {
      const rank = getRank(prices[i])
      const max = tree1.query(0, rank)
      leftMax[i] = max
      tree1.update(rank, profits[i])
    }

    const rightMax = Array<number>(count).fill(0)
    const tree2 = new SegmentTreePointUpdateRangeQuery(count, () => 0, Math.max)
    for (let i = n - 1; ~i; i--) {
      const rank = getRank(prices[i])
      const max = tree2.query(rank + 1, count)
      rightMax[i] = max
      tree2.update(rank, profits[i])
    }

    let res = -1
    for (let i = 1; i < n - 1; i++) {
      if (leftMax[i] === 0 || rightMax[i] === 0) continue
      res = Math.max(res, leftMax[i] + rightMax[i] + profits[i])
    }
    return res
  }
}
