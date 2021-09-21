/**
 * @param {string} num  num.length <= 1000
 * @param {string[]} words  words.length <= 500
 * @return {string[]}
 * 返回匹配单词的列表
 */
const getValidT9Words = function (num: string, words: string[]): string[] {
  const res: string[] = []
  const mapper = [2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 9, 9, 9, 9]
  for (const word of words) {
    const nums: number[] = []
    for (const char of word) {
      nums.push(mapper[char.codePointAt(0)! - 97])
    }
    nums.join('') === num && res.push(word)
  }
  return res
}

console.log(getValidT9Words('8733', ['tree', 'used']))
