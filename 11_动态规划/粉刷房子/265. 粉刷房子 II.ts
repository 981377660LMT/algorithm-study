function minCostII(costs: number[][]): number {
  const n = costs.length
  const k = costs[0].length
  const dp = Array.from<number, number[]>({ length: n }, () => Array(k).fill(0))
  dp[0] = costs[0]

  for (let i = 1; i < n; i++) {
    for (let j = 0; j < k; j++) {
      dp[i][j] =
        Math.min(
          Math.min.apply(null, dp[i - 1].slice(0, j)),
          Math.min.apply(null, dp[i - 1].slice(j + 1))
        ) + costs[i][j]
    }
  }

  return Math.min.apply(null, dp[n - 1])
}

export {}
console.log(
  minCostII([
    [1, 3],
    [2, 4],
  ])
)
