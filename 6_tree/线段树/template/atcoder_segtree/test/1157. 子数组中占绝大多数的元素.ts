// 线段树+摩尔投票维护区间众数

import { bisectLeft, bisectRight } from '../../../../../9_排序和搜索/二分/bisect'
import { SegmentTreePointUpdateRangeQuery } from '../SegmentTreePointUpdateRangeQuery'

type E = [value: number, freq: number]

class MajorityChecker {
  private readonly _freqQuerier: SegmentTreePointUpdateRangeQuery<E>
  private readonly _indexes: Map<number, number[]> = new Map()

  constructor(arr: number[]) {
    const leaves: E[] = Array(arr.length)
    arr.forEach((val, i) => {
      !this._indexes.has(val) && this._indexes.set(val, [])
      this._indexes.get(val)!.push(i)
      leaves[i] = [val, 1]
    })

    this._freqQuerier = new SegmentTreePointUpdateRangeQuery<E>(
      leaves,
      () => [0, 0],
      // !op函数是摩尔投票
      (a, b) => {
        const [aVal, aFreq] = a
        const [bVal, bFreq] = b
        if (aVal === bVal) return [aVal, aFreq + bFreq]
        if (aFreq > bFreq) return [aVal, aFreq - bFreq]
        return [bVal, bFreq - aFreq]
      }
    )
  }

  query(left: number, right: number, threshold: number): number {
    const modeCandicate = this._freqQuerier.query(left, right + 1)[0]
    const indexes = this._indexes.get(modeCandicate)
    if (!indexes) return -1
    const freq = bisectRight(indexes, right) - bisectLeft(indexes, left)
    return freq >= threshold ? modeCandicate : -1
  }
}

/**
 * Your MajorityChecker object will be instantiated and called as such:
 * var obj = new MajorityChecker(arr)
 * var param_1 = obj.query(left,right,threshold)
 */

export {}
