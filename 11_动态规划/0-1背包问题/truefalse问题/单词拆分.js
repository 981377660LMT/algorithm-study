/**
 * @param {string} s
 * @param {string[]} wordDict
 * @return {boolean}
 * @description 是否可以用 wordDict 中的词组合成 s，完全背包问题
 * @description 考虑排列顺序的完全背包问题 容量是s,length worddict是物体
 */
var wordBreak = function (s, wordDict) {
  if (!wordDict) return false
  const set = new Set(wordDict)
  const dp = Array(s.length + 1).fill(false)
  // dp[i]表示以 i 结尾的字符串是否可以被 wordDict 中组合而成
  dp[0] = true

  for (let i = 0; i <= s.length; i++) {
    for (const word of set) {
      if (i >= word.length && set.has(s.slice(i - word.length, i)))
        dp[i] = dp[i] || dp[i - word.length]
    }
  }

  return dp[s.length]
}

console.log(wordBreak('applepenapple', ['apple', 'pen']))
