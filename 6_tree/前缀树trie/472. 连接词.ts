/**
 * @param {string[]} words
 * @return {string[]}
 * @description 返回 words 中的所有 连接词 。
 * 连接词 的定义为：一个字符串完全是由至少两个给定数组中的单词组成的。
 * @summary 可以用trie的思路做
 */
const findAllConcatenatedWordsInADict = function (words: string[]): string[] {
  const res: string[] = []
  const store = new Set(words)

  // 注意s的prefix是s.slice(0,i+1),suffix是s.slice(i+1)
  const isConcat = (s: string): boolean => {
    if (store.has(s)) return true
    for (let i = 0; i < s.length; i++) {
      const prefix = s.slice(0, i + 1)
      if (store.has(prefix)) {
        const suffix = s.slice(i + 1)
        if (isConcat(suffix)) {
          store.add(s)
          return true
        }
      }
    }

    return false
  }

  for (const word of words) {
    store.delete(word)
    if (isConcat(word)) res.push(word)
    store.add(word)
  }

  return res
}

console.log(
  findAllConcatenatedWordsInADict([
    'cat',
    'cats',
    'catsdogcats',
    'dog',
    'dogcatsdog',
    'hippopotamuses',
    'rat',
    'ratcatdogcat',
  ])
)

// ["catsdogcats","dogcatsdog","ratcatdogcat"]

export {}
