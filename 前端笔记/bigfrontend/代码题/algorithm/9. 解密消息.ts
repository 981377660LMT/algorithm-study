// 从左上开始，向右下前进
// 无法前进的时候，向右上前进
// 无法前进的时候，向右下前进
// 2和3的重复
// 无法前进的时候，经过的字符就就是隐藏信息

/**
 * @param {string[][]} message
 * @return {string}
 */
function decode(message: string[][]): string {
  if (!message.length || !message[0].length) return ''
  const m = message.length
  const n = message[0].length

  let x = 0
  let y = 0
  let flag = -1
  const sb: string[] = []
  while (sb.length < n) {
    sb.push(message[x][y])
    if (x === 0 || x === m - 1) flag *= -1
    x += flag
    y++
  }

  return sb.join('')
}

console.log(
  decode([
    ['I', 'B', 'C', 'A', 'L', 'K', 'A'],
    ['D', 'R', 'C', 'A', 'L', 'K', 'A'],
    ['G', 'H', 'C', 'A', 'L', 'K', 'A'],
  ])
)

// 隐藏消息是IRCALKA
