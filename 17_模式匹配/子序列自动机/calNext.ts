/**
 * 子序列自动机.这里的字符集为小写字母.
 * dp[i][c] := i 第i(0<=i<=n)个字符以后第一次出现字符c的位置(不存在的话为n).
 */
function calNext(s: string): number[][] {
  const n = s.length
  const dp: number[][] = Array(n + 1)
  for (let i = 0; i <= n; i++) {
    dp[i] = Array(26).fill(n)
  }
  for (let i = n - 1; ~i; i--) {
    for (let j = 0; j < 26; j++) {
      dp[i][j] = dp[i + 1][j]
    }
    dp[i][s.charCodeAt(i) - 97] = i
  }
  return dp
}

export { calNext }

if (require.main === module) {
  const s = 'ab'
  console.log(calNext(s))
}
