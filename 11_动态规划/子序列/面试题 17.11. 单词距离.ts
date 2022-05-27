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
function findClosest(words: string[], word1: string, word2: string): number {
  const indexMap = new Map<string, number[]>()
  for (const [index, word] of words.entries()) {
    !indexMap.has(word) && indexMap.set(word, [])
    indexMap.get(word)!.push(index)
  }

  return minDiff(indexMap.get(word1) ?? [], indexMap.get(word2) ?? [])

  // 两个有序数组选数差的绝对值最小
  function minDiff(nums1: number[], nums2: number[]): number {
    let res = Infinity
    let i = 0
    let j = 0
    while (i < nums1.length && j < nums2.length) {
      res = Math.min(res, Math.abs(nums1[i] - nums2[j]))
      // 小的向后移
      nums1[i] < nums2[j] ? i++ : j++
    }
    return res
  }
}

// 如果不需要调用多次，一般的双指针即可，遇到即更新
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
