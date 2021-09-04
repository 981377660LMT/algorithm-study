/**
 * @param {number} n
 * @return {number}
 * 在无限的整数序列 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ...中找到第 n 位数字。
 * 1 - 9 : 9 digits
   10 - 99 : 90 * 2 digits
   100 - 999: 900 * 3 digits
 */
const findNthDigit = function (n: number): number {
  let num = 9
  let digit = 1
  while (n > num * digit) {
    n -= num * digit
    num *= 10
    digit++
  }
  const base = 10 ** (digit - 1)
  const offset = Math.ceil(n / digit) - 1
  const indexInDigit = (n - 1) % digit
  console.log(base, offset, indexInDigit)
  return parseInt((base + offset).toString()[indexInDigit])
}

console.log(findNthDigit(11))
// 输出：0
// 解释：第 11 位数字在序列 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, ... 里是 0 ，它是 10 的一部分。
