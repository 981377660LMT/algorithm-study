// 输入n，打印出所有骰子朝上一面的点数之和为s的所有可能的值出现的概率。
function dicesProbability(n: number): number[] {
  const dp = Array.from<unknown, number[]>({ length: n + 1 }, () => Array(6 * n + 1).fill(0))
  for (let i = 1; i <= 6; i++) {
    dp[1][i] = 1
  }

  for (let i = 2; i <= n; i++) {
    for (let j = i; j <= i * 6; j++) {
      for (let k = 1; k <= 6; k++) {
        j - k > 0 && (dp[i][j] += dp[i - 1][j - k])
      }
    }
  }

  const res: number[] = []
  const total = 6 ** n
  for (let i = n; i <= 6 * n; i++) {
    res.push(dp[n][i] / total)
  }

  return res
}

console.log(dicesProbability(2))
// 输入: 2
// 输出: [
//   0.02778, 0.05556, 0.08333, 0.11111, 0.13889, 0.16667, 0.13889, 0.11111, 0.08333, 0.05556, 0.02778,
// ]
// dp[i][j] 代表前 i 个骰子的点数和 j 的概率
