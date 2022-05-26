import { bisectRight } from '../../9_排序和搜索/二分/bisect'

/**
 * @param {string[]} queries
 * @param {string[]} words
 * @return {number[]}
 * 函数 f(s)，统计 s  中（按字典序比较）最小字母的出现频次
 * 对于每次查询 queries[i] ，
 * 需统计 words 中满足 f(queries[i]) < f(W) 的 词的数目 ，W 表示词汇表 words 中的每个词。
 * @summary
 * 离线排序
 */
function numSmallerByFrequency(queries: string[], words: string[]): number[] {
  const res = Array<number>(queries.length).fill(0)
  const freqs = words.map(getLowestFreq).sort((a, b) => a - b)

  for (let i = 0; i < queries.length; i++) {
    const query = queries[i]
    const queryFreq = getLowestFreq(query)
    // f(queries[i]) < f(W) 的 词的数目
    res[i] = words.length - bisectRight(freqs, queryFreq)
  }

  return res

  function getLowestFreq(word: string): number {
    let smallestChar = word.codePointAt(0)!
    let count = 1

    for (let i = 1; i < word.length; i++) {
      const curChar = word[i]
      if (curChar.codePointAt(0)! < smallestChar) {
        smallestChar = curChar.codePointAt(0)!
        count = 1
      } else if (curChar.codePointAt(0)! === smallestChar) {
        count++
      }
    }

    return count
  }
}

console.log(numSmallerByFrequency(['bbb', 'cc'], ['a', 'aa', 'aaa', 'aaaa']))
// 输出：[1,2]
// 解释：第一个查询 f("bbb") < f("aaaa")，第二个查询 f("aaa") 和 f("aaaa") 都 > f("cc")。

export {}
