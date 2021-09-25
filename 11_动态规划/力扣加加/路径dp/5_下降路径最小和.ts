// 从 arr 数组中的每一行选择一个数字，且按顺序选出来的数字中，相邻数字不在原数组的同一列。
function minFallingPathSum(grid: number[][]): number {
  const len = grid.length
  const dp = Array.from<number, number[]>({ length: len }, () => Array(len).fill(0))
  dp[0] = grid[0]

  for (let i = 1; i < len; i++) {
    let min = Math.min.apply(null, dp[i - 1])
    for (let j = 0; j < len; j++) {
      if (dp[i - 1][j] !== min) dp[i][j] = min + grid[i][j]
      else
        dp[i][j] =
          Math.min(
            Math.min.apply(null, dp[i - 1].slice(0, j)),
            Math.min.apply(null, dp[i - 1].slice(j + 1))
          ) + grid[i][j]
    }
  }
  console.table(dp)
  return Math.min.apply(null, dp[len - 1])
}

console.log(
  minFallingPathSum([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
  ])
)
// 下降路径中数字和最小的是 [1,5,7] ，所以答案是 13 。
export {}
