// 所有 字母日志 都排在 数字日志 之前。
// 字母日志 在内容不同时，忽略标识符后，按内容字母顺序排序；在内容相同时，按标识符排序。
function reorderLogFiles(logs: string[]): string[] {
  const comparator = (a: string, b: string): number => {
    const [[mark1, ...content1], [mark2, ...content2]] = [a.split(' '), b.split(' ')]
    const [isDigit1, isDigit2] = [
      Number.isInteger(Number(content1[0])),
      Number.isInteger(Number(content2[0])),
    ]

    return isDigit1 && isDigit2
      ? 0
      : Number(isDigit1) - Number(isDigit2) ||
          content1.join(' ').localeCompare(content2.join(' ')) ||
          mark1.localeCompare(mark2)
  }

  return logs.sort(comparator)
}

console.log(
  reorderLogFiles(['dig1 8 1 5 1', 'let1 art can', 'dig2 3 6', 'let2 own kit dig', 'let3 art zero'])
)
