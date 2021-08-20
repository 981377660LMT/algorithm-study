// 回文子串是要连续的，回文子序列可不是连续的！
/**
 * @param {string} s
 * @return {number}
 * dp[i][j]：字符串s在[i, j]范围内最长的回文子序列的长度为dp[i][j]。
 * 如果s[i]与s[j]相同，那么dp[i][j] = dp[i + 1][j - 1] + 2;
 */
const longestPalindromeSubseq = function (s: string): number {
  const len = s.length
  const dp = Array.from<number, number[]>({ length: len }, () => Array(len).fill(0))
  for (let i = 0; i < len; i++) {
    dp[i][i] = 1
  }

  // 倒序遍历
  for (let i = len - 1; i >= 0; i--) {
    for (let j = i + 1; j < len; j++) {
      if (s[i] === s[j]) {
        dp[i][j] = dp[i + 1][j - 1] + 2
      } else {
        dp[i][j] = Math.max(dp[i + 1][j], dp[i][j - 1])
      }
    }
  }

  return dp[0][len - 1]
}

console.log(longestPalindromeSubseq('bbbab'))

export default 1
