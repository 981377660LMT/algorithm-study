// 最大正方形(全1的最大正方形)

/**
 * 二维矩形区域中最大的正方形.
 * @param grid 二维矩形区域."1"或者1表示有效区域,"0"或者0表示无效区域.
 * @returns 最大正方形的边长和`[r1,r2) x [c1,c2)`区域.
 */
function maxSquare(
  grid: ArrayLike<ArrayLike<string | number | boolean>>
): [len: number, rectangle: [r1: number, r2: number, c1: number, c2: number]] {
  const ROW = grid.length
  const COL = ROW ? grid[0].length : 0
  let res1 = 0
  let res2: [number, number, number, number] = [0, 0, 0, 0]
  let dp = new Uint32Array(COL)

  for (let r = 0; r < ROW; r++) {
    const ndp = new Uint32Array(COL)
    for (let c = 0; c < COL; c++) {
      if (+grid[r][c] === 1) {
        if (c === 0) {
          ndp[c] = 1
        } else {
          ndp[c] = Math.min(dp[c - 1], dp[c], ndp[c - 1]) + 1
        }
        if (ndp[c] > res1) {
          res1 = ndp[c]
          res2 = [r + 1 - res1, r + 1, c + 1 - res1, c + 1]
        }
      }
    }
    dp = ndp
  }

  return [res1, res2]
}

export { maxSquare }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function maximalSquare(matrix: string[][]): number {
    const res = maxSquare(matrix)
    return res[0] * res[0]
  }
}
