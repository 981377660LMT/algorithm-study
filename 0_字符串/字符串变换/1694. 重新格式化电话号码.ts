// 将数组从左到右 每 3 个一组 分块，直到 剩下 4 个或更少数字。剩下的数字将按下述规定再分块：
// 2 个数字：单个含 2 个数字的块。
// 3 个数字：单个含 3 个数字的块。
// 4 个数字：两个分别含 2 个数字的块。

function reformatNumber(number: string): string {
  const digits = number.replace(/\D/g, '')
  const sb: string[] = []

  let cursor = 0
  while (cursor + 4 < digits.length) {
    sb.push(digits.slice(cursor, cursor + 3))
    sb.push('-')
    cursor += 3
  }

  if (cursor + 2 === digits.length) {
    sb.push(digits.slice(cursor, cursor + 2))
  } else if (cursor + 3 === digits.length) {
    sb.push(digits.slice(cursor, cursor + 3))
  } else {
    sb.push(`${digits.slice(cursor, cursor + 2)}-${digits.slice(cursor + 2, cursor + 4)}`)
  }

  return sb.join('')
}

console.log(reformatNumber('1-23-45 6'))
