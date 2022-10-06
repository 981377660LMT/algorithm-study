/**
 *
 * @param s1
 * @param s2
 * @param s3
 */
function isInterleave(s1: string, s2: string, s3: string): boolean {
  const [m, n] = [s1.length, s2.length]
  if (m + n !== s3.length) return false

  const dp = Array.from({ length: m + 1 }, () => new Uint8Array(n + 1))
  dp[0][0] = 1

  // !注意取0个字符的情况
  for (let i = 0; i < m + 1; i++) {
    for (let j = 0; j < n + 1; j++) {
      if (i > 0) {
        if (s1[i - 1] === s3[i + j - 1]) {
          dp[i][j] |= dp[i - 1][j]
        }
      }

      if (j > 0) {
        if (s2[j - 1] === s3[i + j - 1]) {
          dp[i][j] |= dp[i][j - 1]
        }
      }
    }
  }

  return !!dp[m][n]
}

console.log(isInterleave('aabcc', 'dbbca', 'aadbbcbcac'))

export {}
