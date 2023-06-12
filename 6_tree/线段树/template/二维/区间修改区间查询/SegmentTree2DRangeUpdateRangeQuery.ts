/* eslint-disable no-inner-declarations */
/* eslint-disable prefer-destructuring */

// 适用于 ROW*COL<=1e5 稠密的二维矩阵
// !如果ROW/COL很大,需要离散化

const INF = 2e15

interface IRangeUpdateRangeQuery<E, Id> {
  update(start: number, end: number, lazy: Id): void
  query(start: number, end: number): E
  get(index: number): E
  set(index: number, value: E): void
}

/**
 * 二维区间更新，区间查询的线段树(树套树).
 */
class SegmentTree2DRangeUpdateRangeQuery<E = number, Id = number> {
  /**
   * 存储内层的"树"结构.
   */
  private readonly _seg: IRangeUpdateRangeQuery<E, Id>[]

  /**
   * 初始化内层"树"的函数.
   */
  private readonly _init1D: () => IRangeUpdateRangeQuery<E, Id>

  /**
   * 当列数超过行数时,需要对矩阵进行旋转,将列数控制在根号以下.
   */
  private readonly _needRotate: boolean

  /**
   * 原始矩阵的行数(未经旋转).
   */
  private readonly _rawRow: number

  private readonly _size: number

  private readonly _e: () => E
  private readonly _mergeRow: (a: E, b: E) => E

  /**
   * @param row 行数.对时间复杂度贡献为`O(log(row))`.
   * @param col 列数.内部树的大小.列数越小,对内部树的时间复杂度要求越低.
   * @param createRangeUpdatePointGet1D 初始化内层"树"的函数.入参为内层"树"的大小.
   * @param e 初始化元素的函数(内层树的幺元).
   * @param mergeRow 合并两个元素的函数(合并两个内层树的答案).
   */
  constructor(
    row: number,
    col: number,
    createRangeUpdatePointGet1D: (n: number) => IRangeUpdateRangeQuery<E, Id>,
    e: () => E,
    mergeRow: (a: E, b: E) => E
  ) {
    this._rawRow = row
    this._needRotate = row < col
    if (this._needRotate) {
      row ^= col
      col ^= row
      row ^= col
    }

    let size = 1
    while (size < row) size <<= 1
    this._seg = Array(2 * size - 1)
    this._e = e
    this._mergeRow = mergeRow
    this._init1D = () => createRangeUpdatePointGet1D(col)
    this._size = size
  }

  /**
   * 将`[row1,row2)`x`[col1,col2)`的区间值与`lazy`作用.
   */
  update(row1: number, row2: number, col1: number, col2: number, lazy: Id): void {
    if (this._needRotate) {
      const tmp1 = row1
      const tmp2 = row2
      row1 = col1
      row2 = col2
      col1 = this._rawRow - tmp2
      col2 = this._rawRow - tmp1
    }

    this._update(row1, row2, col1, col2, lazy, 0, 0, this._size)
  }

  /**
   * 查询区间`[row1,row2)`x`[col1,col2)`的聚合值.
   */
  query(row1: number, row2: number, col1: number, col2: number): E {
    if (this._needRotate) {
      const tmp1 = row1
      const tmp2 = row2
      row1 = col1
      row2 = col2
      col1 = this._rawRow - tmp2
      col2 = this._rawRow - tmp1
    }

    const res = [this._e()]
    this._query(row1, row2, col1, col2, 0, 0, this._size, res)
    return res[0]
  }

  get(row: number, col: number): E {
    if (this._needRotate) {
      const tmp = row
      row = col
      col = this._rawRow - tmp - 1
    }

    row += this._size - 1
    let res = this._seg[row] ? this._seg[row].get(col) : this._e()
    while (row > 0) {
      row = (row - 1) >> 1
      if (this._seg[row]) res = this._mergeRow(res, this._seg[row].get(col))
    }
    return res
  }

  set(row: number, col: number, value: E): void {
    if (this._needRotate) {
      const tmp = row
      row = col
      col = this._rawRow - tmp - 1
    }

    row += this._size - 1
    if (!this._seg[row]) this._seg[row] = this._init1D()
    this._seg[row].set(col, value)
    while (row > 0) {
      row = (row - 1) >> 1
      if (!this._seg[row]) this._seg[row] = this._init1D()
      this._seg[row].set(col, value)
    }
  }

  private _update(
    R: number,
    C: number,
    start: number,
    end: number,
    lazy: Id,
    pos: number,
    r: number,
    c: number
  ): void {
    if (c <= R || C <= r) return
    if (R <= r && c <= C) {
      if (!this._seg[pos]) this._seg[pos] = this._init1D()
      this._seg[pos].update(start, end, lazy)
    } else {
      const mid = (r + c) >>> 1
      this._update(R, C, start, end, lazy, 2 * pos + 1, r, mid)
      this._update(R, C, start, end, lazy, 2 * pos + 2, mid, c)
    }
  }

  private _query(
    R: number,
    C: number,
    start: number,
    end: number,
    pos: number,
    r: number,
    c: number,
    ref: E[]
  ): void {
    if (c <= R || C <= r) return
    if (R <= r && c <= C) {
      if (this._seg[pos]) {
        ref[0] = this._mergeRow(ref[0], this._seg[pos].query(start, end))
        console.log(777)
      }
      return
    }
    console.log(R, C, start, end, pos, r, c, ref)
    const mid = (r + c) >>> 1
    this._query(R, C, start, end, 2 * pos + 1, r, mid, ref)
    this._query(R, C, start, end, 2 * pos + 2, mid, c, ref)
  }
}

export { SegmentTree2DRangeUpdateRangeQuery }

if (require.main === module) {
  function possibleToStamp(grid: number[][], h: number, w: number): boolean {
    const ROW = grid.length
    const COL = grid[0].length
    const seg = new SegmentTree2DRangeUpdateRangeQuery(
      ROW,
      COL,
      n => new NaiveTree(n),
      () => 0,
      (a, b) => a + b
    )

    for (let r = 0; r + h <= ROW; r++) {
      for (let c = 0; c + w <= COL; c++) {
        if (!seg.query(r, r + h, c, c + w)) {
          seg.update(r, r + h, c, c + w, 1)
        }
      }
    }

    for (let r = 0; r < ROW; r++) {
      for (let c = 0; c < COL; c++) {
        if (!grid[r][c] && !seg.get(r, c)) return false
      }
    }

    return true
  }

  class NaiveTree implements IRangeUpdateRangeQuery<number, number> {
    private readonly _nums: number[]

    constructor(n: number) {
      this._nums = Array(n).fill(0)
    }

    query(start: number, end: number): number {
      let res = 0
      for (let i = start; i < end; i++) res += this._nums[i]
      return res
    }

    update(start: number, end: number, lazy: number): void {
      for (let i = start; i < end; i++) this._nums[i] += lazy
    }

    set(index: number, value: number): void {
      this._nums[index] = value
    }

    get(index: number): number {
      return this._nums[index]
    }
  }

  const seg2d = new SegmentTree2DRangeUpdateRangeQuery(
    3,
    4,
    n => new NaiveTree(n),
    () => 0,
    (a, b) => a + b
  )

  seg2d.update(0, 2, 0, 4, 10)
  printGrid(seg2d, 3, 4)
  seg2d.update(1, 3, 1, 3, 20)
  printGrid(seg2d, 3, 4)
  seg2d.set(1, 1, 0)
  printGrid(seg2d, 3, 4)
  console.log(seg2d.query(0, 3, 0, 4))

  function printGrid(seg: SegmentTree2DRangeUpdateRangeQuery, row: number, col: number): void {
    const res = Array(row)
    for (let r = 0; r < row; r++) {
      res[r] = Array(col)
      for (let c = 0; c < col; c++) {
        res[r][c] = seg.get(r, c)
      }
    }
    console.table(res)
  }
}
