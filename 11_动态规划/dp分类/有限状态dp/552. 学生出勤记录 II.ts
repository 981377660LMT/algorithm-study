/**
 * @param {number} n   1 <= n <= 105
 * @return {number}
 * 返回记录长度为 n 时，可能获得出勤奖励的记录情况 数量
 * dp[i][j][k]表示第 i 天、在 A 为 j 次、连续的 L 为 k 次的方案数
 * 泰波那契数列  (n+2)∗2**(n−1)
 *
 */
var checkRecord = function (n: number): number {
  const MOD = 1e9 + 7
  const dp = Array.from({ length: n + 1 }, () => Array.from({ length: 3 }, () => Array(2).fill(0)))
  dp[0][0][0] = 1
  // console.log(dp[1][2])
  // 六种状态
  for (let i = 1; i <= n; i++) {
    dp[i][0][0] += dp[i - 1][0][0] % MOD
    dp[i][0][0] += dp[i - 1][1][0] % MOD
    dp[i][0][0] += dp[i - 1][2][0] % MOD

    dp[i][0][1] += dp[i - 1][0][1] % MOD
    dp[i][0][1] += dp[i - 1][1][1] % MOD
    dp[i][0][1] += dp[i - 1][2][1] % MOD
    dp[i][0][1] += dp[i - 1][0][0] % MOD
    dp[i][0][1] += dp[i - 1][1][0] % MOD
    dp[i][0][1] += dp[i - 1][2][0] % MOD

    dp[i][1][0] += dp[i - 1][0][0] % MOD

    dp[i][1][1] += dp[i - 1][0][1] % MOD

    dp[i][2][0] += dp[i - 1][1][0] % MOD

    dp[i][2][1] += dp[i - 1][1][1] % MOD
  }

  return (dp[n][0][0] + dp[n][0][1] + dp[n][1][0] + dp[n][1][1] + dp[n][2][0] + dp[n][2][1]) % MOD
}

console.log(checkRecord(2))
// 输出：8
// 解释：
// 有 8 种长度为 2 的记录将被视为可奖励：
// "PP" , "AP", "PA", "LP", "PL", "AL", "LA", "LL"
// 只有"AA"不会被视为可奖励，因为缺勤次数为 2 次（需要少于 2 次）。

export {}
