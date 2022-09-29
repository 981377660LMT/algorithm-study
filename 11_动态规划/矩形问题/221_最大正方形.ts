// 在一个由 '0' 和 '1' 组成的二维矩阵内，找到只包含 '1' 的最大正方形，并返回其面积。

// dp[i][j]代表以[i][j]为右下角的正方形的最大边长 这里直接在原矩阵修改
const maximalSquare = (matrix: string[][]) => {
  // const dp = Array.from(Array(matrix.length), () => Array(matrix[0].length).fill(0))
  // console.table(dp)
  let max = 0
  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < matrix[0].length; j++) {
      if (matrix[i][j] === '1') {
        if (i >= 1 && j >= 1) {
          // @ts-ignore
          matrix[i][j] = Math.min(matrix[i][j - 1], matrix[i - 1][j], matrix[i - 1][j - 1]) + 1
        }
      }

      // @ts-ignore
      max = Math.max(max, matrix[i][j])
    }
  }

  console.table(matrix)
  return max ** 2
}

console.log(
  maximalSquare([
    ['1', '0', '1', '0', '0'],
    ['1', '0', '1', '1', '1'],
    ['1', '1', '1', '1', '1'],
    ['1', '0', '0', '1', '0']
  ])
)

export {}
