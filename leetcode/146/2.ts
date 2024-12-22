export {}

const MOD = 1e9 + 7
function countPathsWithXorValue(grid: number[][], k: number): number {
  const m = grid.length
  const n = grid[0].length

  let dp: number[][] = Array.from({ length: n }, () => Array(16).fill(0))
  dp[0][grid[0][0] % 16] = 1

  for (let j = 1; j < n; j++) {
    for (let x = 0; x < 16; x++) {
      const v = x ^ grid[0][j]
      dp[j][v] = (dp[j][v] + dp[j - 1][x]) % MOD
    }
  }

  for (let i = 1; i < m; i++) {
    const ndp: number[][] = Array.from({ length: n }, () => Array(16).fill(0))
    for (let j = 0; j < n; j++) {
      for (let x = 0; x < 16; x++) {
        if (dp[j][x] > 0) {
          const v = x ^ grid[i][j]
          ndp[j][v] = (ndp[j][v] + dp[j][x]) % MOD
        }
        if (j > 0 && ndp[j - 1][x] > 0) {
          const v = x ^ grid[i][j]
          ndp[j][v] = (ndp[j][v] + ndp[j - 1][x]) % MOD
        }
      }
    }
    dp = ndp
  }

  return dp[n - 1][k] || 0
}

console.log(
  countPathsWithXorValue(
    [
      [2, 1, 5],
      [7, 10, 0],
      [12, 6, 4]
    ],
    11
  )
)
