/**
 * @param {string} s1
 * @param {string} s2
 * @return {number}
 * 231
 */
const minimumDeleteSum = function (s1: string, s2: string): number {
  const dp = Array.from({ length: s1.length + 1 }, () => Array(s2.length + 1).fill(0))
  for (let i = 1; i <= s1.length; i++) {
    dp[i][0] = dp[i - 1][0] + s1[i - 1].codePointAt(0)!
  }

  for (let j = 1; j <= s2.length; j++) {
    dp[0][j] = dp[0][j - 1] + s2[j - 1].codePointAt(0)!
  }
  console.table(dp)
  for (let i = 1; i <= s1.length; i++) {
    for (let j = 1; j <= s2.length; j++) {
      if (s1[i - 1] === s2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
      } else {
        // 左，上，对角线左上
        dp[i][j] = Math.min(
          dp[i - 1][j] + s1[i - 1].codePointAt(0)!,
          dp[i][j - 1] + s2[j - 1].codePointAt(0)!,
          dp[i - 1][j - 1] + 2 + s1[i - 1].codePointAt(0)! + s2[j - 1].codePointAt(0)!
        )
      }
    }
  }

  return dp[s1.length][s2.length]
}

console.log(minimumDeleteSum('sea', 'eat'))
