/**
 * 区间的LCA.
 * @link https://github.com/pranjalssh/CP_codes/blob/master/anta/!RangeLCA.cpp
 */
class RangeLCA {
  private readonly _seg: number[]
  private readonly _n: number
  private readonly _lcaImpl: (x: number, y: number) => number

  constructor(data: number[], lcaImpl: (x: number, y: number) => number) {
    let n = 1
    while (n < data.length) n <<= 1
    const seg = Array(n << 1)
    for (let i = 0; i < data.length; i++) seg[n + i] = data[i]
    for (let i = n - 1; ~i; i--) seg[i] = this._lca(seg[i << 1], seg[(i << 1) | 1])
    this._seg = seg
    this._n = n
    this._lcaImpl = lcaImpl
  }

  // [start, end).
  query(start: number, end: number): number {
    let res = -1
    for (; start && start + (start & -start) <= end; start += start & -start) {
      res = this._lca(res, this._seg[~~((this._n + start) / (start & -start))])
    }
    for (; start < end; end -= end & -end) {
      res = this._lca(res, this._seg[~~((this._n + end) / (end & -end)) - 1])
    }
    return res
  }

  private _lca(x: number, y: number): number {
    if (x === -1 || y === -1) return x === -1 ? y : x
    return this._lcaImpl(x, y)
  }
}

export {}
