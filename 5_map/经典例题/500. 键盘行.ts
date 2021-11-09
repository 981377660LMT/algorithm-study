/**
 * @param {string[]} words
 * @return {string[]}
 * 只返回可以使用在 美式键盘 同一行的字母打印出来的单词
 */
const findWords = function (words: string[]): string[] {
  const lines = [new Set('qwertyuiop'), new Set('asdfghjkl'), new Set('zxcvbnm')]
  return words.filter(word => lines.some(lineSet => isSubSet(new Set(word.toLowerCase()), lineSet)))

  function isSubSet<T>(curSet: Set<T>, targetSet: Set<T>): boolean {
    if (curSet.size > targetSet.size) return false
    console.log(curSet, targetSet)
    for (const value of curSet) {
      if (!targetSet.has(value)) return false
    }

    return true
  }
}
console.log(findWords(['Hello', 'Alaska', 'Dad', 'Peace']))

export default 1
