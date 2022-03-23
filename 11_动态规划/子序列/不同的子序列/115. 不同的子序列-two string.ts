/**
 * @param {string} s   0 <= s.length, t.length <= 1000
 * @param {string} t
 * @return {number}
 * 相对于72. 编辑距离,只有删除操作
 * 不同的子序列数
 * @summary dp[i][j]：以i-1为结尾的s子序列中出现以j-1为结尾的t的个数为dp[i][j]。
 * @link https://leetcode-cn.com/problems/distinct-subsequences/solution/shou-hua-tu-jie-xiang-jie-liang-chong-ji-4r2y/
 */
const numDistinct = function (s: string, t: string): number {
  const dp = Array.from({ length: s.length + 1 }, () => Array(t.length + 1).fill(0))
  for (let i = 0; i < s.length + 1; i++) {
    dp[i][0] = 1
  }

  for (let i = 1; i < s.length + 1; i++) {
    for (let j = 1; j < t.length + 1; j++) {
      if (s[i - 1] === t[j - 1]) {
        dp[i][j] = dp[i - 1][j] + dp[i - 1][j - 1] // 选不选最后一位 选就是dp[i-1][j-1] 不选就是 dp[i-1][j]
      } else {
        dp[i][j] = dp[i - 1][j]
      }
    }
  }

  return dp[s.length][t.length]
}

console.log(numDistinct('rabbbit', 'rabbit'))

export default 1
