// 截断指定长度的字符串，尽可能尊重空格。
truncateStringAtWhitespace('short', 10) // 'short'
truncateStringAtWhitespace('not so short', 10) // 'not so...'
truncateStringAtWhitespace('trying a thing', 10) // 'trying...'
truncateStringAtWhitespace('javascripting', 10) // 'javascr...'

// 加上ending后不超过limit
// 单词要尽量完整
function truncateStringAtWhitespace(str: string, limit: number, ending = '...') {
  if (str.length <= limit) return str
  // 在允许的范围内找到最后一个空格
  const lastSpace = str.slice(0, limit - ending.length).lastIndexOf(' ')
  return str.slice(0, lastSpace > 0 ? lastSpace : limit - ending.length) + ending
}

console.dir([
  truncateStringAtWhitespace('short', 10),
  truncateStringAtWhitespace('not so short', 10),
  truncateStringAtWhitespace('trying a thing', 10),
  truncateStringAtWhitespace('javascripting', 10),
])

const truncateString = (str: string, num: number) =>
  str.length > num ? str.slice(0, num > 3 ? num - 3 : num) + '...' : str
truncateString('boomerang', 7) // 'boom...'
