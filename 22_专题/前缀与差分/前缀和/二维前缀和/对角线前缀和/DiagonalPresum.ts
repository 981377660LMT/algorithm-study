/* eslint-disable prefer-destructuring */
/* eslint-disable max-len */

/**
 * 二维矩阵斜对角线前缀和.
 */
class DiagonalPresum {
  private readonly _preSum1: Float64Array
  private readonly _preSum2: Float64Array
  private readonly _hash: number

  constructor(grid: ArrayLike<ArrayLike<number>>) {
    const row = grid.length
    const col = grid[0].length
    const preSum1 = new Float64Array((row + 1) * (col + 1))
    const preSum2 = new Float64Array((row + 1) * (col + 1))
    for (let r = 0; r < row; r++) {
      const tmp = grid[r]
      for (let c = 0; c < col; c++) {
        preSum1[(r + 1) * (col + 1) + c + 1] = preSum1[r * (col + 1) + c] + tmp[c]
        preSum2[(r + 1) * (col + 1) + c] = preSum2[r * (col + 1) + c + 1] + tmp[c]
      }
    }
    this._preSum1 = preSum1
    this._preSum2 = preSum2
    this._hash = col + 1
  }

  /**
   * 正对角线左上角到右下角的前缀和↘.
   */
  queryDiagnoal(leftUp: [number, number], rightDown: [number, number]): number
  queryDiagnoal(leftUpRow: number, leftUpCol: number, rightDownRow: number, rightDownCol: number): number
  queryDiagnoal(...args: any[]): number {
    if (args.length === 2) {
      const { 0: leftUp, 1: rightDown } = args
      const { 0: r1, 1: c1 } = leftUp
      const { 0: r2, 1: c2 } = rightDown
      return this._preSum1[(r2 + 1) * this._hash + c2 + 1] - this._preSum1[r1 * this._hash + c1]
    }
    const { 0: r1, 1: c1, 2: r2, 3: c2 } = args
    return this._preSum1[(r2 + 1) * this._hash + c2 + 1] - this._preSum1[r1 * this._hash + c1]
  }

  /**
   * 副对角线左下角到右上角的前缀和↙.
   */
  queryAntiDiagonal(leftDown: [number, number], rightUp: [number, number]): number
  queryAntiDiagonal(leftDownRow: number, leftDownCol: number, rightUpRow: number, rightUpCol: number): number
  queryAntiDiagonal(...args: any[]): number {
    if (args.length === 2) {
      const { 0: leftDown, 1: rightUp } = args
      const { 0: r1, 1: c1 } = leftDown
      const { 0: r2, 1: c2 } = rightUp
      return this._preSum2[(r1 + 1) * this._hash + c1] - this._preSum2[r2 * this._hash + c2 + 1]
    }
    const { 0: r1, 1: c1, 2: r2, 3: c2 } = args
    return this._preSum2[(r1 + 1) * this._hash + c1] - this._preSum2[r2 * this._hash + c2 + 1]
  }
}

export { DiagonalPresum }

if (require.main === module) {
  const grid = [
    [3, 7, 1, 4, 5, 9, 2, 8, 6, 0],
    [8, 0, 2, 9, 4, 1, 7, 3, 5, 6],
    [6, 9, 5, 3, 2, 7, 1, 4, 8, 0],
    [7, 1, 8, 5, 9, 2, 6, 0, 3, 4],
    [5, 6, 9, 0, 8, 4, 3, 7, 2, 1],
    [2, 3, 4, 7, 1, 6, 8, 5, 9, 0],
    [9, 2, 6, 8, 7, 3, 5, 1, 0, 4],
    [0, 8, 7, 6, 3, 5, 4, 9, 1, 2],
    [4, 5, 3, 1, 0, 8, 9, 6, 7, 2],
    [1, 4, 0, 2, 6, 5, 3, 9, 7, 8]
  ]
  const S = new DiagonalPresum(grid)
  console.log(S.queryDiagnoal(0, 0, 9, 9))
  console.log(S.queryDiagnoal([0, 0], [9, 9]))
  console.log(S.queryAntiDiagonal(9, 0, 0, 9))
  console.log(S.queryAntiDiagonal([2, 0], [0, 2]))
}
