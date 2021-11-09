/**
 * @param {string} s  1 <= s.length <= 1000  1 <= dictionary.length <= 1000
 * @param {string[]} dictionary
 * @return {string}
 * 找出并返回字典中最长的字符串，该字符串可以通过删除 s 中的某些字符得到。
 * 如果答案不止一个，返回长度最长且字典序最小的字符串
 */
// 需要多次调用时，使用数组+哈希 map (单词很多时，空间开销极大)
// 对每个word,**倒序**记录从该位置开始往后s中的字符第一次出现的位置
const findLongestWord = function (s: string, dictionary: string[]): string {
  dictionary.sort((a, b) => b.localeCompare(a))
  const n = s.length
  const chars = new Set(s)
  const maps = Array.from<number, Map<string, number>>({ length: n + 1 }, () => new Map())
  // 初始化 n表示不存在
  for (const char of chars) {
    maps[n].set(char, n)
  }

  for (let i = n - 1; i >= 0; i--) {
    for (const char of s) {
      if (char === s[i]) maps[i].set(char, i)
      else maps[i].set(char, maps[i + 1].get(char)!)
    }
  }

  let res = ''
  // 统计每个单词的与target的公共长度
  for (const word of dictionary) {
    let index = 0
    let count = 0
    for (const char of word) {
      // 如果下一个字符的下标为 n，则表示该字符不存在
      if (!maps[index].has(char) || maps[index].get(char) === n) break
      count++
      index = maps[index].get(char)! + 1
    }

    // 注意这里：count === word.length判断是否到了最后一个字符
    // 否则该字符串不可以通过删除 s 中的某些字符得到
    if (count === word.length && count >= res.length) res = word
  }

  return res
}

// 只用调用一次时
const findLongestWord2 = function (s: string, dictionary: string[]): string {
  let res = ''

  for (const word of dictionary) {
    let hit = 0
    for (let i = 0; i < s.length; i++) {
      if (s[i] === word[hit]) hit++
      if (hit === word.length) break
    }
    // 可以匹配
    if (hit === word.length) {
      // 等于要判断字典序在前
      if (word.length > res.length || (word.length === res.length && word < res)) res = word
    }
  }

  return res
}
console.log(findLongestWord2('abpcplea', ['ale', 'apple', 'monkey', 'plea']))

export default 1
