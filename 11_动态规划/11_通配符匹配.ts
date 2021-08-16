/**
 * @param {string} s
 * @param {string} p
 * @return {boolean}
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
    console.log(p[j - 1], dp[0][0], j)
    dp[0][j] = dp[0][j - 1] && p[j - 1] === '*'
  }

  for (let i = 1; i < s.length + 1; i++) {
    for (let j = 1; j < p.length + 1; j++) {
      if (p[j - 1] === '*') {
        // 前面已经匹配，加了之后更加匹配
        // 前面未匹配，缩短源字符串看是否匹配 'aa', '*'
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
