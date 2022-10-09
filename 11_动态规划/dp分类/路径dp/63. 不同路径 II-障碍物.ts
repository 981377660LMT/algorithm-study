/**
 * @param {number[][]} obstacleGrid
 * @return {number}
 */
const uniquePathsWithObstacles = (obstacleGrid: number[][]): number => {
  if (obstacleGrid[0][0] === 1) return 0
  const m = obstacleGrid.length
  const n = obstacleGrid[0].length
  const dp: number[][] = Array.from({ length: m }, () => Array(n).fill(0))

  dp[0][0] = obstacleGrid[0][0] === 1 ? 0 : 1

  for (let i = 1; i < m; i++) {
    if (obstacleGrid[i][0] === 1) continue
    dp[i][0] = dp[i - 1][0]
  }

  for (let j = 1; j < n; j++) {
    if (obstacleGrid[0][j] === 1) continue
    dp[0][j] = dp[0][j - 1]
  }

  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      if (obstacleGrid[i][j] === 1) continue
      dp[i][j] = dp[i - 1][j] + dp[i][j - 1]
    }
  }

  console.table(dp)
  return dp[m - 1][n - 1]
}

console.log(
  uniquePathsWithObstacles([
    [0, 1],
    [0, 0]
  ])
)
console.log(
  uniquePathsWithObstacles([
    [0, 0],
    [0, 1]
  ])
)
console.log(
  uniquePathsWithObstacles([
    [0, 0],
    [1, 1],
    [0, 0]
  ])
)
// console.log(uniquePathsWithObstacles([[1, 0]]))
// 输出：2
export {}
