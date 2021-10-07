/**
 * @param {string} s
 * @param {string} p
 * @return {boolean}
 * @description 给你一个字符串 s 和一个字符规律 p，请你来实现一个支持 '.' 和 '*' 的正则表达式匹配。
   '.' 匹配任意单个字符
   '*' 匹配零个或多个前面的那一个元素 例子：a* 可以与空字符串匹配
   由于 Python 等很多语言字符串都是不可变的，因此会内存开销可能比较大。优化的方式也很简单，只需要记录 pattern 和 text 的索引即可。
 */
const isMatch = function (s: string, p: string): boolean {
  const m = s.length
  const n = p.length
  const dp = Array.from<boolean, boolean[]>({ length: m + 1 }, () => Array(n + 1).fill(false))
  // 两个空串，可以匹配
  dp[0][0] = true
  for (let i = 1; i < m + 1; i++) {
    dp[i][0] = false
  }
  // 如果遇到了*,只要判断其对应的前面两个元素的dp值
  for (let j = 2; j < n + 1; j++) {
    if (p[j - 1] === '*') {
      dp[0][j] = dp[0][j - 2]
    }
  }

  for (let i = 1; i < m + 1; i++) {
    for (let j = 1; j < n + 1; j++) {
      if (p[j - 1] === '*') {
        if (dp[i][j - 2]) dp[i][j] = true
        else if (dp[i - 1][j] && s[i - 1] === p[j - 2]) dp[i][j] = true
        else if (dp[i - 1][j] && p[j - 2] === '.') dp[i][j] = true
      } else {
        if (dp[i - 1][j - 1] && s[i - 1] === p[j - 1]) dp[i][j] = true
        else if (dp[i - 1][j - 1] && p[j - 1] === '.') dp[i][j] = true
      }
    }
  }

  // console.table(dp)
  return dp[m][n]
}

console.log(isMatch('aab', 'c*a*b'))
// const a = [1, 2]
// console.log(a[-1])

export {}
