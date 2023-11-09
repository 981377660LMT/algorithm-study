/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// 二维树状数组 单点加 区间查询
// Add : 单点加
// Query : 区间和
// QueryPrefix : 前缀和

interface IAbelianGroup<E> {
  e: () => E
  op: (a: E, b: E) => E
  inv?: (a: E) => E
}

class BIT2DDense<E> {
  private readonly _row: number
  private readonly _col: number
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _inv?: (a: E) => E
  private readonly _data: E[]

  constructor(row: number, col: number, options: IAbelianGroup<E> & ThisType<void>) {
    const { e, op, inv } = options
    this._row = row
    this._col = col
    this._e = e
    this._op = op
    this._inv = inv
    this._data = Array(row * col)
    for (let i = 0; i < this._data.length; i++) this._data[i] = e()
  }

  build(f: (r: number, c: number) => E): void {
    const ROW = this._row
    const COL = this._col
    for (let r = 0; r < ROW; r++) {
      for (let c = 0; c < COL; c++) {
        this._data[r * COL + c] = f(r, c)
      }
    }
    for (let r = 1; r <= ROW; r++) {
      for (let c = 1; c <= COL; c++) {
        const nc = c + (c & -c)
        if (nc <= COL) this._data[this._id(r, nc)] = this._op(this._data[this._id(r, nc)], this._data[this._id(r, c)])
      }
    }
    for (let r = 1; r <= ROW; r++) {
      for (let c = 1; c <= COL; c++) {
        const nr = r + (r & -r)
        if (nr <= ROW) this._data[this._id(nr, c)] = this._op(this._data[this._id(nr, c)], this._data[this._id(r, c)])
      }
    }
  }

  /**
   * 点 (r,c) 的值加上 value.
   */
  update(r: number, c: number, value: E): void {
    r++
    for (; r <= this._row; r += r & -r) {
      this._add(r, c, value)
    }
  }

  /**
   * [0,r) * [0,c) 的和.
   */
  queryPrefix(r: number, c: number): E {
    if (r > this._row) r = this._row
    if (c > this._col) c = this._col
    let res = this._e()
    while (r > 0) {
      res = this._op(res, this._queryPrefix(r, c))
      r -= r & -r
    }
    return res
  }

  /**
   * [r1,r2) * [c1,c2) 的和.
   */
  queryRange(r1: number, r2: number, c1: number, c2: number): E {
    if (r2 > this._row) r2 = this._row
    if (c2 > this._col) c2 = this._col
    if (r1 >= r2 || c1 >= c2) return this._e()
    if (this._inv == undefined) throw new Error('inv must be defined when query is called.')
    let pos = this._e()
    let neg = this._e()
    while (r1 < r2) {
      pos = this._op(pos, this._queryRange(r2, c1, c2))
      r2 -= r2 & -r2
    }
    while (r2 < r1) {
      neg = this._op(neg, this._queryRange(r1, c1, c2))
      r1 -= r1 & -r1
    }
    return this._op(pos, this._inv(neg))
  }

  toString(): string {
    const grid = Array(this._row)
    for (let r = 0; r < this._row; r++) {
      const row = Array(this._col)
      for (let c = 0; c < this._col; c++) {
        row[c] = this.queryRange(r, r + 1, c, c + 1)
      }
      grid[r] = row.join(',')
    }
    return grid.join('\n')
  }

  private _id(r: number, c: number): number {
    return this._col * (r - 1) + (c - 1)
  }

  private _add(r: number, c: number, val: E): void {
    c++
    for (; c <= this._col; c += c & -c) {
      this._data[this._id(r, c)] = this._op(this._data[this._id(r, c)], val)
    }
  }

  private _queryPrefix(r: number, c: number): E {
    let res = this._e()
    while (c > 0) {
      res = this._op(res, this._data[this._id(r, c)])
      c -= c & -c
    }
    return res
  }

  private _queryRange(r: number, c1: number, c2: number): E {
    let pos = this._e()
    let neg = this._e()
    while (c1 < c2) {
      pos = this._op(pos, this._data[this._id(r, c2)])
      c2 -= c2 & -c2
    }
    while (c2 < c1) {
      neg = this._op(neg, this._data[this._id(r, c1)])
      c1 -= c1 & -c1
    }
    return this._op(pos, this._inv!(neg))
  }
}

export { BIT2DDense }

if (require.main === module) {
  const bit = new BIT2DDense(3, 3, { e: () => 0, op: (a, b) => a + b, inv: a => -a })
  bit.build((r, c) => r * 3 + c)
  console.log(bit.toString())

  class NumMatrix {
    private readonly _bit: BIT2DDense<number>
    private readonly _matrix: number[][]
    constructor(matrix: number[][]) {
      this._matrix = matrix
      this._bit = new BIT2DDense(matrix.length, matrix[0].length, { e: () => 0, op: (a, b) => a + b, inv: a => -a })
      this._bit.build((r, c) => matrix[r][c])
    }

    update(row: number, col: number, val: number): void {
      const pre = this._matrix[row][col]
      const diff = val - pre
      this._matrix[row][col] = val
      this._bit.update(row, col, diff)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this._bit.queryRange(row1, row2 + 1, col1, col2 + 1)
    }
  }
}
