/* eslint-disable no-param-reassign */
/* eslint-disable no-shadow */

// https://nyaannyaan.github.io/library/data-structure-2d/2d-segment-tree.hpp
// !二维线段树:单点修改/区间查询
// 如果值域很大,需要预先离散化
//
// addPoint: (row: number, col: number, value: E) => void
// build: () => void
// get: (row: number, col: number) => E
// set: (row: number, col: number, target: E) => void
// query: (row1: number, col1: number, row2: number, col2: number) => E

/**
 * 单点修改，区间查询的二维线段树.
 */
class SegmentTree2DPointUpdateRangeQuery<E> {
  private readonly _row: number
  private readonly _col: number
  private readonly _tree: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  constructor(row: number, col: number, e: () => E, op: (a: E, b: E) => E) {
    this._row = 1
    while (this._row < row) this._row <<= 1
    this._col = 1
    while (this._col < col) this._col <<= 1
    this._tree = Array((this._row * this._col) << 2)
    for (let i = 0; i < this._tree.length; i++) this._tree[i] = e()
    this._e = e
    this._op = op
  }

  /**
   * 在 {@link build} 之前调用，设置初始值.
   * 0 <= row < ROW, 0 <= col < COL.
   */
  addPoint(row: number, col: number, value: E): void {
    this._tree[this._id(row + this._row, col + this._col)] = value
  }

  /**
   * 如果调用了 {@link addPoint} 初始化，则需要调用此方法构建树.
   */
  build(): void {
    for (let c = this._col; c < this._col << 1; c++) {
      for (let r = this._row - 1; ~r; r--) {
        this._tree[this._id(r, c)] = this._op(
          this._tree[this._id(r << 1, c)],
          this._tree[this._id((r << 1) | 1, c)]
        )
      }
    }
    for (let r = 0; r < this._row << 1; r++) {
      for (let c = this._col - 1; ~c; c--) {
        this._tree[this._id(r, c)] = this._op(
          this._tree[this._id(r, c << 1)],
          this._tree[this._id(r, (c << 1) | 1)]
        )
      }
    }
  }

  /** 0 <= row < ROW, 0 <= col < COL. */
  get(row: number, col: number): E {
    return this._tree[this._id(row + this._row, col + this._col)]
  }

  /** 0 <= row < ROW, 0 <= col < COL. */
  set(row: number, col: number, target: E): void {
    let r = row + this._row
    let c = col + this._col
    this._tree[this._id(r, c)] = target
    for (let i = r >>> 1; i; i >>>= 1) {
      this._tree[this._id(i, c)] = this._op(
        this._tree[this._id(i << 1, c)],
        this._tree[this._id((i << 1) | 1, c)]
      )
    }
    for (; r; r >>>= 1) {
      for (let j = c >>> 1; j; j >>>= 1) {
        this._tree[this._id(r, j)] = this._op(
          this._tree[this._id(r, j << 1)],
          this._tree[this._id(r, (j << 1) | 1)]
        )
      }
    }
  }

  /** 0 <= row < ROW, 0 <= col < COL. */
  update(row: number, col: number, value: E): void {
    let r = row + this._row
    let c = col + this._col
    this._tree[this._id(r, c)] = this._op(this._tree[this._id(r, c)], value)
    for (let i = r >>> 1; i; i >>>= 1) {
      this._tree[this._id(i, c)] = this._op(
        this._tree[this._id(i << 1, c)],
        this._tree[this._id((i << 1) | 1, c)]
      )
    }
    for (; r; r >>>= 1) {
      for (let j = c >>> 1; j; j >>>= 1) {
        this._tree[this._id(r, j)] = this._op(
          this._tree[this._id(r, j << 1)],
          this._tree[this._id(r, (j << 1) | 1)]
        )
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
    row1 += this._row
    row2 += this._row
    col1 += this._col
    col2 += this._col
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
    return ((r * this._col) << 1) + c
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

export { SegmentTree2DPointUpdateRangeQuery }

if (require.main === module) {
  // https://leetcode.cn/problems/range-sum-query-2d-mutable/
  class NumMatrix {
    private readonly _ROW: number
    private readonly _COL: number
    private readonly _tree: SegmentTree2DPointUpdateRangeQuery<number>

    constructor(matrix: number[][]) {
      this._ROW = matrix.length
      this._COL = matrix[0].length
      this._tree = new SegmentTree2DPointUpdateRangeQuery(
        this._ROW,
        this._COL,
        () => 0,
        (a, b) => a + b
      )

      for (let r = 0; r < this._ROW; r++) {
        for (let c = 0; c < this._COL; c++) {
          this._tree.addPoint(r, c, matrix[r][c])
        }
      }

      this._tree.build() // !注意如果set了不要忘记 build
    }

    update(row: number, col: number, val: number): void {
      this._tree.set(row, col, val)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this._tree.query(row1, row2 + 1, col1, col2 + 1)
    }
  }

  const seg2d = new SegmentTree2DPointUpdateRangeQuery(
    3,
    4,
    () => 0,
    (a, b) => a + b
  )
  seg2d.update(0, 0, 2)
  seg2d.update(0, 0, 2)
  console.log(seg2d.query(0, 1, 0, 1))
}
