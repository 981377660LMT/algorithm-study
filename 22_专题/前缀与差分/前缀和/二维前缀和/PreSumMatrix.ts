class PreSumMatrix {
  private readonly preSum: number[][]

  constructor(matrix: number[][]) {
    const row = matrix.length
    const col = matrix[0].length
    const preSum = Array.from({ length: row + 1 }, () => Array(col + 1).fill(0))
    for (let i = 1; i < row + 1; i++) {
      for (let j = 1; j < col + 1; j++) {
        // 注意这里的减1
        preSum[i][j] =
          matrix[i - 1][j - 1] + preSum[i - 1][j] + preSum[i][j - 1] - preSum[i - 1][j - 1]
      }
    }

    this.preSum = preSum
  }

  /**
   * @returns 返回 左上角 (row1, col1) 、右下角 (row2, col2) 闭区间所描述的子矩阵的元素 总和 。
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    return (
      this.preSum[row2 + 1][col2 + 1] +
      this.preSum[row1][col1] -
      this.preSum[row2 + 1][col1] -
      this.preSum[row1][col2 + 1]
    )
  }
}

/**
 * @description 二维差分数组 适用于范围更新多 查询少的情景
 */
class DiffMatrix {
  private readonly ROW: number
  private readonly COL: number
  private readonly matrix: number[][]
  private diff: number[][]
  // private preSum: PreSumMatrix
  private dirty = false

  constructor(M: number[][]) {
    this.ROW = M.length
    this.COL = M[0].length
    this.matrix = Array.from({ length: this.ROW }, () => Array(this.COL).fill(0))
    for (let i = 0; i < this.ROW; i++) {
      for (let j = 0; j < this.COL; j++) {
        this.matrix[i][j] = M[i][j]
      }
    }

    // 需要额外大小为(m+2)∗(n+2)的差分数组，第一行第一列不用(始终为0)
    this.diff = Array.from({ length: this.ROW + 2 }, () => Array(this.COL + 2).fill(0))
    // this.preSum = new PreSumMatrix(this.matrix)
  }

  /**
   * @returns 区间更新 左上角 (row1, col1) 到 右下角 (row2, col2) 闭区间所描述的子矩阵的元素
   */
  add(row1: number, col1: number, row2: number, col2: number, delta: number): void {
    this.dirty = true
    this.diff[row1 + 1][col1 + 1] += delta
    this.diff[row1 + 1][col2 + 2] -= delta
    this.diff[row2 + 2][col1 + 1] -= delta
    this.diff[row2 + 2][col2 + 2] += delta
  }

  /**
   * @description 遍历矩阵，还原对应元素的增量
   */
  update(): void {
    this.dirty = false
    for (let r = 0; r < this.ROW; r++) {
      for (let c = 0; c < this.COL; c++) {
        this.diff[r + 1][c + 1] += this.diff[r + 1][c] + this.diff[r][c + 1] - this.diff[r][c]
        this.matrix[r][c] += this.diff[r + 1][c + 1]
      }
    }

    this.diff = Array.from({ length: this.ROW + 2 }, () => Array(this.COL + 2).fill(0))
    // this.preSum = new PreSumMatrix(this.matrix)
  }

  /**
   * @description 获取指定位置的元素值
   */
  query(row: number, col: number): number {
    if (this.dirty) throw new Error('矩阵已经更新，请先执行 update 方法')
    return this.matrix[row][col]
  }

  // /**
  //  * @returns 返回 左上角 (row1, col1) 、右下角 (row2, col2) 闭区间所描述的子矩阵的元素 总和 。
  //  */
  // queryRange(row1: number, col1: number, row2: number, col2: number): number {
  //   if (this.dirty) throw new Error('矩阵已经更新，请先执行 update 方法')
  //   return this.preSum.queryRange(row1, col1, row2, col2)
  // }
}

if (require.main === module) {
  const nm = new PreSumMatrix([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
  ])

  console.log(nm.queryRange(1, 1, 2, 2))
}

export { PreSumMatrix, DiffMatrix }
