/**
 * @param {number[]} data
 * @return {boolean}
 * @description
 * UTF-8 中的一个字符可能的长度为 1 到 4 字节
 * 对于 1 字节的字符，字节的第一位设为 0 ，后面 7 位为这个符号的 unicode 码。
 * 对于 n 字节的字符 (n > 1)，第一个字节的前 n 位都设为1，第 n+1 位设为 0 ，后面字节的前两位一律设为 10 。
 * @summary
 * 用二进制整数：0bxxxxxxxx，表示二进制格式，避免位操作出错。
 */
const validUtf8 = function (data: number[]): boolean {
  const n = data.length
  let i = 0

  while (i < n) {
    const num = data[i]
    let shouldCheck: number
    if ((num & 0b10000000) === 0) {
      shouldCheck = 0
    } else if ((num & 0b11100000) === 0b11000000) {
      shouldCheck = 1
    } else if ((num & 0b11110000) === 0b11100000) {
      shouldCheck = 2
    } else if ((num & 0b11111000) === 0b11110000) {
      shouldCheck = 3
    } else {
      return false
    }

    i++
    while (i < n && shouldCheck) {
      if ((data[i] & 0b11000000) !== 0b10000000) return false
      i++
      shouldCheck--
    }

    if (shouldCheck) return false
  }

  return true
}

console.log(validUtf8([197, 130, 1]))
// 表示 8 位的序列: 11000101 10000010 00000001.
export default 1
