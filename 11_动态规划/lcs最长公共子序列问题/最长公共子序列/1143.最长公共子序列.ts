/**
 * LCS模板
 * @param {string} text1
 * @param {string} text2
 * @return {number}
 */
function longestCommonSubsequence(text1: string, text2: string): number {
  const n1 = text1.length
  const n2 = text2.length
  const dp = Array.from({ length: n1 + 1 }, () => new Uint32Array(n2 + 1))
  for (let i = 1; i < n1 + 1; i++) {
    for (let j = 1; j < n2 + 1; j++) {
      if (text1[i - 1] === text2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1
      } else {
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1])
      }
    }
  }
  return dp[n1][n2]
}

console.log(longestCommonSubsequence('abcde', 'ace'))
console.log(longestCommonSubsequence('abca', 'acba'))

export { longestCommonSubsequence }
