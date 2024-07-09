/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { ErasableHeap } from '../../../8_heap/ErasableHeap'
import { ODT } from '../ODT-fastset'

/**
 * 维护相同元素的最长连续长度.
 */
class LongestRepeating<T> {
  private static readonly _NONE: any = Symbol('NONE')
  private readonly _lens = new ErasableHeap<number>((a, b) => a > b)
  private readonly _n: number
  private readonly _odt: ODT<T>

  constructor(arr: ArrayLike<T>) {
    const n = arr.length
    this._n = n
    this._odt = new ODT(n, LongestRepeating._NONE)

    let pre = 0
    let pos = 0
    while (pos < n) {
      const leader = arr[pos]
      pos++
      while (pos < n && arr[pos] === leader) pos++
      this._lens.push(pos - pre)
      this._odt.set(pre, pos, leader)
      pre = pos
    }
  }

  queryAll(): number {
    return this._lens.peek() || 0
  }

  /**
   * 查询`[start, end)`范围内最长的连续长度.
   * !复杂度与范围内的区间个数成正比.
   */
  query(start: number, end: number): number {
    let max = 0
    this._odt.enumerateRange(start, end, (s, e) => {
      max = Math.max(max, e - s)
    })
    return max
  }

  update(start: number, end: number, value: T): void {
    const leftStart = this._odt.get(start)![0]
    const rightEnd = this._odt.get(end - 1)![1]
    const leftSeg = this._odt.get(leftStart - 1)
    const rightSeg = this._odt.get(rightEnd)
    const first = leftSeg ? leftSeg[0] : 0
    const last = rightSeg ? rightSeg[1] : this._n
    this._odt.enumerateRange(first, last, (s, e) => {
      this._lens.remove(e - s)
    })
    this._odt.set(start, end, value)
    this._odt.enumerateRange(first, last, (s, e) => {
      this._lens.push(e - s)
    })
  }

  set(index: number, value: T): void {
    const [start, end] = this._odt.get(index)!
    const leftSeg = this._odt.get(start - 1)
    const rightSeg = this._odt.get(end)
    const first = leftSeg ? leftSeg[0] : 0
    const last = rightSeg ? rightSeg[1] : this._n
    this._odt.enumerateRange(first, last, (s, e) => {
      this._lens.remove(e - s)
    })
    this._odt.set(index, index + 1, value)
    this._odt.enumerateRange(first, last, (s, e) => {
      this._lens.push(e - s)
    })
  }
}

export { LongestRepeating }

if (require.main === module) {
  // https://leetcode.cn/problems/longest-substring-of-one-repeating-character/
  function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
    const L = new LongestRepeating(s)
    const q = queryCharacters.length
    const res = Array(q)
    for (let i = 0; i < q; i++) {
      const target = queryCharacters[i]
      const pos = queryIndices[i]
      L.set(pos, target)
      res[i] = L.queryAll()
    }
    return res
  }
}
