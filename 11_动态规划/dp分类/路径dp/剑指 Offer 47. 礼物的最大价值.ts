function maxValue(grid: number[][]): number {
  const m = grid.length
  const n = grid[0].length
  const dp = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))

  dp[0][0] = grid[0][0]

  for (let i = 1; i < m; i++) {
    dp[i][0] = grid[i][0] + dp[i - 1][0]
  }

  for (let j = 1; j < n; j++) {
    dp[0][j] = grid[0][j] + dp[0][j - 1]
  }

  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1]) + grid[i][j]
    }
  }

  return dp[m - 1][n - 1]
}

console.log(
  maxValue([
    [1, 3, 1],
    [1, 5, 1],
    [4, 2, 1],
  ])
)
