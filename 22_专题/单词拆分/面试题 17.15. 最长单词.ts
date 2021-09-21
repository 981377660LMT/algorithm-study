/**
 * @param {string[]} words
 * @return {string}
 * 找出其中的最长单词，且该单词由这组单词中的其他单词组合而成(可重复使用)。
 * 若有多个长度相同的结果，返回其中字典序最小的一项，若没有符合要求的单词则返回空字符串。
 */
const longestWord = function (words: string[]): string {
  // 最长单词若有多个长度相同的结果，返回其中字典序最小的一项
  words.sort((a, b) => b.length - a.length || a.localeCompare(b))

  const dfs = (remain: string, alternativeWords: string[]): boolean => {
    if (remain === '') return true
    for (const next of alternativeWords) {
      const len = next.length
      if (remain.slice(0, len) === next) {
        if (dfs(remain.slice(len), alternativeWords)) return true
      }
    }
    return false
  }

  for (let i = 0; i < words.length; i++) {
    if (dfs(words[i], words.slice(i + 1))) return words[i]
  }

  return ''
}

// console.log(longestWord(['cat', 'banana', 'dog', 'nana', 'walk', 'walker', 'dogwalker']))
console.log(
  longestWord(['bbbbb', 'bbbbbb', 'bb', 'bb', 'bbb', 'bbbbb', 'bbbbbbbb', 'bbb', 'bbbbbbb'])
)
