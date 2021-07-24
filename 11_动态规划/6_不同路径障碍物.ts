/**
 * @param {number[][]} obstacleGrid
 * @return {number}
 */
const uniquePathsWithObstacles = (obstacleGrid: number[][]): number => {
  const m = obstacleGrid.length
  const n = obstacleGrid[0].length
  const dp: number[][] = Array.from({ length: m }, () => Array(n).fill(0))
  const zeroOrN = (n: number) => (n === Infinity ? 0 : n)

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (obstacleGrid[i][j] === 1) dp[i][j] = Infinity
    }
  }

  for (let i = 0; i < m; i++) {
    dp[i][0] = dp[i][0] === Infinity ? Infinity : 1
  }
  for (let i = 0; i < n; i++) {
    dp[0][i] = dp[0][i] === Infinity ? Infinity : 1
  }

  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      if (dp[i][j] === Infinity) continue
      dp[i][j] = zeroOrN(dp[i - 1][j]) + zeroOrN(dp[i][j - 1])
    }
  }

  console.table(dp)
  return dp[m - 1][n - 1]
}

console.log(
  uniquePathsWithObstacles([
    [0, 1],
    [0, 0],
  ])
)
// 输出：2
export {}
