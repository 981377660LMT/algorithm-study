export {}

function minimumOperations(grid: number[][]): number {
  const m = grid.length
  const n = grid[0].length
  let res = 0
  for (let c = 0; c < n; c++) {
    for (let r = 1; r < m; r++) {
      if (grid[r][c] <= grid[r - 1][c]) {
        const diff = grid[r - 1][c] + 1 - grid[r][c]
        res += diff
        grid[r][c] += diff
      }
    }
  }
  return res
}
