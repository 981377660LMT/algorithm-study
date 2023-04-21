type NestedArray = Array<number | NestedArray>

/**
 * 任意维度的树状数组.
 * 支持单点修改, 区间查询.
 */
class BITMultiDimension {
  private readonly _dim: number[]
  private readonly _data: number[]

  constructor(dimension: number[]) {
    let n = 1
    dimension.forEach(v => {
      n *= v
    })

    this._dim = dimension.slice()
    this._data = Array(n).fill(0)
  }

  /**
   * 0<=indices[i]<dimension[i]
   */
  add(indices: number[], x: number): void {
    this._addRec(indices, 0, 0, x)
  }

  /**
   * 0<=indices[i]<dimension[i]
   */
  query(indices: number[]): number {
    return this._queryRec(indices, 0, 0)
  }

  /**
   * 0<=a[i]<=b[i]<=dimension[i]
   */
  queryRange(a: number[], b: number[]): number {
    const t = Array(a.length)
    return this._queryRangeRec(0, a, b, t)
  }

  private _addRec(indices: number[], k: number, t: number, x: number): void {
    const d = this._dim[k]
    t *= d
    if (k + 1 === this._dim.length) {
      for (let i = indices[k]; i < d; i |= i + 1) {
        this._data[t + i] += x
      }
    } else {
      for (let i = indices[k]; i < d; i |= i + 1) {
        this._addRec(indices, k + 1, t + i, x)
      }
    }
  }

  private _queryRec(indices: number[], k: number, t: number): number {
    const d = this._dim[k]
    t *= d
    let res = 0
    if (k + 1 === this._dim.length) {
      for (let i = indices[k]; i > 0; i -= i & -i) {
        res += this._data[t + i - 1]
      }
    } else {
      for (let i = indices[k]; i > 0; i -= i & -i) {
        res += this._queryRec(indices, k + 1, t + i - 1)
      }
    }
    return res
  }

  private _queryRangeRec(d: number, a: number[], b: number[], t: number[]): number {
    if (d === this._dim.length) {
      return this._queryRec(t, 0, 0)
    }
    let res = 0
    t[d] = b[d]
    res += this._queryRangeRec(d + 1, a, b, t)
    t[d] = a[d]
    res -= this._queryRangeRec(d + 1, a, b, t)
    return res
  }
}

export { BITMultiDimension }

if (require.main === module) {
  const bit = new BITMultiDimension([10, 10])
  bit.add([1, 1], 1)
  console.log(bit.query([2, 2]))
}
