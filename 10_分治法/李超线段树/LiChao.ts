/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

type Line = {
  k: number
  b: number
}

const INF = 2e9 // 2e15

class LiChaoTree {
  private readonly _n: number
  private readonly _offset: number
  private readonly _lower: number
  private readonly _higher: number
  private readonly _compress: boolean
  private readonly _minimize: boolean
  private readonly _lineIds: Int32Array
  private readonly _xs: number[]
  private readonly _lines: Line[]
  private readonly _evaluate: (line: Line, x: number) => number

  /**
   * 指定查询的 x 值建立李超线段树，采用坐标压缩.
   */
  constructor(
    queryX: Iterable<number>,
    options?: {
      minimize?: boolean
      evaluate?: (line: Line, x: number) => number
    }
  )

  /**
   * 指定查询的 x 值范围建立李超线段树，不采用坐标压缩.
   * higher - lower <= 1e6.
   */
  constructor(
    lower: number,
    higher: number,
    options?: {
      minimize?: boolean
      evaluate?: (line: Line, x: number) => number
    }
  )

  constructor(arg1: any, arg2?: any, arg3?: any) {
    const shouldCompress = typeof arg1 !== 'number'
    if (shouldCompress) {
      const { minimize = true, evaluate = (line: Line, x: number) => line.k * x + line.b } = arg2 || {}
      const allNums = [...new Set(arg1 as Iterable<number>)].sort((a, b) => a - b)
      const n = allNums.length
      let log = 1
      while (1 << log < n) log++
      const offset = 1 << log
      const lineIds = new Int32Array(offset << 1).fill(-1)
      this._n = n
      this._offset = offset
      this._lower = 0
      this._higher = 0
      this._compress = true
      this._minimize = minimize
      this._lineIds = lineIds
      this._xs = allNums
      this._lines = []
      this._evaluate = evaluate
    } else {
      const { minimize = true, evaluate = (line: Line, x: number) => line.k * x + line.b } = arg3 || {}
      const n = arg2 - arg1
      let log = 1
      while (1 << log < n) log++
      const offset = 1 << log
      const lineIds = new Int32Array(offset << 1).fill(-1)
      this._n = n
      this._offset = offset
      this._lower = arg1
      this._higher = arg2
      this._compress = false
      this._minimize = minimize
      this._lineIds = lineIds
      this._xs = []
      this._lines = []
      this._evaluate = evaluate
    }
  }

  /** O(logn). */
  addLine(line: Line): void {
    const id = this._lines.length
    this._lines.push(line)
    this._addLineAt(1, id)
  }

  /** [start, end). O(log^2n) */
  addSegment(startX: number, endX: number, line: Line): void {
    if (startX >= endX) return
    const id = this._lines.length
    this._lines.push(line)
    startX = this._getIndex(startX) + this._offset
    endX = this._getIndex(endX) + this._offset
    while (startX < endX) {
      if (startX & 1) this._addLineAt(startX++, id)
      if (endX & 1) this._addLineAt(--endX, id)
      startX >>>= 1
      endX >>>= 1
    }
  }

  /** O(logn). */
  query(x: number): { value: number; lineId: number } {
    x = this._getIndex(x)
    let pos = x + this._offset
    let resLineId = -1
    let resValue = this._minimize ? INF : -INF
    while (pos > 0) {
      const candId = this._lineIds[pos]
      if (candId !== -1 && candId !== resLineId) {
        const candValue = this._evaluateInner(candId, x)
        if (this._minimize ? candValue < resValue : candValue > resValue) {
          resValue = candValue
          resLineId = candId
        }
      }
      pos >>>= 1
    }
    return { value: resValue, lineId: resLineId }
  }

  private _addLineAt(i: number, fid: number): void {
    const upperBit = 31 - Math.clz32(i)
    let left = (this._offset >>> upperBit) * (i - (1 << upperBit))
    let right = left + (this._offset >>> upperBit)
    const minimize = this._minimize
    while (left < right) {
      const gid = this._lineIds[i]
      const fl = this._evaluateInner(fid, left)
      const fr = this._evaluateInner(fid, right - 1)
      const gl = this._evaluateInner(gid, left)
      const gr = this._evaluateInner(gid, right - 1)
      const bl = minimize ? fl < gl : fl > gl
      const br = minimize ? fr < gr : fr > gr
      if (bl && br) {
        this._lineIds[i] = fid
        return
      }
      if (!bl && !br) {
        return
      }
      const mid = (left + right) >>> 1
      const fm = this._evaluateInner(fid, mid)
      const gm = this._evaluateInner(gid, mid)
      const bm = minimize ? fm < gm : fm > gm
      if (bm) {
        this._lineIds[i] = fid
        fid = gid
        if (!bl) {
          i <<= 1
          right = mid
        } else {
          i = (i << 1) | 1
          left = mid
        }
      } else if (bl) {
        i <<= 1
        right = mid
      } else {
        i = (i << 1) | 1
        left = mid
      }
    }
  }

  private _evaluateInner(fid: number, x: number): number {
    if (fid === -1) return this._minimize ? INF : -INF
    const target = this._compress ? this._xs[Math.min(x, this._n - 1)] : x + this._lower
    return this._evaluate(this._lines[fid], target)
  }

  private _getIndex(x: number): number {
    if (this._compress) return LiChaoTree._lowerBound(this._xs, x)
    if (x < this._lower || x > this._higher) throw new RangeError(`x out of range: ${x}`)
    return x - this._lower
  }

  private static _lowerBound(arr: ArrayLike<number>, x: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] < x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }
}

export { LiChaoTree }

if (require.main === module) {
  checkWithBf()

  function checkWithBf(): void {
    class Mocker {
      private readonly _minimize: boolean
      private readonly _lines: { start: number; end: number; line: Line }[] = []

      constructor(minimize: boolean) {
        this._minimize = minimize
      }

      addLine(line: Line): void {
        this._lines.push({ start: -Infinity, end: Infinity, line })
      }

      addSegment(start: number, end: number, line: Line): void {
        this._lines.push({ start, end, line })
      }

      query(x: number): { value: number; lineId: number } {
        let resValue = this._minimize ? INF : -INF
        let resLineId = -1
        for (const [index, { start, end, line }] of this._lines.entries()) {
          if (x >= start && x < end) {
            const value = line.k * x + line.b
            if (this._minimize ? value < resValue : value > resValue) {
              resValue = value
              resLineId = index
            }
          }
        }
        return { value: resValue, lineId: resLineId }
      }
    }

    const q = 500
    const points = Array(q)
      .fill(0)
      .map(() => Math.floor(-Math.random() * 1e5) + 5e4)
    const tree1 = new LiChaoTree(points, { minimize: false })
    const tree2 = new Mocker(false)
    for (let i = 0; i < q; i++) {
      const k = -Math.floor(Math.random() * 1e5) + 5e4
      const b = -Math.floor(Math.random() * 1e5) + 5e4
      tree1.addLine({ k, b })
      tree2.addLine({ k, b })
      if (Math.random() < 0.5) {
        const k = -Math.floor(Math.random() * 1e5) + 5e4
        const b = -Math.floor(Math.random() * 1e5) + 5e4
        const start = Math.floor(Math.random() * 1e5)
        const end = Math.floor(Math.random() * 1e5)
        tree1.addSegment(start, end, { k, b })
        tree2.addSegment(start, end, { k, b })
      }

      const x = points[i]
      if (tree1.query(x).value !== tree2.query(x).value) {
        throw new Error()
      }
    }

    console.log('pass!')
  }
}
