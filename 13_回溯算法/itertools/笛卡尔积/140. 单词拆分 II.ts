/**
 * @param {string} s
 * @param {string[]} wordDict
 * @return {string[]}
 * @description 在字符串中增加空格来构建一个句子，使得句子中所有的单词都在词典中。
 * !返回所有这些可能的句子。
 * @summary 普通的回溯会超时(s 的长度是 151)
 */
const wordBreak = function (s: string, wordDict: string[]): string[] {
  const store = new Set(wordDict)
  const res: string[] = []

  bt([], s)

  return res

  function bt(path: string[], remain: string) {
    if (remain.length === 0) {
      res.push(path.join(' '))
      return
    }

    for (let i = 0; i < remain.length; i++) {
      const next = remain.slice(0, i + 1)
      if (store.has(next)) {
        path.push(next)
        bt(path, remain.slice(i + 1))
        path.pop()
      }
    }
  }
}

console.log(wordBreak('pineapplepenapple', ['apple', 'pen', 'applepen', 'pine', 'pineapple']))
// [
//   "pine apple pen apple",
//   "pineapple pen apple",
//   "pine applepen apple"
// ]
