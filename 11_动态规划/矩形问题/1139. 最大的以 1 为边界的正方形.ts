// 请你找出边界全部由 1 组成的最大 正方形 子网格，并返回该子网格中的元素数量。如果不存在，则返回 0。
// !如果知道每个点右边/下面有多少个连续的1 会好做很多
function largest1BorderedSquare(matrix: number[][]): number {
  let res = 0
  const ROW = matrix.length
  const COL = matrix[0].length
  const countDown = Array.from<number, number[]>({ length: ROW }, () => Array(COL).fill(0))
  const countRight = Array.from<number, number[]>({ length: ROW }, () => Array(COL).fill(0))

  // 这里可以用python defaultdict的思想简化
  for (let r = ROW - 1; r >= 0; r--) {
    for (let c = COL - 1; c >= 0; c--) {
      if (matrix[r][c] === 1) {
        countDown[r][c] = (countDown[r + 1]?.[c] ?? 0) + 1
        countRight[r][c] = (countRight[r]?.[c + 1] ?? 0) + 1
      }
    }
  }

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (matrix[r][c] !== 1) continue
      const curSize = Math.min(countDown[r][c], countRight[r][c])
      for (let size = curSize; size > res; size--) {
        // 候选大到小找到一个就终止
        // counyDown 保证每列下面都有size个0
        // countRight 保证每行右边都有szie个0
        if (countRight[r + size - 1][c] >= size && countDown[r][c + size - 1] >= size) {
          res = Math.max(res, size)
          break
        }
      }
    }
  }

  return res ** 2
}
