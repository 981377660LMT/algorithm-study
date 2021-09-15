/**
 * @param {number[][]} matrix
 * @return {void} Do not return anything, modify matrix in-place instead.
 * 如果一个元素为 0 ，则将其所在行和列的所有元素都设为 0 。请使用 原地 算法。
 * @description 泡泡堂
 * O(1):使用了第 0 行和第 0 列来保存 matrix[1:M][1:N] 中是否出现了 0
 */
const setZeroes = function (matrix: number[][]): void {
  const m = matrix.length
  const n = matrix[0].length
  const firstRowHasZero = matrix[0].some(v => v === 0)
  const firstColHasZero = matrix.map(row => row[0]).some(v => v === 0)

  // 标记
  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      matrix[i][j] === 0 && ((matrix[0][j] = 0), (matrix[i][0] = 0))
    }
  }

  // 置零
  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      if (matrix[0][j] === 0 || matrix[i][0] === 0) matrix[i][j] = 0
    }
  }

  // 清理
  if (firstRowHasZero) {
    for (let j = 0; j < n; j++) {
      matrix[0][j] = 0
    }
  }

  if (firstColHasZero) {
    for (let i = 0; i < m; i++) {
      matrix[i][0] = 0
    }
  }
}

console.log(
  setZeroes([
    [0, 1, 2, 0],
    [3, 4, 5, 2],
    [1, 3, 1, 5],
  ])
)

export {}
