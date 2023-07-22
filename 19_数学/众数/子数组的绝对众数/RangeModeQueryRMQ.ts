import { SegmentTreePointUpdateRangeQuery } from '../../../6_tree/线段树/template/atcoder_segtree/SegmentTreePointUpdateRangeQuery'
import { bisectLeft, bisectRight } from '../../../9_排序和搜索/二分/bisect'

type E = { value: number; freqDiff: number }

/**
 * 查询区间`绝对众数`.
 */
class RangeMajorityQueryRMQ {
  private readonly _querier: { query: (start: number, end: number) => E }
  private readonly _pos: Map<number, number[]> = new Map()

  constructor(arr: number[]) {
    const leaves: E[] = Array(arr.length)
    for (let i = 0; i < arr.length; ++i) {
      const v = arr[i]
      leaves[i] = { value: v, freqDiff: 1 }
      !this._pos.has(v) && this._pos.set(v, [])
      this._pos.get(v)!.push(i)
    }
    this._querier = new SegmentTreePointUpdateRangeQuery<E>(
      leaves,
      () => ({ value: 0, freqDiff: 0 }),
      // !op函数是摩尔投票
      (a, b) => {
        const { value: aVal, freqDiff: aDiff } = a
        const { value: bVal, freqDiff: bDiff } = b
        if (aVal === bVal) return { value: aVal, freqDiff: aDiff + bDiff }
        if (aDiff > bDiff) return { value: aVal, freqDiff: aDiff - bDiff }
        return { value: bVal, freqDiff: bDiff - aDiff }
      }
    )
  }

  /**
   * O(log(n)) 查询区间 [start, end) 中的`绝对众数`以及其出现次数.
   * 如果不存在,返回`undefined`.
   * @param threshold 阈值.绝对众数的定义为出现次数大于等于`threshold`的数.
   */
  query(
    start: number,
    end: number,
    threshold: number
  ): [majority: number, freq: number] | undefined {
    const len = end - start
    if (threshold <= len >>> 1) {
      throw new Error('threshold must be greater than half of the length of the array')
    }
    const modeCandidate = this._querier.query(start, end).value
    const indices = this._pos.get(modeCandidate)
    if (!indices) return undefined
    const freq = bisectRight(indices, end - 1) - bisectLeft(indices, start)
    return freq >= threshold ? [modeCandidate, freq] : undefined
  }
}

export { RangeMajorityQueryRMQ }

if (require.main === module) {
  class MajorityChecker {
    private readonly _rmq: RangeMajorityQueryRMQ
    constructor(arr: number[]) {
      this._rmq = new RangeMajorityQueryRMQ(arr)
    }

    query(left: number, right: number, threshold: number): number {
      const res = this._rmq.query(left, right + 1, threshold)
      return res ? res[0] : -1
    }
  }

  /**
   * Your MajorityChecker object will be instantiated and called as such:
   * var obj = new MajorityChecker(arr)
   * var param_1 = obj.query(left,right,threshold)
   */
}
