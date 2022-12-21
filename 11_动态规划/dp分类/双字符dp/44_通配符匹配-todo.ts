/**
 * @param {string} s
 * @param {string} p
 * @return {boolean}
 * 注意 这里不是正则的正确用法 只是举例子
 * '?' 可以匹配任何单个字符。
   '*' 可以匹配任意字符串（包括空字符串）。例子: a* 可以与 ab 匹配 ,不能与空字符串匹配
   dp 维度加一在字符串动态规划问题中很常见，特别是要考虑空字符串的情况下，
 */
const isMatch = function (s: string, p: string): boolean {
  const dp = Array.from<boolean, boolean[]>({ length: s.length + 1 }, () =>
    Array(p.length + 1).fill(false)
  )

  dp[0][0] = true
  for (let i = 1; i < s.length + 1; i++) {
    dp[i][0] = false
  }

  for (let j = 1; j < p.length + 1; j++) {
    dp[0][j] = dp[0][j - 1] && p[j - 1] === '*'
  }

  for (let i = 1; i < s.length + 1; i++) {
    for (let j = 1; j < p.length + 1; j++) {
      if (p[j - 1] === '*') {
        // 1.s少一位
        // 2.p 少一位
        dp[i][j] = dp[i - 1][j] || dp[i][j - 1]
      } else if (p[j - 1] === s[i - 1] || p[j - 1] === '?') {
        dp[i][j] = dp[i - 1][j - 1]
      }
    }
  }
  console.table(dp)
  return dp[s.length][p.length]
}

// console.log(isMatch('adceb', '*a*b'))
console.log(isMatch('aa', '*'))

export default 1
