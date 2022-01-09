/**
 * @param {number[][]} mat
 * @param {number} k
 * @return {number[][]}
 * i - k <= r <= i + k,
   j - k <= c <= j + k 且
   (r, c) 在矩阵内。
   类似于深度学习中的卷积操作
 */
const matrixBlockSum = function (mat: number[][], k: number): number[][] {
  // 加一方便处理
  const m = mat.length + 1
  const n = mat[0].length + 1
  const pre = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))
  const res: number[][] = Array.from<number, number[]>({ length: m - 1 }, () =>
    Array(n - 1).fill(0)
  )
  const sumRegion = (row1: number, col1: number, row2: number, col2: number) => {
    return pre[row2 + 1][col2 + 1] - pre[row2 + 1][col1] - pre[row1][col2 + 1] + pre[row1][col1]
  }

  // 求前缀和
  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      pre[i][j] = mat[i - 1][j - 1] + pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1]
    }
  }

  console.table(pre)

  for (let i = 0; i < m - 1; i++) {
    for (let j = 0; j < n - 1; j++) {
      const right = Math.min(j + k, n - 1 - 1)
      const bottom = Math.min(i + k, m - 1 - 1)
      const left = Math.max(0, j - k)
      const top = Math.max(0, i - k)
      res[i][j] = sumRegion(top, left, bottom, right)
    }
  }

  console.table(res)
  return res
}

console.log(
  matrixBlockSum(
    [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9],
    ],
    1
  )
)

export default 1
