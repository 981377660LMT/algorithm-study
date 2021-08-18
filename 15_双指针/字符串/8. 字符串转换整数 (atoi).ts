/**
 * @param {string} s
 * @return {number}
 * 使其能将字符串转换成一个 32 位有符号整数
 */
const myAtoi = function (s: string): number {
  const pattern = /^[-+]?\d+/g
  const match = s.trim().match(pattern)
  return match ? Math.max(Math.min(Number(match[0]), 2 ** 31 - 1), -(2 ** 31)) : 0
}

console.log(parseInt('4193 with words'))
console.log(parseInt('words and 987') || 0)
console.log(myAtoi('4193 with words'))
