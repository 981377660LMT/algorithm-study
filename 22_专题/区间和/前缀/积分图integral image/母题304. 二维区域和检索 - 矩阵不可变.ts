// 计算其子矩形范围内元素的总和，该子矩阵的左上角为 (row1, col1) ，右下角为 (row2, col2) 。
class NumMatrix {
  private matrix: number[][]
  private pre: number[][]

  constructor(matrix: number[][]) {
    this.matrix = matrix
    // 加一便于处理
    const m = matrix.length + 1
    const n = matrix[0].length + 1
    const pre = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))
    for (let i = 1; i < m; i++) {
      for (let j = 1; j < n; j++) {
        // 注意这里的减1
        pre[i][j] = matrix[i - 1][j - 1] + pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1]
      }
    }
    this.pre = pre
  }

  // 注意这里的减1
  sumRegion(row1: number, col1: number, row2: number, col2: number) {
    return (
      this.pre[row2 + 1][col2 + 1] +
      this.pre[row1][col1] -
      this.pre[row2 + 1][col1] -
      this.pre[row1][col2 + 1]
    )
  }
}

const nm = new NumMatrix([
  [1, 2, 3],
  [4, 5, 6],
  [7, 8, 9],
])

console.log(nm.sumRegion(1, 1, 2, 2))
