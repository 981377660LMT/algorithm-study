import { Trie } from './实现trie/1_实现trie'

// 若其中有多个可行的答案，则返回答案中字典序最小的单词。
// 给出一个字符串数组words组成的一本英语词典。
// 从中找出最长的一个单词，该单词是由words词典中其他单词逐步添加一个字母组成。
// const longestWord = (words: string[]) => {
//   if (words.length === 0) return ''
//   const trie = new Trie()
//   let res = ''
//   // 字典序,即localCompare的顺序排序
//   words.sort()
//   console.log(words)
//   words.forEach(word => {
//     // 注意这个条件写法
//     if (word.length === 1 || trie.search(word.slice(0, -1))) {
//       trie.insert(word)
//       if (word.length > res.length) res = word
//     }
//   })

//   return res
// }

// 集合比trie快一些
const longestWord = (words: string[]) => {
  if (words.length === 0) return ''
  const set = new Set()
  let res = ''
  // 字典序,即localCompare的顺序排序
  words.sort()
  console.log(words)
  words.forEach(word => {
    // 注意这个条件写法
    if (word.length === 1 || set.has(word.slice(0, -1))) {
      set.add(word)
      if (word.length > res.length) res = word
    }
  })

  return res
}

console.dir(
  longestWord([
    't',
    'ti',
    'tig',
    'tige',
    'tiger',
    'e',
    'en',
    'eng',
    'engl',
    'engli',
    'englis',
    'english',
    'h',
    'hi',
    'his',
    'hist',
    'histo',
    'histor',
    'history',
  ]),
  { depth: null }
)

export {}
