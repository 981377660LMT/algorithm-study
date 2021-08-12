import { Trie } from './1_实现trie'

/**
 * @param {string[]} words
 * @return {number}
 * @description 给你一个单词数组 words ，返回成功对 words 进行编码的最小助记字符串 s 的长度 。
 * 助记字符串 s 以 '#' 字符结尾
 * @summary 将每个word倒序会发现规律
 * 使用前缀树 + 倒序插入的形式来模拟后缀树
 */
const minimumLengthEncoding = function (words: string[]): number {
  let res = 0
  const trie = new Trie()

  words = words.sort((a, b) => b.length - a.length).map(word => word.split('').reverse().join(''))

  words.forEach(word => {
    if (!trie.startsWith(word)) {
      trie.insert(word)
      // 多了一个分支#和分支的长度
      res += word.length + 1
    }
  })

  return res
}

console.log(minimumLengthEncoding(['time', 'me', 'bell']))
// 可以看到一个#结尾前面的词表示trie的一个分支
// s = "time#bell#" 和 indices = [0, 2, 5] 。
// 长度为10 每个trie分支的长度+trie分支数
