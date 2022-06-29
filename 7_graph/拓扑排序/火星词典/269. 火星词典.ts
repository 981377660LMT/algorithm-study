import { topoSort2 } from '../topoSort'

/**
 * @param {string[]} words
 * @return {string}
 * 如果 s 中的字母在这门外星语言的字母顺序中位于 t 中字母之前，那么 s 的字典顺序小于 t
 * 如果前面 min(s.length, t.length) 字母都相同，那么 s.length < t.length 时，s 的字典顺序也小于 t
 * 请你根据该词典还原出此语言中已知的字母顺序，并 按字母递增顺序 排列
 * 若不存在合法字母顺序，返回 ""
 * @summary
 * 先每两个比较，找出单词的排序序列，然后装入字典，最后拓扑排序
 */
function alienOrder(words: string[]): string {
  if (!words.length) return ''
  const charSet = new Set<string>(words.join(''))

  // 对每个key，value必须在key之后(value入度加1)
  const adjMap = new Map<string, Set<string>>()
  for (let i = 0; i < words.length - 1; i++) {
    const [preWord, nextWord] = [words[i], words[i + 1]]
    const [preLen, nextLen] = [preWord.length, nextWord.length]

    // 判断输入是不是合法的。 比如 'wrd' 应该在'wr' 后面
    if (preLen > nextLen && preWord.slice(0, nextLen) === nextWord) return ''
    for (let k = 0; k < Math.min(preLen, nextLen); k++) {
      if (preWord[k] !== nextWord[k]) {
        !adjMap.has(preWord[k]) && adjMap.set(preWord[k], new Set())
        adjMap.get(preWord[k])!.add(nextWord[k])
        break // 只能判断一对
      }
    }
  }

  return topoSort2(adjMap, charSet).join('')
}

console.log(alienOrder(['wrt', 'wrf', 'er', 'ett', 'rftt']))
// 输出："wertf"

export {}
