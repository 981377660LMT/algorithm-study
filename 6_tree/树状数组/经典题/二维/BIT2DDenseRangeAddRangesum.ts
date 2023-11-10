/* eslint-disable max-len */

class BIT2DDenseRangeAddRangesum {
  private readonly _ROW: number
  private readonly _COL: number
  private readonly _tree1: Float64Array
  private readonly _tree2: Float64Array
  private readonly _tree3: Float64Array
  private readonly _tree4: Float64Array

  constructor(row: number, col: number) {
    this._ROW = row
    this._COL = col
    this._tree1 = new Float64Array((row + 1) * (col + 1))
    this._tree2 = new Float64Array((row + 1) * (col + 1))
    this._tree3 = new Float64Array((row + 1) * (col + 1))
    this._tree4 = new Float64Array((row + 1) * (col + 1))
  }

  /**
   * [row1,row2) x [col1,col2) 的值加上 delta.
   * 0<=row1<=row2<=ROW, 0<=col1<=col2<=COL.
   */
  addRange(row1: number, row2: number, col1: number, col2: number, delta: number): void {
    if (row1 >= row2 || col1 >= col2) return
    this._add(row1, col1, delta)
    this._add(row2, col1, -delta)
    this._add(row1, col2, -delta)
    this._add(row2, col2, delta)
  }

  /**
   * [0,row) x [0,col) 的和.
   * 0<=row<=ROW, 0<=col<=COL.
   */
  queryPrefix(row: number, col: number): number {
    if (row > this._ROW) row = this._ROW
    if (col > this._COL) col = this._COL
    let res = 0
    for (let r = row; r > 0; r -= r & -r) {
      for (let c = col; c > 0; c -= c & -c) {
        const id = this._id(r, c)
        res += row * col * this._tree1[id] - col * this._tree2[id] - row * this._tree3[id] + this._tree4[id]
      }
    }
    return res
  }

  /**
   * [row1,row2) x [col1,col2) 的和.
   * 0<=row1<=row2<=ROW, 0<=col1<=col2<=COL.
   */
  queryRange(row1: number, row2: number, col1: number, col2: number): number {
    if (row1 < 0) row1 = 0
    if (col1 < 0) col1 = 0
    if (row2 > this._ROW) row2 = this._ROW
    if (col2 > this._COL) col2 = this._COL
    if (row1 >= row2 || col1 >= col2) return 0
    return this.queryPrefix(row2, col2) - this.queryPrefix(row1, col2) - this.queryPrefix(row2, col1) + this.queryPrefix(row1, col1)
  }

  toString(): string {
    const res: string[] = []
    for (let r = 0; r < this._ROW; r++) {
      const row: number[] = []
      for (let c = 0; c < this._COL; c++) {
        row.push(this.queryRange(r, r + 1, c, c + 1))
      }
      res.push(row.join(','))
    }
    return res.join('\n')
  }

  private _add(row: number, col: number, delta: number): void {
    row++
    col++
    for (let r = row; r <= this._ROW; r += r & -r) {
      for (let c = col; c <= this._COL; c += c & -c) {
        const id = this._id(r, c)
        this._tree1[id] += delta
        this._tree2[id] += (row - 1) * delta
        this._tree3[id] += (col - 1) * delta
        this._tree4[id] += (row - 1) * (col - 1) * delta
      }
    }
  }

  private _id(row: number, col: number): number {
    return row * (this._COL + 1) + col
  }
}

export { BIT2DDenseRangeAddRangesum }

if (require.main === module) {
  const bit = new BIT2DDenseRangeAddRangesum(3, 3)
  console.log(bit.toString())
  bit.addRange(1, 2, 1, 2, 1)
  console.log(bit.toString())
  console.log(bit.queryPrefix(2, 2), bit.queryPrefix(2, 3), bit.queryPrefix(3, 2), bit.queryPrefix(3, 3))

  // https://leetcode.cn/problems/range-sum-query-2d-mutable/
  class NumMatrix {
    private readonly _bit: BIT2DDenseRangeAddRangesum
    private readonly _matrix: number[][]

    constructor(matrix: number[][]) {
      this._matrix = matrix
      this._bit = new BIT2DDenseRangeAddRangesum(matrix.length, matrix[0].length)
      for (let r = 0; r < matrix.length; r++) {
        for (let c = 0; c < matrix[0].length; c++) {
          this._bit.addRange(r, r + 1, c, c + 1, matrix[r][c])
        }
      }
    }

    update(row: number, col: number, val: number): void {
      const pre = this._matrix[row][col]
      const diff = val - pre
      this._matrix[row][col] = val
      this._bit.addRange(row, row + 1, col, col + 1, diff)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this._bit.queryRange(row1, row2 + 1, col1, col2 + 1)
    }
  }
}
