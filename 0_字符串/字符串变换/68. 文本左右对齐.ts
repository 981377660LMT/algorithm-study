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
function fullJustify(words: string[], maxWidth: number): string[] {
  const res: string[] = []
  const n = words.length
  let ptr = 0
  while (ptr < n) {
    let lineLen = 0
    let end = ptr
    // 计算当前行可以放多少个单词
    while (end < n && lineLen + words[end].length + (end - ptr) <= maxWidth) {
      lineLen += words[end].length
      end++
    }

    const numWords = end - ptr
    let line = ''

    // 如果是最后一行或只有一个单词，左对齐
    if (end === n || numWords === 1) {
      for (let k = ptr; k < end; k++) {
        line += `${words[k]} `
      }
      line = line.trimEnd()
      line += ' '.repeat(maxWidth - line.length)
    } else {
      const totalSpaces = maxWidth - lineLen
      const gaps = numWords - 1
      const spaceBetween = Math.floor(totalSpaces / gaps)
      const extraSpaces = totalSpaces % gaps
      for (let k = ptr; k < end - 1; k++) {
        line += words[k]
        // 前 extraSpaces 个间隔多加一格
        const spaces = spaceBetween + (k - ptr < extraSpaces ? 1 : 0)
        line += ' '.repeat(spaces)
      }
      line += words[end - 1] // 最后一个单词后不加空格
    }

    res.push(line)
    ptr = end
  }

  return res
}

console.log(fullJustify(['This', 'is', 'an', 'example', 'of', 'text', 'justification.'], 16))

export default 1
