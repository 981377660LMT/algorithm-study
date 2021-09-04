/**
 * @param {string[]} words
 * @param {number} maxWidth
 * @return {string[]}
 * 尽可能多地往每行中放置单词。必要时可用空格 ' ' 填充，使得每行恰好有 maxWidth 个字符
 * 要求尽可能均匀分配单词间的空格数量。如果某一行单词间的空格不能均匀分配，则左侧放置的空格数要多于右侧的空格数。

 * @description 编辑器开发
 * 第一，某一行只有一个单词，这个单词需要左对齐；
 * 第二，对于最后一行，不论有多少个单词，都应该左对齐，并且单词与单词之间只插入一个空格。
 */
const fullJustify = function (words: string[], maxWidth: number): string[] {
  const res: string[] = []
  let lineWords: string[] = [] // 当前行的单词列表
  let charLen = 0 // 当前行单词的字符数(不含空格)

  for (const word of words) {
    // 单词数+单词间空格数超出本行允许的最大字符数 计算需要填充空格
    if (charLen + word.length + lineWords.length > maxWidth) {
      // 如果本行只有一个单词，需要左对齐
      if (lineWords.length === 1) {
        res.push(lineWords[0] + ' '.repeat(maxWidth - lineWords[0].length))
      } else {
        const tmp: string[] = []
        const blank = maxWidth - charLen
        const spaceInterval = ~~(blank / (lineWords.length - 1))
        const spaceSpare = blank % (lineWords.length - 1)
        for (let i = 0; i < lineWords.length; i++) {
          tmp.push(lineWords[i])
          // 除了最后一个单词外，其他单词后面都要跟特定长度的空格 spare空间分给前面的几个位置
          if (i < lineWords.length - 1)
            tmp.push(' '.repeat(spaceInterval + (i < spaceSpare ? 1 : 0)))
        }
        res.push(tmp.join(''))
      }

      // 更新新的行与长度
      lineWords = [word]
      charLen = word.length
    } else {
      lineWords.push(word)
      charLen += word.length
    }
  }

  // 剩下的：最后一行需要左对齐，类似处理单个字符的情况
  if (lineWords.length) {
    res.push(lineWords.join(' ') + ' '.repeat(maxWidth - (charLen + lineWords.length - 1)))
  }

  return res
}

console.log(fullJustify(['This', 'is', 'an', 'example', 'of', 'text', 'justification.'], 16))

export default 1
