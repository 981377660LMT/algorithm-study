/**
 * @param {string[]} words
 * @return {string[]}
 * 只返回可以使用在 美式键盘 同一行的字母打印出来的单词
 */
const findWords = function (words: string[]): string[] {
  return words.filter(w => {
    // remove word from array if it fails matching all three rows
    if (!/^[qwertyuiop]*$/i.test(w) && !/^[asdfghjkl]*$/i.test(w) && !/^[zxcvbnm]*$/i.test(w))
      return false

    return true
  })
}
console.log(findWords(['Hello', 'Alaska', 'Dad', 'Peace']))

export default 1
