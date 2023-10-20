/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */
/* eslint-disable no-param-reassign */

import { discretizeSparse } from '../../../../22_专题/前缀与差分/差分数组/离散化/discretize'

/**
 * 单点更新, 区间查询，维护阿贝尔群(或者幺半群)的树状数组.
 * @alias BITMonoid
 */
class BitAbleGroup<E> {
  private _data!: E[]
  private _total!: E
  private readonly _n: number
  private readonly _e: () => E
  private readonly _op: (e1: E, e2: E) => E
  private readonly _inv: (e: E) => E

  /**
   * @param nOrArr 长度或者数组.
   * @param ableGroup 树状数组维护的阿贝尔群.
   * 如果只查询前缀和(不使用区间查询), 可以不指定逆元`inv`.
   * 默认为加法群.
   */
  constructor(nOrArr: number | E[], ableGroup?: { e: () => E; op: (e1: E, e2: E) => E; inv?: (e: E) => E }) {
    let defaultE = () => 0 as any
    let defaultOp = (e1: any, e2: any) => e1 + e2
    let defaultInv = (e: E) => -e as any
    if (ableGroup) {
      defaultE = ableGroup.e
      defaultOp = ableGroup.op
      defaultInv = ableGroup.inv || defaultInv
    }

    let defaultArray: E[]
    if (typeof nOrArr === 'number') {
      const base = Array(nOrArr)
      for (let i = 0; i < nOrArr; i++) base[i] = defaultE()
      defaultArray = base
    } else {
      defaultArray = nOrArr
    }

    this._n = defaultArray.length
    this._e = defaultE
    this._op = defaultOp
    this._inv = defaultInv
    this.build(defaultArray)
  }

  /**
   * 对arr进行修改，构建树状数组.
   */
  build(arr: E[]) {
    const n = arr.length
    for (let i = 1; i <= n; i++) {
      const j = i + (i & -i)
      if (j <= n) arr[j - 1] = this._op(arr[i - 1], arr[j - 1])
    }
    this._data = arr
    this._total = this.queryPrefix(n)
  }

  queryAll(): E {
    return this._total
  }

  /**
   * [0, end)的和.
   */
  queryPrefix(end: number): E {
    if (end > this._n) end = this._n
    let res = this._e()
    while (end > 0) {
      res = this._op(res, this._data[end - 1])
      end &= end - 1
    }
    return res
  }

  /**
   * [start, end)的和.
   */
  queryRange(start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start === 0) return this.queryPrefix(end)
    if (start > end) return this._e()
    let pos = this._e()
    let neg = this._e()
    while (end > start) {
      pos = this._op(pos, this._data[end - 1])
      end &= end - 1
    }
    while (start > end) {
      neg = this._op(neg, this._data[start - 1])
      start &= start - 1
    }
    return this._op(pos, this._inv(neg))
  }

  update(index: number, lazy: E) {
    this._total = this._op(this._total, lazy)
    for (index++; index <= this._n; index += index & -index) {
      this._data[index - 1] = this._op(this._data[index - 1], lazy)
    }
  }

  /**
   * 返回最大的 end 使得 `check(QueryPrefix(end)) == true`.
   * @param check check(preSum, end): preSum 对应的是 [0, end) 的和.
   *
   * @example
   * ```ts
   * const fw = new BitGroup(10)
   * fw.maxRight(preSum => preSum <= 10)
   */
  maxRight(check: (preSum: E, end: number) => boolean): number {
    let i = 0
    let cur = this._e()
    let k = 1
    while (k << 1 <= this._n) k <<= 1
    while (k > 0) {
      if (i + k - 1 < this._n) {
        const t = this._op(cur, this._data[i + k - 1])
        if (check(t, i + k)) {
          i += k
          cur = t
        }
      }
      k >>= 1
    }
    return i
  }

  toString(): string {
    const res: string[] = []
    for (let i = 0; i < this._n; i++) {
      res.push(`${this.queryRange(i, i + 1)}`)
    }
    return `FenwickTree[${res}]`
  }
}

export { BitAbleGroup, BitAbleGroup as BITMonoid }

if (require.main === module) {
  // 最长上传前缀
  // https://leetcode.cn/problems/longest-uploaded-prefix/
  class LUPrefix {
    private readonly _bit: BitAbleGroup<number>
    constructor(n: number) {
      this._bit = new BitAbleGroup(n)
    }

    upload(video: number): void {
      this._bit.update(video, 1)
    }

    longest(): number {
      return this._bit.maxRight((preSum, end) => preSum >= end)
    }
  }
}
