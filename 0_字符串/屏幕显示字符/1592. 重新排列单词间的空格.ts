// 返回 重新排列空格后的字符串 。
// 请你重新排列空格，使每对相邻单词之间的空格数目都 相等 ，并尽可能 最大化 该数目。
// 如果不能重新平均分配所有空格，请 将多余的空格放置在字符串末尾 ，
function reorderSpaces(text: string): string {
  const words = text.split(/\s+/).filter(word => word.length > 0)
  const intervalCount = words.length - 1
  const wordsLength = words.reduce((pre, cur) => pre + cur.length, 0)
  const spaceLength = text.length - wordsLength

  if (intervalCount === 0) return words[0] + ' '.repeat(spaceLength)

  const [div, mod] = [~~(spaceLength / intervalCount), spaceLength % intervalCount]
  const stringBuilder: string[] = []

  for (let i = 0; i < words.length - 1; i++) {
    stringBuilder.push(words[i])
    stringBuilder.push(' '.repeat(div))
  }

  stringBuilder.push(words[words.length - 1])
  stringBuilder.push(' '.repeat(mod))

  return stringBuilder.join('')
}

// console.log(reorderSpaces('  this   is  a sentence '))
// 输出："this   is   a   sentence"
// 解释：总共有 9 个空格和 4 个单词。可以将 9 个空格平均分配到相邻单词之间，相邻单词间空格数为：9 / (4-1) = 3 个。
console.log(reorderSpaces('  hello'))
