/**
 * 二维差分数组,适用于多次范围更新，单次查询的场景.
 */
class DiffMatrix {
  /**
   * 内部矩阵.
   */
  readonly matrix: number[]
  private readonly _row: number
  private readonly _col: number
  private readonly _diff: number[]

  constructor(mat: number[][]) {
    const row = mat.length
    const col = mat[0].length
    const newMat = Array(row * col)
    for (let i = 0; i < row; i++) {
      const matRow = mat[i]
      const offset = i * col
      for (let j = 0; j < col; j++) {
        newMat[offset + j] = matRow[j]
      }
    }

    const diff = Array((row + 2) * (col + 2)).fill(0)

    this.matrix = newMat
    this._row = row
    this._col = col
    this._diff = diff
  }

  /**
   * 区间更新 左上角` (row1, col1)` 到 右下角` (row2, col2) `闭区间所描述的子矩阵的元素
   */
  add(row1: number, col1: number, row2: number, col2: number, delta: number): void {
    const col = this._col + 2
    this._diff[(row1 + 1) * col + col1 + 1] += delta
    this._diff[(row1 + 1) * col + col2 + 2] -= delta
    this._diff[(row2 + 2) * col + col1 + 1] -= delta
    this._diff[(row2 + 2) * col + col2 + 2] += delta
  }

  /**
   * 遍历矩阵，还原对应元素的增量, 并更新矩阵{@link matrix}.
   */
  update(): void {
    const diff = this._diff
    const mat = this.matrix
    for (let i = 0; i < this._row; i++) {
      const offsetD0 = i * (this._col + 2)
      const offsetD1 = (i + 1) * (this._col + 2)
      const offsetM = i * this._col
      for (let j = 0; j < this._col; j++) {
        diff[offsetD1 + j + 1] += diff[offsetD1 + j] + diff[offsetD0 + j + 1] - diff[offsetD0 + j]
        mat[offsetM + j] += diff[offsetD1 + j + 1]
      }
    }
  }

  /**
   * 查询矩阵中指定位置的元素.
   * !查询前需要先调用{@link update}方法.
   */
  query(row: number, col: number): number {
    return this.matrix[row * this._col + col]
  }
}

/**
 * 二维前缀和数组.查询静态矩阵中指定矩形区间的元素总和.
 */
class PreSumMatrix {
  private readonly _preSum: number[]
  private readonly _col: number

  constructor(mat: number[][]) {
    const row = mat.length
    const col = mat[0].length
    const preSum: number[] = Array((row + 1) * (col + 1)).fill(0)

    for (let i = 0; i < row; i++) {
      const offsetP0 = i * (col + 1)
      const offsetP1 = (i + 1) * (col + 1)
      const matRow = mat[i]
      for (let j = 0; j < col; j++) {
        preSum[offsetP1 + j + 1] =
          matRow[j] + preSum[offsetP0 + j + 1] + preSum[offsetP1 + j] - preSum[offsetP0 + j]
      }
    }
    this._preSum = preSum
    this._col = col
  }

  /**
   * 返回 左上角 `(row1, col1)` 、右下角 `(row2, col2)` 闭区间所描述的子矩阵的元素 总和 。
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    const col = this._col + 1
    return (
      this._preSum[(row2 + 1) * col + col2 + 1] -
      this._preSum[(row2 + 1) * col + col1] -
      this._preSum[row1 * col + col2 + 1] +
      this._preSum[row1 * col + col1]
    )
  }
}

if (require.main === module) {
  const nm = new PreSumMatrix([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
  ])

  console.log(nm.queryRange(1, 1, 2, 2))
}

export { PreSumMatrix, DiffMatrix }
