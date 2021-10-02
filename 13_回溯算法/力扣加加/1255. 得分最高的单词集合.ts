/**
 * @param {string[]} words 1 <= words.length <= 14
 * @param {string[]} letters
 * @param {number[]} score
 * @return {number}
 * @description 请你帮忙计算玩家在单词拼写游戏中所能获得的「最高得分」：
 * 能够由 letters 里的字母拼写出的 任意 属于 words 单词子集中，分数最高的单词集合的得分。
 * 单词表 words 中每个单词只能计分（使用）一次。
 * 可以只使用字母表 letters 中的部分字母，但是每个字母最多被使用一次。
 * @summary 暴力枚举解法
 */
const maxScoreWords = function (words: string[], letters: string[], score: number[]) {
  let res = 0
  const len = words.length
  const letterCounter = new Map<string, number>()
  for (const letter of letters) letterCounter.set(letter, (letterCounter.get(letter) || 0) + 1)
  const wordsScore = words.map(word =>
    word
      .split('')
      .map(char => score[char.codePointAt(0)! - 97])
      .reduce((pre, cur) => pre + cur, 0)
  )

  // console.log(wordsScore)
  const countStore = (arr: string[]): number => {
    console.log(arr)
    let res = 0
    const needCounter = new Map<string, number>()
    for (let i = 0; i < arr.length; i++) {
      if (arr[i] === '0') continue
      res += wordsScore[i]
      for (const letter of words[i]) {
        needCounter.set(letter, (needCounter.get(letter) || 0) + 1)
      }
    }
    console.log(needCounter, letterCounter, res)
    for (const letter of needCounter.keys()) {
      if (!letterCounter.has(letter) || letterCounter.get(letter)! < needCounter.get(letter)!)
        return 0
    }
    return res
  }

  // 二进制枚举子集 从全选开始
  for (let i = 2 ** len - 1; i >= 0; i--) {
    const flagArray = i.toString(2).padStart(len, '0').split('')
    res = Math.max(res, countStore(flagArray))
  }

  return res
}

console.log(
  maxScoreWords(
    ['xxxz', 'ax', 'bx', 'cx'],
    ['z', 'a', 'b', 'c', 'x', 'x', 'x'],
    [4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 10]
  )
)
// 输出：27

export {}
