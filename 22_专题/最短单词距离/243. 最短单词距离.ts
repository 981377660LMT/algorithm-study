// word1 不等于 word2, 并且 word1 和 word2 都在列表里
const shortestDistance = (wordsDict: string[], word1: string, word2: string): number => {
  let l = Infinity,
    r = Infinity,
    res = Infinity

  for (let i = 0; i < wordsDict.length; i++) {
    const cur = wordsDict[i]
    cur === word1 && (l = i)
    cur === word2 && (r = i)
    if (l !== Infinity && r !== Infinity) {
      res = Math.min(res, Math.abs(l - r))
    }
  }

  return res
}

export {}
// 假设 words = ["practice", "makes", "perfect", "coding", "makes"]
// 输入: word1 = “coding”, word2 = “practice”
// 输出: 3
// 输入: word1 = "makes", word2 = "coding"
// 输出: 1
