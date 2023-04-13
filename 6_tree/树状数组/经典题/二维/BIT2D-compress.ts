/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// 二维树状数组 单点加 区间查询 (适用于值域很大的情况)
// Add : 单点加
// Query : 区间和
// QueryPrefix : 前缀和

interface Options<E> {
  xs: number[]
  ys: number[]
  ws?: E[]
  e: () => E
  op: (a: E, b: E) => E
  inv: (a: E) => E
}

class BIT2DCompress<E> {
  private static _bisectLeft(
    arr: ArrayLike<number>,
    x: number,
    left: number,
    right: number
  ): number {
    while (left <= right) {
      const mid = (left + right) >> 1
      if (arr[mid] < x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _inv: (a: E) => E
  private _n!: number
  private _keyX!: number[]
  private _keyY!: number[]
  private _indptr!: Uint32Array
  private _data!: E[]

  constructor(options: Options<E>) {
    let { xs, ys, e, op, inv, ws } = options
    this._e = e
    this._op = op
    this._inv = inv
    if (!ws) {
      ws = Array(xs.length)
      for (let i = 0; i < xs.length; i++) {
        ws[i] = e()
      }
    }
    this._build(xs, ys, ws)
  }

  private _build(xs: number[], ys: number[], ws: E[]): void {
    if (xs.length !== ys.length || xs.length !== ws.length) {
      throw new Error('Lengths of X, Y, and wt must be equal.')
    }

    this._keyX = [...new Set(xs)].sort((a, b) => a - b)
    this._n = this._keyX.length
    const n = this._n
    const keyYRaw: number[][] = Array(n)
    const datRaw: E[][] = Array(n)
    for (let i = 0; i < n; i++) {
      keyYRaw[i] = []
      datRaw[i] = []
    }
    const indices = Array(ys.length)
    for (let i = 0; i < ys.length; i++) {
      indices[i] = i
    }
    indices.sort((a, b) => ys[a] - ys[b])

    indices.forEach(i => {
      let ix = this._xtoi(xs[i])
      const y = ys[i]
      while (ix < n) {
        const ky = keyYRaw[ix]
        const tmp = datRaw[ix]
        if (ky.length === 0 || ky[ky.length - 1] < y) {
          keyYRaw[ix].push(y)
          tmp.push(ws[i])
        } else {
          tmp[tmp.length - 1] = this._op(tmp[tmp.length - 1], ws[i])
        }
        ix += (ix + 1) & -(ix + 1)
      }
    })

    this._indptr = new Uint32Array(n + 1)
    for (let i = 0; i < n; i++) {
      this._indptr[i + 1] = this._indptr[i] + keyYRaw[i].length
    }
    this._keyY = Array(this._indptr[n]).fill(0)
    this._data = Array(this._indptr[n])
    for (let i = 0; i < this._data.length; i++) {
      this._data[i] = this._e()
    }

    for (let i = 0; i < n; i++) {
      for (let j = 0; j < this._indptr[i + 1] - this._indptr[i]; j++) {
        this._keyY[this._indptr[i] + j] = keyYRaw[i][j]
        this._data[this._indptr[i] + j] = datRaw[i][j]
      }
    }

    for (let i = 0; i < n; i++) {
      const diff = this._indptr[i + 1] - this._indptr[i]
      for (let j = 0; j < diff - 1; j++) {
        const k = j + ((j + 1) & -(j + 1))
        if (k < diff) {
          this._data[this._indptr[i] + k] = this._op(
            this._data[this._indptr[i] + k],
            this._data[this._indptr[i] + j]
          )
        }
      }
    }
  }

  private _xtoi(x: number): number {
    return BIT2DCompress._bisectLeft(this._keyX, x, 0, this._n - 1)
  }

  /**
   * 点 (x,y) 的值加上 val.
   */
  add(x: number, y: number, val: E): void {
    for (let i = this._xtoi(x); i < this._n; i += (i + 1) & -(i + 1)) {
      this._add(i, y, val)
    }
  }

  /**
   * [lx,rx) * [ly,ry) 的和.
   */
  query(lx: number, rx: number, ly: number, ry: number): E {
    let pos = this._e()
    let neg = this._e()
    let l = this._xtoi(lx) - 1
    let r = this._xtoi(rx) - 1
    while (l < r) {
      pos = this._op(pos, this._prodI(r, ly, ry))
      r -= (r + 1) & -(r + 1)
    }
    while (r < l) {
      neg = this._op(neg, this._prodI(l, ly, ry))
      l -= (l + 1) & -(l + 1)
    }
    return this._op(pos, this._inv(neg))
  }

  /**
   * [0,rx) * [0,ry) 的和.
   */
  queryPrefix(rx: number, ry: number): E {
    let pos = this._e()
    let r = this._xtoi(rx) - 1
    while (r >= 0) {
      pos = this._op(pos, this._prefixProdI(r, ry))
      r -= (r + 1) & -(r + 1)
    }
    return pos
  }

  private _add(i: number, y: number, val: E): void {
    const lid = this._indptr[i]
    const n = this._indptr[i + 1] - this._indptr[i]
    let j = BIT2DCompress._bisectLeft(this._keyY, y, lid, lid + n - 1) - lid
    while (j < n) {
      this._data[lid + j] = this._op(this._data[lid + j], val)
      j += (j + 1) & -(j + 1)
    }
  }

  private _prodI(i: number, ly: number, ry: number): E {
    let pos = this._e()
    let neg = this._e()
    const lid = this._indptr[i]
    const n = this._indptr[i + 1] - this._indptr[i]
    let left = BIT2DCompress._bisectLeft(this._keyY, ly, lid, lid + n - 1) - lid - 1
    let right = BIT2DCompress._bisectLeft(this._keyY, ry, lid, lid + n - 1) - lid - 1
    while (left < right) {
      pos = this._op(pos, this._data[lid + right])
      right -= (right + 1) & -(right + 1)
    }
    while (right < left) {
      neg = this._op(neg, this._data[lid + left])
      left -= (left + 1) & -(left + 1)
    }
    return this._op(pos, this._inv(neg))
  }

  private _prefixProdI(i: number, ry: number): E {
    let pos = this._e()
    const lid = this._indptr[i]
    const n = this._indptr[i + 1] - this._indptr[i]
    let right = BIT2DCompress._bisectLeft(this._keyY, ry, lid, lid + n - 1) - lid - 1
    while (right >= 0) {
      pos = this._op(pos, this._data[lid + right])
      right -= (right + 1) & -(right + 1)
    }
    return pos
  }
}

if (require.main === module) {
  const tree = new BIT2DCompress({
    xs: [0, 0, 0, 1, 1, 1, 2, 2, 2],
    ys: [0, 1, 2, 0, 1, 2, 0, 1, 2],
    ws: [1, 2, 3, 4, 5, 6, 7, 8, 9],
    op: (a, b) => a + b,
    e: () => 0,
    inv: a => -a
  })
  console.log(tree.query(0, 3, 0, 3))

  // https://leetcode.cn/problems/range-sum-query-2d-mutable/submissions/
  class NumMatrix {
    private readonly matrix: number[][]
    private readonly ROW: number
    private readonly COL: number
    private readonly tree: BIT2DCompress<number>

    constructor(matrix: number[][]) {
      this.matrix = matrix
      this.ROW = matrix.length
      this.COL = matrix[0].length
      const xs = Array(this.ROW * this.COL)
      const ys = Array(this.ROW * this.COL)
      const ws = Array(this.ROW * this.COL)
      for (let r = 0; r < this.ROW; r++) {
        for (let c = 0; c < this.COL; c++) {
          xs[r * this.COL + c] = r
          ys[r * this.COL + c] = c
          ws[r * this.COL + c] = matrix[r][c]
        }
      }
      this.tree = new BIT2DCompress({
        xs,
        ys,
        ws,
        e: () => 0,
        op: (a, b) => a + b,
        inv: a => -a
      })
    }

    update(row: number, col: number, val: number): void {
      const pre = this.matrix[row][col]
      const diff = val - pre
      this.matrix[row][col] = val
      this.tree.add(row, col, diff)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this.tree.query(row1, row2 + 1, col1, col2 + 1)
    }
  }
}

export { BIT2DCompress }
