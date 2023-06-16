// 1 <= val, inc, m <= 100
// 0 <= idx <= 1e5
// 总共最多会有 1e5 次对 append，addAll，multAll 和 getIndex 的调用。
// !https://leetcode.cn/problems/fancy-sequence/

import { SegmentTreeRangeUpdateRangeQuery } from '../SegmentTreeRangeUpdateRangeQuery'
import { createRangeAffineRangeSum } from '../SegmentTreeUtils'

const BIGMOD = BigInt(1e9 + 7)

class Fancy {
  private readonly _seg: SegmentTreeRangeUpdateRangeQuery<
    [size: bigint, sum: bigint],
    [mul: bigint, add: bigint]
  > = createRangeAffineRangeSum(1e5 + 10, BIGMOD)
  private _length = 0

  append(val: number): void {
    this._seg.update(this._length, this._length + 1, [1n, BigInt(val)])
    this._length++
  }

  addAll(inc: number): void {
    this._seg.update(0, this._length, [1n, BigInt(inc)])
  }

  multAll(m: number): void {
    this._seg.update(0, this._length, [BigInt(m), 0n])
  }

  getIndex(idx: number): number {
    if (idx >= this._length) return -1
    return Number(this._seg.get(idx)[1])
  }
}

export {}
