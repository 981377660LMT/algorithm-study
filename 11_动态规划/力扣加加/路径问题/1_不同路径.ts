const uniquePath = (m: number, n: number): number => {
  const dp: number[][] = Array.from({ length: m }, () => Array(n).fill(Infinity))

  for (let i = 0; i < n; i++) {
    dp[0][i] = 1
  }

  for (let i = 0; i < m; i++) {
    dp[i][0] = 1
  }

  for (let i = 1; i < m; i++) {
    for (let j = 1; j < n; j++) {
      dp[i][j] = dp[i - 1][j] + dp[i][j - 1]
    }
  }

  console.table(dp)
  return dp[m - 1][n - 1]
}

console.log(uniquePath(3, 7))

export {}
