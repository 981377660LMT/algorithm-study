// 计算其子矩形范围内元素的总和，该子矩阵的左上角为 (row1, col1) ，右下角为 (row2, col2) 。
class PreSumMatrix {
  private preSum: number[][]

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
  sumRegion(row1: number, col1: number, row2: number, col2: number): number {
    return (
      this.preSum[row2 + 1][col2 + 1] +
      this.preSum[row1][col1] -
      this.preSum[row2 + 1][col1] -
      this.preSum[row1][col2 + 1]
    )
  }
}

if (require.main === module) {
  const nm = new PreSumMatrix([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
  ])

  console.log(nm.sumRegion(1, 1, 2, 2))
}

export { PreSumMatrix as NumMatrix }
