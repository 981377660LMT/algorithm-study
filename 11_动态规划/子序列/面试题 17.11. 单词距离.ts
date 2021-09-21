/**
 * @param {string[]} words
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * 有个内含单词的超大文本文件，给定任意两个单词，找出在这个文件中这两个单词的最短距离(相隔单词数)。
 * 如果寻找过程在这个文件中会重复多次，而每次寻找的单词不同，你能对此优化吗?
 * 哈希表存出现的 index
 * ```
 * 参考 11_动态规划\子序列\note.md
 * ```
 */
const findClosest = function (words: string[], word1: string, word2: string): number {
  // 两个有序数组选数差的绝对值最小
  const minDiff = (arr1: number[], arr2: number[]): number => {
    let res = Infinity
    let l1 = 0
    let l2 = 0
    while (l1 < arr1.length && l2 < arr2.length) {
      res = Math.min(res, Math.abs(arr1[l1] - arr2[l2]))
      // 小的向后移
      arr1[l1] < arr2[l2] ? l1++ : l2++
    }
    return res
  }
  const record = new Map<string, number[]>()
  for (let i = 0; i < words.length; i++) {
    const word = words[i]
    !record.has(word) && record.set(word, [])
    record.get(word)!.push(i)
  }

  return minDiff(record.get(word1)!, record.get(word2)!)
}

// 如果不需要调用多次，一般的双指针即可
const findClosest2 = function (words: string[], word1: string, word2: string): number {
  let i = -Infinity
  let j = Infinity
  let res = Infinity
  for (const [index, word] of words.entries()) {
    word1 === word && (i = index)
    word2 === word && (j = index)
    res = Math.min(res, Math.abs(i - j))
  }

  return res
}
console.log(
  findClosest2(
    ['I', 'am', 'a', 'student', 'from', 'a', 'university', 'in', 'a', 'city'],
    'a',
    'student'
  )
)
