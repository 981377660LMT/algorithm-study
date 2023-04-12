/* eslint-disable no-constant-condition */

/**
 * 在线查询区间众数(出现次数最多的数以及出现的次数).
 */
class RangeModeQuery {
  private static _compress(arr: number[]): [sorted: Uint32Array, compressed: Uint32Array] {
    const set = new Set<number>(arr)
    const sorted = new Uint32Array(set.size)
    let index = 0
    set.forEach(v => {
      sorted[index++] = v
    })
    sorted.sort((a, b) => a - b)
    const rank = new Map<number, number>()
    sorted.forEach((v, i) => {
      rank.set(v, i)
    })
    const compressed = new Uint32Array(arr.length)
    for (let i = 0; i < arr.length; i++) {
      compressed[i] = rank.get(arr[i])!
    }
    return [sorted, compressed]
  }

  private readonly _value: Uint32Array
  private readonly _rank: Uint32Array
  private readonly _mode: Uint32Array[]
  private readonly _freq: Uint32Array[]
  private readonly _qs: number[][]
  private readonly _t: number
  private readonly _sorted: Uint32Array

  /**
   * O(nsqrt(n))构建.
   */
  constructor(arr: number[]) {
    const n = arr.length
    const [sorted, value] = RangeModeQuery._compress(arr)
    const t = Math.max(1, Math.ceil(Math.sqrt(n)))
    const rank = new Uint32Array(n)
    const qs: number[][] = Array(n)
    for (let i = 0; i < n; i++) {
      qs[i] = []
    }
    for (let i = 0; i < n; i++) {
      const e = value[i]
      rank[i] = qs[e].length
      qs[e].push(i)
    }

    const bc = ~~(n / t) + 1
    const mode: Uint32Array[] = Array(bc)
    const freq: Uint32Array[] = Array(bc)
    for (let i = 0; i < bc; i++) {
      mode[i] = new Uint32Array(bc)
      freq[i] = new Uint32Array(bc)
    }

    for (let f = 0; f * t <= n; f++) {
      const freq_ = new Uint32Array(n)
      let curMode = 0
      let curFreq = 0
      for (let l = f + 1; l * t <= n; l++) {
        for (let i = (l - 1) * t; i < l * t; i++) {
          const e = value[i]
          freq_[e]++
          if (freq_[e] > curFreq) {
            curMode = e
            curFreq = freq_[e]
          }
        }
        mode[f][l] = curMode
        freq[f][l] = curFreq
      }
    }

    this._value = value
    this._rank = rank
    this._mode = mode
    this._freq = freq
    this._qs = qs
    this._t = t
    this._sorted = sorted
  }

  /**
   * O(sqrt(n))查询区间 [start, end) 中出现次数最多的数mode, 以及出现的次数freq.
   */
  query(start: number, end: number): [mode: number, freq: number] {
    if (start >= end) {
      throw new Error(`invalid range: [${start}, ${end})`)
    }
    if (start < 0) start = 0
    if (end > this._value.length) end = this._value.length

    const T = this._t
    const bf = ~~(start / T) + 1
    const bl = ~~(end / T)
    if (bf >= bl) {
      let resMode = 0
      let resFreq = 0
      for (let i = start; i < end; i++) {
        const rank = this._rank[i]
        const value = this._value[i]
        const v = this._qs[value]
        const lenV = v.length
        while (true) {
          const idx = rank + resFreq
          if (idx >= lenV || v[idx] >= end) {
            break
          }
          resMode = value
          resFreq++
        }
      }
      return [this._sorted[resMode], resFreq]
    }

    let resMode = this._mode[bf][bl]
    let resFreq = this._freq[bf][bl]
    for (let i = start; i < bf * T; i++) {
      const rank = this._rank[i]
      const value = this._value[i]
      const v = this._qs[value]
      const lenV = v.length
      while (true) {
        const idx = rank + resFreq
        if (idx >= lenV || v[idx] >= end) {
          break
        }
        resMode = value
        resFreq++
      }
    }

    for (let i = bl * T; i < end; i++) {
      const rank = this._rank[i]
      const value = this._value[i]
      const v = this._qs[value]
      const lenV = v.length
      while (true) {
        const idx = rank - resFreq
        if (idx < 0 || idx >= lenV || v[idx] < start) {
          break
        }
        resMode = value
        resFreq++
      }
    }

    return [this._sorted[resMode], resFreq]
  }
}

if (require.main === module) {
  // https://leetcode.cn/problems/online-majority-element-in-subarray/submissions/
  class MajorityChecker {
    private readonly _rmq: RangeModeQuery
    constructor(arr: number[]) {
      this._rmq = new RangeModeQuery(arr)
    }

    query(left: number, right: number, threshold: number): number {
      const [mode, freq] = this._rmq.query(left, right + 1)
      return freq >= threshold ? mode : -1
    }
  }
}

export { RangeModeQuery }
