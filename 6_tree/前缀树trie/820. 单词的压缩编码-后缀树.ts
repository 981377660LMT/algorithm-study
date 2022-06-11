import { Trie } from './Trie'

/**
 * @param {string[]} words
 * @return {number}
 * @description 给你一个单词数组 words ，返回成功对 words 进行编码的最小助记字符串 s 的长度 。
 * 助记字符串 s 以 '#' 字符结尾
 * @summary 将每个word倒序会发现规律
 * 使用前缀树 + 倒序插入的形式来模拟后缀树
 */
function minimumLengthEncoding(words: string[]): number {
  let res = 0
  const trie = new Trie()

  // 按长度倒序插入
  words = words.sort((a, b) => b.length - a.length).map(word => word.split('').reverse().join(''))

  words.forEach(w => {
    if (trie.countPre(w) === 0) {
      // console.log(word)
      trie.insert(w)
      // 多了一个分支#和分支的长度
      res += w.length + 1
    }
  })

  return res
}

console.log(minimumLengthEncoding(['time', 'me', 'bell']))
// 可以看到一个#结尾前面的词表示trie的一个分支
// s = "time#bell#" 和 indices = [0, 2, 5] 。
// 长度为10 每个trie分支的长度+trie分支数
export {}
