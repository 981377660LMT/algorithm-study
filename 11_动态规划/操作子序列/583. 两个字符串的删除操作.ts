/**
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * 相对于72. 编辑距离,只有删除操作
 * 找到使得 word1 和 word2 相同所需的最小步数，每步可以删除任意一个字符串中的一个字符。
 * @summary
 * dp[i][j]：以i-1为结尾的字符串word1，和以j-1位结尾的字符串word2，想要达到相等，所需要删除元素的最少次数。
 */
const minDistance = function (word1: string, word2: string): number {
  const dp = Array.from({ length: word1.length + 1 }, () => Array(word2.length + 1).fill(0))
  for (let i = 0; i <= word1.length; i++) {
    dp[i][0] = i
  }
  for (let j = 0; j <= word2.length; j++) {
    dp[0][j] = j
  }

  for (let i = 1; i <= word1.length; i++) {
    for (let j = 1; j <= word2.length; j++) {
      if (word1[i - 1] === word2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
      } else {
        // 左，上，对角线左上
        dp[i][j] = Math.min(dp[i - 1][j] + 1, dp[i][j - 1] + 1, dp[i - 1][j - 1] + 2)
      }
    }
  }

  return dp[word1.length][word2.length]
}

console.log(minDistance('rabbbit', 'rabbit'))

export default 1
