// 你不小心把一个长篇文章中的空格、标点都删掉了，并且大写也弄成了小写。
// 设计一个算法，把文章断开，要求未识别的字符最少，返回未识别的字符数。
// 有点像有序的完全背包
function respace(dictionary: string[], sentence: string): number {
  const len = sentence.length
  if (len === 0) return 0
  const dp = Array<number>(len + 1).fill(0)

  for (let i = 1; i <= len; i++) {
    // 初始化：如果没有匹配那么dp[i]相比于dp[i-1]直接多1
    dp[i] = dp[i - 1] + 1
    // 如果新加一个字符，组成了一个词的情况
    for (const word of dictionary) {
      const len = word.length
      if (sentence.slice(i - len, i) === word && i - len >= 0) dp[i] = Math.min(dp[i], dp[i - len])
    }
  }

  return dp[len]
}

console.log(
  respace(['looked', 'just', 'like', 'her', 'brother'], 'jesslookedjustliketimherbrother')
)

export {}
