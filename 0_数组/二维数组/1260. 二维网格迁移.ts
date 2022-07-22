/**
 * @param {number[][]} grid
 * @param {number} k
 * @return {number[][]}
 */
function shiftGrid(grid: number[][], k: number): number[][] {
  const [row, col] = [grid.length, grid[0].length]
  const n = row * col
  const res: number[][] = Array.from({ length: row }, () => Array(col).fill(0))

  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      const [index, value] = [r * col + c, grid[r][c]]
      const nextIndex = (index + k) % n
      const [nextR, nextC] = [Math.floor(nextIndex / col), nextIndex % col]
      res[nextR][nextC] = value
    }
  }

  return res
}

console.log(
  shiftGrid(
    [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9],
    ],
    1
  )
)
