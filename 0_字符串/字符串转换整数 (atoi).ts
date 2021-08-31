/**
 * @param {string} s
 * @return {number}
 */
const myAtoi1 = function (s: string): number {
  const pattern = /^[-+]?\d+/g
  const match = s.trim().match(pattern)
  return match ? Math.max(Math.min(Number(match[0]), 2 ** 31 - 1), -(2 ** 31)) : 0
}

const myAtoi2 = function (s: string): number {
  // parseInt
  const num = parseInt(s)

  // 判断 parseInt 的结果是否为 NaN，是则返回 0
  if (isNaN(num)) {
    return 0
  } else if (num < Math.pow(-2, 31) || num > Math.pow(2, 31) - 1) {
    // 超出
    return num < Math.pow(-2, 31) ? Math.pow(-2, 31) : Math.pow(2, 31) - 1
  } else {
    return num
  }
}
console.log(myAtoi2('4193 with words'))
console.log(myAtoi2('words and 987'))
export {}

// parseInt(string, radix)
// 解析一个字符串并返回指定基数的十进制整数，
// radix 是2-36之间的整数，表示被解析字符串的基数。
// 如果 parseInt遇到的字符不是指定 radix参数中的数字，
// 它将忽略该字符以及所有后续字符，并返回到该点为止已解析的整数值。
// console.log(parseInt('1s00010'))  1

// parseInt 将数字截断为整数值。 允许前导和尾随空格。
// parseInt 可以理解两个符号。+ 表示正数，- 表示负数（从ECMAScript 1开始）。
// 它是在去掉空格后作为解析的初始步骤进行的。
// 如果没有找到符号，算法将进入下一步；否则，它将删除符号，并对字符串的其余部分进行数字解析。

// 如果第一个字符不能转换为数字，parseInt会返回 NaN。
// console.log(parseInt('s00010'))  NaN

// 为了算术的目的，NaN 值不能作为任何 radix 的数字。
// 你可以调用isNaN函数来确定parseInt的结果是否为 NaN。
// 如果将NaN传递给算术运算，则运算结果也将是 NaN。 -- 来自MDN
