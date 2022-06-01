/**
 * @param {number} d  d 个一样的骰子
 * @param {number} f  每个骰子上都有 f 个面
 * @param {number} target  需要掷出的总点数为 target
 * @return {number}
 * @description
 * f[i][j] 表示考虑前 i个物品组，凑成价值为 j 的方案数。
 */
const numRollsToTarget = function (d: number, f: number, target: number): number {
  const dp = Array.from<number, number[]>({ length: d + 1 }, () => Array(target + 1).fill(0))
  dp[0][0] = 1

  // 枚举物品组（每个骰子）
  for (let i = 1; i <= d; i++) {
    // 枚举背包容量（所掷得的总点数）
    for (let j = 0; j <= target; j++) {
      // 枚举决策（当前骰子所掷得的点数）
      for (let k = 1; k <= f; k++) {
        j - k >= 0 && (dp[i][j] = (dp[i][j] + dp[i - 1][j - k]) % (10 ** 9 + 7))
      }
    }
  }

  return dp[d][target]
}

console.log(numRollsToTarget(2, 6, 7))
