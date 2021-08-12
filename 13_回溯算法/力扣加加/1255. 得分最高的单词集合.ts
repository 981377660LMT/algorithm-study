/**
 * @param {string[]} words 1 <= words.length <= 14
 * @param {string[]} letters
 * @param {number[]} score
 * @return {number}
 * @description 请你帮忙计算玩家在单词拼写游戏中所能获得的「最高得分」：
 * 能够由 letters 里的字母拼写出的 任意 属于 words 单词子集中，分数最高的单词集合的得分。
 * 单词表 words 中每个单词只能计分（使用）一次。
 * 可以只使用字母表 letters 中的部分字母，但是每个字母最多被使用一次。
 */
const maxScoreWords = function (words: string[], letters: string[], score: number[]) {}

console.log(
  maxScoreWords(
    ['xxxz', 'ax', 'bx', 'cx'],
    ['z', 'a', 'b', 'c', 'x', 'x', 'x'],
    [4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 10]
  )
)
// 输出：27

export {}
