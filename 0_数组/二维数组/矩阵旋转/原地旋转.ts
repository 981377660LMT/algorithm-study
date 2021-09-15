// 先转置后镜像对称
/**
 Do not return anything, modify matrix in-place instead.
 NxN的矩阵
 */
const rotate = (matrix: number[][]): void => {
  // 转置
  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < i; j++) {
      ;[matrix[i][j], matrix[j][i]] = [matrix[j][i], matrix[i][j]]
    }
  }

  // 镜像
  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < matrix.length >> 1; j++) {
      ;[matrix[i][j], matrix[i][matrix.length - j - 1]] = [
        matrix[i][matrix.length - j - 1],
        matrix[i][j],
      ]
    }
  }

  console.table(matrix)
}

rotate([
  [1, 2, 3],
  [4, 5, 6],
  [7, 8, 9],
])

export {}
