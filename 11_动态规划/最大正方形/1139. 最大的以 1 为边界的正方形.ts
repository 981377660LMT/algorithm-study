// 请你找出边界全部由 1 组成的最大 正方形 子网格，并返回该子网格中的元素数量。如果不存在，则返回 0。
// 垂线dp
function largest1BorderedSquare(matrix: number[][]): number {
  let res = 0
  const m = matrix.length
  const n = matrix[0].length
  const countDown = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))
  const countRight = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))

  // 这里可以用python defaultdict的思想简化
  for (let i = m - 1; i >= 0; i--) {
    for (let j = n - 1; j >= 0; j--) {
      if (matrix[i][j] === 1) {
        countDown[i][j] = i + 1 >= m ? 1 : countDown[i + 1][j] + 1
        countRight[i][j] = j + 1 >= n ? 1 : countRight[i][j + 1] + 1
      }
    }
  }

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (matrix[i][j] !== 1) continue
      const curSize = Math.min(countDown[i][j], countRight[i][j])
      for (let size = curSize; size > res; size--) {
        // 候选大到小找到一个就终止
        // counyDown 保证每列下面都有size个0
        // countRight 保证每行右边都有szie个0
        if (countRight[i + size - 1][j] >= size && countDown[i][j + size - 1] >= size) {
          res = size
          break
        }
      }
    }
  }

  return res ** 2
}
