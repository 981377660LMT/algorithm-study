// # 1 <= m, n <= 200
// # 最多调用1e4 次 sumRegion 和 update 方法

import { BIT3, BIT4 } from './BIT'

class NumMatrix {
  private readonly matrix: number[][]
  private readonly ROW: number
  private readonly COL: number
  private readonly tree: BIT3

  constructor(matrix: number[][]) {
    this.matrix = matrix
    this.ROW = matrix.length
    this.COL = matrix[0].length
    this.tree = new BIT3(this.ROW, this.COL)
    for (let r = 0; r < this.ROW; r++) {
      for (let c = 0; c < this.COL; c++) {
        this.tree.update(r, c, matrix[r][c])
      }
    }
  }

  update(row: number, col: number, val: number): void {
    const delta = val - this.matrix[row][col]
    this.matrix[row][col] = val
    this.tree.update(row, col, delta)
  }

  sumRegion(row1: number, col1: number, row2: number, col2: number): number {
    return this.tree.queryRange(row1, col1, row2, col2)
  }
}

class NumMatrix2 {
  private readonly matrix: number[][]
  private readonly ROW: number
  private readonly COL: number
  private readonly tree: BIT4

  constructor(matrix: number[][]) {
    this.matrix = matrix
    this.ROW = matrix.length
    this.COL = matrix[0].length
    this.tree = new BIT4(this.ROW, this.COL)
    for (let r = 0; r < this.ROW; r++) {
      for (let c = 0; c < this.COL; c++) {
        this.tree.updateRange(r, c, r, c, matrix[r][c])
      }
    }
  }

  update(row: number, col: number, val: number): void {
    const delta = val - this.matrix[row][col]
    this.matrix[row][col] = val
    this.tree.updateRange(row, col, row, col, delta)
  }

  sumRegion(row1: number, col1: number, row2: number, col2: number): number {
    return this.tree.queryRange(row1, col1, row2, col2)
  }
}
/**
 * Your NumMatrix object will be instantiated and called as such:
 * var obj = new NumMatrix(matrix)
 * obj.update(row,col,val)
 * var param_2 = obj.sumRegion(row1,col1,row2,col2)
 */
