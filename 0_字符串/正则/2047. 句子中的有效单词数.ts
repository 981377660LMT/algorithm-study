function countValidWords(sentence: string): number {
  return sentence
    .trim()
    .split(/\s+/)
    .filter(str => {
      return str.match(/^([a-z]+(-[a-z]+)?)?[!\.,]?$/)
    }).length
}
// 如果一个 token 同时满足下述条件，则认为这个 token 是一个有效单词：

// 仅由小写字母、连字符和/或标点（不含数字）。
// 至多一个 连字符 '-' 。如果存在，连字符两侧应当都存在小写字母（"a-b" 是一个有效单词，但 "-ab" 和 "ab-" 不是有效单词）。
// 至多一个 标点符号。如果存在，标点符号应当位于 token 的 末尾 。
