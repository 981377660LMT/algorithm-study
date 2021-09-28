// 返回列表中这两个单词之间的最短距离。
// word1 和 word2 是有可能相同的，并且它们将分别表示为列表中两个独立的单词。
// 你可以假设 word1 和 word2 都在列表里。
function shortestWordDistance(wordsDict: string[], word1: string, word2: string): number {
  let l = Infinity,
    r = Infinity,
    res = Infinity

  if (word1 === word2) {
    // 此时只需要一个l判断即可
    for (let i = 0; i < wordsDict.length; i++) {
      const cur = wordsDict[i]
      if (cur === word1 && l !== Infinity) res = Math.min(res, i - l)
      if (cur === word1) l = i
    }
  } else {
    // 243题 代码
    for (let i = 0; i < wordsDict.length; i++) {
      const cur = wordsDict[i]
      cur === word1 && (l = i)
      cur === word2 && (r = i)
      if (l !== Infinity && r !== Infinity) {
        res = Math.min(res, Math.abs(l - r))
      }
    }
  }

  return res
}

export default 1

console.log(
  shortestWordDistance(['practice', 'makes', 'perfect', 'coding', 'makes'], 'makes', 'makes')
)
// 输出: 3
