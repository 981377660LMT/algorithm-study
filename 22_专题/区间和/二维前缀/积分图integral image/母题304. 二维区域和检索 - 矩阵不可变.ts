// 计算其子矩形范围内元素的总和，该子矩阵的左上角为 (row1, col1) ，右下角为 (row2, col2) 。
class NumMatrix {
  private preSum: number[][]

  constructor(matrix: number[][]) {
    const m = matrix.length
    const n = matrix[0].length
    const preSum = Array.from<number, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
    for (let i = 0; i < m; i++) {
      for (let j = 0; j < n; j++) {
        // 注意这里的减1
        preSum[i + 1][j + 1] = matrix[i][j] + preSum[i][j + 1] + preSum[i + 1][j] - preSum[i][j]
      }
    }

    this.preSum = preSum
  }

  /**
   *
   * @param row1
   * @param col1
   * @param row2
   * @param col2
   * @returns 返回 左上角 (row1, col1) 、右下角 (row2, col2) 闭区间所描述的子矩阵的元素 总和 。
   */
  sumRegion(row1: number, col1: number, row2: number, col2: number) {
    return (
      this.preSum[row2 + 1][col2 + 1] +
      this.preSum[row1][col1] -
      this.preSum[row2 + 1][col1] -
      this.preSum[row1][col2 + 1]
    )
  }
}

if (require.main === module) {
  const nm = new NumMatrix([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
  ])

  console.log(nm.sumRegion(1, 1, 2, 2))
}

export { NumMatrix }
