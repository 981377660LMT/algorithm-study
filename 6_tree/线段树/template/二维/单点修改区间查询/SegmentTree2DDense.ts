/* eslint-disable no-param-reassign */
/* eslint-disable no-shadow */

// https://nyaannyaan.github.io/library/data-structure-2d/2d-segment-tree.hpp
// !二维线段树:单点修改/区间查询
// 如果值域很大,需要预先离散化
// build: (f: (r: number, c: number) => E) => void
// get: (row: number, col: number) => E
// set: (row: number, col: number, target: E) => void
// query: (row1: number, col1: number, row2: number, col2: number) => E

/**
 * 单点修改，区间查询的二维线段树.
 */
class SegmentTree2DDense<E> {
  private readonly _row: number
  private readonly _col: number
  private readonly _rowOffset: number
  private readonly _colOffset: number
  private readonly _tree: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  constructor(row: number, col: number, e: () => E, op: (a: E, b: E) => E) {
    this._row = row
    this._col = col
    this._rowOffset = 1
    while (this._rowOffset < row) this._rowOffset <<= 1
    this._colOffset = 1
    while (this._colOffset < col) this._colOffset <<= 1
    this._tree = Array((this._rowOffset * this._colOffset) << 2)
    for (let i = 0; i < this._tree.length; i++) this._tree[i] = e()
    this._e = e
    this._op = op
  }

  build(f: (r: number, c: number) => E): void {
    for (let r = 0; r < this._row; r++) {
      for (let c = 0; c < this._col; c++) {
        this._tree[this._id(r + this._rowOffset, c + this._colOffset)] = f(r, c)
      }
    }
    for (let c = this._colOffset; c < this._colOffset << 1; c++) {
      for (let r = this._rowOffset - 1; ~r; r--) {
        this._tree[this._id(r, c)] = this._op(this._tree[this._id(r << 1, c)], this._tree[this._id((r << 1) | 1, c)])
      }
    }
    for (let r = 0; r < this._rowOffset << 1; r++) {
      for (let c = this._colOffset - 1; ~c; c--) {
        this._tree[this._id(r, c)] = this._op(this._tree[this._id(r, c << 1)], this._tree[this._id(r, (c << 1) | 1)])
      }
    }
  }

  /** 0 <= row < ROW, 0 <= col < COL. */
  get(row: number, col: number): E {
    return this._tree[this._id(row + this._rowOffset, col + this._colOffset)]
  }

  /** 0 <= row < ROW, 0 <= col < COL. */
  set(row: number, col: number, target: E): void {
    let r = row + this._rowOffset
    let c = col + this._colOffset
    this._tree[this._id(r, c)] = target
    for (let i = r >>> 1; i; i >>>= 1) {
      this._tree[this._id(i, c)] = this._op(this._tree[this._id(i << 1, c)], this._tree[this._id((i << 1) | 1, c)])
    }
    for (; r; r >>>= 1) {
      for (let j = c >>> 1; j; j >>>= 1) {
        this._tree[this._id(r, j)] = this._op(this._tree[this._id(r, j << 1)], this._tree[this._id(r, (j << 1) | 1)])
      }
    }
  }

  /** 0 <= row < ROW, 0 <= col < COL. */
  update(row: number, col: number, value: E): void {
    let r = row + this._rowOffset
    let c = col + this._colOffset
    this._tree[this._id(r, c)] = this._op(this._tree[this._id(r, c)], value)
    for (let i = r >>> 1; i; i >>>= 1) {
      this._tree[this._id(i, c)] = this._op(this._tree[this._id(i << 1, c)], this._tree[this._id((i << 1) | 1, c)])
    }
    for (; r; r >>>= 1) {
      for (let j = c >>> 1; j; j >>>= 1) {
        this._tree[this._id(r, j)] = this._op(this._tree[this._id(r, j << 1)], this._tree[this._id(r, (j << 1) | 1)])
      }
    }
  }

  /**
   * 查询区间 `[row1, row2)` x `[col1, col2)` 的聚合值.
   * 0 <= row1 <= row2 <= ROW.
   * 0 <= col1 <= col2 <= COL.
   */
  query(row1: number, row2: number, col1: number, col2: number): E {
    if (row1 >= row2 || col1 >= col2) return this._e()
    let res = this._e()
    row1 += this._rowOffset
    row2 += this._rowOffset
    col1 += this._colOffset
    col2 += this._colOffset
    for (; row1 < row2; row1 >>>= 1, row2 >>>= 1) {
      if (row1 & 1) {
        res = this._op(res, this._query(row1, col1, col2))
        row1++
      }
      if (row2 & 1) {
        row2--
        res = this._op(res, this._query(row2, col1, col2))
      }
    }
    return res
  }

  private _id(r: number, c: number): number {
    return ((r * this._colOffset) << 1) + c
  }

  private _query(r: number, c1: number, c2: number): E {
    let res = this._e()
    for (; c1 < c2; c1 >>>= 1, c2 >>>= 1) {
      if (c1 & 1) {
        res = this._op(res, this._tree[this._id(r, c1)])
        c1++
      }
      if (c2 & 1) {
        c2--
        res = this._op(res, this._tree[this._id(r, c2)])
      }
    }
    return res
  }
}

export { SegmentTree2DDense }

if (require.main === module) {
  // https://leetcode.cn/problems/range-sum-query-2d-mutable/
  class NumMatrix {
    private readonly _ROW: number
    private readonly _COL: number
    private readonly _tree: SegmentTree2DDense<number>

    constructor(matrix: number[][]) {
      this._ROW = matrix.length
      this._COL = matrix[0].length
      this._tree = new SegmentTree2DDense(
        this._ROW,
        this._COL,
        () => 0,
        (a, b) => a + b
      )

      this._tree.build((r, c) => matrix[r][c])
    }

    update(row: number, col: number, val: number): void {
      this._tree.set(row, col, val)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this._tree.query(row1, row2 + 1, col1, col2 + 1)
    }
  }

  const seg2d = new SegmentTree2DDense(
    3,
    4,
    () => 0,
    (a, b) => a + b
  )
  seg2d.update(0, 0, 2)
  seg2d.update(0, 0, 2)
  console.log(seg2d.query(0, 1, 0, 1))
}
