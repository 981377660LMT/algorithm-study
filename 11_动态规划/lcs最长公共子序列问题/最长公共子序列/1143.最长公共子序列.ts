/**
 * @param {string} text1
 * @param {string} text2
 * @return {number}
 */
function longestCommonSubsequence(text1: string, text2: string): number {
  let res = 0
  const m = text1.length
  const n = text2.length
  const dp = Array.from<number, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i < m + 1; i++) {
    for (let j = 1; j < n + 1; j++) {
      if (text1[i - 1] === text2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1
        res = Math.max(res, dp[i][j])
      } else {
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1])
      }
    }
  }
  return res
}

console.log(longestCommonSubsequence('abcde', 'ace'))
console.log(longestCommonSubsequence('abca', 'acba'))
