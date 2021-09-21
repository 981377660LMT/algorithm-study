/**
 * 
 * @param words 1 <= words.length <= 5000
 * 给定一组 互不相同 的单词， 找出所有 不同 的索引对 (i, j)，
   使得列表中的两个单词， words[i] + words[j] ，可拼接成回文串。
   @summary
   对于示例1的一个单词 lls，如何找其他单词把它拼成回文串。
   两种方式，一种是在左边加，比如有个单词s，那么可以拼车 slls；一种右边加，单词为sssll，拼成 llssssll
   由上面简单想法，可以把单词分成两部分考虑
 */
function palindromePairs(words: string[]): number[][] {
  const reversed = (str: string) => str.split('').reverse().join('')
  const isPalindrome = (str: string) => reversed(str) === str

  const res: number[][] = []
  const wordToIndexMap = new Map<string, number>()
  const palindromStr: number[] = []

  for (const [index, word] of words.entries()) {
    wordToIndexMap.set(reversed(word), index)
    if (isPalindrome(word)) palindromStr.push(index)
  }

  for (let i = 0; i < words.length; i++) {
    const word = words[i]

    if (word) {
      for (let j = 0; j < words.length - 1; j++) {
        const [left, right] = [word.slice(0, j), word.slice(j)]
        console.log(left, right)
        // 是否存在在单词左边加 使得成为回文串
        if (isPalindrome(left) && wordToIndexMap.has(right) && i !== wordToIndexMap.get(right)!) {
          res.push([wordToIndexMap.get(right)!, i])
        }

        // 是否存在在单词右边加 使得成为回文串
        if (isPalindrome(right) && wordToIndexMap.has(left) && i !== wordToIndexMap.get(left)!) {
          res.push([i, wordToIndexMap.get(left)!])
        }
      }
    } else {
      // 空字符串
      for (const index of palindromStr) {
        if (index !== i) {
          res.push([index, i])
          res.push([i, index])
        }
      }
    }
  }

  return res
}

console.log(palindromePairs(['abcd', 'dcba', 'lls', 's', 'sssll']))
// 输出：[[0,1],[1,0],[3,2],[2,4]]
// 解释：可拼接成的回文串为 ["dcbaabcd","abcddcba","slls","llssssll"]

export {}
