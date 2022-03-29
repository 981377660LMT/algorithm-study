// 从 arr 数组中的每一行选择一个数字，且按顺序选出来的数字中，相邻数字不在原数组的同一列。
function minFallingPathSum(matrix: number[][]): number {
  const len = matrix.length
  // dp[i][j], 走到i，j时的最小矩阵元素和
  const dp = Array.from<number, number[]>({ length: len }, () => Array(len).fill(0))
  dp[0] = matrix[0]

  for (let i = 1; i < len; i++) {
    for (let j = 0; j < len; j++) {
      dp[i][j] =
        Math.min(dp[i - 1][j], dp[i - 1][j - 1] ?? Infinity, dp[i - 1][j + 1] ?? Infinity) +
        matrix[i][j]
    }
  }

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
