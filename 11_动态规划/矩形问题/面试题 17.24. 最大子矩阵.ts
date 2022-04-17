class NumMatrix {
  private pre: number[][]

  constructor(matrix: number[][]) {
    // 加一便于处理
    const m = matrix.length
    const n = matrix[0].length
    const pre = Array.from<number, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
    for (let i = 1; i <= m; i++) {
      for (let j = 1; j <= n; j++) {
        // 注意这里的减1
        pre[i][j] = matrix[i - 1][j - 1] + pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1]
      }
    }
    this.pre = pre
  }

  sumRegion(row1: number, col1: number, row2: number, col2: number) {
    return (
      this.pre[row2 + 1][col2 + 1] +
      this.pre[row1][col1] -
      this.pre[row2 + 1][col1] -
      this.pre[row1][col2 + 1]
    )
  }
}

// 找出元素总和最大的子矩阵。
// 返回一个数组 [r1, c1, r2, c2]，
// 其中 r1, c1 分别代表子矩阵左上角的行号和列号，
// r2, c2 分别代表右下角的行号和列号。
function getMaxMatrix(matrix: number[][]): number[] {
  let globalMax = -Infinity
  let res = [0, 0, 0, 0]
  const m = matrix.length
  const n = matrix[0].length
  const numMatrix = new NumMatrix(matrix)

  // 先固定上下两条边
  for (let top = 0; top < m; top++) {
    for (let bottom = top; bottom < m; bottom++) {
      let localMax = 0
      let left = 0

      // 然后从左往右一遍扫描找最大子序和
      for (let right = 0; right < n; right++) {
        localMax = numMatrix.sumRegion(top, left, bottom, right)

        if (localMax > globalMax) {
          res = [top, left, bottom, right]
          globalMax = localMax
        }

        // left直接跳到下一个，不取这个
        if (localMax < 0) {
          localMax = 0
          left = right + 1
        }
      }
    }
  }

  return res
}

console.log(
  getMaxMatrix([
    [-1, 0],
    [0, -1],
  ])
)

export {}
