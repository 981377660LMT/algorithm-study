/**
 * @param {string[]} words 0 <= sum(words[i].length) <= 6 * 10**5  不能dp但是有参考价值
 * @return {string[]}
 * @description 返回 words 中的所有 连接词 。
 * 连接词 的定义为：一个字符串完全是由至少两个给定数组中的单词组成的。
 * @summary
 * 对于下面的极端用例会超时 哈希表遍历字符串消耗很大(字符串的问题)
 * 需要用trie
 */
const findAllConcatenatedWordsInADict = function (words: string[]): string[] {
  const res: string[] = []
  const preWords: string[] = []
  words.sort((a, b) => a.length - b.length)

  // 139_单词拆分1
  const wordBreak = function (s: string, wordDict: string[]) {
    if (wordDict.length === 0) return false
    // 背包
    const set = new Set(wordDict)
    // dp[i]表示以 i 结尾的字符串是否可以被 wordDict 中组合而成
    const dp = Array(s.length + 1).fill(false)
    dp[0] = true

    for (let i = 0; i <= s.length; i++) {
      for (const word of set) {
        if (i >= word.length && dp[i - word.length] && s.slice(i - word.length, i) === word)
          dp[i] = true
      }
    }

    return dp[s.length]
  }

  for (const word of words) {
    if (wordBreak(word, preWords)) res.push(word)
    preWords.push(word)
  }

  return res
}

console.log(
  findAllConcatenatedWordsInADict([
    'cat',
    'cats',
    'catsdogcats',
    'dog',
    'dogcatsdog',
    'hippopotamuses',
    'rat',
    'ratcatdogcat',
  ])
)

// ["catsdogcats","dogcatsdog","ratcatdogcat"]

export {}
